package hotbff

import (
	"errors"
	"io/fs"
	"net/http"
	"strings"

	"github.com/navikt/hotbff/httpx"
)

type spa struct {
	root  http.FileSystem
	index http.Handler
	fs    http.Handler
}

func (s *spa) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	upath := req.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		req.URL.Path = upath
	}

	switch upath {
	// index.html might need decoration, http.FileServer does not perform decoration
	case "/", "/index.html":
		s.index.ServeHTTP(w, req)
	default:
		f, err := s.root.Open(upath)
		if s.handleError(w, req, err) {
			return
		}
		defer f.Close()

		info, err := f.Stat()
		if s.handleError(w, req, err) {
			return
		}

		if info.IsDir() {
			s.index.ServeHTTP(w, req)
			return
		}

		// serve file
		s.fs.ServeHTTP(w, req)
	}
}

func (s *spa) handleError(w http.ResponseWriter, req *http.Request, err error) bool {
	if err == nil {
		return false
	}

	switch {
	case errors.Is(err, fs.ErrNotExist):
		s.index.ServeHTTP(w, req)
	case errors.Is(err, fs.ErrPermission):
		httpx.ErrorCode(w, http.StatusForbidden)
	default:
		httpx.ErrorCode(w, http.StatusInternalServerError)
	}

	return true
}

func newSPAHandler(rootDir string, index http.Handler) http.Handler {
	root := http.Dir(rootDir)
	if index == nil {
		index = http.NotFoundHandler()
	}
	return &spa{root, index, http.FileServer(root)}
}
