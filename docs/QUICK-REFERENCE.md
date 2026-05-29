# Quick Reference: CI/CD and Deployment

## Essential Commands

### Local Development

```bash
# Setup and run locally
make dev
# or
docker-compose up -d

# Check health
curl http://localhost:8080/health

# View logs
docker-compose logs -f
```

### Build and Test

```bash
# Build
make build

# Run tests
make test

# Run linters
make lint

# Format code
make format
```

### Docker

```bash
# Build image
make docker-build

# Run container
make docker-run

# Stop container
make docker-stop
```

### Release

```bash
# Create and push version tag
make release VERSION=1.0.0
# This triggers the release workflow automatically
```

## Workflow Triggers

| Workflow         | Trigger                  | Action                |
| ---------------- | ------------------------ | --------------------- |
| build.yml        | Push to main/develop, PR | Build, test, lint     |
| code-quality.yml | Push to main/develop, PR | Coverage, size checks |
| security.yml     | Push, PR, Daily 2 AM UTC | Security scans        |
| docker.yml       | Push, tags, PR           | Docker build/push     |
| release.yml      | Tag v\* pattern          | Create release        |

## CI/CD Status

View status at: https://github.com/zuudevs/service-order-api/actions

### Status Checks

- All commits must pass build and code-quality checks before merge
- Security scan failures on PRs prevent merge
- Docker image scan results visible in security tab

## Deployment

### Development

```bash
docker-compose up -d
```

### Production (Docker)

```bash
docker pull ghcr.io/zuudevs/service-order-api:latest
docker run -d -p 8080:8080 \
  -e INTERNAL_API_TOKEN_HASH=your_hash \
  -v /var/lib/service-order-api:/app/storage \
  ghcr.io/zuudevs/service-order-api:latest
```

### Production (Binary)

```bash
# Download from release
wget https://github.com/zuudevs/service-order-api/releases/download/v1.0.0/service-order-api-linux-amd64
chmod +x service-order-api-linux-amd64
./service-order-api-linux-amd64
```

### Production (Kubernetes)

```bash
kubectl apply -f k8s-deployment.yaml
```

## Environment Setup

### Local Development (.env)

```env
PORT=8080
INTERNAL_API_TOKEN_HASH=your_bcrypt_hash
```

### Docker Secrets

Set INTERNAL_API_TOKEN_HASH via:

- Environment variable
- Docker secrets
- Kubernetes secrets
- Container orchestration platform

## Health Check

```bash
# Local
curl http://localhost:8080/health

# Kubernetes
kubectl get pods -n service-order

# Docker
docker ps --filter "name=service-order-api"
```

## Documentation

- **Full CI/CD Guide**: `docs/CI-CD.md`
- **Deployment Guide**: `docs/DEPLOYMENT.md`
- **Available Commands**: `make help`
- **API Reference**: `docs/API.md`
- **Architecture**: `docs/ARCHITECTURE.md`
- **Build Instructions**: `docs/BUILD.md`

## Release Checklist

Before creating a release:

- [ ] All tests passing
- [ ] Code review complete
- [ ] Version updated (if applicable)
- [ ] CHANGELOG updated

```bash
# Create release
make release VERSION=1.0.0
```

After push:

- [ ] GitHub Actions workflows execute
- [ ] Release page created with binaries
- [ ] Docker image pushed to registry
- [ ] All platforms built successfully

## Troubleshooting

### Workflow Failed?

1. Check Actions tab for logs
2. View detailed error messages
3. Fix and push again

### Docker Won't Start?

```bash
docker logs service-order-api
# Check env vars and volume mounts
```

### Tests Failing?

```bash
make test-verbose
# Run with detailed output
```

### Can't Connect to API?

```bash
# Check if running
docker ps | grep service-order-api
# Check health
curl -v http://localhost:8080/health
# Check logs
docker logs service-order-api
```

## Security

- API requires INTERNAL_API_TOKEN_HASH in Authorization header
- Use strong, randomly generated tokens
- Rotate tokens regularly
- Store secrets securely (not in code)
- Container images scanned for vulnerabilities

## Performance

- Multi-platform builds (Linux, macOS, Windows)
- Docker layer caching for faster builds
- Binary size monitoring
- Resource limits configured in Kubernetes

## Support

- GitHub Issues: Report bugs
- GitHub Discussions: Ask questions
- Pull Requests: Submit improvements
- Email: zuudevs@gmail.com

---

Last Updated: 2026-05-29
