# Package Ordering Service

A Golang backend service (Gin) and logic to optimally pack orders using configurable pack sizes.

- REST API (Gin), flexible pack sizes, simple web UI coming soon.

## Getting Started

1. Install Go 1.24+
2. Run:

   ```sh
   go mod tidy
   go run ./cmd/api/main.go
   ```

Healthcheck endpoint: http://localhost:8080/health