package texas

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/navikt/hotbff/internal/assert"
)

func TestRedirectActiveToken(t *testing.T) {
	res := callRedirectHandler(t, "userToken", true)
	//goland:noinspection GoUnhandledErrorResult
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.Equal(t, getLocation(t, res), "")
}

func TestRedirectInactiveToken(t *testing.T) {
	res := callRedirectHandler(t, "userToken", false)
	//goland:noinspection GoUnhandledErrorResult
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusTemporaryRedirect)
	assert.Equal(t, getLocation(t, res), "/oauth2/login?redirect="+url.QueryEscape("/"))
}

func TestRedirectMissingToken(t *testing.T) {
	res := callRedirectHandler(t, "", true)
	//goland:noinspection GoUnhandledErrorResult
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusTemporaryRedirect)
	assert.Equal(t, getLocation(t, res), "/oauth2/login?redirect="+url.QueryEscape("/"))
}

func callRedirectHandler(t *testing.T, userToken string, active bool) *http.Response {
	t.Helper()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if userToken != "" {
		req.Header.Set(HeaderAuthorization, "Bearer "+userToken)
	}

	idp := NewTestIDP(userToken, active, nil)
	h := Authenticate(idp, Redirect("/", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))
	h.ServeHTTP(w, req)

	return w.Result()
}

func getLocation(t *testing.T, res *http.Response) string {
	t.Helper()
	return res.Header.Get("Location")
}

func TestLoginRedirect(t *testing.T) {
	tests := []struct {
		name                     string
		basePath                 string
		requestURI               string
		expectedLocationContains []string
	}{
		{
			name:       "simple redirect",
			basePath:   "/app",
			requestURI: "/users",
			expectedLocationContains: []string{
				"/app/oauth2/login",
				"redirect=%2Fusers",
			},
		},
		{
			name:       "redirect with query string",
			basePath:   "/app",
			requestURI: "/accounts?key=value",
			expectedLocationContains: []string{
				"/app/oauth2/login",
				"redirect=%2Faccounts%3Fkey%3Dvalue",
			},
		},
		{
			name:       "redirect from root",
			basePath:   "/",
			requestURI: "/",
			expectedLocationContains: []string{
				"/oauth2/login",
				"redirect=%2F",
			},
		},
		{
			name:       "redirect with trailing slash",
			basePath:   "/app/",
			requestURI: "/customers",
			expectedLocationContains: []string{
				"/app/oauth2/login",
				"redirect=%2Fcustomers",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tt.requestURI, nil)

			loginRedirect(w, req, tt.basePath)

			assert.Equal(t, w.Code, http.StatusTemporaryRedirect)

			location := w.Header().Get("Location")
			assert.NotEmpty(t, location)

			for _, expected := range tt.expectedLocationContains {
				assert.Contains(t, location, expected)
			}
		})
	}
}
