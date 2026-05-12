package hotbff

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/navikt/hotbff/decorator"
	"github.com/navikt/hotbff/httpx"
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
	BasePath      string                 // The base path to serve the application on (defaults to "/").
	RootDir       string                 // The directory to serve static files from (defaults to "dist").
	DecoratorOpts *decorator.Options     // Options for the HTML decorator.
	Proxy         proxy.Map              // Map of proxy options keyed by URL prefix.
	IDP           texas.IdentityProvider // Identity provider to use for token validation and exchange (if nil, no validation or exhange is performed).
	PublicPaths   []string               // Paths that should be publicly accessible without login redirection (only relevant if IDP is set).
	EnvKeys       []string               // Environment variables that should be exposed to the frontend (via "/{basePath}/settings.js").
}

func (opts *Options) LogValue() slog.Value {
	var address string
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		// If we can't split, fallback to default
		address = "http://localhost" + addr
	}
	// Handle empty or wildcard host
	if host == "" || host == "0.0.0.0" || host == "127.0.0.1" {
		host = "localhost"
	}
	address = fmt.Sprintf("http://%s:%s", host, port)
	return slog.GroupValue(
		slog.String("address", address),
		slog.String("basePath", opts.BasePath),
		slog.String("rootDir", opts.RootDir),
		slog.Any("idp", opts.IDP),
	)
}

func (opts *Options) WithDefaults() *Options {
	if opts == nil {
		opts = &Options{}
	}
	if opts.BasePath == "" {
		opts.BasePath = "/"
	}
	if opts.RootDir == "" {
		opts.RootDir = "dist"
	}
	return opts
}

// Start starts the HTTP server with the given options.
func Start(mux *http.ServeMux, opts *Options) {
	if mux == nil {
		mux = http.DefaultServeMux
	}

	if opts == nil {
		opts = &Options{}
	}
	opts = opts.WithDefaults()

	err := routes(mux, opts)
	if err != nil {
		slog.Error("route configuration failed", "error", err)
		os.Exit(1)
	}

	slog.Info("starting server", "options", opts)
	err = http.ListenAndServe(addr, mux)
	if err != nil {
		slog.Error("server startup failed", "error", err)
		os.Exit(1)
	}
}

func routes(mux *http.ServeMux, opts *Options) error {
	// /isalive (without base path)
	mux.Handle("GET /isalive", httpx.Text("ALIVE"))
	// /isready (without base path)
	mux.Handle("GET /isready", httpx.Text("READY"))

	// /{basePath}/*
	baseMux := http.NewServeMux()

	// /{basePath}/auth/status
	baseMux.Handle("GET /auth/status", texas.Validate(opts.IDP))

	// /{basePath}/settings.js
	baseMux.Handle("GET /settings.js", settingsHandler(opts.BasePath, opts.EnvKeys))

	// /{basePath}/{prefix}/*
	err := proxy.Configure(baseMux, opts.Proxy, opts.IDP)
	if err != nil {
		return err
	}

	// /{basePath}/
	index, err := indexHandler(opts.RootDir, opts.DecoratorOpts)
	if err != nil {
		return err
	}
	if opts.IDP != nil {
		index = texas.Redirect(opts.IDP, index, opts.BasePath, opts.PublicPaths)
	}
	baseMux.Handle("/", rootServer(opts.RootDir, index))

	if opts.BasePath == "/" {
		mux.Handle(opts.BasePath, baseMux)
	} else {
		prefix := httpx.NormalizePath(opts.BasePath)
		mux.Handle(opts.BasePath, http.StripPrefix(prefix, baseMux))
	}

	return nil
}
