package calc

import (
	"fmt"

	"github.com/berquerant/dql/ast"
	"github.com/berquerant/dql/cc"
	"github.com/berquerant/dql/logger"
)

type Debugger struct {
	Lexer      cc.Lexer
	Calculator Calculator
	Verbose    int
}

func (s *Debugger) Init() {
	if s.Verbose >= 0 {
		logger.SetDebugFlags()
		s.Lexer.Debug(s.Verbose)
	}
}

func (s *Debugger) Calc() {
	s.Init()

	status := cc.Parse(s.Lexer)
	logger.Info("[debugger] parser exit with %d", status)
	if err := s.Lexer.Err(); err != nil {
		logger.Error(" %v", err)
		return
	}
	stmt := s.Lexer.Result().(*ast.Statement)
	expr := stmt.SelectSection.Terms.Terms[0].Target.Expr
	logger.Info("[debugger] expr %s", logger.JSON(expr))
	v, err := s.Calculator.Data(expr)
	if err != nil {
		logger.Error("[debugger] failed to calc result %s %v", logger.JSON(expr), err)
		return
	}
	fmt.Printf("%s\n", logger.JSON(v))
}
