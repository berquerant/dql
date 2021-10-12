package eval_test

import (
	"context"
	"sort"
	"testing"

	"github.com/berquerant/dql/ast"
	"github.com/berquerant/dql/data"
	"github.com/berquerant/dql/env"
	"github.com/berquerant/dql/errors"
	"github.com/berquerant/dql/eval"
	"github.com/berquerant/dql/logger"
	"github.com/stretchr/testify/assert"
)

func resultToGRows(resultC <-chan eval.GRow) []eval.GRow {
	r := []eval.GRow{}
	for v := range resultC {
		r = append(r, v)
	}
	return r
}

func TestGroupBy(t *testing.T) {
	var (
		errMockRow = errors.New("error mock row")
		errRow     = eval.NewErrRow(errMockRow)
		tmpRow     = eval.NewRow(&mockInfo{name: "tmp"})
		yield      = func(rows ...eval.Row) <-chan eval.Row {
			c := make(chan eval.Row, len(rows))
			for _, row := range rows {
				c <- row
			}
			close(c)
			return c
		}
	)

	t.Run("group", func(t *testing.T) {
		t.Run("grouped", func(t *testing.T) {
			infos := []eval.Info{
				&mockInfo{name: "a", size: 1},
				&mockInfo{name: "b", size: 1},
				&mockInfo{name: "c", size: 2},
			}
			newRows := func(v ...eval.Info) []eval.Row {
				r := make([]eval.Row, len(v))
				for i, x := range v {
					r[i] = eval.NewRow(x)
				}
				return r
			}

			for _, tc := range []*struct {
				title string
				input []eval.Row
				want  map[interface{}][]eval.Row
			}{
				{
					title: "a row",
					input: newRows(infos[0]),
					want: map[interface{}][]eval.Row{
						1: newRows(infos[0]),
					},
				},
				{
					title: "2 rows",
					input: newRows(infos[:2]...),
					want: map[interface{}][]eval.Row{
						1: newRows(infos[:2]...),
					},
				},
				{
					title: "2 keys",
					input: newRows(infos...),
					want: map[interface{}][]eval.Row{
						1: newRows(infos[:2]...),
						2: newRows(infos[2]),
					},
				},
			} {
				tc := tc
				t.Run(tc.title, func(t *testing.T) {
					gotRows := resultToGRows(eval.NewGroupBy("size").Group(context.TODO(), env.New(), yield(tc.input...)))
					got := map[interface{}][]eval.Row{}
					for _, row := range gotRows {
						if row.Type() != eval.GroupedRowType {
							t.Fatal("got not grouped row")
						}
						g := row.Grouped()
						assert.Equal(t, "size", g.Key())
						got[g.Value().Value()] = g.Rows()
					}
					assert.Equal(t, len(tc.want), len(got))
					for k, w := range tc.want {
						g, ok := got[k]
						if !ok {
							t.Fatalf("%s not found", logger.JSON(k))
						}
						assert.Equal(t, len(w), len(g))
						sort.SliceStable(w, func(i, j int) bool { return w[i].Info().Name() < w[j].Info().Name() })
						sort.SliceStable(g, func(i, j int) bool { return g[i].Info().Name() < g[j].Info().Name() })
						for i, ww := range w {
							gg := g[i]
							assert.Equal(t, ww.Info().Name(), gg.Info().Name())
						}
					}
				})
			}
		})

		t.Run("err row", func(t *testing.T) {
			got := resultToGRows(eval.NewGroupBy("name").Group(context.TODO(), env.New(), yield(errRow)))
			assert.Equal(t, 1, len(got))
			assert.ErrorIs(t, got[0].Err(), errMockRow)
		})

		t.Run("canceled", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.TODO())
			cancel()
			got := resultToGRows(eval.NewGroupBy("name").Group(ctx, env.New(), yield(tmpRow)))
			assert.Equal(t, 1, len(got))
			assert.ErrorIs(t, got[0].Err(), context.Canceled)
		})

		t.Run("invalid ref", func(t *testing.T) {
			t.Run("invalid ident", func(t *testing.T) {
				got := resultToGRows(eval.NewGroupBy("key").Group(context.TODO(), env.New(), yield(tmpRow)))
				assert.Equal(t, 1, len(got))
				assert.ErrorIs(t, got[0].Err(), eval.ErrInvalidIdent)
			})

			for _, tc := range []*struct {
				title   string
				envData env.Data
				err     error
			}{
				{
					title:   "data",
					envData: env.FromData(data.FromBool(true)),
					err:     eval.ErrInvalidIdentRef,
				},
				{
					title: "not ident",
					envData: env.FromExpr(&ast.IntLit{
						Value: 1,
					}),
					err: eval.ErrInvalidIdentRef,
				},
			} {
				tc := tc
				t.Run(tc.title, func(t *testing.T) {
					envMap := env.New()
					envMap.Set("key", tc.envData)
					got := resultToGRows(eval.NewGroupBy("key").Group(context.TODO(), envMap, yield(tmpRow)))
					assert.Equal(t, 1, len(got))
					assert.ErrorIs(t, got[0].Err(), tc.err)
				})
			}
		})
	})

	t.Run("noop", func(t *testing.T) {
		var (
			info = &mockInfo{
				name:    "mockName",
				size:    10,
				mode:    "mockMode",
				modTime: 100,
				isDir:   true,
			}
			infoRow = eval.NewRow(info)
		)
		t.Run("raw", func(t *testing.T) {
			got := resultToGRows(eval.NewGroupBy("").Group(context.TODO(), nil, yield(infoRow)))
			assert.Equal(t, 1, len(got))
			assert.Equal(t, eval.RawRowType, got[0].Type())
			assert.Equal(t, info.name, got[0].Raw().Info().Name())
		})

		t.Run("err row", func(t *testing.T) {
			got := resultToGRows(eval.NewGroupBy("").Group(context.TODO(), nil, yield(errRow)))
			assert.Equal(t, 1, len(got))
			assert.ErrorIs(t, got[0].Err(), errMockRow)
		})

		t.Run("canceled", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.TODO())
			cancel()
			got := resultToGRows(eval.NewGroupBy("").Group(ctx, nil, yield(infoRow)))
			assert.Equal(t, 1, len(got))
			assert.ErrorIs(t, got[0].Err(), context.Canceled)
		})
	})
}
