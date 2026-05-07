package texas

import (
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

func loginRedirect(w http.ResponseWriter, req *http.Request, basePath string) {
	ctx := req.Context()
	loginURL, _ := url.JoinPath(basePath, "/oauth2/login")

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
