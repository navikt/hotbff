package texas

import "strings"

func (opts *Options) isPublic(urlPath string, basePath string) (bool, string) {
	relativePath := strings.TrimPrefix(urlPath, strings.TrimSuffix(basePath, "/"))
	// Check file extensions
	for extension := range opts.PublicExtensions {
		if strings.HasSuffix(relativePath, extension) {
			return true, "extension match: " + extension
		}
	}
	// Check exact path matches
	for path := range opts.PublicPaths {
		if relativePath == path {
			return true, "exact path match: " + path
		}
	}
	// Check path prefixes
	for prefix := range opts.PublicPrefixes {
		if strings.HasPrefix(relativePath, prefix) {
			return true, "prefix match: " + prefix
		}
	}
	return false, ""
}
