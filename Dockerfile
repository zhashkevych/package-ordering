## Multi-stage build for small production image

FROM golang:1.24-alpine AS builder
WORKDIR /src

# Install build deps
RUN apk add --no-cache ca-certificates git build-base

# Cache go mod
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build static binary
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o /out/app ./cmd/api

FROM alpine:3.20 AS runner
WORKDIR /app
RUN apk add --no-cache ca-certificates && adduser -D -H app
USER app

COPY --from=builder /out/app /app/app
COPY web ./web

EXPOSE 8080
ENTRYPOINT ["/app/app"]

