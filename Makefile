.PHONY: run test test-cover build clean swagger docker-build docker-up docker-down docker-logs dev-db dev-db-down lint

# Development
run:
	go run cmd/api/main.go

dev-db:
	docker-compose -f docker-compose.dev.yml up -d

dev-db-down:
	docker-compose -f docker-compose.dev.yml down

# Testing
test:
	go test ./... -v

test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Build
build:
	go build -o bin/api cmd/api/main.go

clean:
	rm -rf bin/ coverage.out coverage.html

# Documentation
swagger:
	swag init -g cmd/api/main.go -o docs

# Docker
docker-build:
	docker-compose build --no-cache

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f api

docker-restart:
	docker-compose down && docker-compose up -d

# Lint (optional - install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
lint:
	golangci-lint run ./...