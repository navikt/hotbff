package texas

import (
	"log/slog"
	"net/http"
)

func Protected(idp IdentityProvider, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token, ok := TokenFromRequest(req)
		if !ok {
			slog.DebugContext(req.Context(), "unauthorized: token missing")
			LoginRedirect(w, req)
			return
		}
		i, err := IntrospectToken(idp, token)
		if err != nil {
			slog.ErrorContext(req.Context(), "unauthorized: error", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !i.Active {
			slog.DebugContext(req.Context(), "unauthorized: token invalid")
			LoginRedirect(w, req)
			return
		}
		ctx := NewContext(req.Context(), &User{Authenticated: true, Token: token})
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}

func LoginRedirect(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "/oauth2/login", http.StatusTemporaryRedirect)
}
