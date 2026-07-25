# Configuration

The program calls `godotenv.Load()` at startup, then reads environment
variables with `os.Getenv`. Existing process environment values are retained
by the dotenv loader; a local `.env` supplies values that are otherwise absent.

## Service variables

| Variable | Required | Code default | Purpose |
|---|---:|---|---|
| `AUTH_SERVER` | Yes in practice | Empty | Full URL used for the outbound HTTP GET. An empty value causes the process to exit without listening. |
| `ACCOUNT_SIGNER_SEED` | Yes for production | An embedded development seed exists | Account-signing NKey seed used to sign user and authorization-response JWTs. The unsafe default is intentionally not reproduced. |
| `NATS_URL` | No | A localhost URL with embedded development credentials | NATS connection URL. The credential-bearing default is intentionally not reproduced. |

`ACCOUNT_SIGNER_PUB` appears in `env.example` and is expanded by `nats.conf` as
the auth-callout issuer. The Go process does not read it.

`NATS_SUBJECT` is read only by `publish-request.py`; it defaults to `requests`.
It does not configure the redirect service.

## Safe example

```dotenv
AUTH_SERVER=https://auth.example.internal/v1/nats/authorize
ACCOUNT_SIGNER_SEED=<account-signing-seed>
NATS_URL=tls://auth:<password>@nats.example.internal:4222
ACCOUNT_SIGNER_PUB=<matching-account-public-key>
```

Keep `.env` out of version control and restrict access to the signer seed. The
repository `.gitignore` excludes `.env` and `nats-keys`, but operators must
also protect environment inspection, process logs, CI output, and container
configuration.

## NATS development configuration

The supplied `nats.conf` demonstrates:

- an auth callout whose issuer is `$ACCOUNT_SIGNER_PUB`;
- `AUTH` as the callout account;
- `auth` as an exempt auth user used by the redirect service;
- `APP` and `SYS` accounts;
- `SYS` as the system account;
- an unencrypted WebSocket listener on port 8080.

The file contains fixed development users/passwords. Replace them before any
non-local deployment. The service's `NATS_URL` credentials must correspond to
an auth user that can connect and receive auth-callout requests.

## HTTP response configuration

The Go decoder accepts this JSON structure:

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

`pub` and `sub` are optional. Unknown JSON fields are ignored. Missing or empty
`account` is not rejected by the redirect service before claim encoding; the
resulting NATS behavior is not covered by tests.

## Fixed runtime settings

The following values are not configurable in the current implementation:

- subscription subject: `$SYS.REQ.USER.AUTH`;
- queue group: `auth-workers`;
- NATS connection name: `auth-validator`;
- reconnect attempts: unlimited;
- reconnect wait: 500 milliseconds;
- HTTP method: GET;
- HTTP client: `http.DefaultClient`;
- HTTP timeout: none explicitly configured.

## Unresolved configuration questions

- Should `AUTH_SERVER` be validated at startup rather than during a request?
- Should signer and NATS URL defaults be removed?
- Should queue name, connection name, HTTP timeout, and status policy become
  configurable?
- Should the HTTP response use top-level permissions or the nested shape in
  `server.py`?
