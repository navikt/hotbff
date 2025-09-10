package hotbff

import (
	"fmt"
	"net/http"
)

func statusHandler(status string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, err := fmt.Fprint(w, status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
