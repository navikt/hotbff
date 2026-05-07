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

// Options for the proxy.
type Options struct {
	Target      string               `json:"target"`      // the URL to proxy to (backend)
	StripPrefix bool                 `json:"stripPrefix"` // whether to strip the prefix from the request URL
	IDP         texas.TokenExchanger `json:"idp"`         // identity provider for token exchange (if unset, no token exchange is performed)
	IDPTarget   string               `json:"idpTarget"`   // the target audience used in the token exchange (required if IDP is set)
}

// Map is a map of proxy [Options] keyed by URL prefix.
type Map map[string]*Options

// Configure adds proxy handlers to the given [http.ServeMux] based on the provided [Map].
func Configure(mux *http.ServeMux, proxy Map) error {
	if mux == nil {
		mux = http.DefaultServeMux
	}
	if proxy == nil {
		slog.Info("proxy: no proxy")
		return nil
	}
	for prefix, opts := range proxy {
		if opts == nil {
			slog.Warn("proxy: skipping proxy", "prefix", prefix)
			continue
		}
		slog.Info("proxy: adding proxy", "prefix", prefix, "target", opts.Target)
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

func newReverseProxy(opts *Options) (h http.Handler, err error) {
	var t *url.URL
	t, err = url.Parse(opts.Target)
	if err != nil {
		return nil, fmt.Errorf("proxy: invalid target: %w", err)
	}
	if opts.IDP == nil || !opts.IDP.Enabled() {
		h, err = publicBackend(t), nil
	} else {
		h, err = protectedBackend(t, opts.IDP, opts.IDPTarget), nil
	}
	return
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
				slog.WarnContext(ctx, "proxy: user unauthenticated", "idp", idp, "idpTarget", idpTarget)
				return
			}
			ts, err := idp.ExchangeToken(ctx, idpTarget, user.Token)
			if err != nil {
				if !errors.Is(err, context.Canceled) {
					slog.ErrorContext(ctx, "proxy: token exchange error", "idp", idp, "idpTarget", idpTarget, "error", err)
				}
				return
			}
			r.Out.Header.Set(texas.HeaderAuthorization, "Bearer "+ts.AccessToken)
		},
	}
}
