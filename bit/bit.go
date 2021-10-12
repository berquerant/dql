package bit

import (
	"strconv"

	"github.com/berquerant/dql/errors"
)

// Calculator provides bit operations.
type Calculator interface {
	Not(arg interface{}) (int, error)
	And(left, right interface{}) (int, error)
	Or(left, right interface{}) (int, error)
	Xor(left, right interface{}) (int, error)
}

// New returns a new Calculator.
func New() Calculator { return &calculator{} }

type calculator struct{}

var (
	ErrInvalidArgument = errors.New("invalid argument")
)

func (*calculator) fetch(v interface{}) (int, error) {
	switch v := v.(type) {
	case int:
		return v, nil
	case string:
		return FromBinaryString(v)
	default:
		return 0, ErrInvalidArgument
	}
}

func (s *calculator) fetchDouble(left, right interface{}) (int, int, error) {
	l, err := s.fetch(left)
	if err != nil {
		return 0, 0, errors.Wrap(err, "left %v", left)
	}
	r, err := s.fetch(right)
	if err != nil {
		return 0, 0, errors.Wrap(err, "right %v", right)
	}
	return l, r, nil
}

func (s *calculator) Not(arg interface{}) (int, error) {
	v, err := s.fetch(arg)
	if err != nil {
		return 0, errors.Wrap(err, "bit not")
	}
	return ^v, nil
}

func (s *calculator) And(left, right interface{}) (int, error) {
	l, r, err := s.fetchDouble(left, right)
	if err != nil {
		return 0, errors.Wrap(err, "bit and")
	}
	return l & r, nil
}

func (s *calculator) Or(left, right interface{}) (int, error) {
	l, r, err := s.fetchDouble(left, right)
	if err != nil {
		return 0, errors.Wrap(err, "bit or")
	}
	return l | r, nil
}

func (s *calculator) Xor(left, right interface{}) (int, error) {
	l, r, err := s.fetchDouble(left, right)
	if err != nil {
		return 0, errors.Wrap(err, "bit xor")
	}
	return l ^ r, nil
}

// FromBinaryString parses a string as a binary digits.
func FromBinaryString(v string) (int, error) {
	x, err := strconv.ParseInt(v, 2, 32)
	if err != nil {
		return 0, errors.Wrap(err, "from binary string")
	}
	return int(x), nil
}

// ToBinaryString converts an integer into a binary digits.
func ToBinaryString(v int) string { return strconv.FormatInt(int64(v), 2) }
