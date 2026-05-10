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
		return nil, fmt.Errorf("failed parsing template %q: %w", name, err)
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
			switch {
			case errors.Is(err, context.Canceled):
				return
			case errors.Is(err, context.DeadlineExceeded):
				w.WriteHeader(http.StatusGatewayTimeout)
			default:
				serveError(ctx, w, "failed fetching elements", err)
			}
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, elems); err != nil {
			serveError(ctx, w, "failed executing template", err)
		}
	}), nil
}

func serveError(ctx context.Context, w http.ResponseWriter, msg string, err error) {
	log.ErrorContext(ctx, msg, "error", err)
	http.Error(w, msg, http.StatusInternalServerError)
}
