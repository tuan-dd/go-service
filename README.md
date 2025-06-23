# Go Microservices Architecture

A scalable microservices architecture built with Go, featuring URL shortener service and shared packages.

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- Make
- Air (for hot reloading)

### Run All Services

```bash
# Start all infrastructure and services
make run-all
```

### Setup URL Shortener Service

```bash
# Navigate to URL service
cd services/url

# Initialize service (install deps, run migrations)
make init

# Run with hot reloading
air
```

### Test the API

1. **Import Postman Collection**
   - Open Postman
   - Click **Import** → Select `interview.postman_collection.json`
   - All API endpoints will be available

2. **Test Endpoints**
   ```bash
   # Health check
   curl http://localhost:8080/health
   
   # Create short URL
   curl -X POST http://localhost:8080/shortURLs/ \
     -H "Content-Type: application/json" \
     -d '{"url": "https://google.com"}'
   
   # Access short URL (redirects)
   curl -L http://localhost:8080/{short_code}/redirect
   ```

## 📁 Project Structure

```
go-service/
├── pkg/                          # Shared packages
│   ├── appLogger/               # Logging utilities
│   ├── caching/                 # Cache implementations (Redis, Memory)
│   ├── common/                  # Common utilities and types
│   ├── database/                # Database connections (MySQL, PostgreSQL)
│   ├── fiberzap-logger/         # Fiber + Zap integration
│   ├── grpc-pkg/                # gRPC utilities
│   ├── http-client/             # HTTP client utilities
│   ├── orm/ent/                 # Ent ORM mixins and utilities
│   ├── queues/                  # Queue implementations (NATS, RabbitMQ, Asynq)
│   └── settings/                # Configuration management
├── services/
│   └── url/                     # URL Shortener Service
│       ├── cmd/app/             # Application entry point
│       ├── internal/            # Private application code
│       ├── tests/               # Test files
│       ├── configs/             # Configuration files
│       ├── migrations/          # Database migrations
│       ├── Makefile            # Service-specific commands
│       └── README.md           # Service documentation
├── scripts/                     # Build and deployment scripts
├── docker-compose.dev.yml      # Development environment
├── go.work                     # Go workspace
├── Makefile                    # Global commands
├── interview.postman_collection.json  # API testing collection
└── README.md                   # This file
```

## 🏗️ Architecture

### Microservices
- **URL Shortener**: High-performance URL shortening service with caching

### Shared Packages
- **Logging**: Structured logging with Zap
- **Caching**: Redis and in-memory cache implementations
- **Database**: PostgreSQL and MySQL support with connection pooling
- **Messaging**: NATS, RabbitMQ, and Asynq queue implementations
- **HTTP**: Fiber-based HTTP utilities and middleware
- **gRPC**: gRPC server and client utilities
- **ORM**: Ent ORM with soft delete and audit mixins

## 🛠️ Development

### Available Commands

```bash
# Global commands (from root)
make run-all          # Start all services
make stop-all         # Stop all services
make clean           # Clean up containers and volumes
make test-all        # Run tests for all services

# Service-specific commands (from services/url)
make init            # Initialize service
make run             # Run service
make test            # Run tests
make build           # Build binary
make migrate-up      # Run database migrations
make migrate-down    # Rollback migrations
```

### Development Workflow

1. **Start Infrastructure**
   ```bash
   make run-all
   ```

2. **Setup Service**
   ```bash
   cd services/url
   make init
   ```

3. **Develop with Hot Reload**
   ```bash
   air
   ```

4. **Run Tests**
   ```bash
   make test
   ```

## 🧪 Testing

### Import Postman Collection
1. Open Postman
2. Click **Import**
3. Select `interview.postman_collection.json`
4. Collection includes all API endpoints with examples

### Test Coverage
- Unit tests for all business logic
- Integration tests for API endpoints
- Performance tests for high-load scenarios

```bash
# Run all tests
cd services/url && make test

# Run with coverage
go test -v -cover ./tests/unit/

# Run specific tests
go test -v ./tests/unit/ -run TestUrlController
```

## 🚀 Deployment

### Docker Compose
```bash
# Development
docker-compose -f docker-compose.dev.yml up

# Production
docker-compose up -d
```
