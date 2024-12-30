# Project variables
APP_NAME := goshort
BIN_DIR := bin
BUILD_DIR := $(BIN_DIR)/$(APP_NAME)
PKG_LIST := $(shell go list ./... | grep -v /tests)

# Commands
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOLINT := golangci-lint
GOFMT := gofmt
GOMODTIDY := $(GOCMD) mod tidy
DOCKER_COMPOSE := docker compose

# Build variables
BUILD_FLAGS := -ldflags "-s -w"

# Targets
.PHONY: all build test lint fmt run clean docker-up docker-down docker-build help

all: build ## Build the application

build: clean ## Build the application binaries
	@echo "Building the $(APP_NAME)..."
	$(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR) ./cmd/server
	@echo "Build completed."

test: ## Run unit tests
	@echo "Running tests..."
	$(GOTEST) -v -race ./...
	@echo "Tests completed."

lint: ## Run lint checks
	@echo "Running lint checks..."
	$(GOLINT) run ./...
	@echo "Linting completed."

fmt: ## Format the code
	@echo "Formatting code..."
	$(GOFMT) -s -w $(PKG_LIST)
	@echo "Code formatted."

run: build ## Build and run the application
	@echo "Running the application..."
	$(BUILD_DIR)

clean: ## Clean up build artifacts
	@echo "Cleaning up..."
	rm -rf $(BIN_DIR)
	@echo "Cleanup completed."

docker-up: ## Start Docker containers (app and database)
	@echo "Starting Docker containers..."
	$(DOCKER_COMPOSE) up -d
	@echo "Docker containers started."

docker-down: ## Stop Docker containers (app and database)
	@echo "Stopping Docker containers..."
	$(DOCKER_COMPOSE) down
	@echo "Docker containers stopped."

docker-build: ## Build Docker images for the application
	@echo "Building Docker images..."
	$(DOCKER_COMPOSE) build
	@echo "Docker images built."

help: ## Display this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z0-9_-]+:.*?##' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "} {printf "  %-15s %s\n", $$1, $$2}'
