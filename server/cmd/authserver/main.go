// Package main is the entry point for the authserver binary: the OIDC provider
// and WebAuthn Relying Party (see README.md, plans/epics/E00-foundation.md).
//
// Why a placeholder: E00-T01 only proves the Go toolchain and test runner work.
// Real behaviour (HTTP routes, passkeys, tokens) arrives in later tickets.
package main

import "fmt"

// banner returns the startup line. It is a separate function so the test can
// check it without starting a server.
func banner() string {
	return "authserver: not implemented yet (see plans/README.md)"
}

// main prints the banner and exits. Real startup code replaces this in E05-T01.
func main() {
	fmt.Println(banner())
}
