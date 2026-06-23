.PHONY: run test test-unit lint tidy install docker-build docker-run docker-compose-up docker-compose-down

run:
	go run ./cmd/api

test:
	go test -race ./...

test-unit:
	go test -short ./...

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

install:
	go mod download
	go install github.com/air-verse/air@latest

# Docker commands
docker-build:
	docker build -t go-clean-template:latest .

docker-run:
	docker run --rm -p 8080:8080 --name go-clean-template go-clean-template:latest

docker-compose-up:
	docker compose up --build -d

docker-compose-down:
	docker compose down

# OpenTelemetry example (requires otel-collector service uncommented in docker-compose.yaml)
otel-up:
	OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector:4318 docker compose up --build -d

include .env
