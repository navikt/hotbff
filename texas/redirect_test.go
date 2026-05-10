package texas

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/navikt/hotbff/internal/assert"
)

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
