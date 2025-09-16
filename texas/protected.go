package texas

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"path"
)

// Protected is a middleware that protects the given handler with token validation.
// If the identity provider is not set, the handler is returned as is.
// If the token is missing or invalid, the user is redirected to the login page.
func Protected(idp IdentityProvider, basePath string, next http.Handler) http.Handler {
	if idp == "" {
		slog.Warn("texas: identity provider not set, token validation disabled")
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		token, ok := TokenFromRequest(req)
		if !ok {
			slog.DebugContext(ctx, "texas: unauthorized: token missing", "idp", idp)
			loginRedirect(w, req, basePath)
			return
		}
		ti, err := IntrospectToken(ctx, idp, token)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				w.WriteHeader(http.StatusRequestTimeout)
			} else {
				slog.ErrorContext(ctx, "texas: error", "idp", idp, "error", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		if !ti.Active {
			slog.DebugContext(ctx, "texas: unauthorized: token invalid", "idp", idp)
			loginRedirect(w, req, basePath)
			return
		}
		ctx = NewContext(ctx, &User{Authenticated: true, Token: token})
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}

func loginRedirect(w http.ResponseWriter, req *http.Request, basePath string) {
	ctx := req.Context()
	url := path.Join(basePath, "/oauth2/login")
	if basePath != "/" {
		url = url + "?redirect=" + basePath
	}
	slog.DebugContext(ctx, "texas: login redirect", "url", url)
	http.Redirect(w, req, url, http.StatusTemporaryRedirect)
}
