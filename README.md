# TaskFlow

A RESTful API for project and task management with real-time WebSocket support, built with Go, Gin, and PostgreSQL.

## Features

- **User Authentication** - JWT-based authentication with registration and login
- **Project Management** - Create and manage projects with members
- **Task Management** - Create, update, and track tasks within projects
- **Real-time Updates** - WebSocket support for real-time task updates
- **Role-based Access** - Owner and member roles for projects
- **Docker Support** - Easy deployment with Docker Compose

## Tech Stack

- **Language**: Go 1.25
- **Web Framework**: Gin
- **Database**: PostgreSQL 16
- **Authentication**: JWT (JSON Web Tokens)
- **Migration Tool**: go-migrate
- **Containerization**: Docker & Docker Compose

## Project Structure

```
taskflow/
├── cmd/
│   └── api/
│       └── main.go          # Application entry point
├── docs/
│   └── API.md               # API documentation
├── internal/
│   ├── domain/              # Domain models and interfaces
│   │   ├── project.go
│   │   ├── task.go
│   │   └── user.go
│   ├── handler/             # HTTP request handlers
│   ├── middleware/          # Gin middleware (CORS, auth)
│   ├── repository/          # Database access layer
│   ├── service/             # Business logic layer
│   └── websocket/           # WebSocket hub and handlers
├── migrations/              # Database migrations
├── pkg/
│   ├── config/              # Configuration management
│   ├── database/            # Database connection
│   ├── response/            # Response wrappers
│   └── token/               # JWT token utilities
├── .env                     # Environment variables (not included)
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── Makefile
```

## Getting Started

### Prerequisites

- Go 1.25+
- Docker & Docker Compose (optional)
- PostgreSQL 16 (optional, if not using Docker)

### Environment Setup

1. Clone the repository:
```bash
git clone <repository-url>
cd taskflow
```

2. Create a `.env` file in the root directory:
```env
APP_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=taskflow
JWT_SECRET=your_jwt_secret_key_here
```

### Running Locally

1. Install dependencies:
```bash
go mod download
```

2. Run database migrations:
```bash
make migrate-up
```

3. Start the development server:
```bash
make run
```

The server will start on `http://localhost:8080`

### Running with Docker

1. Start all services:
```bash
make docker-up
```

2. View logs:
```bash
make docker-logs
```

3. Stop services:
```bash
make docker-down
```

## API Documentation

See [docs/API.md](docs/API.md) for complete API endpoint documentation.

### Quick Overview

#### Public Endpoints
- `GET /health` - Health check
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login and get JWT token

#### Protected Endpoints (require JWT token)
- `GET /api/me` - Get current user
- `POST /api/projects` - Create project
- `GET /api/projects` - Get user's projects
- `GET /api/projects/:id` - Get project details
- `POST /api/projects/:id/members` - Add project member
- `POST /api/projects/:id/tasks` - Create task
- `PATCH /api/tasks/:task_id` - Update task
- `DELETE /api/tasks/:task_id` - Delete task

## Makefile Commands

| Command | Description |
|---------|-------------|
| `make run` | Run the application with air (hot reload) |
| `make build` | Build the binary |
| `make test` | Run tests with coverage |
| `make migrate-up` | Run database migrations up |
| `make migrate-down` | Rollback database migrations |
| `make docker-up` | Start Docker Compose services |
| `make docker-down` | Stop Docker Compose services |
| `make docker-logs` | View Docker logs |

## Database Migrations

Migrations are located in the `migrations/` directory. To create a new migration:

```bash
migrate create -ext sql -dir migrations -seq <migration_name>
```

## Task Status Values

Tasks can have the following statuses:
- `todo` - Task is pending
- `in_progress` - Task is being worked on
- `done` - Task is completed
