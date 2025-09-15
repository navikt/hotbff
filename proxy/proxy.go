package proxy

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/navikt/hotbff/texas"
)

type Options struct {
	Target      string                 `json:"target"`
	StripPrefix bool                   `json:"stripPrefix"`
	IDP         texas.IdentityProvider `json:"idp"`
	IDPTarget   string                 `json:"idpTarget"`
}

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

type Map map[string]*Options

func Configure(m *Map, mux *http.ServeMux) {
	if m == nil {
		slog.Info("proxy: no proxy")
		return
	}
	for prefix, opts := range *m {
		slog.Info("proxy: adding proxy", "prefix", prefix, "target", opts.Target)
		proxyHandler := opts.Handler()
		if opts.StripPrefix {
			mux.Handle(prefix, http.StripPrefix(prefix, proxyHandler))
		} else {
			mux.Handle(prefix, proxyHandler)
		}
	}
}

func newTokenExchangeReverseProxy(target *url.URL, idp texas.IdentityProvider, idpTarget string) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(target)
			user := texas.FromContext(r.In.Context())
			if !user.Authenticated {
				slog.WarnContext(r.In.Context(), "proxy: user unauthenticated", "idp", idp, "idpTarget", idpTarget)
				return
			}
			ts, err := texas.ExchangeToken(idp, idpTarget, user.Token)
			if err != nil {
				slog.ErrorContext(r.In.Context(), "proxy: token exchange error", "idp", idp, "idpTarget", idpTarget, "error", err)
				return
			}
			r.Out.Header.Set("Authorization", "Bearer "+ts.AccessToken)
		},
	}
}
