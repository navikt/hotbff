package texas

import (
	"log/slog"
	"net/http"
)

func Protected(idp IdentityProvider, next http.Handler) http.Handler {
	if idp == "" {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token, ok := TokenFromRequest(req)
		if !ok {
			slog.DebugContext(req.Context(), "texas: unauthorized: token missing")
			loginRedirect(w, req)
			return
		}
		ti, err := IntrospectToken(idp, token)
		if err != nil {
			slog.ErrorContext(req.Context(), "texas: error", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !ti.Active {
			slog.DebugContext(req.Context(), "texas: unauthorized: token invalid")
			loginRedirect(w, req)
			return
		}
		ctx := NewContext(req.Context(), &User{Authenticated: true, Token: token})
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}

func loginRedirect(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "/oauth2/login", http.StatusTemporaryRedirect)
}
