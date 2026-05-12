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
	"time"
)

const (
	ContextPrivatperson      = "privatperson"
	ContextArbeidsgiver      = "arbeidsgiver"
	ContextSamarbeidspartner = "samarbeidspartner"
)

var (
	cluster = os.Getenv("NAIS_CLUSTER_NAME")

	decoratorURLProd = "http://nav-dekoratoren.personbruker/dekoratoren/ssr"
	decoratorURLDev  = "https://dekoratoren.ekstern.dev.nav.no/dekoratoren/ssr"
	decoratorURL     = getDecoratorURL()

	log = slog.Default().With("package", "decorator", "decoratorURL", decoratorURL)

	client = &http.Client{
		Timeout: 5 * time.Second,
	}
)

// Fetch retrieves decorator elements using the given options.
func Fetch(ctx context.Context, opts *Options) (*Elements, error) {
	if opts == nil {
		opts = &Options{}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, decoratorURL, nil)
	if err != nil {
		return nil, fmt.Errorf("decorator: %w", err)
	}
	req.URL.RawQuery = opts.Query().Encode()

	log.DebugContext(ctx, "fetching elements", "url", req.URL)
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("decorator: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("decorator: unexpected response, statusCode: %d", res.StatusCode)
	}

	var elems Elements
	if err := json.NewDecoder(res.Body).Decode(&elems); err != nil {
		return nil, fmt.Errorf("decorator: %w", err)
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
	Context            string              // The context, e.g. "privatperson" | "arbeidsgiver" | "samarbeidspartner".
	Chatbot            *bool               // Enable the chatbot if true.
	Language           string              // Locale, e.g. "nb".
	AvailableLanguages []AvailableLanguage // Available languages for the language selector in the decorator.
	LogoutWarning      *bool               // Show a logout warning if true.
}

// Query is the decorator options expressed as URL query parameters.
func (opts *Options) Query() url.Values {
	q := url.Values{}
	if opts == nil {
		return q
	}

	q.Set("context", opts.Context)

	if opts.Chatbot != nil {
		q.Set("chatbot", strconv.FormatBool(*opts.Chatbot))
	}

	if opts.Language != "" {
		q.Set("language", opts.Language)
	}

	if len(opts.AvailableLanguages) > 0 {
		b, _ := json.Marshal(opts.AvailableLanguages)
		q.Set("availableLanguages", string(b))
	}

	if opts.LogoutWarning != nil {
		q.Set("logoutWarning", strconv.FormatBool(*opts.LogoutWarning))
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

func getDecoratorURL() string {
	switch cluster {
	case "", "local", "test":
		return decoratorURLDev
	default:
		return decoratorURLProd
	}
}
