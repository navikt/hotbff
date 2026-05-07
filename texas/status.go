package texas

import (
	"net/http"
)

func (idp IdentityProvider) Status() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		user := FromContext(ctx)
		if user.Authenticated {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	})
}
