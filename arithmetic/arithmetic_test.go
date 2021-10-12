package arithmetic_test

import (
	"testing"

	"github.com/berquerant/dql/arithmetic"
	"github.com/stretchr/testify/assert"
)

func TestCalculator(t *testing.T) {
	const delta = 0.000001

	t.Run("Pow", func(t *testing.T) {
		for _, tc := range []*struct {
			title       string
			left, right interface{}
			want        float64
		}{
			{
				title: "int ^ int",
				left:  2,
				right: 3,
				want:  8,
			},
			{
				title: "int ^ float",
				left:  2,
				right: 1.2,
				want:  2.2973967,
			},
			{
				title: "float ^ int",
				left:  1.1,
				right: 2,
				want:  1.21,
			},
			{
				title: "float * float",
				left:  1.1,
				right: 1.1,
				want:  1.1105342,
			},
		} {
			tc := tc
			t.Run(tc.title, func(t *testing.T) {
				got, err := arithmetic.New().Pow(tc.left, tc.right)
				assert.Nil(t, err)
				assert.InDelta(t, tc.want, got, delta)
			})
		}
	})

	t.Run("Divide", func(t *testing.T) {
		t.Run("zero division", func(t *testing.T) {
			_, err := arithmetic.New().Divide(1, 0)
			assert.NotNil(t, err)
		})

		for _, tc := range []*struct {
			title       string
			left, right interface{}
			want        float64
		}{
			{
				title: "int / int",
				left:  1,
				right: 2,
				want:  0.5,
			},
			{
				title: "int / float",
				left:  3,
				right: 1.5,
				want:  2,
			},
			{
				title: "float / int",
				left:  1.1,
				right: 2,
				want:  0.55,
			},
			{
				title: "float + float",
				left:  3.6,
				right: 1.2,
				want:  3,
			},
		} {
			tc := tc
			t.Run(tc.title, func(t *testing.T) {
				got, err := arithmetic.New().Divide(tc.left, tc.right)
				assert.Nil(t, err)
				assert.InDelta(t, tc.want, got, delta)
			})
		}
	})

	t.Run("Multiply", func(t *testing.T) {
		for _, tc := range []*struct {
			title       string
			left, right interface{}
			want        float64
		}{
			{
				title: "int * int",
				left:  1,
				right: 2,
				want:  2,
			},
			{
				title: "int * float",
				left:  1,
				right: 2.1,
				want:  2.1,
			},
			{
				title: "float * int",
				left:  1.1,
				right: 2,
				want:  2.2,
			},
			{
				title: "float * float",
				left:  1.1,
				right: 2.1,
				want:  2.31,
			},
		} {
			tc := tc
			t.Run(tc.title, func(t *testing.T) {
				got, err := arithmetic.New().Multiply(tc.left, tc.right)
				assert.Nil(t, err)
				assert.InDelta(t, tc.want, got, delta)
			})
		}
	})

	t.Run("Subtract", func(t *testing.T) {
		for _, tc := range []*struct {
			title       string
			left, right interface{}
			want        float64
		}{
			{
				title: "int - int",
				left:  1,
				right: 2,
				want:  -1,
			},
			{
				title: "int - float",
				left:  1,
				right: 2.1,
				want:  -1.1,
			},
			{
				title: "float - int",
				left:  1.1,
				right: 2,
				want:  -0.9,
			},
			{
				title: "float - float",
				left:  1.1,
				right: 2.1,
				want:  -1,
			},
		} {
			tc := tc
			t.Run(tc.title, func(t *testing.T) {
				got, err := arithmetic.New().Subtract(tc.left, tc.right)
				assert.Nil(t, err)
				assert.InDelta(t, tc.want, got, delta)
			})
		}
	})

	t.Run("Add", func(t *testing.T) {
		for _, tc := range []*struct {
			title       string
			left, right interface{}
			want        float64
		}{
			{
				title: "int + int",
				left:  1,
				right: 2,
				want:  3,
			},
			{
				title: "int + float",
				left:  1,
				right: 2.1,
				want:  3.1,
			},
			{
				title: "float + int",
				left:  1.1,
				right: 2,
				want:  3.1,
			},
			{
				title: "float + float",
				left:  1.1,
				right: 2.1,
				want:  3.2,
			},
		} {
			tc := tc
			t.Run(tc.title, func(t *testing.T) {
				got, err := arithmetic.New().Add(tc.left, tc.right)
				assert.Nil(t, err)
				assert.InDelta(t, tc.want, got, delta)
			})
		}
	})
}

func TestIsInt(t *testing.T) {
	t.Run("not int", func(t *testing.T) {
		assert.False(t, arithmetic.IsInt(1.1))
	})
	t.Run("is int", func(t *testing.T) {
		assert.True(t, arithmetic.IsInt(float64(1)))
	})
}
