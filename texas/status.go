package texas

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
)

func (idp IdentityProvider) Status() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		token, ok := TokenFromRequest(req)
		if !ok {
			slog.DebugContext(ctx, "texas: unauthorized: token missing", "idp", idp)
			w.WriteHeader(http.StatusUnauthorized)
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
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}
