package test

import (
	"context"
	"log/slog"

	"github.com/navikt/hotbff/texas"
)

type IdentityProvider struct {
	AccessToken string
	Active      bool
	Err         error
}

func (idp *IdentityProvider) GetToken(_ context.Context, _ string) (*texas.TokenSet, error) {
	return &texas.TokenSet{AccessToken: idp.AccessToken}, idp.Err
}

func (idp *IdentityProvider) ExchangeToken(_ context.Context, _ string, _ string) (*texas.TokenSet, error) {
	return &texas.TokenSet{AccessToken: idp.AccessToken}, idp.Err
}

func (idp *IdentityProvider) IntrospectToken(_ context.Context, _ string) (*texas.TokenIntrospection, error) {
	return &texas.TokenIntrospection{Active: idp.Active}, idp.Err
}

func (idp *IdentityProvider) Set() bool {
	return idp != nil
}

func (idp *IdentityProvider) LogValue() slog.Value {
	return slog.StringValue("TestIdentityProvider")
}

func (idp *IdentityProvider) String() string {
	return "TestIdentityProvider"
}

func NewIdentityProvider(accessToken string, active bool, err error) *IdentityProvider {
	return &IdentityProvider{AccessToken: accessToken, Active: active, Err: err}
}
