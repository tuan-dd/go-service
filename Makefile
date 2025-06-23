# dev:
# 	go run ./cmd/market-data-service
# sql:
# 	sqlc generate

# # example run: make new_migration name=00001_init_schema
# new_migration:
# 	migrate create -ext sql -dir sql/schema -seq $(name)

# migrate_up:
# 	migrate -path sql/schema -database "postgresql://postgres:Kiloma123@@localhost:5432/Crypto?sslmode=disable" -verbose up

# migrate_down:
# 	migrate -path sql/schema -database "postgresql://postgres:Kiloma123@@localhost:5432/Crypto?sslmode=disable" -verbose down


# Global Makefile for Go Microservices Project

.PHONY: run-all stop-all clean test-all help

# Start all services and infrastructure
run-all:
	@echo "🚀 Starting all services..."
	@docker compose -f docker-compose.dev.yml up -d
	@echo "✅ All services started successfully!"
	@echo "📋 Services:"
	@echo "   - PostgreSQL: localhost:5432"
	@echo "   - Redis: localhost:6379" 
	@echo "   - URL Service will be available at: localhost:8080"
	@echo ""
	@echo "Next steps:"
	@echo "   cd services/url"
	@echo "   make init"
	@echo "   air"

# Stop all services
stop-all:
	@echo "🛑 Stopping all services..."
	@docker compose -f docker-compose.dev.yml down
	@echo "✅ All services stopped!"

# Clean up containers and volumes
clean:
	@echo "🧹 Cleaning up containers and volumes..."
	@docker compose -f docker-compose.dev.yml down -v --remove-orphans
	@docker system prune -f
	@echo "✅ Cleanup completed!"

# Run tests for all services
test-all:
	@echo "🧪 Running tests for all services..."
	@cd services/url && make test

# Show help
help:
	@echo "Available commands:"
	@echo "  run-all    - Start all services and infrastructure"
	@echo "  stop-all   - Stop all services"
	@echo "  clean      - Clean up containers and volumes"
	@echo "  test-all   - Run tests for all services"
	@echo "  help       - Show this help message"
	@echo ""
	@echo "Service-specific commands:"
	@echo "  cd services/url && make help"

# Lint Go code (existing functionality)
go-lint:
	@for dir in $(DIRS); do \
		if [ "$$(basename $$dir)" != "proto" ]; then \
			echo "Running golangci-lint in $$dir"; \
			cd $$dir && golangci-lint run && cd ..; \
		else \
			echo "Skipping $$dir"; \
		fi \
	done

# Legacy variables and directories
WORKER_IMAGE=1.24.0-alpine3.21
DIRS := $(shell find service-internal -mindepth 1 -maxdepth 1 -type d ! -name 'no_get' 2>/dev/null || echo "")
DONT_STOP := db redis