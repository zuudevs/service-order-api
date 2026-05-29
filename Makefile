.PHONY: help build run test lint clean docker-build docker-run docker-stop release deploy

help:
	@echo "Service Order API - Make Commands"
	@echo ""
	@echo "Development:"
	@echo "  make build              - Build binary for current OS"
	@echo "  make run                - Run the application"
	@echo "  make test               - Run tests with coverage"
	@echo "  make test-verbose       - Run tests with verbose output"
	@echo "  make lint               - Run linters (golangci-lint, go vet, gofmt)"
	@echo "  make format             - Format code with gofmt"
	@echo "  make clean              - Remove build artifacts"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build       - Build Docker image locally"
	@echo "  make docker-run         - Run container locally"
	@echo "  make docker-stop        - Stop running container"
	@echo "  make docker-logs        - View container logs"
	@echo "  make docker-compose-up  - Start services with docker-compose"
	@echo "  make docker-compose-down - Stop services with docker-compose"
	@echo ""
	@echo "Dependencies:"
	@echo "  make deps-download      - Download Go dependencies"
	@echo "  make deps-verify        - Verify dependency integrity"
	@echo "  make deps-update        - Update all dependencies"
	@echo ""
	@echo "Database:"
	@echo "  make db-reset           - Reset database (delete storage/app.db)"
	@echo "  make db-backup          - Create database backup"
	@echo ""
	@echo "Release:"
	@echo "  make release VERSION=1.0.0 - Create release tag and push"
	@echo ""
	@echo "Health:"
	@echo "  make health-check       - Check API health endpoint"
	@echo ""

# Development Commands
build:
	@echo "Building for current OS..."
	@mkdir -p bin
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X main.Version=$$(git describe --tags --always --dirty 2>/dev/null || echo 'dev')" -o bin/service-order-api ./cmd/server
	@echo "Build complete: bin/service-order-api"

run: build
	@echo "Running service-order-api..."
	./bin/service-order-api

test:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

test-verbose:
	@echo "Running tests (verbose)..."
	go test -v -race -timeout=10m ./...

lint:
	@echo "Checking format..."
	@if [ -n "$$(gofmt -l .)" ]; then echo "Code needs formatting"; gofmt -l .; exit 1; fi
	@echo "Running go vet..."
	go vet ./...
	@echo "Running golangci-lint..."
	golangci-lint run ./...

format:
	@echo "Formatting code..."
	gofmt -s -w .
	go mod tidy
	@echo "Format complete"

clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean
	@echo "Clean complete"

# Docker Commands
docker-build:
	@echo "Building Docker image..."
	docker build -t service-order-api:latest .
	@echo "Docker build complete"

docker-run: docker-build
	@echo "Running Docker container..."
	docker run -d \
		-p 8080:8080 \
		-e INTERNAL_API_TOKEN_HASH=$${INTERNAL_API_TOKEN_HASH} \
		-v $$(pwd)/storage:/app/storage \
		--name service-order-api \
		service-order-api:latest
	@echo "Container started. Use 'make docker-logs' to view logs"

docker-stop:
	@echo "Stopping Docker container..."
	docker stop service-order-api
	docker rm service-order-api
	@echo "Container stopped"

docker-logs:
	docker logs -f service-order-api

docker-compose-up:
	@echo "Starting services with docker-compose..."
	docker-compose up -d
	@echo "Services started. Check 'docker-compose ps'"

docker-compose-down:
	@echo "Stopping services..."
	docker-compose down

# Dependency Commands
deps-download:
	@echo "Downloading dependencies..."
	go mod download
	@echo "Dependencies downloaded"

deps-verify:
	@echo "Verifying dependencies..."
	go mod verify
	@echo "Dependencies verified"

deps-update:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy
	@echo "Dependencies updated"

# Database Commands
db-reset:
	@echo "Resetting database..."
	rm -f storage/app.db
	@echo "Database reset. It will be recreated on next run."

db-backup:
	@echo "Creating database backup..."
	@mkdir -p storage
	cp storage/app.db storage/app.db.backup.$$(date +%Y%m%d-%H%M%S)
	@echo "Backup created"

# Release Commands
release:
	@if [ -z "$(VERSION)" ]; then echo "Usage: make release VERSION=1.0.0"; exit 1; fi
	@echo "Creating release v$(VERSION)..."
	git tag v$(VERSION)
	git push origin v$(VERSION)
	@echo "Release v$(VERSION) created and pushed"

# Health Check
health-check:
	@echo "Checking API health..."
	@curl -s -o /dev/null -w "HTTP Status: %{http_code}\n" http://localhost:8080/health || echo "API is not responding"

# Combined Commands
all: clean deps-download lint test build
	@echo "All tasks completed successfully"

dev: clean deps-download docker-compose-up
	@echo "Development environment ready"
	@echo "API available at http://localhost:8080"
	@echo "Use 'docker-compose logs -f' to view logs"

ci: deps-verify lint test build
	@echo "CI pipeline completed successfully"
