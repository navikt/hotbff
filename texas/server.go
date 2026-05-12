package texas

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/navikt/hotbff/httpx"
)

// Authenticate is a middleware that checks for a bearer token in the Authorization header,
// and if present, validates it using the provided token introspector.
// If the token is valid, the user is marked as authenticated in the request context.
// If the token is invalid or absent, the user is marked as unauthenticated.
func Authenticate(idp TokenIntrospector, next http.Handler) http.Handler {
	if idp == nil {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		user := &User{Authenticated: false}

		if token, ok := TokenFromRequest(req); ok {
			ti, err := idp.IntrospectToken(ctx, token)
			if err != nil {
				if !errors.Is(err, context.Canceled) {
					log.ErrorContext(ctx, "token validation failed", "error", err)
				}
			} else if ti.Active {
				user = &User{Authenticated: true, Token: token}
			}
		} else {
			log.DebugContext(ctx, "token absent")
		}

		ctx = NewContext(ctx, user)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}

// Protected is a middleware that checks if the user is authenticated.
// If not, it returns 401 Unauthorized.
// If the user is authenticated, it calls the next handler.
func Protected(idp TokenIntrospector, next http.Handler) http.Handler {
	if idp == nil {
		return httpx.Unauthorized
	}
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		if isAuthenticated(ctx) {
			next.ServeHTTP(w, req)
		} else {
			log.DebugContext(ctx, "unauthenticated")
			httpx.ErrorCode(w, http.StatusUnauthorized)
		}
	})
}

// Redirect is a middleware that checks if the user is authenticated.
// If not, it redirects to the login page.
// If the user is authenticated, it calls the next handler.
func Redirect(idp TokenIntrospector, next http.Handler, basePath string) http.Handler {
	if idp == nil {
		return httpx.Unauthorized
	}
	return Authenticate(idp, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		if isAuthenticated(ctx) {
			next.ServeHTTP(w, req)
		} else {
			log.DebugContext(ctx, "unauthenticated")
			loginRedirect(w, req, basePath)
		}
	}))
}

// Validate is a handler that returns 200 OK if the user is authenticated, and 401 Unauthorized otherwise.
func Validate(idp IdentityProvider) http.Handler {
	if idp == nil {
		return httpx.Unauthorized
	}
	return Authenticate(idp, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		if isAuthenticated(ctx) {
			w.WriteHeader(http.StatusOK)
		} else {
			log.DebugContext(ctx, "unauthenticated")
			httpx.ErrorCode(w, http.StatusUnauthorized)
		}
	}))
}

func isAuthenticated(ctx context.Context) bool {
	user := FromContext(ctx)
	return user.Authenticated
}

func loginRedirect(w http.ResponseWriter, req *http.Request, basePath string) {
	ctx := req.Context()

	loginPath, err := url.JoinPath("/", basePath, "oauth2/login")
	if err != nil {
		http.Error(w, "invalid login path", http.StatusInternalServerError)
		return
	}

	loginURL, err := url.Parse(loginPath)
	if err != nil {
		http.Error(w, "invalid login url", http.StatusInternalServerError)
		return
	}

	q := loginURL.Query()
	q.Set("redirect", req.URL.RequestURI())
	loginURL.RawQuery = q.Encode()

	log.DebugContext(ctx, "login redirect", "loginURL", loginURL)

	http.Redirect(w, req, loginURL.String(), http.StatusTemporaryRedirect)
}
