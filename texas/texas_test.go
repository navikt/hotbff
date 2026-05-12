package texas

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/navikt/hotbff/internal/assert"
)

func TestGetToken(t *testing.T) {
	target := Target{"a", "b", "c"}.String()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		assert.Nil(t, err)
		assert.Equal(t, req.FormValue(providerFormKey), entraIDProvider)
		assert.Equal(t, req.FormValue(targetFormKey), target)
		_, _ = w.Write([]byte(`{"access_token":"accessToken"}`))
	}))
	defer server.Close()

	tokenURL = server.URL

	ts, err := EntraID.GetToken(t.Context(), target)
	assert.Nil(t, err)
	assert.Equal(t, ts.AccessToken, "accessToken")
}

func TestExchangeToken(t *testing.T) {
	target := Target{"a", "b", "c"}.String()
	userToken := "userToken"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		assert.Nil(t, err)
		assert.Equal(t, req.FormValue(providerFormKey), tokenXProvider)
		assert.Equal(t, req.FormValue(targetFormKey), target)
		assert.Equal(t, req.FormValue(userTokenFormKey), userToken)
		_, _ = w.Write([]byte(`{"access_token":"accessToken"}`))
	}))
	defer server.Close()

	tokenExchangeURL = server.URL

	ts, err := TokenX.ExchangeToken(t.Context(), target, userToken)
	assert.Nil(t, err)
	assert.Equal(t, ts.AccessToken, "accessToken")
}

func TestIntrospectToken(t *testing.T) {
	token := "token"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		assert.Nil(t, err)
		assert.Equal(t, req.FormValue(providerFormKey), idPortenProvider)
		assert.Equal(t, req.FormValue(tokenFormKey), token)
		_, _ = w.Write([]byte(`{"active":true}`))
	}))
	defer server.Close()

	tokenIntrospectionURL = server.URL

	ti, err := IDPorten.IntrospectToken(t.Context(), token)
	assert.Nil(t, err)
	assert.True(t, ti.Active)
}
