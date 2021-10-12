package function

import (
	"context"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/berquerant/dql/arithmetic"
	"github.com/berquerant/dql/bit"
	"github.com/berquerant/dql/cast"
	"github.com/berquerant/dql/chrono"
	"github.com/berquerant/dql/data"
	"github.com/berquerant/dql/errors"
	"github.com/berquerant/dql/logger"
	"github.com/berquerant/gogrep"
)

func NormalFunctionNames() []string {
	return []string{
		"now",
		"cast",
		"int2bin",
		"bin2int",
		"ext",
		"dir",
		"base",
		"floor",
		"ceil",
		"pow",
		"grep",
		"len",
		"depth",
	}
}

// NewGrep returns a new grep function.
// It greps args[1] by args[0] and returns a count of selected lines.
func NewGrep(grepper gogrep.Grepper) Function {
	return &grep{
		grepper: grepper,
	}
}

type grep struct {
	grepper gogrep.Grepper
}

func (*grep) Name() string { return "grep" }
func (s *grep) Call(args ...data.Data) (data.Data, error) {
	if len(args) != 2 {
		return nil, errors.Wrap(ErrInvalidArgument, "arg len want 2 but got %d", len(args))
	}
	var (
		pattern  = args[0]
		filename = args[1]
	)
	if !(pattern.Type() == data.TypeString && filename.Type() == data.TypeString) {
		return nil, errors.Wrap(ErrInvalidArgument,
			"arg[0] want string and arg[1] want string but got %s %s", pattern.Type(), filename.Type())
	}
	f, err := os.Open(filename.String())
	if err != nil {
		return nil, errors.Wrap(err, "cannot grep")
	}
	defer f.Close()
	resultC, err := s.grepper.Grep(context.Background(), pattern.String(), f)
	if err != nil {
		return nil, errors.Wrap(err, "failed to grep")
	}
	var c int
	for r := range resultC {
		if err := r.Err(); err != nil {
			return nil, errors.Wrap(err, "failure on grep")
		}
		c++
	}
	return data.FromInt(c), nil
}

// NewDepth returns a new depth function.
// It returns the depth of the path.
func NewDepth() Function { return &depth{} }

type depth struct{}

func (*depth) Name() string { return "depth" }
func (*depth) Call(args ...data.Data) (data.Data, error) {
	if len(args) != 1 {
		return nil, errors.Wrap(ErrInvalidArgument, "arg len want 1 but got %d", len(args))
	}
	arg := args[0]
	if arg.Type() != data.TypeString {
		return nil, errors.Wrap(ErrInvalidArgument, "arg want string but got %s", arg.Type())
	}
	return data.FromInt(strings.Count(arg.String(), "/")), nil
}

// NewNow returns a new now function.
// It returns the current time as a timestamp.
func NewNow() Function { return &now{} }

type now struct{}

func (*now) Name() string { return "now" }
func (*now) Call(args ...data.Data) (data.Data, error) {
	if len(args) > 0 {
		return nil, errors.Wrap(ErrInvalidArgument, "arg len want 0 but got %d", len(args))
	}
	return data.FromInt(chrono.Now().Unix()), nil
}

// NewCast returns a new cast function.
// Converts args[0] as args[1].
func NewCast(caster cast.Caster) Function {
	return &casting{
		caster: caster,
	}
}

type casting struct {
	caster cast.Caster
}

func (*casting) Name() string { return "cast" }
func (s *casting) Call(args ...data.Data) (data.Data, error) {
	if len(args) != 2 {
		return nil, errors.Wrap(ErrInvalidArgument, "arg len want 2 but got %d", len(args))
	}
	var (
		arg = args[0]
		to  = args[1]
	)
	if to.Type() != data.TypeString {
		return nil, errors.Wrap(ErrInvalidArgument, "cast type want string but got %s", logger.JSON(to))
	}
	r, err := s.caster.Cast(arg, s.toType(to.String()))
	if err != nil {
		return nil, errors.Wrap(err, "cannot cast")
	}
	return r, nil
}

func (*casting) toType(t string) cast.Type {
	switch strings.ToLower(t) {
	case "int":
		return cast.TypeInt
	case "float":
		return cast.TypeFloat
	case "string":
		return cast.TypeString
	case "bool":
		return cast.TypeBool
	case "timestamp":
		return cast.TypeTimestamp
	case "time":
		return cast.TypeTime
	case "duration":
		return cast.TypeDuration
	default:
		return cast.TypeUndefined
	}
}

// NewInt2Bin returns a new int2bin function.
// It parses the integer as binary digits.
func NewInt2Bin() Function { return &int2bin{} }

type int2bin struct{}

func (*int2bin) Name() string { return "int2bin" }
func (s *int2bin) Call(args ...data.Data) (data.Data, error) {
	if len(args) != 1 {
		return nil, errors.Wrap(ErrInvalidArgument, "arg len want 1 but got %d", len(args))
	}
	arg := args[0]
	if arg.Type() != data.TypeInt {
		return nil, errors.Wrap(ErrInvalidArgument, "arg type want int but got %s", arg.Type())
	}
	return data.FromString(bit.ToBinaryString(arg.Int())), nil
}

// NewBin2Int returns a new bin2int function.
// It parses the string as binary digits.
func NewBin2Int() Function { return &bin2int{} }

type bin2int struct{}

func (*bin2int) Name() string { return "bin2int" }
func (*bin2int) Call(args ...data.Data) (data.Data, error) {
	if len(args) != 1 {
		return nil, errors.Wrap(ErrInvalidArgument, "arg len want 1 but got %d", len(args))
	}
	arg := args[0]
	if arg.Type() != data.TypeString {
		return nil, errors.Wrap(ErrInvalidArgument, "arg type want string but got %s", arg.Type())
	}
	r, err := bit.FromBinaryString(arg.String())
	if err != nil {
		return nil, errors.Wrap(ErrInvalidArgument, "arg %s", arg.String())
	}
	return data.FromInt(r), nil
}

// NewExt returns a new ext function.
// It returns the extension of given file path.
func NewExt() Function { return &ext{} }

type ext struct{}

func (*ext) Name() string { return "ext" }
func (*ext) Call(args ...data.Data) (data.Data, error) {
	if len(args) != 1 {
		return nil, errors.Wrap(ErrInvalidArgument, "arg len want 1 but got %d", len(args))
	}
	arg := args[0]
	if arg.Type() != data.TypeString {
		return nil, errors.Wrap(ErrInvalidArgument, "arg type want string but got %s", arg.Type())
	}
	return data.FromString(filepath.Base(arg.String())), nil
}

// NewBase returns a new dir function.
// It returns the directory of given file path.
func NewDir() Function { return &dir{} }

type dir struct{}

func (*dir) Name() string { return "dir" }
func (*dir) Call(args ...data.Data) (data.Data, error) {
	if len(args) != 1 {
		return nil, errors.Wrap(ErrInvalidArgument, "arg len want 1 but got %d", len(args))
	}
	arg := args[0]
	if arg.Type() != data.TypeString {
		return nil, errors.Wrap(ErrInvalidArgument, "arg type want string but got %s", arg.Type())
	}
	return data.FromString(filepath.Dir(arg.String())), nil
}

// NewBase returns a new base function.
// It returns the last element of given file path.
func NewBase() Function { return &base{} }

type base struct{}

func (*base) Name() string { return "base" }
func (*base) Call(args ...data.Data) (data.Data, error) {
	if len(args) != 1 {
		return nil, errors.Wrap(ErrInvalidArgument, "arg len want 1 but got %d", len(args))
	}
	arg := args[0]
	if arg.Type() != data.TypeString {
		return nil, errors.Wrap(ErrInvalidArgument, "arg type want string but got %s", arg.Type())
	}
	return data.FromString(filepath.Base(arg.String())), nil
}

// NewLen returns a new len function.
// It returns the length of the string.
func NewLen() Function { return &length{} }

type length struct{}

func (*length) Name() string { return "len" }
func (*length) Call(args ...data.Data) (data.Data, error) {
	if len(args) != 1 {
		return nil, errors.Wrap(ErrInvalidArgument, "arg len want 1 but got %d", len(args))
	}
	arg := args[0]
	if arg.Type() != data.TypeString {
		return nil, errors.Wrap(ErrInvalidArgument, "arg type want string but got %s", arg.Type())
	}
	return data.FromInt(len(arg.String())), nil
}

// NewFloor returns a new floor function.
func NewFloor() Function { return &floor{} }

type floor struct{}

func (*floor) Name() string { return "floor" }
func (*floor) Call(args ...data.Data) (data.Data, error) {
	if len(args) != 1 {
		return nil, errors.Wrap(ErrInvalidArgument, "arg len want 1 but got %d", len(args))
	}
	arg := args[0]
	switch arg.Type() {
	case data.TypeInt:
		return arg, nil
	case data.TypeFloat:
		return data.FromInt(int(math.Floor(arg.Float()))), nil
	default:
		return nil, errors.Wrap(ErrInvalidArgument, "arg type want int or float but got %s", arg.Type())
	}
}

// NewCeil returns a new ceil function.
func NewCeil() Function { return &ceil{} }

type ceil struct{}

func (*ceil) Name() string { return "ceil" }
func (*ceil) Call(args ...data.Data) (data.Data, error) {
	if len(args) != 1 {
		return nil, errors.Wrap(ErrInvalidArgument, "arg len want 1 but got %d", len(args))
	}
	arg := args[0]
	switch arg.Type() {
	case data.TypeInt:
		return arg, nil
	case data.TypeFloat:
		return data.FromInt(int(math.Ceil(arg.Float()))), nil
	default:
		return nil, errors.Wrap(ErrInvalidArgument, "arg type want int or float but got %s", arg.Type())
	}
}

// NewPow returns new pow function.
func NewPow(calculator arithmetic.Calculator) Function {
	return &pow{
		calculator: calculator,
	}
}

type pow struct {
	calculator arithmetic.Calculator
}

func (*pow) Name() string { return "pow" }
func (s *pow) Call(args ...data.Data) (data.Data, error) {
	if len(args) != 2 {
		return nil, errors.Wrap(ErrInvalidArgument, "arg len want 2 but got %d", len(args))
	}
	r, err := s.calculator.Pow(args[0].Value(), args[1].Value())
	if err != nil {
		return nil, errors.Wrap(err, "left %v right %v", args[0].Value(), args[1].Value())
	}
	if arithmetic.IsInt(r) {
		return data.FromInt(int(r)), nil
	}
	return data.FromFloat(r), nil
}
