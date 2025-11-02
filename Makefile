.PHONY: help build run test clean migrate-up migrate-down docker-up docker-down

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building application..."
	@go build -o bin/api cmd/api/main.go

run: ## Run the application
	@echo "Running application..."
	@go run cmd/api/main.go

dev: ## Run with hot reload (requires air)
	@echo "Running with hot reload..."
	@air

test: ## Run tests
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

test-coverage: test ## Run tests with coverage report
	@echo "Generating coverage report..."
	@go tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report generated at coverage.html"

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run ./...

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@goimports -w .

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

tidy: ## Tidy dependencies
	@echo "Tidying dependencies..."
	@go mod tidy

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf coverage.txt coverage.html

migrate-up: ## Run database migrations
	@echo "Running migrations..."
	@go run cmd/api/main.go migrate

migrate-create: ## Create a new migration file
	@read -p "Enter migration name: " name; \
	go run cmd/api/main.go migrate create $$name

docker-up: ## Start docker containers (development)
	@echo "Starting docker containers..."
	@docker-compose up -d

docker-up-staging: ## Start docker containers (staging)
	@echo "Starting staging docker containers..."
	@docker-compose -f docker-compose.staging.yml up -d

docker-down: ## Stop docker containers
	@echo "Stopping docker containers..."
	@docker-compose down

docker-down-staging: ## Stop staging docker containers
	@echo "Stopping staging docker containers..."
	@docker-compose -f docker-compose.staging.yml down

docker-logs: ## Show docker logs
	@docker-compose logs -f

docker-logs-staging: ## Show staging docker logs
	@docker-compose -f docker-compose.staging.yml logs -f

install-tools: ## Install development tools
	@echo "Installing tools..."
	@go install github.com/air-verse/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest

db-reset: ## Reset database (drop and recreate)
	@echo "Resetting database..."
	@PGPASSWORD=$(DB_PASSWORD) psql -h localhost -U salju -c "DROP DATABASE IF EXISTS bunely;"
	@PGPASSWORD=$(DB_PASSWORD) psql -h localhost -U salju -c "CREATE DATABASE bunely;"
	@echo "Database reset complete"

.DEFAULT_GOAL := help
