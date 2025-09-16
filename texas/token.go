package texas

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

const HeaderAuthorization = "Authorization"

// TokenFromRequest extracts a bearer token from the Authorization header of an HTTP request.
// It returns the token string and a boolean indicating whether a bearer token was present.
// This function does not validate the token in any way.
func TokenFromRequest(req *http.Request) (token string, ok bool) {
	token, ok = strings.CutPrefix(req.Header.Get(HeaderAuthorization), "Bearer ")
	if token == "" {
		ok = false
	}
	if !ok {
		token = ""
	}
	return
}

type JWT struct {
	Header    map[string]any
	Claims    map[string]any
	Signature []byte
}

// ParseJWT parses a JWT string into its components: header, claims, and signature.
// It returns a JWT struct and an error if the parsing fails.
func ParseJWT(jwtStr string) (*JWT, error) {
	parts := strings.Split(jwtStr, ".")
	if len(parts) != 3 {
		return nil, errors.New("texas: invalid jwt")
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
