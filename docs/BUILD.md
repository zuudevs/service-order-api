# Build Guide

This document provides comprehensive instructions for building the Service Order API from source.

## Prerequisites

### Required

- **Go 1.25.8 or later** - [Download](https://golang.org/dl/)
- **Git** - For version control
- **SQLite3** - Usually included with Go dependencies

### Optional

- **Make** - For build automation
- **Docker** - For containerized builds

## Building from Source

### Basic Build

```bash
# Navigate to project root
cd service-order-api

# Download dependencies
go mod download

# Build the executable
go build -o service-order-api ./cmd/server
```

The compiled binary will be created as `service-order-api` (or `service-order-api.exe` on Windows).

### Build with Custom Output

```bash
# Specify output directory
go build -o ./bin/service-order-api ./cmd/server

# Build with specific name
go build -o my-api ./cmd/server
```

## Cross-Platform Compilation

Build for different operating systems and architectures.

### Linux (64-bit)

```bash
GOOS=linux GOARCH=amd64 go build -o service-order-api-linux ./cmd/server
```

### macOS (ARM64 - Apple Silicon)

```bash
GOOS=darwin GOARCH=arm64 go build -o service-order-api-macos ./cmd/server
```

### macOS (Intel)

```bash
GOOS=darwin GOARCH=amd64 go build -o service-order-api-macos ./cmd/server
```

### Windows (64-bit)

```bash
GOOS=windows GOARCH=amd64 go build -o service-order-api.exe ./cmd/server
```

### Windows (32-bit)

```bash
GOOS=windows GOARCH=386 go build -o service-order-api.exe ./cmd/server
```

## Build with Version Information

Embed version information during build:

```bash
VERSION=1.0.0
COMMIT=$(git rev-parse --short HEAD)
BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ')

go build \
  -ldflags="-X main.Version=$VERSION -X main.Commit=$COMMIT -X main.BuildDate=$BUILD_DATE" \
  -o service-order-api ./cmd/server
```

## Optimized Release Build

Create an optimized, smaller binary for production:

```bash
CGO_ENABLED=0 \
  go build \
  -trimpath \
  -ldflags="-s -w" \
  -o service-order-api ./cmd/server
```

**Flags explained:**

- `CGO_ENABLED=0` - Disable C bindings for portable binary
- `-trimpath` - Remove file system paths from binary
- `-s` - Strip symbol table
- `-w` - Strip debugging information

## Dependency Management

### Download Dependencies

```bash
# Download all dependencies
go mod download

# Download with verification
go mod download -json
```

### View Dependencies

```bash
# List all dependencies
go list -m all

# Display dependency graph
go mod graph
```

### Update Dependencies

```bash
# Check for updates
go list -u -m all

# Update specific dependency
go get -u github.com/package/name

# Update all dependencies
go get -u ./...
```

### Clean Dependencies

```bash
# Remove unused dependencies
go mod tidy

# Verify module
go mod verify
```

## Code Quality

### Format Code

```bash
# Format all Go files
go fmt ./...

# Or use gofmt
gofmt -s -w .
```

### Lint Code

If `golangci-lint` is installed:

```bash
golangci-lint run ./...
```

### Type Check

```bash
# Type check without building
go vet ./...
```

## Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Build Verification

After building, verify the binary works:

### Linux/macOS

```bash
# Test the binary
./service-order-api &

# Check if running
curl http://localhost:8080/health

# Stop the server
pkill service-order-api
```

### Windows

```cmd
# Test the binary
service-order-api.exe

# Check if running (in another terminal)
curl http://localhost:8080/health

# Stop with Ctrl+C
```

## Docker Build

### Build Docker Image

Create a `Dockerfile`:

```dockerfile
FROM golang:1.25.8-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o service-order-api ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/service-order-api .
COPY --from=builder /app/storage ./storage
EXPOSE 8080
CMD ["./service-order-api"]
```

Build and run:

```bash
# Build image
docker build -t service-order-api:latest .

# Run container
docker run -p 8080:8080 -e PORT=8080 service-order-api:latest
```

## CI/CD Integration

### GitHub Actions Example

Create `.github/workflows/build.yml`:

```yaml
name: Build

on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        os: [linux, darwin, windows]
        arch: [amd64, arm64]

    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: "1.25.8"

      - name: Build
        run: |
          GOOS=${{ matrix.os }} \
          GOARCH=${{ matrix.arch }} \
          go build -o service-order-api-${{ matrix.os }}-${{ matrix.arch }} \
          ./cmd/server
```

## Troubleshooting

### Build Fails with Module Errors

```bash
# Clear Go cache
go clean -modcache

# Reinitialize modules
rm go.sum
go mod tidy
go mod download
```

### CGO Compilation Errors

```bash
# Ensure C compiler is installed
# Windows: Install MinGW or Visual Studio
# macOS: Install Xcode Command Line Tools
# Linux: Install build-essential

# Or disable CGO
CGO_ENABLED=0 go build -o service-order-api ./cmd/server
```

### Permission Denied on Linux/macOS

```bash
# Make binary executable
chmod +x service-order-api

# Then run
./service-order-api
```

## Performance Tips

1. **Use `-ldflags="-s -w"`** for smaller binary size
2. **Enable `-trimpath`** to remove file paths
3. **Use `GOOS` and `GOARCH`** for target-specific optimization
4. **Test with `-race`** flag to detect race conditions

## Build Checklist

- [ ] Dependencies downloaded (`go mod download`)
- [ ] Code formatted (`go fmt ./...`)
- [ ] Code linted (`go vet ./...`)
- [ ] Tests pass (`go test ./...`)
- [ ] Binary builds successfully (`go build`)
- [ ] Binary runs without errors
- [ ] Health endpoint responds (`/health`)
- [ ] All required environment variables set

---

For more information, see [Go Build Documentation](https://golang.org/cmd/go/#hdr-Compile_packages_and_dependencies).
