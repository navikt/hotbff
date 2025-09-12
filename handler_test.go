package hotbff

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/navikt/hotbff/internal/assert"
)

func TestStaticHandler(t *testing.T) {
	rootDir := t.TempDir()
	indexPath := filepath.Join(rootDir, "index.html")
	err := os.WriteFile(indexPath, []byte("<!DOCTYPE html><html><body>test</body></html>"), 0644)
	assert.Nil(t, err)

	r := http.NewServeMux()
	r.Handle("/test/", http.StripPrefix("/test", staticHandler(rootDir, nil)))

	h := http.NewServeMux()
	h.Handle("/", r)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test/test", nil)

	h.ServeHTTP(w, req)

	res := w.Result()
	assert.Equal(t, res.StatusCode, http.StatusOK)
}
