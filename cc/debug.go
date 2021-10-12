package cc

import (
	"encoding/json"
	"fmt"

	"github.com/berquerant/dql/logger"
	"github.com/berquerant/dql/token"
)

type Debugger struct {
	Lexer   Lexer
	Verbose int
}

func (s *Debugger) Init() {
	if s.Verbose >= 0 {
		logger.SetDebugFlags()
		s.Lexer.Debug(s.Verbose)
	}
}

func (s *Debugger) Lex() {
	s.Init()

	for {
		t := s.Lexer.Scan()
		if err := s.Lexer.Err(); err != nil {
			logger.Error("[debugger] %v", err)
			return
		}
		if t == EOF {
			logger.Info("[debugger] lexer finished successfully")
			return
		}
		tok := token.New(t, s.Lexer.Buffer())
		s.Lexer.ResetBuffer()
		v, err := json.Marshal(tok)
		if err != nil {
			logger.Error("[debugger] failed to marshal token %s %v", tok, err)
			continue
		}
		logger.Info("[debugger][token] %s", v)
	}
}

func (s *Debugger) Parse() {
	s.Init()

	status := Parse(s.Lexer)
	logger.Info("[debugger] parser exit with %d", status)
	if err := s.Lexer.Err(); err != nil {
		logger.Error(" %v", err)
		return
	}
	fmt.Printf("%s\n", logger.JSON(s.Lexer.Result()))
}

func (s *Debugger) Unparse() {
	s.Init()
	status := Parse(s.Lexer)
	logger.Info("[debugger] parser exit with %d", status)
	if err := s.Lexer.Err(); err != nil {
		logger.Error(" %v", err)
		return
	}
	fmt.Println(s.Lexer.Result().String())
}
