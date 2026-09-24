# Anime Events API

A REST API for browsing and booking anime-community events —
conventions, doujin markets, screenings, cosplay contests, game
tournaments, and meetups. Built with Go, Gin, and GORM.

## Features

- **Auth** — JWT-based sign-up/sign-in, with a self-service role
  choice (`attendee` or `organizer`) at sign-up
- **Events** — categorized (convention, doujin market, screening,
  cosplay contest, game tournament, meetup), with free-form
  fandom/genre tags (e.g. "Jujutsu Kaisen", "shounen")
- **Organizers** — create, update, and delete their own events;
  `organizer`/`admin` role required, enforced at the route level
- **Bookings** — attendees book an event with a phone number and get
  a booking code back; can view and cancel their own bookings
- **Profiles** — favorite-series/interest tags, reusing the same tag
  pool events are tagged with
- **Privacy-aware event detail** — an event's full attendee list
  (phone numbers, booking codes) is only visible to that event's
  organizer or an admin; everyone else sees just a count. A viewer
  can always see their own booking on an event, regardless

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
  model            database entities (User, Event, Booking, Tag, Role, Category)
  dto              API request/response shapes
  middleware       auth (required + optional), role guard, CORS
  router           route registration
  config           typed env config
  apperror         typed errors → HTTP status mapping
  response         standard JSON envelope
docs               API reference + generated OpenAPI spec
```

## Getting started

**Prerequisites:** Go 1.25+, Docker, an [ImageKit](https://imagekit.io)
account (for event image uploads).

1. Copy the env template and fill in your values:
   ```bash
   cp .env.example .env
   ```
   `CLIENT_ORIGIN` should match wherever the frontend runs
   (`http://localhost:5173` by default) — CORS blocks everything else.
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

## Known limitations

- Deleting an account soft-deletes the user row only — their created
  events and bookings aren't cascade-deleted or reassigned, so they'll
  display with a blank organizer/attendee once loaded.
- There's no way to change your role (`attendee`/`organizer`) after
  signing up.

## Development

Regenerate the OpenAPI spec after changing any handler's `@swag` annotations:
```bash
swag init -g cmd/server/main.go -o docs
```
