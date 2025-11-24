package texas

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// Protected wraps a handler with token-based authentication using the provided [IdentityProvider].
// If the identity provider is not set, the handler is returned as is.
// If the token is missing or invalid, the user is redirected to the login page.
func Protected(idp IdentityProvider, basePath string, whitelist *WhitelistConfig, next http.Handler) http.Handler {
	if idp == "" {
		slog.Warn("texas: identity provider not set, token validation disabled")
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		if whitelisted, reason := isWhitelisted(req.URL.Path, basePath, whitelist); whitelisted {
			slog.DebugContext(ctx, "texas: path whitelisted, skipping authentication", "path", req.URL.Path, "reason", reason)
			next.ServeHTTP(w, req.WithContext(ctx))
			return
		}

		token, ok := TokenFromRequest(req)
		if !ok {
			slog.DebugContext(ctx, "texas: unauthorized: token missing", "idp", idp)
			loginRedirect(w, req, basePath)
			return
		}
		ti, err := IntrospectToken(ctx, idp, token)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				w.WriteHeader(http.StatusRequestTimeout)
			} else {
				slog.ErrorContext(ctx, "texas: error", "idp", idp, "error", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		if !ti.Active {
			slog.DebugContext(ctx, "texas: unauthorized: token invalid", "idp", idp)
			loginRedirect(w, req, basePath)
			return
		}
		ctx = NewContext(ctx, &User{Authenticated: true, Token: token})
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}

func isWhitelisted(urlPath string, basepath string, config *WhitelistConfig) (bool, string) {
	relativePath := strings.TrimPrefix(urlPath, strings.TrimSuffix(basepath, "/"))
	// Check exact path matches
	for _, whitelistPath := range config.WhitelistPaths {
		if relativePath == whitelistPath {
			return true, "exact path match: " + whitelistPath
		}
	}
	// Check path prefixes
	for _, prefix := range config.WhitelistPrefixes {
		if strings.HasPrefix(relativePath, prefix) {
			return true, "prefix match: " + prefix
		}
	}

	// Check file extensions
	for _, ext := range config.WhitelistExtensions {
		if strings.HasSuffix(relativePath, ext) {
			return true, "extension match: " + ext
		}
	}

	return false, ""
}

func loginRedirect(w http.ResponseWriter, req *http.Request, basePath string) {
	ctx := req.Context()
	loginURL := path.Join(basePath, "/oauth2/login")

	returnToURL := req.URL.Path
	if req.URL.RawQuery != "" {
		returnToURL = returnToURL + "?" + req.URL.RawQuery
	}
	redirectBase := strings.TrimSuffix(basePath, "/")
	loginURL = loginURL + "?redirect=" + redirectBase + url.QueryEscape(returnToURL)

	slog.InfoContext(ctx, "texas: login redirect",
		"loginURL", loginURL,
		"returnToURL", returnToURL,
		"req.URL.Path", req.URL.Path,
		"basePath", basePath)
	http.Redirect(w, req, loginURL, http.StatusTemporaryRedirect)
}
