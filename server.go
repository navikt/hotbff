package hotbff

import (
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

// Options are the options for the server.
type Options struct {
	// BasePath is the base path to serve the application on. Defaults to "/".
	BasePath string
	// RootDir is the directory to serve static files from. Defaults to "dist".
	RootDir string
	// DecoratorOpts are the options for the HTML decorator.
	DecoratorOpts *decorator.Options
	// Proxy is the map of proxy options.
	Proxy *proxy.Map
	// IDP is the identity provider to use for token validation. If empty, no validation is performed.
	IDP texas.IdentityProvider
	// EnvKeys is the list of environment variable keys to expose to the frontend.
	EnvKeys *[]string
}

// Start starts the HTTP server with the given options.
func Start(opts *Options) {
	slog.Info("hotbff: starting server", "address", addr, "basePath", opts.BasePath, "rootDir", opts.RootDir)
	err := http.ListenAndServe(addr, Handler(opts))
	if err != nil {
		slog.Error("hotbff: server startup failed", "error", err)
		os.Exit(1)
	}
}

// Handler returns an http.Handler that serves the application with the given options.
func Handler(opts *Options) http.Handler {
	basePath := opts.BasePath
	rootDir := opts.RootDir

	if basePath == "" {
		basePath = "/"
	}
	if rootDir == "" {
		rootDir = "dist"
	}

	// / (public)
	rootMux := http.NewServeMux()
	rootMux.Handle("GET /isalive", healthHandler("ALIVE"))
	rootMux.Handle("GET /isready", healthHandler("READY"))

	// /base/path/ (public)
	baseMux := http.NewServeMux()
	baseMux.Handle("GET /settings.js", settingsHandler(opts.EnvKeys))

	// /base/path/ (protected)
	protectedMux := http.NewServeMux()
	protectedMux.Handle("/", staticHandler(rootDir, opts.DecoratorOpts))

	// /base/path/proxy/prefix/ (protected)
	proxy.Configure(opts.Proxy, protectedMux)

	baseMux.Handle("/", texas.Protected(opts.IDP, basePath, protectedMux))
	rootMux.Handle(basePath, maybeStripPrefix(path.Join(basePath), baseMux))
	return rootMux
}

func maybeStripPrefix(prefix string, h http.Handler) http.Handler {
	if prefix == "/" {
		return h
	}
	return http.StripPrefix(prefix, h)
}
