package texas

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/navikt/hotbff/internal/assert"
)

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "valid token",
			token: "valid_token",
		},
		{
			name:  "invalid token",
			token: "invalid_token",
		},
		{
			name:  "missing token",
			token: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var user *User

			active := tt.token == "valid_token"

			idp := NewTestIDP("", active, nil)
			handler := Authenticate(idp, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				user = FromContext(req.Context())
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			}))

			req := httptest.NewRequest("GET", "/", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			assert.NotNil(t, user)
			assert.Equal(t, user.Authenticated, active)
		})
	}
}
