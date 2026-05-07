package texas

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

var ErrUnauthorized = errors.New("unauthorized")

func authorizeRequest(req *http.Request, idp TokenIntrospector) (*User, error) {
	ctx := req.Context()
	token, ok := TokenFromRequest(req)
	if !ok {
		return nil, fmt.Errorf("token missing: %w", ErrUnauthorized)
	}
	ti, err := idp.IntrospectToken(ctx, token)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			return nil, fmt.Errorf("token introspection canceled: %w", err)
		case errors.Is(err, context.DeadlineExceeded):
			return nil, fmt.Errorf("token introspection timed out: %w", err)
		default:
			return nil, fmt.Errorf("token introspection failed: %w", err)
		}
	}
	if !ti.Active {
		return nil, fmt.Errorf("token inactive: %w", ErrUnauthorized)
	}
	return &User{Authenticated: true, Token: token}, nil
}

func Authenticate(idp TokenIntrospector, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		user := &User{Authenticated: false}
		if token, ok := TokenFromRequest(req); ok {
			ti, err := idp.IntrospectToken(ctx, token)
			if err != nil {
				switch {
				case errors.Is(err, context.Canceled):
					slog.DebugContext(ctx, "token introspection canceled")
				case errors.Is(err, context.DeadlineExceeded):
					slog.WarnContext(ctx, "token introspection timed out", "error", err)
				default:
					slog.ErrorContext(ctx, "token introspection failed", "error", err)
				}
			} else if ti.Active {
				user = &User{Authenticated: true, Token: token}
			}
		} else {
			slog.DebugContext(ctx, "token missing")
		}
		ctx = NewContext(ctx, user)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}
