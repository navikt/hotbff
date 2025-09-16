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

func TestHandlerIsAlive(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/isalive", nil)

	res := callHandler(t, req)
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.Equal(t, string(data), "ALIVE")
}

func TestHandlerIsReady(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/isready", nil)

	res := callHandler(t, req)
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.Equal(t, string(data), "READY")
}

func callHandler(t *testing.T, req *http.Request) *http.Response {
	t.Helper()

	rootDir := t.TempDir()
	indexPath := filepath.Join(rootDir, "index.html")
	err := os.WriteFile(indexPath, []byte("<!DOCTYPE html><html><body>test</body></html>"), 0644)
	assert.Nil(t, err)

	h := Handler(&Options{
		BasePath: "/",
		RootDir:  rootDir,
	})

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	return w.Result()
}
