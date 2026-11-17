package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHealthzReturnsOK calls the handler through the real router using httptest,
// which runs an HTTP handler in memory without opening a network port.
func TestHealthzReturnsOK(t *testing.T) {
	rec := httptest.NewRecorder()
	newMux().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := rec.Body.String(); got != "ok\n" {
		t.Fatalf("body = %q, want %q", got, "ok\n")
	}
}

// TestHealthzRejectsOtherMethods checks that only GET is routed. Why: a health
// endpoint should not quietly accept writes. Go 1.22+ patterns like "GET /healthz"
// answer other methods with 405 Method Not Allowed.
func TestHealthzRejectsOtherMethods(t *testing.T) {
	rec := httptest.NewRecorder()
	newMux().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/healthz", nil))

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

// TestHealthcheckURL covers how a listen address becomes the self-check URL.
func TestHealthcheckURL(t *testing.T) {
	cases := []struct{ addr, want string }{
		{":8080", "http://127.0.0.1:8080/healthz"},
		{"0.0.0.0:9000", "http://127.0.0.1:9000/healthz"},
		{"127.0.0.1:1234", "http://127.0.0.1:1234/healthz"},
	}
	for _, c := range cases {
		got, err := healthcheckURL(c.addr)
		if err != nil || got != c.want {
			t.Errorf("healthcheckURL(%q) = %q, %v; want %q", c.addr, got, err, c.want)
		}
	}
	if _, err := healthcheckURL("not-an-address"); err == nil {
		t.Error("healthcheckURL(\"not-an-address\") returned no error")
	}
}

// TestRunHealthcheck runs the real self-check against a live in-process server:
// healthy server -> nil error, and a server answering 500 -> error.
func TestRunHealthcheck(t *testing.T) {
	healthy := httptest.NewServer(newMux())
	defer healthy.Close()
	if err := runHealthcheck(healthy.Listener.Addr().String()); err != nil {
		t.Fatalf("healthy server reported error: %v", err)
	}

	broken := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "boom")
	}))
	defer broken.Close()
	if err := runHealthcheck(broken.Listener.Addr().String()); err == nil {
		t.Fatal("broken server reported healthy")
	}
}
