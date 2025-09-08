package decorator

import (
	"html/template"
	"net/http"
	"path"
)

func Handler(rootDir string, opts *Options) http.Handler {
	tmpl, err := template.ParseFiles(path.Join(rootDir, "index.html"))
	if err != nil {
		panic(err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		elem, err := Get(opts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = tmpl.Execute(w, &elem)
	})
}
