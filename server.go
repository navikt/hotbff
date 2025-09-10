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

type ServerOptions struct {
	BasePath      string
	RootDir       string
	DecoratorOpts *decorator.Options
	Proxy         *proxy.Map
	IDP           texas.IdentityProvider
	EnvKeys       *[]string
}

func StartServer(opts *ServerOptions) {
	basePath := opts.BasePath
	if basePath == "" {
		basePath = "/"
	}
	rootDir := opts.RootDir
	if rootDir == "" {
		rootDir = "dist"
	}

	http.Handle("GET /isalive", statusHandler("ALIVE"))
	http.Handle("GET /isready", statusHandler("READY"))

	envKeys := []string{}
	if opts.EnvKeys != nil {
		envKeys = *opts.EnvKeys
	}
	http.Handle(fmt.Sprintf("GET %s", path.Join(basePath, "/settings.js")), settingsJS(envKeys))

	mux := http.NewServeMux()
	if opts.Proxy != nil {
		for prefix, pOpts := range *opts.Proxy {
			slog.Info("hotbff: adding proxy", "prefix", prefix, "target", pOpts.Target)
			mux.Handle(prefix, pOpts.Handler(prefix))
		}
	}
	mux.Handle("GET /", rootHandler(rootDir, opts.DecoratorOpts))

	var h http.Handler = mux
	if opts.IDP != "" {
		h = texas.Protected(opts.IDP, mux)
	}
	if basePath == "/" {
		http.Handle(basePath, h)
	} else {
		http.Handle(basePath, http.StripPrefix(basePath, h))
	}

	addr := os.Getenv("BIND_ADDRESS")
	if addr == "" {
		addr = ":9000"
	}
	slog.Info("hotbff: starting server", "address", addr, "basePath", basePath, "rootDir", rootDir)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		slog.Error("hotbff: server startup failed", "error", err)
		os.Exit(1)
	}
}
