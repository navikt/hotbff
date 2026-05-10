package hotbff

import (
	"errors"
	"io/fs"
	"net/http"
	"strings"
)

type spaHandler struct {
	root  http.FileSystem
	index http.Handler
	fs    http.Handler
}

func (h *spaHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	upath := req.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		req.URL.Path = upath
	}

	switch upath {
	// index.html might need decoration, http.FileServer does not perform decoration
	case "/", "/index.html":
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
			h.index.ServeHTTP(w, req)
			return
		}

		// serve file
		h.fs.ServeHTTP(w, req)
	}
}

func (h *spaHandler) handleError(w http.ResponseWriter, req *http.Request, err error) bool {
	if err == nil {
		return false
	}

	switch {
	case errors.Is(err, fs.ErrNotExist):
		h.index.ServeHTTP(w, req)
	case errors.Is(err, fs.ErrPermission):
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	default:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	return true
}

func newSPAHandler(rootDir string, index http.Handler) http.Handler {
	root := http.Dir(rootDir)
	if index == nil {
		index = http.NotFoundHandler()
	}
	return &spaHandler{root, index, http.FileServer(root)}
}
