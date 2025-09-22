# lensflare-common

Common Go utilities used across Lensflare services: observability (Sentry), HTTP middleware (Gin error handling), database initialization (GORM + Postgres), and Consul service registration.

This repository is a library/module (not an executable). Import the packages you need into your Go service.

## Stack
- Language: Go (module path: `github.com/pyatnitsev/lensflare-common`)
- Go version target: `go 1.25.0` (as declared in go.mod)
  - Note: If your local toolchain does not support 1.25 yet, use the latest available Go version and update go.mod if needed.
- Package manager: Go Modules
- Libraries used:
  - Sentry: `github.com/getsentry/sentry-go`
  - HTTP framework: `github.com/gin-gonic/gin`
  - Consul API: `github.com/hashicorp/consul/api`
  - ORM / DB driver: `gorm.io/gorm`, `gorm.io/driver/postgres`

## Overview of Packages
- `observability/sentry`: Sentry initialization helper.
- `middleware`: Gin middleware for uniform error handling and Sentry reporting.
- `db`: Postgres connection initialization via GORM with a reasonable connection pool.
- `consul`: Service registration in Consul with basic HTTP health check.

## Requirements
- Go toolchain compatible with the module (see go.mod; currently `1.25.0`).
- For DB usage: reachable Postgres instance and correct DSN.
- For Consul usage: reachable Consul agent and service metadata.
- For Sentry usage: valid DSN and optional env/release setup.

## Installation
Add the module to your project using Go modules:

```
go get github.com/pyatnitsev/lensflare-common@latest
```

Then import the packages you need, for example:

```go
import (
    lcdb "github.com/pyatnitsev/lensflare-common/db"
    lcmw "github.com/pyatnitsev/lensflare-common/middleware"
    lcobs "github.com/pyatnitsev/lensflare-common/observability"
    lcconsul "github.com/pyatnitsev/lensflare-common/consul"
)
```

## Usage Examples

### Sentry initialization
```go
// At app startup:
shutdownSentry := lcobs.InitSentry()
// defer graceful flush on shutdown
defer shutdownSentry()
```
Environment variables used by Sentry:
- `SENTRY_DSN` (required to enable; if empty, Sentry is skipped)
- `SENTRY_SAMPLE_RATE` (optional, float; defaults to 1.0)
- `SENTRY_ENV` (optional)
- `SENTRY_RELEASE` (optional)

### Gin error middleware
```go
r := gin.Default()
r.Use(lcmw.ErrorMiddleware())

r.GET("/hello", func(c *gin.Context) {
    // ... your handler
})

_ = r.Run() // start server
```
The middleware logs Gin errors and, if Sentry is initialized and bound to the request context, captures exceptions.

### Database (GORM + Postgres)
```go
db, err := lcdb.Init()
if err != nil {
    log.Fatalf("db init: %v", err)
}
// use `db` (type *gorm.DB)
```
Environment variables used by DB:
- `PG_URL` (required) — Postgres DSN, e.g. `postgres://user:pass@host:5432/dbname?sslmode=disable`

The pool is configured with:
- Max open conns: 20
- Max idle conns: 5
- Conn max lifetime: 30m

### Consul service registration
```go
if err := lcconsul.RegisterService(); err != nil {
    log.Fatalf("consul register: %v", err)
}
```
Environment variables used by Consul registration:
- `CONSUL_HTTP_ADDR` (required) — Consul agent address, e.g. `127.0.0.1:8500` or `http://127.0.0.1:8500`
- `HTTP_PORT` (required) — Service HTTP port (integer)
- `CONSUL_SERVICE_ID` (required) — Unique instance ID
- `CONSUL_SERVICE_NAME` (required) — Service name
- `HOST_ADDRESS` (required) — Host/IP where the service listens (used for health check)

The registration uses an HTTP health check at `http://HOST_ADDRESS:HTTP_PORT/health` with interval `10s` and timeout `3s`. Ensure your service exposes this endpoint.

## Scripts and Entry Points
- This repository provides library packages only. There are no executables or CLI entry points.
- No custom scripts are defined in this repo.
- Common Go commands you may use in consumers:
  - `go mod tidy` — ensure dependencies are resolved.
  - `go build ./...` — build your project.
  - `go test ./...` — run your project tests.

## Environment Variables (summary)
- Sentry: `SENTRY_DSN`, `SENTRY_SAMPLE_RATE`, `SENTRY_ENV`, `SENTRY_RELEASE`
- Database: `PG_URL`
- Consul registration: `CONSUL_HTTP_ADDR`, `HTTP_PORT`, `CONSUL_SERVICE_ID`, `CONSUL_SERVICE_NAME`, `HOST_ADDRESS`

## Project Structure
```
.
├── consul/
│   └── registration.go
├── db/
│   └── db.go
├── middleware/
│   └── error.go
├── observability/
│   └── sentry.go
├── go.mod
├── go.sum
└── README.md
```

## Tests
Unit tests are included for the following packages:
- consul: registration input validation
- db: database initialization behavior
- observability: Sentry initialization behavior without DSN
- middleware: Gin error middleware behavior and Sentry capture

How to run tests:
- Run all tests in the module:
  - `go test ./...`
- Run tests for a single package:
  - `go test ./consul`
- Verbose output:
  - `go test -v ./...`
- With coverage report (terminal):
  - `go test -cover ./...`
- Generate an HTML coverage report:
  - `go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out`

## Development
- Update dependencies: `go get -u` (selectively) and `go mod tidy`.
- Lint/Vet (optional): `go vet ./...`; integrate your preferred linter.

## Versioning
- This module follows standard Go module versioning. Pin a specific version in your consumers as needed.

## License
This project is licensed under the MIT License. See the LICENSE file for details.

## Changelog
- TODO: If needed, add a CHANGELOG.md and keep track of notable changes.
