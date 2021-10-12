package eval

import (
	"context"

	"github.com/berquerant/dql/ast"
	"github.com/berquerant/dql/buf"
	"github.com/berquerant/dql/calc"
	"github.com/berquerant/dql/dig"
	"github.com/berquerant/dql/env"
	"github.com/berquerant/dql/errors"
	"github.com/berquerant/dql/logger"
	"github.com/berquerant/dql/preprocessor"
)

const (
	AllSelectSymbol = "all"
)

type (
	Runner interface {
		Run(ctx context.Context, names ...string) <-chan SRow
		Headers() []string
	}

	runner struct {
		stmt *ast.Statement
	}
)

func NewRunner(stmt *ast.Statement) Runner {
	r := &runner{
		stmt: stmt,
	}
	r.init()
	return r
}

func (s *runner) init() {
	if err := s.preprocess(); err != nil {
		logger.Error(err.Error())
	}
}

func (s *runner) Run(ctx context.Context, names ...string) <-chan SRow {
	var (
		table   = s.prepareEnv()
		where   = func(sourceC <-chan Row) <-chan Row { return s.where(ctx, table, sourceC) }
		groupBy = func(sourceC <-chan Row) <-chan GRow { return s.groupBy(ctx, table, sourceC) }
		having  = func(sourceC <-chan GRow) <-chan GRow { return s.having(ctx, table, sourceC) }
		orderBy = func(sourceC <-chan GRow) <-chan GRow { return s.orderBy(ctx, table, sourceC) }
		limit   = func(sourceC <-chan GRow) <-chan GRow { return s.limit(ctx, sourceC) }
		selekt  = func(sourceC <-chan GRow) <-chan SRow { return s.selekt(ctx, table, sourceC) }
	)
	return selekt(limit(orderBy(having(groupBy(where(NewSource(dig.New()).Yield(ctx, names...)))))))
}

func (s *runner) Headers() []string {
	b := buf.NewStrings()
	for _, t := range s.stmt.SelectSection.Terms.Terms {
		if t.As != nil {
			b.Add(t.As.Value)
			continue
		}
		b.Add(t.Target.Expr.String())
	}
	return b.Get()
}

func (s *runner) preprocess() error {
	ps := []preprocessor.PreProcessor{
		preprocessor.NewSelectAll(AllSelectSymbol),
	}
	for _, p := range ps {
		if err := p.PreProcess(s.stmt); err != nil {
			return errors.Wrap(err, "runner")
		}
	}
	return nil
}

func (s *runner) selekt(ctx context.Context, table env.Map, sourceC <-chan GRow) <-chan SRow {
	exprs := make([]ast.Expr, len(s.stmt.SelectSection.Terms.Terms))
	for i, t := range s.stmt.SelectSection.Terms.Terms {
		exprs[i] = t.Target.Expr
	}
	// TODO: implement distinct
	return NewSelect(calc.NewAggregation, exprs).Select(ctx, table, sourceC)
}

func (s *runner) limit(ctx context.Context, sourceC <-chan GRow) <-chan GRow {
	if s.stmt.LimitSection == nil {
		return sourceC
	}
	var (
		lim    = s.stmt.LimitSection.Limit.Value
		offset = 0
	)
	if s.stmt.LimitSection.Offset != nil {
		offset = s.stmt.LimitSection.Offset.Value
	}
	return NewLimit().Limit(ctx, lim, offset, sourceC)
}

func (s *runner) orderBy(ctx context.Context, table env.Map, sourceC <-chan GRow) <-chan GRow {
	if s.stmt.OrderBySection == nil {
		return sourceC
	}
	// TODO: multiple orderby
	t := s.stmt.OrderBySection.Terms.Terms[0]
	return NewOrderBy(calc.NewAggregation).Sort(ctx, table, t.Expr, t.Option.IsDesc, sourceC)
}

func (s *runner) having(ctx context.Context, table env.Map, sourceC <-chan GRow) <-chan GRow {
	if s.stmt.HavingSection == nil {
		return sourceC
	}
	return NewHaving(calc.NewAggregation).Filter(ctx, table, s.stmt.HavingSection.Condition.Expr, sourceC)
}

func (s *runner) groupBy(ctx context.Context, table env.Map, sourceC <-chan Row) <-chan GRow {
	if s.stmt.GroupBySection == nil {
		return NewGroupBy("").Group(ctx, table, sourceC)
	}
	// TODO: multiple groupby
	ident := s.stmt.GroupBySection.Terms.Terms[0].Expr.(*ast.BoolPrimaryPredicate).Pred.(*ast.PredicateBitExpr).Expr.(*ast.BitExprSimpleExpr).Expr.(*ast.Ident)
	return NewGroupBy(ident.Value).Group(ctx, table, sourceC)
}

func (s *runner) where(ctx context.Context, table env.Map, sourceC <-chan Row) <-chan Row {
	if s.stmt.WhereSection == nil {
		return sourceC
	}
	return NewWhere(calc.NewNormal).Filter(ctx, table, s.stmt.WhereSection.Condition.Expr, sourceC)
}

func (s *runner) prepareEnv() env.Map {
	x := env.New()
	// aliases, e.g. select size as x
	for _, t := range s.stmt.SelectSection.Terms.Terms {
		if t.As != nil {
			x.Set(t.As.Value, env.FromExpr(t.Target.Expr))
		}
	}
	return x
}
