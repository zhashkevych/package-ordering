APP_NAME := package-ordering
BIN := bin/app
PKG := ./cmd/api

.PHONY: all tidy test run build clean docker-build docker-run fmt

all: tidy test build

tidy:
	go mod tidy

fmt:
	go fmt ./...

test:
	go test ./... -cover

run:
	go run $(PKG)

build:
	@mkdir -p bin
	CGO_ENABLED=0 go build -o $(BIN) $(PKG)

clean:
	rm -rf bin

docker-build:
	docker build -t $(APP_NAME) .

docker-run:
	docker run --rm -p 8080:8080 $(APP_NAME)


