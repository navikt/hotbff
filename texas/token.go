package texas

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// ErrInvalidJWT is returned when a JWT-string cannot be parsed.
var ErrInvalidJWT = errors.New("texas: invalid jwt")

// HeaderAuthorization is the HTTP Authorization header name.
const HeaderAuthorization string = "Authorization"

// TokenFromRequest extracts the bearer token from the Authorization header of an [http.Request].
// It returns the token string and a boolean indicating whether a bearer token was found.
//
// NB! This function does not validate the token in any way.
func TokenFromRequest(req *http.Request) (token string, ok bool) {
	h := req.Header.Get(HeaderAuthorization)
	token, ok = strings.CutPrefix(h, "Bearer ")
	if token == "" {
		ok = false
	}
	if !ok {
		token = ""
	}
	return
}

// JWT represents the parsed components of a JSON Web Token.
type JWT struct {
	Header    map[string]any // Header contains the JWT header claims.
	Claims    map[string]any // Claims contains the JWT payload claims.
	Signature []byte         // Signature is the JWT signature.
}

// ParseJWT parses a JWT-string into its components: header, claims, and signature.
// It returns a [JWT] or an [error] if the parsing fails.
func ParseJWT(jwtStr string) (*JWT, error) {
	parts := strings.Split(jwtStr, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidJWT
	}
	h, err := parseJWTPart(parts[0])
	if err != nil {
		return nil, err
	}
	c, err := parseJWTPart(parts[1])
	if err != nil {
		return nil, err
	}
	s, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, err
	}
	return &JWT{Header: h, Claims: c, Signature: s}, nil
}

func parseJWTPart(base64Str string) (map[string]any, error) {
	data, err := base64.RawURLEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, err
	}
	v := make(map[string]any)
	err = json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}
