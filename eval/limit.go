package eval

import (
	"context"

	"github.com/berquerant/dql/async"
	"github.com/berquerant/dql/errors"
)

type (
	Limit interface {
		Limit(ctx context.Context, limit, offset int, sourceC <-chan GRow) <-chan GRow
	}

	limit struct{}
)

func NewLimit() Limit { return &limit{} }

func (*limit) Limit(ctx context.Context, limit, offset int, sourceC <-chan GRow) <-chan GRow {
	resultC := make(chan GRow, resultCBufferSize)
	if limit < 1 || offset < 0 {
		resultC <- NewErrGRow(errors.Wrap(ErrInvalidLimit, "limit %d offset %d", limit, offset))
		close(resultC)
		return resultC
	}
	go func() {
		defer close(resultC)
		var (
			i, c int
		)
		for r := range sourceC {
			if async.IsDone(ctx) {
				resultC <- NewErrGRow(errors.Wrap(ctx.Err(), "limit"))
				return
			}
			if i >= offset && c < limit {
				resultC <- r
				c++
			}
			i++
		}
	}()
	return resultC
}
