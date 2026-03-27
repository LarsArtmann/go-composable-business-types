# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [Unreleased]

### Added

- **Policy Documentation**: Comprehensive library policies ([POLICY.md](./POLICY.md))
  - Versioning policy (Semantic Versioning)
  - Breaking changes policy
  - Deprecation policy with migration guidelines
  - API stability guarantees
  - Security vulnerability reporting process
  - Contribution guidelines with conventional commits
  - Release process and quality gates
  - Dependency management policy
- **Support Documentation**: User support guide ([SUPPORT.md](./SUPPORT.md))
- **Security Documentation**: Security policy ([SECURITY.md](./SECURITY.md))
- **GitHub Templates**: Issue templates for bugs, features, and breaking changes
- **Pull Request Template**: Structured PR template with quality gates
- **Implementation Checklist**: Policy implementation tracking ([docs/planning/POLICY_IMPLEMENTATION_CHECKLIST.md](./docs/planning/POLICY_IMPLEMENTATION_CHECKLIST.md))
- **Dependabot Configuration**: Automated dependency updates (`.github/dependabot.yml`)
- **CI Security Scanning**: Added govulncheck to GitHub Actions workflow

### Changed

- **CI Workflow**: Updated Go versions to 1.26, 1.27, 1.28; added security job
- **golangci.yml**: Optimized for memory efficiency with reduced linter set

### Deprecated

### Removed

### Fixed

### Security

## [0.1.0] - 2026-01-01

### Added

- Initial release
