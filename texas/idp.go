package texas

import "context"

type IdentityProvider string

const (
	EntraID      IdentityProvider = "azuread"
	IDPorten     IdentityProvider = "idporten"
	Maskinporten IdentityProvider = "maskinporten"
	TokenX       IdentityProvider = "tokenx"

	None IdentityProvider = ""
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
	// It returns a [TokenSet] containing the new token.
	GetToken(ctx context.Context, target string) (*TokenSet, error)
	Enabled() bool
}

type TokenExchanger interface {
	// ExchangeToken exchanges the user's token for a new token from the identity provider for the given target audience.
	// It returns a [TokenSet] containing the new token.
	ExchangeToken(ctx context.Context, target string, userToken string) (*TokenSet, error)
	Enabled() bool
}

type TokenIntrospector interface {
	// IntrospectToken validates the given token from the identity provider.
	// It returns a [TokenIntrospection] indicating whether the token is active.
	IntrospectToken(ctx context.Context, token string) (*TokenIntrospection, error)
	Enabled() bool
}
