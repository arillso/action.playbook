
# Set PROJECT_DIR to the CI-provided project directory if available; otherwise fallback to the current directory.
ifndef CI_PROJECT_DIR
	ifndef GITHUB_WORKSPACE
		PROJECT_DIR := $(shell pwd)
	else
		PROJECT_DIR := $(GITHUB_WORKSPACE)
	endif
else
	PROJECT_DIR := $(CI_PROJECT_DIR)
endif

# go-build-firewallguard:
# Build firewallguard for Linux AMD64.
go-build-ansible:
	GO111MODULE=on GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" main.go

# lint-wallpaper:
# Run GolangCI-Lint on the wallpaper project.
lint-ansible:
	@docker run --rm -v $(shell pwd):/workspace -w /workspace golangci/golangci-lint bash -c "go fmt && go mod tidy && GOARCH=amd64 golangci-lint run --fix"

# run-megalinter:
# Run Megalinter locally to check code quality across multiple languages.
run-megalinter:
	@docker run --rm --name megalint -v $(shell pwd):/tmp/lint busybox rm -rf /tmp/lint/megalinter-reports /tmp/lint/packages/firewallguard/assets/static/js/cdn.min.js /tmp/lint/assets/abuild/6696f7cf.rsa
	@docker run --rm --name megalint -v $(shell pwd):/tmp/lint -e MARKDOWN_SUMMARY_REPORTER=true oxsecurity/megalinter:v8.4.2


format-code: ## Format code files using Prettier via Docker.
	@docker run --rm --name prettier -v $(PROJECT_DIR):$(PROJECT_DIR) -w /$(PROJECT_DIR) node:alpine npx prettier . --write

format-all: format-code ## Run both format-code and format-eclint.
	@echo "Formatting completed."

help: ## Show an overview of available targets.
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?##' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

action-build: ## Build the Ansible Docker image.
	@docker build \
		-t action:latest \
		-f Dockerfile .
