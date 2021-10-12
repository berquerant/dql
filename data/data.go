package data

import (
	"encoding/json"
	"fmt"

	"github.com/berquerant/dql/arithmetic"
)

type Type int

//go:generate stringer -type Type -output data_type_stringer_generated.go

const (
	TypeInt Type = iota
	TypeFloat
	TypeString
	TypeBool
)

func (s Type) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, s.String())), nil
}

// Data is a data container of this repository.
type Data interface {
	// Type returns the type of the content.
	Type() Type
	// Value returns the raw value of the content.
	Value() interface{}
	// Int returns an int content.
	// Returns 0 if the content is not an int.
	Int() int
	// Float returns a float content.
	// Returns 0 if the content is not a float.
	Float() float64
	// String returns a string content.
	// Returns an empty string if the content is not a string.
	String() string
	// Bool returns a bool content.
	// Rturns false if the content is not a bool.
	Bool() bool
	Clone() Data
}

// FromInt returns a new Data with int.
func FromInt(v int) Data {
	return &data{
		typ: TypeInt,
		val: v,
	}
}

// FromFloat returns a new Data with float.
func FromFloat(v float64) Data {
	return &data{
		typ: TypeFloat,
		val: v,
	}
}

// FromString returns a new Data with string.
func FromString(v string) Data {
	return &data{
		typ: TypeString,
		val: v,
	}
}

// FromBool returns a new Data with bool.
func FromBool(v bool) Data {
	return &data{
		typ: TypeBool,
		val: v,
	}
}

func FromInterface(v interface{}) (Data, bool) {
	switch v := v.(type) {
	case int:
		return FromInt(v), true
	case float64:
		if arithmetic.IsInt(v) {
			return FromInt(int(v)), true
		}
		return FromFloat(v), true
	case float32:
		return FromInterface(float64(v))
	case string:
		return FromString(v), true
	case bool:
		return FromBool(v), true
	default:
		return nil, false
	}
}

type data struct {
	typ Type
	val interface{}
}

func (s *data) Type() Type         { return s.typ }
func (s *data) Value() interface{} { return s.val }

func (s *data) Clone() Data {
	switch s.typ {
	case TypeInt:
		return FromInt(s.Int())
	case TypeFloat:
		return FromFloat(s.Float())
	case TypeString:
		return FromString(s.String())
	case TypeBool:
		return FromBool(s.Bool())
	default:
		panic("unreachable: unknown data type")
	}
}

func (s *data) Int() int {
	if x, ok := s.val.(int); ok {
		return x
	}
	return 0
}

func (s *data) Float() float64 {
	if x, ok := s.val.(float64); ok {
		return x
	}
	return 0
}

func (s *data) String() string {
	if x, ok := s.val.(string); ok {
		return x
	}
	return ""
}

func (s *data) Bool() bool {
	if x, ok := s.val.(bool); ok {
		return x
	}
	return false
}

func (s *data) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  s.typ,
		"value": s.val,
	})
}
