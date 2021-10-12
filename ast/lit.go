package ast

import (
	"fmt"
	"strconv"
)

type (
	// Lit is a literal.
	Lit interface {
		Expr
		IsLit()
	}
)

//go:generate marker -method IsNode,IsLit,IsExpr -type IntLit,FloatLit,StringLit -output lit_marker_generated.go

type (
	IntLit struct {
		Value int `json:"value"`
	}

	FloatLit struct {
		Value float64 `json:"value"`
	}

	StringLit struct {
		Value string `json:"value"`
	}
)

func (s *IntLit) String() string    { return strconv.Itoa(s.Value) }
func (s *FloatLit) String() string  { return fmt.Sprint(s.Value) }
func (s *StringLit) String() string { return fmt.Sprintf(`"%s"`, s.Value) }
