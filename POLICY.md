# Library Policy

> **go-composable-business-types** — Governance, stability, and operational policies

**Version:** 1.0.0  
**Effective Date:** 2026-03-27  
**Last Updated:** 2026-03-27

---

## Table of Contents

1. [Versioning Policy](#1-versioning-policy)
2. [Breaking Changes Policy](#2-breaking-changes-policy)
3. [Deprecation Policy](#3-deprecation-policy)
4. [API Stability Guarantees](#4-api-stability-guarantees)
5. [Security Policy](#5-security-policy)
6. [Contribution Guidelines](#6-contribution-guidelines)
7. [Release Process](#7-release-process)
8. [Dependency Management](#8-dependency-management)
9. [Quality Gates](#9-quality-gates)
10. [Support Policy](#10-support-policy)
11. [License & Legal](#11-license--legal)

---

## 1. Versioning Policy

### 1.1 Semantic Versioning

This library follows [Semantic Versioning 2.0.0](https://semver.org/) (SemVer):

| Version   | Format              | Meaning                                           |
| --------- | ------------------- | ------------------------------------------------- |
| **Major** | `v1.0.0` → `v2.0.0` | Breaking API changes requiring code modifications |
| **Minor** | `v1.0.0` → `v1.1.0` | New backward-compatible functionality added       |
| **Patch** | `v1.0.0` → `v1.0.1` | Backward-compatible bug fixes                     |

### 1.2 Pre-Release Versions

| Tag              | Meaning           | Stability                                         |
| ---------------- | ----------------- | ------------------------------------------------- |
| `v0.x.x`         | Development phase | No stability guarantees; breaking changes allowed |
| `v1.0.0-alpha.x` | Early testing     | Unstable; subject to change                       |
| `v1.0.0-beta.x`  | Feature-complete  | Likely stable; seeking feedback                   |
| `v1.0.0-rc.x`    | Release candidate | Stable unless critical issues found               |

### 1.3 Module Path

- **v0.x.x and v1.x.x**: `github.com/larsartmann/go-composable-business-types`
- **v2.x.x+**: `github.com/larsartmann/go-composable-business-types/v2`

```go
// v1 import path
go get github.com/larsartmann/go-composable-business-types

// v2 import path (when released)
go get github.com/larsartmann/go-composable-business-types/v2
```

### 1.4 Go Version Support

| Library Version | Minimum Go | Supported Go Versions |
| --------------- | ---------- | --------------------- |
| Current (v0.x)  | 1.26       | 1.26+                 |
| v1.0.0+         | 1.26       | 1.26, 1.27+           |

**Policy**: Support the current and previous major Go releases, following [Go's release policy](https://go.dev/doc/devel/release#policy).

---

## 2. Breaking Changes Policy

### 2.1 What Constitutes a Breaking Change

Breaking changes **require a major version increment**:

| Change Type                         | Breaking? | Example                                |
| ----------------------------------- | --------- | -------------------------------------- |
| Removing exported function/type     | **Yes**   | Deleting `NewID()`                     |
| Renaming exported function/type     | **Yes**   | `ActorChain` → `ActorList`             |
| Changing function signature         | **Yes**   | Adding required parameter              |
| Adding method to exported interface | **Yes**   | Breaks existing implementations        |
| Changing behavior contract          | **Yes**   | `IsZero()` returning different results |
| Removing exported constant          | **Yes**   | Removing `TriggerManual`               |
| Changing JSON serialization format  | **Yes**   | `ID` serializing as object vs string   |

### 2.2 What Is NOT a Breaking Change

These changes are **backward-compatible** and require only minor/patch:

| Change Type                         | Breaking? | Example                                  |
| ----------------------------------- | --------- | ---------------------------------------- |
| Adding new exported function/type   | No        | Adding `NewActorChainFromSlice()`        |
| Adding new method to concrete type  | No        | Adding `Filter()` to `ActorChain`        |
| Adding new optional field to struct | No        | New field in `DataPoint[T]` (unexported) |
| Adding new enum value               | No        | Adding `TriggerBatch`                    |
| Adding new package                  | No        | Creating `validate/` package             |
| Performance improvements            | No        | Optimizing `Compare()` implementation    |
| Bug fixes                           | No        | Fixing `BoundedString` validation        |

### 2.3 Breaking Change Process

1. **RFC Required**: Open GitHub issue with `breaking-change` label
2. **Community Input**: 30-day comment period minimum
3. **Migration Guide**: Document upgrade path
4. **Deprecation Period**: Maintain old API for at least 2 minor versions
5. **Major Version**: Release as v2.x.x with module path change

### 2.4 Exceptions

The following may bypass major version increment with explicit justification:

- **Security vulnerabilities**: Immediate fix with clear documentation
- **Specification errors**: When documented behavior was incorrect
- **Unspecified behavior**: Changes to undocumented behavior

---

## 3. Deprecation Policy

### 3.1 Deprecation Markers

Use standard Go deprecation comments:

```go
// Deprecated: Use [NewActorChain] instead. This function will be removed in v2.0.0.
func OldActorChainConstructor() ActorChain[string] { ... }

// Deprecated: ActorKindLegacy is replaced by [ActorKindSystem].
// To migrate: replace all ActorKindLegacy with ActorKindSystem.
const ActorKindLegacy = ActorKindSystem
```

### 3.2 Deprecation Timeline

| Phase                | Duration         | Actions                                        |
| -------------------- | ---------------- | ---------------------------------------------- |
| **Announcement**     | Immediate        | Add deprecation comment; document in CHANGELOG |
| **Soft Deprecation** | 2 minor versions | Functionality maintained; warnings in tooling  |
| **Hard Deprecation** | Until next major | Functionality removed in next major version    |

### 3.3 Deprecation Documentation Requirements

Every deprecation must include:

1. **Replacement API**: Clear pointer to new approach
2. **Migration example**: Before/after code sample
3. **Timeline**: When functionality will be removed
4. **Rationale**: Why the change is being made

Example:

```go
// Deprecated: BoundedStringOf is replaced by [NewBoundedString] which provides
// clearer error messages and better validation.
//
// Migration:
//   // Before
//   fn := bounded.BoundedStringOf(1, 100)
//   s, _ := fn("value")
//
//   // After
//   s, err := bounded.NewBoundedString(1, 100, "value")
//
// Will be removed in v2.0.0.
func BoundedStringOf(min, max int) func(string) (BoundedString, error) { ... }
```

---

## 4. API Stability Guarantees

### 4.1 Stability Levels

| Level            | Package                                                                                                                                                | Stability Guarantee               |
| ---------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------ | --------------------------------- |
| **Stable**       | `actor/`, `datapoint/`, `temporal/`, `types/`, `enums/`, `money/`, `locale/`, `nanoid/`, `bounded/`, `importance/`, `tag/`, `validate/`, `pkg/errors/` | No breaking changes before v2.0.0 |
| **Experimental** | N/A                                                                                                                                                    | No current experimental packages  |

### 4.2 Compatibility Commitments

**Source Compatibility**: Code written against v1.x.x will compile against any v1.y.y where y ≥ x.

**Behavioral Compatibility**: Stable APIs maintain documented behavior across minor/patch versions.

**Serialization Compatibility**: JSON, SQL, and binary formats are stable within major versions.

### 4.3 Generic Type Stability

Generic types (`ActorEntry[T]`, `DataPoint[T]`, `ActorChain[T]`, `Reference[T]`, `Cause[T]`) follow these rules:

- Type parameters may be added (backward-compatible for inference)
- Type constraints may be relaxed (e.g., `comparable` to `any`)
- Type constraints may NOT be tightened

### 4.4 Interface Stability

Exported interfaces are **permanent commitments**:

| Interface          | Status | Notes                      |
| ------------------ | ------ | -------------------------- |
| `sql.Scanner`      | Stable | Standard library interface |
| `driver.Valuer`    | Stable | Standard library interface |
| `json.Marshaler`   | Stable | Standard library interface |
| `json.Unmarshaler` | Stable | Standard library interface |
| `fmt.Formatter`    | Stable | Standard library interface |
| `fmt.GoStringer`   | Stable | Standard library interface |

**Rule**: New methods are never added to exported interfaces within a major version.

---

## 5. Security Policy

### 5.1 Supported Versions

Security updates are provided for:

| Version        | Support Status                |
| -------------- | ----------------------------- |
| Latest minor   | Full security support         |
| Previous minor | Security support for 6 months |
| Older versions | No security support           |

### 5.2 Reporting Vulnerabilities

**Email**: security@lars.software

**Subject Line**: `[SECURITY] go-composable-business-types: <brief description>`

**Required Information**:

- Description of vulnerability
- Steps to reproduce
- Affected versions
- Proposed fix (if available)
- Contact for follow-up

### 5.3 Disclosure Timeline

| Phase                 | Timeline              | Actions                                |
| --------------------- | --------------------- | -------------------------------------- |
| **Acknowledgment**    | Within 48 hours       | Confirm receipt; assign CVE if needed  |
| **Assessment**        | Within 7 days         | Determine severity; develop fix        |
| **Fix Development**   | Within 30 days        | Create and test patch                  |
| **Pre-announcement**  | 7 days before release | Notify affected users with mitigations |
| **Public Disclosure** | With release          | Publish CVE, fix, and post-mortem      |

### 5.4 Severity Classification

| Severity     | Criteria                                    | Response Time          |
| ------------ | ------------------------------------------- | ---------------------- |
| **Critical** | Remote code execution; data corruption      | 48 hours               |
| **High**     | Authentication bypass; privilege escalation | 7 days                 |
| **Medium**   | Information disclosure; DoS                 | 30 days                |
| **Low**      | Minor security improvements                 | Next scheduled release |

### 5.5 Security Best Practices

This library implements:

- **FIPS 140-2**: Uses `sixafter/nanoid` for cryptographically secure IDs
- **No panic on invalid input**: All constructors return errors
- **No logging of sensitive data**: IDs and values are never logged internally
- **Bounds checking**: All numeric types validate ranges
- **SQL injection prevention**: Scanner/Valuer interfaces use parameterized queries

---

## 6. Contribution Guidelines

### 6.1 Getting Started

1. **Fork** the repository
2. **Clone** your fork: `git clone https://github.com/YOUR_USERNAME/go-composable-business-types.git`
3. **Install dependencies**: `just install-deps`
4. **Run tests**: `just check`

### 6.2 Development Workflow

```bash
# Create feature branch
git checkout -b feature/descriptive-name

# Make changes
# ... edit files ...

# Run quality checks
just check

# Commit with conventional commits
git commit -m "feat(id): add CompareOrZero method"

# Push and open PR
git push origin feature/descriptive-name
```

### 6.3 Commit Message Convention

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

**Types**:

| Type       | Use For            | Example                                       |
| ---------- | ------------------ | --------------------------------------------- |
| `feat`     | New features       | `feat(id): add CompareOrZero method`          |
| `fix`      | Bug fixes          | `fix(bounded): handle empty string correctly` |
| `docs`     | Documentation      | `docs(readme): update usage examples`         |
| `style`    | Formatting         | `style(types): fix indentation`               |
| `refactor` | Code restructuring | `refactor(actor): simplify chain traversal`   |
| `test`     | Adding tests       | `test(datapoint): add bitemporal fuzz tests`  |
| `chore`    | Maintenance        | `chore(deps): update golangci-lint`           |
| `perf`     | Performance        | `perf(actor): optimize chain append`          |
| `security` | Security fixes     | `security(nanoid): validate alphabet`         |

**Scopes**: Package name (`actor`, `datapoint`, `types`, `enums`, `money`, `locale`, `nanoid`, `bounded`, `temporal`, `importance`, `tag`, `pkg/errors`)

### 6.4 Code Review Process

| Check             | Requirement                             |
| ----------------- | --------------------------------------- |
| **CLA**           | Must sign Contributor License Agreement |
| **Tests**         | All tests must pass (`just test`)       |
| **Coverage**      | Maintain or improve coverage            |
| **Linting**       | Zero lint errors (`just lint`)          |
| **Documentation** | Update relevant docs                    |
| **Reviewers**     | Minimum 1 approval from maintainer      |

### 6.5 Pull Request Template

```markdown
## Summary

Brief description of changes

## Changes

- Change 1
- Change 2

## Testing

- [ ] Unit tests added/updated
- [ ] Integration tests pass
- [ ] Manual testing performed

## Documentation

- [ ] README updated (if needed)
- [ ] CHANGELOG updated
- [ ] Code comments added

## Breaking Changes

- [ ] No breaking changes
- [ ] Breaking changes (explain below)

## Checklist

- [ ] `just check` passes
- [ ] Commits follow conventional format
- [ ] PR title follows convention
```

### 6.6 Code Style

Follow existing patterns:

- **Functional style**: Prefer immutable types with `With*` methods
- **Error handling**: Return errors; use `Must*` only for constants
- **Generics**: Use type parameters for brand types
- **Naming**: Export minimal surface; use descriptive names
- **Documentation**: Every exported type/function must have doc comment

---

## 7. Release Process

### 7.1 Release Checklist

- [ ] All tests pass (`just check`)
- [ ] CHANGELOG.md updated
- [ ] Version bumped in relevant files
- [ ] Documentation updated
- [ ] Git tag created
- [ ] GitHub release created
- [ ] Announcement drafted

### 7.2 Version Bump Locations

| File           | Format         | Example                   |
| -------------- | -------------- | ------------------------- |
| `CHANGELOG.md` | Section header | `## [1.2.0] - 2026-03-27` |
| Git tag        | Semantic       | `v1.2.0`                  |

### 7.3 Release Commands

```bash
# 1. Update CHANGELOG.md
# 2. Commit changes
git add CHANGELOG.md
git commit -m "chore(release): prepare v1.2.0"

# 3. Create tag
git tag -a v1.2.0 -m "Release v1.2.0"

# 4. Push
git push origin main --tags

# 5. Create GitHub release (manual or via CI)
```

### 7.4 Release Types

| Type          | Trigger         | Example                      |
| ------------- | --------------- | ---------------------------- |
| **Patch**     | Bug fix         | `v1.0.0` → `v1.0.1`          |
| **Minor**     | New feature     | `v1.0.0` → `v1.1.0`          |
| **Major**     | Breaking change | `v1.0.0` → `v2.0.0`          |
| **Emergency** | Security fix    | `v1.0.0` → `v1.0.1-security` |

---

## 8. Dependency Management

### 8.1 Dependency Philosophy

**Minimal dependencies**: Prefer standard library; add external deps only when justified.

### 8.2 Current Dependencies

| Package                      | Purpose           | Policy               |
| ---------------------------- | ----------------- | -------------------- |
| `github.com/bojanz/currency` | ISO 4217 currency | Pin to minor version |
| `github.com/sixafter/nanoid` | FIPS-140 NanoID   | Pin to minor version |
| `golang.org/x/text`          | Locale/language   | Pin to minor version |
| `github.com/abice/go-enum`   | Enum generation   | Development only     |

### 8.3 Dependency Update Policy

| Update Type          | Frequency | Approval               |
| -------------------- | --------- | ---------------------- |
| **Patch updates**    | Monthly   | Automated (Dependabot) |
| **Minor updates**    | Quarterly | Manual review          |
| **Major updates**    | As needed | RFC required           |
| **Security updates** | Immediate | Emergency review       |

### 8.4 Vulnerability Scanning

```bash
# Check for vulnerabilities
go list -json -m all | nancy sleuth

# Or use govulncheck
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```

---

## 9. Quality Gates

### 9.1 Pre-Commit Requirements

| Gate      | Command         | Threshold          |
| --------- | --------------- | ------------------ |
| **Build** | `just build`    | Zero errors        |
| **Test**  | `just test`     | 100% pass          |
| **Lint**  | `just lint`     | Zero errors        |
| **Race**  | Built into test | No race conditions |
| **Vet**   | `go vet ./...`  | Zero issues        |

### 9.2 Coverage Requirements

| Metric                | Minimum | Target |
| --------------------- | ------- | ------ |
| **Overall**           | 80%     | 90%    |
| **Critical packages** | 90%     | 95%    |
| **New code**          | 80%     | 90%    |

### 9.3 Performance Benchmarks

| Package      | Benchmark                   | Target     |
| ------------ | --------------------------- | ---------- |
| `actor/`     | `BenchmarkActorChainAppend` | < 100ns/op |
| `nanoid/`    | `BenchmarkNew`              | < 1μs/op   |
| `datapoint/` | `BenchmarkDataPointJSON`    | < 5μs/op   |

### 9.4 Continuous Integration

GitHub Actions runs on every PR:

1. Build with GOEXPERIMENT=jsonv2
2. Test with race detector
3. Lint with golangci-lint
4. Coverage report
5. Benchmark comparison (if benchmarks exist)

---

## 10. Support Policy

### 10.1 Maintenance Commitment

| Version        | Maintenance Level  | Response Time |
| -------------- | ------------------ | ------------- |
| Latest minor   | Active development | 48 hours      |
| Previous minor | Maintenance mode   | 7 days        |
| Older minors   | Best effort        | 30 days       |

### 10.2 Support Channels

| Channel            | Purpose               | Response Time |
| ------------------ | --------------------- | ------------- |
| GitHub Issues      | Bug reports, features | 7 days        |
| GitHub Discussions | Questions, help       | 14 days       |
| Security Email     | Vulnerabilities       | 48 hours      |

### 10.3 Commercial Support

For enterprise support, custom development, or consulting:

**Contact**: enterprise@lars.software

Services available:

- Priority bug fixes
- Custom type development
- Migration assistance
- Architecture review
- Training

---

## 11. License & Legal

### 11.1 License

This library is released under the **MIT License**:

```
MIT License

Copyright (c) 2026 Lars Artmann

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

### 11.2 Third-Party Licenses

Dependencies and their licenses:

| Package                      | License      |
| ---------------------------- | ------------ |
| `github.com/bojanz/currency` | MIT          |
| `github.com/sixafter/nanoid` | Apache-2.0   |
| `golang.org/x/text`          | BSD-3-Clause |
| `github.com/abice/go-enum`   | MIT (dev)    |

### 11.3 Contributor License Agreement

By contributing to this project, you agree that:

1. You have the right to submit the code
2. Your contribution is licensed under the MIT License
3. You grant Lars Artmann a perpetual license to use your contribution

### 11.4 Disclaimer

THIS SOFTWARE IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND. USE AT YOUR OWN RISK.

---

## Appendix A: Policy Change Log

| Date       | Version | Changes                 |
| ---------- | ------- | ----------------------- |
| 2026-03-27 | 1.0.0   | Initial policy document |

---

## Appendix B: Quick Reference

### Breaking Change Checklist

Before releasing a breaking change:

- [ ] Is this absolutely necessary?
- [ ] Has an RFC been opened?
- [ ] Is there a deprecation path?
- [ ] Is migration documented?
- [ ] Will this be a major version bump?

### Release Checklist

- [ ] `just check` passes
- [ ] CHANGELOG updated
- [ ] Version tags created
- [ ] GitHub release published
- [ ] Documentation synced

### Security Checklist

- [ ] No hardcoded secrets
- [ ] Input validated
- [ ] No panics on user input
- [ ] Dependencies scanned
- [ ] Security policy reviewed

---

_This policy is a living document. Updates are announced in CHANGELOG.md and GitHub releases._
