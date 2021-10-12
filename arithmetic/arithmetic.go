package arithmetic

import (
	"math"

	"github.com/berquerant/dql/errors"
)

// Calculator provides basic arithmetic operations.
type Calculator interface {
	Add(left, right interface{}) (float64, error)
	Subtract(left, right interface{}) (float64, error)
	Multiply(left, right interface{}) (float64, error)
	Divide(left, right interface{}) (float64, error)
	Pow(left, right interface{}) (float64, error)
}

// New returns a new Calculator.
func New() Calculator { return &calculator{} }

type calculator struct{}

var (
	ErrInvalidArgument = errors.New("invalid argument")
)

func (*calculator) fetch(left, right interface{}) (float64, float64, error) {
	get := func(v interface{}) (float64, error) {
		switch v := v.(type) {
		case int:
			return float64(v), nil
		case float64:
			return v, nil
		default:
			return 0, ErrInvalidArgument
		}
	}
	l, err := get(left)
	if err != nil {
		return 0, 0, errors.Wrap(err, "left %v", left)
	}
	r, err := get(right)
	if err != nil {
		return 0, 0, errors.Wrap(err, "right %v", right)
	}
	return l, r, nil
}

func (s *calculator) Add(left, right interface{}) (float64, error) {
	l, r, err := s.fetch(left, right)
	if err != nil {
		return 0, errors.Wrap(err, "art add")
	}
	return l + r, nil
}

func (s *calculator) Subtract(left, right interface{}) (float64, error) {
	l, r, err := s.fetch(left, right)
	if err != nil {
		return 0, errors.Wrap(err, "art subtract")
	}
	return l - r, nil
}

func (s *calculator) Multiply(left, right interface{}) (float64, error) {
	l, r, err := s.fetch(left, right)
	if err != nil {
		return 0, errors.Wrap(err, "art multiply")
	}
	return l * r, nil
}

func (s *calculator) Divide(left, right interface{}) (float64, error) {
	l, r, err := s.fetch(left, right)
	if err != nil {
		return 0, errors.Wrap(err, "art divide")
	}
	if r == 0 {
		return 0, errors.Wrap(ErrInvalidArgument, "art divide zero division left %v right %v", l, r)
	}
	return l / r, nil
}

func (s *calculator) Pow(left, right interface{}) (float64, error) {
	l, r, err := s.fetch(left, right)
	if err != nil {
		return 0, errors.Wrap(err, "art pow")
	}
	return math.Pow(l, r), nil
}

// IsInt returns true if v is an integer.
func IsInt(v float64) bool { return v == math.Floor(v) }
