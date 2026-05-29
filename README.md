# Service Order API

A robust, RESTful API service for managing service orders, persons, contacts, tasks, and transactions. Built with Go using the Chi router framework and SQLite database, with automatic Google Drive backup capabilities.

## Overview

Service Order API is a production-ready backend service designed to manage service order workflows. It provides endpoints for:

- Person Management: Create and manage customer/service provider profiles
- Contact Management: Store and manage contact information linked to persons
- Order Management: Create and track service orders with status tracking
- Task Management: Create and assign tasks related to service orders
- Transaction Management: Record and track financial transactions for orders
- Detail Tasks: Track detailed task information and progress

### Key Features

- Built with Go & Chi Router - High-performance HTTP routing framework
- SQLite Database - Lightweight, file-based relational database
- Internal API Token Authentication - Secure internal API access
- Google Drive Backup - Automatic database backup to Google Drive every 6 hours
- Health Check Endpoint - Monitor API availability
- RESTful API Design - Standard HTTP methods (GET, POST, PUT, PATCH, DELETE)

## Table of Contents

- [Quick Start](#quick-start)
- [Requirements](#requirements)
- [Installation](#installation)
- [Configuration](#configuration)
- [Building](#building)
- [Running](#running)
- [API Documentation](#api-documentation)
- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [License](#license)

## Quick Start

### Prerequisites

- Go 1.25.8 or later
- SQLite (included in dependencies)
- (Optional) Google Drive credentials for backup feature

### Installation

```bash
# Clone the repository
git clone https://github.com/zuudevs/service-order-api.git
cd service-order-api

# Install dependencies
go mod download

# Copy environment variables template
cp .env.example .env

# Configure your environment (see Configuration section)
# Edit .env with your settings
```

### Running the Server

```bash
# Build and run
go run ./cmd/server

# Or build then run
go build -o service-order-api ./cmd/server
./service-order-api
```

The server will start on `http://localhost:8080` (or the port specified in `PORT` env variable).

### Verify Server is Running

```bash
curl http://localhost:8080/health
# Response: ok
```

## Configuration

Create a `.env` file in the project root with the following variables:

```env
# Server Configuration
PORT=8080

# Authentication
INTERNAL_API_TOKEN_HASH=YOUR_INTERNAL_API_TOKEN_HASH

# Google Drive Configuration (optional)
GOOGLE_DRIVE_CREDENTIALS=YOUR_GOOGLE_DRIVE_CREDENTIALS
GOOGLE_DRIVE_TOKEN=YOUR_GOOGLE_DRIVE_TOKEN
GOOGLE_DRIVE_BACKUP_FOLDER_ID=YOUR_GOOGLE_DRIVE_BACKUP_FOLDER_ID
GOOGLE_DRIVE_DB_FILE_ID=YOUR_GOOGLE_DRIVE_DB_FILE_ID
```

### Configuration Details

| Variable                        | Description                                        | Required | Default |
| ------------------------------- | -------------------------------------------------- | -------- | ------- |
| `PORT`                          | HTTP server port                                   | No       | 8080    |
| `INTERNAL_API_TOKEN_HASH`       | Hash of your internal API token for authentication | Yes      | -       |
| `GOOGLE_DRIVE_CREDENTIALS`      | Google Drive OAuth2 credentials JSON               | No       | -       |
| `GOOGLE_DRIVE_TOKEN`            | Google Drive access token                          | No       | -       |
| `GOOGLE_DRIVE_BACKUP_FOLDER_ID` | Google Drive folder ID for backups                 | No       | -       |
| `GOOGLE_DRIVE_DB_FILE_ID`       | Google Drive file ID for database backup           | No       | -       |

## Building

### Build the Project

```bash
# Build executable for current OS
go build -o service-order-api ./cmd/server

# Build for specific OS (Linux example)
GOOS=linux GOARCH=amd64 go build -o service-order-api ./cmd/server
```

### Build with Version Information

```bash
VERSION=1.0.0
go build -ldflags "-X main.Version=$VERSION" -o service-order-api ./cmd/server
```

## Project Structure

```
service-order-api/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── database/
│   │   ├── schema.sql             # Database schema definitions
│   │   └── sqlite.go              # SQLite connection and initialization
│   ├── handlers/                   # HTTP request handlers
│   │   ├── person_handler.go
│   │   ├── contact_handler.go
│   │   ├── order_handler.go
│   │   ├── task_handler.go
│   │   ├── transaction_handler.go
│   │   └── detail_task_handler.go
│   ├── middlewares/                # HTTP middleware
│   │   └── auth_middleware.go      # Token authentication
│   ├── models/                     # Data models
│   │   ├── person.go
│   │   ├── contact.go
│   │   ├── order.go
│   │   ├── task.go
│   │   ├── transaction.go
│   │   └── detail_task.go
│   ├── repositories/               # Data access layer
│   │   ├── person_repository.go
│   │   ├── contact_repository.go
│   │   └── ...
│   ├── services/                   # Business logic layer
│   │   ├── person_service.go
│   │   ├── contact_service.go
│   │   └── ...
│   └── routes/
│       └── routes.go              # Route registration
├── pkg/                            # Public packages (if any)
├── storage/
│   └── database.db                # SQLite database file (generated)
├── configs/
│   └── google/                    # Google Drive configuration
├── scripts/
│   ├── set-env.ps1                # PowerShell set environtment script
│   └── generate-database.ps1      # Database generation script
├── .env.example                   # Environment variables template
├── go.mod                         # Go module dependencies
├── go.sum                         # Go module checksums
└── README.md                      # This file
```

## API Documentation

All API endpoints require authentication via the `INTERNAL_API_TOKEN_HASH` header.

### Authentication

Include the following header in all API requests:

```
Authorization: Bearer YOUR_API_TOKEN
```

### Base URL

```
http://localhost:8080
```

### Health Check

**GET** `/health`

- Returns server status
- **Response**: `ok`

### Persons

| Method | Endpoint        | Description                  |
| ------ | --------------- | ---------------------------- |
| POST   | `/persons`      | Create a new person          |
| GET    | `/persons`      | List all persons             |
| GET    | `/persons/{id}` | Get person by ID             |
| PUT    | `/persons/{id}` | Replace entire person record |
| PATCH  | `/persons/{id}` | Update specific fields       |
| DELETE | `/persons/{id}` | Delete a person              |

### Contacts

| Method | Endpoint         | Description            |
| ------ | ---------------- | ---------------------- |
| POST   | `/contacts`      | Create a new contact   |
| GET    | `/contacts`      | List all contacts      |
| GET    | `/contacts/{id}` | Get contact by ID      |
| PUT    | `/contacts/{id}` | Replace entire contact |
| PATCH  | `/contacts/{id}` | Update specific fields |
| DELETE | `/contacts/{id}` | Delete a contact       |

### Orders

| Method | Endpoint       | Description            |
| ------ | -------------- | ---------------------- |
| POST   | `/orders`      | Create a new order     |
| GET    | `/orders`      | List all orders        |
| GET    | `/orders/{id}` | Get order by ID        |
| PUT    | `/orders/{id}` | Replace entire order   |
| PATCH  | `/orders/{id}` | Update specific fields |
| DELETE | `/orders/{id}` | Delete an order        |

### Tasks

| Method | Endpoint      | Description            |
| ------ | ------------- | ---------------------- |
| POST   | `/tasks`      | Create a new task      |
| GET    | `/tasks`      | List all tasks         |
| GET    | `/tasks/{id}` | Get task by ID         |
| PUT    | `/tasks/{id}` | Replace entire task    |
| PATCH  | `/tasks/{id}` | Update specific fields |
| DELETE | `/tasks/{id}` | Delete a task          |

### Transactions

| Method | Endpoint             | Description                |
| ------ | -------------------- | -------------------------- |
| POST   | `/transactions`      | Create a new transaction   |
| GET    | `/transactions`      | List all transactions      |
| GET    | `/transactions/{id}` | Get transaction by ID      |
| PUT    | `/transactions/{id}` | Replace entire transaction |
| PATCH  | `/transactions/{id}` | Update specific fields     |
| DELETE | `/transactions/{id}` | Delete a transaction       |

### Detail Tasks

| Method | Endpoint             | Description                |
| ------ | -------------------- | -------------------------- |
| POST   | `/detail-tasks`      | Create a new detail task   |
| GET    | `/detail-tasks`      | List all detail tasks      |
| GET    | `/detail-tasks/{id}` | Get detail task by ID      |
| PUT    | `/detail-tasks/{id}` | Replace entire detail task |
| PATCH  | `/detail-tasks/{id}` | Update specific fields     |
| DELETE | `/detail-tasks/{id}` | Delete a detail task       |

See [API.md](./docs/API.md) for detailed request/response examples.

## Architecture

The application follows a **layered architecture** pattern:

```
┌─────────────────────────────┐
│    HTTP Request/Response    │
│    (REST Endpoints)         │
└──────────────┬──────────────┘
               │
┌──────────────▼──────────────┐
│    Route Layer              │
│    (routes.go)              │
└──────────────┬──────────────┘
               │
┌──────────────▼──────────────┐
│    Middleware Layer         │
│    (Authentication, etc)    │
└──────────────┬──────────────┘
               │
┌──────────────▼──────────────┐
│    Handler Layer            │
│    (Business Logic Entry)   │
└──────────────┬──────────────┘
               │
┌──────────────▼──────────────┐
│    Service Layer            │
│    (Business Logic)         │
└──────────────┬──────────────┘
               │
┌──────────────▼──────────────┐
│    Repository Layer         │
│    (Data Access)            │
└──────────────┬──────────────┘
               │
┌──────────────▼──────────────┐
│    Database Layer           │
│    (SQLite)                 │
└─────────────────────────────┘
```

### Design Patterns

- **Repository Pattern**: Abstracts data access logic
- **Service Layer**: Encapsulates business logic
- **Dependency Injection**: Controllers receive dependencies via constructors
- **Middleware Pattern**: Cross-cutting concerns (authentication)
- **RESTful API**: Standard HTTP conventions

See [ARCHITECTURE.md](./docs/ARCHITECTURE.md) for detailed architecture documentation.

## Dependencies

Key dependencies used in this project:

- **[Chi Router](https://github.com/go-chi/chi)** v5.3.0 - Lightweight HTTP router
- **[SQLite](https://modernc.org/sqlite)** v1.39.1 - Embedded SQL database
- **[Google API Client](https://google.golang.org/api)** v0.281.0 - Google Drive integration
- **[OpenTelemetry](https://opentelemetry.io/)** - Tracing and monitoring

Full dependency list available in [go.mod](./go.mod).

## Development

### Running in Development Mode

```bash
go run ./cmd/server
```

### Database Initialization

The database is automatically initialized on first run. To manually generate the database:

```bash
# Using PowerShell (Windows)
.\scripts\generate-database.ps1

# Using Go directly
go run ./cmd/server
```

### Google Drive Backup

The application includes automatic backup to Google Drive:

1. Configure Google Drive credentials in `.env`
2. First upload: Set `GOOGLE_DRIVE_DB_FILE_ID=""` to upload new backup
3. Subsequent uploads: Automatic every 6 hours
4. After first upload, save the returned file ID to `GOOGLE_DRIVE_DB_FILE_ID`

## Security

- **Authentication**: Internal API token validation via middleware
- **Database**: SQLite with proper SQL parameterization to prevent injection
- **HTTPS**: Recommended for production deployments
- **Token Management**: Use strong, randomly generated tokens

## License

Copyright (c) 2026 ZUU Devs. Licensed under the [MIT License](./LICENSE).

## Contributing

Contributions are welcome! Please follow the existing code style and patterns.

## Support

For issues, questions, or suggestions:

- Email: zuudevs@gmail.com
- GitHub: https://github.com/zuudevs/service-order-api

---

**Version**: 0.1.0  
**Last Updated**: 2026-05-29
