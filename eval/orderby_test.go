package eval_test

import (
	"context"
	"testing"

	"github.com/berquerant/dql/ast"
	"github.com/berquerant/dql/calc"
	"github.com/berquerant/dql/data"
	"github.com/berquerant/dql/env"
	"github.com/berquerant/dql/eval"
	"github.com/stretchr/testify/assert"
)

type mockMultipleCalculator struct {
	values []data.Data
	i      int
}

func (s *mockMultipleCalculator) Data(_ ast.Expr) (data.Data, error) {
	defer func() {
		s.i++
	}()
	return s.values[s.i], nil
}

func TestOrderBy(t *testing.T) {
	var (
		factory = func(c calc.Calculator) func(env.Map) calc.Calculator {
			return func(_ env.Map) calc.Calculator {
				return c
			}
		}
		makeRows = func(names ...string) []eval.GRow {
			rows := make([]eval.GRow, len(names))
			for i, x := range names {
				rows[i] = eval.NewRawGRow(eval.NewRow(&mockInfo{
					name: x,
				}))
			}
			return rows
		}
		yield = func(rows []eval.GRow) <-chan eval.GRow {
			c := make(chan eval.GRow, len(rows))
			for _, r := range rows {
				c <- r
			}
			close(c)
			return c
		}
	)

	for _, tc := range []*struct {
		title  string
		names  []string
		isDesc bool
		want   []string
	}{
		{
			title: "zero",
		},
		{
			title: "an element",
			names: []string{"a"},
			want:  []string{"a"},
		},
		{
			title: "two elements",
			names: []string{"b", "a"},
			want:  []string{"a", "b"},
		},
		{
			title:  "three elements desc",
			names:  []string{"b", "a", "c"},
			want:   []string{"c", "b", "a"},
			isDesc: true,
		},
		{
			title: "three elements no changes",
			names: []string{"a", "b", "c"},
			want:  []string{"a", "b", "c"},
		},
		{
			title:  "three elements desc no changes",
			names:  []string{"c", "b", "a"},
			want:   []string{"c", "b", "a"},
			isDesc: true,
		},
	} {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			ret := make([]data.Data, len(tc.names))
			for i, n := range tc.names {
				ret[i] = data.FromString(n)
			}
			got := resultToGRows(eval.NewOrderBy(factory(&mockMultipleCalculator{
				values: ret,
			})).Sort(context.TODO(), env.New(), nil, tc.isDesc, yield(makeRows(tc.names...))))
			assert.Equal(t, len(tc.want), len(got))
			for i, g := range got {
				n := g.Raw().Info().Name()
				assert.Equal(t, tc.want[i], n)
			}
		})
	}
}
