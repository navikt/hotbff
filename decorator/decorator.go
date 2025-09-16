package decorator

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"net/url"
	"os"
)

// Fetch fetches the decorator elements using the given options.
// It returns an Elements struct containing HTML snippets.
func Fetch(ctx context.Context, opts *Options) (*Elements, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getDecoratorURL(), nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = opts.Query().Encode()
	slog.Debug("decorator: fetching elements", "url", req.URL.String())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("decorator: unexpected statusCode: %d", res.StatusCode)
	}
	var elems *Elements
	err = json.NewDecoder(res.Body).Decode(&elems)
	if err != nil {
		return nil, err
	}
	return elems, nil
}

type Options struct {
	Context string
}

func (o *Options) Query() url.Values {
	q := url.Values{}
	q.Set("context", o.Context)
	return q
}

type Elements struct {
	HeadAssets template.HTML `json:"headAssets"`
	Header     template.HTML `json:"header"`
	Footer     template.HTML `json:"footer"`
	Scripts    template.HTML `json:"scripts"`
}

var (
	cluster         = os.Getenv("NAIS_CLUSTER_NAME")
	decoratorURL    = "http://nav-dekoratoren.personbruker/dekoratoren/ssr"
	decoratorURLDev = "https://dekoratoren.ekstern.dev.nav.no/dekoratoren/ssr"
	//decoratorURLProd = "https://www.nav.no/dekoratoren/ssr"
)

func getDecoratorURL() string {
	switch cluster {
	case "", "local", "test":
		return decoratorURLDev
	default:
		return decoratorURL
	}
}
