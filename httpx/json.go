package httpx

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, v any) {
	w.Header().Set(HeaderContentType, ContentTypeApplicationJSON)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ServeJSON(v any) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		JSON(w, v)
	})
}
