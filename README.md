# Package Ordering Service

A Golang backend service (Gin) and logic to optimally pack orders using configurable pack sizes.

- REST API (Gin), flexible pack sizes, simple web UI.

## Getting Started

1. Install Go 1.24+
2. Run:

   ```sh
   go mod tidy
   go run ./cmd/api/main.go
   ```

Healthcheck endpoint: http://localhost:8080/health

UI: http://localhost:8080/ui/

---

## Calculation Logic

Goal: For an order size `amount` and available `packSizes`, ship only whole packs such that:

1) Minimize total shipped items (allow overfill but keep it minimal).
2) If multiple options ship the same minimal amount, minimize the number of packs.

Implementation:

- The core calculation uses a dynamic programming approach that considers all combinations of pack sizes to strictly minimize total shipped items (overfill), and—among options with the same minimal overfill—minimizes the number of packs.
- The API and UI are designed to be independent of the calculation method, so improvements to the core logic require no changes to clients.

Core types live in `internal/order`:

- `CalculatePacks(amount int, packSizes []int) map[int]int` → returns a size→quantity allocation.

Notes:
- Map order is irrelevant; UI sorts output for presentation.
- Pack sizes are configurable at runtime via the API.

---

## API

Base URL: `http://localhost:8080`

- GET `/health` → `{ "status": "ok" }`

- GET `/packs` → `{ "packSizes": [250,500,1000,...] }` (sorted ascending for readability)

- PUT `/packs`
  - Request: `{ "packSizes": [250,500,1000,2000,5000] }`
  - Behavior: de-duplicates, filters non-positive, stores sorted desc internally.
  - Response: `{ "packSizes": [5000,2000,1000,500,250] }`

- POST `/calculate`
  - Request: `{ "amount": 501 }`
  - Response:
    ```json
    {
      "amount": 501,
      "packSizes": [5000,2000,1000,500,250],
      "allocation": {"500":1, "250":1},
      "totalItems": 750,
      "totalPacks": 2,
      "overfill": 249
    }
    ```

---

## UI

Static HTML (no framework) served from `/ui`.

- Edit pack sizes as a comma-separated list and save (PUT `/packs`).
- Enter an amount and calculate (POST `/calculate`).
- Results table shows pack size and quantity; summary shows totals and overfill.

Open: `http://localhost:8080/ui/`

---

## Development

Makefile targets:

- `make tidy`   – go mod tidy
- `make fmt`    – go fmt ./...
- `make test`   – run unit tests with coverage
- `make run`    – run API locally
- `make build`  – build to `bin/app`
- `make clean`  – remove build artifacts
- `make docker-build` / `make docker-run` – build and run container

Run tests:

```sh
make test
```

---

## Docker

Multi-stage Dockerfile builds a small Alpine image:

```sh
docker build -t package-ordering .
docker run --rm -p 8080:8080 package-ordering
```

UI: `http://localhost:8080/ui/`

---

## Deployment (free options)

- Render (Web Service with Docker): connect GitHub repo, port 8080, health `/health`.
- Fly.io (uses Dockerfile): `flyctl launch && flyctl deploy`.
- Cloud Run (Google): deploy container, set port 8080, allow unauthenticated.

All require no persistent storage; pack sizes are in-memory and configurable via API.