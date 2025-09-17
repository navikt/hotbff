package hotbff

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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
	assert.True(t, strings.HasPrefix(js, "window.appSettings = {\n"))
	assert.True(t, strings.Contains(js, `"API_URL": null`))
	assert.True(t, strings.Contains(js, `"BASE_PATH": "/"`))
	assert.True(t, strings.HasSuffix(js, "}\n"))
}
