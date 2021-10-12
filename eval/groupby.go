package eval

import (
	"context"

	"github.com/berquerant/dql/ast"
	"github.com/berquerant/dql/async"
	"github.com/berquerant/dql/data"
	"github.com/berquerant/dql/env"
	"github.com/berquerant/dql/errors"
	"github.com/berquerant/dql/logger"
)

type (
	GroupBy interface {
		Group(ctx context.Context, table env.Map, sourceC <-chan Row) <-chan GRow
	}

	groupByNoop struct{}

	groupByKey struct {
		key string
	}
)

func NewGroupBy(key string) GroupBy {
	if key == "" {
		return &groupByNoop{}
	}
	return &groupByKey{
		key: key,
	}
}

func (s *groupByKey) getKey(table env.Map) (string, error) {
	v, ok := table.Get(s.key)
	if !ok {
		return s.key, nil
	}
	switch v.Type() {
	case env.TypeExpr:
		if ident, ok := v.Expr().(*ast.Ident); ok {
			return ident.Value, nil
		}
		return "", errors.Wrap(ErrInvalidIdentRef, "group key %s %s", s.key, logger.JSON(v.Expr()))
	case env.TypeData:
		return "", errors.Wrap(ErrInvalidIdentRef, "group key %s %s", s.key, logger.JSON(v.Data()))
	default:
		return "", errors.Wrap(ErrInvalidIdentRef, "group key %s %s", s.key, logger.JSON(v))
	}
}

func (s *groupByKey) Group(ctx context.Context, table env.Map, sourceC <-chan Row) <-chan GRow {
	resultC := make(chan GRow, resultCBufferSize)
	key, err := s.getKey(table)
	if err != nil {
		resultC <- NewErrGRow(err)
		close(resultC)
		return resultC
	}

	go func() {
		defer close(resultC)
		d := map[interface{}][]Row{}
		for r := range sourceC {
			if async.IsDone(ctx) {
				resultC <- NewErrGRow(errors.Wrap(ctx.Err(), "group by"))
				return
			}
			if err := r.Err(); err != nil {
				resultC <- NewErrGRow(errors.Wrap(err, "group by"))
				return
			}
			m := r.Info().ToMap()
			k, exist := m[key]
			if !exist {
				resultC <- NewErrGRow(errors.Wrap(ErrInvalidIdent, "group by key %s", key))
				return
			}
			kv := k.Value()
			if _, ok := d[kv]; !ok {
				d[kv] = []Row{}
			}
			d[kv] = append(d[kv], r)
		}
		for k, rows := range d {
			kk, _ := data.FromInterface(k)
			resultC <- NewGroupedGRow(NewGroupedRow(key, kk, rows))
		}
	}()
	return resultC
}

func (*groupByNoop) Group(ctx context.Context, _ env.Map, sourceC <-chan Row) <-chan GRow {
	resultC := make(chan GRow, 1000)
	go func() {
		defer close(resultC)
		for r := range sourceC {
			if async.IsDone(ctx) {
				resultC <- NewErrGRow(errors.Wrap(ctx.Err(), "group by (noop)"))
				return
			}
			if err := r.Err(); err != nil {
				resultC <- NewErrGRow(errors.Wrap(err, "group by (noop)"))
				return
			}
			resultC <- NewRawGRow(r)
		}
	}()
	return resultC
}
