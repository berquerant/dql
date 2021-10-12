package preprocessor

import (
	"github.com/berquerant/dql/errors"

	"github.com/berquerant/dql/ast"
)

var (
	ErrAliasAll = errors.New("alias all")
)

type (
	selectAll struct {
		allSymbol string
	}
)

func NewSelectAll(allSymbol string) PreProcessor {
	return &selectAll{
		allSymbol: allSymbol,
	}
}

func (s *selectAll) PreProcess(stmt *ast.Statement) error {
	terms := []*ast.SelectTerm{}
	for _, t := range stmt.SelectSection.Terms.Terms {
		e := t.Target.Expr
		if id, ok := s.fetchIdent(e); ok {
			if id == s.allSymbol {
				if t.As != nil {
					return errors.Wrap(ErrAliasAll, "selectAll")
				}
				terms = append(terms, s.unzip()...)
				continue
			}
		}
		terms = append(terms, t)
	}
	stmt.SelectSection.Terms.Terms = terms
	return nil
}

func (s *selectAll) unzip() []*ast.SelectTerm {
	targets := []string{"name", "size", "mode", "mod_time", "is_dir"}
	terms := make([]*ast.SelectTerm, len(targets))
	for i, t := range targets {
		terms[i] = &ast.SelectTerm{
			Target: &ast.SelectTarget{
				Expr: s.identExpr(t),
			},
		}
	}
	return terms
}

func (*selectAll) fetchIdent(expr ast.Expr) (string, bool) {
	bp, ok := expr.(*ast.BoolPrimaryPredicate)
	if !ok {
		return "", false
	}
	pb, ok := bp.Pred.(*ast.PredicateBitExpr)
	if !ok {
		return "", false
	}
	se, ok := pb.Expr.(*ast.BitExprSimpleExpr)
	if !ok {
		return "", false
	}
	if id, ok := se.Expr.(*ast.Ident); ok {
		return id.Value, true
	}
	return "", false
}

func (*selectAll) identExpr(value string) ast.Expr {
	return &ast.BoolPrimaryPredicate{
		Pred: &ast.PredicateBitExpr{
			Expr: &ast.BitExprSimpleExpr{
				Expr: &ast.Ident{
					Value: value,
				},
			},
		},
	}
}
