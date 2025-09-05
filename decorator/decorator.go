package decorator

import (
	"encoding/json"
	"net/http"
)

var (
	decoratorURL     = "http://nav-dekoratoren.personbruker/dekoratoren/ssr"
	decoratorURLDev  = "https://dekoratoren.ekstern.dev.nav.no/dekoratoren/ssr"
	decoratorURLProd = "https://www.nav.no/dekoratoren/dekoratoren/ssr"
)

//todo:caching
//todo:templating

func Get(r *Request) (*Response, error) {
	req, err := http.NewRequest(http.MethodGet, decoratorURLDev, nil)
	if err != nil {
		return nil, err
	}
	req.URL.Query().Set("context", r.Context)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer res.Body.Close()
	var p *Response
	err = json.NewDecoder(res.Body).Decode(&p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

type Request struct {
	Context string
}

type Response struct {
	HeadAssets string `json:"headAssets"`
	Header     string `json:"header"`
	Footer     string `json:"footer"`
	Scripts    string `json:"scripts"`
}
