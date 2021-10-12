package function_test

import (
	"math"
	"testing"

	"github.com/berquerant/dql/compare"
	"github.com/berquerant/dql/data"
	"github.com/berquerant/dql/function"
	"github.com/stretchr/testify/assert"
)

func prepareData(n int) []data.Data {
	v := make([]data.Data, n)
	for i := 0; i < n; i++ {
		v[i] = data.FromInt(i)
	}
	return v
}

func TestCount(t *testing.T) {
	for _, tc := range []*struct {
		title string
		args  []data.Data
		want  int
	}{
		{
			title: "no args",
			want:  0,
		},
		{
			title: "1 arg",
			args:  prepareData(1),
			want:  1,
		},
		{
			title: "2 args",
			args:  prepareData(2),
			want:  2,
		},
	} {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			got, err := function.NewCount().Call(tc.args...)
			assert.Nil(t, err)
			assert.Equal(t, data.TypeInt, got.Type())
			assert.Equal(t, tc.want, got.Int())
		})
	}
}

type mockComparer struct {
	result compare.Result
}

func (s *mockComparer) Compare(_, _ interface{}) compare.Result {
	return s.result
}
func (*mockComparer) In(_, _ interface{}) compare.Result {
	return compare.ResultUndefined
}
func (*mockComparer) Between(_, _, _ interface{}) compare.Result {
	return compare.ResultUndefined
}
func (*mockComparer) Like(_, _ interface{}) compare.Result {
	return compare.ResultUndefined
}

func TestMin(t *testing.T) {
	t.Run("no args", func(t *testing.T) {
		_, err := function.NewMin(nil).Call()
		assert.ErrorIs(t, err, function.ErrInvalidArgument)
	})
	t.Run("1 arg", func(t *testing.T) {
		r, err := function.NewMin(nil).Call(data.FromInt(1))
		assert.Nil(t, err)
		assert.Equal(t, r.Int(), 1)
	})
	t.Run("2 args got error", func(t *testing.T) {
		_, err := function.NewMin(&mockComparer{
			result: compare.ResultUndefined,
		}).Call(prepareData(2)...)
		assert.ErrorIs(t, err, function.ErrInvalidArgument)
	})
	t.Run("2 args", func(t *testing.T) {
		args := prepareData(2)
		for _, tc := range []*struct {
			title  string
			result compare.Result
			want   int
		}{
			{
				title:  "first",
				result: compare.ResultLessThan,
				want:   0,
			},
			{
				title:  "second",
				result: compare.ResultGreaterThan,
				want:   1,
			},
		} {
			tc := tc
			t.Run(tc.title, func(t *testing.T) {
				r, err := function.NewMin(&mockComparer{
					result: tc.result,
				}).Call(args...)
				assert.Nil(t, err)
				assert.Equal(t, r.Int(), tc.want)
			})
		}
	})
}

func TestMax(t *testing.T) {
	t.Run("no args", func(t *testing.T) {
		_, err := function.NewMax(nil).Call()
		assert.ErrorIs(t, err, function.ErrInvalidArgument)
	})
	t.Run("1 arg", func(t *testing.T) {
		r, err := function.NewMax(nil).Call(data.FromInt(1))
		assert.Nil(t, err)
		assert.Equal(t, r.Int(), 1)
	})
	t.Run("2 args got error", func(t *testing.T) {
		_, err := function.NewMax(&mockComparer{
			result: compare.ResultUndefined,
		}).Call(prepareData(2)...)
		assert.ErrorIs(t, err, function.ErrInvalidArgument)
	})
	t.Run("2 args", func(t *testing.T) {
		args := prepareData(2)
		for _, tc := range []*struct {
			title  string
			result compare.Result
			want   int
		}{
			{
				title:  "first",
				result: compare.ResultGreaterThan,
				want:   0,
			},
			{
				title:  "second",
				result: compare.ResultLessThan,
				want:   1,
			},
		} {
			tc := tc
			t.Run(tc.title, func(t *testing.T) {
				r, err := function.NewMax(&mockComparer{
					result: tc.result,
				}).Call(args...)
				assert.Nil(t, err)
				assert.Equal(t, r.Int(), tc.want)
			})
		}
	})
}

type mockArtCalculator struct{}

func (*mockArtCalculator) Add(left, right interface{}) (float64, error) {
	return left.(float64) + right.(float64), nil
}
func (*mockArtCalculator) Subtract(left, right interface{}) (float64, error) {
	return left.(float64) - right.(float64), nil
}
func (*mockArtCalculator) Multiply(left, right interface{}) (float64, error) {
	return left.(float64) * right.(float64), nil
}
func (*mockArtCalculator) Divide(left, right interface{}) (float64, error) {
	if r, ok := right.(int); ok {
		return left.(float64) / float64(r), nil
	}
	return left.(float64) / right.(float64), nil
}
func (*mockArtCalculator) Pow(left, right interface{}) (float64, error) {
	return math.Pow(left.(float64), right.(float64)), nil
}

func TestSum(t *testing.T) {
	t.Run("no args", func(t *testing.T) {
		_, err := function.NewSum(nil).Call()
		assert.ErrorIs(t, err, function.ErrInvalidArgument)
	})
	t.Run("do", func(t *testing.T) {
		for _, tc := range []*struct {
			title string
			args  []int
			want  int
		}{
			{
				title: "1 arg",
				args:  []int{1},
				want:  1,
			},
			{
				title: "2 args",
				args:  []int{1, 2},
				want:  3,
			},
		} {
			tc := tc
			t.Run(tc.title, func(t *testing.T) {
				args := make([]data.Data, len(tc.args))
				for i, x := range tc.args {
					args[i] = data.FromFloat(float64(x))
				}
				r, err := function.NewSum(&mockArtCalculator{}).Call(args...)
				assert.Nil(t, err)
				assert.Equal(t, tc.want, r.Int())
			})
		}
	})
}

func TestProduct(t *testing.T) {
	t.Run("no args", func(t *testing.T) {
		_, err := function.NewProduct(nil).Call()
		assert.ErrorIs(t, err, function.ErrInvalidArgument)
	})
	t.Run("do", func(t *testing.T) {
		for _, tc := range []*struct {
			title string
			args  []int
			want  int
		}{
			{
				title: "1 arg",
				args:  []int{1},
				want:  1,
			},
			{
				title: "2 args",
				args:  []int{1, 2},
				want:  2,
			},
		} {
			tc := tc
			t.Run(tc.title, func(t *testing.T) {
				args := make([]data.Data, len(tc.args))
				for i, x := range tc.args {
					args[i] = data.FromFloat(float64(x))
				}
				r, err := function.NewProduct(&mockArtCalculator{}).Call(args...)
				assert.Nil(t, err)
				assert.Equal(t, tc.want, r.Int())
			})
		}
	})
}

type mockAggregation struct {
	v   data.Data
	err error
}

func (*mockAggregation) IsAggregation() {}
func (*mockAggregation) Name() string   { return "mockAggregation" }
func (s *mockAggregation) Call(_ ...data.Data) (data.Data, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.v, nil
}

func TestAvg(t *testing.T) {
	t.Run("no args", func(t *testing.T) {
		_, err := function.NewAvg(nil, nil).Call()
		assert.ErrorIs(t, err, function.ErrInvalidArgument)
	})
	t.Run("do", func(t *testing.T) {
		for _, tc := range []*struct {
			title string
			args  []data.Data
			sum   data.Data
			want  float64
		}{
			{
				title: "1 arg",
				args:  prepareData(1),
				sum:   data.FromFloat(1.4),
				want:  1.4,
			},
			{
				title: "2 args",
				args:  prepareData(2),
				sum:   data.FromFloat(2.4),
				want:  1.2,
			},
		} {
			tc := tc
			t.Run(tc.title, func(t *testing.T) {
				r, err := function.NewAvg(
					&mockArtCalculator{},
					&mockAggregation{
						v: tc.sum,
					},
				).Call(tc.args...)
				assert.Nil(t, err)
				assert.Equal(t, tc.want, r.Float())
			})
		}
	})
}
