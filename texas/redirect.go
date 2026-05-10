package texas

import (
	"net/http"
	"net/url"
)

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
