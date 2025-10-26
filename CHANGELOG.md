# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.3.0] - 2025-10-26

### Added

- Comprehensive SSH key format tests for CI/CD validation
  - Test Unix LF format keys
  - Test Windows CRLF format keys with auto-normalization
  - Test keys missing trailing newlines with auto-fix
  - Test invalid key format with proper validation
- Missing `.github/scripts/utils/generate-summary.sh` utility script for workflow reports

### Changed

- **Major:** Upgraded `go.ansible` dependency from v1.0.0 to v1.2.0
  - Automatic CRLF → LF conversion for SSH private keys
  - Automatic trailing newline addition per RFC 7468
  - SSH key format validation (RSA, OpenSSH, EC, DSA)
- Streamlined SSH Authentication documentation (70% reduction: 171 → 53 lines)
- Disabled MegaLinter auto-push to prevent CI failures
  - Set `APPLY_FIXES_MODE` to "none" (report only)
  - Set `APPLY_FIXES_EVENT` to "none"

### Fixed

- YAML indentation in test playbooks (4 spaces → 2 spaces, Ansible convention)
- Markdown line length violations (max 175 characters)
- Shellcheck SC2129: Use grouped redirects for better performance
- Shellcheck SC2086: Quote `$GITHUB_ENV` to prevent globbing
- KICS security scanner false positives with suppression comments
- Replaced `unix2dos` with `sed` for CRLF conversion (compatibility with GitHub Runners)
- SSH key handling with bastion hosts and ProxyCommand configurations

## [0.2.2] - 2025-06-21

### Fixed

- Uses explicit `docker://` prefix for container image reference

## [0.2.1] - 2025-06-21

### Changed

- Updates container publish workflow for improved reliability

## [0.2.0] - 2025-06-21

### Added

- Enhanced action testing with comprehensive CI workflows
- Improved workflow organization and test coverage

### Changed

- Updated GitHub Actions dependencies
  - `actions/attest-build-provenance` to v2
  - `aquasecurity/trivy-action` to v0.31.0
- Updated Docker base images
  - `arillso/ansible:2.18.6` with pinned digests

### Fixed

- Tag matching regex for version detection
- Upgraded `github.com/urfave/cli/v3` to v3.3.8

## [0.1.1] - 2025-06-15

### Added

- Comprehensive CI/CD workflows for testing and security
- Enhanced security scanning and code quality checks
- Renovate configuration for automated dependency updates

### Changed

- Migrated from `urfave/cli/v2` to `urfave/cli/v3` (v3.3.3)
- Updated to Go 1.24.2
- Updated shadow package version
- Improved Renovate configuration
- Pinned GitHub Actions to commit SHAs for security

### Fixed

- Various dependency updates and security patches
- Docker image digest updates

## [0.1.0] - 2025-04-23

### Added

- Enhanced CI/CD and code quality workflows
- New Galaxy options for Ansible configuration:
  - `galaxy_api_key`
  - `galaxy_api_server_url`
  - `galaxy_collections_path`
  - `galaxy_disable_gpg_verify`
  - `galaxy_force_with_deps`
  - `galaxy_ignore_certs`
  - `galaxy_ignore_signature_status_codes`
  - `galaxy_keyring`
  - `galaxy_offline`
  - `galaxy_pre`
  - `galaxy_required_valid_signature_count`
  - `galaxy_signature`
  - `galaxy_timeout`
  - `galaxy_upgrade`
  - `galaxy_no_deps`

### Changed

- Enhanced codebase with linting and formatting improvements
- Major refactoring of CI/CD workflows

## [0.0.9] - 2023-11-11

### Changed

- **Dependabot Configuration:** Improved formatting in `.github/dependabot.yml` for enhanced readability and consistency.
- **Changelog Formatting:** Updated the changelog format in `CHANGELOG.md` for better clarity and adherence to the human-readable standard.
- **Dependency Updates:** Modified versions in `go.mod` and `go.sum` to incorporate newer dependency versions.
- **`main.go` Adjustments:** Extensive revision of flag descriptions and comments to improve user-friendliness and comprehensibility.

## [0.0.8] - 2023-11-09

### Changed

- Add 'GalaxyForce' boolean variable.

## [0.0.7] - 2023-11-02

### Changed

- Bump arillso/ansible from 2.15.4 to 2.15.5.

## [0.0.6] - 2023-10-30

### Changed

- Bump go from 1.14 to 1.18.
- Bump arillso/ansible from 2.12.4 to 2.15.4.
- Bump github.com/arillso/go.ansible from 0.0.1 to 0.0.2.
- Bump github.com/joho/godotenv from 1.4.0 to 1.5.1.
- Bump github.com/urfave/cli/v2 from 2.11.1 to 2.25.7.

## [0.0.5] - 2022-07-28

### Changed

- Bump github.com/urfave/cli/v2 from 2.3.0 to 2.11.1.

## [0.0.4] - 2021-11-28

### Changed

- Bump arillso/ansible from 2.10.3 to 2.12.0.

## [0.0.3] - 2020-08-01

### Changed

- Bump arillso/ansible from 2.9.7 to 2.9.9.

## [0.0.2] - 2020-05-02

### Changed

- Bump arillso/ansible from 2.9.6 to 2.9.7.

## [0.0.1] - 2020-04-06

### Added

- Initial Commit.
