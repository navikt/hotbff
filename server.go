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
	mux.Handle("/", rootHandler(rootDir, opts.DecoratorOpts))
	if opts.Proxy != nil {
		for prefix, t := range *opts.Proxy {
			slog.Info("hotbff: adding proxy", "prefix", prefix, "target", t.Target)
			mux.Handle(prefix, t.Handler(prefix))
		}
	}

	var r http.Handler
	if opts.IDP == "" {
		r = mux
	} else {
		r = texas.Protected(opts.IDP, mux)
	}
	http.Handle(basePath, http.StripPrefix(basePath, r))

	addr := os.Getenv("BIND_ADDRESS")
	if addr == "" {
		addr = ":5000"
	}
	slog.Info("hotbff: starting server", "address", addr, "basePath", basePath, "rootDir", rootDir)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		slog.Error("hotbff: server startup failed", "error", err)
		os.Exit(1)
	}
}
