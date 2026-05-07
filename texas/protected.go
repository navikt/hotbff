package texas

import (
	"net/http"
)

// Protected TODO
func Protected(basePath string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		user := FromContext(ctx)
		if user.Authenticated {
			next.ServeHTTP(w, req)
		} else {
			loginRedirect(w, req, basePath)
		}
	})
}
