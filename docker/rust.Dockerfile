# rust.Dockerfile: runs the Rust checks (clippy and tests) inside a clean container.
#
# Why a container: a test that passes on your Mac may secretly rely on something
# installed there (a library, an environment variable, a file). A fresh image has
# only what this file installs, so a pass here is real proof.
#
# Why the checks run DURING `docker build`: a RUN step that fails stops the build.
# So a successful `docker build --target test` means every check passed, and the
# Makefile needs no extra "did it pass?" logic.

# ---- Stage: test ----
# Pinned to a Rust version so a new release can't change results without us
# noticing. The edition in authenticator/Cargo.toml (2024) needs Rust 1.85+.
FROM rust:1.98-bookworm AS test
WORKDIR /src

# The official rust image ships without clippy (the linter) to keep it small.
# Why add it here: `make lint` uses clippy, so the container must check the same way.
RUN rustup component add clippy

# Copy the workspace. Only authenticator/ is needed: the Go server is tested
# by its own container step (see Makefile `test-docker`).
COPY authenticator/ ./

# Cache mounts keep the downloaded crates (registry) and the compiled build
# output (target) between builds, so re-runs are fast. They live on the build
# machine, not inside the image. Why not copy target/ from your Mac: it was built
# for macOS, and the container is Linux, so it would be useless here.
#
# -D warnings: clippy warnings count as failures, same as `make lint`.
RUN --mount=type=cache,target=/usr/local/cargo/registry \
    --mount=type=cache,target=/src/target \
    cargo clippy --workspace --all-targets -- -D warnings \
 && cargo test --workspace
