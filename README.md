# mon-go

MongoDB + Go HTTP API. Uses the [official MongoDB Go driver](https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo) and [Chi](https://github.com/go-chi/chi) for routing.

## Layout

- **`cmd/server`** – main entry; wires config, DB, handlers, and server
- **`configs/`** – config file folder (`configs/config.yaml`)
- **`internal/config`** – config loader (Viper: `AddConfigPath` + `SetConfigName`)
- **`internal/model`** – domain types (e.g. `Item`)
- **`internal/store`** – MongoDB connection and repository-style access (e.g. `ItemStore`)
- **`internal/handler`** – HTTP handlers (JSON in/out)
- **`internal/server`** – Chi router and HTTP server setup

Adding a new resource: add a model, a store in `internal/store`, a handler in `internal/handler`, and register routes in `internal/server/server.go`.

## Run

### Local (Go only)

1. Start MongoDB (e.g. `docker run -d -p 27017:27017 mongo`).
2. From repo root:

   ```bash
   make run
   # or: go run ./cmd/server
   ```

### Docker Compose (app + Mongo)

Services use the **`mon-go`** profile. Start with:

```bash
make run-docker   # build images and start containers
# or: make build-docker && make up
# or: docker compose --profile mon-go up -d
```

App: http://localhost:8080 — Mongo: localhost:27017. Stop with `make down`.

To run **only MongoDB** (e.g. for local development with `make run`): `make up-mongo`; stop with `make down-mongo`.

### Config

- **Path (hardcoded):** Config is read from the `configs/` folder via `AddConfigPath` and `SetConfigName` (file: `configs/config.yaml`). No env overrides.
- **Local:** Run from repo root so `configs/config.yaml` is found.
- **Docker:** Compose injects config via the `configs` key; the config is mounted at `/app/configs/config.yaml` in the image.

## API (example: items)

- `GET /health` – health check
- `GET /ping` – returns "Hello World"
- `POST/GET /items.create` – create item (POST body: `{"name": "..."}`; GET query: `?name=...`)
- `POST/GET /items.list` – list items
- `POST/GET /items.get` – get one item (query: `?id=...`)
- `POST/GET /items.delete` – delete item (query: `?id=...`)

## Makefile

| Command        | Description                    |
|----------------|--------------------------------|
| `make build`   | Build binary to `_build/server`   |
| `make tidy`    | Run `go mod tidy`              |
| `make vet`     | Run `go vet ./...`             |
| `make run`     | Run server locally             |
| `make build-docker` | Build Docker images        |
| `make up`      | Start app + Mongo (mon-go profile, detached) |
| `make down`    | Stop and remove containers (mon-go profile) |
| `make up-mongo`   | Start only MongoDB (mongo profile)     |
| `make down-mongo` | Stop and remove MongoDB (mongo profile) |
| `make run-docker` | Build images and start containers (mon-go profile) |
