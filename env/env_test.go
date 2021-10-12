package env_test

import (
	"testing"

	"github.com/berquerant/dql/data"
	"github.com/berquerant/dql/env"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	s := env.New()
	// key str => value null
	t.Run("set", func(t *testing.T) {
		s.Set("str", env.FromData(data.FromString("null")))
		t.Run("get", func(t *testing.T) {
			v, ok := s.Get("str")
			assert.True(t, ok)
			assert.Equal(t, env.TypeData, v.Type())
			assert.Equal(t, "null", v.Data().String())
		})
		t.Run("miss", func(t *testing.T) {
			_, ok := s.Get("int")
			assert.False(t, ok)
		})
	})
	// key int => key str
	t.Run("ref", func(t *testing.T) {
		s.Ref("int", "str")
		t.Run("follow ref", func(t *testing.T) {
			v, ok := s.Get("int")
			assert.True(t, ok)
			assert.Equal(t, env.TypeData, v.Type())
			assert.Equal(t, "null", v.Data().String())
		})
		t.Run("direct", func(t *testing.T) {
			v, ok := s.Get("str")
			assert.True(t, ok)
			assert.Equal(t, env.TypeData, v.Type())
			assert.Equal(t, "null", v.Data().String())
		})
	})
	// key str => value any
	t.Run("overwrite", func(t *testing.T) {
		s.Set("str", env.FromData(data.FromString("any")))
		t.Run("follow ref", func(t *testing.T) {
			v, ok := s.Get("int")
			assert.True(t, ok)
			assert.Equal(t, env.TypeData, v.Type())
			assert.Equal(t, "any", v.Data().String())
		})
		t.Run("direct", func(t *testing.T) {
			v, ok := s.Get("str")
			assert.True(t, ok)
			assert.Equal(t, env.TypeData, v.Type())
			assert.Equal(t, "any", v.Data().String())
		})
	})
	// key complex => key int
	t.Run("double ref", func(t *testing.T) {
		s.Ref("complex", "int")
		t.Run("cannot follow ref twice", func(t *testing.T) {
			_, ok := s.Get("complex")
			assert.False(t, ok)
		})
	})
	// key int => key str => value some
	t.Run("set by ref", func(t *testing.T) {
		s.Set("int", env.FromData(data.FromString("some")))
		t.Run("follow ref", func(t *testing.T) {
			v, ok := s.Get("int")
			assert.True(t, ok)
			assert.Equal(t, env.TypeData, v.Type())
			assert.Equal(t, "some", v.Data().String())
		})
		t.Run("direct", func(t *testing.T) {
			v, ok := s.Get("int")
			assert.True(t, ok)
			assert.Equal(t, env.TypeData, v.Type())
			assert.Equal(t, "some", v.Data().String())
		})
	})
}
