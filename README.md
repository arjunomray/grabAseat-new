# GrabMySeat

An event ticketing platform where people can create events, browse them, and claim seats. Built as a small, end-to-end full-stack project: a Go/Gin API, a React/Vite single-page app, and PostgreSQL.

## Tech stack

| Layer | Technology |
|---|---|
| Backend | Go 1.23, Gin, GORM, PostgreSQL 16 |
| Auth | JWT bearer tokens |
| Frontend | React 18, Vite 6, React Router, Tailwind CSS, `qrcode.react` |
| Tooling | Nix flakes, GNU Make, goose (migrations), golangci-lint, govulncheck |
| Containers (optional) | Docker / Docker Compose |

## Architecture

```
┌──────────────────────────┐
│  React SPA (Vite :5173)  │
│  fetch + Bearer JWT      │
└────────────┬─────────────┘
             │ JSON / HTTP
             ▼
┌──────────────────────────┐
│  Go API (Gin :8080)      │
│  handlers → GORM         │
└────────────┬─────────────┘
             │ SQL
             ▼
┌──────────────────────────┐
│  PostgreSQL 16 (:5432)   │
└──────────────────────────┘
```

In development the SPA is served by the Vite dev server and the API runs as a normal Go process. `make dev` starts PostgreSQL, the API, and the UI together.

## Repository layout

```
.
├── main.go               # API entrypoint and routes
├── handlers/             # HTTP handlers (users, events, tickets)
├── middleware/           # JWT auth middleware
├── models/               # GORM models
├── dbConfig/             # Database connection
├── utils/                # JWT helpers
├── frontend/             # React + Vite app
├── Makefile              # Local dev and database lifecycle
├── flake.nix             # Reproducible dev shell and builds
└── docker-compose.yml    # Optional container workflow
```

## Getting started with Nix (recommended)

Requires [Nix](https://nixos.org) with flakes enabled.

```bash
# 1. Enter the development shell (Go, Node, PostgreSQL, goose, make, ...)
nix develop

# 2. Initialise and start the local PostgreSQL cluster, then run everything
make dev
```

- API: <http://localhost:8080>
- UI: <http://localhost:5173>

`make dev` creates a project-local database cluster in `.postgres/` (trust auth, local only), starts it, and then runs the API and the Vite server side by side. Stop the database with `make db-stop`.

You can also start everything in one shot, without entering the shell first:

```bash
nix run   # runs `make dev` with the toolchain bundled
```

### Make targets

| Target | Description |
|---|---|
| `make dev` | PostgreSQL + API + UI together |
| `make backend` | Run the API only |
| `make frontend` | Run the Vite dev server only |
| `make db-init` / `db-start` | Initialise and start the local PostgreSQL cluster |
| `make db-stop` / `db-reset` | Stop it / stop and wipe it |
| `make migrate` | Apply goose migrations (once `migrations/` exists) |
| `make test` | Run the Go tests |
| `make build` | Build `bin/eventer` and the Vite production bundle |
| `make help` | List all targets |

## Docker (optional)

A `docker-compose.yml` is provided that builds production-style images: the backend applies migrations on startup and runs as a non-root user, and the frontend is served by nginx (which also proxies `/api` to the backend).

```bash
JWT_SECRET=change-me docker compose up --build
```

- UI: <http://localhost:5173>
- API: <http://localhost:8080>

Note: the Docker path is provided as an alternative to Nix and is not exercised by CI.

## Environment variables

The service loads an optional `.env` file at startup, then reads these variables from the process environment. Copy `.env.example` to `.env` and adjust; the Makefile supplies the same defaults, so no `.env` file is required for `make dev`.

| Variable | Default | Description |
|---|---|---|
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `postgres` | Database user |
| `DB_PASSWORD` | `postgres` | Database password |
| `DB_NAME` | `eventer` | Database name |
| `PORT` | `8080` | API port |
| `JWT_SECRET` | – | Required; secret used to sign JWTs |
| `CORS_ORIGINS` | `http://localhost:5173` | Comma-separated allowed origins |

The frontend reads `VITE_API_URL` (see `frontend/.env.example`). Leave it unset in development and the Vite dev server proxies `/api` to the backend; set it when the UI and API are served from different origins.

## API

All routes are JSON over HTTP. Protected routes require `Authorization: Bearer <token>`.

| Method | Path | Auth | Description |
|---|---|---|---|
| `GET` | `/healthz` | – | Liveness probe: the process is running |
| `GET` | `/readyz` | – | Readiness probe: the database is reachable |
| `GET` | `/api/v1/stream` | – | Server-Sent Events stream of seat-availability updates |
| `POST` | `/api/v1/register` | – | Create an account |
| `POST` | `/api/v1/login` | – | Exchange credentials for a JWT |
| `POST` | `/api/v1/events` | ✓ | Create an event |
| `GET` | `/api/v1/events` | ✓ | List events |
| `GET` | `/api/v1/events/:id` | ✓ | Get one event |
| `POST` | `/api/v1/tickets` | ✓ | Claim a ticket (supports an `Idempotency-Key` header for safe retries) |
| `GET` | `/api/v1/tickets` | ✓ | List the caller's tickets |
| `GET` | `/api/v1/tickets/me` | ✓ | Alias for the caller's tickets |
| `POST` | `/api/v1/tickets/verify` | ✓ | Mark a ticket as used (QR check-in) |
| `DELETE` | `/api/v1/tickets/:id` | ✓ | Cancel one of your tickets |
| `GET` | `/api/v1/users/:id` | ✓ | User profile |
| `GET` | `/api/v1/users/:id/events` | ✓ | Events a user created or registered for |
| `GET` | `/api/v1/users/:id/tickets` | ✓ | Tickets for a user |

> API routes are versioned under `/api/v1`; `/healthz` and `/readyz` stay unversioned.

## Data model

- **User** — `id` (UUID), `name`, `email` (unique), `password`
- **Event** — `id`, `created_by`, `name`, `description`, `time` (`timestamptz`), `location`, `price_cents` (integer minor units), `seats_available`, `tags`
- **Ticket** — `id`, `event_id`, `user_id` (a registration connecting a user to an event)

## Testing

```bash
make test        # starts Postgres, migrates a dedicated eventer_test DB, runs go test ./...
nix flake check  # builds the backend and frontend
```

Handler tests run against a real PostgreSQL database (`eventer_test`), so `make test` is the supported entry point — it creates and migrates that database first.

## Design notes and tradeoffs

- **Single service, single SPA.** No microservices; the API and UI are simple to run and reason about, at the cost of one deployable handling everything.
- **GORM over hand-written SQL.** Faster to build CRUD, but schema management is moving from `AutoMigrate` to versioned goose migrations so changes are reviewable and reversible.
- **JWT in `localStorage`.** Simple and stateless; it does mean XSS is a real threat and there is no server-side revocation or refresh yet.
- **Money as `float64`, time as `string`.** Convenient initially, but both will become integer minor units and `timestamptz` respectively.
- **Dev-first tooling.** The Nix flake and Makefile are the primary workflow; production container images are a later step.

## License

Released under the [MIT License](LICENSE).
