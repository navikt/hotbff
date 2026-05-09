package proxy

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/navikt/hotbff/texas"
)

var (
	log = slog.Default().With("package", "proxy")
)

// Options for the proxy.
type Options struct {
	Target      string               `json:"target"`      // The URL to proxy to (backend).
	StripPrefix bool                 `json:"stripPrefix"` // Whether to strip the prefix from the request URL or not.
	IDP         texas.TokenExchanger `json:"idp"`         // Identity provider for token exchange (if unset, no token exchange is performed).
	IDPTarget   string               `json:"idpTarget"`   // The target audience used in the token exchange (required if IDP is set).
}

// Map is a map of proxy [Options] keyed by URL prefix.
type Map map[string]*Options

// Configure adds proxy handlers to the given [http.ServeMux] based on the provided [Map].
func Configure(mux *http.ServeMux, proxy Map) error {
	if mux == nil {
		mux = http.DefaultServeMux
	}

	if proxy == nil {
		log.Info("no proxy endpoints")
		return nil
	}

	for prefix, opts := range proxy {
		if opts == nil {
			log.Warn("skipping proxy", "prefix", prefix)
			continue
		}

		log.Info("adding proxy", "prefix", prefix, "target", opts.Target)
		h, err := newReverseProxy(opts)
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

func newReverseProxy(opts *Options) (http.Handler, error) {
	t, err := url.Parse(opts.Target)
	if err != nil {
		return nil, fmt.Errorf("invalid target: %w", err)
	}

	if opts.IDP == nil || !opts.IDP.Enabled() {
		return publicBackend(t), nil
	}

	return protectedBackend(t, opts.IDP, opts.IDPTarget), nil
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
				log.WarnContext(ctx, "user unauthenticated", "idp", idp, "idpTarget", idpTarget)
				return
			}

			ts, err := idp.ExchangeToken(ctx, idpTarget, user.Token)
			if err != nil {
				if !errors.Is(err, context.Canceled) {
					log.ErrorContext(ctx, "token exchange failed", "idp", idp, "idpTarget", idpTarget, "error", err)
				}
				return
			}

			r.Out.Header.Set(texas.HeaderAuthorization, "Bearer "+ts.AccessToken)
		},
	}
}
