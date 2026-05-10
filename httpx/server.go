package httpx

import (
	"net/http"
	"path"
	"strings"
)

// MaybeStripPrefix uses [http.StripPrefix] to remove the given
// prefix from the request URL before passing it to the handler.
// If the prefix is "/", it returns the original handler as is.
func MaybeStripPrefix(prefix string, h http.Handler) http.Handler {
	if prefix == "/" {
		return h
	}
	return http.StripPrefix(prefix, h)
}

// NormalizePath cleans the given path and ensures it starts with a "/".
// If the cleaned path is empty, it returns "/".
func NormalizePath(p string) string {
	out := path.Clean("/" + strings.TrimSpace(p))
	if out == "." {
		return "/"
	}
	return out
}
