# Deployment

## Local process

Build and run the Go service:

```sh
go build -o nats-auth-redirect .

AUTH_SERVER='http://127.0.0.1:8000' \
ACCOUNT_SIGNER_SEED='<account-signing-seed>' \
NATS_URL='nats://auth:<password>@127.0.0.1:4222' \
./nats-auth-redirect
```

For the repository demonstration:

```sh
docker compose up nats
python3 server.py
go run .
```

These commands require a correctly populated, private `.env` for values
expanded by `nats.conf` and used by the Go process. The compose file starts
only the NATS container. Start the redirect service and example HTTP server
separately.

## Container image

The Dockerfile:

- uses a Go 1.25.5 build stage;
- downloads modules before copying the remaining source;
- cross-compiles a static `runner` binary with CGO disabled;
- passes version, commit, and date through linker flags;
- copies the binary into a distroless Debian 12 image;
- runs as the image's non-root user;
- starts `/runner` as its entrypoint.

Build locally:

```sh
docker build -t nats-auth-redirect:local .
```

Run with secrets supplied at runtime:

```sh
docker run --rm \
  -e AUTH_SERVER='https://auth.example.internal/v1/nats/authorize' \
  -e ACCOUNT_SIGNER_SEED='<account-signing-seed>' \
  -e NATS_URL='tls://auth:<password>@nats.example.internal:4222' \
  nats-auth-redirect:local
```

The image has no declared port because the service initiates NATS and HTTP
connections and exposes no server endpoint.

## Multi-architecture publishing

`build-docker.sh` invokes Docker Buildx for `linux/amd64` and `linux/arm64`,
tags an image name fixed by the script, and pushes it. `TAG` selects the image
tag and defaults to `latest`.

Running this script mutates an external registry and requires:

- an active Buildx builder supporting both platforms;
- registry authentication;
- permission to push to the script's configured image repository;
- network access.

Review and change the destination image before publishing. The script was not
executed as part of documentation verification.

## Compose topology

`docker-compose.yml` creates:

- one `nats:latest` container;
- host mappings for ports 4222 and 8080;
- a read-only mount of `./nats.conf`;
- a named network, `os-msa-poc-network`;
- environment values loaded from `.env`.

Port 8080 is commented as a monitoring port in the compose file, but
`nats.conf` configures it as an unencrypted WebSocket listener. This conflict
is unresolved. The manifest does not define the redirect service, HTTP API,
health checks, restart policy, persistent storage, or TLS materials.

## Production assumptions and recommendations

The following are operational recommendations, not implemented guarantees:

- Pin container images instead of using `latest`.
- Store the signer seed and NATS credentials in a secret manager.
- Use TLS for both NATS and the HTTP API.
- Restrict network access so only the redirect service can reach the API and
  auth-callout NATS account.
- Redact tokens, bodies, request JWTs, and response JWTs from logs before
  production use.
- Add an explicit HTTP timeout and external liveness/readiness supervision.
- Validate the HTTP-provided account and permissions against an allow list.
- Run multiple redirect instances only after testing queue-group failover and
  request timeouts.
- Monitor process exits: signer parsing and initial NATS connection failures
  terminate the process.

## Build and release verification

Available repository checks:

```sh
go test ./...
go vet ./...
go build ./...
docker build -t nats-auth-redirect:local .
```

The Go repository currently contains no automated tests. A successful build
does not verify the NATS-to-HTTP authentication exchange. Before production,
add unit tests for HTTP status/body handling and claim construction, then add
an integration test with NATS auth callout enabled.

## Deployment gaps

- No Kubernetes, systemd, or complete compose deployment is supplied.
- No health/readiness endpoint or graceful shutdown behavior is implemented.
- No resource limits or capacity guidance are established.
- No HTTP timeout, retry policy, or maximum response size is configured.
- No automated end-to-end authentication test exists.
- Linker metadata targets in the Dockerfile have no corresponding variables in
  the indexed `main` package, so their intended observability is unclear.
