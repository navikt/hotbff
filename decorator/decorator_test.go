package decorator

import (
	"testing"

	"github.com/navikt/hotbff/internal/assert"
)

func TestGet(t *testing.T) {
	elems, err := Fetch(t.Context(), &Options{Context: "privatperson"})
	assert.Nil(t, err)
	assert.NotEqual(t, elems.HeadAssets, "")
	assert.NotEqual(t, elems.Header, "")
	assert.NotEqual(t, elems.Footer, "")
	assert.NotEqual(t, elems.Scripts, "")
}
