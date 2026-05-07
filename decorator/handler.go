package decorator

import (
	"context"
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"os"
)

// Handler returns a [http.Handler] that renders the named template file
// decorated with [Elements] fetched using the given [Options].
// If fetching the elements fails, it returns a 500 Internal Server Error.
func Handler(name string, opts *Options) http.Handler {
	tmpl, err := template.ParseFiles(name)
	if err != nil {
		slog.Error("failed parsing decorator template", "name", name, "error", err)
		os.Exit(1)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		cookie, err := req.Cookie("decorator-language")
		if err == nil && (cookie.Value == "nb" || cookie.Value == "nn") {
			opts.Language = cookie.Value
		}

		elems, err := Fetch(ctx, opts)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				w.WriteHeader(http.StatusGatewayTimeout)
			} else {
				slog.ErrorContext(ctx, "failed fetching decorator elements", "error", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, elems); err != nil {
			slog.ErrorContext(ctx, "failed executing decorator template", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
