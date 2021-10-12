package eval

import (
	"context"

	"github.com/berquerant/dql/ast"
	"github.com/berquerant/dql/async"
	"github.com/berquerant/dql/calc"
	"github.com/berquerant/dql/data"
	"github.com/berquerant/dql/env"
	"github.com/berquerant/dql/errors"
	"github.com/berquerant/dql/function"
)

type (
	Select interface {
		Select(ctx context.Context, table env.Map, sourceC <-chan GRow) <-chan SRow
	}

	selectImpl struct {
		exprs         []ast.Expr
		calcFactory   func(env.Map) calc.Calculator
		isAggregation bool
	}
)

func NewSelect(calcFactory func(env.Map) calc.Calculator, exprs []ast.Expr) Select {
	x := &selectImpl{
		exprs:       exprs,
		calcFactory: calcFactory,
	}
	x.isAggregation = x.containsAggregation()
	return x
}

func (s *selectImpl) containsAggregation() bool {
	for _, expr := range s.exprs {
		if DetectAggregation(expr) {
			return true
		}
	}
	return false
}

func (s *selectImpl) Select(ctx context.Context, table env.Map, sourceC <-chan GRow) <-chan SRow {
	resultC := make(chan SRow, resultCBufferSize)
	go func() {
		defer close(resultC)

		var (
			isRawAggregation bool
			rowCount         int
			rawRows          = []Row{}
		)
		for r := range sourceC {
			rowCount++
			if async.IsDone(ctx) {
				resultC <- NewErrSRow(errors.Wrap(ctx.Err(), "select"))
				return
			}
			if s.isAggregation && r.Type() == RawRowType {
				isRawAggregation = true
				rawRows = append(rawRows, r.Raw())
				continue
			}

			values := make([]data.Data, len(s.exprs))
			for i, expr := range s.exprs {
				v, err := s.evalRow(table, expr, r)
				if err != nil {
					resultC <- NewErrSRow(errors.Wrap(err, "select"))
					return
				}
				values[i] = v
			}
			resultC <- NewSRow(values)
		}

		if !isRawAggregation {
			return
		}
		// aggregation below
		if len(rawRows) != rowCount {
			resultC <- NewErrSRow(errors.Wrap(ErrInvalidSelectSource,
				"select got %d rows but %d raw on aggregation", rowCount, len(rawRows)))
			return
		}
		values := make([]data.Data, len(s.exprs))
		for i, expr := range s.exprs {
			v, err := s.evalAggregation(table, expr, rawRows)
			if err != nil {
				resultC <- NewErrSRow(errors.Wrap(err, "select"))
				return
			}
			values[i] = v
		}
		resultC <- NewSRow(values)
	}()
	return resultC
}

func (s *selectImpl) evalAggregation(table env.Map, expr ast.Expr, rows []Row) (data.Data, error) {
	return s.calcFactory(AppendRowsToEnv(table, rows)).Data(expr)
}

func (s *selectImpl) evalRow(table env.Map, expr ast.Expr, row GRow) (data.Data, error) {
	t, err := AppendGRowToEnv(table, row)
	if err != nil {
		return nil, err
	}
	return s.calcFactory(t).Data(expr)
}

// DetectAggregation returns true if expr contains aggregation function call.
func DetectAggregation(expr ast.Expr) bool {
	aggregations := function.AggregationFunctionNames()
	aggSet := make(map[string]bool, len(aggregations))
	for _, x := range aggregations {
		aggSet[x] = true
	}

	var (
		isDetected bool
		detector   = func(x ast.Expr) bool {
			if f, ok := x.(*ast.FunctionCall); ok {
				if aggSet[f.FunctionName.Value] {
					isDetected = true
					return false
				}
			}
			return true
		}
	)
	expr.Accept(ast.NewBaseVisitor(detector))
	return isDetected
}
