package proxy

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/navikt/hotbff/internal/assert"
	"github.com/navikt/hotbff/texas"
)

func TestHandler(t *testing.T) {
	user := &texas.User{
		Authenticated: true,
		Token:         "userToken",
	}
	accessToken := "accessToken"
	idpTarget := "api://test.test.test/.default"

	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.Header.Get(texas.HeaderAuthorization), "Bearer "+accessToken)
		_, _ = w.Write([]byte("backend"))
	}))
	defer backend.Close()

	p := &Options{
		Target:      backend.URL,
		StripPrefix: false,
		IDP: exchangeTokenFunc(func(_ context.Context, target string, userToken string) (*texas.TokenSet, error) {
			assert.Equal(t, target, idpTarget)
			assert.Equal(t, userToken, user.Token)
			return &texas.TokenSet{
				AccessToken: accessToken,
			}, nil
		}),
		IDPTarget: idpTarget,
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(texas.NewContext(req.Context(), user))
	req.Header.Set(texas.HeaderAuthorization, "Bearer "+user.Token)

	h := p.Handler()
	h.ServeHTTP(w, req)

	res := w.Result()
	//goland:noinspection GoUnhandledErrorResult
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusOK)
}

type exchangeTokenFunc func(ctx context.Context, target string, userToken string) (*texas.TokenSet, error)

func (f exchangeTokenFunc) ExchangeToken(ctx context.Context, target string, userToken string) (*texas.TokenSet, error) {
	return f(ctx, target, userToken)
}
