# Comprehensive Status Report — 2026-05-07 19:19

> **Session Focus:** README review → deep reflection → execution → documentation cleanup
> **Previous Status:** [2026-05-07_18-38_public-release-complete.md](./2026-05-07_18-38_public-release-complete.md)

---

## Executive Summary

The library is in **good shape** — all 19 testable packages pass with 86.9% coverage, the codebase is ~10.8K lines of Go, and the architecture is clean. This session focused on **documentation accuracy** and **code quality** after the initial README review revealed cascading issues from the `go-branded-id` extraction.

**Verdict:** Production-ready. Documentation now matches reality. Code has minor DRY and consistency gaps worth addressing.

---

## A. Fully Done ✅

### Code Changes (this session)

| Commit | What | Impact |
|--------|------|--------|
| `61cfca6` | `Email.Validate()` and `URL.Validate()` now delegate to constructors | Eliminated 27 lines of pure duplication |
| `bfd0965` | `Duration.Compare()` uses existing `compare[T]()` generic | Eliminated hand-rolled comparison |
| `9462a5b` | SQL Scan lambdas replaced with direct constructor refs | Linter cleanup (`gocritic/unlambda`) |

### Documentation Changes (this session)

| Commit | What | Impact |
|--------|------|--------|
| `8a819a4` | README.md: fix broken `id/` imports, add missing types/enums/deps | **Critical** — examples were uncompilable |
| `5ac2c8f` | AGENTS.md: complete package tree, add all missing packages | Accurate project map |
| `73d377e` | Add `// Package` docs to 6 packages (importance, tag, version, programminglanguage, projectcore, testutil) | pkg.go.dev now renders these |
| `eeb636c` | POLICY.md: remove all stale `id/` references | No longer claims `id/` is an internal stable package |
| `0d1d451` | PARTS.md: mark `id/` as extracted | Historical accuracy |

### Prior Work (previous sessions, still accurate)

- ✅ `go-branded-id` extraction complete and published
- ✅ Selective import structure (Go 1.26 subpackages)
- ✅ All types have JSON, SQL, validation support
- ✅ `DataPoint[T]` with full audit trail
- ✅ Enum code generation via `go-enum`
- ✅ MIT license (transitioned from PROPRIETARY)
- ✅ CI/CD with GitHub Actions (test, lint, security, release)
- ✅ `flake.nix` build system
- ✅ Examples directory (`basic/`, `datapoint/`)

---

## B. Partially Done ⚠️

| Area | Status | What's Missing |
|------|--------|----------------|
| `Cents` vs `Money` relationship | Undocumented overlap | No guidance on when to use which. `Cents` is currency-unaware; `Money` wraps `bojanz/currency`. Consumers may silently mix USD cents with EUR cents. |
| `validate.Validator` interface consistency | Partially implemented | `Email`, `URL`, `Cents`, `Percentage` implement it. `Timestamp`, `Duration`, `NanoID`, `Locale`, `Tag`, `Importance` do **not**. Either all types should implement it, or it shouldn't be part of the package contract. |
| `emailRegex` double validation | Identified but not fixed | `NewEmail` runs both `mail.ParseAddress` (RFC 5322) AND `emailRegex`. The regex is strictly less capable than the stdlib parser. One or the other suffices. |
| `Timestamp`/`Duration` `Scan()` methods | Not using `scanutil` helpers | `Duration.Scan()` has a hand-rolled 50-line type switch while `Email.Scan()` uses the clean `scanStringType` helper. Inconsistent patterns. |
| `Timestamp`/`Duration` `Compare()` methods | `Duration.Compare()` uses `compare[T]()` ✅ but `Timestamp.Compare()` delegates to `time.Time.Compare()` ✅ — both are correct but the patterns differ. |

---

## C. Not Started ❌

| # | Item | Priority | Effort |
|---|------|----------|--------|
| 1 | Add `Timestamp.Validate()` and `Duration.Validate()` to implement `validate.Validator` | Medium | Tiny |
| 2 | Remove `emailRegex` — trust `mail.ParseAddress` alone, or keep regex and remove `mail.ParseAddress` | Medium | Tiny |
| 3 | Document `Cents` vs `Money` usage guidance in README | Medium | Small |
| 4 | Refactor `Duration.Scan()` and `Timestamp.Scan()` to use `scanutil` helpers | Low | Small |
| 5 | Add `Locale` to README Usage examples | Low | Tiny |
| 6 | Add `Importance` to README Usage examples | Low | Tiny |
| 7 | Add `Tag` to README Usage examples | Low | Tiny |
| 8 | `report/` package — purpose unknown, appears empty or unused | Unknown | Unknown |
| 9 | `BDD_TESTS_REVIEW.md` — stale document from earlier session | Low | Tiny |
| 10 | `MIGRATION_TO_NIX_FLAKES_PROPOSAL.md` — proposal exists but not executed | Low | Large |
| 11 | `PROJECT_SPLIT_EXECUTIVE_REPORT.md` — historical, could be archived | Low | Tiny |

---

## D. Totally Fucked Up 💥

| Item | Severity | Details |
|------|----------|---------|
| **CI billing failure** | 🔴 Critical | GitHub Actions billing is broken — all CI runs fail with billing/spending limit error. This is an account-level issue, not a code problem. No PR can be validated by CI. |
| **`justfile` is deprecated** but still exists | 🟡 Minor | AGENTS.md says "use flake.nix", POLICY.md still references `just check`, `just test`, `just lint`, `just release`. Mixed signals. |
| **`depguard` linter misconfigured** | 🟡 Minor | Multiple packages get `depguard: import is not allowed from list 'Main'` warnings for perfectly valid intra-project imports. Linter config needs updating. |
| **~260 golangci-lint warnings** | 🟡 Moderate | Mostly `revive` (missing doc comments), `err113` (dynamic errors), `recvcheck` (mixed receivers), `gochecknoglobals`. Not blocking, but noisy. |

---

## E. What We Should Improve 🏗️

### Architecture & Design

1. **`validate.Validator` contract consistency** — Either all domain types implement it or we remove the compile-time assertions in `types_sql.go:174-178`. Having 4 of 6 types implement it is a split brain.

2. **`Cents` and `Money` relationship** — `Cents` is a bare `int64` with no currency. `Money` wraps `bojanz/currency.Amount`. Consumers can silently do `Cents(1099)` (USD? EUR?) and pass it alongside `Money` (which knows its currency). Consider documenting guidance or deprecating `Cents` in favor of `Money`.

3. **`emailRegex` belt-and-suspenders** — Running both `mail.ParseAddress` AND a custom regex is either defensive or redundant. Pick one. `mail.ParseAddress` is more correct.

4. **`Scan()` pattern consistency** — `Duration.Scan()` and `Timestamp.Scan()` have 50-line hand-rolled type switches. `Email.Scan()` and `URL.Scan()` use clean `scanStringType` generic helpers. The time types should use `scanutil` too.

### Documentation

5. **All exported symbols need doc comments** — 260+ linter warnings are mostly `revive: exported X should have comment`. This directly impacts pkg.go.dev rendering.

6. **POLICY.md references `just` commands** — Should reference `flake.nix` equivalents or update to match actual toolchain.

7. **Examples are minimal** — `examples/basic/` and `examples/datapoint/` are good starts but don't cover `money`, `locale`, `bounded`, `temporal`, `importance`, `tag`.

### Code Quality

8. **`err113` — dynamic errors** — Multiple packages create errors with `errors.New()` or `fmt.Errorf()` inside functions. Should use sentinel errors from `pkg/errors/`.

9. **Mixed receiver types** — `recvcheck` warns about types with both pointer and value receivers (Email, URL, Locale, Cents, Percentage, Timestamp, Duration). Should standardize.

10. **`testutil/parse.go` has unused parameter** — `RunStringTests` takes a `name` parameter it never uses.

---

## F. Top 25 Things To Do Next

Sorted by impact × effort (Pareto ranking):

### P0 — Do Now (high impact, low effort)

| # | Task | Effort | Impact |
|---|------|--------|--------|
| 1 | **Fix CI billing** — resolve GitHub Actions spending limit | 5min (account settings) | Critical — no CI = no safety net |
| 2 | **Add `Validate()` to `Timestamp` and `Duration`** | 5min | Consistency — completes `validate.Validator` contract |
| 3 | **Remove or justify `emailRegex`** | 5min | DRY — eliminates redundant double validation |
| 4 | **Add `CauseKind` to `enums` table in PARTS.md** | 2min | Accuracy — already in code, just missing from analysis |

### P1 — Do Soon (good impact, moderate effort)

| # | Task | Effort | Impact |
|---|------|--------|--------|
| 5 | **Document `Cents` vs `Money` guidance in README** | 15min | Prevents consumer confusion |
| 6 | **Add usage examples for `Locale`, `Importance`, `Tag` in README** | 20min | Discoverability — 3 types completely unexampled |
| 7 | **Refactor `Duration.Scan()` and `Timestamp.Scan()` to use `scanutil`** | 30min | DRY — eliminates 80 lines of hand-rolled type switches |
| 8 | **Fix `depguard` linter config** — allow intra-project imports | 10min | Eliminates ~20 false-positive warnings |
| 9 | **Add `Locale.Validate()`, `Tag.Validate()` consistency** | 10min | All types should implement `validate.Validator` |
| 10 | **Update POLICY.md to reference `flake.nix` instead of `just`** | 15min | Eliminates tooling confusion |

### P2 — Do Eventually (nice to have)

| # | Task | Effort | Impact |
|---|------|--------|--------|
| 11 | **Add doc comments to all exported symbols** (~50 items) | 2hr | pkg.go.dev rendering, eliminates ~100 linter warnings |
| 12 | **Replace dynamic errors with sentinels from `pkg/errors/`** | 2hr | Consistent error handling, eliminates `err113` warnings |
| 13 | **Standardize receiver types** (value vs pointer) | 1hr | Code consistency, eliminates `recvcheck` warnings |
| 14 | **Add `Money` example to `examples/`** | 30min | Most complex type, deserves dedicated example |
| 15 | **Add `BoundedString` example to `examples/`** | 15min | Common type, no example exists |
| 16 | **Add `Temporal` example to `examples/`** | 15min | Complex type, no example exists |
| 17 | **Archive `BDD_TESTS_REVIEW.md`** to `docs/status/archive/` | 1min | Housekeeping |
| 18 | **Archive `PROJECT_SPLIT_EXECUTIVE_REPORT.md`** to `docs/` | 1min | Housekeeping |
| 19 | **Investigate `report/` package** — what is it? Used? | 10min | Either document or remove |
| 20 | **Fix `testutil/parse.go` unused `name` parameter** | 2min | Eliminates linter warning |
| 21 | **Add `Locale` SQL round-trip tests** | 15min | `Locale.Scan()`/`Value()` lack direct test coverage |
| 22 | **Add `Importance` SQL round-trip tests** | 15min | Same as above |
| 23 | **Consider adding `encoding.TextAppender` to types** (Go 1.24+) | 1hr | Modern Go patterns, already identified in earlier analysis |
| 24 | **Consider adding `encoding.BinaryAppender` to types** (Go 1.24+) | 1hr | Modern Go patterns |
| 25 | **Consider `MIGRATION_TO_NIX_FLAKES_PROPOSAL.md` execution** | 4hr | Build system unification |

---

## G. Top #1 Question I Cannot Figure Out Myself 🤔

**Is `report/` package intentionally empty/unused, or is it leftover from a previous plan?**

The directory exists with no `.go` files visible. It doesn't compile as a package. It's not referenced anywhere in the codebase. Should it be removed, or is it a placeholder for future work?

---

## Package Coverage Matrix

| Package | Coverage | Status |
|---------|----------|--------|
| `actor/` | **100.0%** | ✅ Excellent |
| `nanoid/` | **100.0%** | ✅ Excellent |
| `programminglanguage/` | **100.0%** | ✅ Excellent |
| `validate/` | **100.0%** | ✅ Excellent |
| `enums/` | **98.9%** | ✅ Excellent |
| `bounded/` | **97.8%** | ✅ Excellent |
| `scanutil/` | **97.1%** | ✅ Excellent |
| `importance/` | **97.6%** | ✅ Excellent |
| `tag/` | **93.4%** | ✅ Good |
| `temporal/` | **95.1%** | ✅ Good |
| `datapoint/` | **90.1%** | ✅ Good |
| `money/` | **90.0%** | ✅ Good |
| `types/` | **89.8%** | ✅ Good |
| `locale/` | **88.1%** | 🟡 Adequate |
| `pkg/errors/` | **87.5%** | 🟡 Adequate |
| `projectcore/` | **84.4%** | 🟡 Adequate |
| `version/` | **81.0%** | 🟡 Adequate |
| **Total** | **86.9%** | ✅ Above 80% target |

## Build & Test Status

- ✅ `go build ./...` — passes
- ✅ `go test -race ./...` — all 19 packages pass
- ✅ `go vet ./...` — passes
- ⚠️ `golangci-lint run` — ~260 warnings (mostly doc comments, non-blocking)

## Key Metrics

| Metric | Value |
|--------|-------|
| Go packages | 18 (excluding examples/testutil) |
| Total Go LOC | ~10,820 |
| Test coverage | 86.9% |
| Dependencies | 4 direct (`go-branded-id`, `bojanz/currency`, `sixafter/nanoid`, `golang.org/x/text`) + 1 dev (`abice/go-enum`) |
| License | MIT |
| Go version | 1.26.2 |

---

_Report generated: 2026-05-07 19:19_
_Session commits: 8 (61cfca6..8a819a4)_
