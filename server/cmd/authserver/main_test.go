package main

import (
	"strings"
	"testing"
)

// TestBanner checks that the placeholder banner mentions the binary name.
// Why: this is the one trivial test that proves `go test` runs and is wired
// into `make test-go`. Later tickets replace it with real tests.
func TestBanner(t *testing.T) {
	got := banner()
	if !strings.HasPrefix(got, "authserver") {
		t.Fatalf("banner() = %q, want prefix %q", got, "authserver")
	}
}
