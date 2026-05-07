package texas

import (
	"context"
	"errors"
	"net/http"
)

// Authenticate is a middleware that validates bearer tokens by introspecting them with the provided TokenIntrospector.
// If a valid bearer token is present and token introspection succeeds with an active token, the [User] is added to the context
// with Authenticated set to true. Otherwise, a non-authenticated User is added to the context.
func Authenticate(idp TokenIntrospector, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		user := &User{Authenticated: false}
		if token, ok := TokenFromRequest(req); ok {
			ti, err := idp.IntrospectToken(ctx, token)
			if err != nil {
				switch {
				case errors.Is(err, context.Canceled):
					log.DebugContext(ctx, "token introspection canceled", "error", err)
				case errors.Is(err, context.DeadlineExceeded):
					log.WarnContext(ctx, "token introspection timed out", "error", err)
				default:
					log.ErrorContext(ctx, "token introspection failed", "error", err)
				}
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
