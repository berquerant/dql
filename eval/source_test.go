package eval_test

import (
	"context"
	"errors"
	"io/fs"
	"testing"
	"time"

	"github.com/berquerant/dql/dig"
	"github.com/berquerant/dql/eval"
	"github.com/stretchr/testify/assert"
)

type mockFileInfo struct {
	name    string
	size    int64
	mode    fs.FileMode
	modTime time.Time
	isDir   bool
}

func (s *mockFileInfo) Name() string       { return s.name }
func (s *mockFileInfo) Size() int64        { return s.size }
func (s *mockFileInfo) Mode() fs.FileMode  { return s.mode }
func (s *mockFileInfo) ModTime() time.Time { return s.modTime }
func (s *mockFileInfo) IsDir() bool        { return s.isDir }

type mockDigger struct {
	infos []dig.FileInfo
	err   error
}

func (s *mockDigger) Dig(_ string, handler dig.FileInfoHandler) error {
	if s.err != nil {
		return s.err
	}
	for _, x := range s.infos {
		if handler(x) == dig.InstrCancel {
			break
		}
	}
	return nil
}

func newFileInfos(names ...string) []dig.FileInfo {
	r := make([]dig.FileInfo, len(names))
	for i, n := range names {
		r[i] = &mockFileInfo{
			name: n,
		}
	}
	return r
}

func resultToRows(resultC <-chan eval.Row) []eval.Row {
	r := []eval.Row{}
	for v := range resultC {
		r = append(r, v)
	}
	return r
}

func TestSource(t *testing.T) {
	t.Run("err", func(t *testing.T) {
		errSource := errors.New("error source")
		got := resultToRows(eval.NewSource(&mockDigger{
			err: errSource,
		}).Yield(context.TODO(), ""))
		assert.Equal(t, 1, len(got))
		assert.ErrorIs(t, got[0].Err(), errSource)
	})

	t.Run("cancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.TODO())
		cancel()
		got := resultToRows(eval.NewSource(&mockDigger{
			infos: newFileInfos("a"),
		}).Yield(ctx, ""))
		assert.Equal(t, 1, len(got))
		assert.ErrorIs(t, got[0].Err(), context.Canceled)
	})

	t.Run("yield", func(t *testing.T) {
		got := resultToRows(eval.NewSource(&mockDigger{
			infos: newFileInfos("a", "b"),
		}).Yield(context.TODO(), ""))
		assert.Equal(t, 2, len(got))
		assert.Equal(t, "a", got[0].Info().Name())
		assert.Equal(t, "b", got[1].Info().Name())
	})
}
