package texas

import (
	"net/http"
	"net/url"
)

// Redirect is a middleware that checks if the user is authenticated. If not, it redirects to the login page.
func Redirect(basePath string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		user := FromContext(ctx)
		if user.Authenticated {
			next.ServeHTTP(w, req)
		} else {
			log.DebugContext(ctx, "user unauthenticated, redirecting to login")
			loginRedirect(w, req, basePath)
		}
	})
}

func loginRedirect(w http.ResponseWriter, req *http.Request, basePath string) {
	ctx := req.Context()

	loginPath, err := url.JoinPath(basePath, "oauth2/login")
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
