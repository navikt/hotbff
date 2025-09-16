package texas

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/navikt/hotbff/internal/assert"
)

func TestProtectedActiveToken(t *testing.T) {
	res := callProtectedHandler(t, "userToken", true)
	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.Equal(t, getLocation(t, res), "")
}

func TestProtectedInactiveToken(t *testing.T) {
	res := callProtectedHandler(t, "userToken", false)
	assert.Equal(t, res.StatusCode, http.StatusTemporaryRedirect)
	assert.Equal(t, getLocation(t, res), "/oauth2/login")
}

func TestProtectedMissingToken(t *testing.T) {
	res := callProtectedHandler(t, "", true)
	assert.Equal(t, res.StatusCode, http.StatusTemporaryRedirect)
	assert.Equal(t, getLocation(t, res), "/oauth2/login")
}

func callProtectedHandler(t *testing.T, userToken string, active bool) *http.Response {
	t.Helper()
	server := texasIntrospectionServer(t, active)
	defer server.Close()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if userToken != "" {
		req.Header.Set(HeaderAuthorization, "Bearer "+userToken)
	}

	h := Protected(TokenX, "/", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	h.ServeHTTP(w, req)

	return w.Result()
}

func texasIntrospectionServer(t *testing.T, active bool) *httptest.Server {
	t.Helper()
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		assert.Nil(t, err)
		if active {
			_, _ = w.Write([]byte(`{"active":true}`))
		} else {
			_, _ = w.Write([]byte(`{"active":false}`))
		}
	}))
	tokenIntrospectionURL = s.URL
	return s
}

func getLocation(t *testing.T, res *http.Response) string {
	t.Helper()
	return res.Header.Get("Location")
}
