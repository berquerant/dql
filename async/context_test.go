package async_test

import (
	"context"
	"testing"

	"github.com/berquerant/dql/async"
	"github.com/stretchr/testify/assert"
)

func TestIsDone(t *testing.T) {
	t.Run("continue", func(t *testing.T) {
		assert.False(t, async.IsDone(context.TODO()))
	})

	t.Run("canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.TODO())
		cancel()
		assert.True(t, async.IsDone(ctx))
	})
}
