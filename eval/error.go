package eval

import "github.com/berquerant/dql/errors"

var (
	errFiltered            = errors.New("filtered")
	ErrNotBoolExpr         = errors.New("not bool expr")
	ErrInvalidIdentRef     = errors.New("invalid ident ref")
	ErrInvalidIdent        = errors.New("invalid ident")
	ErrInvalidHaving       = errors.New("invalid having")
	ErrUnknownRowType      = errors.New("unknown row type")
	ErrUnknownDataType     = errors.New("unknown data type")
	ErrInvalidLimit        = errors.New("invalid limit")
	ErrInvalidSelectSource = errors.New("invalid select source")
)
