package decorator

import (
	"context"
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"os"
)

// Handler returns a handler that renders the named template file
// decorated with [Elements] fetched using the given [Options].
// If fetching the elements fails, it returns a 500 Internal Server Error.
func Handler(name string, opts *Options) http.Handler {
	tmpl, err := template.ParseFiles(name)
	if err != nil {
		slog.Error("decorator: failed parsing template", "name", name, "error", err)
		os.Exit(1)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		elems, err := Fetch(ctx, opts)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				w.WriteHeader(http.StatusRequestTimeout)
			} else {
				slog.ErrorContext(ctx, "decorator: failed fetching elements", "error", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		w.Header().Set("Cache-Control", "max-age=3600, private")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, &elems); err != nil {
			slog.ErrorContext(ctx, "decorator: failed executing template", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
