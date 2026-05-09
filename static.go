package hotbff

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/navikt/hotbff/decorator"
)

func staticHandler(rootDir string, index http.Handler) http.Handler {
	if index == nil {
		index = http.NotFoundHandler()
	}
	fs := http.FileServer(http.Dir(rootDir))
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		slog.DebugContext(ctx, "hotbff: serving static file", "path", req.URL.Path)
		switch req.URL.Path {
		// index.html might need decoration, http.FileServer does not perform decoration
		case "", "/", "index.html", "/index.html":
			index.ServeHTTP(w, req)
		default:
			r := &statusCodeRecorder{ResponseWriter: w}
			fs.ServeHTTP(r, req)
			// we have client side routing, serve index.html instead of 404
			if r.statusCode == http.StatusNotFound {
				index.ServeHTTP(w, req)
			}
		}
	})
}

func indexHandler(rootDir string, opts *decorator.Options) (http.Handler, error) {
	name := filepath.Join(rootDir, "index.html")

	// index.html with decoration
	if opts != nil {
		h, err := decorator.Handler(name, opts)
		if err != nil {
			return nil, fmt.Errorf("decorator handler %q: %w", name, err)
		}
		return h, nil
	}

	// index.html without decoration
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("read %q: %w", name, err)
	}

	info, err := os.Stat(name)
	if err != nil {
		return nil, fmt.Errorf("stat %q: %w", name, err)
	}
	modTime := info.ModTime()

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeContent(
			w,
			req,
			"index.html",
			modTime,
			bytes.NewReader(data),
		)
	}), nil
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
