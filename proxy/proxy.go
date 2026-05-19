package proxy

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/navikt/hotbff/httpx"
	"github.com/navikt/hotbff/texas"
)

var (
	log = slog.Default().With("package", "proxy")
)

// Options for the proxy.
type Options struct {
	Target      string `json:"target"`      // The URL to proxy to (backend).
	StripPrefix bool   `json:"stripPrefix"` // Whether to strip the prefix from the request URL.
	IDPTarget   string `json:"idpTarget"`   // The target audience used in the token exchange (enables authentication if not empty).
}

// Map is a map of proxy Options keyed by URL prefix.
type Map map[string]*Options

// AddRoutes adds proxy handlers to the given mux based on the provided map.
func AddRoutes(mux *http.ServeMux, proxy Map, idp texas.IdentityProvider) error {
	if mux == nil {
		mux = http.DefaultServeMux
	}
	if proxy == nil {
		return nil
	}

	for prefix, opts := range proxy {
		if opts == nil {
			continue
		}
		h, err := reverseProxy(idp, opts)
		if err != nil {
			return err
		}
		if opts.StripPrefix {
			mux.Handle(prefix, http.StripPrefix(prefix, h))
		} else {
			mux.Handle(prefix, h)
		}
	}

	return nil
}

func reverseProxy(idp texas.IdentityProvider, opts *Options) (http.Handler, error) {
	t, err := url.Parse(opts.Target)
	if err != nil {
		return nil, fmt.Errorf("proxy: invalid target: %w", err)
	}

	if opts.IDPTarget == "" {
		return publicBackend(t), nil
	}
	return texas.Protected(idp, protectedBackend(t, idp, opts.IDPTarget)), nil
}

func publicBackend(target *url.URL) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(target)
		},
	}
}

func protectedBackend(target *url.URL, idp texas.TokenExchanger, idpTarget string) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(target)

			ctx := r.In.Context()
			user := texas.FromContext(ctx)
			if !user.Authenticated {
				log.WarnContext(ctx, "unauthenticated", "idp", idp, "idpTarget", idpTarget)
				return
			}

			ts, err := idp.ExchangeToken(ctx, idpTarget, user.Token)
			if err != nil {
				if !errors.Is(err, context.Canceled) {
					log.ErrorContext(ctx, "token exchange failed", "idp", idp, "idpTarget", idpTarget, "error", err)
				}
				return
			}

			r.Out.Header.Set(httpx.HeaderAuthorization, "Bearer "+ts.AccessToken)
		},
	}
}
