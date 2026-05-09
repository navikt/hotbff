package texas

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var client = &http.Client{
	Timeout: 5 * time.Second,
}

func (idp IdentityProvider) GetToken(ctx context.Context, target string) (*TokenSet, error) {
	fv := newFormValues(idp)
	fv.Set(targetFormKey, target)
	var ts TokenSet
	if err := post(ctx, tokenURL, fv, &ts); err != nil {
		return nil, err
	}
	return &ts, nil
}

func (idp IdentityProvider) ExchangeToken(ctx context.Context, target string, userToken string) (*TokenSet, error) {
	fv := newFormValues(idp)
	fv.Set(targetFormKey, target)
	fv.Set(userTokenFormKey, userToken)
	var ts TokenSet
	if err := post(ctx, tokenExchangeURL, fv, &ts); err != nil {
		return nil, err
	}
	return &ts, nil
}

func (idp IdentityProvider) IntrospectToken(ctx context.Context, token string) (*TokenIntrospection, error) {
	fv := newFormValues(idp)
	fv.Set(tokenFormKey, token)
	var ti TokenIntrospection
	if err := post(ctx, tokenIntrospectionURL, fv, &ti); err != nil {
		return nil, err
	}
	return &ti, nil
}

func (idp IdentityProvider) Enabled() bool {
	return idp == EntraID || idp == IDPorten || idp == Maskinporten || idp == TokenX
}

func (idp IdentityProvider) LogValue() slog.Value {
	return slog.StringValue(string(idp))
}

func (idp IdentityProvider) String() string {
	return string(idp)
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

func newFormValues(idp IdentityProvider) url.Values {
	fv := url.Values{}
	fv.Set(idpFormKey, string(idp))
	return fv
}

func post(ctx context.Context, url string, fv url.Values, v any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(fv.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
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
