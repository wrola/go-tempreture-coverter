.PHONY: help build test lint lint-fix fmt clean

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the temperature converter binary
	go build -o temp-converter .

test: ## Run all tests
	go test ./...

bench: ## Run benchmarks
	go test -bench=. ./converter

lint: ## Run golangci-lint
	golangci-lint run

lint-fix: ## Run golangci-lint with auto-fix
	golangci-lint run --fix

fmt: ## Format code with gofmt
	gofmt -s -w .

clean: ## Clean build artifacts
	rm -f temp-converter

all: clean fmt lint test build ## Run all checks and build
