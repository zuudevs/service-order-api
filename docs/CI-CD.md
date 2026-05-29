# Continuous Integration / Continuous Deployment (CI/CD) Documentation

## Overview

This project uses GitHub Actions to automate testing, building, and deployment workflows. All workflows are triggered automatically on push and pull requests.

## Workflows

### 1. build.yml - Build and Test

**Triggers:** Push to main/develop, Pull Requests

**Jobs:**

- build: Main build job with linting, testing, and binary creation
- build-matrix: Cross-platform builds (Linux, macOS, Windows)

**Steps:**

1. Checkout code
2. Set up Go 1.25.8
3. Download and verify dependencies
4. Format check (gofmt)
5. Lint code (golangci-lint)
6. Vet code (go vet)
7. Run tests with coverage
8. Upload coverage to Codecov
9. Build optimized binary
10. Verify binary
11. Upload artifacts

**Artifacts:**

- service-order-api-linux-amd64
- service-order-api-macos-amd64
- service-order-api-macos-arm64
- service-order-api-windows-amd64.exe

**Status:** SUCCESS indicates:

- Code passes format check
- Code passes linting
- Code passes vetting
- All tests pass
- Binary builds successfully

### 2. docker.yml - Docker Build and Push

**Triggers:** Push to main/develop, Tags (v\*), Pull Requests

**Jobs:**

- build-and-push: Builds and pushes Docker image
- security-scan: Scans Docker image for vulnerabilities

**Steps:**

1. Checkout code
2. Set up Docker Buildx
3. Log in to Container Registry (GHCR)
4. Extract metadata (tags, labels)
5. Build and push Docker image
6. Scan image with Trivy
7. Upload security results

**Docker Registry:** ghcr.io (GitHub Container Registry)

**Image Tags:**

- branch-based: ghcr.io/zuudevs/service-order-api:main
- version-based: ghcr.io/zuudevs/service-order-api:v1.0.0
- short SHA: ghcr.io/zuudevs/service-order-api:abc1234
- latest: ghcr.io/zuudevs/service-order-api:latest (on main branch)

**Security Scan:**
Uses Trivy to scan for vulnerabilities in the Docker image and reports findings to GitHub Security tab.

### 3. release.yml - Release Creation

**Triggers:** Push with tag (v\*)

**Jobs:**

- create-release: Creates GitHub Release
- build-release-artifacts: Builds binaries for all platforms
- publish-docker: Publishes Docker image to registry

**Release Artifacts:**

- service-order-api-linux-amd64 + .sha256
- service-order-api-linux-arm64 + .sha256
- service-order-api-macos-amd64 + .sha256
- service-order-api-macos-arm64 + .sha256
- service-order-api-windows-amd64.exe + .sha256

**Release Information:**

- Automatically creates release notes
- Attaches all compiled binaries
- Includes SHA256 checksums
- Publishes Docker image with version tag

**Usage:**

```bash
git tag v1.0.0
git push origin v1.0.0
```

This triggers:

1. Release page creation on GitHub
2. Binary compilation for all platforms
3. Docker image publish with version tag

### 4. security.yml - Security Scanning

**Triggers:** Push to main/develop, Pull Requests, Daily Schedule (2 AM UTC)

**Jobs:**

- code-analysis: CodeQL analysis for Go
- gosec-scan: Go security issues detection
- dependency-check: Vulnerable dependencies detection
- secrets-scan: Secret/credential detection (TruffleHog)
- sast-scan: Static application security testing
- dependency-review: Pull request dependency review

**Security Tools:**

1. **CodeQL** - GitHub's code analysis engine
2. **Gosec** - Go-specific security scanner
3. **govulncheck** - Vulnerable dependencies check
4. **TruffleHog** - Secret detection
5. **StaticCheck** - Go linter
6. **golangci-lint** - Multiple Go linters

**Results:**

- Displayed in GitHub Security tab
- Fails on critical vulnerabilities in PRs
- Daily scheduled scans for proactive detection

### 5. code-quality.yml - Code Quality

**Triggers:** Push to main/develop, Pull Requests

**Jobs:**

- quality: Code formatting, linting, testing
- docs: Documentation validation
- binary-size: Binary size monitoring

**Quality Checks:**

1. Format check (gofmt)
2. Linting (golangci-lint)
3. Unit tests with race detection
4. Code coverage analysis
5. Documentation verification
6. Binary size check (warns if > 50MB)

**Coverage Reports:**

- Uploaded to Codecov
- HTML coverage report as artifact
- Retention: 30 days

**Metrics:**

- Code coverage percentage
- Test pass rate
- Lint issues count
- Binary size

## Workflow Status Badges

Add these to your README.md:

```markdown
[![Build and Test](https://github.com/zuudevs/service-order-api/actions/workflows/build.yml/badge.svg)](https://github.com/zuudevs/service-order-api/actions/workflows/build.yml)
[![Docker Build](https://github.com/zuudevs/service-order-api/actions/workflows/docker.yml/badge.svg)](https://github.com/zuudevs/service-order-api/actions/workflows/docker.yml)
[![Security](https://github.com/zuudevs/service-order-api/actions/workflows/security.yml/badge.svg)](https://github.com/zuudevs/service-order-api/actions/workflows/security.yml)
[![Code Quality](https://github.com/zuudevs/service-order-api/actions/workflows/code-quality.yml/badge.svg)](https://github.com/zuudevs/service-order-api/actions/workflows/code-quality.yml)
```

## Environment Variables and Secrets

### Required Secrets (for Docker push):

- GITHUB_TOKEN: Automatically provided by GitHub Actions

### Optional Secrets:

- CODECOV_TOKEN: For Codecov integration (optional)
- DOCKER_REGISTRY_USERNAME: Custom registry username (if not GHCR)
- DOCKER_REGISTRY_PASSWORD: Custom registry password (if not GHCR)

## Branch Protection Rules

Recommended branch protection settings for main:

1. Require status checks to pass:
   - build
   - build-matrix
   - quality
   - security-scan
   - sast-scan

2. Require code reviews: 1

3. Require branches to be up to date before merging: Yes

4. Require conversation resolution: Yes

## Common Issues and Solutions

### Issue: Workflow fails with "go mod verify" error

**Solution:**

```bash
go mod tidy
git add go.mod go.sum
git commit -m "Update dependencies"
```

### Issue: Docker push fails with authentication error

**Solution:**
GitHub Token is automatically provided. Ensure GHCR is accessible in repository settings.

### Issue: Test fails in CI but passes locally

**Solution:**

- Use `-race` flag locally: `go test -race ./...`
- Check environment variables in CI workflow
- Verify test isolation and cleanup

### Issue: Binary size warning

**Solution:**
Build was successful but binary is larger than expected.
Check build flags and ensure `-ldflags="-s -w"` is used.

## Performance Optimization

### Build Times

- Caching enabled for Go modules and build artifacts
- Parallel matrix builds for cross-platform compilation
- Docker layer caching for faster rebuilds

### Cost Optimization

- Public runners used (free for public repos)
- Artifact retention limited to 7-30 days
- Docker image layer caching

## Monitoring and Alerts

### GitHub Actions Dashboard

- Navigate to Actions tab in repository
- View workflow runs and logs
- Check status checks on Pull Requests

### Email Notifications

GitHub sends emails for:

- Workflow failures
- Successful releases
- Security alerts

## Deployment

### Development (main branch)

- Automatic Docker build and push
- All artifacts available
- Health checks enabled

### Staging

Configure in repository settings or deployment action

### Production

Use release tags (v\*):

```bash
git tag v1.0.0
git push origin v1.0.0
```

This automatically:

1. Creates GitHub Release
2. Builds all platform binaries
3. Publishes versioned Docker images
4. Generates checksums

## CI/CD Improvement Ideas

Future enhancements:

1. Integration tests with Docker Compose
2. Performance benchmarking
3. Database migration testing
4. API contract testing
5. Load testing
6. Automated deployments to staging/production
7. Slack notifications for failures
8. Dependency update automation (Dependabot)
9. Code generation verification
10. Documentation generation and publishing

## Troubleshooting Workflows

### View Workflow Logs

1. Go to Actions tab
2. Select workflow run
3. Click job name to view logs
4. Search for specific errors

### Re-run Failed Workflows

1. Click the failed workflow run
2. Click "Re-run failed jobs" or "Re-run all jobs"
3. Monitor execution

### Debug Workflow

Add `debug` step:

```yaml
- name: Debug
  run: |
    echo "Environment: $GITHUB_ENV"
    go version
    go env
```

## References

- GitHub Actions: https://docs.github.com/en/actions
- Go Testing: https://golang.org/pkg/testing/
- golangci-lint: https://golangci-lint.run/
- Docker Best Practices: https://docs.docker.com/develop/dev-best-practices/
- Codecov: https://codecov.io/

---

Last Updated: 2026-05-29
