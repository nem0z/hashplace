// Package api implements the HTTP transport of the hashplace server.
package api

import (
	"encoding/json"
	"net/http"
)

// NewHandler returns the HTTP handler serving every API route.
func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handleHealthz)

	return mux
}

func handleHealthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
