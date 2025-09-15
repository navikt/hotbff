package hotbff

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func settingsHandler(envKeys *[]string) http.Handler {
	s := make(map[string]any)
	allEnvKeys := defaultEnvKeys
	if envKeys != nil {
		allEnvKeys = append(defaultEnvKeys, *envKeys...)
	}
	for _, key := range allEnvKeys {
		s[key] = parseEnv(key)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		_, err := fmt.Fprint(w, "window.appSettings = ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		err = enc.Encode(&s)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
