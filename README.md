# Task Management Service- Submission- Gunish

A microservice-based Task Management System built with Go, implementing a clean architecture pattern with clear separation of concerns.
The system provides a RESTful API for managing tasks with features like pagination, filtering, and event-driven architecture.

## Features

- **CRUD Operations**: Create, Read, Update, and Delete tasks
- **Pagination**: Efficient task listing with pagination support
- **Filtering**: Filter tasks by status (pending, in_progress, completed)
- **Event-Driven Architecture**: Task events are published to Kafka for asynchronous processing, can be consumed by other services.
- **Clean Architecture**: Clear separation of concerns with domain-driven design
- **Docker Support**: Containerized application with Docker and Docker Compose
- **SQLite Database**: Lightweight database for task storage, initially started with in-memory implementation for MVP
- **RESTful API**: Well-defined API endpoints following REST principles 
- Graceful handling(closing) of resources like db, server and kafka

## Architecture

The system follows a clean architecture pattern with the following layers:

1. **API Layer** (`internal/api/`): Handles HTTP requests and responses
2. **Service Layer** (`internal/service/`): Contains business logic
3. **Repository Layer** (`internal/domain/repository/`): Manages data persistence
4. **Domain Layer** (`internal/domain/`): Contains core business models and interfaces
5. **Common Layer** (`internal/common/`): Shared utilities and configurations (logging,middlewares)

### Microservices Design

The system is designed as a microservice with the following characteristics:

- **Single Responsibility**: Each component has a specific responsibility
- **Event-Driven Communication**: Uses Kafka for event publishing and consumption
- **Scalability**: Can be horizontally scaled by running multiple instances
- **Independent Deployment**: Each component can be deployed independently

## API Documentation

### Endpoints

#### Create Task
```http
POST /tasks
Content-Type: application/json

{
    "title": "Task Title",
    "description": "Task Description",
    "due_date": "2024-03-20T00:00:00Z"
}
```

#### Get Task
```http
GET /tasks/{id}
```

#### Update Task
```http
PUT /tasks/{id}
Content-Type: application/json

{
    "title": "Updated Title",
    "description": "Updated Description",
    "status": "in_progress",
    "due_date": "2024-03-21T00:00:00Z"
}
```

#### Delete Task
```http
DELETE /tasks/{id}
```

#### List Tasks
```http
GET /tasks?status=completed&page=1&page_size=10
```

Query Parameters:
- `status`: Filter by status (pending, in_progress, completed)
- `page`: Page number (default: 1)
- `page_size`: Items per page (default: 10, max: 100)

## Getting Started

### Prerequisites

- Go 1.24 or higher
- Docker and Docker Compose
- Make (optional, for using Makefile commands)

### Running Locally

1. Clone the repository:
```bash
git clone https://github.com/pvnptl/alle-task-manager-gunish.git
cd alle-task-manager-gunish
```

2. Install dependencies:
```bash
go mod download
```

3. Run the application:
```bash
go run main.go
```

### Running with Docker

1. Build and start the containers:
```bash
docker-compose up --build
```

The application will be available at `http://localhost:8080`

## Configuration

The application can be configured using environment variables:

- `SERVER_PORT`: Server port (default: 8080)
- `SERVER_READ_TIMEOUT`: Read timeout in seconds (default: 10)
- `SERVER_WRITE_TIMEOUT`: Write timeout in seconds (default: 10)
- `DB_DRIVER`: Database driver (default: sqlite)
- `SQLITE_DB_PATH`: SQLite database path (default: tasks.db)
- `KAFKA_BROKERS`: Kafka broker addresses (default: localhost:9092)
- `KAFKA_TOPIC`: Kafka topic for task events (default: task-events)
- `KAFKA_GROUP_ID`: Kafka consumer group ID (default: task-management-group)

## Testing

Run the tests:
```bash
 cd internal && go test -v ./...
```

## Project Structure

```
.
├── cmd/                    # Application entry points
├── internal/              # Private application code
│   ├── api/              # API handlers and middleware
│   ├── common/           # Shared utilities and configurations
│   ├── domain/           # Core business models and interfaces
│   └── service/          # Business logic implementation
├── pkg/                  # Public libraries
├── Dockerfile           # Docker build configuration
├── docker-compose.yaml  # Docker Compose configuration
├── go.mod              # Go module definition
└── go.sum              # Go module dependencies checksum
```

Design Patterns Used:
* Singleton pattern in logger.go



