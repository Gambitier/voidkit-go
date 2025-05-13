.PHONY: help dev up down build test tests clean proto logs prepare

# Default target
help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## ' Makefile | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

dev: ## Start services and run the Go server locally (depends on `up` command)
serve: ## Run the Go server locally
up: ## Start all Docker services
down: ## Stop all Docker services
build: ## Build the Go binary
test: ## Run particular Go test
tests: ## Run all Go tests
clean: ## Remove build artifacts and stop services
proto: ## Generate protobuf and gRPC code
logs: ## Tail logs from all services
prepare: ## Create necessary directories for volumes

# Development commands
dev: up serve

serve:
	@echo "Starting server..."
	@go run cmd/server/main.go --config="default.yaml" --env=development 

up:
	@echo "Starting services..."
	@docker-compose --env-file docker.env -f docker-compose.dev.yml up -d

down:
	@echo "Stopping services..."
	@docker-compose --env-file docker.env -f docker-compose.dev.yml down

build:
	@echo "Building service..."
	@go mod tidy
	@go build -o bin/voidkitgo cmd/server/main.go

test:
	@echo "Running tests with debug output..."
	@TEST_ENV_VARIABLE=true \
	go test -v -count=1 ./tests/... -run $(TEST)

tests:
	@echo "Running tests..."
	@go test -v -count=1 ./... | grep -v "no test files"

clean:
	@echo "Cleaning up..."
	@rm -rf bin
	@make down

proto:
	@echo "Generating protobuf code..."
	@protoc -I. \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/**/*.proto

logs:
	@echo "Showing logs..."
	@docker-compose --env-file docker.env -f docker-compose.dev.yml logs -f
