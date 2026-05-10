package texas

import (
	"net/http"

	"github.com/navikt/hotbff/httpx"
)

func (idp *idp) Status() http.Handler {
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

func Status(idp IdentityProvider) http.Handler {
	if idp == nil {
		return httpx.Unauthorized
	}
	return Authenticate(idp, idp.Status())
}
