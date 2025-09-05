package texas

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		_ = req.ParseForm()
		fmt.Println(req.FormValue("identity_provider"))
		fmt.Println(req.FormValue("target"))

		_, _ = w.Write([]byte(`{"access_token":"token"}`))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tokenURL = server.URL

	token, err := GetToken(IdentityProviderEntraId, "test")
	if err != nil {
		t.Fatal(err)
	}
	want := "token"
	if token.AccessToken != want {
		t.Errorf("GetToken() = %q, want %q", token.AccessToken, want)
	}
}
