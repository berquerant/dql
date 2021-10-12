package cc

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/berquerant/dql/ast"
	"github.com/berquerant/dql/errors"
	"github.com/berquerant/dql/logger"
	"github.com/berquerant/dql/position"
	"github.com/berquerant/dql/token"
)

const EOF = -1

type Lexer interface {
	// Err returns an error caused during scanning.
	Err() error
	// Scan scans the next token.
	// This returns a token type defined by dql.y or EOF if reached EOF or caused some error.
	Scan() int
	// Buffer returns the value of the token buffered by scanning.
	Buffer() string
	// ResetBuffer clears buffer.
	ResetBuffer()
	// Debug enables debug logs.
	// level is yydebug value. If negative level, noop.
	Debug(level int)

	/* Utility methods. Expect called from dql.y. */

	// ParseInt parses a string as an integer.
	// Returns 0 and reports the error if failed.
	ParseInt(x string) int
	// ParseFloat parses a string as a floating point.
	// Returns 0 and reports the error if failed.
	ParseFloat(x string) float64

	/* Some consts in dql.y to ast consts translations. Expect called from dql.y. */

	AsPrefixOperatorType(op int) ast.PrefixOperatorType
	AsBitOperatorType(op int) ast.BitOperatorType
	AsArithmeticOperatorType(op int) ast.ArithmeticOperatorType
	AsComparisonType(op int) ast.ComparisonType

	/* Implements yyLexer. */

	Lex(lval *yySymType) int
	Error(msg string) // yyerror

	/* Keeps the result of pasing. */

	// SetResult stores the result of parsing.
	// Expect called from dql.y.
	SetResult(result ast.Node)
	// Result returns the stored result.
	Result() ast.Node
}

type lexer struct {
	position position.Position
	reader   *bufio.Reader
	buf      bytes.Buffer
	result   ast.Node
	err      error
	isDebug  bool
}

func NewLexer(r io.Reader) Lexer {
	yyErrorVerbose = true // YYERROR_VERBOSE
	return &lexer{
		position: position.New(1, 0, 0),
		reader:   bufio.NewReader(r),
	}
}

func (s *lexer) Result() ast.Node          { return s.result }
func (s *lexer) SetResult(result ast.Node) { s.result = result }
func (s *lexer) Err() error                { return s.err }
func (s *lexer) Debug(level int) {
	if level < 0 {
		return
	}
	s.isDebug = true
	yyDebug = level // YYDEBUG
}

func (s *lexer) Error(msg string) {
	s.err = errors.New("[lex] %s at %s", msg, s.position)
	s.debugf("%v", s.err)
}

func (s *lexer) errorf(format string, v ...interface{}) { s.Error(fmt.Sprintf(format, v...)) }
func (s *lexer) debugf(format string, v ...interface{}) {
	if s.isDebug {
		x := fmt.Sprintf(format, v...)
		logger.Debug("[lex][%s][%q] %s", s.position, s.buf.String(), x)
	}
}

func (s *lexer) AsComparisonType(op int) ast.ComparisonType {
	switch op {
	case EQ:
		return ast.CmpEqual
	case NE:
		return ast.CmpNotEqual
	case GT:
		return ast.CmpGreaterThan
	case GQ:
		return ast.CmpGreaterEqual
	case LT:
		return ast.CmpLessThan
	case LQ:
		return ast.CmpLessEqual
	}
	s.errorf("cannot translate comparison type %d", op)
	return ast.CmpEqual
}

func (s *lexer) AsArithmeticOperatorType(op int) ast.ArithmeticOperatorType {
	switch op {
	case PLUS:
		return ast.ArtOpAdd
	case MINUS:
		return ast.ArtOpSubtract
	case AST:
		return ast.ArtOpMultiply
	case SLASH:
		return ast.ArtOpDivide
	}
	s.errorf("cannot translate arithmetic operator %d", op)
	return ast.ArtOpAdd
}

func (s *lexer) AsBitOperatorType(op int) ast.BitOperatorType {
	switch op {
	case AMP:
		return ast.BitOpAnd
	case PIPE:
		return ast.BitOpOr
	case HAT:
		return ast.BitOpXor
	}
	s.errorf("cannot translate bit operator %d", op)
	return ast.BitOpAnd
}

func (s *lexer) AsPrefixOperatorType(op int) ast.PrefixOperatorType {
	switch op {
	case PLUS:
		return ast.PreOpPlus
	case MINUS:
		return ast.PreOpMinus
	case TILDE:
		return ast.PreOpBitNot
	case NOT:
		return ast.PreOpNot
	}
	s.errorf("cannot translate prefix operator %d", op)
	return ast.PreOpPlus
}

func (s *lexer) ParseInt(x string) int {
	r, err := ParseInt(x)
	if err != nil {
		s.errorf("cannot parse %s as int %v", x, err)
		return 0
	}
	return r
}

func (s *lexer) ParseFloat(x string) float64 {
	r, err := ParseFloat(x)
	if err != nil {
		s.errorf("cannot parse %s as float %v", x, err)
		return 0
	}
	return r
}

type scannedDigitType int

const (
	scannedDigitUnknown = iota
	scannedDigitInt
	scannedDigitFloat
)

func (s *lexer) scanDigit() scannedDigitType {
	if !IsDigit(s.Peek()) {
		return scannedDigitUnknown
	}
	for x := s.Peek(); IsDigit(x); x = s.Peek() {
		_ = s.Next()
	}
	if s.Peek() != '.' {
		return scannedDigitInt
	}
	// Read floating point
	_ = s.Next()
	for x := s.Peek(); IsDigit(x); x = s.Peek() {
		_ = s.Next()
	}
	return scannedDigitFloat
}

func (s *lexer) scanIdent() {
	x := s.Peek()
	if !IsIdentHead(x) {
		s.errorf("unexpected rune for ident %q", x)
		return
	}
	_ = s.Next()
	for x = s.Peek(); IsIdentTail(x); x = s.Peek() {
		_ = s.Next()
	}
}

// scanSpaces skips spaces.
// Returns true if skipped spaces exist.
func (s *lexer) scanSpaces() bool {
	read := false
	for x := s.Peek(); IsSpace(x); x = s.Peek() {
		_ = s.Next()
		read = true
	}
	return read
}

// scanString reads a token like 'ground' (when stop is ').
// When this called, firstly it peeks the next rune of the first stop rune, like `g` for 'ground'.
func (s *lexer) scanString(stop rune) {
	for {
		switch s.Peek() {
		case EOF:
			s.errorf("unclosed string, expect %v but reached EOF", stop)
			return
		case stop:
			s.Discard()
			return
		case '\\':
			s.Discard()
			switch s.Peek() {
			case 'b':
				s.Discard()
				s.buf.WriteRune('\b')
			case 'n':
				s.Discard()
				s.buf.WriteRune('\n')
			case 'r':
				s.Discard()
				s.buf.WriteRune('\r')
			case 't':
				s.Discard()
				s.buf.WriteRune('\t')
			}
		}
		if unicode.IsControl(s.Peek()) {
			s.Error("string cannot accept control characters")
			return
		}
		if s.Peek() == EOF {
			s.errorf("unclosed string, expect %v but reached EOF", stop)
			return
		}
		_ = s.Next()
	}
}

func (s *lexer) Scan() int {
	switch s.Peek() {
	case EOF:
		return EOF
	case '\'':
		s.Discard()
		s.scanString('\'')
		return STRING
	case '"':
		s.Discard()
		s.scanString('"')
		return STRING
	case ',':
		_ = s.Next()
		return COMMA
	case ';':
		_ = s.Next()
		return SCOLON
	case '(':
		_ = s.Next()
		return LPAR
	case ')':
		_ = s.Next()
		return RPAR
	case '+':
		_ = s.Next()
		return PLUS
	case '-':
		_ = s.Next()
		return MINUS
	case '*':
		_ = s.Next()
		return AST
	case '/':
		_ = s.Next()
		return SLASH
	case '&':
		_ = s.Next()
		return AMP
	case '|':
		_ = s.Next()
		return PIPE
	case '^':
		_ = s.Next()
		return HAT
	case '~':
		_ = s.Next()
		return TILDE
	case '=':
		_ = s.Next()
		return EQ
	case '<':
		_ = s.Next()
		switch s.Peek() {
		case '>':
			_ = s.Next()
			return NE
		case '=':
			_ = s.Next()
			return LQ
		default:
			return LT
		}
	case '>':
		_ = s.Next()
		switch s.Peek() {
		case '=':
			_ = s.Next()
			return GQ
		default:
			return GT
		}
	}

	// skip spaces
	if s.scanSpaces() {
		s.ResetBuffer()
		return s.Scan()
	}

	switch s.scanDigit() {
	case scannedDigitInt:
		return INT
	case scannedDigitFloat:
		return FLOAT
	}

	s.scanIdent()
	// keywords are case-insensitive.
	switch strings.ToLower(s.Buffer()) {
	case "select":
		return SELECT
	case "distinct":
		return DISTINCT
	case "where":
		return WHERE
	case "having":
		return HAVING
	case "group":
		return GROUP
	case "by":
		return BY
	case "order":
		return ORDER
	case "limit":
		return LIMIT
	case "as":
		return AS
	case "asc":
		return ASC
	case "desc":
		return DESC
	case "like":
		return LIKE
	case "in":
		return IN
	case "not":
		return NOT
	case "and":
		return AND
	case "or":
		return OR
	case "xor":
		return XOR
	case "between":
		return BETWEEN
	case "offset":
		return OFFSET
	}
	return IDENT
}

func (s *lexer) Buffer() string { return s.buf.String() }
func (s *lexer) ResetBuffer()   { s.buf.Reset() }

// Discard skips the next character.
func (s *lexer) Discard() rune {
	r, size, err := s.reader.ReadRune()
	s.debugf("[Discard] %v %d %v", r, size, err)
	if err != nil {
		if err != io.EOF {
			s.errorf("[Discard] from reader %v", err)
		}
		return EOF
	}
	if r == '\n' {
		s.position = position.New(s.position.Line()+1, 0, s.position.Offset())
	}
	s.position = s.position.AddColumn(size).AddOffset(size)
	return r
}

// Peek peeks the next character.
func (s *lexer) Peek() rune {
	r, _, err := s.reader.ReadRune()
	s.debugf("[Peek] %v %v", r, err)
	if err != nil {
		if err != io.EOF {
			s.errorf("[Peek] from reader %v", err)
		}
		return EOF
	}
	if err := s.reader.UnreadRune(); err != nil {
		s.errorf("[Peek] failed to unread %v", err)
		return EOF
	}
	return r
}

// Next reads the next character.
func (s *lexer) Next() rune {
	r, size, err := s.reader.ReadRune()
	s.debugf("[Next] %v %d %v", r, size, err)
	if err != nil {
		if err != io.EOF {
			s.errorf("[Next] from reader %v", err)
		}
		return EOF
	}
	if r == '\n' {
		s.position = position.New(s.position.Line()+1, 0, s.position.Offset())
	}
	s.position = s.position.AddColumn(size).AddOffset(size)
	if _, err := s.buf.WriteRune(r); err != nil {
		s.errorf("[Next] failed to write buffer %v", err)
		return EOF
	}
	return r
}

func (s *lexer) Lex(lval *yySymType) int {
	if s.err != nil {
		return EOF
	}
	t := s.Scan()
	v := s.Buffer()
	lval.token = token.New(t, v)
	s.ResetBuffer()
	return t
}
