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
	// Target is the target URL to proxy to.
	Target string `json:"target"`
	// StripPrefix indicates whether to strip the prefix from the request URL.
	StripPrefix bool `json:"stripPrefix"`
	// IDP is the identity provider to use for token exchange. If empty, no token exchange is performed.
	IDP texas.IdentityProvider `json:"idp"`
	// IDPTarget is the target audience used in the token exchange. Required if IDP is set.
	IDPTarget string `json:"idpTarget"`
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

// Map is a map of proxy Options keyed by URL prefix.
type Map map[string]*Options

// Configure adds proxy handlers to the given ServeMux based on the provided Map.
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
