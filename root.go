package hotbff

import (
	"errors"
	"io/fs"
	"net/http"
	"strings"

	"github.com/navikt/hotbff/httpx"
)

type rootHandler struct {
	root       http.FileSystem
	index      http.Handler
	fileServer http.Handler
}

func (h *rootHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	upath := req.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		req.URL.Path = upath
	}

	switch upath {
	case "/", "/index.html":
		// serve index.html
		h.index.ServeHTTP(w, req)
	default:
		f, err := h.root.Open(upath)
		if h.handleError(w, req, err) {
			return
		}
		defer f.Close()

		info, err := f.Stat()
		if h.handleError(w, req, err) {
			return
		}

		if info.IsDir() {
			// serve index.html
			h.index.ServeHTTP(w, req)
			return
		}

		// serve file
		h.fileServer.ServeHTTP(w, req)
	}
}

func (h *rootHandler) handleError(w http.ResponseWriter, req *http.Request, err error) bool {
	if err == nil {
		return false
	}

	switch {
	case errors.Is(err, fs.ErrNotExist):
		// serve index.html
		h.index.ServeHTTP(w, req)
	case errors.Is(err, fs.ErrPermission):
		httpx.ErrorCode(w, http.StatusForbidden)
	default:
		httpx.ErrorCode(w, http.StatusInternalServerError)
	}

	return true
}

func rootServer(rootDir string, index http.Handler) http.Handler {
	root := http.Dir(rootDir)
	if index == nil {
		index = http.NotFoundHandler()
	}
	return &rootHandler{root, index, http.FileServer(root)}
}
