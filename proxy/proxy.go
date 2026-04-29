package proxy

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/navikt/hotbff/texas"
)

// Options for the proxy.
type Options struct {
	Target      string               `json:"target"`      // the URL to proxy to (backend)
	StripPrefix bool                 `json:"stripPrefix"` // whether to strip the prefix from the request URL
	IDP         texas.TokenExchanger `json:"idp"`         // identity provider for token exchange (if nil, no token exchange is performed)
	IDPTarget   string               `json:"idpTarget"`   // the target audience used in the token exchange (required if IDP is set)
}

// Handler returns a handler that proxies requests to the target URL.
func (opts *Options) Handler() http.Handler {
	target, err := url.Parse(opts.Target)
	if err != nil {
		slog.Error("proxy: invalid target", "target", opts.Target, "error", err)
		os.Exit(1)
	}
	if opts.IDP == nil {
		return &httputil.ReverseProxy{
			Rewrite: func(r *httputil.ProxyRequest) {
				r.SetURL(target)
			},
		}
	}
	return newTokenExchangeReverseProxy(target, opts.IDP, opts.IDPTarget)
}

func newTokenExchangeReverseProxy(target *url.URL, idp texas.TokenExchanger, idpTarget string) *httputil.ReverseProxy {
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

// Map is a map of proxy [Options] keyed by URL prefix.
type Map map[string]*Options

// Configure adds proxy handlers to the given [http.ServeMux] based on the provided [Map].
func Configure(proxy Map, mux *http.ServeMux) {
	if proxy == nil {
		slog.Info("proxy: no proxy")
		return
	}
	for prefix, opts := range proxy {
		slog.Info("proxy: adding proxy", "prefix", prefix, "target", opts.Target)
		proxyHandler := opts.Handler()
		if opts.StripPrefix {
			mux.Handle(prefix, http.StripPrefix(prefix, proxyHandler))
		} else {
			mux.Handle(prefix, proxyHandler)
		}
	}
}
