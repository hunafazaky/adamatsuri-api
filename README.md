# AdaMatsuri API

A REST API for browsing events and managing bookings, built with Go, Gin,
and GORM.

## Stack

- **Go** / **Gin** — HTTP layer
- **GORM** / **PostgreSQL** — persistence
- **JWT** — authentication
- **ImageKit** — event image storage
- **swaggo** — OpenAPI documentation

## Architecture

Layered: `handler → service → repository`, each depending on the layer
below through an interface. See [`docs/API_RESPONSES.md`](docs/API_RESPONSES.md)
for the full response contract.

```
cmd/server        entrypoint — config, DB, dependency wiring
internal/
  handler          HTTP layer: bind request, call service, write response
  service          business rules, validation, DTO mapping
  repository       GORM queries, no business logic
  model            database entities
  dto              API request/response shapes
  middleware       auth
  router           route registration
  config           typed env config
  apperror         typed errors → HTTP status mapping
  response         standard JSON envelope
docs               API reference + generated OpenAPI spec
```

## Getting started

**Prerequisites:** Go 1.22+, Docker, an [ImageKit](https://imagekit.io)
account (for event image uploads).

1. Copy the env template and fill in your values:
   ```bash
   cp .env.example .env
   ```
2. Start the API and database:
   ```bash
   docker compose up
   ```
   The API is available at `http://localhost:8080` (or whatever `PORT`
   you set).

## API documentation

Interactive docs (Scalar UI, with request/response schemas and
try-it-out): **`/docs`**

For a quick-reference summary without running the server, see
[`docs/API_RESPONSES.md`](docs/API_RESPONSES.md).

## Development

Regenerate the OpenAPI spec after changing any handler's `@swag` annotations:
```bash
swag init -g cmd/server/main.go -o docs
```
