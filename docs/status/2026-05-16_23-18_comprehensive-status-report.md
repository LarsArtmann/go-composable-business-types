# Comprehensive Status Report — go-composable-business-types

**Date:** 2026-05-16 23:18 CEST
**Commit:** `c9bda50` (master)
**Go Version:** 1.26.2

---

## Executive Summary

The library is in **good shape** — 86.8% test coverage, zero lint issues, clean build, all tests pass with race detector. The `programminglanguage/` package was just removed this session in favor of recommending `go-enry`. The codebase is ~11,600 lines across 53 Go files (26 test files), with 17 packages.

**Verdict: Solid foundation, needs cleanup of stale docs/config and a v1.0.0 release planning.**

---

## Project Metrics

| Metric                | Value                                                                                              |
| --------------------- | -------------------------------------------------------------------------------------------------- |
| Total Go files        | 53                                                                                                 |
| Test files            | 26                                                                                                 |
| Total Go lines        | ~11,631                                                                                            |
| Packages              | 17 (incl. examples, testutil, pkg/errors)                                                          |
| Test coverage         | **86.8%**                                                                                          |
| Lint issues           | **0**                                                                                              |
| Build status          | **PASS**                                                                                           |
| Race conditions       | **0**                                                                                              |
| Go version            | 1.26.2                                                                                             |
| Dependencies (direct) | 5 (`go-branded-id`, `bojanz/currency`, `sixafter/nanoid`, `stretchr/testify`, `golang.org/x/text`) |
| Dependencies (dev)    | 1 (`abice/go-enum`)                                                                                |
| License               | MIT                                                                                                |
| CI                    | GitHub Actions (billing broken — account-level issue)                                              |

---

## Per-Package Coverage

| Package        | Coverage   | Status    |
| -------------- | ---------- | --------- |
| `actor/`       | **100.0%** | Excellent |
| `enums/`       | **100.0%** | Excellent |
| `nanoid/`      | **100.0%** | Excellent |
| `validate/`    | **100.0%** | Excellent |
| `bounded/`     | **97.8%**  | Excellent |
| `importance/`  | **97.6%**  | Excellent |
| `scanutil/`    | **97.1%**  | Excellent |
| `temporal/`    | **95.1%**  | Great     |
| `tag/`         | **93.4%**  | Great     |
| `datapoint/`   | **90.1%**  | Good      |
| `money/`       | **90.0%**  | Good      |
| `types/`       | **89.9%**  | Good      |
| `locale/`      | **88.1%**  | Good      |
| `pkg/errors/`  | **87.5%**  | Good      |
| `projectcore/` | **82.9%**  | OK        |
| `version/`     | **81.0%**  | OK        |
| `examples/`    | 0%         | N/A       |
| `testutil/`    | 0%         | N/A       |

---

## A) FULLY DONE ✅

| What                           | Details                                                                          |
| ------------------------------ | -------------------------------------------------------------------------------- |
| Core type system               | `ID[B,V]` extracted to `go-branded-id`, all types use branded IDs                |
| Domain primitives              | Email, URL, Percentage, Cents, Timestamp, Duration — all with JSON/SQL support   |
| BoundedString                  | Full validation, factory pattern, SQL interfaces                                 |
| Money                          | ISO 4217 via `bojanz/currency`, locale formatting                                |
| Locale                         | BCP 47 via `golang.org/x/text`, common constants                                 |
| NanoID                         | FIPS-140 via `sixafter/nanoid`, validation, SQL interfaces                       |
| Enums                          | Generated via `go-enum` — ActorKind, Priority, Status, Trigger, CauseKind        |
| ActorChain[T]                  | Generic audit trail chain with origin/current/filter                             |
| DataPoint[T]                   | Self-contained data unit with full audit trail, bitemporal, references, causes   |
| Bitemporal                     | ValidFrom/ValidUntil/Recorded with point-in-time queries                         |
| Importance                     | 0-100 with named levels (Low/Medium/High/Critical)                               |
| Tag                            | Validated string labels with alphanumeric+hyphen constraint                      |
| ProjectCore                    | Composite project metadata (name, path, languages, importance, tags)             |
| Error package                  | Sentinel errors + structured errors with Is/As support                           |
| SQL support                    | Scanner/Valuer for all domain types via `scanutil/`                              |
| JSON support                   | Custom Marshal/Unmarshal for all domain types                                    |
| Validation                     | `validate.Validator` interface, self-validating types                            |
| Lint                           | Zero issues with 90+ golangci-lint linters enabled                               |
| CI pipeline                    | Full GitHub Actions: test (3 Go versions), lint, security, generate, benchmark   |
| Release pipeline               | git-cliff changelog, SemVer tags, GitHub Releases                                |
| Code formatting                | gofumpt, golines, gci via golangci-lint formatters                               |
| Depguard                       | Strict import allowlist — only approved dependencies                             |
| Open source                    | MIT license, SECURITY.md, POLICY.md, SUPPORT.md, issue/PR templates              |
| `programminglanguage/` removal | Done this session — replaced with `[]string` in projectcore, go-enry recommended |
| PARTS.md                       | Component analysis with extraction recommendations                               |

---

## B) PARTIALLY DONE 🔶

| What                       | Status                                               | What's Missing                                                                                                       |
| -------------------------- | ---------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| **CHANGELOG.md**           | Structure exists but only has `v0.1.0` (placeholder) | Needs real entries for all changes since Jan 2026, git-cliff generation                                              |
| **Examples**               | `examples/basic/` and `examples/datapoint/` exist    | No examples for money, bounded, actor, temporal, locale, importance, tag                                             |
| **flake.nix**              | Working devShell + checks                            | Missing: formatter check, lint app, test app, benchmark app; justfile still exists alongside it                      |
| **go-enry recommendation** | Added to projectcore docs and AGENTS.md              | Not mentioned in README.md as a recommended companion library                                                        |
| **GOEXPERIMENT=jsonv2**    | Set in flake.nix, CI, justfile                       | Not actually using `encoding/json/v2` in any code — still on `encoding/json` v1. The flag is set but has zero effect |
| **POLICY.md**              | Comprehensive, well-structured                       | References `just` commands throughout; should reference `flake.nix` or `go` directly per AGENTS.md convention        |
| **docs/ecosystem/**        | Extensive planning docs from May 4 migration         | Many references to `programminglanguage/` are now stale                                                              |

---

## C) NOT STARTED ⬜

| What                            | Priority | Notes                                                          |
| ------------------------------- | -------- | -------------------------------------------------------------- |
| **v1.0.0 release**              | HIGH     | Library is stable enough. Needs CHANGELOG, tag, GitHub release |
| **Adopt encoding/json/v2**      | MED      | GOEXPERIMENT=jsonv2 is set but v2 API not actually used        |
| **Benchmarks for all packages** | MED      | CI runs benchmarks but few packages have them                  |
| **godoc/pkg.go.dev polish**     | MED      | Every exported type needs example functions                    |
| **Fuzz tests**                  | LOW      | Only `nanoid/` has fuzz testing historically                   |
| **BDD tests (ginkgo)**          | LOW      | BDD_TESTS_REVIEW.md exists but no ginkgo tests written         |
| **go-bitemporal extraction**    | LOW      | Identified in PARTS.md as strong standalone candidate          |
| **go-actor-chain extraction**   | LOW      | Identified in PARTS.md as standalone candidate                 |
| **SQL repository helpers**      | LOW      | scanutil provides Scanner/Valuer but no query builders         |
| **OpenTelemetry integration**   | LOW      | ActorChain could bridge to OTEL traces                         |

---

## D) TOTALLY FUCKED UP 💥

| Issue                                               | Severity | Details                                                                                                                                                                                       |
| --------------------------------------------------- | -------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **CI billing broken**                               | HIGH     | All GitHub Actions runs fail with billing/spending limit error. Account-level issue, not code. Blocks all CI/CD.                                                                              |
| **Stale `programminglanguage/` in `.golangci.yml`** | LOW      | Line ~136 still has `- path: programminglanguage/` exclusion for a deleted package. Harmless but messy.                                                                                       |
| **`justfile` coexists with `flake.nix`**            | MED      | AGENTS.md says "justfile is deprecated" and "never create new justfiles", but it still exists and POLICY.md/SUPPORT.md reference `just` commands. Conflicting guidance.                       |
| **GOEXPERIMENT=jsonv2 is theater**                  | MED      | The flag is set everywhere (flake.nix, CI, justfile) but no code uses `encoding/json/v2`. The v1 `encoding/json` is used throughout. This is misleading — either adopt v2 or remove the flag. |
| **~50 stale status/planning docs**                  | LOW      | `docs/status/` has 50+ historical status reports. Some reference deleted packages, old architectures, completed work. Consider archiving or cleaning.                                         |
| **`docs/` references to `programminglanguage/`**    | LOW      | 10+ markdown files still reference the deleted package. Mostly historical/planning docs so not blocking, but confusing.                                                                       |

---

## E) WHAT WE SHOULD IMPROVE 🔧

### High Impact

1. **Fix CI billing** — Without working CI, there's no automated quality gate. This is an account-level issue that only Lars can resolve.
2. **Remove `justfile` or migrate all commands to `flake.nix`** — Having both with conflicting documentation is confusing. Pick one.
3. **Either adopt `encoding/json/v2` or remove the GOEXPERIMENT flag** — Currently it's set but unused, which is misleading.
4. **Generate real CHANGELOG.md** — Run `git-cliff` to produce actual changelog entries. Current CHANGELOG is a placeholder.
5. **Clean `.golangci.yml`** — Remove `programminglanguage/` exclusion, audit other stale entries.

### Medium Impact

6. **Update POLICY.md and SUPPORT.md** — Replace `just` command references with direct `go` commands or `nix` equivalents.
7. **Add go-enry to README.md** — As a recommended companion library for language detection.
8. **Add more examples** — Missing examples for: money, bounded, actor, temporal, locale, importance, tag, projectcore.
9. **Increase coverage on low packages** — `projectcore/` (82.9%), `version/` (81.0%) are below the 85% target.
10. **Write benchmark tests** — Currently only some packages have benchmarks. CI runs them but there's little to compare against.

### Low Impact

11. **Archive old status docs** — Move 40+ historical status reports to `docs/status/archive/` (some already are, but many aren't).
12. **Clean stale docs** — Update or remove references to deleted `programminglanguage/` in ecosystem/planning docs.
13. **Remove `.auto-deduplicate/`** — Contains old false-positive entries from a dedup tool; paths reference macOS (`/Users/larsartmann/`).
14. **Remove `report/` directory** — Empty directory at project root.
15. **Remove `BDD_TESTS_REVIEW.md`** — Review document, no BDD tests exist or are planned.
16. **Remove `MIGRATION_TO_NIX_FLAKES_PROPOSAL.md`** — If flake.nix is already adopted, the proposal is historical.
17. **Remove `PROJECT_SPLIT_EXECUTIVE_REPORT.md`** — Historical planning doc, now completed.

---

## F) TOP 25 THINGS TO DO NEXT

| #   | Task                                                                | Impact   | Effort      | Category       |
| --- | ------------------------------------------------------------------- | -------- | ----------- | -------------- |
| 1   | Fix GitHub Actions billing                                          | CRITICAL | 5min        | Infrastructure |
| 2   | Remove `justfile` (migrate to flake.nix or drop)                    | HIGH     | 30min       | Cleanup        |
| 3   | Clean `.golangci.yml` — remove `programminglanguage/` exclusion     | HIGH     | 5min        | Cleanup        |
| 4   | Run `git-cliff` to generate real CHANGELOG.md                       | HIGH     | 15min       | Release prep   |
| 5   | Decide on `encoding/json/v2`: adopt it or remove GOEXPERIMENT flag  | HIGH     | 2hr or 5min | Technical debt |
| 6   | Update POLICY.md — replace `just` refs with `go`/`nix` commands     | MED      | 30min       | Docs           |
| 7   | Update SUPPORT.md — replace `just` refs                             | MED      | 15min       | Docs           |
| 8   | Add go-enry recommendation to README.md                             | MED      | 5min        | Docs           |
| 9   | Add example for `money/` package                                    | MED      | 30min       | Docs           |
| 10  | Add example for `actor/` package                                    | MED      | 30min       | Docs           |
| 11  | Add example for `datapoint/` with bitemporal tracking               | MED      | 30min       | Docs           |
| 12  | Increase `projectcore/` test coverage to 85%+                       | MED      | 30min       | Quality        |
| 13  | Increase `version/` test coverage to 85%+                           | MED      | 30min       | Quality        |
| 14  | Add benchmark tests for all core packages                           | MED      | 2hr         | Quality        |
| 15  | Update PARTS.md — remove `programminglanguage/`, mark as removed    | LOW      | 10min       | Docs           |
| 16  | Clean stale `programminglanguage/` refs in docs/ecosystem/          | LOW      | 15min       | Docs           |
| 17  | Archive old status reports to `docs/status/archive/`                | LOW      | 15min       | Housekeeping   |
| 18  | Remove empty `report/` directory                                    | LOW      | 1min        | Housekeeping   |
| 19  | Remove `.auto-deduplicate/false-positives.json` (stale macOS paths) | LOW      | 1min        | Housekeeping   |
| 20  | Remove `BDD_TESTS_REVIEW.md` (no BDD tests planned)                 | LOW      | 1min        | Housekeeping   |
| 21  | Remove `MIGRATION_TO_NIX_FLAKES_PROPOSAL.md` (already migrated)     | LOW      | 1min        | Housekeeping   |
| 22  | Remove `PROJECT_SPLIT_EXECUTIVE_REPORT.md` (completed)              | LOW      | 1min        | Housekeeping   |
| 23  | Add fuzz tests for `bounded/`, `types/`, `money/`                   | LOW      | 2hr         | Quality        |
| 24  | Plan v1.0.0 release — API stability review, migration guide         | MED      | 3hr         | Strategy       |
| 25  | Add `go` doc examples (example functions) for key types             | LOW      | 3hr         | Docs           |

---

## G) TOP #1 QUESTION I CANNOT ANSWER MYSELF 🤔

**What is the plan for `GOEXPERIMENT=jsonv2`?**

The flag is set in flake.nix, CI workflows, and justfile — but no code in the entire project actually imports or uses `encoding/json/v2`. Everything uses standard `encoding/json` (v1). This creates a confusing situation:

- **Option A:** Adopt `encoding/json/v2` — requires rewriting all MarshalJSON/UnmarshalJSON methods to use the v2 API. This is a significant effort but would future-proof the library. The v2 API has different semantics (no `interface{}` unmarshaling, different error handling, `json.Value` type).
- **Option B:** Remove the flag entirely — it's theater right now. The v1 API works fine and is stable. Remove `GOEXPERIMENT=jsonv2` from all config files.
- **Option C:** Keep it as-is — the flag is harmless (v1 json still works with the flag set), and it might enable v2 opt-in for consumers. But it's misleading.

This is a **strategic decision** that affects API design going forward, especially for a potential v1.0.0 release. Only Lars can make this call.

---

_Generated by Crush — 2026-05-16 23:18 CEST_
