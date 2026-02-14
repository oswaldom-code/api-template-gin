# API Template Gin

A production-ready Go REST API template built with the [Gin](https://gin-gonic.com) framework, following **Clean Architecture (Hexagonal)** principles. Designed as a starting point for building scalable, maintainable APIs.

## Features

- **Clean Architecture** — Hexagonal layers (adapters, application, domain) with clear dependency flow
- **Route Groups** — Public and protected route groups with middleware support
- **Basic Auth Middleware** — Base64-encoded secret validation on protected routes
- **Prometheus Metrics** — Built-in metrics endpoint at `/metrics` via gin-metrics
- **CLI Support** — Cobra-based CLI with subcommands (`server`, `cli`)
- **Docker Ready** — Multi-stage Dockerfile and docker-compose with PostgreSQL
- **Structured Logging** — Logrus-based logger with configurable log levels
- **GORM + PostgreSQL** — Thread-safe singleton repository with connection pooling
- **Graceful Shutdown** — SIGINT signal handling for clean server termination
- **OpenAPI Spec** — API defined in `swagger/swagger.yml` (OpenAPI 3.0.3)

## Project Structure

```
.
├── main.go                          # Entrypoint (Cobra root command)
├── cmd/api/
│   └── main.go                      # Server and CLI subcommand registration
├── src/
│   ├── adapters/                    # External interfaces (driven/driving)
│   │   ├── http/rest/
│   │   │   ├── handlers/            # Gin HTTP handlers (implement ServerInterface)
│   │   │   ├── dto/                 # Request/response DTOs (Response, ResponseWithData)
│   │   │   └── infrastructure/      # Gin engine setup, route registration, middleware
│   │   ├── repository/              # GORM data access (PostgreSQL, sync.Once singleton)
│   │   └── cli/                     # CLI adapter (Cobra subcommand)
│   ├── application/                 # Use cases and business logic
│   │   └── system_services/
│   │       ├── health.go            # Health service
│   │       └── ports/               # Interfaces/contracts (e.g., Store)
│   └── domain/                      # Domain entities and services (ready for expansion)
├── pkg/                             # Shared reusable packages
│   ├── config/                      # Configuration (godotenv + pflag + env vars)
│   └── log/                         # Structured logging wrapper (logrus)
├── swagger/
│   └── swagger.yml                  # OpenAPI 3.0.3 specification
├── Dockerfile                       # Multi-stage build (alpine)
├── docker-compose.yml               # App + PostgreSQL with healthcheck
└── .env.example                     # Environment variables template
```

## Quick Start

### Prerequisites

- [Go](https://golang.org/doc/install) 1.24 or higher
- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/) (optional, for containerized setup)

### Setup

```bash
# Clone the repository
git clone https://github.com/oswaldom-code/api-template-gin.git
cd api-template-gin

# Copy environment variables
cp .env.example .env
# Edit .env with your configuration

# Install dependencies
go mod tidy

# Run the server
go run main.go server
```

The server starts at `http://localhost:9000` by default.

## Commands

All commands are available via `make`. Run `make help` to list all targets.

| Make Target | Description |
|-------------|-------------|
| `make run` | Start the HTTP server (development) |
| `make build` | Build production binary |
| `make test` | Run all tests |
| `make test-cover` | Run tests with coverage report |
| `make lint` | Run vet + check formatting |
| `make fmt` | Format all Go source files |
| `make tidy` | Tidy module dependencies |
| `make docker-up` | Run app + PostgreSQL with Docker |
| `make docker-down` | Stop Docker services |
| `make swagger` | Generate API docs from OpenAPI spec |
| `make clean` | Remove build artifacts |
| `make help` | Show all available targets |

<details>
<summary>Direct Go commands</summary>

| Command | Description |
|---------|-------------|
| `go run main.go server` | Start the HTTP server (development) |
| `go run main.go cli -f test` | Run CLI utilities (e.g., test DB connection) |
| `go test ./...` | Run all tests |
| `go test ./pkg/log -v` | Run tests for a specific package |
| `go test ./pkg/log -run TestSetLogLevel -v` | Run a specific test |
| `CGO_ENABLED=0 go build -o bin/api main.go` | Build production binary |
| `docker-compose up --build` | Run app + PostgreSQL with Docker |
| `./scripts/swagger.sh` | Generate API docs from OpenAPI spec |

</details>

## Configuration

All configuration is loaded via `.env` file, CLI flags (pflag), and environment variables. See `.env.example` for the full template.

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_HOST` | Server bind address | `localhost` |
| `SERVER_PORT` | Server port | `9000` |
| `SERVER_SCHEME` | URL scheme | `http` |
| `SERVER_MODE` | Gin mode (`debug` / `release`) | `debug` |
| `DB_USER` | Database user | — |
| `DB_PASSWORD` | Database password | — |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_DATABASE` | Database name | — |
| `DB_ENGINE` | Database engine | `postgres` |
| `DB_SSL_MODE` | SSL mode | `disable` |
| `DB_MAX_CONNECTIONS` | Max connection pool size | `10` |
| `DB_LOG_MODE` | GORM log level | `debug` |
| `LOG_LEVEL` | Application log level | `info` |
| `LOG_ERROR_LOG_FILE` | Error log file path | — |
| `AUTH_SECRET` | Secret for Basic Auth middleware | — |
| `ENVIRONMENT` | Runtime environment | `development` |

## Architecture

This project follows Clean Architecture (Hexagonal) with clear separation of concerns.

**Dependency flow:**

```
Handlers → Application Services → Ports (interfaces) ← Repository (implementation)
```

- **Handlers** receive HTTP requests and delegate to application services
- **Application Services** contain business logic and depend on port interfaces
- **Ports** define contracts that adapters must implement (Dependency Inversion)
- **Repository** implements port interfaces for data persistence

For detailed architecture diagrams (layers, request flow, auth middleware, startup sequence, and class diagrams), see [docs/diagrams.md](doc/diagrams.md).

## API Endpoints

| Method | Path | Auth | Description | Response |
|--------|------|------|-------------|----------|
| `GET` | `/ping` | No | Health check / ping | `{"status": true, "message": "pong"}` |
| `GET` | `/metrics` | No | Prometheus metrics | Prometheus text format |

Protected routes use Basic Auth — send the `Authorization` header with `Basic <base64-encoded AUTH_SECRET>`.

The full API specification is available in [`swagger/swagger.yml`](swagger/swagger.yml).

## Docker

### Docker Compose (recommended)

Runs the application with a PostgreSQL database:

```bash
# Copy and configure environment variables
cp .env.example .env

# Start all services
docker-compose up --build
```

This starts:
- **app** — API server on port `9000`
- **db** — PostgreSQL 16 on port `5432` with healthcheck and persistent volume

### Standalone Docker

```bash
# Build the image
docker build -t api-template-app .

# Run the container
docker run -d -p 9000:9000 --env-file .env api-template-app
```

## Testing

This project uses `testing` + [`testify`](https://github.com/stretchr/testify) for assertions.

```bash
# Run all tests
go test ./...

# Run with verbose output
go test ./... -v

# Run a specific package
go test ./pkg/log -v

# Run a specific test
go test ./pkg/log -run TestSetLogLevel -v

# Run with coverage
go test ./... -cover
```

## Load Test Results

Stress tests executed against the `/ping` endpoint using Go's `net/http` + `httptest` (no external tools required). Run them with:

```bash
make test                           # includes load tests
go test ./src/adapters/http/rest/infrastructure/ -run TestLoadPing -v   # load test only
go test ./src/adapters/http/rest/infrastructure/ -bench=. -benchmem     # benchmark only
```

### Concurrent Load Test — `/ping`

| Metric | Value |
|--------|-------|
| Total Requests | 10,000 |
| Concurrency | 100 |
| Duration | 276ms |
| Throughput | **36,240 req/s** |
| Success | 10,000 (100%) |
| Failures | 0 |
| Latency avg | 2.516ms |
| Latency p50 | 1.879ms |
| Latency p95 | 5.846ms |
| Latency p99 | 15.583ms |

### Benchmark — `/ping` Handler

| Metric | Value |
|--------|-------|
| Iterations | 593,370 |
| Throughput | **500,585 req/s** |
| Latency avg | 1.998 us/op |
| Memory | 6,186 B/op |
| Allocations | 20 allocs/op |

### Auth Middleware Load Test — `/secure-ping` (no token)

| Metric | Value |
|--------|-------|
| Total Requests | 5,000 |
| Concurrency | 50 |
| Duration | 176ms |
| Throughput | **28,463 req/s** |
| Rejected (401) | 5,000 (100%) |

> Results from: Intel Core i7-1255U, Linux 6.17, Go test with `httptest.NewServer`. Date: 2026-02-14.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/my-feature`)
3. Commit your changes (`git commit -m 'Add my feature'`)
4. Push to the branch (`git push origin feature/my-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.
