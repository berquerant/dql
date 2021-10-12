package ast

import (
	"fmt"
	"strings"

	"github.com/berquerant/dql/buf"
)

type (
	Section interface {
		Node
		IsSection()
	}
)

//go:generate marker -method IsSection -type SelectSection,WhereSection,HavingSection,GroupBySection,OrderBySection,LimitSection -output section_marker_generated.go

//go:generate marker -method IsNode -type SelectSection,SelectTerms,SelectTerm,SelectOption,SelectTarget,WhereSection,WhereCondition,GroupBySection,GroupByTerms,GroupByTerm,HavingSection,OrderBySection,OrderByTerms,OrderByTerm,OrderByTermOption,LimitSection -output section_node_marker_generated.go

type (
	SelectSection struct {
		Terms  *SelectTerms  `json:"terms,omitempty"`
		Option *SelectOption `json:"option,omitempty"`
	}
	SelectTerms struct {
		Terms []*SelectTerm `json:"terms,omitempty"`
	}
	SelectTerm struct {
		Target *SelectTarget `json:"target,omitempty"`
		As     *Ident        `json:"as,omitempty"`
	}
	SelectOption struct {
		IsDistinct bool `json:"distinct,omitempty"`
	}
	SelectTarget struct {
		Expr Expr `json:"expr,omitempty"`
	}

	WhereSection struct {
		Condition *WhereCondition `json:"condition,omitempty"`
	}
	WhereCondition struct {
		Expr Expr `json:"expr,omitempty"`
	}

	GroupBySection struct {
		Terms *GroupByTerms `json:"terms,omitempty"`
	}
	GroupByTerms struct {
		Terms []*GroupByTerm `json:"terms,omitempty"`
	}
	GroupByTerm struct {
		Expr Expr `json:"expr,omitempty"`
	}

	HavingSection struct {
		Condition *WhereCondition `json:"condition,omitempty"`
	}

	OrderBySection struct {
		Terms *OrderByTerms `json:"terms,omitempty"`
	}
	OrderByTerms struct {
		Terms []*OrderByTerm `json:"terms,omitempty"`
	}
	OrderByTerm struct {
		Expr   Expr               `json:"expr,omitempty"`
		Option *OrderByTermOption `json:"option,omitempty"`
	}
	OrderByTermOption struct {
		IsDesc bool `json:"desc,omitempty"`
	}

	LimitSection struct {
		Limit  *IntLit `json:"limit,omitempty"`
		Offset *IntLit `json:"offset,omitempty"`
	}
)

func (s *SelectSection) String() string {
	b := buf.NewStrings()
	b.Add("select")
	if s.Option != nil {
		if s.Option.String() != "" {
			b.Add(s.Option.String())
		}
	}
	b.Add(s.Terms.String())
	return strings.Join(b.Get(), " ")
}

func (s *SelectTerms) String() string {
	b := buf.NewStrings()
	for _, t := range s.Terms {
		b.Add(t.String())
	}
	return strings.Join(b.Get(), ", ")
}

func (s *SelectTerm) String() string {
	b := buf.NewStrings()
	b.Add(s.Target.String())
	if s.As != nil {
		b.Add(s.As.String())
	}
	return strings.Join(b.Get(), " ")
}

func (s *SelectOption) String() string {
	if s.IsDistinct {
		return "distinct"
	}
	return ""
}

func (s *SelectTarget) String() string {
	return s.Expr.String()
}

func (s *WhereSection) String() string {
	return fmt.Sprintf("where %s", s.Condition)
}

func (s *WhereCondition) String() string {
	return s.Expr.String()
}

func (s *GroupBySection) String() string {
	return fmt.Sprintf("group by %s", s.Terms)
}

func (s *GroupByTerms) String() string {
	b := buf.NewStrings()
	for _, x := range s.Terms {
		b.Add(x.String())
	}
	return strings.Join(b.Get(), ", ")
}

func (s *GroupByTerm) String() string {
	return s.Expr.String()
}

func (s *HavingSection) String() string {
	return fmt.Sprintf("having %s", s.Condition)
}

func (s *OrderBySection) String() string {
	return fmt.Sprintf("order by %s", s.Terms)
}

func (s *OrderByTerms) String() string {
	b := buf.NewStrings()
	for _, x := range s.Terms {
		b.Add(x.String())
	}
	return strings.Join(b.Get(), ", ")
}

func (s *OrderByTerm) String() string {
	b := buf.NewStrings()
	b.Add(s.Expr.String())
	if s.Option != nil {
		if s.Option.String() != "" {
			b.Add(s.Option.String())
		}
	}
	return strings.Join(b.Get(), " ")
}

func (s *OrderByTermOption) String() string {
	if s.IsDesc {
		return "desc"
	}
	return ""
}

func (s *LimitSection) String() string {
	b := buf.NewStrings()
	b.Add("limit")
	b.Add(s.Limit.String())
	if s.Offset != nil {
		b.Add("offset")
		b.Add(s.Offset.String())
	}
	return strings.Join(b.Get(), " ")
}
