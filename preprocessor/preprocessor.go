package preprocessor

import "github.com/berquerant/dql/ast"

type (
	// PreProcessor provides an interface for tree translation.
	PreProcessor interface {
		PreProcess(stmt *ast.Statement) error
	}
)
