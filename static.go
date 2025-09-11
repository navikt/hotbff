package hotbff

import (
	"log/slog"
	"net/http"
	"path"

	"github.com/navikt/hotbff/decorator"
)

func staticHandler(rootDir string, opts *decorator.Options) http.Handler {
	idx := indexHandler(rootDir, opts)
	fs := http.FileServer(http.Dir(rootDir))
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		slog.Debug("hotbff: serving static file", "path", req.URL.Path)
		switch req.URL.Path {
		// index.html må kanskje dekoreres, fs vil ikke gjøre dette
		case "", "/", "index.html", "/index.html":
			idx.ServeHTTP(w, req)
		default:
			r := &statusCodeRecorder{ResponseWriter: w}
			fs.ServeHTTP(r, req)
			// vi har client side routing, overstyr 404 til index.html
			if r.statusCode == http.StatusNotFound {
				idx.ServeHTTP(w, req)
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
	return decorator.TemplateHandler(name, opts)
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
