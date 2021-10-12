package eval

import (
	"context"
	"sort"

	"github.com/berquerant/dql/ast"
	"github.com/berquerant/dql/async"
	"github.com/berquerant/dql/calc"
	"github.com/berquerant/dql/data"
	"github.com/berquerant/dql/env"
	"github.com/berquerant/dql/errors"
)

type (
	OrderBy interface {
		Sort(ctx context.Context, table env.Map, expr ast.Expr, isDesc bool, sourceC <-chan GRow) <-chan GRow
	}

	orderBy struct {
		calcFactory func(env.Map) calc.Calculator
	}

	orderByRow struct {
		row   GRow
		value data.Data
	}
)

func NewOrderBy(calcFactory func(env.Map) calc.Calculator) OrderBy {
	return &orderBy{
		calcFactory: calcFactory,
	}
}

func (s *orderBy) Sort(
	ctx context.Context, table env.Map, expr ast.Expr, isDesc bool, sourceC <-chan GRow,
) <-chan GRow {
	resultC := make(chan GRow, resultCBufferSize)
	go func() {
		defer close(resultC)
		rows := []*orderByRow{}
		for r := range sourceC {
			if async.IsDone(ctx) {
				resultC <- NewErrGRow(errors.Wrap(ctx.Err(), "order by"))
				return
			}
			v, err := s.evalRow(table, expr, r)
			if err != nil {
				resultC <- NewErrGRow(errors.Wrap(err, "order by"))
				return
			}
			rows = append(rows, v)
		}
		if len(rows) == 0 {
			return
		}
		sf, err := s.getSortFunc(rows, isDesc)
		if err != nil {
			resultC <- NewErrGRow(err)
			return
		}
		sort.SliceStable(rows, sf)
		for _, r := range rows {
			resultC <- r.row
		}
	}()
	return resultC
}

func (s *orderBy) getSortFunc(rows []*orderByRow, isDesc bool) (func(int, int) bool, error) {
	f, err := s.sortFunc(rows)
	if err != nil {
		return nil, err
	}
	if isDesc {
		return func(i, j int) bool { return f(j, i) }, nil
	}
	return f, nil
}

func (*orderBy) sortFunc(rows []*orderByRow) (func(int, int) bool, error) {
	switch rows[0].value.Type() {
	case data.TypeInt:
		return func(i, j int) bool { return rows[i].value.Int() < rows[j].value.Int() }, nil
	case data.TypeString:
		return func(i, j int) bool { return rows[i].value.String() < rows[j].value.String() }, nil
	case data.TypeFloat:
		return func(i, j int) bool { return rows[i].value.Float() < rows[j].value.Float() }, nil
	case data.TypeBool:
		return func(i, j int) bool { return !rows[i].value.Bool() && rows[j].value.Bool() }, nil
	default:
		return nil, errors.Wrap(ErrUnknownDataType, "order by")
	}
}

func (s *orderBy) evalRow(table env.Map, expr ast.Expr, row GRow) (*orderByRow, error) {
	t, err := AppendGRowToEnv(table, row)
	if err != nil {
		return nil, err
	}
	v, err := s.calcFactory(t).Data(expr)
	if err != nil {
		return nil, err
	}
	return &orderByRow{
		row:   row,
		value: v,
	}, nil
}
