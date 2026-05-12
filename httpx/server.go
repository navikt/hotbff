package httpx

import (
	"net/http"
	"path"
	"strings"
)

// NormalizePath cleans the given path and ensures it starts with a "/".
// If the cleaned path is empty, it returns "/".
func NormalizePath(p string) string {
	return path.Clean("/" + p)
}

// PathMatches checks if the request URL path matches any of the given paths.
// It normalizes both the request path and the provided paths before comparison.
// A match occurs if the requested path is equal to a provided path or starts with a provided path followed by a "/".
func PathMatches(req *http.Request, paths []string) bool {
	reqPath := NormalizePath(req.URL.Path)
	for _, pubPath := range paths {
		pubPath = NormalizePath(pubPath)
		if reqPath == pubPath || strings.HasPrefix(reqPath, pubPath+"/") {
			return true
		}
	}
	return false
}
