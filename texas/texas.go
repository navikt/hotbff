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

	"github.com/navikt/hotbff/httpx"
)

const (
	entraIDProvider  = "azuread"
	idPortenProvider = "idporten"
	tokenXProvider   = "tokenx"

	providerFormKey  = "identity_provider"
	targetFormKey    = "target"
	tokenFormKey     = "token"
	userTokenFormKey = "user_token"
)

var (
	tokenURL              = os.Getenv("NAIS_TOKEN_ENDPOINT")
	tokenExchangeURL      = os.Getenv("NAIS_TOKEN_EXCHANGE_ENDPOINT")
	tokenIntrospectionURL = os.Getenv("NAIS_TOKEN_INTROSPECTION_ENDPOINT")

	client = &http.Client{
		Timeout: 5 * time.Second,
	}
)

type idp struct {
	provider         string
	exchangeProvider string
}

func (idp *idp) GetToken(ctx context.Context, target string) (*TokenSet, error) {
	fv := newFormValues(idp.provider)
	fv.Set(targetFormKey, target)
	var ts TokenSet
	if err := post(ctx, tokenURL, fv, &ts); err != nil {
		return nil, fmt.Errorf("texas: token retrieval failed: %w", err)
	}
	return &ts, nil
}

func (idp *idp) ExchangeToken(ctx context.Context, target string, userToken string) (*TokenSet, error) {
	fv := newFormValues(idp.exchangeProvider)
	fv.Set(targetFormKey, target)
	fv.Set(userTokenFormKey, userToken)
	var ts TokenSet
	if err := post(ctx, tokenExchangeURL, fv, &ts); err != nil {
		return nil, fmt.Errorf("texas: token exchange failed: %w", err)
	}
	return &ts, nil
}

func (idp *idp) IntrospectToken(ctx context.Context, token string) (*TokenIntrospection, error) {
	fv := newFormValues(idp.provider)
	fv.Set(tokenFormKey, token)
	var ti TokenIntrospection
	if err := post(ctx, tokenIntrospectionURL, fv, &ti); err != nil {
		return nil, fmt.Errorf("texas: token validation failed: %w", err)
	}
	return &ti, nil
}

func (idp *idp) LogValue() slog.Value {
	return slog.StringValue(idp.provider)
}

func newFormValues(provider string) url.Values {
	fv := url.Values{}
	fv.Set(providerFormKey, provider)
	return fv
}

func post(ctx context.Context, url string, fv url.Values, v any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(fv.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set(httpx.HeaderAccept, httpx.ContentTypeApplicationJSON)
	req.Header.Set(httpx.HeaderContentType, httpx.ContentTypeFormUrlEncoded)
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
