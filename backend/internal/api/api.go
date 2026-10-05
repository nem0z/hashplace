// Package api implements the HTTP transport of the hashplace server.
package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

// NewHandler returns the HTTP handler serving every API route.
// now is the server clock, sent with every response.
func NewHandler(now func() time.Time) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handleHealthz)

	return withServerClock(mux, now)
}

// withServerClock sets the Hashplace-Now header (Unix seconds) on every response.
func withServerClock(next http.Handler, now func() time.Time) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Hashplace-Now", strconv.FormatInt(now().Unix(), 10))
		next.ServeHTTP(w, r)
	})
}

func handleHealthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
