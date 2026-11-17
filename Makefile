# Makefile: one entry point for every task in this project.
#
# Why a Makefile: the project mixes Go, Rust, Docker and shell scripts. Make
# gives one memorable command per task (`make test`) regardless of language,
# and ships with macOS, so nothing extra to install.
#
# Recipe lines MUST start with a TAB character, not spaces (a Make rule).
#
# .PHONY lists targets that are names of actions, not files on disk. Why:
# without it, a file named `test` in this folder would stop `make test` running.
.PHONY: test test-go test-rust lint fmt up down logs

# test: run the unit tests for both languages. Fast, no Docker needed.
test: test-go test-rust

# test-go: run every Go test under server/. `./...` means "this folder and all subfolders".
test-go:
	cd server && go test ./...

# test-rust: run every Rust test (unit tests and doc tests) in the workspace.
test-rust:
	cd authenticator && cargo test --workspace

# lint: check code for likely bugs without running it.
#   go vet      = Go's built-in bug finder
#   cargo clippy = Rust's linter; -D warnings turns every warning into an error
lint:
	cd server && go vet ./...
	cd authenticator && cargo clippy --workspace --all-targets -- -D warnings

# fmt: rewrite files into the standard layout so style is never debated.
fmt:
	cd server && gofmt -w .
	cd authenticator && cargo fmt --all

# up: build the image and start every service in the background (-d).
#   --build  rebuild the image if the code changed
#   --wait   don't return until each service's healthcheck reports "healthy",
#            so `make up && curl ...` never races a still-starting server
up:
	docker compose up -d --build --wait

# down: stop and remove the containers (images and your files are kept).
down:
	docker compose down

# logs: follow the services' output; press Ctrl+C to stop watching.
logs:
	docker compose logs -f
