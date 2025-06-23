# URL Shortener Service Makefile

.PHONY: help init run build test clean migrate-up migrate-down code-gen

# Default target
help:
	@echo "URL Shortener Service Commands:"
	@echo "  init         - Initialize service (deps, codegen, migrations, install air)"
	@echo "  run          - Run service in development mode"
	@echo "  build        - Build the service binary"
	@echo "  test         - Run all tests"
	@echo "  test-cov     - Run tests with coverage"
	@echo "  clean        - Clean build artifacts"
	@echo ""
	@echo "Database Commands:"
	@echo "  migrate-up   - Run database migrations"
	@echo "  migrate-down - Rollback migrations"
	@echo "  migrate-new  - Create new migration"
	@echo "  code-gen     - Generate Ent code"
	@echo ""
	@echo "Development:"
	@echo "  air          - Run with hot reloading (after make init)"
	@echo "  go-lint      - Run linter"

# Initialize service (install dependencies, generate code, run migrations)
init:
	@echo "🔧 Initializing URL Shortener Service..."
	@go mod tidy
	@echo "📝 Generating database code..."
	@make code-gen
	@echo "🗄️  Running database migrations..."
	@make migrate-up
	@echo "🔥 Installing Air for hot reloading..."
	@go install github.com/air-verse/air@latest
	@echo "✅ Service initialized successfully!"
	@echo ""
	@echo "Next step: Run 'air' for hot reloading development"

# Run service
run:
	@echo "🚀 Starting URL Shortener Service..."
	@go run cmd/app/main.go

# Build service binary
build:
	@echo "🔨 Building URL Shortener Service..."
	@go build -o bin/url-shortener cmd/app/main.go
	@echo "✅ Build completed: bin/url-shortener"

# Generate Ent database code
code-gen:
	@go generate cmd/ent/generate.go
	@echo "🎉 Database code generated successfully!"

# generate a blank migration
migrate-gen: code-gen
	@read -p "Enter migration name (e.g., create-todo-table): " name; \
			go run cmd/migrate/main.go ent_add $$name

# generate a blank migration
migrate-new:
	@read -p "Enter migration name (e.g., create-todo-table): " name; \
		go run cmd/migrate/main.go add $$name

# deploy migration
migrate-up:
	go run cmd/migrate/main.go up

# rollback migration
migrate-down:
	@read -p "Enter migration version (e.g., 20250304025256): " version; \
		go run cmd/migrate/main.go down $$version

# re-generate migration directory integrity hash 
migrate-hash:
	go run cmd/migrate/main.go hash

# check migration status
migrate-status:
	go run cmd/migrate/main.go statusa

SHELL := /bin/bash
CI_DOCKER_IMAGE=github
service=user

# build images from all services
.PHONY: build-image-local
build-image-local: build-linux
	docker build --platform linux/arm64 -f build/Dockerfile.dev -t user_service:latest . --build-arg service=${service}

.PHONY: build-linux
build-linux: download-go-mod
	GOOS=linux GOARCH=amd64 go build -mod=readonly -o out/cmd/main cmd/app/main.go

.PHONY: download-go-mod
download-go-mod: go.mod
	go mod download all
