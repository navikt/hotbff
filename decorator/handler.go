package decorator

import (
	"html/template"
	"net/http"
	"path"

	"github.com/navikt/hotbff/common"
)

func ServeIndex(rootDir string, opts *Options) http.Handler {
	tmpl, err := template.ParseFiles(path.Join(rootDir, "index.html"))
	if err != nil {
		common.Fatal("failed to parse template", "rootDir", rootDir, "error", err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		elems, err := Get(opts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, &elems); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
