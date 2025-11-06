# Makefile for Go Fiber API Project

.PHONY: help migrate migrate-up migrate-down migrate-new migrate-apply migrate-status build run test clean deps lint format air docker-build docker-run docker-up docker-down docker-dev docker-logs docker-clean docker-ps docker-setup stop install-tools install-atlas

# Variables
BINARY_NAME=main
BINARY_PATH=./tmp/$(BINARY_NAME)
GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_RUN=$(GO_CMD) run
GO_TEST=$(GO_CMD) test
GO_MOD=$(GO_CMD) mod
GO_FMT=$(GO_CMD) fmt
GO_VET=$(GO_CMD) vet
MAIN_PATH=main.go
ENV_FILE?=.env.dev

# Colors for output
CYAN=\033[0;36m
GREEN=\033[0;32m
YELLOW=\033[0;33m
RED=\033[0;31m
NC=\033[0m # No Color

# Default target
.DEFAULT_GOAL := help

## help: Show this help message
help:
	@echo "$(CYAN)Available commands:$(NC)"
	@echo ""
	@echo "$(GREEN)Setup & Dependencies:$(NC)"
	@echo "  make deps          - Download and install dependencies"
	@echo "  make tidy          - Clean up dependencies"
	@echo "  make install-tools - Install required tools (air, etc.)"
	@echo "  make install-atlas - Install Atlas migration tool"
	@echo ""
	@echo "$(GREEN)Database Migrations (Atlas with GORM):$(NC)"
	@echo "  make migrate-diff   - Generate migrations from GORM models"
	@echo "  make migrate-apply  - Apply pending migrations"
	@echo "  make migrate-down   - Rollback last migration"
	@echo "  make migrate-status - Check migration status"
	@echo ""
	@echo "$(GREEN)Build & Run:$(NC)"
	@echo "  make build         - Build the application"
	@echo "  make run           - Run the application with go run"
	@echo "  make air           - Run the application with air (hot reload)"
	@echo "  make run-prod      - Run the built binary"
	@echo ""
	@echo "$(GREEN)Testing:$(NC)"
	@echo "  make test          - Run all tests"
	@echo "  make test-verbose  - Run tests with verbose output"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo ""
	@echo "$(GREEN)Code Quality:$(NC)"
	@echo "  make fmt           - Format code"
	@echo "  make vet           - Run go vet"
	@echo "  make lint          - Run linter (if installed)"
	@echo ""
	@echo "$(GREEN)Cleanup:$(NC)"
	@echo "  make clean         - Remove build artifacts"
	@echo "  make clean-all     - Remove all generated files and logs"
	@echo ""
	@echo "$(GREEN)Docker:$(NC)"
	@echo "  make docker-build     - Build Docker image"
	@echo "  make docker-run       - Run Docker container"
	@echo "  make docker-up        - Start services with docker-compose (production)"
	@echo "  make docker-dev       - Start development environment with docker-compose"
	@echo "  make docker-down      - Stop services with docker-compose"
	@echo "  make docker-down-dev  - Stop development services"
	@echo "  make docker-logs      - View docker-compose logs"
	@echo "  make docker-logs-dev  - View development docker-compose logs"
	@echo "  make docker-restart   - Restart services"
	@echo "  make docker-clean      - Remove containers, volumes, and networks"
	@echo "  make docker-ps         - Show running containers"
	@echo "  make docker-setup      - Setup with Docker (start services)"
	@echo ""
	@echo "$(GREEN)Utilities:$(NC)"
	@echo "  make check-env     - Check if .env.dev file exists"
	@echo "  make stop         - Stop running processes"
	@echo ""

## deps: Download and install dependencies
deps:
	@echo "$(CYAN)Installing dependencies...$(NC)"
	@$(GO_MOD) download
	@$(GO_MOD) tidy
	@echo "$(GREEN)✓ Dependencies installed$(NC)"

## tidy: Clean up dependencies
tidy:
	@echo "$(CYAN)Cleaning up dependencies...$(NC)"
	@$(GO_MOD) tidy
	@echo "$(GREEN)✓ Dependencies cleaned$(NC)"

## install-tools: Install required development tools
install-tools:
	@echo "$(CYAN)Installing development tools...$(NC)"
	@if ! command -v air > /dev/null; then \
		echo "$(YELLOW)Installing air...$(NC)"; \
		$(GO_CMD) install github.com/cosmtrek/air@latest; \
	else \
		echo "$(GREEN)✓ air already installed$(NC)"; \
	fi
	@echo "$(GREEN)✓ Tools installed$(NC)"

## install-atlas: Install Atlas migration tool
install-atlas:
	@echo "$(CYAN)Installing Atlas...$(NC)"
	@if command -v atlas > /dev/null; then \
		echo "$(GREEN)✓ Atlas already installed$(NC)"; \
	else \
		echo "$(YELLOW)Installing Atlas...$(NC)"; \
		curl -sSf https://atlasgo.sh | sh; \
	fi
	@echo "$(GREEN)✓ Atlas installed$(NC)"

## migrate-apply: Apply pending migrations using Atlas
migrate-apply: check-env
	@echo "$(CYAN)Applying migrations with Atlas...$(NC)"
	@if ! command -v atlas > /dev/null; then \
		echo "$(RED)✗ Atlas is not installed. Run 'make install-atlas' first$(NC)"; \
		exit 1; \
	fi
	@if [ -f "$(ENV_FILE)" ]; then \
		export $$(cat $(ENV_FILE) | grep -v '^#' | xargs); \
		atlas migrate apply --env local --var "db_url=postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=$$DB_SSLMODE"; \
	else \
		echo "$(YELLOW)⚠ Warning: $(ENV_FILE) not found, using default connection$(NC)"; \
		atlas migrate apply --env local --var "db_url=$(db_url)"; \
	fi
	@echo "$(GREEN)✓ Migrations applied$(NC)"

## migrate-down: Rollback last migration
migrate-down: check-env
	@echo "$(CYAN)Rolling back last migration...$(NC)"
	@if ! command -v atlas > /dev/null; then \
		echo "$(RED)✗ Atlas is not installed. Run 'make install-atlas' first$(NC)"; \
		exit 1; \
	fi
	@if [ -f "$(ENV_FILE)" ]; then \
		export $$(cat $(ENV_FILE) | grep -v '^#' | xargs); \
		atlas migrate down --env local --var "db_url=postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=$$DB_SSLMODE" 1; \
	else \
		echo "$(YELLOW)⚠ Warning: $(ENV_FILE) not found, using default connection$(NC)"; \
		atlas migrate down --env local 1; \
	fi
	@echo "$(GREEN)✓ Migration rolled back$(NC)"

## migrate-status: Check migration status
migrate-status: check-env
	@echo "$(CYAN)Checking migration status...$(NC)"
	@if ! command -v atlas > /dev/null; then \
		echo "$(RED)✗ Atlas is not installed. Run 'make install-atlas' first$(NC)"; \
		exit 1; \
	fi
	@if [ -f "$(ENV_FILE)" ]; then \
		export $$(cat $(ENV_FILE) | grep -v '^#' | xargs); \
		atlas migrate status --env local --var "db_url=postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=$$DB_SSLMODE"; \
	else \
		echo "$(YELLOW)⚠ Warning: $(ENV_FILE) not found, using default connection$(NC)"; \
		atlas migrate status --env local; \
	fi

## migrate-diff: Generate migrations from GORM models
migrate-diff: check-env
	@echo "$(CYAN)Generating migrations from GORM models...$(NC)"
	@if ! command -v atlas > /dev/null; then \
		echo "$(RED)✗ Atlas is not installed. Run 'make install-atlas' first$(NC)"; \
		exit 1; \
	fi
	@if [ -f "$(ENV_FILE)" ]; then \
		export $$(cat $(ENV_FILE) | grep -v '^#' | xargs); \
		atlas migrate diff --env gorm --var "db_url=postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=$$DB_SSLMODE"; \
	else \
		echo "$(YELLOW)⚠ Warning: $(ENV_FILE) not found, using default connection$(NC)"; \
		atlas migrate diff --env gorm; \
	fi
	@echo "$(GREEN)✓ Migrations generated from GORM models$(NC)"

## build: Build the application
build:
	@echo "$(CYAN)Building application...$(NC)"
	@mkdir -p tmp
	@$(GO_BUILD) -o $(BINARY_PATH).exe $(MAIN_PATH)
	@echo "$(GREEN)✓ Build complete: $(BINARY_PATH).exe$(NC)"

## build-linux: Build for Linux
build-linux:
	@echo "$(CYAN)Building for Linux...$(NC)"
	@mkdir -p tmp
	@GOOS=linux GOARCH=amd64 $(GO_BUILD) -o $(BINARY_PATH)-linux $(MAIN_PATH)
	@echo "$(GREEN)✓ Linux build complete$(NC)"

## build-windows: Build for Windows
build-windows:
	@echo "$(CYAN)Building for Windows...$(NC)"
	@mkdir -p tmp
	@GOOS=windows GOARCH=amd64 $(GO_BUILD) -o $(BINARY_PATH).exe $(MAIN_PATH)
	@echo "$(GREEN)✓ Windows build complete$(NC)"

## build-mac: Build for macOS
build-mac:
	@echo "$(CYAN)Building for macOS...$(NC)"
	@mkdir -p tmp
	@GOOS=darwin GOARCH=amd64 $(GO_BUILD) -o $(BINARY_PATH)-mac $(MAIN_PATH)
	@echo "$(GREEN)✓ macOS build complete$(NC)"

## run: Run the application
run:
	@echo "$(CYAN)Running application...$(NC)"
	@$(GO_RUN) $(MAIN_PATH)

## run-prod: Run the built binary
run-prod: build
	@echo "$(CYAN)Running production binary...$(NC)"
	@$(BINARY_PATH).exe

## air: Run with air (hot reload)
air:
	@echo "$(CYAN)Starting application with air (hot reload)...$(NC)"
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "$(RED)✗ air is not installed. Run 'make install-tools' first$(NC)"; \
		exit 1; \
	fi

## test: Run tests
test:
	@echo "$(CYAN)Running tests...$(NC)"
	@$(GO_TEST) -v ./...

## test-verbose: Run tests with verbose output
test-verbose:
	@echo "$(CYAN)Running tests with verbose output...$(NC)"
	@$(GO_TEST) -v -race ./...

## test-coverage: Run tests with coverage
test-coverage:
	@echo "$(CYAN)Running tests with coverage...$(NC)"
	@$(GO_TEST) -v -coverprofile=coverage.out ./...
	@$(GO_CMD) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage report generated: coverage.html$(NC)"

## fmt: Format code
fmt:
	@echo "$(CYAN)Formatting code...$(NC)"
	@$(GO_FMT) ./...
	@echo "$(GREEN)✓ Code formatted$(NC)"

## vet: Run go vet
vet:
	@echo "$(CYAN)Running go vet...$(NC)"
	@$(GO_VET) ./...
	@echo "$(GREEN)✓ Vet check complete$(NC)"

## lint: Run linter (requires golangci-lint)
lint:
	@echo "$(CYAN)Running linter...$(NC)"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "$(YELLOW)⚠ golangci-lint not installed. Install it from: https://golangci-lint.run/$(NC)"; \
		exit 1; \
	fi

## clean: Remove build artifacts
clean:
	@echo "$(CYAN)Cleaning build artifacts...$(NC)"
	@rm -rf tmp/
	@rm -f $(BINARY_NAME)
	@rm -f $(BINARY_NAME).exe
	@echo "$(GREEN)✓ Clean complete$(NC)"

## clean-all: Remove all generated files
clean-all: clean
	@echo "$(CYAN)Cleaning all generated files...$(NC)"
	@rm -f coverage.out coverage.html
	@rm -f *.log
	@rm -f air.log
	@rm -f build-errors.log
	@rm -f server.log server_error.log
	@echo "$(GREEN)✓ Deep clean complete$(NC)"

## check-env: Check if .env.dev file exists
check-env:
	@if [ ! -f $(ENV_FILE) ]; then \
		echo "$(RED)✗ $(ENV_FILE) file not found!$(NC)"; \
		echo "$(YELLOW)Copy .env.dev.example to .env.dev and update with your values$(NC)"; \
		exit 1; \
	else \
		echo "$(GREEN)✓ $(ENV_FILE) file exists$(NC)"; \
	fi

## stop: Stop running processes
stop:
	@echo "$(CYAN)Stopping running processes...$(NC)"
	@pkill -f "air" || true
	@pkill -f "main.exe" || true
	@pkill -f "main" || true
	@echo "$(GREEN)✓ Processes stopped$(NC)"

## docker-build: Build Docker image
docker-build:
	@echo "$(CYAN)Building Docker image...$(NC)"
	@docker build -t go-fiber-api:latest .
	@echo "$(GREEN)✓ Docker image built$(NC)"

## docker-run: Run Docker container
docker-run:
	@echo "$(CYAN)Running Docker container...$(NC)"
	@docker run -p 8080:8080 --env-file .env.dev go-fiber-api:latest

## docker-stop: Stop Docker container
docker-stop:
	@echo "$(CYAN)Stopping Docker container...$(NC)"
	@docker stop $$(docker ps -q --filter ancestor=go-fiber-api:latest) || true
	@echo "$(GREEN)✓ Docker container stopped$(NC)"

## docker-up: Start services with docker-compose (production)
docker-up:
	@echo "$(CYAN)Starting services with docker-compose...$(NC)"
	@docker-compose up -d
	@echo "$(GREEN)✓ Services started$(NC)"
	@echo "$(YELLOW)API: http://localhost:8080$(NC)"
	@echo "$(YELLOW)PostgreSQL: localhost:5432$(NC)"

## docker-up-dev: Start development environment with docker-compose
docker-dev:
	@echo "$(CYAN)Starting development environment with docker-compose...$(NC)"
	@docker-compose -f docker-compose.dev.yml up -d
	@echo "$(GREEN)✓ Development services started$(NC)"
	@echo "$(YELLOW)API: http://localhost:8080 (with hot reload)$(NC)"
	@echo "$(YELLOW)PostgreSQL: localhost:5432$(NC)"

## docker-down: Stop services with docker-compose
docker-down:
	@echo "$(CYAN)Stopping services with docker-compose...$(NC)"
	@docker-compose down
	@echo "$(GREEN)✓ Services stopped$(NC)"

## docker-down-dev: Stop development services
docker-down-dev:
	@echo "$(CYAN)Stopping development services...$(NC)"
	@docker-compose -f docker-compose.dev.yml down
	@echo "$(GREEN)✓ Development services stopped$(NC)"

## docker-logs: View docker-compose logs
docker-logs:
	@echo "$(CYAN)Viewing docker-compose logs...$(NC)"
	@docker-compose logs -f

## docker-logs-dev: View development docker-compose logs
docker-logs-dev:
	@echo "$(CYAN)Viewing development docker-compose logs...$(NC)"
	@docker-compose -f docker-compose.dev.yml logs -f

## docker-restart: Restart services
docker-restart: docker-down docker-up
	@echo "$(GREEN)✓ Services restarted$(NC)"

## docker-clean: Remove containers, volumes, and networks
docker-clean:
	@echo "$(CYAN)Cleaning Docker resources...$(NC)"
	@docker-compose down -v --remove-orphans || true
	@docker-compose -f docker-compose.dev.yml down -v --remove-orphans || true
	@echo "$(GREEN)✓ Docker resources cleaned$(NC)"

## docker-ps: Show running containers
docker-ps:
	@echo "$(CYAN)Running containers:$(NC)"
	@docker-compose ps

## setup: Initial setup (install tools, deps, check env)
setup: install-tools deps check-env
	@echo "$(GREEN)✓ Setup complete!$(NC)"
	@echo "$(CYAN)You can now run:$(NC)"
	@echo "$(CYAN)  - 'make air' to start the server locally$(NC)"
	@echo "$(CYAN)  - 'make docker-dev' to start with Docker Compose$(NC)"

## docker-setup: Setup with Docker (start services)
docker-setup: docker-up
	@echo "$(GREEN)✓ Docker setup complete!$(NC)"
	@echo "$(CYAN)Services are running. API available at http://localhost:8080$(NC)"

## dev: Development mode (run with air after checking env)
dev: check-env air

## prod: Production mode (build and run)
prod: check-env build run-prod

