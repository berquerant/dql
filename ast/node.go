package ast

import (
	"strings"

	"github.com/berquerant/dql/buf"
)

type (
	Node interface {
		IsNode()
		String() string
	}
)

//go:generate marker -method IsNode -type Statement -output statement_marker_generated.go

type (
	Statement struct {
		SelectSection  *SelectSection  `json:"select,omitempty"`
		WhereSection   *WhereSection   `json:"where,omitempty"`
		HavingSection  *HavingSection  `json:"having,omitempty"`
		GroupBySection *GroupBySection `json:"group_by,omitempty"`
		OrderBySection *OrderBySection `json:"order_by,omitempty"`
		LimitSection   *LimitSection   `json:"limit,omitempty"`
	}
)

func (s *Statement) String() string {
	b := buf.NewStrings()
	b.Add(s.SelectSection.String())
	if s.WhereSection != nil {
		b.Add(s.WhereSection.String())
	}
	if s.GroupBySection != nil {
		b.Add(s.GroupBySection.String())
	}
	if s.HavingSection != nil {
		b.Add(s.HavingSection.String())
	}
	if s.OrderBySection != nil {
		b.Add(s.OrderBySection.String())
	}
	if s.LimitSection != nil {
		b.Add(s.LimitSection.String())
	}
	return strings.Join(b.Get(), " ") + ";"
}
