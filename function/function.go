package function

import (
	"github.com/berquerant/dql/data"
	"github.com/berquerant/dql/errors"
)

type (
	Caller interface {
		Call(name string, args ...data.Data) (data.Data, error)
		Func(name string) (Function, bool)
	}

	Function interface {
		Name() string
		Call(args ...data.Data) (data.Data, error)
	}
)

func NewCaller(functions ...Function) Caller {
	set := map[string]Function{}
	for _, f := range functions {
		set[f.Name()] = f
	}
	return &caller{
		functions: set,
	}
}

type caller struct {
	functions map[string]Function
}

var (
	ErrNotFound        = errors.New("not found")
	ErrInvalidArgument = errors.New("invalid argument")
)

func (s *caller) Call(name string, args ...data.Data) (data.Data, error) {
	f, found := s.functions[name]
	if !found {
		return nil, errors.Wrap(ErrNotFound, "function %s", name)
	}
	r, err := f.Call(args...)
	if err != nil {
		return nil, errors.Wrap(err, "from function %s", name)
	}
	return r, nil
}

func (s *caller) Func(name string) (Function, bool) {
	f, exist := s.functions[name]
	return f, exist
}

func NewCallerWithNames(builder FactoryBuilder, functionNames ...string) Caller {
	functions := []Function{}
	for _, name := range functionNames {
		if f, ok := builder.Factory(name); ok {
			functions = append(functions, f())
		}
	}
	return NewCaller(functions...)
}
