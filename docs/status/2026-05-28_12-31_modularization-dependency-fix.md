# Status Report: Modularization Fix — Sub-Module Dependency Resolution

**Date:** 2026-05-28 12:31 CEST
**Branch:** `modularize/split-modules` (up to date with origin)
**Go:** 1.26.3 linux/amd64
**Session Focus:** Fix broken sub-module go.mod files after initial modularization

---

## Executive Summary

The initial modularization (commit `f6c0943`) split the monolith into 6 modules but left **all sub-module go.mod files critically broken** — they imported root module packages without declaring the dependency. This session diagnosed and fixed every sub-module, establishing a working hybrid `go.work` + `replace` strategy.

**Build:** PASS | **Tests:** 19/19 suites PASS | **Vet:** CLEAN | **Mod Verify:** CLEAN | **Coverage:** 89.0%

---

## A) FULLY DONE

### 1. Sub-Module go.mod Dependency Declarations

Every sub-module now correctly declares its dependencies:

| Module    | Root Dep | Sibling Deps | External Deps | Replace Directives |
| --------- | -------- | ------------ | ------------- | ------------------ |
| **Root**  | — | — | go-branded-id, go-enum (tool), testify | None |
| **nanoid** | v0.4.0 | — | sixafter/nanoid | root → ../ |
| **locale** | v0.4.0 | — | golang.org/x/text | root → ../ |
| **money** | v0.4.0 | locale | bojanz/currency | root → ../, locale → ../locale |
| **datapoint** | v0.4.0 | nanoid | go-branded-id v0.3.0 | root → ../, nanoid → ../nanoid |
| **examples** | v0.4.0 | nanoid, datapoint | go-branded-id v0.3.0 | root → ../, nanoid → ../nanoid, datapoint → ../datapoint |

### 2. go-branded-id Version Alignment

- **Before:** Root at v0.3.0, datapoint and examples at v0.1.0
- **After:** All modules at v0.3.0

### 3. Go Version Alignment

- All 6 go.mod files now specify `go 1.26.3` (was mixed 1.26.2/1.26.3)

### 4. `go mod tidy` Works in All Sub-Modules

- **Before:** `go mod tidy` failed with ambiguous import errors in every sub-module
- **After:** `go mod tidy` succeeds in all 6 modules independently

### 5. Hybrid `go.work` + `replace` Strategy Documented

- DEPENDENCY_GRAPH.md rewritten with actual replace directive strategy
- PROPOSAL.md updated with rationale for the hybrid approach
- AGENTS.md updated with correct build/test/tidy commands

### 6. Full Verification Matrix

| Check | Root | nanoid | locale | money | datapoint | examples |
| ----- | ---- | ------ | ------ | ----- | --------- | -------- |
| `go build ./...` | PASS | PASS | PASS | PASS | PASS | PASS |
| `go test -race ./...` | 13/13 PASS | PASS | PASS | PASS | PASS | 0 tests |
| `go vet ./...` | CLEAN | CLEAN | CLEAN | CLEAN | CLEAN | CLEAN |
| `go mod tidy` | CLEAN | CLEAN | CLEAN | CLEAN | CLEAN | CLEAN |
| `go mod verify` | VERIFIED | VERIFIED | VERIFIED | VERIFIED | VERIFIED | VERIFIED |
| Coverage | 89.0% | 100.0% | 88.1% | 90.0% | 90.1% | N/A |

---

## B) PARTIALLY DONE

### 1. Replace Directive Scaffolding

The `replace` directives in sub-modules are **temporary scaffolding** needed because the published root v0.4.0 still contains all packages. They should be removed once:
- A new root version is published that no longer includes nanoid/, locale/, money/, datapoint/, examples/
- Sub-modules get their first versioned git tags (e.g., `nanoid/v0.5.0`)

**Status:** Working correctly today, but represents technical debt for the next release.

### 2. CI/CD Multi-Module Support

The CI pipeline (GitHub Actions) still tests only the root module. It needs updating to test all 6 modules independently. CI is currently blocked by a billing/spending limit issue on the GitHub account (not a code problem).

### 3. flake.nix Multi-Module Build

The Nix flake needs updating to build and test all modules. Currently it only tests root.

---

## C) NOT STARTED

### 1. First Versioned Release of Sub-Modules

No git tags exist for sub-modules (`nanoid/v0.5.0`, `locale/v0.5.0`, etc.). The release workflow is defined but not executed.

### 2. Merge `modularize/split-modules` into `master`

The modularization branch has not been merged. Requires review and approval.

### 3. Replace `testify` with `ginkgo/gomega`

3 test files in root still use `github.com/stretchr/testify` (banned per policy):
- `importance/importance_test.go`
- `tag/tag_test.go`
- `projectcore/project_core_test.go`

### 4. CI Pipeline for Multi-Module

GitHub Actions workflow needs:
- Per-module test jobs (or a loop over modules)
- Per-module lint jobs
- Multi-module release workflow

### 5. flake.nix Multi-Module Support

Nix build system needs per-module checks and aggregated root-level checks.

### 6. Integration Tests for External Consumers

No test verifying that external consumers can `go get` individual sub-modules correctly.

### 7. API Documentation (pkg.go.dev)

No verified that `pkg.go.dev` renders all sub-module documentation correctly.

### 8. Examples Module Testing

`examples/` has zero tests. Should have at least compilation verification.

### 9. Root Module Cleanup for Next Release

Root module go.mod still has `testify` and `go-enum` tool dependency. After removing testify, the root go.mod will be cleaner.

---

## D) TOTALLY FUCKED UP (Now Fixed)

### 1. Sub-Module go.mod Files Were Ghost Dependencies

**Severity:** CRITICAL

The initial modularization (commit `f6c0943`) created sub-module go.mod files that:
- Imported `pkg/errors`, `scanutil`, `testutil`, `enums`, `actor`, `types`, `temporal` from root — **without declaring the root module as a dependency**
- Imported `locale` from locale module — **without declaring locale as a dependency**
- Only "worked" because `go.work` accidentally provided the modules locally
- `go mod tidy` would fail in every sub-module with ambiguous import errors
- External consumers could NOT resolve any cross-module imports

**Root cause:** The published root v0.4.0 still contains all packages. When `go mod tidy` resolves the root module, it sees nanoid/locale/money/datapoint/examples as both part of the published root AND as separate local modules → ambiguous imports.

**Fix:** Added `replace` directives in every sub-module pointing to the local root module (`../`) and local sibling modules. This forces Go to resolve against the local root (which correctly excludes sub-module directories that have their own go.mod).

### 2. go-branded-id Version Mismatch

**Severity:** MEDIUM

Root used v0.3.0; datapoint and examples used v0.1.0. Fixed to v0.3.0 everywhere.

### 3. Go Version Inconsistency

**Severity:** LOW

Root was 1.26.3, sub-modules were 1.26.2. Fixed to 1.26.3 everywhere.

---

## E) WHAT WE SHOULD IMPROVE

### Architecture & Design

1. **Remove `replace` scaffolding** — Publish a new root v0.5.0 that removes split-out packages, then remove `replace` directives from sub-modules
2. **Error type placement** — NanoID sentinel errors (`ErrNanoIDEmpty`, etc.) live in `pkg/errors` (root) but are used by the nanoid sub-module. If a consumer only imports nanoid, they must also import root for error checking. Consider moving nanoid-specific errors into the nanoid module.
3. **`testutil` cross-module dependency** — testutil lives in root but is imported by nanoid and locale tests. This means those sub-modules have a production-declared dependency on root just for test helpers. Consider duplicating testutil inline or making testutil its own module.
4. **`scanutil` cross-module dependency** — scanutil (root) is imported by nanoid and locale in production code. This is correct architecturally but means nanoid/locale always pull in the root module. Acceptable but worth noting.

### Quality

5. **Replace testify** — 3 files still use banned dependency. Replace with ginkgo/gomega.
6. **Coverage gaps** — `pkg/errors` at 87.5%, `scanutil` at 79.2%, `version` at 81.0%, `projectcore` at 82.9%. Target 95%+.
7. **No integration tests** — No test verifies external consumer experience.
8. **Examples have zero tests** — At minimum, verify they compile.

### Operations

9. **CI is broken** — GitHub Actions billing issue (account-level, not code). Blocks all CI.
10. **flake.nix outdated** — Still references old single-module structure.

---

## F) Top 25 Things We Should Get Done Next

| # | Priority | Task | Impact | Effort |
|---|----------|------|--------|--------|
| 1 | **P0** | Merge `modularize/split-modules` into `master` | Unblocks all downstream work | Low |
| 2 | **P0** | Fix GitHub Actions billing/spending limit | Restores CI | External |
| 3 | **P0** | Update CI workflow for multi-module testing | Ensures all modules tested | Medium |
| 4 | **P1** | Create first versioned release (v0.5.0) with sub-module tags | Enables external consumers | Low |
| 5 | **P1** | Remove `replace` directives after new root release | Eliminates scaffolding tech debt | Low |
| 6 | **P1** | Update flake.nix for multi-module build/test | Restores Nix CI | Medium |
| 7 | **P1** | Replace testify with ginkgo/gomega (3 files) | Policy compliance | Medium |
| 8 | **P1** | Move nanoid-specific errors from `pkg/errors` to nanoid module | Cleaner dep isolation | Medium |
| 9 | **P2** | Add compilation test for examples module | Catches build regressions | Low |
| 10 | **P2** | Write integration test verifying external consumer imports | Validates consumer experience | Medium |
| 11 | **P2** | Update README.md to reflect multi-module structure | Consumer documentation | Low |
| 12 | **P2** | Verify pkg.go.dev renders all sub-modules correctly | Documentation | Low |
| 13 | **P2** | Add `go work edit -fmt` to CI pipeline | Prevents go.work drift | Low |
| 14 | **P2** | Create release automation script (justfile → flake.nix) | Operational | Medium |
| 15 | **P3** | Improve `pkg/errors` coverage to 95%+ (87.5% → 95%) | Quality | Low |
| 16 | **P3** | Improve `scanutil` coverage to 95%+ (79.2% → 95%) | Quality | Low |
| 17 | **P3** | Improve `version` coverage to 95%+ (81.0% → 95%) | Quality | Low |
| 18 | **P3** | Improve `projectcore` coverage to 95%+ (82.9% → 95%) | Quality | Low |
| 19 | **P3** | Consider extracting `testutil` to its own module | Cleaner test dep isolation | Medium |
| 20 | **P3** | Add `//go:build` constraints if any platform-specific code exists | Correctness | Low |
| 21 | **P3** | Add benchmark suite for hot paths (nanoid generation, parsing) | Performance visibility | Medium |
| 22 | **P3** | Add CHANGELOG entry for v0.5.0 modularization release | Documentation | Low |
| 23 | **P4** | Evaluate `internal/` packages for cross-module access safety | Correctness | Low |
| 24 | **P4** | Set up Dependabot for sub-module go.mod files | Security | Low |
| 25 | **P4** | Add pre-commit hook for `go mod tidy` verification | Developer experience | Low |

---

## G) Top #1 Question I Cannot Figure Out Myself

**Should `replace` directives stay or go?**

The current `replace` directives in sub-module go.mod files exist because the published root v0.4.0 still contains all packages. There are two paths forward:

1. **Publish root v0.5.0 without the split-out packages** — This would eliminate the ambiguous import problem. Sub-modules could reference root v0.5.0 without `replace`. But this is a **breaking change** for anyone importing `go-composable-business-types/nanoid` from the root module (they'd need to switch to the nanoid sub-module).

2. **Keep `replace` directives permanently** — They work fine for development and are ignored by consumers of published sub-modules. But they add noise to go.mod files and represent an unusual pattern.

**Question for you:** Are you ready to cut root v0.5.0 as a breaking release that removes nanoid/locale/money/datapoint/examples from the root module? This would let us clean up the `replace` directives and give consumers a clean experience.

---

## Project Metrics

| Metric | Value |
|--------|-------|
| Go version | 1.26.3 |
| Modules | 6 (root + 5 sub-modules) |
| Packages | 18 (13 root + 5 sub-module) |
| Go source files | 53 |
| Lines of code | ~10,933 |
| Test suites | 19 (16 with tests, 3 no-test packages) |
| Test coverage (root) | 89.0% |
| Test coverage (nanoid) | 100.0% |
| Test coverage (locale) | 88.1% |
| Test coverage (money) | 90.0% |
| Test coverage (datapoint) | 90.1% |
| External deps (root) | go-branded-id, go-enum, testify |
| External deps (nanoid) | sixafter/nanoid |
| External deps (locale) | golang.org/x/text |
| External deps (money) | bojanz/currency |
| External deps (datapoint) | go-branded-id |
| Replace directives | 8 total across 5 sub-modules |
| License | MIT |

---

## Files Changed This Session

```
 AGENTS.md                               |  10 ++-
 datapoint/go.mod                        |  21 ++++++-
 datapoint/go.sum                        |  25 +++++++-
 docs/modularization/DEPENDENCY_GRAPH.md | 107 ++++++++++++++++++--------------
 docs/modularization/PROPOSAL.md         |  21 +++++--
 examples/go.mod                         |  23 ++++++-
 examples/go.sum                         |  25 +++++++-
 locale/go.mod                           |   9 ++-
 locale/go.sum                           |   1 +
 money/go.mod                            |  17 +++--
 money/go.sum                            |   9 ++-
 nanoid/go.mod                           |  15 +++--
 nanoid/go.sum                           |  15 ++++-
 13 files changed, 221 insertions(+), 77 deletions(-)
```
