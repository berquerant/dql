package bit_test

import (
	"testing"

	"github.com/berquerant/dql/bit"
	"github.com/stretchr/testify/assert"
)

func TestCalculator(t *testing.T) {
	t.Run("Xor", func(t *testing.T) {
		for _, tc := range []*struct {
			title       string
			left, right interface{}
			want        int
		}{
			{
				title: "int and int",
				left:  11,
				right: 6,
				want:  13,
			},
			{
				title: "int and string",
				left:  11,
				right: "0110",
				want:  13,
			},
			{
				title: "string and int",
				left:  "1011",
				right: 6,
				want:  13,
			},
			{
				title: "string and string",
				left:  "1011",
				right: "0110",
				want:  13,
			},
		} {
			tc := tc
			t.Run(tc.title, func(t *testing.T) {
				got, err := bit.New().Xor(tc.left, tc.right)
				assert.Nil(t, err)
				assert.Equal(t, tc.want, got)
			})
		}
	})

	t.Run("Or", func(t *testing.T) {
		for _, tc := range []*struct {
			title       string
			left, right interface{}
			want        int
		}{
			{
				title: "int and int",
				left:  11,
				right: 6,
				want:  15,
			},
			{
				title: "int and string",
				left:  11,
				right: "0110",
				want:  15,
			},
			{
				title: "string and int",
				left:  "1011",
				right: 6,
				want:  15,
			},
			{
				title: "string and string",
				left:  "1011",
				right: "0110",
				want:  15,
			},
		} {
			tc := tc
			t.Run(tc.title, func(t *testing.T) {
				got, err := bit.New().Or(tc.left, tc.right)
				assert.Nil(t, err)
				assert.Equal(t, tc.want, got)
			})
		}
	})

	t.Run("And", func(t *testing.T) {
		for _, tc := range []*struct {
			title       string
			left, right interface{}
			want        int
		}{
			{
				title: "int and int",
				left:  11,
				right: 6,
				want:  2,
			},
			{
				title: "int and string",
				left:  11,
				right: "0110",
				want:  2,
			},
			{
				title: "string and int",
				left:  "1011",
				right: 6,
				want:  2,
			},
			{
				title: "string and string",
				left:  "1011",
				right: "0110",
				want:  2,
			},
		} {
			tc := tc
			t.Run(tc.title, func(t *testing.T) {
				got, err := bit.New().And(tc.left, tc.right)
				assert.Nil(t, err)
				assert.Equal(t, tc.want, got)
			})
		}
	})

	t.Run("Not", func(t *testing.T) {
		for _, tc := range []*struct {
			title string
			arg   interface{}
			want  int
		}{
			{
				title: "int",
				arg:   10,
				want:  -11,
			},
			{
				title: "string",
				arg:   "1010",
				want:  -11,
			},
		} {
			tc := tc
			t.Run(tc.title, func(t *testing.T) {
				got, err := bit.New().Not(tc.arg)
				assert.Nil(t, err)
				assert.Equal(t, tc.want, got)
			})
		}
	})
}
