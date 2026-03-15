# Comprehensive Execution Plan

**Created:** 2026-03-15 16:00 CET
**Project:** go-composable-business-types
**Goal:** Production-ready library with high test coverage, clean code, and excellent documentation

---

## Current State Summary

| Metric | Value | Target |
|--------|-------|--------|
| Build Status | ✅ Pass (after json fix) | ✅ Pass |
| Test Status | ✅ Pass | ✅ Pass |
| Test Coverage | 34.9% | 80%+ |
| Linter Issues | ~300+ | 0 |
| Package Docs | 50% | 100% |
| Examples | 2 | 4+ |

---

## Task Categories

### P0 - CRITICAL (Must fix now - breaks build/tests)
### P1 - HIGH (Important for quality)
### P2 - MEDIUM (Good to have)
### P3 - LOW (Nice to have)

---

## COMPREHENSIVE TASK LIST

| # | Category | Task | Impact | Effort | Value | Est. |
|---|----------|------|--------|--------|-------|------|
| **PHASE 1: STABILIZATION** |
| 1 | P0 | Fix json/v2 imports → json (id/id.go) | Critical | 1m | High | 1m |
| 2 | P0 | Fix json/v2 imports → json (temporal/temporal.go) | Critical | 1m | High | 1m |
| 3 | P0 | Fix json/v2 imports → json (bounded/bounded.go) | Critical | 1m | High | 1m |
| 4 | P0 | Verify build passes with go build ./... | Critical | 2m | High | 2m |
| 5 | P0 | Verify tests pass with go test -race ./... | Critical | 3m | High | 3m |
| 6 | P0 | Commit json/v2 → json fixes | Critical | 2m | High | 2m |
| **PHASE 2: LINTER CONFIG** |
| 7 | P1 | Fix golangci.yml depguard pattern (add /**) | High | 1m | High | 1m |
| 8 | P1 | Run golangci-lint to verify config works | High | 2m | Medium | 2m |
| 9 | P1 | Commit golangci.yml fix | High | 2m | Medium | 2m |
| **PHASE 3: PACKAGE DOCUMENTATION** |
| 10 | P1 | Add package comment to enums/enums.go | High | 2m | Medium | 2m |
| 11 | P1 | Add package comment to money/money.go | High | 2m | Medium | 2m |
| 12 | P1 | Add package comment to locale/locale.go | High | 2m | Medium | 2m |
| 13 | P1 | Add package comment to temporal/temporal.go | High | 2m | Medium | 2m |
| 14 | P1 | Add package comment to datapoint/datapoint.go | High | 2m | Medium | 2m |
| 15 | P1 | Verify godoc renders correctly | Medium | 3m | Medium | 3m |
| 16 | P1 | Commit package documentation | High | 2m | Medium | 2m |
| **PHASE 4: TEST COVERAGE - TYPES PACKAGE (25.9% → 60%)** |
| 17 | P1 | Add tests for Percentage.Compare | High | 3m | High | 3m |
| 18 | P1 | Add tests for Cents.Compare | High | 3m | High | 3m |
| 19 | P1 | Add tests for Timestamp.Compare | High | 3m | High | 3m |
| 20 | P1 | Add tests for Duration.Compare | High | 3m | High | 3m |
| 21 | P1 | Add tests for Duration.Scan (all cases) | High | 5m | High | 5m |
| 22 | P1 | Add tests for Duration.Value | High | 3m | High | 3m |
| 23 | P1 | Add tests for Email.Scan | High | 4m | High | 4m |
| 24 | P1 | Add tests for Email.Value | High | 3m | High | 3m |
| 25 | P1 | Add tests for URL.Scan | High | 4m | High | 4m |
| 26 | P1 | Add tests for URL.Value | High | 3m | High | 4m |
| 27 | P1 | Add tests for Cents.Scan | High | 4m | High | 4m |
| 28 | P1 | Add tests for Cents.Value | High | 3m | High | 4m |
| 29 | P1 | Add tests for Timestamp.Scan | High | 4m | High | 4m |
| 30 | P1 | Add tests for Timestamp.Value | High | 3m | High | 4m |
| 31 | P1 | Run tests and verify types coverage | High | 2m | High | 2m |
| 32 | P1 | Commit types test improvements | High | 2m | High | 2m |
| **PHASE 5: TEST COVERAGE - ENUMS PACKAGE (6.8% → 50%)** |
| 33 | P1 | Add tests for ActorKind enum methods | High | 5m | High | 5m |
| 34 | P1 | Add tests for Priority enum methods | High | 5m | High | 5m |
| 35 | P1 | Add tests for Status enum methods | High | 5m | High | 5m |
| 36 | P1 | Add tests for Trigger enum methods | High | 5m | High | 5m |
| 37 | P1 | Add tests for MustParse* functions | High | 4m | High | 4m |
| 38 | P1 | Add tests for SQL Scan/Value | High | 5m | High | 5m |
| 39 | P1 | Run tests and verify enums coverage | High | 2m | High | 2m |
| 40 | P1 | Commit enums test improvements | High | 2m | High | 2m |
| **PHASE 6: TEST COVERAGE - LOCALE PACKAGE (28.9% → 60%)** |
| 41 | P1 | Add tests for Locale.Scan | High | 4m | High | 4m |
| 42 | P1 | Add tests for Locale.Value | High | 3m | High | 3m |
| 43 | P1 | Add tests for Locale JSON marshal/unmarshal | High | 4m | High | 4m |
| 44 | P1 | Add tests for Locale.Text marshal/unmarshal | High | 4m | High | 4m |
| 45 | P1 | Add tests for edge cases (empty, invalid) | High | 4m | High | 4m |
| 46 | P1 | Run tests and verify locale coverage | High | 2m | High | 2m |
| 47 | P1 | Commit locale test improvements | High | 2m | High | 2m |
| **PHASE 7: TEST COVERAGE - BOUNDED PACKAGE (43.8% → 60%)** |
| 48 | P1 | Add tests for BoundedString.MarshalJSON | High | 4m | High | 4m |
| 49 | P1 | Add tests for BoundedString.UnmarshalJSON | High | 4m | High | 4m |
| 50 | P1 | Add tests for BoundedString.Scan | High | 4m | High | 4m |
| 51 | P1 | Add tests for BoundedString.Value | High | 3m | High | 3m |
| 52 | P1 | Add tests for MustBoundedString | High | 3m | High | 3m |
| 53 | P1 | Run tests and verify bounded coverage | High | 2m | High | 2m |
| 54 | P1 | Commit bounded test improvements | High | 2m | High | 2m |
| **PHASE 8: TEST COVERAGE - ID PACKAGE (41.9% → 55%)** |
| 55 | P1 | Add tests for ID.Compare with all types | High | 8m | High | 8m |
| 56 | P1 | Add tests for ID.Or | High | 3m | High | 3m |
| 57 | P1 | Add tests for ID binary marshal/unmarshal | High | 6m | High | 6m |
| 58 | P1 | Add tests for ID gob encode/decode | High | 4m | High | 4m |
| 59 | P1 | Run tests and verify id coverage | High | 2m | High | 2m |
| 60 | P1 | Commit id test improvements | High | 2m | High | 2m |
| **PHASE 9: TEST COVERAGE - DATAPOINT PACKAGE (50% → 65%)** |
| 61 | P1 | Add tests for Cause constructors | High | 4m | High | 4m |
| 62 | P1 | Add tests for Reference constructors | High | 4m | High | 4m |
| 63 | P1 | Add tests for Context.With* methods | High | 5m | High | 5m |
| 64 | P1 | Add tests for DataPoint.IsZero | High | 3m | High | 3m |
| 65 | P1 | Add tests for DataPoint JSON round-trip | High | 5m | High | 5m |
| 66 | P1 | Run tests and verify datapoint coverage | High | 2m | High | 2m |
| 67 | P1 | Commit datapoint test improvements | High | 2m | High | 2m |
| **PHASE 10: STATUS REPORT & FINALIZATION** |
| 68 | P1 | Write comprehensive status report | High | 5m | High | 5m |
| 69 | P1 | Commit status report | High | 2m | Medium | 2m |
| 70 | P1 | Push all commits to remote | High | 2m | Medium | 2m |

---

## SUMMARY BY PHASE

| Phase | Tasks | Total Est. | Priority |
|-------|-------|------------|----------|
| 1. Stabilization | 6 | 11m | P0 |
| 2. Linter Config | 3 | 5m | P1 |
| 3. Package Docs | 7 | 15m | P1 |
| 4. Types Tests | 16 | 52m | P1 |
| 5. Enums Tests | 8 | 31m | P1 |
| 6. Locale Tests | 7 | 21m | P1 |
| 7. Bounded Tests | 7 | 20m | P1 |
| 8. ID Tests | 6 | 24m | P1 |
| 9. DataPoint Tests | 7 | 24m | P1 |
| 10. Finalization | 3 | 9m | P1 |
| **TOTAL** | **70** | **~3h** | - |

---

## EXECUTION ORDER (By Impact/Effort Ratio)

Top 20 highest value tasks (sorted by Impact/Effort):

| Rank | # | Task | Impact | Effort | Ratio |
|------|---|------|--------|--------|-------|
| 1 | 1-3 | Fix json/v2 imports | Critical | 3m | 10.0 |
| 2 | 7 | Fix golangci.yml depguard | High | 1m | 9.0 |
| 3 | 10-14 | Add package comments | High | 10m | 8.0 |
| 4 | 4-5 | Verify build/test | Critical | 5m | 8.0 |
| 5 | 17-20 | Add Compare tests | High | 12m | 7.0 |
| 6 | 33-36 | Add enum tests | High | 20m | 5.0 |
| 7 | 41-45 | Add locale tests | High | 19m | 4.5 |
| 8 | 48-52 | Add bounded tests | High | 18m | 4.0 |
| 9 | 55-58 | Add id tests | High | 21m | 3.5 |
| 10 | 61-65 | Add datapoint tests | High | 21m | 3.5 |

---

## NOTES

1. All tasks are designed to be completable in ≤12 minutes
2. Each commit should be self-contained and atomic
3. Run `go test -race ./...` after each phase
4. Total estimated time: ~3 hours for full completion
5. Minimum viable completion: Phases 1-3 (~30 minutes)

---

_Plan created: 2026-03-15 16:00 CET_
