.PHONY: help run build test clean migrate-up migrate-down

help:
	@echo "Available commands:"
	@echo "  make run          - Run the application in development mode"
	@echo "  make build        - Build the application binary"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make migrate-up   - Run database migrations"
	@echo "  make deps         - Download dependencies"

run:
	go run cmd/main.go

build:
	go build -o bin/bruvela-api cmd/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/
	go clean

migrate-up:
	psql -U postgres -d bruvela_db -f migrations/schema.sql

deps:
	go mod download
	go mod tidy

dev:
	air
