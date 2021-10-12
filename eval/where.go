package eval

import (
	"context"

	"github.com/berquerant/dql/ast"
	"github.com/berquerant/dql/async"
	"github.com/berquerant/dql/calc"
	"github.com/berquerant/dql/data"
	"github.com/berquerant/dql/env"
	"github.com/berquerant/dql/errors"
)

type (
	Where interface {
		Filter(ctx context.Context, table env.Map, expr ast.Expr, sourceC <-chan Row) <-chan Row
	}

	where struct {
		calcFactory func(env.Map) calc.Calculator
	}
)

func NewWhere(calcFactory func(env.Map) calc.Calculator) Where {
	return &where{
		calcFactory: calcFactory,
	}
}

func (s *where) Filter(ctx context.Context, table env.Map, expr ast.Expr, sourceC <-chan Row) <-chan Row {
	resultC := make(chan Row, resultCBufferSize)
	go func() {
		defer close(resultC)
		for row := range sourceC {
			if async.IsDone(ctx) {
				resultC <- NewErrRow(errors.Wrap(ctx.Err(), "where"))
				return
			}
			if err := row.Err(); err != nil {
				resultC <- NewErrRow(errors.Wrap(err, "where"))
				return
			}
			t := AppendRowToEnv(table, row)
			err := s.filter(t, expr)
			switch {
			case err == nil:
				// accepted
				resultC <- row
			case errors.Is(err, errFiltered):
				// denied
				continue
			default:
				resultC <- NewErrRow(errors.Wrap(err, "where"))
				return
			}
		}
	}()
	return resultC
}

func (s *where) filter(table env.Map, expr ast.Expr) error {
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
