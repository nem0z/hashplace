// Package api implements the HTTP transport of the hashplace server.
package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

// NewHandler returns the HTTP handler serving every API route.
func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handleHealthz)

	return withServerClock(mux)
}

// withServerClock sets the Hashplace-Now header (Unix seconds) on every response.
func withServerClock(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Hashplace-Now", strconv.FormatInt(time.Now().Unix(), 10))
		next.ServeHTTP(w, r)
	})
}

func handleHealthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
