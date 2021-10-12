package eval_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/berquerant/dql/eval"
	"github.com/stretchr/testify/assert"
)

func TestLimit(t *testing.T) {
	var (
		makeRows = func(n int) []eval.GRow {
			rows := make([]eval.GRow, n)
			for i := 0; i < n; i++ {
				rows[i] = eval.NewRawGRow(eval.NewRow(&mockInfo{
					name: fmt.Sprint(i),
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

	t.Run("limit", func(t *testing.T) {
		for _, tc := range []*struct {
			title         string
			rows          []eval.GRow
			limit, offset int
			want          []string
		}{
			{
				title: "limit 1",
				rows:  makeRows(10),
				limit: 1,
				want:  []string{"0"},
			},
			{
				title: "just limit",
				rows:  makeRows(3),
				limit: 3,
				want:  []string{"0", "1", "2"},
			},
			{
				title: "over limit",
				rows:  makeRows(3),
				limit: 10,
				want:  []string{"0", "1", "2"},
			},
			{
				title:  "offset 1",
				rows:   makeRows(3),
				limit:  10,
				offset: 1,
				want:   []string{"1", "2"},
			},
			{
				title:  "just offset",
				rows:   makeRows(3),
				limit:  10,
				offset: 3,
			},
			{
				title:  "limit and offset",
				rows:   makeRows(10),
				limit:  3,
				offset: 5,
				want:   []string{"5", "6", "7"},
			},
		} {
			tc := tc
			t.Run(tc.title, func(t *testing.T) {
				got := resultToGRows(eval.NewLimit().Limit(context.TODO(), tc.limit, tc.offset, yield(tc.rows)))
				assert.Equal(t, len(tc.want), len(got))
				for i, w := range tc.want {
					g := got[i].Raw().Info().Name()
					assert.Equal(t, w, g)
				}
			})
		}
	})

	t.Run("invalid limit", func(t *testing.T) {
		got := resultToGRows(eval.NewLimit().Limit(context.TODO(), 0, 0, yield(makeRows(1))))
		assert.Equal(t, 1, len(got))
		assert.ErrorIs(t, got[0].Err(), eval.ErrInvalidLimit)
	})

	t.Run("invalid offset", func(t *testing.T) {
		got := resultToGRows(eval.NewLimit().Limit(context.TODO(), 1, -1, yield(makeRows(1))))
		assert.Equal(t, 1, len(got))
		assert.ErrorIs(t, got[0].Err(), eval.ErrInvalidLimit)
	})
}
