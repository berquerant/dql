package eval

import (
	"context"

	"github.com/berquerant/dql/async"
	"github.com/berquerant/dql/dig"
	"github.com/berquerant/dql/errors"
)

type (
	Source interface {
		Yield(ctx context.Context, names ...string) <-chan Row
	}

	source struct {
		digger dig.Digger
	}
)

func NewSource(digger dig.Digger) Source {
	return &source{
		digger: digger,
	}
}

func (s *source) Yield(ctx context.Context, names ...string) <-chan Row {
	resultC := make(chan Row, resultCBufferSize)
	go func() {
		defer close(resultC)
		for _, name := range names {
			if err := s.digger.Dig(name, func(v dig.FileInfo) dig.Instr {
				if async.IsDone(ctx) {
					resultC <- NewErrRow(errors.Wrap(ctx.Err(), "yield"))
					return dig.InstrCancel
				}
				resultC <- NewRow(NewInfo(v))
				return dig.InstrContinue
			}); err != nil {
				resultC <- NewErrRow(errors.Wrap(err, "yield"))
				return
			}
		}
	}()
	return resultC
}
