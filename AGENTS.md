# action.playbook

GitHub Action for running Ansible Playbooks with advanced configuration options.

## Context

This is a Docker-based GitHub Action written in Go that executes Ansible playbooks. The action runs inside a container based on `arillso/ansible` and provides extensive configuration options for Galaxy, SSH, Vault, and privilege escalation.

## Structure

```text
.
├── main.go              # Go entrypoint - CLI argument parsing
├── action.yml           # GitHub Action definition
├── Dockerfile           # Container build (based on arillso/ansible)
├── go.mod / go.sum      # Go dependencies
├── tests/               # Test playbooks and inventory files
└── .github/workflows/   # CI/CD pipelines
```

## Conventions

- **Language**: Go 1.25+
- **Formatting**: gofmt (built-in)
- **Linting**: golangci-lint with staticcheck, govet, gosimple
- **Dependencies**: Minimal - only go.ansible, godotenv, urfave/cli
- **Commits**: Conventional Commits format

## Key Patterns

- SSH keys are automatically normalized (CRLF → LF, trailing newlines)
- All Ansible options are exposed as action inputs
- Docker image is published to ghcr.io/arillso/action.playbook

## Do Not

- Add unnecessary dependencies
- Break backwards compatibility of action inputs
- Commit secrets or private keys
