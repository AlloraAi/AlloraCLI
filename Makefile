# AlloraCLI Makefile

# Build variables
BINARY_NAME=allora
BINARY_PATH=bin/$(BINARY_NAME)
MAIN_PATH=./cmd/allora

# Version variables
VERSION?=dev
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

# Go variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

# Default target
.PHONY: all
all: clean deps lint test build

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all          - Clean, download deps, lint, test, and build"
	@echo "  build        - Build the application"
	@echo "  build-all    - Build for all platforms"
	@echo "  test         - Run tests"
	@echo "  test-cover   - Run tests with coverage"
	@echo "  lint         - Run linters"
	@echo "  clean        - Clean build artifacts"
	@echo "  deps         - Download dependencies"
	@echo "  install      - Install to GOPATH/bin"
	@echo "  release      - Create release builds"

# Download dependencies
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Build the application
.PHONY: build
build:
	@echo "Building AlloraCLI..."
	@mkdir -p bin
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "Build completed: $(BINARY_PATH)"

# Build for Windows
.PHONY: build-windows
build-windows:
	@echo "Building AlloraCLI for Windows..."
	@mkdir -p bin
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_PATH).exe $(MAIN_PATH)
	@echo "Build completed: $(BINARY_PATH).exe"

# Build for Linux
.PHONY: build-linux
build-linux:
	@echo "Building AlloraCLI for Linux..."
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_PATH)-linux $(MAIN_PATH)
	@echo "Build completed: $(BINARY_PATH)-linux"

# Build for macOS
.PHONY: build-macos
build-macos:
	@echo "Building AlloraCLI for macOS..."
	@mkdir -p bin
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_PATH)-macos $(MAIN_PATH)
	@echo "Build completed: $(BINARY_PATH)-macos"

# Build for all platforms
.PHONY: build-all
build-all: build-windows build-linux build-macos

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@echo "Clean completed"

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Lint code
.PHONY: lint
lint:
	@echo "Linting code..."
	golangci-lint run ./...

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Tidy dependencies
.PHONY: tidy
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

# Run the application
.PHONY: run
run: build
	@echo "Running AlloraCLI..."
	./$(BINARY_PATH) --help

# Install the application
.PHONY: install
install: build
	@echo "Installing AlloraCLI..."
	@mkdir -p $(HOME)/.local/bin
	@cp $(BINARY_PATH) $(HOME)/.local/bin/$(BINARY_NAME)
	@echo "AlloraCLI installed to $(HOME)/.local/bin/$(BINARY_NAME)"

# Development build (faster, no optimizations)
.PHONY: dev
dev:
	@echo "Building development version..."
	@mkdir -p bin
	go build -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "Development build completed: $(BINARY_PATH)"

# Release build (optimized)
.PHONY: release
release:
	@echo "Building release version..."
	@mkdir -p bin
	go build -a -installsuffix cgo $(LDFLAGS) -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "Release build completed: $(BINARY_PATH)"

# Docker build
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t alloracli:$(VERSION) .

# Docker run
.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	docker run --rm -it alloracli:$(VERSION)

# Generate documentation
.PHONY: docs
docs:
	@echo "Generating documentation..."
	go run docs/generate.go

# Check for security vulnerabilities
.PHONY: security
security:
	@echo "Checking for security vulnerabilities..."
	govulncheck ./...

# Benchmark
.PHONY: bench
bench:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...
