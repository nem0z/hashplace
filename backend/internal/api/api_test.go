package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nem0z/hashplace/backend/internal/api"
)

func TestHealthz(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	api.NewHandler().ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/healthz", nil))

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
	api.NewHandler().ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/healthz", nil))

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}
