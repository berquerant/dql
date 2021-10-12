package ast

//go:generate mkvisitor -type OrExpr,AndExpr,XorExpr,NotExpr,BoolPrimaryComparison,BoolPrimaryPredicate,Exprs,PredicateIn,PredicateBetween,PredicateLike,PredicateBitExpr,BitExprBitOp,BitExprArtOp,BitExprSimpleExpr,SimpleExprPrefixOp,SimpleExprLit,Ident,FunctionCall,SimpleExprExpr,IntLit,FloatLit,StringLit -vType ExprVisitor -output expr_mkvisitor_generated.go

type (
	// VisitorCallback is the callback function for BaseVisitor.
	// Cancel waking if returns false.
	VisitorCallback func(expr Expr) bool

	// BaseVisitor visits expr and invokes callback recursively.
	BaseVisitor interface {
		ExprVisitor
		Init()
	}

	baseVisitor struct {
		callback VisitorCallback
		isDone   bool
	}
)

func NewBaseVisitor(callback VisitorCallback) BaseVisitor {
	return &baseVisitor{
		callback: callback,
	}
}

func (s *baseVisitor) Init() { s.isDone = false }

func (s *baseVisitor) run(v interface{}) {
	if s.isDone {
		return
	}
	if !s.callback(v.(Expr)) {
		s.isDone = true
	}
}

func (s *baseVisitor) VisitOrExpr(v *OrExpr)                               { s.VisitBinaryOp(v) }
func (s *baseVisitor) VisitAndExpr(v *AndExpr)                             { s.VisitBinaryOp(v) }
func (s *baseVisitor) VisitXorExpr(v *XorExpr)                             { s.VisitBinaryOp(v) }
func (s *baseVisitor) VisitNotExpr(v *NotExpr)                             { s.VisitUnaryOp(v) }
func (s *baseVisitor) VisitBoolPrimaryComparison(v *BoolPrimaryComparison) { s.VisitBinaryOp(v) }
func (s *baseVisitor) VisitBoolPrimaryPredicate(v *BoolPrimaryPredicate) {
	s.run(v)
	s.visit(v.Pred)
}
func (s *baseVisitor) VisitExprs(v *Exprs) {
	s.run(v)
	for _, x := range v.Exprs {
		s.visit(x)
	}
}
func (s *baseVisitor) VisitPredicateIn(v *PredicateIn) {
	s.run(v)
	s.visit(v.Target)
	s.visit(v.List)
}
func (s *baseVisitor) VisitPredicateBetween(v *PredicateBetween) {
	s.run(v)
	s.visit(v.Target)
	s.visit(v.Left)
	s.visit(v.Right)
}
func (s *baseVisitor) VisitPredicateLike(v *PredicateLike) {
	s.run(v)
	s.visit(v.Target)
	s.visit(v.Pattern)
}
func (s *baseVisitor) VisitPredicateBitExpr(v *PredicateBitExpr) {
	s.run(v)
	s.visit(v.Expr)
}
func (s *baseVisitor) VisitBitExprBitOp(v *BitExprBitOp) { s.VisitBinaryOp(v) }
func (s *baseVisitor) VisitBitExprArtOp(v *BitExprArtOp) { s.VisitBinaryOp(v) }
func (s *baseVisitor) VisitBitExprSimpleExpr(v *BitExprSimpleExpr) {
	s.run(v)
	s.visit(v.Expr)
}
func (s *baseVisitor) VisitSimpleExprPrefixOp(v *SimpleExprPrefixOp) { s.VisitUnaryOp(v) }
func (s *baseVisitor) VisitSimpleExprLit(v *SimpleExprLit) {
	s.run(v)
	s.visit(v.Lit)
}
func (s *baseVisitor) VisitIdent(v *Ident) { s.run(v) }
func (s *baseVisitor) VisitFunctionCall(v *FunctionCall) {
	s.run(v)
	s.visit(v.FunctionName)
	s.visit(v.Arguments)
}
func (s *baseVisitor) VisitSimpleExprExpr(v *SimpleExprExpr) {
	s.run(v)
	s.visit(v.Expr)
}
func (s *baseVisitor) VisitIntLit(v *IntLit)       { s.run(v) }
func (s *baseVisitor) VisitFloatLit(v *FloatLit)   { s.run(v) }
func (s *baseVisitor) VisitStringLit(v *StringLit) { s.run(v) }

func (s *baseVisitor) VisitBinaryOp(v BinaryOp) {
	s.run(v)
	s.visit(v.LeftArg())
	s.visit(v.RightArg())
}

func (s *baseVisitor) VisitUnaryOp(v UnaryOp) {
	s.run(v)
	s.visit(v.Arg())
}

func (s *baseVisitor) visit(v Expr) { VisitSwitch(s, v) }
