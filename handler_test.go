package hotbff

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestRootHandler(t *testing.T) {
	rootDir := t.TempDir()
	indexPath := filepath.Join(rootDir, "index.html")
	err := os.WriteFile(indexPath, []byte("<!DOCTYPE html><html><body>test</body></html>"), 0644)
	if err != nil {
		t.Fatalf("failed to create index.html: %v", err)
	}
	r := http.NewServeMux()
	r.Handle("/test/", http.StripPrefix("/test/", staticHandler(rootDir, nil)))
	h := http.NewServeMux()
	h.Handle("/", r)
	req := httptest.NewRequest("GET", "/test/test", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}
}
