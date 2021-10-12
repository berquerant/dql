package function

import (
	"github.com/berquerant/dql/arithmetic"
	"github.com/berquerant/dql/compare"
	"github.com/berquerant/dql/data"
	"github.com/berquerant/dql/errors"
	"github.com/berquerant/dql/logger"
)

// Aggregation expects to aggregate data from rows.
type Aggregation interface {
	Function
	IsAggregation()
}

func AggregationFunctionNames() []string {
	return []string{
		"count",
		"min",
		"max",
		"product",
		"sum",
		"avg",
	}
}

//go:generate marker -method IsAggregation -output aggregation_marker_generated.go -type count,min,max,sum,product,avg

// NewCount returns a new count function.
// It counts up the length of the arguments.
func NewCount() Aggregation { return &count{} }

type count struct{}

func (*count) Name() string { return "count" }
func (*count) Call(args ...data.Data) (data.Data, error) {
	return data.FromInt(len(args)), nil
}

// NewMin returns a new min function.
// It returns the minimum value of the arguments.
func NewMin(comparer compare.Comparer) Aggregation {
	return &min{
		comparer: comparer,
	}
}

type min struct {
	comparer compare.Comparer
}

func (*min) Name() string { return "min" }
func (s *min) Call(args ...data.Data) (data.Data, error) {
	switch len(args) {
	case 0:
		return nil, errors.Wrap(ErrInvalidArgument, "arg len want positive but got 0")
	case 1:
		return args[0], nil
	default:
		v := args[0]
		for _, x := range args[1:] {
			switch s.comparer.Compare(v.Value(), x.Value()) {
			case compare.ResultGreaterThan:
				v = x
			case compare.ResultEqual, compare.ResultLessThan:
				continue
			default:
				return nil, errors.Wrap(ErrInvalidArgument, "cannot compare %s and %s", logger.JSON(v), logger.JSON(x))
			}
		}
		return v, nil
	}
}

// NewMax returns a new max function.
// It returns the maximum value of the arguments.
func NewMax(comparer compare.Comparer) Aggregation {
	return &max{
		comparer: comparer,
	}
}

type max struct {
	comparer compare.Comparer
}

func (*max) Name() string { return "max" }
func (s *max) Call(args ...data.Data) (data.Data, error) {
	switch len(args) {
	case 0:
		return nil, errors.Wrap(ErrInvalidArgument, "arg len want positive but got 0")
	case 1:
		return args[0], nil
	default:
		v := args[0]
		for _, x := range args[1:] {
			switch s.comparer.Compare(v.Value(), x.Value()) {
			case compare.ResultLessThan:
				v = x
			case compare.ResultEqual, compare.ResultGreaterThan:
				continue
			default:
				return nil, errors.Wrap(ErrInvalidArgument, "cannot compare %s and %s", logger.JSON(v), logger.JSON(x))
			}
		}
		return v, nil
	}
}

// NewProduct returns a new product function.
// It returns the product of the arguments.
func NewProduct(calculator arithmetic.Calculator) Aggregation {
	return &product{
		calculator: calculator,
	}
}

type product struct {
	calculator arithmetic.Calculator
}

func (*product) Name() string { return "product" }
func (s *product) Call(args ...data.Data) (data.Data, error) {
	if len(args) == 0 {
		return nil, errors.Wrap(ErrInvalidArgument, "arg len want positive but got 0")
	}
	var (
		acc float64 = 1
		err error
	)
	for i, a := range args {
		if acc, err = s.calculator.Multiply(acc, a.Value()); err != nil {
			return nil, errors.Wrap(err, "args[%d] %v", i, a)
		}
	}
	if arithmetic.IsInt(acc) {
		return data.FromInt(int(acc)), nil
	}
	return data.FromFloat(acc), nil
}

// NewSum returns a new sum function.
// It returns the sum of the arguments.
func NewSum(calculator arithmetic.Calculator) Aggregation {
	return &sum{
		calculator: calculator,
	}
}

type sum struct {
	calculator arithmetic.Calculator
}

func (*sum) Name() string { return "sum" }
func (s *sum) Call(args ...data.Data) (data.Data, error) {
	if len(args) == 0 {
		return nil, errors.Wrap(ErrInvalidArgument, "arg len want positive but got 0")
	}
	var (
		acc float64
		err error
	)
	for i, a := range args {
		if acc, err = s.calculator.Add(acc, a.Value()); err != nil {
			return nil, errors.Wrap(err, "args[%d] %v", i, a)
		}
	}
	if arithmetic.IsInt(acc) {
		return data.FromInt(int(acc)), nil
	}
	return data.FromFloat(acc), nil
}

// NewAvg returns a new avg function.
// It returns the average of the arguments.
func NewAvg(calculator arithmetic.Calculator, sum Aggregation) Aggregation {
	return &avg{
		sum:        sum,
		calculator: calculator,
	}
}

type avg struct {
	sum        Aggregation
	calculator arithmetic.Calculator
}

func (*avg) Name() string { return "avg" }
func (s *avg) Call(args ...data.Data) (data.Data, error) {
	if len(args) == 0 {
		return nil, errors.Wrap(ErrInvalidArgument, "arg len want positive but got 0")
	}
	v, err := s.sum.Call(args...)
	if err != nil {
		return nil, err
	}
	r, err := s.calculator.Divide(v.Value(), len(args))
	if err != nil {
		return nil, err
	}
	if arithmetic.IsInt(r) {
		return data.FromInt(int(r)), nil
	}
	return data.FromFloat(r), nil
}
