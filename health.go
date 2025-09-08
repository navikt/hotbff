package hotbff

import (
	"fmt"
	"net/http"
)

func Health() {
	http.HandleFunc("GET /isalive", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprint(w, "ALIVE")
	})
	http.HandleFunc("GET /isready", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprint(w, "READY")
	})
}
