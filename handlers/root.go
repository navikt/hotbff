package handlers

import (
	"net/http"
	"path"

	"github.com/navikt/hotbff/decorator"
)

func Root(rootDir string, opts *decorator.Options) http.Handler {
	var idx http.Handler
	if opts == nil {
		idx = ServeIndex(rootDir)
	} else {
		idx = decorator.ServeIndex(rootDir, opts)
	}
	fs := http.FileServer(http.Dir(rootDir))
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "", "/", "/index.html":
			idx.ServeHTTP(w, req)
			return
		}
		r := &statusCodeRecorder{ResponseWriter: w}
		fs.ServeHTTP(r, req)
		if r.statusCode == http.StatusNotFound {
			idx.ServeHTTP(w, req)
			return
		}
	})
}

func ServeIndex(rootDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeFile(w, req, path.Join(rootDir, "index.html"))
	})
}

type statusCodeRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusCodeRecorder) Write(data []byte) (int, error) {
	if r.statusCode != http.StatusNotFound {
		return r.ResponseWriter.Write(data)
	}
	return len(data), nil
}

func (r *statusCodeRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	if statusCode != http.StatusNotFound {
		r.ResponseWriter.WriteHeader(statusCode)
	}
}
