# Go Get Images - Makefile
# Cross-platform build and test automation

# Variables
BINARY_NAME=go-get-imgs
MAIN_PATH=cmd/go-get-imgs
BUILD_DIR=build
DIST_DIR=dist
VERSION?=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD)

# Go build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

# Platforms to build for
PLATFORMS=darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64 windows/arm64

# Default target
.PHONY: all
all: clean build

# Build for current platform
.PHONY: build
build:
	@echo "Building for current platform..."
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./$(MAIN_PATH)

# Build for all platforms
.PHONY: build-all
build-all: clean
	@echo "Building for all platforms..."
	@mkdir -p $(BUILD_DIR)
	@for platform in $(PLATFORMS); do \
		IFS='/' read -r GOOS GOARCH <<< "$$platform"; \
		BINARY_NAME_FULL="$(BINARY_NAME)-$$GOOS-$$GOARCH"; \
		if [ "$$GOOS" = "windows" ]; then \
			BINARY_NAME_FULL="$$BINARY_NAME_FULL.exe"; \
		fi; \
		echo "Building for $$GOOS/$$GOARCH..."; \
		GOOS=$$GOOS GOARCH=$$GOARCH go build $(LDFLAGS) -o $(BUILD_DIR)/$$BINARY_NAME_FULL ./$(MAIN_PATH); \
	done

# Build for specific platforms
.PHONY: build-darwin build-linux build-windows
build-darwin:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./$(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./$(MAIN_PATH)

build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./$(MAIN_PATH)
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./$(MAIN_PATH)

build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./$(MAIN_PATH)
	GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-arm64.exe ./$(MAIN_PATH)

# Test targets
.PHONY: test test-coverage test-race
test:
	@echo "Running tests..."
	go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o docs/coverage.html
	@echo "Coverage report generated at docs/coverage.html"

test-race:
	@echo "Running tests with race detection..."
	go test -race -v ./...

# Run specific test categories
.PHONY: test-unit test-integration test-cross-platform
test-unit:
	@echo "Running unit tests..."
	go test -v -run TestUnit ./...

test-integration:
	@echo "Running integration tests..."
	go test -v -run TestIntegration ./...

test-cross-platform:
	@echo "Running cross-platform tests..."
	go test -v -run TestCrossPlatform ./...

# Benchmarks
.PHONY: bench
bench:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -rf $(DIST_DIR)
	rm -f coverage.out

# Create distribution packages
.PHONY: dist
dist: build-all
	@echo "Creating distribution packages..."
	@mkdir -p $(DIST_DIR)
	@cd $(BUILD_DIR) && for file in *; do \
		tar -czf ../$(DIST_DIR)/$$file.tar.gz $$file; \
	done

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
.PHONY: lint
lint:
	@echo "Linting code..."
	golangci-lint run

# Run the application
.PHONY: run
run: build
	@echo "Running application..."
	./$(BUILD_DIR)/$(BINARY_NAME) examples/sample.csv 3

# Docker targets
.PHONY: docker-build docker-test docker-run
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME) docker/

docker-test:
	@echo "Running tests in Docker..."
	docker-compose -f docker/docker-compose.yml up tester

docker-run:
	@echo "Running application in Docker..."
	docker-compose -f docker/docker-compose.yml up runtime

# Development helpers
.PHONY: dev-setup
dev-setup: deps
	@echo "Setting up development environment..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build          - Build for current platform"
	@echo "  build-all      - Build for all platforms"
	@echo "  build-darwin   - Build for macOS"
	@echo "  build-linux    - Build for Linux"
	@echo "  build-windows  - Build for Windows"
	@echo "  test           - Run all tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  test-race      - Run tests with race detection"
	@echo "  bench          - Run benchmarks"
	@echo "  clean          - Clean build artifacts"
	@echo "  dist           - Create distribution packages"
	@echo "  deps           - Install dependencies"
	@echo "  fmt            - Format code"
	@echo "  lint           - Lint code"
	@echo "  run            - Run the application"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-test    - Run tests in Docker"
	@echo "  docker-run     - Run application in Docker"
	@echo "  dev-setup      - Setup development environment"
	@echo "  help           - Show this help" 