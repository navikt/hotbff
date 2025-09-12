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
		assert.Equal(t, req.FormValue(idpFormKey), string(EntraId))
		assert.Equal(t, req.FormValue(targetFormKey), target)
		_, _ = w.Write([]byte(`{"access_token":"accessToken"}`))
	}))
	defer server.Close()

	tokenURL = server.URL

	ts, err := GetToken(EntraId, target)
	assert.Nil(t, err)
	assert.Equal(t, ts.AccessToken, "accessToken")
}

func TestExchangeToken(t *testing.T) {
	target := Target{"a", "b", "c"}.String()
	userToken := "userToken"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		assert.Nil(t, err)
		assert.Equal(t, req.FormValue(idpFormKey), string(TokenX))
		assert.Equal(t, req.FormValue(targetFormKey), target)
		assert.Equal(t, req.FormValue(userTokenFormKey), userToken)
		_, _ = w.Write([]byte(`{"access_token":"accessToken"}`))
	}))
	defer server.Close()

	tokenExchangeURL = server.URL

	ts, err := ExchangeToken(TokenX, target, userToken)
	assert.Nil(t, err)
	assert.Equal(t, ts.AccessToken, "accessToken")
}

func TestIntrospectToken(t *testing.T) {
	token := "token"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		assert.Nil(t, err)
		assert.Equal(t, req.FormValue(idpFormKey), string(IdPorten))
		assert.Equal(t, req.FormValue(tokenFormKey), token)
		_, _ = w.Write([]byte(`{"active":true}`))
	}))
	defer server.Close()

	tokenIntrospectionURL = server.URL

	ti, err := IntrospectToken(IdPorten, token)
	assert.Nil(t, err)
	assert.True(t, ti.Active)
}
