package eval_test

import (
	"context"
	"testing"

	"github.com/berquerant/dql/ast"
	"github.com/berquerant/dql/calc"
	"github.com/berquerant/dql/data"
	"github.com/berquerant/dql/env"
	"github.com/berquerant/dql/errors"
	"github.com/berquerant/dql/eval"
	"github.com/stretchr/testify/assert"
)

type mockCalculator struct {
	value data.Data
	err   error
}

func (s *mockCalculator) Data(_ ast.Expr) (data.Data, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.value, nil
}

type mockRow struct {
	err  error
	info eval.Info
}

func (s *mockRow) Err() error      { return s.err }
func (s *mockRow) Info() eval.Info { return s.info }

type mockInfo struct {
	name    string
	size    int
	mode    string
	modTime int
	isDir   bool
}

func (s *mockInfo) Name() string { return s.name }
func (s *mockInfo) Size() int    { return s.size }
func (s *mockInfo) Mode() string { return s.mode }
func (s *mockInfo) ModTime() int { return s.modTime }
func (s *mockInfo) IsDir() bool  { return s.isDir }
func (s *mockInfo) ToMap() map[string]data.Data {
	return map[string]data.Data{
		"name":     data.FromString(s.name),
		"size":     data.FromInt(s.size),
		"mode":     data.FromString(s.mode),
		"mod_time": data.FromInt(s.modTime),
		"is_dir":   data.FromBool(s.isDir),
	}
}

func TestWhere(t *testing.T) {
	var (
		factory = func(c calc.Calculator) func(env.Map) calc.Calculator {
			return func(_ env.Map) calc.Calculator {
				return c
			}
		}
		infoRow = &mockRow{
			info: &mockInfo{
				name: "mock",
			},
		}
		errMockRow = errors.New("mock row")
		errRow     = &mockRow{
			err: errMockRow,
		}
		yield = func(row eval.Row) <-chan eval.Row {
			c := make(chan eval.Row, 1)
			c <- row
			close(c)
			return c
		}
	)

	t.Run("accept", func(t *testing.T) {
		resultC := eval.NewWhere(factory(&mockCalculator{
			value: data.FromBool(true),
		})).Filter(context.TODO(), env.New(), nil, yield(infoRow))
		got := resultToRows(resultC)
		assert.Equal(t, 1, len(got))
		assert.Equal(t, "mock", got[0].Info().Name())
	})

	t.Run("deny", func(t *testing.T) {
		resultC := eval.NewWhere(factory(&mockCalculator{
			value: data.FromBool(false),
		})).Filter(context.TODO(), env.New(), nil, yield(infoRow))
		got := resultToRows(resultC)
		assert.Equal(t, 0, len(got))
	})

	t.Run("invalid expr type", func(t *testing.T) {
		resultC := eval.NewWhere(factory(&mockCalculator{
			value: data.FromString("true"),
		})).Filter(context.TODO(), env.New(), nil, yield(infoRow))
		got := resultToRows(resultC)
		assert.Equal(t, 1, len(got))
		assert.ErrorIs(t, got[0].Err(), eval.ErrNotBoolExpr)
	})

	t.Run("calc error", func(t *testing.T) {
		resultC := eval.NewWhere(factory(&mockCalculator{
			err: errMockRow,
		})).Filter(context.TODO(), env.New(), nil, yield(infoRow))
		got := resultToRows(resultC)
		assert.Equal(t, 1, len(got))
		assert.ErrorIs(t, got[0].Err(), errMockRow)
	})

	t.Run("error row", func(t *testing.T) {
		resultC := eval.NewWhere(factory(&mockCalculator{
			value: data.FromBool(true),
		})).Filter(context.TODO(), nil, nil, yield(errRow))
		got := resultToRows(resultC)
		assert.Equal(t, 1, len(got))
		assert.ErrorIs(t, got[0].Err(), errMockRow)
	})

	t.Run("canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.TODO())
		cancel()
		resultC := eval.NewWhere(factory(&mockCalculator{
			value: data.FromBool(true),
		})).Filter(ctx, nil, nil, yield(infoRow))
		got := resultToRows(resultC)
		assert.Equal(t, 1, len(got))
		assert.ErrorIs(t, got[0].Err(), context.Canceled)
	})
}
