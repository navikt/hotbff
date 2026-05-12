package texas

import (
	"context"
)

var (
	EntraID  IdentityProvider = &idp{provider: entraIDProvider, exchangeProvider: entraIDProvider}
	IDPorten IdentityProvider = &idp{provider: idPortenProvider, exchangeProvider: tokenXProvider}
	TokenX   IdentityProvider = &idp{provider: tokenXProvider, exchangeProvider: tokenXProvider}
)

type TokenSet struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type TokenIntrospection struct {
	Active bool `json:"active"`
}

type TokenGetter interface {
	// GetToken retrieves a token from the identity provider for the given target audience.
	// It returns a TokenSet containing the new token.
	GetToken(ctx context.Context, target string) (*TokenSet, error)
}

type TokenExchanger interface {
	// ExchangeToken exchanges the user's token for a new token from the identity provider for the given target audience.
	// It returns a TokenSet containing the new token.
	ExchangeToken(ctx context.Context, target string, userToken string) (*TokenSet, error)
}

type TokenIntrospector interface {
	// IntrospectToken validates the given token from the identity provider.
	// It returns a TokenIntrospection indicating whether the token is active.
	IntrospectToken(ctx context.Context, token string) (*TokenIntrospection, error)
}

type IdentityProvider interface {
	TokenGetter
	TokenExchanger
	TokenIntrospector
}
