package hotbff

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/navikt/hotbff/decorator"
	"github.com/navikt/hotbff/httpx"
)

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
		w.Header().Set(httpx.HeaderContentType, httpx.ContentTypeTextHTML)
		w.Header().Set("Cache-Control", "no-cache, must-revalidate")
		http.ServeContent(
			w,
			req,
			"index.html",
			modTime,
			bytes.NewReader(data),
		)
	}), nil
}
