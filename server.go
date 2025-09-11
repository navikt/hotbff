package hotbff

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"

	"github.com/navikt/hotbff/decorator"
	"github.com/navikt/hotbff/proxy"
	"github.com/navikt/hotbff/texas"
)

var (
	addr = os.Getenv("BIND_ADDRESS")
)

func init() {
	if addr == "" {
		addr = ":9000"
	}
}

type ServerOptions struct {
	BasePath      string
	RootDir       string
	DecoratorOpts *decorator.Options
	Proxy         *proxy.Map
	IDP           texas.IdentityProvider
	EnvKeys       *[]string
}

func (o ServerOptions) basePath() string {
	if o.BasePath == "" {
		return "/"
	}
	return o.BasePath
}

func (o ServerOptions) rootDir() string {
	if o.RootDir == "" {
		return "dist"
	}
	return o.RootDir
}

func StartServer(opts *ServerOptions) {
	basePath := opts.basePath()
	rootDir := opts.rootDir()

	// public routes
	http.Handle("GET /isalive", healthHandler("ALIVE"))
	http.Handle("GET /isready", healthHandler("READY"))

	envKeys := []string{}
	if opts.EnvKeys != nil {
		envKeys = *opts.EnvKeys
	}
	http.Handle(fmt.Sprintf("GET %s", path.Join(basePath, "settings.js")), settingsJS(envKeys))

	// (potentially) protected routes
	mux := http.NewServeMux()
	mux.Handle(basePath, maybeStripPrefix(basePath, rootHandler(rootDir, opts.DecoratorOpts)))

	if opts.Proxy != nil {
		for proxyPrefix, proxyOpts := range *opts.Proxy {
			proxyPath := path.Join(basePath, proxyPrefix)
			slog.Info("hotbff: adding proxy", "prefix", proxyPrefix, "proxyPath", proxyPath, "target", proxyOpts.Target)
			mux.Handle(proxyPath, maybeStripPrefix(basePath, proxyOpts.Handler(proxyPrefix)))
		}
	}

	http.Handle("/", texas.Protected(opts.IDP, basePath, mux))

	slog.Info("hotbff: starting server", "address", addr, "basePath", basePath, "rootDir", rootDir)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		slog.Error("hotbff: server startup failed", "error", err)
		os.Exit(1)
	}
}

func maybeStripPrefix(prefix string, h http.Handler) http.Handler {
	if prefix == "/" {
		return h
	}
	return http.StripPrefix(prefix, h)
}
