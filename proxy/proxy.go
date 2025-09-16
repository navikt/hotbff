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

// Options are the options for the proxy.
type Options struct {
	Target      string                 `json:"target"`      // the URL to proxy to (backend)
	StripPrefix bool                   `json:"stripPrefix"` // whether to strip the prefix from the request URL
	IDP         texas.IdentityProvider `json:"idp"`         // IDP for token exchange, if empty, no token exchange is performed
	IDPTarget   string                 `json:"idpTarget"`   // the target audience used in the token exchange, required if IDP is set
}

// Handler returns a handler that proxies requests to the target URL.
func (t *Options) Handler() http.Handler {
	target, err := url.Parse(t.Target)
	if err != nil {
		slog.Error("proxy: invalid target", "target", t.Target, "error", err)
		os.Exit(1)
	}
	if t.IDP == "" {
		return &httputil.ReverseProxy{
			Rewrite: func(r *httputil.ProxyRequest) {
				r.SetURL(target)
			},
		}
	}
	return newTokenExchangeReverseProxy(target, t.IDP, t.IDPTarget)
}

func newTokenExchangeReverseProxy(target *url.URL, idp texas.IdentityProvider, idpTarget string) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(target)
			ctx := r.In.Context()
			user := texas.FromContext(ctx)
			if !user.Authenticated {
				slog.WarnContext(ctx, "proxy: user unauthenticated", "idp", idp, "idpTarget", idpTarget)
				return
			}
			ts, err := texas.ExchangeToken(ctx, idp, idpTarget, user.Token)
			if err != nil {
				if !errors.Is(err, context.Canceled) {
					slog.ErrorContext(ctx, "proxy: token exchange error", "idp", idp, "idpTarget", idpTarget, "error", err)
				}
				return
			}
			r.Out.Header.Set("Authorization", "Bearer "+ts.AccessToken)
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
