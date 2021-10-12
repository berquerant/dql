package eval

type RowType int

//go:generate stringer -type RowType -output row_stringer_generated.go

const (
	ErrRowType RowType = iota
	RawRowType
	GroupedRowType
)

const (
	resultCBufferSize = 1000
)
