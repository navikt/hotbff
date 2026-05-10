package hotbff

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/navikt/hotbff/internal/assert"
	"github.com/navikt/hotbff/proxy"
	"github.com/navikt/hotbff/texas"
)

func TestHandlerIsAlive(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/isalive", nil)

	res := callHandler(t, req)
	//goland:noinspection GoUnhandledErrorResult
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.Equal(t, string(data), "ALIVE")
}

func TestHandlerIsReady(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/isready", nil)

	res := callHandler(t, req)
	//goland:noinspection GoUnhandledErrorResult
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

	h := http.NewServeMux()
	Configure(h, &Options{
		BasePath: "/test/",
		RootDir:  rootDir,
	})

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	return w.Result()
}

func TestStatusUnauthorizedWithoutIDP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test/auth/status", nil)

	res := callHandler(t, req)
	defer res.Body.Close()

	assert.Equal(t, res.StatusCode, http.StatusUnauthorized)
}

func TestStatusOKWithAuthorizedUser(t *testing.T) {
	h := newConfiguredHandler(t, &Options{
		BasePath: "/test/",
		RootDir:  withIndexFile(t),
		IDP:      texas.NewTestIDP("", true, nil),
	})

	res := callRequest(h, http.MethodGet, "/test/auth/status", "valid_token")
	defer res.Body.Close()

	assert.Equal(t, res.StatusCode, http.StatusOK)
}

func TestStatusUnauthorizedWithMissingToken(t *testing.T) {
	h := newConfiguredHandler(t, &Options{
		BasePath: "/test/",
		RootDir:  withIndexFile(t),
		IDP:      texas.NewTestIDP("", true, nil),
	})

	res := callRequest(h, http.MethodGet, "/test/auth/status", "")
	defer res.Body.Close()

	assert.Equal(t, res.StatusCode, http.StatusUnauthorized)
}

func TestIndexRedirectsWhenProtectedAndMissingToken(t *testing.T) {
	h := newConfiguredHandler(t, &Options{
		BasePath: "/test/",
		RootDir:  withIndexFile(t),
		IDP:      texas.NewTestIDP("", true, nil),
	})

	res := callRequest(h, http.MethodGet, "/test/", "")
	defer res.Body.Close()

	assert.Equal(t, res.StatusCode, http.StatusTemporaryRedirect)
	assert.Equal(t, res.Header.Get("Location"), "/test/oauth2/login?redirect=%2F")
}

func TestIndexSkipsAuthOnPublicPath(t *testing.T) {
	rootDir := withIndexFile(t)
	h := newConfiguredHandler(t, &Options{
		BasePath:    "/test/",
		RootDir:     rootDir,
		IDP:         texas.NewTestIDP("", true, nil),
		PublicPaths: []string{"/open"},
	})

	res := callRequest(h, http.MethodGet, "/test/open", "")
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.Contains(t, string(body), "test")
}

func TestProxyUnauthorizedWithoutRedirectWhenMissingToken(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Header.Get(texas.HeaderAuthorization) == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	h := newConfiguredHandler(t, &Options{
		BasePath: "/test/",
		RootDir:  withIndexFile(t),
		IDP:      texas.NewTestIDP("backend_token", true, nil),
		Proxy: proxy.Map{
			"/api/": {
				Target: backend.URL,
				// IDP:    texas.NewTestIDP("backend_token", true, nil),
			},
		},
	})

	res := callRequest(h, http.MethodGet, "/test/api/resource", "")
	defer res.Body.Close()

	assert.Equal(t, res.StatusCode, http.StatusUnauthorized)
	assert.Equal(t, res.Header.Get("Location"), "")
}

func TestSPAStaticFileRemainsPublic(t *testing.T) {
	rootDir := t.TempDir()
	indexPath := filepath.Join(rootDir, "index.html")
	err := os.WriteFile(indexPath, []byte("<!DOCTYPE html><html><body>test</body></html>"), 0644)
	assert.Nil(t, err)

	jsPath := filepath.Join(rootDir, "assets", "app.js")
	err = os.MkdirAll(filepath.Dir(jsPath), 0755)
	assert.Nil(t, err)
	err = os.WriteFile(jsPath, []byte("console.log('ok');"), 0644)
	assert.Nil(t, err)

	h := newConfiguredHandler(t, &Options{
		BasePath: "/test/",
		RootDir:  rootDir,
		IDP:      texas.NewTestIDP("", true, nil),
	})

	res := callRequest(h, http.MethodGet, "/test/assets/app.js", "")
	defer res.Body.Close()

	body, readErr := io.ReadAll(res.Body)
	assert.Nil(t, readErr)
	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.Contains(t, string(body), "console.log('ok');")
}

func newConfiguredHandler(t *testing.T, opts *Options) *http.ServeMux {
	t.Helper()
	h := http.NewServeMux()
	Configure(h, opts)
	return h
}

func callRequest(h *http.ServeMux, method, urlPath, token string) *http.Response {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, urlPath, nil)
	if token != "" {
		req.Header.Set(texas.HeaderAuthorization, "Bearer "+token)
	}
	h.ServeHTTP(w, req)
	return w.Result()
}

func withIndexFile(t *testing.T) string {
	t.Helper()
	rootDir := t.TempDir()
	indexPath := filepath.Join(rootDir, "index.html")
	err := os.WriteFile(indexPath, []byte("<!DOCTYPE html><html><body>test</body></html>"), 0644)
	assert.Nil(t, err)
	return rootDir
}
