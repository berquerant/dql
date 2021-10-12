package eval

import (
	"context"
	"fmt"

	"github.com/berquerant/dql/ast"
	"github.com/berquerant/dql/cc"
	"github.com/berquerant/dql/logger"
)

type Debugger struct {
	Lexer     cc.Lexer
	Verbose   int
	FileNames []string
}

func (s *Debugger) Init() {
	if s.Verbose >= 0 {
		logger.SetDebugFlags()
		s.Lexer.Debug(s.Verbose)
	}
}

func (s *Debugger) Eval() {
	s.Init()

	status := cc.Parse(s.Lexer)
	logger.Info("[debugger] parser exit with %d", status)
	if err := s.Lexer.Err(); err != nil {
		logger.Error("%v", err)
		return
	}
	stmt := s.Lexer.Result().(*ast.Statement)
	for r := range NewRunner(stmt).Run(context.Background(), s.FileNames...) {
		if err := r.Err(); err != nil {
			logger.Error("%v", err)
			return
		}
		fmt.Printf("%s\n", logger.JSON(r))
	}
}
