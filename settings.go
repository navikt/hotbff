package hotbff

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

var defaultEnvVarKeys = []string{
	"NAIS_APP_NAME",
	"NAIS_CLUSTER_NAME",
	"USE_MSW",
	"GIT_COMMIT",
}

func SettingsHandler(envVarKeys []string) http.Handler {
	s := make(map[string]any)
	for _, key := range defaultEnvVarKeys {
		s[key] = parseEnv(key)
	}
	for _, key := range envVarKeys {
		s[key] = parseEnv(key)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		_, _ = fmt.Fprint(w, "window.appSettings = ")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		_ = enc.Encode(&s)
	})
}

func parseEnv(key string) any {
	v := os.Getenv(key)
	switch v {
	case "true":
		return true
	case "false":
		return false
	default:
		return v
	}
}
