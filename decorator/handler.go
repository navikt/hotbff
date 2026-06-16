package decorator

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"slices"

	"github.com/navikt/hotbff/httpx"
)

var validCookieLanguages = []string{"nb", "nn"}

// Handler returns an HTTP handler that serves the provided template with decorator elements.
// The template will be executed with the fetched elements as data.
func Handler(name string, opts *Options) (http.Handler, error) {
	tmpl, err := template.ParseFiles(name)
	if err != nil {
		return nil, fmt.Errorf("decorator: failed parsing template %q: %w", name, err)
	}

	if opts == nil {
		opts = &Options{}
	}
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		reqOpts := opts.Clone()

		cookie, err := req.Cookie("decorator-language")
		if err == nil && slices.Contains(validCookieLanguages, cookie.Value) {
			reqOpts.Language = cookie.Value
		}

		elems, err := Fetch(ctx, reqOpts)
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

		w.Header().Set(httpx.HeaderContentType, httpx.ContentTypeTextHTML)
		w.Header().Set("Cache-Control", "no-cache, must-revalidate")
		if err := tmpl.Execute(w, elems); err != nil {
			serveError(ctx, w, "failed executing template", err)
		}
	}), nil
}

func serveError(ctx context.Context, w http.ResponseWriter, msg string, err error) {
	log.ErrorContext(ctx, msg, "error", err)
	http.Error(w, msg, http.StatusInternalServerError)
}
