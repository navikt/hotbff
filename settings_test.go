package hotbff

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/navikt/hotbff/internal/assert"
)

func TestSettingsHandler(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/settings.js", nil)

	h := settingsHandler("/", []string{"API_URL"})
	h.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusOK)

	data, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	js := string(data)
	assert.HasPrefix(t, js, "window.appSettings = {\n")
	assert.Contains(t, js, `"API_URL": null`)
	assert.Contains(t, js, `"BASE_PATH": "/"`)
	assert.HasSuffix(t, js, "}\n")
}
