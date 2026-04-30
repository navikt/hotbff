package test

import (
	"context"
	"testing"

	"github.com/navikt/hotbff/internal/assert"
	"github.com/navikt/hotbff/texas"
)

type GetTokenFunc func(ctx context.Context, target string) (*texas.TokenSet, error)

func (f GetTokenFunc) GetToken(ctx context.Context, target string) (*texas.TokenSet, error) {
	return f(ctx, target)
}

func NewTokenGetter(t *testing.T, expectedTarget, accessToken string) texas.TokenGetter {
	t.Helper()
	return GetTokenFunc(func(ctx context.Context, target string) (*texas.TokenSet, error) {
		assert.Equal(t, target, expectedTarget)
		return &texas.TokenSet{AccessToken: accessToken}, nil
	})
}

type ExchangeTokenFunc func(ctx context.Context, target string, userToken string) (*texas.TokenSet, error)

func (f ExchangeTokenFunc) ExchangeToken(ctx context.Context, target string, userToken string) (*texas.TokenSet, error) {
	return f(ctx, target, userToken)
}

func NewTokenExchanger(t *testing.T, expectedTarget, expectedUserToken, accessToken string) texas.TokenExchanger {
	t.Helper()
	return ExchangeTokenFunc(func(ctx context.Context, target string, userToken string) (*texas.TokenSet, error) {
		assert.Equal(t, target, expectedTarget)
		assert.Equal(t, userToken, expectedUserToken)
		return &texas.TokenSet{AccessToken: accessToken}, nil
	})
}

type IntrospectTokenFunc func(ctx context.Context, token string) (*texas.TokenIntrospection, error)

func (f IntrospectTokenFunc) IntrospectToken(ctx context.Context, token string) (*texas.TokenIntrospection, error) {
	return f(ctx, token)
}

func NewTokenIntrospector(t *testing.T, expectedToken string, active bool) texas.TokenIntrospector {
	t.Helper()
	return IntrospectTokenFunc(func(ctx context.Context, token string) (*texas.TokenIntrospection, error) {
		assert.Equal(t, token, expectedToken)
		return &texas.TokenIntrospection{Active: active}, nil
	})
}
