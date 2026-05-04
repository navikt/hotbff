package texas

import (
	"github.com/navikt/hotbff/internal/set"
)

type Options struct {
	// PublicExtensions are file extensions that don't require authentication (e.g. ".js", ".png", ".css").
	PublicExtensions set.Set[string]
	// PublicPaths are exact paths that don't require authentication (e.g. "/", "/about").
	PublicPaths set.Set[string]
	// PublicPrefixes are path prefixes that don't require authentication (e.g. "/assets/", "/public/").
	PublicPrefixes set.Set[string]
}

func (opts *Options) AddPublicExtensions(extensions ...string) {
	opts.PublicExtensions.AddAll(extensions...)
}

func (opts *Options) AddPublicPaths(paths ...string) {
	opts.PublicPaths.AddAll(paths...)
}

func (opts *Options) AddPublicPrefixes(prefixes ...string) {
	opts.PublicPrefixes.AddAll(prefixes...)
}

func (opts *Options) AddDefaults() {
	opts.PublicExtensions.AddAll(
		".css",
		".eot",
		".gif",
		".html",
		".ico",
		".jpeg",
		".jpg",
		".js",
		".js.map",
		".json",
		".otf",
		".png",
		".svg",
		".ttf",
		".txt",
		".webp",
		".woff",
		".woff2",
		".xml",
	)
	opts.PublicPrefixes.AddAll(
		"/assets/",
		"/public/",
	)
}

func (opts *Options) Clone() *Options {
	return &Options{
		PublicExtensions: opts.PublicExtensions.Clone(),
		PublicPaths:      opts.PublicPaths.Clone(),
		PublicPrefixes:   opts.PublicPrefixes.Clone(),
	}
}

func NewOptions() *Options {
	return &Options{
		PublicExtensions: set.New[string](),
		PublicPaths:      set.New[string](),
		PublicPrefixes:   set.New[string](),
	}
}

func DefaultOptions() *Options {
	opts := NewOptions()
	opts.AddDefaults()
	return opts
}
