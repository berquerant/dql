package ast

import (
	"fmt"
	"strings"

	"github.com/berquerant/dql/buf"
)

type (
	Expr interface {
		Node
		IsExpr()
		Accept(ExprVisitor)
	}
	BinaryOp interface {
		Expr
		LeftArg() Expr
		RightArg() Expr
	}
	UnaryOp interface {
		Expr
		Arg() Expr
	}
)

func BinaryOpToString(op BinaryOp, opName string) string {
	return fmt.Sprintf("%s %s %s", op.LeftArg(), opName, op.RightArg())
}

//go:generate marker -method IsNode,IsExpr -type OrExpr,AndExpr,XorExpr,NotExpr,BoolPrimaryComparison,BoolPrimaryPredicate,Exprs,PredicateIn,PredicateBetween,PredicateLike,PredicateBitExpr,BitExprBitOp,BitExprArtOp,BitExprSimpleExpr,SimpleExprPrefixOp,SimpleExprLit,Ident,FunctionCall,SimpleExprExpr -output expr_marker_generated.go

type (
	OrExpr struct {
		Left  Expr `json:"or_left"`
		Right Expr `json:"or_right"`
	}
	AndExpr struct {
		Left  Expr `json:"and_left"`
		Right Expr `json:"and_right"`
	}
	XorExpr struct {
		Left  Expr `json:"xor_left"`
		Right Expr `json:"xor_right"`
	}
	NotExpr struct {
		Expr Expr `json:"not_expr"`
	}
)

func (s *OrExpr) LeftArg() Expr   { return s.Left }
func (s *OrExpr) RightArg() Expr  { return s.Right }
func (s *AndExpr) LeftArg() Expr  { return s.Left }
func (s *AndExpr) RightArg() Expr { return s.Right }
func (s *XorExpr) LeftArg() Expr  { return s.Left }
func (s *XorExpr) RightArg() Expr { return s.Right }
func (s *NotExpr) Arg() Expr      { return s.Expr }

func (s *OrExpr) String() string  { return BinaryOpToString(s, "or") }
func (s *AndExpr) String() string { return BinaryOpToString(s, "and") }
func (s *XorExpr) String() string { return BinaryOpToString(s, "xor") }
func (s *NotExpr) String() string { return fmt.Sprintf("not %s", s.Expr) }

type ComparisonType int

//go:generate stringer -type ComparisonType -output comparison_type_stringer_generated.go

const (
	CmpEqual ComparisonType = iota
	CmpNotEqual
	CmpGreaterThan
	CmpGreaterEqual
	CmpLessThan
	CmpLessEqual
)

func (s ComparisonType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, s.String())), nil
}

func (s ComparisonType) Readable() string {
	switch s {
	case CmpEqual:
		return "="
	case CmpNotEqual:
		return "<>"
	case CmpGreaterThan:
		return ">"
	case CmpGreaterEqual:
		return ">="
	case CmpLessThan:
		return "<"
	case CmpLessEqual:
		return "<="
	default:
		panic(fmt.Sprintf("unknown comparison type %d", s))
	}
}

type (
	BoolPrimary interface {
		Expr
		IsBoolPrimary()
	}
)

//go:generate marker -method IsBoolPrimary -type BoolPrimaryComparison,BoolPrimaryPredicate -output bool_primary_marker_generated.go

type (
	BoolPrimaryComparison struct {
		Op    ComparisonType `json:"cmp_op"`
		Left  BoolPrimary    `json:"left,omitempty"`
		Right Predicate      `json:"right,omitempty"`
	}
	BoolPrimaryPredicate struct {
		Pred Predicate `json:"pred"`
	}
)

func (s *BoolPrimaryComparison) LeftArg() Expr  { return s.Left }
func (s *BoolPrimaryComparison) RightArg() Expr { return s.Right }

func (s *BoolPrimaryComparison) String() string { return BinaryOpToString(s, s.Op.Readable()) }
func (s *BoolPrimaryPredicate) String() string  { return s.Pred.String() }

type (
	Predicate interface {
		Expr
		IsPredicate()
	}
)

//go:generate marker -method IsPredicate -type PredicateIn,PredicateBetween,PredicateLike,PredicateBitExpr -output predicate_marker_generated.go

type (
	Exprs struct {
		Exprs []Expr `json:"exprs,omitempty"`
	}
	PredicateIn struct {
		IsNot  bool    `json:"in_not,omitempty"`
		Target BitExpr `json:"in_target"`
		List   *Exprs  `json:"in_list"`
	}
	PredicateBetween struct {
		IsNot  bool      `json:"between_not,omitempty"`
		Target BitExpr   `json:"between_target"`
		Left   BitExpr   `json:"between_lower"`
		Right  Predicate `json:"between_upper"`
	}
	PredicateLike struct {
		IsNot   bool       `json:"like_not,omitempty"`
		Target  BitExpr    `json:"like_target"`
		Pattern SimpleExpr `json:"like_pattern"`
	}
	PredicateBitExpr struct {
		Expr BitExpr `json:"bit_expr"`
	}
)

func (s *Exprs) String() string {
	b := buf.NewStrings()
	for _, x := range s.Exprs {
		b.Add(x.String())
	}
	return strings.Join(b.Get(), ", ")
}

func (s *PredicateIn) String() string {
	b := buf.NewStrings()
	b.Add(s.Target.String())
	if s.IsNot {
		b.Add("not")
	}
	b.Add("in")
	b.Add(fmt.Sprintf("(%s)", s.List))
	return strings.Join(b.Get(), " ")
}

func (s *PredicateBetween) String() string {
	b := buf.NewStrings()
	b.Add(s.Target.String())
	if s.IsNot {
		b.Add("not")
	}
	b.Add("between")
	b.Add(s.Left.String())
	b.Add("and")
	b.Add(s.Right.String())
	return strings.Join(b.Get(), " ")
}

func (s *PredicateLike) String() string {
	b := buf.NewStrings()
	b.Add(s.Target.String())
	if s.IsNot {
		b.Add("not")
	}
	b.Add("like")
	b.Add(s.Pattern.String())
	return strings.Join(b.Get(), " ")
}

func (s *PredicateBitExpr) String() string {
	return s.Expr.String()
}

type (
	ArithmeticOperatorType int
	BitOperatorType        int
)

//go:generate stringer -type ArithmeticOperatorType -output arithmetic_operator_type_stringer_generated.go
//go:generate stringer -type BitOperatorType -output bit_operator_type_stringer_generated.go

const (
	ArtOpAdd ArithmeticOperatorType = iota
	ArtOpSubtract
	ArtOpMultiply
	ArtOpDivide
)

const (
	BitOpAnd BitOperatorType = iota
	BitOpOr
	BitOpXor
)

func (s ArithmeticOperatorType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, s.String())), nil
}

func (s ArithmeticOperatorType) Readable() string {
	switch s {
	case ArtOpAdd:
		return "+"
	case ArtOpSubtract:
		return "-"
	case ArtOpMultiply:
		return "*"
	case ArtOpDivide:
		return "/"
	default:
		panic(fmt.Sprintf("unknown art op type %d", s))
	}
}

func (s BitOperatorType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, s.String())), nil
}

func (s BitOperatorType) Readable() string {
	switch s {
	case BitOpAnd:
		return "&"
	case BitOpOr:
		return "|"
	case BitOpXor:
		return "^"
	default:
		panic(fmt.Sprintf("unknown bit op type %d", s))
	}
}

type (
	BitExpr interface {
		Expr
		IsBitExpr()
	}
)

//go:generate marker -method IsBitExpr -type BitExprBitOp,BitExprArtOp,BitExprSimpleExpr -output bit_expr_marker_generated.go

type (
	BitExprBitOp struct {
		Op    BitOperatorType `json:"bit_op"`
		Left  BitExpr         `json:"left"`
		Right BitExpr         `json:"right"`
	}
	BitExprArtOp struct {
		Op    ArithmeticOperatorType `json:"art_op"`
		Left  BitExpr                `json:"left"`
		Right BitExpr                `json:"right"`
	}
	BitExprSimpleExpr struct {
		Expr SimpleExpr `json:"simple_expr,omitempty"`
	}
)

func (s *BitExprBitOp) LeftArg() Expr  { return s.Left }
func (s *BitExprBitOp) RightArg() Expr { return s.Right }
func (s *BitExprArtOp) LeftArg() Expr  { return s.Left }
func (s *BitExprArtOp) RightArg() Expr { return s.Right }

func (s *BitExprBitOp) String() string      { return BinaryOpToString(s, s.Op.Readable()) }
func (s *BitExprArtOp) String() string      { return BinaryOpToString(s, s.Op.Readable()) }
func (s *BitExprSimpleExpr) String() string { return s.Expr.String() }

type PrefixOperatorType int

//go:generate stringer -type PrefixOperatorType -output prefix_operator_stringer_generated.go

const (
	PreOpPlus PrefixOperatorType = iota
	PreOpMinus
	PreOpBitNot
	PreOpNot
)

func (s PrefixOperatorType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, s.String())), nil
}

func (s PrefixOperatorType) Readable() string {
	switch s {
	case PreOpPlus:
		return "+"
	case PreOpMinus:
		return "-"
	case PreOpBitNot:
		return "~"
	case PreOpNot:
		return "not"
	default:
		panic(fmt.Sprintf("unknown prefix op %d", s))
	}
}

type (
	SimpleExpr interface {
		Expr
		IsSimpleExpr()
	}
)

//go:generate marker -method IsSimpleExpr -type SimpleExprPrefixOp,SimpleExprLit,Ident,FunctionCall,SimpleExprExpr -output simple_expr_marker_generated.go

type (
	SimpleExprPrefixOp struct {
		Op   PrefixOperatorType `json:"pre_op"`
		Expr SimpleExpr         `json:"expr"`
	}
	SimpleExprLit struct {
		Lit Lit `json:"lit"`
	}
	Ident struct {
		Value string `json:"ident"`
	}
	FunctionCall struct {
		FunctionName *Ident `json:"function_name"`
		Arguments    *Exprs `json:"args,omitempty"`
	}
	SimpleExprExpr struct {
		Expr Expr `json:"expr,omitempty"`
	}
)

func (s *SimpleExprPrefixOp) Arg() Expr { return s.Expr }

func (s *SimpleExprPrefixOp) String() string { return fmt.Sprintf("%s%s", s.Op.Readable(), s.Expr) }
func (s *SimpleExprLit) String() string      { return s.Lit.String() }
func (s *Ident) String() string              { return s.Value }
func (s *FunctionCall) String() string       { return fmt.Sprintf("%s(%s)", s.FunctionName, s.Arguments) }
func (s *SimpleExprExpr) String() string     { return s.Expr.String() }
