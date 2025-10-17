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
	"strconv"
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

// AvailableLanguage represents a language option in the decorator.
type AvailableLanguage struct {
	Locale      string `json:"locale"`
	HandleInApp bool   `json:"handleInApp"`
}

// Options for the decorator.
type Options struct {
	Context            string // "privatperson" | "arbeidsgiver" | "samarbeidspartner"
	Chatbot            *bool
	Language           string // Locale, e.g. "nb"
	AvailableLanguages []AvailableLanguage
	LogoutWarning      *bool // Show logout warning if true
}

// Query is the decorator [Options] expressed as URL query parameters.
func (o *Options) Query() url.Values {
	q := url.Values{}
	q.Set("context", o.Context)
	if o.Chatbot != nil {
		q.Set("chatbot", strconv.FormatBool(*o.Chatbot))
	}
	if o.Language != "" {
		q.Set("language", o.Language)
	}
	if len(o.AvailableLanguages) > 0 {
		b, _ := json.Marshal(o.AvailableLanguages)
		q.Set("availableLanguages", string(b))
	}
	if o.LogoutWarning != nil {
		q.Set("logoutWarning", strconv.FormatBool(*o.LogoutWarning))
	}
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
