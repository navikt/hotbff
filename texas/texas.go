package texas

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type IdentityProvider string

const (
	EntraId      IdentityProvider = "azuread"
	IdPorten     IdentityProvider = "idporten"
	Maskinporten IdentityProvider = "maskinporten"
	TokenX       IdentityProvider = "tokenx"
)

func GetToken(idp IdentityProvider, target string) (*TokenSet, error) {
	fv := newFormValues(idp)
	fv.Set(targetFormKey, target)
	var ts *TokenSet
	err := post(tokenURL, fv, &ts)
	if err != nil {
		return nil, err
	}
	return ts, nil
}

func ExchangeToken(idp IdentityProvider, target string, userToken string) (*TokenSet, error) {
	fv := newFormValues(idp)
	fv.Set(targetFormKey, target)
	fv.Set(userTokenFormKey, userToken)
	var ts *TokenSet
	err := post(tokenExchangeURL, fv, &ts)
	if err != nil {
		return nil, err
	}
	return ts, nil
}

func IntrospectToken(idp IdentityProvider, token string) (*TokenIntrospection, error) {
	fv := newFormValues(idp)
	fv.Set(tokenFormKey, token)
	var ti *TokenIntrospection
	err := post(tokenIntrospectionURL, fv, &ti)
	if err != nil {
		return nil, err
	}
	return ti, nil
}

type TokenSet struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type TokenIntrospection struct {
	Active bool `json:"active"`
}

const (
	idpFormKey       = "identity_provider"
	targetFormKey    = "target"
	tokenFormKey     = "token"
	userTokenFormKey = "user_token"
)

var (
	tokenURL              = os.Getenv("NAIS_TOKEN_ENDPOINT")
	tokenExchangeURL      = os.Getenv("NAIS_TOKEN_EXCHANGE_ENDPOINT")
	tokenIntrospectionURL = os.Getenv("NAIS_TOKEN_INTROSPECTION_ENDPOINT")
)

func newFormValues(idp IdentityProvider) *url.Values {
	data := &url.Values{}
	data.Set(idpFormKey, string(idp))
	return data
}

func post(url string, data *url.Values, v any) error {
	res, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("texas: unexpected statusCode: %d", res.StatusCode)
	}
	return json.NewDecoder(res.Body).Decode(v)
}
