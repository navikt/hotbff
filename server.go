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
	"github.com/navikt/hotbff/middleware"
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
	BasePath      string                  // the base path to serve the application on (defaults to "/")
	RootDir       string                  // the directory to serve static files from (defaults to "dist")
	DecoratorOpts *decorator.Options      // options for the HTML decorator
	Proxy         proxy.Map               // map of proxy options keyed by URL prefix
	IDP           texas.TokenIntrospector // identity provider to use for token introspection (if nil, no validation is performed)
	TexasOpts     *texas.Options          // options for Texas
	EnvKeys       []string                // list of environment variable keys to expose to the frontend (via "/settings.js")
}

// Start starts the HTTP server with the given [Options].
func Start(mux *http.ServeMux, opts *Options) {
	if mux == nil {
		mux = http.DefaultServeMux
	}
	slog.Info("hotbff: starting server", "address", strings.Join([]string{bindAddressToLog(addr), opts.BasePath}, ""), "basePath", opts.BasePath, "rootDir", opts.RootDir)
	Configure(mux, opts)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		slog.Error("hotbff: server startup failed", "error", err)
		os.Exit(1)
	}
}

// Configure configures the given [http.ServeMux] with the [Options] provided.
func Configure(mux *http.ServeMux, opts *Options) {
	if mux == nil {
		mux = http.DefaultServeMux
	}

	basePath := opts.BasePath
	if basePath == "" {
		basePath = "/"
	}
	rootDir := opts.RootDir
	if rootDir == "" {
		rootDir = "dist"
	}

	texasOpts := opts.TexasOpts
	if texasOpts == nil {
		texasOpts = texas.DefaultOptions()
	} else {
		texasOpts = texasOpts.Clone()
		texasOpts.AddDefaults()
	}

	// / (public)
	mux.Handle("GET /isalive", middleware.Health("ALIVE"))
	mux.Handle("GET /isready", middleware.Health("READY"))

	// /base/path/ (public)
	baseMux := http.NewServeMux()
	baseMux.Handle("GET /settings.js", settingsHandler(basePath, opts.EnvKeys))
	// baseMux.Handle("GET /auth/status", opts.IDP.Status())

	// /base/path/ (protected)
	protectedMux := http.NewServeMux()
	protectedMux.Handle("/", staticHandler(rootDir, opts.DecoratorOpts))

	// /base/path/proxy/prefix/ (protected)
	err := proxy.Configure(protectedMux, opts.Proxy)
	if err != nil {
		slog.Error("hotbff: failed to configure proxy", "error", err)
		os.Exit(1)
	}

	baseMux.Handle("/", texas.Protected(opts.IDP, opts.TexasOpts, basePath, protectedMux))
	mux.Handle(basePath, maybeStripPrefix(path.Join(basePath), baseMux))
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
