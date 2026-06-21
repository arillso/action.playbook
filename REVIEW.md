# Code Review Guidelines

## Scope

In scope:

- Go source changes (`main.go`, `*.go`)
- `action.yml` input/output contract changes
- `Dockerfile` changes (build, base image, multi-arch)
- CI/CD workflow changes
- Renovate configuration updates

Out of scope:

- Renovate dependency-only PRs (patch/minor with automerge enabled)
- Generated changelog entries from release automation

## Required checks

- No secrets committed — no credentials, tokens, or keys in source or workflows
- `gofmt` clean and `golangci-lint` passes
- `go test ./...` passes (with `-race`)
- Action integration tests pass (the `test-*` jobs in pull-request.yml)
- Security scans pass (gitleaks, trivy, CodeQL)
- Backwards-compatible `action.yml` inputs — no breaking changes to existing input names without a major bump
- Numeric inputs are bounds-validated; secrets (SSH keys, passphrases) never interpolated into shell

## Severity levels

| Level        | Meaning                                             | Merge impact       |
| ------------ | --------------------------------------------------- | ------------------ |
| Bug          | Incorrect behavior or broken contract               | Blocks merge       |
| Nit          | Minor issue — suboptimal but not incorrect          | Non-blocking       |
| Pre-existing | Issue present before this PR; flagged for awareness | No action required |

## Skip

- Renovate PRs with `automerge: true` (patch/minor) after CI passes
- Documentation-only changes with no functional impact
