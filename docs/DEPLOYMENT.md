# Deployment Guide

## Overview

This guide covers deploying the Service Order API to various environments using Docker, GitHub Actions, and manual deployment methods.

## Prerequisites

- Docker and Docker Compose (for containerized deployments)
- GitHub repository with Actions enabled
- Server with SSH access (for manual deployments)
- Proper environment variables configured

## Deployment Methods

### 1. Docker Compose (Local/Development)

#### Quick Start

```bash
# Clone and setup
git clone https://github.com/zuudevs/service-order-api.git
cd service-order-api

# Configure environment
cp .env.example .env
# Edit .env with your settings

# Start services
docker-compose up -d

# Verify deployment
curl http://localhost:8080/health

# View logs
docker-compose logs -f service-order-api

# Stop services
docker-compose down
```

#### Docker Compose Configuration

The included `docker-compose.yml` provides:

- Service container running on port 8080
- Volume mount for persistent storage
- Health checks configured
- Environment variable support

### 2. Docker Image Deployment

#### Build Locally

```bash
# Build Docker image
docker build -t service-order-api:latest .

# Run container
docker run -d \
  -p 8080:8080 \
  -e INTERNAL_API_TOKEN_HASH=your_token_hash \
  -v $(pwd)/storage:/app/storage \
  --name service-order-api \
  service-order-api:latest

# Verify running
docker ps
curl http://localhost:8080/health
```

#### Using GitHub Container Registry

```bash
# Pull image from GHCR
docker pull ghcr.io/zuudevs/service-order-api:latest

# Run container
docker run -d \
  -p 8080:8080 \
  -e INTERNAL_API_TOKEN_HASH=your_token_hash \
  -v /var/lib/service-order-api:/app/storage \
  ghcr.io/zuudevs/service-order-api:latest
```

### 3. Kubernetes Deployment

#### Create Namespace

```bash
kubectl create namespace service-order
```

#### Create ConfigMap for Configuration

```bash
kubectl create configmap service-order-config \
  --from-literal=PORT=8080 \
  -n service-order
```

#### Create Secret for Sensitive Data

```bash
kubectl create secret generic service-order-secrets \
  --from-literal=internal-api-token-hash=YOUR_TOKEN_HASH \
  -n service-order
```

#### Deploy Using Manifest

Create `k8s-deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: service-order-api
  namespace: service-order
spec:
  replicas: 3
  selector:
    matchLabels:
      app: service-order-api
  template:
    metadata:
      labels:
        app: service-order-api
    spec:
      containers:
        - name: service-order-api
          image: ghcr.io/zuudevs/service-order-api:latest
          imagePullPolicy: IfNotPresent
          ports:
            - containerPort: 8080
          env:
            - name: PORT
              valueFrom:
                configMapKeyRef:
                  name: service-order-config
                  key: PORT
            - name: INTERNAL_API_TOKEN_HASH
              valueFrom:
                secretKeyRef:
                  name: service-order-secrets
                  key: internal-api-token-hash
          volumeMounts:
            - name: storage
              mountPath: /app/storage
          livenessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 10
            periodSeconds: 30
          readinessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 10
          resources:
            requests:
              memory: "128Mi"
              cpu: "100m"
            limits:
              memory: "512Mi"
              cpu: "500m"
      volumes:
        - name: storage
          persistentVolumeClaim:
            claimName: service-order-pvc
---
apiVersion: v1
kind: Service
metadata:
  name: service-order-api
  namespace: service-order
spec:
  type: LoadBalancer
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
  selector:
    app: service-order-api
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: service-order-pvc
  namespace: service-order
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
```

Deploy:

```bash
kubectl apply -f k8s-deployment.yaml

# Verify deployment
kubectl get deployment -n service-order
kubectl get pods -n service-order
kubectl get svc -n service-order

# View logs
kubectl logs -f deployment/service-order-api -n service-order
```

### 4. Manual Binary Deployment

#### Build Binary

```bash
# Linux
GOOS=linux GOARCH=amd64 go build \
  -ldflags="-s -w -X main.Version=1.0.0" \
  -o service-order-api \
  ./cmd/server

# Or download from GitHub Releases
wget https://github.com/zuudevs/service-order-api/releases/download/v1.0.0/service-order-api-linux-amd64
chmod +x service-order-api
```

#### Setup Systemd Service

Create `/etc/systemd/system/service-order-api.service`:

```ini
[Unit]
Description=Service Order API
After=network.target
Wants=network-online.target

[Service]
Type=simple
User=api
WorkingDirectory=/opt/service-order-api
ExecStart=/opt/service-order-api/service-order-api
Restart=on-failure
RestartSec=10

Environment="PORT=8080"
Environment="INTERNAL_API_TOKEN_HASH=your_token_hash"

StandardOutput=journal
StandardError=journal
SyslogIdentifier=service-order-api

[Install]
WantedBy=multi-user.target
```

#### Deploy

```bash
# Copy binary
sudo cp service-order-api /opt/service-order-api/

# Set permissions
sudo chown -R api:api /opt/service-order-api
sudo chmod 755 /opt/service-order-api/service-order-api

# Enable and start service
sudo systemctl daemon-reload
sudo systemctl enable service-order-api
sudo systemctl start service-order-api

# Check status
sudo systemctl status service-order-api

# View logs
sudo journalctl -u service-order-api -f
```

### 5. GitHub Actions Automatic Deployment

The CI/CD pipeline automatically:

1. **Build and Test**
   - Runs on every push to main/develop
   - Executes linting, tests, and builds

2. **Docker Build**
   - Builds Docker image
   - Pushes to GitHub Container Registry
   - Runs security scans

3. **Release**
   - Creates GitHub Release on version tags
   - Publishes binaries for all platforms
   - Pushes versioned Docker images

#### Trigger Deployment

```bash
# Create and push a version tag
git tag v1.0.0
git push origin v1.0.0
```

This automatically:

- Creates GitHub Release
- Builds release binaries
- Publishes Docker image with version tag

## Environment Configuration

### Production Environment Variables

```env
# Server
PORT=8080

# Authentication
INTERNAL_API_TOKEN_HASH=bcrypt_hashed_token

# Google Drive (optional)
GOOGLE_DRIVE_CREDENTIALS={"type":"service_account",...}
GOOGLE_DRIVE_TOKEN=access_token
GOOGLE_DRIVE_BACKUP_FOLDER_ID=folder_id
GOOGLE_DRIVE_DB_FILE_ID=file_id

# Logging (optional)
LOG_LEVEL=info
LOG_FORMAT=json
```

### Generate Token Hash

```bash
# Using bcrypt-cli
go install github.com/bitnami/bcrypt-cli@latest
echo -n "your_token_here" | bcrypt-cli hash

# Or in Python
python3 -c "import bcrypt; print(bcrypt.hashpw(b'your_token_here', bcrypt.gensalt()).decode())"
```

## Health Checks and Monitoring

### Health Endpoint

```bash
curl -s http://localhost:8080/health
```

Response: `ok`

### Docker Health Check

The Docker image includes a HEALTHCHECK instruction that monitors the `/health` endpoint every 30 seconds.

```bash
# Check container health
docker ps --filter "name=service-order-api"
# Status shows: "healthy" or "unhealthy"
```

### Kubernetes Probes

- **Liveness Probe**: Checks if container is alive
- **Readiness Probe**: Checks if container can accept traffic

Both configured to use `/health` endpoint.

## Backup and Recovery

### Database Backup

```bash
# Manual backup
cp storage/app.db storage/app.db.backup.$(date +%Y%m%d-%H%M%S)

# List backups
ls -lh storage/*.backup
```

### Google Drive Backup

Configured in `.env`:

- Automatic backups every 6 hours
- Files stored in Google Drive folder
- Manual trigger by restarting service

## Scaling

### Horizontal Scaling

#### Docker Swarm

```bash
# Initialize swarm
docker swarm init

# Deploy service
docker service create \
  --name service-order-api \
  --replicas 3 \
  -p 8080:8080 \
  -e INTERNAL_API_TOKEN_HASH=token \
  ghcr.io/zuudevs/service-order-api:latest

# Scale service
docker service scale service-order-api=5
```

#### Kubernetes

```bash
# Scale deployment
kubectl scale deployment service-order-api --replicas=5 -n service-order

# Auto-scaling (requires metrics server)
kubectl autoscale deployment service-order-api \
  --min=3 --max=10 \
  --cpu-percent=80 \
  -n service-order
```

### Load Balancing

- **Docker Swarm**: Built-in load balancing across replicas
- **Kubernetes**: Service type LoadBalancer distributes traffic
- **Manual**: Use Nginx, HAProxy, or cloud provider load balancer

## Rollback

### Docker Rollback

```bash
# Pull previous version
docker pull ghcr.io/zuudevs/service-order-api:v1.0.0

# Stop current container
docker stop service-order-api

# Run previous version
docker run -d \
  -p 8080:8080 \
  -e INTERNAL_API_TOKEN_HASH=token \
  -v /var/lib/service-order-api:/app/storage \
  --name service-order-api \
  ghcr.io/zuudevs/service-order-api:v1.0.0
```

### Kubernetes Rollback

```bash
# View rollout history
kubectl rollout history deployment/service-order-api -n service-order

# Rollback to previous version
kubectl rollout undo deployment/service-order-api -n service-order

# Rollback to specific version
kubectl rollout undo deployment/service-order-api --to-revision=2 -n service-order
```

## Monitoring and Logs

### Docker Logs

```bash
# View logs
docker logs service-order-api

# Follow logs
docker logs -f service-order-api

# Last 100 lines
docker logs --tail=100 service-order-api
```

### Kubernetes Logs

```bash
# Pod logs
kubectl logs deployment/service-order-api -n service-order

# Follow logs
kubectl logs -f deployment/service-order-api -n service-order

# Previous crash logs
kubectl logs --previous pod-name -n service-order
```

### Systemd Logs

```bash
# View logs
journalctl -u service-order-api

# Follow logs
journalctl -u service-order-api -f

# Last hour
journalctl -u service-order-api --since "1 hour ago"
```

## Troubleshooting

### Container won't start

```bash
# Check logs
docker logs service-order-api

# Check environment variables
docker inspect service-order-api | grep -A 20 "Env"

# Test binary directly
docker run -it ghcr.io/zuudevs/service-order-api:latest /bin/sh
```

### Connection refused

```bash
# Check port binding
docker port service-order-api
netstat -tulpn | grep 8080

# Test endpoint
curl -v http://localhost:8080/health
```

### Database errors

```bash
# Check volume mount
docker inspect service-order-api | grep -A 10 "Mounts"

# Check database file
ls -la /var/lib/service-order-api/

# Restore from backup if needed
cp /var/lib/service-order-api/app.db.backup /var/lib/service-order-api/app.db
docker restart service-order-api
```

## Security Best Practices

1. **Keep Docker Images Updated**

   ```bash
   docker pull ghcr.io/zuudevs/service-order-api:latest
   docker run ... ghcr.io/zuudevs/service-order-api:latest
   ```

2. **Use Secrets Management**
   - Never commit `.env` to repository
   - Use Kubernetes Secrets for production
   - Use environment variables or secret management services

3. **Network Security**
   - Use HTTPS/TLS in production
   - Configure firewall rules
   - Use VPN or private networks for sensitive operations

4. **Database Security**
   - Keep backups encrypted
   - Use file permissions (600 for database files)
   - Regular backup testing

5. **Token Rotation**
   - Periodically rotate API tokens
   - Use strong, random tokens
   - Monitor token usage

## Performance Tuning

### Go Application Tuning

```bash
# Set GOMAXPROCS for CPU cores
export GOMAXPROCS=4

# Memory limits
export GOMEMLIMIT=512MiB
```

### Docker Resource Limits

```bash
docker run -d \
  -p 8080:8080 \
  --memory=512m \
  --cpus=1 \
  ...
```

### Kubernetes Resource Requests

See K8s deployment manifest above for resource requests and limits.

## References

- Docker: https://docs.docker.com/
- Kubernetes: https://kubernetes.io/docs/
- GitHub Actions: https://docs.github.com/en/actions
- Go Binary Distribution: https://golang.org/doc/install/source

---

Last Updated: 2026-05-29
