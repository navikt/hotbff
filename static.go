package hotbff

import (
	"log/slog"
	"net/http"
	"path"

	"github.com/navikt/hotbff/decorator"
)

func staticHandler(rootDir string, opts *decorator.Options) http.Handler {
	index := indexHandler(rootDir, opts)
	fs := http.FileServer(http.Dir(rootDir))
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		slog.DebugContext(ctx, "hotbff: serving static file", "path", req.URL.Path)
		switch req.URL.Path {
		// index.html might need decoration, http.FileServer will not do that
		case "", "/", "index.html", "/index.html":
			index.ServeHTTP(w, req)
		default:
			r := &statusCodeRecorder{ResponseWriter: w}
			fs.ServeHTTP(r, req)
			// we have client side routing, override 404 to index.html
			if r.statusCode == http.StatusNotFound {
				index.ServeHTTP(w, req)
			}
		}
	})
}

func indexHandler(rootDir string, opts *decorator.Options) http.Handler {
	name := path.Join(rootDir, "index.html")
	if opts == nil {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			http.ServeFile(w, req, name)
		})
	}
	return decorator.Handler(name, opts)
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
