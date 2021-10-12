package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/berquerant/dql"
	"github.com/berquerant/dql/ast"
	"github.com/berquerant/dql/cc"
	"github.com/berquerant/dql/eval"
	"github.com/berquerant/dql/logger"
)

var (
	verbose   = flag.Int("v", -1, "Verbose logging level. Enable debug logs if not negative level.")
	asJSON    = flag.Bool("j", false, "Print result as json.")
	noHeaders = flag.Bool("H", false, "Print no header line.")
)

const usage = `Usage of sql:
  dql QUERY files... directory...
Flags:`

func Usage() {
	fmt.Fprintln(os.Stderr, dql.Meta())
	fmt.Fprintln(os.Stderr, usage)
	flag.PrintDefaults()
}

func main() {
	flag.Usage = Usage
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		flag.Usage()
		os.Exit(2)
	}

	var (
		query   = args[0]
		targets = args[1:]
	)

	lexer := cc.NewLexer(strings.NewReader(query))
	if status := cc.Parse(lexer); status != 0 {
		logger.Error("failed parser; exit status %d", status)
		os.Exit(status)
	}
	if err := lexer.Err(); err != nil {
		logger.Error("lexer got error %v", err)
		os.Exit(1)
	}
	stmt := lexer.Result().(*ast.Statement)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	err := printResult(ctx, eval.NewRunner(stmt), targets)
	stop()
	if err != nil {
		logger.Error("%v", err)
		os.Exit(1)
	}
}

func printResult(ctx context.Context, runner eval.Runner, targets []string) error {
	if *asJSON {
		return NewJSONWriter(runner, targets).Write(ctx, os.Stdout)
	}
	return NewCSVWriter(runner, targets, *noHeaders).Write(ctx, os.Stdout)
}
