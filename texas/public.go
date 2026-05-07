package texas

import (
	"fmt"
	"net/http"
	"path"
	"strings"
)

func isPublic(req *http.Request, basePath string, opts *Options) (bool, string) {
	urlPath := req.URL.Path
	relativePath := strings.TrimPrefix(urlPath, strings.TrimSuffix(basePath, "/"))
	// Check file extensions
	if extension := path.Ext(urlPath); opts.PublicExtensions.Has(extension) {
		return true, fmt.Sprintf("extension match: %q", extension)
	}
	// Check exact path matches
	if opts.PublicPaths.Has(relativePath) {
		return true, fmt.Sprintf("exact path match: %q", relativePath)
	}
	// Check path prefixes
	for prefix := range opts.PublicPrefixes {
		if strings.HasPrefix(relativePath, prefix) {
			return true, fmt.Sprintf("prefix match: %q", prefix)
		}
	}
	return false, ""
}
