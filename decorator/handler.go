package decorator

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"slices"
)

var validCookieLanguages = []string{"nb", "nn"}

// Handler returns an [http.Handler] that renders the named template file
// decorated with [Elements] fetched using the given [Options].
// If fetching the elements fails, it returns a 500 Internal Server Error.
func Handler(name string, opts *Options) (http.Handler, error) {
	tmpl, err := template.ParseFiles(name)
	if err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}
	if opts == nil {
		opts = &Options{}
	}
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		cookie, err := req.Cookie("decorator-language")
		if err == nil && slices.Contains(validCookieLanguages, cookie.Value) {
			opts.Language = cookie.Value
		}

		elems, err := Fetch(ctx, opts)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				w.WriteHeader(http.StatusGatewayTimeout)
			} else {
				httpError(ctx, w, "fetching elements", err)
			}
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, elems); err != nil {
			httpError(ctx, w, "executing template", err)
		}
	}), nil
}

func httpError(ctx context.Context, w http.ResponseWriter, msg string, err error) {
	log.ErrorContext(ctx, msg, "error", err)
	http.Error(w, msg, http.StatusInternalServerError)
}
