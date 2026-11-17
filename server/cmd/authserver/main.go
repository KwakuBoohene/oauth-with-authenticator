// Package main is the entry point for the authserver binary: the OIDC provider
// and WebAuthn Relying Party (see README.md, plans/epics/E00-foundation.md).
//
// E00-T03 only adds the bare minimum a container needs to prove it is alive:
//   - GET /healthz answers "ok", and
//   - the -healthcheck flag makes the binary call its OWN /healthz and exit 0/1.
//
// Why the self-check flag: the final Docker image is "distroless", meaning it has
// no shell and no curl. Docker's HEALTHCHECK therefore has to run our own binary.
// Real routes, config and graceful shutdown arrive in later tickets (E05-T01).
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

// newMux builds the HTTP routes. It is a function (not a global) so tests can
// create a fresh router without starting a real server.
func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthz)
	return mux
}

// healthz answers "ok" so Docker (and humans with curl) can tell the process is serving.
func healthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "ok")
}

// healthcheckURL turns a listen address such as ":8080" into the URL the
// self-check should call. An empty host (":8080") means "all interfaces", so the
// check talks to loopback (127.0.0.1) instead.
func healthcheckURL(addr string) (string, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "", fmt.Errorf("bad address %q: %w", addr, err)
	}
	if host == "" || host == "0.0.0.0" {
		host = "127.0.0.1"
	}
	return "http://" + net.JoinHostPort(host, port) + "/healthz", nil
}

// runHealthcheck calls our own /healthz once and returns nil only on a 200.
// The short timeout matters: Docker would otherwise wait on a hung server.
func runHealthcheck(addr string) error {
	url, err := healthcheckURL(addr)
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("healthz returned %s", resp.Status)
	}
	return nil
}

func main() {
	addr := flag.String("addr", ":8080", "address to listen on (or check, with -healthcheck)")
	check := flag.Bool("healthcheck", false, "call this server's /healthz and exit 0 (healthy) or 1")
	flag.Parse()

	if *check {
		if err := runHealthcheck(*addr); err != nil {
			fmt.Fprintln(os.Stderr, "unhealthy:", err)
			os.Exit(1)
		}
		return
	}

	// ReadHeaderTimeout stops "slowloris" clients that open a connection and
	// dribble header bytes forever to tie up the server.
	srv := &http.Server{Addr: *addr, Handler: newMux(), ReadHeaderTimeout: 5 * time.Second}
	log.Printf("authserver listening on %s", *addr)
	log.Fatal(srv.ListenAndServe())
}
