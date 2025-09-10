package texas

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetToken(t *testing.T) {
	target := Target{"test", "test", "test"}.String()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v := req.FormValue(idpFormKey); v != string(EntraId) {
			t.Errorf("got identity_provider %q, want %q", v, EntraId)
		}
		if v := req.FormValue(targetFormKey); v != target {
			t.Errorf("got target %q, want %q", v, target)
		}
		_, _ = w.Write([]byte(`{"access_token":"accessToken"}`))
	}))
	defer server.Close()

	tokenURL = server.URL

	token, err := GetToken(EntraId, target)
	if err != nil {
		t.Fatal(err)
	}
	want := "accessToken"
	if token.AccessToken != want {
		t.Errorf("GetToken() = %q, want %q", token.AccessToken, want)
	}
}
