package handlers

import (
	"fmt"
	"net/http"
)

func Health() {
	http.HandleFunc("GET /isalive", func(w http.ResponseWriter, _ *http.Request) {
		_, err := fmt.Fprint(w, "ALIVE")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	http.HandleFunc("GET /isready", func(w http.ResponseWriter, _ *http.Request) {
		_, err := fmt.Fprint(w, "READY")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
