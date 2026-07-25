# AGENTS.md

This file guides AI coding agents working in this repository.

## Scope and source of truth

The Go service is the product. Python programs and `nats-ws-test/` are
demonstration utilities. Ground behavioral claims and changes in this order:

1. Go implementation and its types;
2. tests (none currently exist);
3. `nats.conf`, `Dockerfile`, and deployment manifests;
4. scripts and example clients.

Do not infer behavior that these sources do not establish. Mark conflicts and
unknowns explicitly.

## Repository map

- `main.go`: loads `.env`, selects the HTTP redirect, and starts the service.
- `base/natsAuth.go`: NATS connection, auth-callout decoding, claim creation,
  signing, error responses, and reply publication.
- `redirects/http.go`: HTTP Bearer-token forwarding and response decoding.
- `nats.conf`: development auth-callout and WebSocket configuration.
- `Dockerfile`, `build-docker.sh`, `docker-compose.yml`: container artifacts.
- `server.py`: example HTTP authorization endpoint.
- `publish-request.py`: example Python publisher.
- `nats-ws-test/`: example React/NATS WebSocket subscriber.

## Investigation workflow

Use repository maps and symbol-level retrieval before reading whole files.
Start with these symbols:

- `main.main`
- `base.GetConnectionDataEnv`
- `base.Listen`
- `base.badError`
- `redirects.SetupHttpAuth`

Trace both the success path and every early return. Inspect configuration and
examples only after confirming the implementing symbol.

## Build and verification

For Go changes:

```sh
go test ./...
go vet ./...
go build ./...
```

The current test command succeeds with `[no test files]`; it does not establish
behavioral coverage. Add focused tests for behavior changes before modifying
authentication, HTTP status handling, claim construction, or error responses.

For browser-client changes, from `nats-ws-test/`:

```sh
npm install
npm run lint
npm run build
```

Do not claim Docker, NATS, or end-to-end behavior was verified unless those
services were actually run and the observed result is reported.

## Change rules

- Preserve the NATS request/reply subject contract and JWT signing chain.
- Keep HTTP response fields aligned with `base.ResponseData`.
- Document new environment variables in `docs/configuration.md`.
- Document flow or failure changes in `docs/authentication-flow.md`.
- Keep production recommendations distinct from implemented behavior.
- Do not silently “correct” example/configuration conflicts; explain or test
  them first.
- Do not edit generated NKey material.

## Secrets and sensitive data

Never print, copy into documentation, or commit:

- account signer seeds or private NKeys;
- connection tokens, passwords, or authorization JWTs;
- `.env` contents;
- generated files under `nats-keys/`.

Use obvious placeholders such as `<account-signing-seed>`. The implementation
currently logs tokens, response bodies, request data, and signed JWTs; avoid
repeating real log output in issues, tests, or documentation.

## Known gaps

- There are no automated tests.
- The Python example's nested `perms` response does not match the Go
  `ResponseData` JSON shape.
- HTTP responses are rejected only when their status is greater than 300.
- HTTP requests use `http.DefaultClient` without an explicit timeout.
- Some errors return signed authorization responses; other decode/encode
  failures return plain JSON.
- The compose file starts NATS only and labels port 8080 as monitoring even
  though `nats.conf` configures it for WebSockets.
- Docker linker flags target version metadata variables not present in the
  indexed `main` package.

Treat these as unresolved until code or tests establish intended behavior.
