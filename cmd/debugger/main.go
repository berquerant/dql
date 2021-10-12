package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/berquerant/dql/calc"
	"github.com/berquerant/dql/cc"
	"github.com/berquerant/dql/data"
	"github.com/berquerant/dql/env"
	"github.com/berquerant/dql/eval"
	"github.com/berquerant/dql/logger"
)

var (
	verbose   = flag.Int("v", -1, "Verbose logging level. Enable debug logs if not negative level.")
	doLex     = flag.Bool("l", false, "Print result of lex.")
	doParse   = flag.Bool("p", false, "Print result of parse.")
	doCalc    = flag.Bool("c", false, "Print result of calc expr. Stdin should be `select EXPR` if this specified.")
	envStr    = flag.String("env", "{}", "Environment for calc.")
	doEval    = flag.Bool("e", false, "Print result of eval.")
	doUnparse = flag.Bool("u", false, "Print unparsed query.")
)

const usage = `Usage of debugger:
  echo QUERY | debugger -e files... directory...  # eval query
  echo QUERY | debugger -c -env ENV  # calc expr, limited query
  echo QUERY | debugger -p  # parse QUERY
  echo QUERY | debugger -l  # lex QUERY
  echo QUERY | debugger -u  # parse and unparse QUERY
Flags:`

func Usage() {
	fmt.Fprintln(os.Stderr, usage)
	flag.PrintDefaults()
}

func main() {
	flag.Usage = Usage
	flag.Parse()
	if !(*doLex || *doParse || *doCalc || *doEval || *doUnparse) {
		fmt.Fprintf(os.Stderr, "Specify -l or -p or -c or -e or -u\n")
		flag.Usage()
		os.Exit(2)
	}
	now := time.Now()
	defer func() {
		logger.Info("[debugger] elapsed %v", time.Since(now))
	}()

	switch {
	case *doEval:
		debugger := &eval.Debugger{
			Lexer:     cc.NewLexer(os.Stdin),
			Verbose:   *verbose,
			FileNames: flag.Args(),
		}
		debugger.Eval()
	case *doCalc:
		table, err := parseEnv()
		if err != nil {
			panic(err)
		}
		debugger := &calc.Debugger{
			Lexer:      cc.NewLexer(os.Stdin),
			Verbose:    *verbose,
			Calculator: calc.NewAggregation(table),
		}
		debugger.Init()
		logger.Info("[debugger] env %s", logger.JSON(table))
		debugger.Calc()
	default:
		debugger := &cc.Debugger{
			Lexer:   cc.NewLexer(os.Stdin),
			Verbose: *verbose,
		}
		switch {
		case *doLex:
			debugger.Lex()
		case *doParse:
			debugger.Parse()
		case *doUnparse:
			debugger.Unparse()
		}
	}
}

func parseEnv() (env.Map, error) {
	var d map[string]interface{}
	if err := json.Unmarshal([]byte(*envStr), &d); err != nil {
		return nil, err
	}
	m := env.New()
	for k, v := range d {
		if reflect.TypeOf(v).Kind() == reflect.Slice {
			t := reflect.ValueOf(v)
			vs := make([]data.Data, t.Len())
			for i := 0; i < t.Len(); i++ {
				if d, ok := data.FromInterface(t.Index(i).Interface()); ok {
					vs[i] = d
					continue
				}
				logger.Error("cannot accept %#v at %d in %s", t.Index(i).Interface(), i, logger.JSON(v))
			}
			m.Set(k, env.FromDataList(vs))
			continue
		}
		if d, ok := data.FromInterface(v); ok {
			m.Set(k, env.FromData(d))
			continue
		}
		logger.Error("cannot accept %#v", v)
	}
	return m, nil
}
