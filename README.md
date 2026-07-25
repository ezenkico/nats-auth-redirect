# nats-auth-redirect

`nats-auth-redirect` is a NATS auth-callout service. It listens for NATS
authentication requests, forwards the client's connection token to an HTTP
endpoint as a Bearer token, and turns an accepted HTTP response into a signed
NATS authorization response.

This repository also contains a small example HTTP server and two example NATS
clients. They are development aids, not production components.

## How it works

The service:

1. connects to NATS and queue-subscribes to `$SYS.REQ.USER.AUTH`;
2. decodes the NATS authorization-request JWT;
3. sends `GET $AUTH_SERVER` with `Authorization: Bearer <connection-token>`;
4. decodes the HTTP JSON response into an account and optional permissions;
5. signs a user JWT and an authorization-response JWT with the account signer;
6. publishes the response JWT to the request's NATS reply subject.

See [Authentication flow](docs/authentication-flow.md) for the complete success
and failure paths.

## Requirements

- Go 1.23 or later, as declared by `go.mod`
- A NATS server configured with an auth callout
- An account-signing NKey seed matching the issuer configured in NATS
- An HTTP authorization endpoint

Docker is optional. The supplied Dockerfile builds with Go 1.25.5 and runs the
binary as a non-root user in a distroless image.

## Build and run

```sh
go build -o nats-auth-redirect .

export AUTH_SERVER=http://127.0.0.1:8000
export ACCOUNT_SIGNER_SEED='<account-signing-seed>'
export NATS_URL='nats://auth:<password>@127.0.0.1:4222'
./nats-auth-redirect
```

The program also loads a local `.env` file when present. Do not commit that
file. `AUTH_SERVER` is required in practice: when it is empty, the program
prints `No listener specified` and exits without starting a listener.

For a local demonstration, `python3 server.py` starts the example HTTP server
at `http://localhost:8000`. The supplied `docker-compose.yml` starts only NATS;
it does not start the redirect service or HTTP server. It also expects
`nats.conf` and `.env` to exist in the repository root.

See [Configuration](docs/configuration.md) and
[Deployment](docs/deployment.md) before operating the service.

## HTTP API contract

The service sends:

```http
GET /configured/path HTTP/1.1
Authorization: Bearer <NATS connection token>
```

The response type implemented by the Go service is:

```json
{
  "account": "APP",
  "pub": {
    "allow": ["events.>"],
    "deny": ["events.internal.>"]
  },
  "sub": {
    "allow": ["requests"],
    "deny": []
  }
}
```

`account` becomes the user JWT audience. `pub` and `sub` are optional
`jwt.Permission` values from `github.com/nats-io/jwt/v2`.

The included `server.py` uses a different, nested `perms` shape. Go's JSON
decoder ignores that unknown field, so the example currently demonstrates
account selection but not permission assignment. See
[Documentation gaps](docs/architecture.md#documentation-gaps).

## Development commands

```sh
go test ./...
go vet ./...
go build ./...
```

There are currently no Go test files.

The optional browser client has separate commands:

```sh
cd nats-ws-test
npm install
npm run lint
npm run build
npm run dev
```

## Security warning

The current implementation logs authentication request data, connection
tokens, HTTP response bodies, and generated JWTs. Treat its logs as sensitive.
Use TLS-protected NATS and HTTP connections in production, inject signer seeds
through a secret manager, and replace all example credentials.

Some repository examples contain credential- or seed-shaped development
values. They are not reproduced here and must not be treated as safe defaults.

## Documentation

- [Architecture](docs/architecture.md)
- [Configuration](docs/configuration.md)
- [Authentication flow](docs/authentication-flow.md)
- [Deployment](docs/deployment.md)
- [AI coding-agent guidance](AGENTS.md)

## License

See [LICENSE](LICENSE).
