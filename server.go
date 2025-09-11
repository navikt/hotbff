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

type Options struct {
	BasePath      string
	RootDir       string
	DecoratorOpts *decorator.Options
	Proxy         *proxy.Map
	IDP           texas.IdentityProvider
	EnvKeys       *[]string
}

func Start(opts *Options) {
	slog.Info("hotbff: starting server", "address", addr, "basePath", opts.BasePath, "rootDir", opts.RootDir)
	err := http.ListenAndServe(addr, Handler(opts))
	if err != nil {
		slog.Error("hotbff: server startup failed", "error", err)
		os.Exit(1)
	}
}

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
	root := http.NewServeMux()
	root.Handle("GET /isalive", healthHandler("ALIVE"))
	root.Handle("GET /isready", healthHandler("READY"))

	// /base/path/ (public)
	base := http.NewServeMux()
	base.Handle("GET /settings.js", settingsHandler(opts.EnvKeys))

	// /base/path/ (protected)
	protected := http.NewServeMux()
	protected.Handle("/", staticHandler(rootDir, opts.DecoratorOpts))

	// /base/path/proxy/prefix/ (protected)
	if opts.Proxy != nil {
		for prefix, proxyOpts := range *opts.Proxy {
			slog.Info("hotbff: adding proxy", "prefix", prefix, "target", proxyOpts.Target)
			proxyHandler := proxyOpts.Handler()
			if proxyOpts.StripPrefix {
				protected.Handle(prefix, http.StripPrefix(prefix, proxyHandler))
			} else {
				protected.Handle(prefix, proxyHandler)
			}
		}
	}

	base.Handle("/", texas.Protected(opts.IDP, basePath, protected))

	root.Handle(basePath, maybeStripPrefix(path.Join(basePath), base))
	return root
}

func maybeStripPrefix(prefix string, h http.Handler) http.Handler {
	if prefix == "/" {
		return h
	}
	return http.StripPrefix(prefix, h)
}
