package decorator

import (
	"html/template"
	"log/slog"
	"net/http"
	"os"
)

func TemplateHandler(name string, opts *Options) http.Handler {
	tmpl, err := template.ParseFiles(name)
	if err != nil {
		slog.Error("decorator: failed to parse template", "name", name, "error", err)
		os.Exit(1)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		elems, err := GetElements(opts)
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
