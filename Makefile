.PHONY: run test test-unit lint tidy install docker-build docker-run docker-compose-up docker-compose-down lgtm-up lgtm-down lgtm-logs lgtm-ps

run:
	LOG_LEVEL=${LOG_LEVEL:-INFO} OTEL_EXPORTER_OTLP_ENDPOINT=${OTEL_EXPORTER_OTLP_ENDPOINT:-} go run ./cmd/api

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

# LGTM stack commands
lgtm-up:
	docker compose up --build -d

lgtm-down:
	docker compose down -v

lgtm-logs:
	docker compose logs -f api alloy otel-collector

lgtm-ps:
	docker compose ps

# OpenTelemetry dev helpers
otel-up:
	OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector:4318 docker compose up --build -d

include .env
