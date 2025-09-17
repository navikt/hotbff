package hotbff

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func settingsHandler(basePath string, envKeys []string) http.Handler {
	s := map[string]any{"BASE_PATH": basePath}
	for _, key := range defaultEnvKeys {
		s[key] = parseEnv(key)
	}
	for _, key := range envKeys {
		s[key] = parseEnv(key)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		if _, err := fmt.Fprint(w, "window.appSettings = "); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		if err := enc.Encode(s); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func parseEnv(key string) any {
	v := os.Getenv(key)
	switch v {
	case "":
		return nil
	case "true":
		return true
	case "false":
		return false
	default:
		return v
	}
}

var defaultEnvKeys = []string{
	"NAIS_APP_NAME",
	"NAIS_CLUSTER_NAME",
	"USE_MSW",
	"GIT_COMMIT",
}
