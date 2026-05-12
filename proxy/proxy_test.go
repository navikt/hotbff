package proxy

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/navikt/hotbff/httpx"
	"github.com/navikt/hotbff/internal/assert"
	"github.com/navikt/hotbff/texas"
)

func TestReverseProxy(t *testing.T) {
	target := "api://test.test.test/.default"
	user := &texas.User{
		Authenticated: true,
		Token:         "userToken",
	}
	accessToken := "accessToken"

	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.Header.Get(httpx.HeaderAuthorization), "Bearer "+accessToken)
		_, _ = w.Write([]byte("backend"))
	}))
	defer backend.Close()

	opts := &Options{
		Target:      backend.URL,
		StripPrefix: false,
		IDPTarget:   target,
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(texas.NewContext(req.Context(), user))
	req.Header.Set(httpx.HeaderAuthorization, "Bearer "+user.Token)

	idp := texas.NewTestIDP("accessToken", true, nil)
	h, err := reverseProxy(idp, opts)
	assert.Nil(t, err)
	h.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusOK)
}
