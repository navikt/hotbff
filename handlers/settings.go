package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

var defaultEnvKeys = []string{
	"NAIS_APP_NAME",
	"NAIS_CLUSTER_NAME",
	"USE_MSW",
	"GIT_COMMIT",
}

func Settings(envKeys []string) http.Handler {
	s := make(map[string]any)
	keys := append(defaultEnvKeys, envKeys...)
	for _, key := range keys {
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
