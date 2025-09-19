package decorator

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/navikt/hotbff/internal/assert"
)

func TestHandler(t *testing.T) {
	rootDir := t.TempDir()
	indexPath := filepath.Join(rootDir, "index.html")
	err := os.WriteFile(indexPath, []byte(testTemplate), 0644)
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	h := Handler(indexPath, &Options{Context: "privatperson"})
	h.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusOK)

	data, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	html := string(data)
	assert.Contains(t, html, "decorator-header")
	assert.Contains(t, html, "decorator-footer")
}

const testTemplate = `
<!DOCTYPE html>
<html>
	<head>
		{{.HeadAssets}}
	</head>
	<body>
		{{.Header}}
		<main>test</main>
		{{.Footer}}
		{{.Scripts}}
	</body>
</html>
`
