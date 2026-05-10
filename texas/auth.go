package texas

import (
	"context"
	"errors"
	"net/http"

	"github.com/navikt/hotbff/httpx"
)

// Authenticate is a middleware that checks for a valid token using the provided [TokenIntrospector].
// If the token is valid, an authenticated user is added to the request context and the next handler is called.
// If the idp is nil or disabled, the middleware simply calls the next handler without performing any validation.
func Authenticate(idp TokenIntrospector, next http.Handler) http.Handler {
	if idp == nil {
		log.Warn("idp nil", "idp", idp)
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		user := &User{Authenticated: false}

		if token, ok := TokenFromRequest(req); ok {
			ti, err := idp.IntrospectToken(ctx, token)
			if err != nil && !errors.Is(err, context.Canceled) {
				log.ErrorContext(ctx, "token introspection failed", "error", err)
			} else if ti.Active {
				user = &User{Authenticated: true, Token: token}
			}
		} else {
			log.DebugContext(ctx, "token missing")
		}

		ctx = NewContext(ctx, user)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}

// RequireAuthentication is a middleware that requires a valid token using the provided [TokenIntrospector].
// It redirects unauthenticated users to the login page.
func RequireAuthentication(idp TokenIntrospector, basePath string, next http.Handler) http.Handler {
	if idp == nil {
		return httpx.Unauthorized
	}
	return Authenticate(idp, Redirect(basePath, next))
}
