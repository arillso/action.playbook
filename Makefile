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

go-build-ansible: ## Build the Ansible application for Linux AMD64.
	GO111MODULE=on GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" main.go

lint-ansible: ## Run GolangCI-Lint to format code, tidy modules, and automatically fix lint issues for the Ansible project.
	@docker run --rm -v $(shell pwd):/workspace -w /workspace golangci/golangci-lint bash -c "go fmt && go mod tidy && GOARCH=amd64 golangci-lint run --fix"

run-megalinter: ## Execute Megalinter to perform code quality checks across multiple programming languages.
	@docker run --rm --name megalint -v $(shell pwd):/tmp/lint busybox rm -rf /tmp/lint/megalinter-reports /tmp/lint/packages/firewallguard/assets/static/js/cdn.min.js /tmp/lint/assets/abuild/6696f7cf.rsa
	@docker run --rm --name megalint -v $(shell pwd):/tmp/lint -e MARKDOWN_SUMMARY_REPORTER=true oxsecurity/megalinter:v8.4.2

format-code: ## Format all code files in the project using Prettier via Docker.
	@docker run --rm --name prettier -v $(PROJECT_DIR):$(PROJECT_DIR) -w /$(PROJECT_DIR) node:alpine npx prettier . --write

format-all: format-code ## Execute all available code formatting tasks.
	@echo "Formatting completed."

action-build: ## Build the Docker image for the Ansible action using the specified Dockerfile.
	@docker build \
		-t action:latest \
		-f Dockerfile .

tests: action-build ## Run action tests with the built Docker image using the provided Ansible playbook, inventory, and Galaxy file.
	@docker run --rm \
		-v "$(PROJECT_DIR):/github/workspace" \
		-w "/github/workspace" \
		-e ANSIBLE_PLAYBOOK=tests/basic_playbook.yml -e ANSIBLE_INVENTORY=tests/hosts.yml -e ANSIBLE_GALAXY_FILE=tests/requirements.yml \
		action:latest

help: ## Display a list of all available make targets along with their descriptions.
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?##' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'
