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

// Fetch retrieves decorator [Elements] using the given [Options].
func Fetch(ctx context.Context, opts *Options) (*Elements, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getDecoratorURL(), nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = opts.Query().Encode()
	slog.Debug("decorator: fetching elements", "url", req.URL)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("decorator: unexpected statusCode: %d", res.StatusCode)
	}
	var elems Elements
	if err := json.NewDecoder(res.Body).Decode(&elems); err != nil {
		return nil, err
	}
	return &elems, nil
}

// Options for the decorator.
type Options struct {
	Context string // "privatperson" | "arbeidsgiver" | "samarbeidspartner"
}

// Query is the decorator [Options] expressed as URL query parameters.
func (o *Options) Query() url.Values {
	q := url.Values{}
	q.Set("context", o.Context)
	return q
}

// Elements fetched from the decorator.
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
)

func getDecoratorURL() string {
	switch cluster {
	case "", "local", "test":
		return decoratorURLDev
	default:
		return decoratorURL
	}
}
