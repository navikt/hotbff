package texas

import (
	"context"
	"log/slog"
	"net/http"
)

type testIDP struct {
	AccessToken string
	Active      bool
	Err         error
}

func (idp *testIDP) GetToken(_ context.Context, _ string) (*TokenSet, error) {
	return &TokenSet{AccessToken: idp.AccessToken}, idp.Err
}

func (idp *testIDP) ExchangeToken(_ context.Context, _ string, _ string) (*TokenSet, error) {
	return &TokenSet{AccessToken: idp.AccessToken}, idp.Err
}

func (idp *testIDP) IntrospectToken(_ context.Context, _ string) (*TokenIntrospection, error) {
	return &TokenIntrospection{Active: idp.Active}, idp.Err
}

func (idp *testIDP) Status() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		user := FromContext(ctx)
		if user.Authenticated {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	})
}

func (idp *testIDP) Enabled() bool {
	return idp != nil
}

func (idp *testIDP) LogValue() slog.Value {
	return slog.StringValue("testIDP")
}

func (idp *testIDP) String() string {
	return "testIDP"
}

func NewTestIDP(accessToken string, active bool, err error) *testIDP {
	return &testIDP{AccessToken: accessToken, Active: active, Err: err}
}
