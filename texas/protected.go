package texas

import (
	"log/slog"
	"net/http"
	"path"
)

func Protected(idp IdentityProvider, basePath string, next http.Handler) http.Handler {
	if idp == "" {
		slog.Warn("texas: identity provider not set, token validation disabled")
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token, ok := TokenFromRequest(req)
		if !ok {
			slog.DebugContext(req.Context(), "texas: unauthorized: token missing", "idp", idp)
			loginRedirect(w, req, basePath)
			return
		}
		ti, err := IntrospectToken(idp, token)
		if err != nil {
			slog.ErrorContext(req.Context(), "texas: error", "idp", idp, "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !ti.Active {
			slog.DebugContext(req.Context(), "texas: unauthorized: token invalid", "idp", idp)
			loginRedirect(w, req, basePath)
			return
		}
		ctx := NewContext(req.Context(), &User{Authenticated: true, Token: token})
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}

func loginRedirect(w http.ResponseWriter, req *http.Request, basePath string) {
	url := path.Join(basePath, "/oauth2/login")
	if basePath != "/" {
		url = url + "?redirect=" + basePath
	}
	slog.DebugContext(req.Context(), "texas: login redirect", "url", url)
	http.Redirect(w, req, url, http.StatusTemporaryRedirect)
}
