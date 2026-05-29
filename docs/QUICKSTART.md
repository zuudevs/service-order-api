# Quick Start Guide

Get the Service Order API up and running in minutes!

## 5-Minute Setup

### Step 1: Clone the Repository

```bash
git clone https://github.com/zuudevs/service-order-api.git
cd service-order-api
```

### Step 2: Configure Environment

```bash
# Copy the example environment file
cp .env.example .env

# Edit .env with your settings (optional for quick test)
# At minimum, set INTERNAL_API_TOKEN_HASH
```

### Step 3: Run the Server

```bash
go run ./cmd/server
```

You should see:

```
server running on :8080
```

### Step 4: Test the API

In a new terminal:

```bash
# Check server health
curl http://localhost:8080/health

# Response: ok
```

SUCCESS! The API is running!

---

## Making Your First API Call

### Create a Person

```bash
curl -X POST http://localhost:8080/persons \
  -H "Content-Type: application/json" \
  -d '{
    "firstname": "John",
    "middlename": "Michael",
    "lastname": "Doe"
  }'
```

Response:

```json
{
  "id": 1,
  "firstname": "John",
  "middlename": "Michael",
  "lastname": "Doe",
  "created_at": "2026-05-29T11:50:32Z"
}
```

### List All Persons

```bash
curl http://localhost:8080/persons
```

### Get a Specific Person

```bash
curl http://localhost:8080/persons/1
```

### Update a Person (PATCH)

```bash
curl -X PATCH http://localhost:8080/persons/1 \
  -H "Content-Type: application/json" \
  -d '{
    "lastname": "Smith"
  }'
```

### Delete a Person

```bash
curl -X DELETE http://localhost:8080/persons/1
```

---

## Working with Related Resources

### Create an Order for a Person

First, ensure you have a person (e.g., ID: 1):

```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "person_id": 1,
    "order_date": "2026-05-29T11:50:32Z",
    "status": 0,
    "total_price": 10000
  }'
```

### Create a Contact for a Person

```bash
curl -X POST http://localhost:8080/contacts \
  -H "Content-Type: application/json" \
  -d '{
    "person_id": 1,
    "contact_type": "email",
    "contact_value": "john@example.com"
  }'
```

### Create a Task

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "task_name": "Install service",
    "task_description": "Install the service at customer location"
  }'
```

---

## Development Workflow

### 1. Make Code Changes

Edit files in internal/ directory

### 2. Rebuild and Test

```bash
# Option 1: Run directly (auto-rebuilds)
go run ./cmd/server

# Option 2: Build then run
go build -o service-order-api ./cmd/server
./service-order-api
```

### 3. Test Changes

Use curl or your preferred REST client to test endpoints

### 4. Check Database

The SQLite database is stored in storage/database.db:

```bash
# Using sqlite3 (if installed)
sqlite3 storage/database.db

# List tables
.tables

# View persons table
SELECT * FROM persons;
```

---

## Using Docker

### Build Docker Image

```bash
docker build -t service-order-api:latest .
```

### Run Docker Container

```bash
docker run -p 8080:8080 \
  -e PORT=8080 \
  -e INTERNAL_API_TOKEN_HASH=your_token_hash \
  service-order-api:latest
```

### Docker Compose

Create docker-compose.yml:

```yaml
version: "3.8"

services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      PORT: 8080
      INTERNAL_API_TOKEN_HASH: ${INTERNAL_API_TOKEN_HASH}
    volumes:
      - ./storage:/app/storage
```

Run it:

```bash
docker-compose up
```

---

## Environment Variables

Essential variables for quick start:

```env
# Server port (default: 8080)
PORT=8080

# Required: Internal API authentication token hash
INTERNAL_API_TOKEN_HASH=your_api_token_hash

# Optional: Google Drive backup
GOOGLE_DRIVE_CREDENTIALS=your_credentials_json
GOOGLE_DRIVE_TOKEN=your_access_token
GOOGLE_DRIVE_BACKUP_FOLDER_ID=your_folder_id
GOOGLE_DRIVE_DB_FILE_ID=your_db_file_id
```

---

## Database Initialization

### Automatic

The database is automatically created on first run.

### Manual

```bash
# Using the generate script (Windows PowerShell)
.\scripts\generate-database.ps1

# Or run the server (it initializes on startup)
go run ./cmd/server
```

### View Database Schema

```bash
cat internal/database/schema.sql
```

---

## Testing API Endpoints

### Using curl

Create:

```bash
curl -X POST http://localhost:8080/persons \
  -H "Content-Type: application/json" \
  -d '{"firstname":"Jane"}'
```

Read:

```bash
curl http://localhost:8080/persons/1
```

Update:

```bash
curl -X PATCH http://localhost:8080/persons/1 \
  -H "Content-Type: application/json" \
  -d '{"lastname":"Smith"}'
```

Delete:

```bash
curl -X DELETE http://localhost:8080/persons/1
```

### Using Postman

1. Import endpoints from /docs/API.md
2. Set base URL: http://localhost:8080
3. Add header: Authorization: Bearer YOUR_TOKEN
4. Test each endpoint

### Using VS Code REST Client

Create test.http:

```http
### Health Check
GET http://localhost:8080/health

### Create Person
POST http://localhost:8080/persons
Content-Type: application/json

{
  "firstname": "John",
  "lastname": "Doe"
}

### List Persons
GET http://localhost:8080/persons

### Get Person
GET http://localhost:8080/persons/1
```

Install the REST Client extension and click "Send Request".

---

## Common Tasks

### Add a New API Endpoint

1. Create handler method in internal/handlers/
2. Register route in internal/routes/routes.go
3. Test with curl/Postman

### Add a Service Feature

1. Implement logic in internal/services/
2. Call from handler or existing service
3. Add repository method if needed

### Modify Database Schema

1. Edit internal/database/schema.sql
2. Delete storage/database.db to reset
3. Restart server to regenerate database

---

## Troubleshooting

### Server will not start

```bash
# Check port is not in use
netstat -ano | findstr :8080  # Windows

# Or use different port
PORT=3000 go run ./cmd/server
```

### Database errors

```bash
# Clear database and reinitialize
rm storage/database.db
go run ./cmd/server
```

### API returns 401 Unauthorized

```bash
# Ensure INTERNAL_API_TOKEN_HASH is set in .env
# Or disable auth middleware for testing
```

### Module not found error

```bash
# Download dependencies
go mod download

# Or tidy modules
go mod tidy
```

---

## Next Steps

1. [+] Server is running
2. [+] Made first API call
3. Now explore:
   - Read API.md for detailed endpoint documentation
   - Read ARCHITECTURE.md to understand the codebase
   - Read BUILD.md for deployment options
   - Explore source code in internal/

---

## Useful Links

- API Documentation: API.md
- Architecture Guide: ARCHITECTURE.md
- Build Guide: BUILD.md
- GitHub: https://github.com/zuudevs/service-order-api
- Go Documentation: https://golang.org/doc/
- Chi Router Docs: https://go-chi.io/

---

## Getting Help

- Check logs in terminal for error messages
- Review .env.example for all available configuration
- Consult API.md for endpoint documentation
- Email: zuudevs@gmail.com

---

Ready to dive deeper? Check out the API Documentation (API.md) for complete endpoint reference!
