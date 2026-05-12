package hotbff

import "net/http"

type CSPHeaderOptions struct {
}

func cspHeader(opts *CSPHeaderOptions, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Security-Policy", "")
		next.ServeHTTP(w, req)
	})
}
