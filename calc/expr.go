package calc

import (
	"github.com/berquerant/dql/arithmetic"
	"github.com/berquerant/dql/ast"
	"github.com/berquerant/dql/bit"
	"github.com/berquerant/dql/compare"
	"github.com/berquerant/dql/data"
	"github.com/berquerant/dql/env"
	"github.com/berquerant/dql/errors"
	"github.com/berquerant/dql/function"
	"github.com/berquerant/dql/logger"
)

type Calculator interface {
	Data(expr ast.Expr) (data.Data, error)
}

func New(
	comparer compare.Comparer,
	artCalculator arithmetic.Calculator,
	bitCalculator bit.Calculator,
	funcCaller function.Caller,
	env env.Map,
) Calculator {
	return &calculator{
		comparer:      comparer,
		artCalculator: artCalculator,
		bitCalculator: bitCalculator,
		funcCaller:    funcCaller,
		env:           env,
	}
}

type calculator struct {
	comparer          compare.Comparer
	artCalculator     arithmetic.Calculator
	bitCalculator     bit.Calculator
	funcCaller        function.Caller
	env               env.Map
	onAggregation     bool
	aggregationIndex  int
	aggregationTarget string
}

var (
	ErrUnknownExpr  = errors.New("unknown expr")
	ErrTypeMismatch = errors.New("type mismatch")
)

func (s *calculator) Data(expr ast.Expr) (data.Data, error) {
	data, err := s.data(expr)
	if err != nil {
		return nil, errors.Wrap(err, "cannot calculate %s", logger.JSON(expr))
	}
	return data, nil
}

func (s *calculator) data(expr ast.Expr) (data.Data, error) {
	data, err := func() (data.Data, error) {
		switch expr := expr.(type) {
		case *ast.OrExpr, *ast.AndExpr, *ast.XorExpr, *ast.NotExpr:
			return s.dataLogicalOperator(expr)
		case ast.BoolPrimary:
			return s.dataBoolPrimary(expr)
		case ast.Predicate:
			return s.dataPredicate(expr)
		case ast.BitExpr:
			return s.dataBitExpr(expr)
		case ast.SimpleExpr:
			return s.dataSimpleExpr(expr)
		default:
			return nil, ErrUnknownExpr
		}
	}()
	if err != nil {
		return nil, errors.Wrap(err, "cannot calculate %s", logger.JSON(expr))
	}
	return data, nil
}

func (s *calculator) dataSimpleExpr(expr ast.SimpleExpr) (data.Data, error) {
	switch expr := expr.(type) {
	case *ast.SimpleExprPrefixOp:
		return s.dataPrefixOp(expr.Op, expr.Expr)
	case *ast.SimpleExprLit:
		return s.dataLit(expr.Lit)
	case *ast.Ident:
		if s.onAggregation {
			return s.dataIdentAggregation(expr)
		}
		return s.dataIdentNormal(expr)
	case *ast.FunctionCall:
		f, exist := s.funcCaller.Func(expr.FunctionName.Value)
		if !exist {
			return nil, errors.Wrap(ErrUnknownExpr, "function call function %s not found", expr.FunctionName.Value)
		}
		if _, ok := f.(function.Aggregation); ok {
			if s.onAggregation {
				return nil, errors.Wrap(ErrUnknownExpr, "aggregation cannot be nested")
			}
			return s.dataFunctionCallAggregation(expr)
		}
		return s.dataFunctionCallNormal(expr)
	case *ast.SimpleExprExpr:
		return s.data(expr.Expr)
	default:
		return nil, errors.Wrap(ErrUnknownExpr, "simple expr %s", logger.JSON(expr))
	}
}

/* Aggregation expression must contain just 1 ident (column),
 * aggregations must not be nested.
 * Aggregation process are below:
 * 1. Set true to onAggregation.
 * 2. Calculate for each element in the ident value.
 * 3. Set false to onAggregation.
 * 4. Aggregate them.
 */

func (s *calculator) prepareAggregation(start bool) {
	s.onAggregation = start
	s.aggregationIndex = 0
	s.aggregationTarget = ""
}

func (s *calculator) dataIdentNormal(expr *ast.Ident) (data.Data, error) {
	if v, ok := s.env.Get(expr.Value); ok {
		switch v.Type() {
		case env.TypeData:
			return v.Data(), nil
		case env.TypeExpr:
			d, err := s.data(v.Expr())
			if err != nil {
				return nil, errors.Wrap(err, "failed to calc ident %s %s", expr.Value, logger.JSON(v))
			}
			s.env.Set(expr.Value, env.FromData(d))
			return d, nil
		default:
			return nil, errors.Wrap(ErrUnknownExpr, "unknown ident %s %s", expr.Value, logger.JSON(v))
		}
	}
	return nil, errors.Wrap(ErrUnknownExpr, "cannot find ident %s", expr.Value)
}

func (s *calculator) dataIdentAggregation(expr *ast.Ident) (data.Data, error) {
	v, ok := s.env.Get(expr.Value)
	if !ok {
		return nil, errors.Wrap(ErrUnknownExpr, "cannot find ident %s", expr.Value)
	}
	if v.Type() != env.TypeDataList {
		return nil, errors.Wrap(ErrTypeMismatch, "cannot find data list %s", expr.Value)
	}
	if s.aggregationTarget == "" {
		s.aggregationTarget = expr.Value
	} else if s.aggregationTarget != expr.Value {
		return nil, errors.Wrap(ErrUnknownExpr, "aggregation cannot depend on multiple columns %s %s",
			s.aggregationTarget, expr.Value)
	}
	list := v.DataList() // fetch rows
	if s.aggregationIndex < 0 || s.aggregationIndex >= len(list) {
		return nil, errors.Wrap(ErrUnknownExpr, "out of index %d in %s", s.aggregationIndex, logger.JSON(list))
	}
	t := list[s.aggregationIndex]
	s.aggregationIndex++ // next row
	if s.aggregationIndex == len(list) {
		// all elements are scanned, end aggregation
		s.onAggregation = false
	}
	return t, nil
}

func (s *calculator) dataFunctionCallAggregation(expr *ast.FunctionCall) (data.Data, error) {
	if len(expr.Arguments.Exprs) != 1 {
		return nil, errors.Wrap(ErrUnknownExpr,
			"number of aggregation function arguments must be 1 but got %d", len(expr.Arguments.Exprs))
	}
	s.prepareAggregation(true) // start aggregation
	defer s.prepareAggregation(false)
	args := []data.Data{}
	for s.onAggregation {
		v, err := s.data(expr.Arguments.Exprs[0])
		if err != nil {
			return nil, errors.Wrap(err, "failed to calc %s[%d] on aggregation %s",
				s.aggregationTarget, s.aggregationIndex, expr.FunctionName.Value,
			)
		}
		args = append(args, v)
	}
	r, err := s.funcCaller.Call(expr.FunctionName.Value, args...)
	if err != nil {
		return nil, errors.Wrap(err, "function call %s %s", logger.JSON(args), logger.JSON(expr))
	}
	return r, nil
}

func (s *calculator) dataFunctionCallNormal(expr *ast.FunctionCall) (data.Data, error) {
	args := make([]data.Data, len(expr.Arguments.Exprs))
	for i, a := range expr.Arguments.Exprs {
		v, err := s.data(a)
		if err != nil {
			return nil, errors.Wrap(err, "function call %s args[%d] %s", logger.JSON(expr), i, logger.JSON(a))
		}
		args[i] = v
	}
	r, err := s.funcCaller.Call(expr.FunctionName.Value, args...)
	if err != nil {
		return nil, errors.Wrap(err, "function call %s %s", logger.JSON(args), logger.JSON(expr))
	}
	return r, nil
}

func (s *calculator) dataLit(expr ast.Lit) (data.Data, error) {
	switch expr := expr.(type) {
	case *ast.IntLit:
		return data.FromInt(expr.Value), nil
	case *ast.FloatLit:
		return data.FromFloat(expr.Value), nil
	case *ast.StringLit:
		return data.FromString(expr.Value), nil
	default:
		return nil, errors.Wrap(ErrUnknownExpr, "literal %s", expr)
	}
}

func (s *calculator) dataPrefixOp(op ast.PrefixOperatorType, expr ast.SimpleExpr) (data.Data, error) {
	v, err := s.data(expr)
	if err != nil {
		return nil, errors.Wrap(err, "arg %s", logger.JSON(expr))
	}
	switch op {
	case ast.PreOpPlus:
		return v, nil
	case ast.PreOpMinus:
		switch v.Type() {
		case data.TypeInt:
			return data.FromInt(-v.Int()), nil
		case data.TypeFloat:
			return data.FromFloat(-v.Float()), nil
		default:
			return nil, errors.Wrap(ErrUnknownExpr, "prefix op minus arg %s", logger.JSON(v))
		}
	case ast.PreOpBitNot:
		r, err := s.bitCalculator.Not(v.Value())
		if err != nil {
			return nil, errors.Wrap(err, "prefix op bit not arg %s", logger.JSON(v))
		}
		return data.FromInt(r), nil
	case ast.PreOpNot:
		switch v.Type() {
		case data.TypeBool:
			return data.FromBool(!v.Bool()), nil
		default:
			return nil, errors.Wrap(ErrUnknownExpr, "prefix op not arg %s", logger.JSON(v))
		}
	default:
		return nil, errors.Wrap(ErrUnknownExpr, "prefix op unknown operation %s", op)
	}
}

func (s *calculator) dataBitExpr(expr ast.BitExpr) (data.Data, error) {
	var err error
	switch expr := expr.(type) {
	case *ast.BitExprSimpleExpr:
		return s.data(expr.Expr)
	case ast.BinaryOp:
		var (
			left, right data.Data
		)
		if left, err = s.Data(expr.LeftArg()); err != nil {
			return nil, errors.Wrap(err, "left arg %s", logger.JSON(expr.LeftArg()))
		}
		if right, err = s.Data(expr.RightArg()); err != nil {
			return nil, errors.Wrap(err, "right arg %s", logger.JSON(expr.RightArg()))
		}
		switch expr := expr.(type) {
		case *ast.BitExprBitOp:
			return s.dataBitOp(expr.Op, left, right)
		case *ast.BitExprArtOp:
			return s.dataArithmeticOp(expr.Op, left, right)
		default:
			return nil, errors.Wrap(ErrUnknownExpr, "unknown operation %s", logger.JSON(expr))
		}
	}
	return nil, errors.Wrap(ErrUnknownExpr, "%s", logger.JSON(expr))
}

func (s *calculator) dataBitOp(op ast.BitOperatorType, left, right data.Data) (data.Data, error) {
	r, err := func() (int, error) {
		left := left.Value()
		right := right.Value()
		switch op {
		case ast.BitOpAnd:
			return s.bitCalculator.And(left, right)
		case ast.BitOpOr:
			return s.bitCalculator.Or(left, right)
		case ast.BitOpXor:
			return s.bitCalculator.Xor(left, right)
		default:
			return 0, errors.Wrap(ErrUnknownExpr, "unknown operation %s", op)
		}
	}()
	if err != nil {
		return nil, errors.Wrap(err, "bit operation left %s right %s op %s", logger.JSON(left), logger.JSON(right), op)
	}
	return data.FromInt(r), nil
}

func (s *calculator) dataArithmeticOp(op ast.ArithmeticOperatorType, left, right data.Data) (data.Data, error) {
	r, err := func() (float64, error) {
		left := left.Value()
		right := right.Value()
		switch op {
		case ast.ArtOpAdd:
			return s.artCalculator.Add(left, right)
		case ast.ArtOpSubtract:
			return s.artCalculator.Subtract(left, right)
		case ast.ArtOpMultiply:
			return s.artCalculator.Multiply(left, right)
		case ast.ArtOpDivide:
			return s.artCalculator.Divide(left, right)
		default:
			return 0, errors.Wrap(ErrUnknownExpr, "unknown operation %s", op)
		}
	}()
	if err != nil {
		return nil, errors.Wrap(err, "arithmetic operation left %s right %s op %s", logger.JSON(left), logger.JSON(right), op)
	}
	if arithmetic.IsInt(r) {
		return data.FromInt(int(r)), nil
	}
	return data.FromFloat(r), nil
}

func (s *calculator) dataPredicate(expr ast.Predicate) (data.Data, error) {
	switch expr := expr.(type) {
	case *ast.PredicateIn:
		r, err := s.dataPredicateIn(expr)
		if err != nil {
			return nil, err
		}
		if expr.IsNot {
			return data.FromBool(!r.Bool()), nil
		}
		return r, nil
	case *ast.PredicateBetween:
		r, err := s.dataPredicateBetween(expr)
		if err != nil {
			return nil, err
		}
		if expr.IsNot {
			return data.FromBool(!r.Bool()), nil
		}
		return r, nil
	case *ast.PredicateLike:
		r, err := s.dataPredicateLike(expr)
		if err != nil {
			return nil, err
		}
		if expr.IsNot {
			return data.FromBool(!r.Bool()), nil
		}
		return r, nil
	case *ast.PredicateBitExpr:
		return s.data(expr.Expr)
	default:
		return nil, errors.Wrap(ErrUnknownExpr, "%s", logger.JSON(expr))
	}
}

func (s *calculator) dataPredicateLike(expr *ast.PredicateLike) (data.Data, error) {
	t, err := s.data(expr.Target)
	if err != nil {
		return nil, errors.Wrap(err, "target %s", logger.JSON(expr.Target))
	}
	p, err := s.data(expr.Pattern)
	if err != nil {
		return nil, errors.Wrap(err, "pattern %s", logger.JSON(expr.Target))
	}
	switch s.comparer.Like(t.Value(), p.Value()) {
	case compare.ResultMatched:
		return data.FromBool(true), nil
	case compare.ResultNotMatched:
		return data.FromBool(false), nil
	default:
		return nil, errors.Wrap(ErrUnknownExpr, "%s", logger.JSON(expr))
	}
}

func (s *calculator) dataPredicateBetween(expr *ast.PredicateBetween) (data.Data, error) {
	t, err := s.data(expr.Target)
	if err != nil {
		return nil, errors.Wrap(err, "target %s", logger.JSON(expr.Target))
	}
	l, err := s.data(expr.Left)
	if err != nil {
		return nil, errors.Wrap(err, "left %s", logger.JSON(expr.Left))
	}
	r, err := s.data(expr.Right)
	if err != nil {
		return nil, errors.Wrap(err, "right %s", logger.JSON(expr.Right))
	}
	switch s.comparer.Between(t.Value(), l.Value(), r.Value()) {
	case compare.ResultIn:
		return data.FromBool(true), nil
	case compare.ResultNotIn:
		return data.FromBool(false), nil
	default:
		return nil, errors.Wrap(ErrUnknownExpr, "%s", logger.JSON(expr))
	}
}

func (s *calculator) dataPredicateIn(expr *ast.PredicateIn) (data.Data, error) {
	t, err := s.data(expr.Target)
	if err != nil {
		return nil, errors.Wrap(err, "target %s", logger.JSON(expr.Target))
	}
	switch t.Type() {
	case data.TypeBool:
		list := make([]bool, len(expr.List.Exprs))
		for i, e := range expr.List.Exprs {
			v, err := s.data(e)
			if err != nil {
				return nil, errors.Wrap(err, "list[%d] %s", i, e)
			}
			if v.Type() != data.TypeBool {
				return nil, errors.Wrap(ErrTypeMismatch, "list[%d] %s: expected bool but got %s", i, e, v.Type())
			}
			list[i] = v.Bool()
		}
		result := s.comparer.In(t.Value(), list)
		switch result {
		case compare.ResultIn:
			return data.FromBool(true), nil
		case compare.ResultNotIn:
			return data.FromBool(false), nil
		default:
			return nil, errors.Wrap(ErrUnknownExpr, "in unknown result %s args %s %s", result, logger.JSON(t), logger.JSON(list))
		}
	case data.TypeInt:
		list := make([]int, len(expr.List.Exprs))
		for i, e := range expr.List.Exprs {
			v, err := s.data(e)
			if err != nil {
				return nil, errors.Wrap(err, "list[%d] %s", i, e)
			}
			if v.Type() != data.TypeInt {
				return nil, errors.Wrap(ErrTypeMismatch, "list[%d] %s: expected int but got %s", i, e, v.Type())
			}
			list[i] = v.Int()
		}
		result := s.comparer.In(t.Value(), list)
		switch result {
		case compare.ResultIn:
			return data.FromBool(true), nil
		case compare.ResultNotIn:
			return data.FromBool(false), nil
		default:
			return nil, errors.Wrap(ErrUnknownExpr, "in unknown result %s args %s %s", result, logger.JSON(t), logger.JSON(list))
		}
	case data.TypeFloat:
		list := make([]float64, len(expr.List.Exprs))
		for i, e := range expr.List.Exprs {
			v, err := s.data(e)
			if err != nil {
				return nil, errors.Wrap(err, "list[%d] %s", i, e)
			}
			if v.Type() != data.TypeFloat {
				return nil, errors.Wrap(ErrTypeMismatch, "list[%d] %s: expected float but got %s", i, e, v.Type())
			}
			list[i] = v.Float()
		}
		result := s.comparer.In(t.Value(), list)
		switch result {
		case compare.ResultIn:
			return data.FromBool(true), nil
		case compare.ResultNotIn:
			return data.FromBool(false), nil
		default:
			return nil, errors.Wrap(ErrUnknownExpr, "in unknown result %s args %s %s", result, logger.JSON(t), logger.JSON(list))
		}
	case data.TypeString:
		list := make([]string, len(expr.List.Exprs))
		for i, e := range expr.List.Exprs {
			v, err := s.data(e)
			if err != nil {
				return nil, errors.Wrap(err, "list[%d] %s", i, e)
			}
			if v.Type() != data.TypeString {
				return nil, errors.Wrap(ErrTypeMismatch, "list[%d] %s: expected string but got %s", i, e, v.Type())
			}
			list[i] = v.String()
		}
		result := s.comparer.In(t.Value(), list)
		switch result {
		case compare.ResultIn:
			return data.FromBool(true), nil
		case compare.ResultNotIn:
			return data.FromBool(false), nil
		default:
			return nil, errors.Wrap(ErrUnknownExpr, "in unknown result %s args %s %s", result, logger.JSON(t), logger.JSON(list))
		}
	default:
		return nil, errors.Wrap(ErrUnknownExpr, "%s", logger.JSON(expr))
	}
}

func (s *calculator) dataBoolPrimary(expr ast.BoolPrimary) (data.Data, error) {
	switch expr := expr.(type) {
	case *ast.BoolPrimaryComparison:
		left, err := s.data(expr.Left)
		if err != nil {
			return nil, errors.Wrap(err, "left arg %s", logger.JSON(expr.Left))
		}
		right, err := s.data(expr.Right)
		if err != nil {
			return nil, errors.Wrap(err, "right arg %s", logger.JSON(expr.Right))
		}
		result, err := s.compareData(expr.Op, left, right)
		if err != nil {
			return nil, errors.Wrap(err, "arg %s left %s right %s", logger.JSON(expr), logger.JSON(left), logger.JSON(right))
		}
		return result, nil
	case *ast.BoolPrimaryPredicate:
		return s.data(expr.Pred)
	default:
		return nil, errors.Wrap(ErrUnknownExpr, "%s", logger.JSON(expr))
	}
}

func (s *calculator) compareData(op ast.ComparisonType, left, right data.Data) (data.Data, error) {
	result := s.comparer.Compare(left.Value(), right.Value())
	if result == compare.ResultUndefined {
		return nil, errors.New("cannot compare left %s right %s", logger.JSON(left), logger.JSON(right))
	}
	r, err := func() (bool, error) {
		switch op {
		case ast.CmpEqual:
			return result == compare.ResultEqual, nil
		case ast.CmpNotEqual:
			return result != compare.ResultEqual, nil
		case ast.CmpGreaterThan:
			return result == compare.ResultGreaterThan, nil
		case ast.CmpGreaterEqual:
			return result != compare.ResultLessThan, nil
		case ast.CmpLessThan:
			return result == compare.ResultLessThan, nil
		case ast.CmpLessEqual:
			return result != compare.ResultGreaterThan, nil
		default:
			return false, errors.Wrap(ErrUnknownExpr, "unknown operation %s", op)
		}
	}()
	if err != nil {
		return nil, errors.Wrap(err, "compare left %s right %s op %s", logger.JSON(left), logger.JSON(right), op)
	}
	return data.FromBool(r), nil
}

func (s *calculator) dataLogicalOperator(expr ast.Expr) (data.Data, error) {
	var err error
	switch expr := expr.(type) {
	case ast.BinaryOp:
		var (
			left, right data.Data
		)
		if left, err = s.Data(expr.LeftArg()); err != nil {
			return nil, errors.Wrap(err, "left arg %s", logger.JSON(expr.LeftArg()))
		}
		if left.Type() != data.TypeBool {
			return nil, errors.Wrap(ErrTypeMismatch,
				"left data: expected bool but got %s from %s", logger.JSON(left), logger.JSON(expr.LeftArg()))
		}
		if right, err = s.Data(expr.RightArg()); err != nil {
			return nil, errors.Wrap(err, "right arg %s", logger.JSON(expr.RightArg()))
		}
		if right.Type() != data.TypeBool {
			return nil, errors.Wrap(ErrTypeMismatch,
				"right data: expected bool but got %s from %s", logger.JSON(right), logger.JSON(expr.RightArg()))
		}
		switch expr.(type) {
		case *ast.OrExpr:
			return data.FromBool(left.Bool() || right.Bool()), nil
		case *ast.AndExpr:
			return data.FromBool(left.Bool() && right.Bool()), nil
		case *ast.XorExpr:
			return data.FromBool(left.Bool() != right.Bool()), nil
		default:
			return nil, errors.Wrap(ErrUnknownExpr,
				"%s left %s right %s", logger.JSON(expr), logger.JSON(left), logger.JSON(right))
		}
	case ast.UnaryOp:
		var arg data.Data
		if arg, err = s.Data(expr.Arg()); err != nil {
			return nil, errors.Wrap(err, "arg %s", logger.JSON(expr.Arg()))
		}
		if arg.Type() != data.TypeBool {
			return nil, errors.Wrap(ErrTypeMismatch,
				"arg: expected bool but got %s from %s", logger.JSON(arg), logger.JSON(expr.Arg()))
		}
		switch expr.(type) {
		case *ast.NotExpr:
			return data.FromBool(!arg.Bool()), nil
		default:
			return nil, errors.Wrap(ErrUnknownExpr, "%s arg %s", logger.JSON(expr), logger.JSON(arg))
		}
	default:
		return nil, errors.Wrap(ErrUnknownExpr, "%s", logger.JSON(expr))
	}
}
