.PHONY: run test test-unit lint tidy install

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

include .env
