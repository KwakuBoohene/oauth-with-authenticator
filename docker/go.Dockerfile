# go.Dockerfile: builds every Go binary in server/cmd/ and ships them in a tiny image.
#
# Why multi-stage: stage 1 ("build") has the whole Go toolchain (~800 MB) so it
# can compile. Stage 2 copies ONLY the finished binaries into a near-empty image.
# Anything not copied (compiler, source code, shell) is not in the final image.
# That is smaller to ship and leaves an attacker almost nothing to use.

# ---- Stage 1: compile ----
FROM golang:1.26 AS build
WORKDIR /src

# go.mod first, as its own layer: Docker caches layers, so this only re-runs
# when go.mod changes (matters later, once there are dependencies to download).
COPY server/go.mod ./
COPY server/ ./

# CGO_ENABLED=0 -> a fully static binary with no C library, which is what lets
# it run on "distroless/static" (it has no libc). -trimpath removes your local
# folder paths from the binary; -s -w drop debug info to shrink it.
# `./cmd/...` builds every program under cmd/ into /out/.
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/ ./cmd/...

# ---- Stage 2: run ----
# "distroless" = no shell, no package manager, no curl. "nonroot" = runs as an
# unprivileged user (UID 65532) instead of root.
FROM gcr.io/distroless/static:nonroot
COPY --from=build /out/ /app/

# Said explicitly (even though the image already defaults to it) so a reader
# sees that the process never runs as root.
USER nonroot:nonroot

# CMD, not ENTRYPOINT: later services (app-a, app-b, ...) reuse this image and
# simply override the command in docker-compose.yml.
CMD ["/app/authserver"]
