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
	IDP           texas.IdentityProvider // Identity provider to use for token validation and exchange (if nil, no validation is performed).
	PublicPaths   []string               // Paths that should be publicly accessible without authentication (only relevant if IDP is set).
	EnvKeys       []string               // Environment variables that should be exposed to the frontend (via "/settings.js").
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
	)
}

// Start starts the HTTP server with the given [Options].
func Start(mux *http.ServeMux, opts *Options) {
	if mux == nil {
		mux = http.DefaultServeMux
	}
	if opts == nil {
		opts = &Options{}
	}

	err := Configure(mux, opts)
	if err != nil {
		slog.Error("server configuration failed", "error", err)
		os.Exit(1)
	}

	slog.Info("starting server", "options", opts)
	err = http.ListenAndServe(addr, mux)
	if err != nil {
		slog.Error("server startup failed", "error", err)
		os.Exit(1)
	}
}

// Configure configures the given [http.ServeMux] with the [Options] provided.
func Configure(mux *http.ServeMux, opts *Options) error {
	if mux == nil {
		mux = http.DefaultServeMux
	}
	if opts == nil {
		opts = &Options{}
	}

	basePath := opts.BasePath
	if basePath == "" {
		basePath = "/"
	}

	rootDir := opts.RootDir
	if rootDir == "" {
		rootDir = "dist"
	}

	// create index handler
	index, err := indexHandler(rootDir, opts.DecoratorOpts)
	if err != nil {
		return fmt.Errorf("failed to create index handler: %w", err)
	}
	index = protectedIndexHandler(basePath, opts.IDP, opts.PublicPaths, index)

	// / (public)
	mux.Handle("GET /isalive", httpx.Health("ALIVE"))
	mux.Handle("GET /isready", httpx.Health("READY"))

	// /base/path/ (public)
	baseMux := http.NewServeMux()
	baseMux.Handle("GET /settings.js", settingsHandler(basePath, opts.EnvKeys))
	baseMux.Handle("GET /auth/status", texas.Validate(opts.IDP))
	baseMux.Handle("/", newSPAHandler(rootDir, index))

	// /base/path/proxy/prexix
	err = proxy.Configure(baseMux, opts.Proxy, opts.IDP)
	if err != nil {
		return fmt.Errorf("failed to configure proxy: %w", err)
	}

	mux.Handle(basePath, httpx.MaybeStripPrefix(path.Join(basePath), baseMux))

	return nil
}

func protectedIndexHandler(basePath string, idp texas.TokenIntrospector, publicPaths []string, index http.Handler) http.Handler {
	if idp == nil {
		return index
	}
	protected := texas.Redirect(idp, index, basePath)
	if len(publicPaths) == 0 {
		return protected
	}
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if isPublicPath(req.URL.Path, publicPaths) {
			index.ServeHTTP(w, req)
			return
		}
		protected.ServeHTTP(w, req)
	})
}

func isPublicPath(requestPath string, publicPaths []string) bool {
	requested := httpx.NormalizePath(requestPath)
	for _, p := range publicPaths {
		public := httpx.NormalizePath(p)
		if public == "/" || requested == public || strings.HasPrefix(requested, public+"/") {
			return true
		}
	}
	return false
}
