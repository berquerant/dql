package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"

	"github.com/berquerant/dql/errors"
	"github.com/berquerant/dql/eval"
)

type ResultWriter interface {
	Write(ctx context.Context, w io.Writer) error
}

func NewJSONWriter(runner eval.Runner, targets []string) ResultWriter {
	return &jsonWriter{
		runner:  runner,
		targets: targets,
	}
}

type jsonWriter struct {
	runner  eval.Runner
	targets []string
}

func (s *jsonWriter) Write(ctx context.Context, w io.Writer) error {
	for r := range s.runner.Run(ctx, s.targets...) {
		if err := r.Err(); err != nil {
			return err
		}
		if len(s.runner.Headers()) != r.Len() {
			return errors.New(fmt.Sprintf("%d headers but got row %d columns", len(s.runner.Headers()), r.Len()))
		}
		d := make(map[string]interface{}, r.Len())
		for i, h := range s.runner.Headers() {
			d[h] = r.Get(i).Value()
		}
		b, err := json.Marshal(d)
		if err != nil {
			return errors.Wrap(err, "row %#v", d)
		}
		fmt.Fprintf(w, "%s\n", b)
	}
	return nil
}

func NewCSVWriter(runner eval.Runner, targets []string, noHeaders bool) ResultWriter {
	return &csvWriter{
		runner:    runner,
		targets:   targets,
		noHeaders: noHeaders,
	}
}

type csvWriter struct {
	runner    eval.Runner
	targets   []string
	noHeaders bool
}

func (s *csvWriter) Write(ctx context.Context, w io.Writer) error {
	writer := csv.NewWriter(w)
	if !s.noHeaders {
		if err := writer.Write(s.runner.Headers()); err != nil {
			return errors.Wrap(err, "header")
		}
	}
	for r := range s.runner.Run(ctx, s.targets...) {
		if err := r.Err(); err != nil {
			return err
		}
		if len(s.runner.Headers()) != r.Len() {
			return errors.New(fmt.Sprintf("%d headers but got row %d columns", len(s.runner.Headers()), r.Len()))
		}
		values := make([]string, r.Len())
		for i := 0; i < r.Len(); i++ {
			values[i] = fmt.Sprintf("%v", r.Get(i).Value())
		}
		if err := writer.Write(values); err != nil {
			return errors.Wrap(err, "row %#v", values)
		}
	}
	writer.Flush()
	return writer.Error()
}
