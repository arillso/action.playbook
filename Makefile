.DEFAULT_GOAL := help

.PHONY: help lint lint-go lint-yaml format test build build-docker action-test clean

## Linting
lint: lint-go lint-yaml ## Run all linters

lint-go: ## Run Go linter (golangci-lint)
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found, skipping"; \
	fi

lint-yaml: ## Run YAML linter
	@if command -v yamllint >/dev/null 2>&1; then \
		if [ -f .yamllint.yml ]; then yamllint -c .yamllint.yml .; else yamllint .; fi; \
	else \
		echo "yamllint not found, skipping"; \
	fi

## Formatting
format: ## Format Go code
	gofmt -s -w .
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	else \
		echo "goimports not found, skipping"; \
	fi

## Testing
test: ## Run Go tests
	go test -v ./...

action-test: build-docker ## Run action tests with test playbook
	@docker run --rm \
		-v $(shell pwd):/github/workspace \
		-w /github/workspace \
		action-playbook:local \
		--playbook tests/basic_playbook.yml \
		--inventory tests/hosts.yml \
		--galaxy-requirements tests/requirements.yml

## Building
build: ## Build Go binary
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main .

build-docker: ## Build Docker image
	docker build -t action-playbook:local .

## Cleanup
clean: ## Remove build artifacts
	go clean
	rm -rf build/ dist/ megalinter-reports/ main

## Help
help: ## Show available targets
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?##' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'
