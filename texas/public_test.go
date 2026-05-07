package texas

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/navikt/hotbff/internal/assert"
)

func Test_isPublic(t *testing.T) {
	opts := DefaultOptions()
	opts.AddPublicPaths("/about")

	type args struct {
		urlPath  string
		basePath string
		opts     *Options
	}
	tests := []struct {
		name   string
		args   args
		public bool
		reason string
	}{
		{
			name: "extension match",
			args: args{
				urlPath:  "/image.png",
				basePath: "/",
				opts:     opts,
			},
			public: true,
			reason: "extension match: \".png\"",
		},
		{
			name: "path match",
			args: args{
				urlPath:  "/about",
				basePath: "/",
				opts:     opts,
			},
			public: true,
			reason: "exact path match: \"/about\"",
		},
		{
			name: "prefix match",
			args: args{
				urlPath:  "/assets/image",
				basePath: "/",
				opts:     opts,
			},
			public: true,
			reason: "prefix match: \"/assets/\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.args.urlPath, nil)
			public, reason := isPublic(req, tt.args.basePath, tt.args.opts)
			assert.Equal(t, public, tt.public)
			assert.Equal(t, reason, tt.reason)
		})
	}
}
