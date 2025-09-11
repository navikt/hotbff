package texas

import (
	"log/slog"
	"net/http"
	"path"
)

func Protected(idp IdentityProvider, basePath string, next http.Handler) http.Handler {
	if idp == "" {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token, ok := TokenFromRequest(req)
		if !ok {
			slog.DebugContext(req.Context(), "texas: unauthorized: token missing")
			loginRedirect(w, req, basePath)
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
			loginRedirect(w, req, basePath)
			return
		}
		ctx := NewContext(req.Context(), &User{Authenticated: true, Token: token})
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}

func loginRedirect(w http.ResponseWriter, req *http.Request, basePath string) {
	http.Redirect(w, req, path.Join(basePath, "/oauth2/login"), http.StatusTemporaryRedirect)
}
