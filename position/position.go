package position

import "fmt"

type (
	Position interface {
		Line() int
		Column() int
		Offset() int
		AddOffset(x int) Position
		AddLine(x int) Position
		AddColumn(x int) Position
		Clone() Position
	}

	position struct {
		line, column, offset int
	}
)

func New(line, column, offset int) Position {
	return &position{
		line:   line,
		column: column,
		offset: offset,
	}
}

func (s *position) Clone() Position { return s.clone() }

func (s *position) clone() *position {
	return &position{
		line:   s.line,
		column: s.column,
		offset: s.offset,
	}
}

func (s *position) AddOffset(x int) Position {
	p := s.clone()
	p.offset += x
	return p
}

func (s *position) AddColumn(x int) Position {
	p := s.clone()
	p.column += x
	return p
}

func (s *position) AddLine(x int) Position {
	p := s.clone()
	p.line += x
	return p
}

func (s *position) Line() int   { return s.line }
func (s *position) Column() int { return s.column }
func (s *position) Offset() int { return s.offset }
func (s *position) String() string {
	return fmt.Sprintf("line %d column %d at %d bytes", s.line, s.column, s.offset)
}
