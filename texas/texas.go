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
	data := &url.Values{}
	data.Set(idpKey, string(idp))
	data.Set(targetKey, target)
	var v *TokenSet
	err := post(tokenURL, data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func ExchangeToken(idp IdentityProvider, target string, userToken string) (*TokenSet, error) {
	data := &url.Values{}
	data.Set(idpKey, string(idp))
	data.Set(targetKey, target)
	data.Set(userTokenKey, userToken)
	var v *TokenSet
	err := post(tokenExchangeURL, data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func IntrospectToken(idp IdentityProvider, token string) (*TokenIntrospection, error) {
	data := &url.Values{}
	data.Set(idpKey, string(idp))
	data.Set(tokenKey, token)
	var v *TokenIntrospection
	err := post(tokenIntrospectionURL, data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
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
	idpKey       = "identity_provider"
	targetKey    = "target"
	tokenKey     = "token"
	userTokenKey = "user_token"
)

var (
	tokenURL              = os.Getenv("NAIS_TOKEN_ENDPOINT")
	tokenExchangeURL      = os.Getenv("NAIS_TOKEN_EXCHANGE_ENDPOINT")
	tokenIntrospectionURL = os.Getenv("NAIS_TOKEN_INTROSPECTION_ENDPOINT")
)

func post(url string, data *url.Values, v any) error {
	res, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected statusCode: %d", res.StatusCode)
	}
	return json.NewDecoder(res.Body).Decode(v)
}
