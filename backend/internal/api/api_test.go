package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nem0z/hashplace/backend/internal/api"
)

func fixedNow() time.Time { return time.Unix(1791072000, 0) }

func TestHealthz(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	api.NewHandler(fixedNow).ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/healthz", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}

	if got, want := rec.Body.String(), `{"status":"ok"}`+"\n"; got != want {
		t.Errorf("body = %q, want %q", got, want)
	}
}

func TestHealthzRejectsOtherMethods(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	api.NewHandler(fixedNow).ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/healthz", nil))

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestServerClockHeader(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{name: "ok", method: http.MethodGet, path: "/healthz"},
		{name: "method not allowed", method: http.MethodPost, path: "/healthz"},
		{name: "not found", method: http.MethodGet, path: "/missing"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			api.NewHandler(fixedNow).ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), tt.method, tt.path, nil))

			if got := rec.Header().Get("Hashplace-Now"); got != "1791072000" {
				t.Errorf("Hashplace-Now = %q, want 1791072000", got)
			}
		})
	}
}
