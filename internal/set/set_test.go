package set

import (
	"testing"

	"github.com/navikt/hotbff/internal/assert"
)

func TestAdd(t *testing.T) {
	s := New[string]()
	assert.Equal(t, s.Len(), 0)

	v := "first"
	assert.False(t, s.Has(v))

	s.Add(v)
	assert.Equal(t, s.Len(), 1)
	assert.True(t, s.Has(v))

	s.Add(v)
	assert.Equal(t, s.Len(), 1)
	assert.True(t, s.Has(v))
}

func TestAddAll(t *testing.T) {
	s := New[string]()
	assert.Equal(t, s.Len(), 0)

	v1 := "first"
	assert.False(t, s.Has(v1))
	v2 := "second"
	assert.False(t, s.Has(v2))

	s.AddAll(v1, v2)
	assert.Equal(t, s.Len(), 2)
	assert.True(t, s.Has(v1))
	assert.True(t, s.Has(v2))

	s.AddAll()
	assert.Equal(t, s.Len(), 2)
}
