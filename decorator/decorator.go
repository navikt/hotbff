package decorator

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

func Get(opts *Options) (*Elements, error) {
	req, err := http.NewRequest(http.MethodGet, getDecoratorURL(), nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("context", opts.Context)
	req.URL.RawQuery = q.Encode()
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

type Elements struct {
	HeadAssets template.HTML `json:"headAssets"`
	Header     template.HTML `json:"header"`
	Footer     template.HTML `json:"footer"`
	Scripts    template.HTML `json:"scripts"`
}

var (
	decoratorURL    = "http://nav-dekoratoren.personbruker/dekoratoren/ssr"
	decoratorURLDev = "https://dekoratoren.ekstern.dev.nav.no/dekoratoren/ssr"
	//decoratorURLProd = "https://www.nav.no/dekoratoren/ssr"
)

func getDecoratorURL() string {
	switch os.Getenv("NAIS_CLUSTER_NAME") {
	case "", "local", "test":
		return decoratorURLDev
	default:
		return decoratorURL
	}
}
