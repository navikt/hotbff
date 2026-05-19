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

	h := settingsHandler("/test/", []string{"API_URL"})
	h.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusOK)

	data, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	js := string(data)
	assert.Equal(t, js, expectedJs)
}

const expectedJs = `window.appSettings = {
  "API_URL": null,
  "BASE_PATH": "/test/",
  "GIT_COMMIT": null,
  "NAIS_APP_NAME": null,
  "NAIS_CLUSTER_NAME": null,
  "USE_MSW": null
}
`
