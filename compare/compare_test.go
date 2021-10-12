package compare_test

import (
	"testing"

	"github.com/berquerant/dql/compare"
	"github.com/stretchr/testify/assert"
)

func TestComparer(t *testing.T) {
	t.Run("Like", func(t *testing.T) {
		for _, tc := range []*struct {
			title           string
			target, pattern interface{}
			want            compare.Result
		}{
			{
				title: "all nil",
				want:  compare.ResultUndefined,
			},
			{
				title:  "pattern nil",
				target: "interface",
				want:   compare.ResultUndefined,
			},
			{
				title:   "target nil",
				pattern: "interface",
				want:    compare.ResultUndefined,
			},
			{
				title:   "pattern is not string",
				target:  "interface",
				pattern: true,
				want:    compare.ResultUndefined,
			},
			{
				title:   "target is not string",
				target:  10,
				pattern: "interface",
				want:    compare.ResultUndefined,
			},
			{
				title:   "not matched",
				target:  "for whom the bell tolls",
				pattern: "sea",
				want:    compare.ResultNotMatched,
			},
			{
				title:   "matched",
				target:  "for whom the bell tolls",
				pattern: "^for",
				want:    compare.ResultMatched,
			},
		} {
			tc := tc
			t.Run(tc.title, func(t *testing.T) {
				assert.Equal(t, tc.want, compare.New().Like(tc.target, tc.pattern))
			})
		}
	})

	t.Run("Between", func(t *testing.T) {
		for _, tc := range []*struct {
			title                string
			target, lower, upper interface{}
			want                 compare.Result
		}{
			{
				title: "all nil",
				want:  compare.ResultUndefined,
			},
			{
				title: "target is nil",
				lower: 1,
				upper: 2,
				want:  compare.ResultUndefined,
			},
			{
				title:  "lower is nil",
				target: 1,
				upper:  2,
				want:   compare.ResultUndefined,
			},
			{
				title:  "upper is nil",
				lower:  1,
				target: 2,
				want:   compare.ResultUndefined,
			},
			{
				title:  "type mismatch target",
				target: "1",
				lower:  1,
				upper:  2,
				want:   compare.ResultUndefined,
			},
			{
				title:  "type mismatch lower",
				target: 1,
				lower:  "1",
				upper:  2,
				want:   compare.ResultUndefined,
			},
			{
				title:  "type mismatch upper",
				target: 1,
				lower:  1,
				upper:  "2",
				want:   compare.ResultUndefined,
			},
			{
				title:  "invalid type",
				target: true,
				lower:  false,
				upper:  true,
				want:   compare.ResultUndefined,
			},
			{
				title:  "contains int upper edge",
				target: 5,
				lower:  1,
				upper:  5,
				want:   compare.ResultIn,
			},
			{
				title:  "contains int lower edge",
				target: 5,
				lower:  5,
				upper:  10,
				want:   compare.ResultIn,
			},
			{
				title:  "contains int equal",
				target: 5,
				lower:  5,
				upper:  5,
				want:   compare.ResultIn,
			},
			{
				title:  "not contain int",
				target: 6,
				lower:  1,
				upper:  5,
				want:   compare.ResultNotIn,
			},
			{
				title:  "not contain int invalid range",
				target: 6,
				lower:  10,
				upper:  5,
				want:   compare.ResultNotIn,
			},
			{
				title:  "contains float upper edge",
				target: float64(5),
				lower:  float64(1),
				upper:  float64(5),
				want:   compare.ResultIn,
			},
			{
				title:  "contains float lower edge",
				target: float64(5),
				lower:  float64(5),
				upper:  float64(10),
				want:   compare.ResultIn,
			},
			{
				title:  "contains float equal",
				target: float64(5),
				lower:  float64(5),
				upper:  float64(5),
				want:   compare.ResultIn,
			},
			{
				title:  "not contain float",
				target: float64(6),
				lower:  float64(1),
				upper:  float64(5),
				want:   compare.ResultNotIn,
			},
			{
				title:  "not contain float invalid range",
				target: float64(6),
				lower:  float64(10),
				upper:  float64(5),
				want:   compare.ResultNotIn,
			},
			{
				title:  "contains string upper edge",
				target: "b",
				lower:  "a",
				upper:  "b",
				want:   compare.ResultIn,
			},
			{
				title:  "contains string lower edge",
				target: "a",
				lower:  "a",
				upper:  "b",
				want:   compare.ResultIn,
			},
			{
				title:  "contains string equal",
				target: "a",
				lower:  "a",
				upper:  "a",
				want:   compare.ResultIn,
			},
			{
				title:  "not contain string",
				target: "c",
				lower:  "a",
				upper:  "b",
				want:   compare.ResultNotIn,
			},
			{
				title:  "not contain string invalid range",
				target: "b",
				lower:  "c",
				upper:  "a",
				want:   compare.ResultNotIn,
			},
		} {
			tc := tc
			t.Run(tc.title, func(t *testing.T) {
				assert.Equal(t, tc.want, compare.New().Between(tc.target, tc.lower, tc.upper))
			})
		}
	})

	t.Run("In", func(t *testing.T) {
		for _, tc := range []*struct {
			title        string
			target, list interface{}
			want         compare.Result
		}{
			{
				title:  "invalid list",
				target: 1,
				list:   1,
				want:   compare.ResultUndefined,
			},
			{
				title:  "type mismatched list",
				target: 1,
				list:   []string{"1"},
				want:   compare.ResultUndefined,
			},
			{
				title:  "empty list",
				target: 1,
				list:   []int{},
				want:   compare.ResultNotIn,
			},
			{
				title:  "contains int",
				target: 1,
				list:   []int{1},
				want:   compare.ResultIn,
			},
			{
				title:  "not contain int",
				target: 1,
				list:   []int{2},
				want:   compare.ResultNotIn,
			},
			{
				title:  "contains float",
				target: float64(1),
				list:   []float64{1.0, 2.1},
				want:   compare.ResultIn,
			},
			{
				title:  "not contain float",
				target: float64(1),
				list:   []float64{2.1},
				want:   compare.ResultNotIn,
			},
			{
				title:  "contains bool",
				target: true,
				list:   []bool{true},
				want:   compare.ResultIn,
			},
			{
				title:  "not contain bool",
				target: true,
				list:   []bool{false},
				want:   compare.ResultNotIn,
			},
			{
				title:  "contains string",
				target: "s",
				list:   []string{"s", "t"},
				want:   compare.ResultIn,
			},
			{
				title:  "not contain string",
				target: "s",
				list:   []string{"t"},
				want:   compare.ResultNotIn,
			},
		} {
			tc := tc
			t.Run(tc.title, func(t *testing.T) {
				assert.Equal(t, tc.want, compare.New().In(tc.target, tc.list))
			})
		}
	})

	t.Run("Compare", func(t *testing.T) {
		for _, tc := range []*struct {
			title       string
			left, right interface{}
			want        compare.Result
		}{
			{
				title: "nil",
				want:  compare.ResultUndefined,
			},
			{
				title: "type mismatch",
				left:  1,
				right: "1",
				want:  compare.ResultUndefined,
			},
			{
				title: "bool less",
				left:  false,
				right: true,
				want:  compare.ResultLessThan,
			},
			{
				title: "bool greater",
				left:  true,
				right: false,
				want:  compare.ResultGreaterThan,
			},
			{
				title: "bool equal true",
				left:  true,
				right: true,
				want:  compare.ResultEqual,
			},
			{
				title: "bool equal false",
				left:  false,
				right: false,
				want:  compare.ResultEqual,
			},
			{
				title: "int less",
				left:  1,
				right: 10,
				want:  compare.ResultLessThan,
			},
			{
				title: "int greater",
				left:  10,
				right: 1,
				want:  compare.ResultGreaterThan,
			},
			{
				title: "int equal",
				left:  10,
				right: 10,
				want:  compare.ResultEqual,
			},
			{
				title: "float less",
				left:  float64(1),
				right: float64(10),
				want:  compare.ResultLessThan,
			},
			{
				title: "float greater",
				left:  float64(10),
				right: float64(1),
				want:  compare.ResultGreaterThan,
			},
			{
				title: "string less",
				left:  "a",
				right: "b",
				want:  compare.ResultLessThan,
			},
			{
				title: "string greater",
				left:  "b",
				right: "a",
				want:  compare.ResultGreaterThan,
			},
			{
				title: "string equal",
				left:  "a",
				right: "a",
				want:  compare.ResultEqual,
			},
		} {
			tc := tc
			t.Run(tc.title, func(t *testing.T) {
				assert.Equal(t, tc.want, compare.New().Compare(tc.left, tc.right))
			})
		}
	})
}
