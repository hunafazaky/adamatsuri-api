# AdaMatsuri API

REST API for discovering and booking anime-community events such as conventions, doujin markets, screenings, cosplay contests, game tournaments, and meetups. The platform supports organizer-managed event creation, attendee bookings, role-based access, JWT authentication, and Swagger-generated API documentation.

---

## Tech Stack

- **Runtime / Language:** Go 1.27
- **Framework:** Gin
- **Database:** PostgreSQL
- **Tooling & Infrastructure:** Docker, Docker Compose, GORM, Swagger, Render-ready production image

---

## Prerequisites

Ensure you have the following installed on your machine before getting started:

- Go >= 1.27
- Docker & Docker Compose
- PostgreSQL (or use the Docker Compose service included in this project)
- An ImageKit account if you plan to upload event images

---

## Getting Started

### 1. Clone the Repository

```bash
git clone <repository-url>
cd adamatsuri-api
```

### 2. Environment Setup

Copy the example environment file and configure your local variables:

```bash
cp .env.example .env
```

At minimum, set the values used by the app at runtime. The project supports either a single `DB_URI` or the individual `POSTGRES_*` settings:

```env
PORT=8080
DB_URI=postgres://admin:adminpassword@localhost:5432/mydb?sslmode=disable
# or:
# POSTGRES_HOST=localhost
# POSTGRES_PORT=5432
# POSTGRES_USER=admin
# POSTGRES_PASSWORD=adminpassword
# POSTGRES_DB_NAME=mydb
# POSTGRES_SSLMODE=disable

CLIENT_ORIGIN=http://localhost:5173
JWT_SECRET=your_jwt_secret
IMAGEKIT_PRIVATE_KEY=your_imagekit_private_key
PUBLIC_HOST=
```

> `CLIENT_ORIGIN` should match the frontend origin that is allowed to call the API from the browser.

### 3. Install Dependencies

```bash
go mod download
```

---

## Running the Application

### Development Mode

Start the API directly with Go:

```bash
go run ./cmd/server
```

Or run the app and PostgreSQL together with Docker Compose:

```bash
docker compose up --build
```

The API will be available at:

- http://localhost:8080
- API docs: http://localhost:8080/docs

### Production Mode

Build the production container and run it with the same environment variables:

```bash
docker build -f Dockerfile.prod -t adamatsuri-api .
docker run --rm -p 8080:8080 --env-file .env adamatsuri-api
```

The project is also configured for deployment on Render via [render.yaml](render.yaml).

---

## Available Scripts

| Script | Description |
| :--- | :--- |
| `go run ./cmd/server` | Start the API locally in development mode. |
| `docker compose up --build` | Start the PostgreSQL database and API using Docker Compose. |
| `swag init -g cmd/server/main.go -o docs` | Regenerate the Swagger/OpenAPI specification from handler annotations. |
| `docker build -f Dockerfile.prod -t adamatsuri-api .` | Build the production Docker image. |

---

## API Documentation

The project exposes Swagger and Scalar-based documentation through the built-in routes:

- API docs: http://localhost:8080/docs
- Swagger UI: http://localhost:8080/swagger/index.html
- OpenAPI JSON: http://localhost:8080/openapi.json

Generated spec files are stored in the [docs](docs) directory and include the JSON/YAML exports for the API contract.

---

## Testing

There is no automated test suite currently committed to this repository. When tests are added, the standard Go command is:

```bash
go test ./...
```

---

## Project Structure

```text
.
├── cmd/
│   └── server/
│       └── main.go
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal/
│   ├── apperror/
│   ├── config/
│   ├── dto/
│   ├── handler/
│   ├── middleware/
│   ├── model/
│   ├── repository/
│   ├── response/
│   ├── router/
│   └── service/
├── .air.toml
├── .env.example
├── docker-compose.yaml
├── Dockerfile
├── Dockerfile.prod
├── go.mod
├── go.sum
├── README.md
├── readme-context.md
├── render.yaml
└── tmp/
``` 
