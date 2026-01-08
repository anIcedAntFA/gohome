# Output binary file name
BINARY_NAME=gohome
# Main file path
MAIN_PATH=./cmd/gohome/main.go

# Version information
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags to inject version info
LDFLAGS=-ldflags "-X github.com/anIcedAntFA/gohome/internal/version.Version=$(VERSION) \
                   -X github.com/anIcedAntFA/gohome/internal/version.Commit=$(COMMIT) \
                   -X github.com/anIcedAntFA/gohome/internal/version.Date=$(BUILD_DATE)"

# System variables
GO=go
GOTEST=$(GO) test
GOBUILD=$(GO) build
GOMOD=$(GO) mod
GOLINT=golangci-lint

.PHONY: all build run test clean lint tidy help

# Default target when running 'make'
default: help

## build: Build binary for current operating system
build:
	@echo "  >  Building binary..."
	$(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME) $(MAIN_PATH)
	@echo "  >  Build successful! Binary is at bin/$(BINARY_NAME)"

## run: Run application directly (dev mode)
run:
	$(GO) run $(MAIN_PATH)

## install: Install tool to $GOPATH/bin (to run from anywhere)
install:
	@echo "  >  Installing..."
	$(GOBUILD) $(LDFLAGS) -o $(shell go env GOPATH)/bin/$(BINARY_NAME) $(MAIN_PATH)
	@echo "  >  Installed successfully! You can run '$(BINARY_NAME)' now."

## clean: Remove old build artifacts
clean:
	@echo "  >  Cleaning build cache..."
	$(GO) clean
	rm -rf bin/
	@echo "  >  Cleaned."

## test: Run unit tests
test:
	@echo "  >  Running tests..."
	$(GOTEST) -v ./...

## lint: Run golangci-lint to check code quality
lint:
	@echo "  >  Linting code..."
	$(GOLINT) run

## tidy: Clean up go.mod and go.sum
tidy:
	$(GOMOD) tidy

## help: Display this help message
help:
	@echo "Choose a command run in $(BINARY_NAME):"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'