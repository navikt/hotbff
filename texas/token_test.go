package texas

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/navikt/hotbff/internal/assert"
)

const jwtStr = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.KMUFsIDTnFmyG3nMiGM6H9FNFUROf3wh7SmqJp-QV30"

func TestTokenFromRequestPresent(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	req.Header.Set(HeaderAuthorization, "Bearer "+jwtStr)
	token, ok := TokenFromRequest(req)
	assert.Equal(t, token, jwtStr)
	assert.True(t, ok)
}

func TestTokenFromRequestMissing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	req.Header.Set(HeaderAuthorization, "Bearer ")
	token, ok := TokenFromRequest(req)
	assert.Equal(t, token, "")
	assert.False(t, ok)

	req.Header.Set(HeaderAuthorization, "Bearer")
	token, ok = TokenFromRequest(req)
	assert.Equal(t, token, "")
	assert.False(t, ok)

	req.Header.Del(HeaderAuthorization)
	token, ok = TokenFromRequest(req)
	assert.Equal(t, token, "")
	assert.False(t, ok)
}

func TestParseJWT(t *testing.T) {
	j, err := ParseJWT(jwtStr)
	assert.Nil(t, err)
	assert.Equal(t, j.Header["alg"], "HS256")
	assert.Equal(t, j.Claims["sub"], "1234567890")
}
