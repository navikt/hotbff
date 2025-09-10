package proxy

import (
	"fmt"
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

type Map map[string]*Options

func (t *Options) Handler(prefix string) http.Handler {
	target, err := url.Parse(t.Target)
	if err != nil {
		slog.Error("proxy: invalid target", "error", err)
		os.Exit(1)
	}
	var h http.Handler
	if t.IDP == "" {
		h = newReverseProxy(target)
	} else {
		h = newTokenExchangeReverseProxy(target, t.IDP, t.IDPTarget)
	}
	if t.StripPrefix {
		return http.StripPrefix(prefix, h)
	}
	return h
}

func newReverseProxy(target *url.URL) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(target)
		},
	}
}

func newTokenExchangeReverseProxy(target *url.URL, idp texas.IdentityProvider, idpTarget string) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(target)
			u := texas.FromContext(r.In.Context())
			if !u.Authenticated {
				slog.WarnContext(r.In.Context(), "proxy: unauthenticated")
				return
			}
			ts, err := texas.ExchangeToken(idp, idpTarget, u.Token)
			if err != nil {
				slog.ErrorContext(r.In.Context(), "proxy: token error", "error", err)
				return
			}
			r.Out.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ts.AccessToken))
		},
	}
}
