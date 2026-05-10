package hotbff

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/navikt/hotbff/internal/assert"
)

func TestSPAHandler(t *testing.T) {
	rootDir := t.TempDir()

	indexHTML := "<!DOCTYPE html><html><body>index-content</body></html>"
	indexPath := filepath.Join(rootDir, "index.html")
	err := os.WriteFile(indexPath, []byte(indexHTML), 0644)
	assert.Nil(t, err)

	jsBody := "console.log('asset-content');"
	jsPath := filepath.Join(rootDir, "assets", "index.js")
	err = os.MkdirAll(filepath.Dir(jsPath), 0755)
	assert.Nil(t, err)
	err = os.WriteFile(jsPath, []byte(jsBody), 0644)
	assert.Nil(t, err)

	index, err := indexHandler(rootDir, nil)
	assert.Nil(t, err)

	r := http.NewServeMux()
	r.Handle("/test/", http.StripPrefix("/test", newSPAHandler(rootDir, index)))

	h := http.NewServeMux()
	h.Handle("/", r)

	tests := []struct {
		name         string
		path         string
		wantStatus   int
		wantContains string
		dontWant     string
		wantCTPrefix string
	}{
		{
			name:         "direct hit index.html serves index",
			path:         "/test/index.html",
			wantStatus:   http.StatusOK,
			wantContains: "index-content",
			dontWant:     "asset-content",
			wantCTPrefix: "text/html",
		},
		{
			name:         "not found route falls back to index",
			path:         "/test/some/route",
			wantStatus:   http.StatusOK,
			wantContains: "index-content",
			dontWant:     "asset-content",
			wantCTPrefix: "text/html",
		},
		{
			name:         "not found asset falls back to index",
			path:         "/test/some/asset.js",
			wantStatus:   http.StatusOK,
			wantContains: "index-content",
			dontWant:     "asset-content",
			wantCTPrefix: "text/html",
		},
		{
			name:         "existing js asset is served as file",
			path:         "/test/assets/index.js",
			wantStatus:   http.StatusOK,
			wantContains: "asset-content",
			dontWant:     "index-content",
			wantCTPrefix: "text/javascript",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)

			h.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, tc.wantStatus)

			body, err := io.ReadAll(res.Body)
			assert.Nil(t, err)
			bodyStr := string(body)

			assert.Contains(t, bodyStr, tc.wantContains)
			assert.False(t, bodyStr == tc.dontWant)
			assert.HasPrefix(t, res.Header.Get("Content-Type"), tc.wantCTPrefix)
		})
	}
}
