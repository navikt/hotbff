package httpx

import "net/http"

func ErrorCode(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

func Error(code int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		ErrorCode(w, code)
	})
}

var Unauthorized = Error(http.StatusUnauthorized)
