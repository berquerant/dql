package cast

import (
	"fmt"
	"math"
	"strconv"

	"github.com/berquerant/dql/chrono"
	"github.com/berquerant/dql/data"
	"github.com/berquerant/dql/errors"
	"github.com/berquerant/dql/logger"
)

type Type int

//go:generate stringer -type Type -output cast_stringer_generated.go

const (
	TypeUndefined Type = iota
	// TypeInt is an integer.
	TypeInt
	// TypeFloat is a floating point.
	TypeFloat
	// TypeString is a string.
	TypeString
	// TypeBool is a bool.
	TypeBool
	// TypeTimestamp is UNIX time.
	TypeTimestamp
	// TypeTime is a time string.
	TypeTime
	// TypeDuration is a duration string.
	TypeDuration
)

// Caster provides cast operations.
type Caster interface {
	// Cast converts values to specified type.
	Cast(v data.Data, t Type) (data.Data, error)
}

var (
	ErrCannotCast = errors.New("cannot cast")
)

// New returns a new Caster.
func New() Caster { return &caster{} }

type caster struct{}

func (s *caster) Cast(v data.Data, t Type) (data.Data, error) {
	r, err := s.cast(v, t)
	if err != nil {
		return nil, errors.Wrap(err, "cast %s to %s", logger.JSON(v), t)
	}
	return r, nil
}

func (s *caster) cast(v data.Data, t Type) (data.Data, error) {
	switch t {
	case TypeInt:
		return s.toInt(v)
	case TypeFloat:
		return s.toFloat(v)
	case TypeString:
		return s.toString(v)
	case TypeBool:
		return s.toBool(v)
	case TypeTimestamp:
		return s.toTimestamp(v)
	case TypeTime:
		return s.toTime(v)
	case TypeDuration:
		return s.toDuration(v)
	default:
		return nil, ErrCannotCast
	}
}

func (*caster) toInt(v data.Data) (data.Data, error) {
	switch v.Type() {
	case data.TypeInt:
		return v, nil
	case data.TypeFloat:
		return data.FromInt(int(math.Floor(v.Float()))), nil
	case data.TypeString:
		r, err := strconv.Atoi(v.String())
		if err != nil {
			return nil, errors.Wrap(err, "%s to int", v.String())
		}
		return data.FromInt(r), nil
	case data.TypeBool:
		if v.Bool() {
			return data.FromInt(1), nil
		}
		return data.FromInt(0), nil
	default:
		return nil, errors.Wrap(ErrCannotCast, "unknown data")
	}
}

func (*caster) toFloat(v data.Data) (data.Data, error) {
	switch v.Type() {
	case data.TypeInt:
		return data.FromFloat(float64(v.Int())), nil
	case data.TypeFloat:
		return v, nil
	case data.TypeString:
		r, err := strconv.ParseFloat(v.String(), 64)
		if err != nil {
			return nil, errors.Wrap(err, "%s to float", v.String())
		}
		return data.FromFloat(r), nil
	case data.TypeBool:
		if v.Bool() {
			return data.FromFloat(1), nil
		}
		return data.FromFloat(0), nil
	default:
		return nil, errors.Wrap(ErrCannotCast, "unknown data")
	}
}

func (*caster) toString(v data.Data) (data.Data, error) {
	return data.FromString(fmt.Sprint(v.Value())), nil
}

func (*caster) toBool(v data.Data) (data.Data, error) {
	switch v.Type() {
	case data.TypeInt:
		return data.FromBool(v.Int() != 0), nil
	case data.TypeFloat:
		return data.FromBool(v.Float() != 0), nil
	case data.TypeString:
		return data.FromBool(v.String() != ""), nil
	case data.TypeBool:
		return v, nil
	default:
		return nil, errors.Wrap(ErrCannotCast, "unknown data")
	}
}

func (*caster) toTimestamp(v data.Data) (data.Data, error) {
	switch v.Type() {
	case data.TypeString:
		t, err := chrono.TimeFromString(v.String())
		if err != nil {
			return nil, errors.Wrap(err, "%s to timestamp", v.String())
		}
		return data.FromInt(t.Unix()), nil
	default:
		return nil, errors.Wrap(ErrCannotCast, "unknown data")
	}
}

func (*caster) toTime(v data.Data) (data.Data, error) {
	switch v.Type() {
	case data.TypeInt:
		return data.FromString(chrono.TimeFromTimestamp(v.Int()).String()), nil
	case data.TypeFloat:
		return data.FromString(chrono.TimeFromTimestamp(int(math.Floor(v.Float()))).String()), nil
	default:
		return nil, errors.Wrap(ErrCannotCast, "unknown data")
	}
}

func (*caster) toDuration(v data.Data) (data.Data, error) {
	switch v.Type() {
	case data.TypeInt:
		return data.FromString(chrono.DurationFromInt(v.Int()).String()), nil
	case data.TypeFloat:
		return data.FromString(chrono.DurationFromInt(int(math.Floor(v.Float()))).String()), nil
	default:
		return nil, errors.Wrap(ErrCannotCast, "unknown data")
	}
}
