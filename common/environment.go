package common

import (
	"net/url"
	"os"
)

var (
	AppName     = os.Getenv("NAIS_APP_NAME")
	ClusterName = os.Getenv("NAIS_CLUSTER_NAME")
	Namespace   = os.Getenv("NAIS_NAMESPACE")

	Port        = os.Getenv("PORT")
	BindAddress = os.Getenv("BIND_ADDRESS")
)

func MustGetEnvURL(key string) *url.URL {
	urlStr := os.Getenv(key)
	if urlStr == "" {
		Fatal("missing environment variable", "key", key)
	}
	u, err := url.Parse(urlStr)
	if err != nil {
		Fatal("invalid url", "urlStr", urlStr, "key", key, "error", err)
	}
	return u
}
