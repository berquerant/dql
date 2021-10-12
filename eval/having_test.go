package eval_test

import (
	"context"
	"testing"

	"github.com/berquerant/dql/calc"
	"github.com/berquerant/dql/data"
	"github.com/berquerant/dql/env"
	"github.com/berquerant/dql/errors"
	"github.com/berquerant/dql/eval"
	"github.com/stretchr/testify/assert"
)

func TestHaving(t *testing.T) {
	var (
		factory = func(c calc.Calculator) func(env.Map) calc.Calculator {
			return func(_ env.Map) calc.Calculator {
				return c
			}
		}
		infoRawGRow     = eval.NewRawGRow(eval.NewRow(&mockInfo{name: "mock"}))
		infoGroupedGRow = eval.NewGroupedGRow(eval.NewGroupedRow("mock", data.FromString("mmm"), []eval.Row{
			eval.NewRow(&mockInfo{name: "mock"}),
		}))
		errMockRow = errors.New("mock row")
		errGRow    = eval.NewErrGRow(errMockRow)
		yield      = func(row eval.GRow) <-chan eval.GRow {
			c := make(chan eval.GRow, 1)
			c <- row
			close(c)
			return c
		}
	)

	t.Run("accepted", func(t *testing.T) {
		got := resultToGRows(eval.NewHaving(factory(&mockCalculator{
			value: data.FromBool(true),
		})).Filter(context.TODO(), env.New(), nil, yield(infoGroupedGRow)))
		assert.Equal(t, 1, len(got))
		assert.Equal(t, "mock", got[0].Grouped().Key())
		assert.Equal(t, "mmm", got[0].Grouped().Value().String())
	})

	t.Run("denied", func(t *testing.T) {
		got := resultToGRows(eval.NewHaving(factory(&mockCalculator{
			value: data.FromBool(false),
		})).Filter(context.TODO(), env.New(), nil, yield(infoGroupedGRow)))
		assert.Equal(t, 0, len(got))
	})

	t.Run("not bool", func(t *testing.T) {
		got := resultToGRows(eval.NewHaving(factory(&mockCalculator{
			value: data.FromString("true"),
		})).Filter(context.TODO(), env.New(), nil, yield(infoGroupedGRow)))
		assert.Equal(t, 1, len(got))
		assert.ErrorIs(t, got[0].Err(), eval.ErrNotBoolExpr)
	})

	t.Run("calc error", func(t *testing.T) {
		got := resultToGRows(eval.NewHaving(factory(&mockCalculator{
			err: errMockRow,
		})).Filter(context.TODO(), env.New(), nil, yield(infoGroupedGRow)))
		assert.Equal(t, 1, len(got))
		assert.ErrorIs(t, got[0].Err(), errMockRow)
	})

	t.Run("not grouped row", func(t *testing.T) {
		got := resultToGRows(eval.NewHaving(factory(&mockCalculator{
			value: data.FromBool(true),
		})).Filter(context.TODO(), nil, nil, yield(infoRawGRow)))
		assert.Equal(t, 1, len(got))
		assert.ErrorIs(t, got[0].Err(), eval.ErrInvalidHaving)
	})

	t.Run("err row", func(t *testing.T) {
		got := resultToGRows(eval.NewHaving(factory(&mockCalculator{
			value: data.FromBool(true),
		})).Filter(context.TODO(), nil, nil, yield(errGRow)))
		assert.Equal(t, 1, len(got))
		assert.ErrorIs(t, got[0].Err(), errMockRow)
	})

	t.Run("canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.TODO())
		cancel()
		got := resultToGRows(eval.NewHaving(factory(&mockCalculator{
			value: data.FromBool(true),
		})).Filter(ctx, nil, nil, yield(infoGroupedGRow)))
		assert.Equal(t, 1, len(got))
		assert.ErrorIs(t, got[0].Err(), context.Canceled)
	})
}
