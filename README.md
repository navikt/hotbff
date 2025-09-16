# hotbff

Frackend (backend for frontend) for DigiHoTs apper.

Støtter følgende funksjoner:
1. Server for statiske filer (HTML, JavaScript, bilder etc.)
2. Client side routing
    * Svarer med index.html hvis request-URL ikke peker på en faktisk fil (og aldri 404)
3. Token validering med Texas
    * Redirect til login hvis manglende eller ugyldig token
4. Reverse proxy
    * Token exchange med Texas hvis backend krever dette
5. Dekorering av templates med Dekoratøren
6. CSP (TODO)

## Eksempel på bruk:
```go
package main

import (
	"os"

	"github.com/navikt/hotbff"
	"github.com/navikt/hotbff/proxy"
	"github.com/navikt/hotbff/texas"
)

func main() {
	opts := &hotbff.Options{
		BasePath: "/",
		RootDir:  "dist",
		Proxy: &proxy.Map{
			"/api/": &proxy.Options{
				Target:      os.Getenv("API_URL"),
				StripPrefix: false,
                // Identity provider for token exchange.
				IDP:         texas.TokenX,
				IDPTarget:   os.Getenv("API_SCOPE"),
			},
			"/public-api/": &proxy.Options{
				Target:      os.Getenv("PUBLIC_API_URL"),
                // true hvis kall som f.eks. /public-api/api skal skrives om til bare /api mot backend.
				StripPrefix: true,
			},
		},
        // Identity provider for validering av token.
		IDP: texas.IDPorten,
        // Disse blir tilgengelige i window.appSettings hvis f.eks. index.html laster inn /{BasePath}/settings.js.
		EnvKeys: &[]string{
            "SOME_ENV_A",
            "SOME_ENV_B",
            "SOME_ENV_C",
		},
	}
	hotbff.Start(opts)
}
```
