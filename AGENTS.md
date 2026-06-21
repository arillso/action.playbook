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

## Base image coupling (arillso/ansible)

The production stage builds `FROM arillso/ansible:<ansible-core-version>@sha256:<digest>`
(produced by the `arillso/docker.ansible` repo).

- **Tag = the ansible-core version** the image ships (e.g. `2.21.1`), not a
  docker.ansible-specific semver. docker.ansible publishes image tags `latest`,
  `<ansible-core-version>` and `<version>-<sha>`; its git releases are
  date-tagged (rolling release), but the image version tag tracks ansible-core.
- **Pin tag + digest.** The tag documents which Ansible the action runs; the
  `@sha256` digest makes the build reproducible.
- **Renovate keeps it current.** The `dockerfile` manager (via the
  `renovate-actions` preset, group `docker-images`) bumps both the tag and the
  digest when a newer `arillso/ansible` version ships — no manual cadence to
  track. Verified with `renovate --dry-run`: it offers the tag+digest update.
- **On an ansible-core major/minor bump,** review the action's behaviour
  against the new Ansible before merging the Renovate PR (flag/output changes).

## Do Not

- Add unnecessary dependencies
- Break backwards compatibility of action inputs
- Commit secrets or private keys
