package hotbff

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path"
	"strings"

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

// Options for the server.
type Options struct {
	BasePath      string                 // the base path to serve the application on (defaults to "/")
	RootDir       string                 // the directory to serve static files from (defaults to "dist")
	DecoratorOpts *decorator.Options     // options for the HTML decorator
	Proxy         proxy.Map              // map of proxy options keyed by URL prefix
	IDP           texas.IdentityProvider // identity provider to use for token validation (if empty, no validation is performed)
	EnvKeys       []string               // list of environment variable keys to expose to the frontend (via "/settings.js")
}

// Start starts the HTTP server with the given [Options].
func Start(opts *Options) {
	slog.Info("hotbff: starting server", "address", strings.Join([]string{bindAddressToLog(addr), opts.BasePath}, ""), "basePath", opts.BasePath, "rootDir", opts.RootDir)
	err := http.ListenAndServe(addr, Handler(opts))
	if err != nil {
		slog.Error("hotbff: server startup failed", "error", err)
		os.Exit(1)
	}
}

// Handler returns a handler that serves the application with the given [Options].
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
	baseMux.Handle("GET /settings.js", settingsHandler(basePath, opts.EnvKeys))

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

func bindAddressToLog(bindAddr string) string {
	host, port, err := net.SplitHostPort(bindAddr)
	if err != nil {
		// If we can't split, fallback to default
		return "http://localhost" + bindAddr
	}

	// Handle empty or wildcard host
	if host == "" || host == "0.0.0.0" || host == "127.0.0.1" {
		host = "localhost"
	}

	return fmt.Sprintf("http://%s:%s", host, port)
}
