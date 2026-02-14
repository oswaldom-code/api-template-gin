# API Template Gin — Makefile
# Usage: make <target>
# Run 'make' or 'make help' to list available targets.

APP_NAME   := api-template-app
BINARY     := bin/api
MAIN       := main.go

.DEFAULT_GOAL := help

.PHONY: run cli-test build test test-cover vet lint fmt tidy \
        docker-up docker-down docker-build swagger clean help

run: ## Start the HTTP server (development)
	go run $(MAIN) server

cli-test: ## Test DB connection via CLI adapter
	go run $(MAIN) cli -f test

build: ## Build production binary (CGO disabled)
	CGO_ENABLED=0 go build -o $(BINARY) $(MAIN)

test: ## Run all tests with verbose output
	go test ./... -v

test-cover: ## Run tests with coverage report
	go test ./... -cover -coverprofile=coverage.out
	@echo "Coverage report: coverage.out"
	@echo "Run 'go tool cover -html=coverage.out' to view in browser"

vet: ## Run static analysis (go vet)
	go vet ./...

lint: vet ## Run vet + check formatting
	@test -z "$$(gofmt -l .)" || { echo "Files need formatting:"; gofmt -l .; exit 1; }

fmt: ## Format all Go source files
	gofmt -w .

tidy: ## Tidy module dependencies
	go mod tidy

docker-up: ## Start app + PostgreSQL with docker-compose
	docker-compose up --build

docker-down: ## Stop docker-compose services
	docker-compose down

docker-build: ## Build standalone Docker image
	docker build -t $(APP_NAME) .

swagger: ## Generate API docs from OpenAPI spec
	./scripts/swagger.sh

clean: ## Remove build artifacts and coverage reports
	rm -rf bin/ coverage.out

help: ## Show this help
	@printf "\nUsage: make \033[36m<target>\033[0m\n\n"
	@printf "Targets:\n"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
