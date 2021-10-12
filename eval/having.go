package eval

import (
	"context"

	"github.com/berquerant/dql/ast"
	"github.com/berquerant/dql/async"
	"github.com/berquerant/dql/calc"
	"github.com/berquerant/dql/data"
	"github.com/berquerant/dql/env"
	"github.com/berquerant/dql/errors"
	"github.com/berquerant/dql/logger"
)

type (
	Having interface {
		Filter(ctx context.Context, table env.Map, expr ast.Expr, sourceC <-chan GRow) <-chan GRow
	}

	having struct {
		calcFactory func(env.Map) calc.Calculator
	}
)

func NewHaving(calcFactory func(env.Map) calc.Calculator) Having {
	return &having{
		calcFactory: calcFactory,
	}
}

func (s *having) Filter(ctx context.Context, table env.Map, expr ast.Expr, sourceC <-chan GRow) <-chan GRow {
	resultC := make(chan GRow, resultCBufferSize)
	go func() {
		defer close(resultC)
		for row := range sourceC {
			if async.IsDone(ctx) {
				resultC <- NewErrGRow(errors.Wrap(ctx.Err(), "having"))
				return
			}
			if err := row.Err(); err != nil {
				resultC <- NewErrGRow(errors.Wrap(err, "having"))
				return
			}
			if row.Type() != GroupedRowType {
				// having requires group by
				resultC <- NewErrGRow(errors.Wrap(ErrInvalidHaving, "row type %s %s", row.Type(), logger.JSON(row)))
				return
			}
			t := AppendGroupedRowToEnv(table, row.Grouped())
			err := s.filter(t, expr)
			switch {
			case err == nil:
				// accepted
				resultC <- row
			case errors.Is(err, errFiltered):
				// denied
				continue
			default:
				resultC <- NewErrGRow(errors.Wrap(err, "having"))
				return
			}
		}
	}()
	return resultC
}

func (s *having) filter(table env.Map, expr ast.Expr) error {
	r, err := s.calcFactory(table).Data(expr)
	if err != nil {
		return err
	}
	if r.Type() != data.TypeBool {
		return ErrNotBoolExpr
	}
	if r.Bool() {
		return nil
	}
	return errFiltered
}
