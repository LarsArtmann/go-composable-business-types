# Project Status Report

**Date:** 2026-02-25 12:00
**Project:** go-composable-business-types
**Report Type:** Deduplication Refactoring Completion & Quality Assessment

---

## Executive Summary

Successfully completed the deduplication refactoring task based on semantic analysis. Production code duplications eliminated, test helpers created, and major test files refactored. Reduced total lines by 566 while maintaining 82.5% test coverage.

---

## Work Completed

### A) FULLY DONE

| Task                                  | Files Changed                | Lines Removed | Impact                        |
| ------------------------------------- | ---------------------------- | ------------- | ----------------------------- |
| Extract `ID.reset()` helper           | `id.go`, `id_jsonv2.go`      | 12            | Eliminates 4 duplicate blocks |
| Refactor `NewBitemporalWithRange`     | `datapoint_temporal.go`      | 8             | Uses `withCorrection` helper  |
| Create test assertion helpers         | `test_helpers_test.go` (new) | -             | 10 helper functions           |
| Refactor `datapoint_test.go`          | -                            | ~250          | 60+ assertions converted      |
| Refactor `sql_test.go`                | -                            | ~280          | All Scan tests converted      |
| Refactor `cbt_test.go` tag assertions | -                            | 12            | Tag patterns unified          |

### B) PARTIALLY DONE

| Task                           | Status      | Remaining                                         |
| ------------------------------ | ----------- | ------------------------------------------------- |
| `cbt_test.go` full refactoring | 5% done     | 100+ assertions could use helpers                 |
| Static error definitions       | Not started | 20 dynamic errors in production code              |
| Generic Scan test helper       | Not started | `TestEmail_Scan` vs `TestLocale_Scan` duplication |

### C) NOT STARTED

| Task                                              | Priority | Effort                  |
| ------------------------------------------------- | -------- | ----------------------- |
| Static error variables (`err113`)                 | Medium   | 30 min                  |
| Wrap external errors (`wrapcheck`)                | Medium   | 20 min                  |
| Modernize benchmark loops (`b.N` → `b.Loop()`)    | Low      | 10 min                  |
| Fix unchecked type assertions (`forcetypeassert`) | Low      | 15 min (but safe as-is) |

### D) TOTALLY FUCKED UP

| Issue | Severity | Notes                          |
| ----- | -------- | ------------------------------ |
| None  | -        | All tests pass, no regressions |

---

## Linter Analysis (389 issues)

### By Category

| Category           | Count | Severity   | Action                                   |
| ------------------ | ----- | ---------- | ---------------------------------------- |
| `paralleltest`     | 50    | Low        | Tests don't use t.Parallel()             |
| `revive`           | 50    | Low        | Style preferences                        |
| `wsl_v5`           | 50    | Noise      | Whitespace formatting                    |
| `varnamelen`       | 43    | Noise      | Short variable names                     |
| `exhaustruct`      | 35    | Noise      | Intentional partial initialization       |
| `err113`           | 21    | **Medium** | Dynamic errors → static errors           |
| `wrapcheck`        | 17    | **Medium** | Unwrapped external errors                |
| `recvcheck`        | 13    | Low        | Receiver type consistency                |
| `forcetypeassert`  | 10    | Low        | Safe due to type switch guards           |
| `noinlineerr`      | 10    | Noise      | Inline error handling style              |
| `intrange`         | 10    | Low        | `for i := 0; i < n; i++` → `for range n` |
| `gochecknoglobals` | 8     | Noise      | Intentional enum values                  |
| `cyclop`           | 8     | Low        | Test function complexity                 |
| `gocritic`         | 7     | Low        | Minor code suggestions                   |
| `godot`            | 7     | Noise      | Comment punctuation                      |
| `nlreturn`         | 11    | Noise      | Newline before return                    |
| `dupl`             | 2     | **Low**    | TestEmail_Scan ≈ TestLocale_Scan         |
| `funcorder`        | 2     | Low        | Method ordering                          |
| `depguard`         | 3     | Noise      | Intentional dependencies                 |
| `gocognit`         | 2     | Low        | Cognitive complexity                     |
| `goconst`          | 2     | Low        | Repeated strings                         |
| `golines`          | 3     | Low        | Long lines                               |
| `ireturn`          | 2     | Noise      | Interface return style                   |
| `mnd`              | 6     | Noise      | Magic numbers                            |
| `nilnil`           | 8     | Low        | nil, nil returns                         |
| `staticcheck`      | 3     | Low        | Minor issues                             |
| `testpackage`      | 4     | Noise      | Test package naming                      |
| `unused`           | 2     | Low        | Unused code                              |

### Actionable Linter Issues (45 total)

**High Impact / Low Effort:**

1. `dupl: 2` - Create generic `testScanFunc[T]` helper for Email/Locale Scan tests
2. `intrange: 10` - Modernize `for range n` syntax (Go 1.22+)
3. `noinlineerr: 10` - Extract inline errors to variables

**Medium Impact / Medium Effort:** 4. `err113: 20` - Convert dynamic errors to static sentinel errors 5. `wrapcheck: 17` - Wrap external package errors with context

---

## Top 25 Things to Do Next

### High Impact / Low Effort (Do Now)

| #   | Task                                     | Effort | Impact | Files             |
| --- | ---------------------------------------- | ------ | ------ | ----------------- |
| 1   | Create static error variables            | 30 min | High   | `errors.go` (new) |
| 2   | Modernize `for range n` loops            | 10 min | Medium | `*_test.go`       |
| 3   | Fix `dupl` with generic Scan test helper | 15 min | Medium | `sql_test.go`     |
| 4   | Add `t.Parallel()` to test functions     | 15 min | Medium | `*_test.go`       |
| 5   | Commit and push changes                  | 5 min  | High   | -                 |

### Medium Impact / Medium Effort (Do Soon)

| #   | Task                                        | Effort | Impact | Files         |
| --- | ------------------------------------------- | ------ | ------ | ------------- |
| 6   | Wrap external errors (`wrapcheck`)          | 20 min | Medium | `*.go`        |
| 7   | Refactor remaining `cbt_test.go` assertions | 45 min | Low    | `cbt_test.go` |
| 8   | Add `b.Loop()` for benchmarks (Go 1.24)     | 10 min | Low    | `cbt_test.go` |
| 9   | Fix `recvcheck` receiver type issues        | 15 min | Low    | `*.go`        |
| 10  | Create errors package with typed errors     | 30 min | Medium | `errors/`     |

### Lower Priority (Nice to Have)

| #   | Task                               | Effort | Impact | Files                                 |
| --- | ---------------------------------- | ------ | ------ | ------------------------------------- |
| 11  | Add `nilnil` return validation     | 20 min | Low    | `*.go`                                |
| 12  | Fix `gocritic` suggestions         | 15 min | Low    | `*.go`                                |
| 13  | Fix `staticcheck` issues           | 10 min | Low    | `*.go`                                |
| 14  | Address `goconst` repeated strings | 15 min | Low    | `*.go`                                |
| 15  | Fix `funcorder` method ordering    | 5 min  | Low    | `bounded.go`, `datapoint_temporal.go` |
| 16  | Add godot comment punctuation      | 10 min | Low    | `*.go`                                |
| 17  | Review `ireturn` interface returns | 10 min | Low    | `*.go`                                |
| 18  | Fix `nlreturn` newline issues      | 15 min | Low    | `*.go`                                |
| 19  | Address `mnd` magic numbers        | 15 min | Low    | `*.go`                                |
| 20  | Review `exhaustruct` partial init  | 20 min | Low    | `*.go`                                |

### Architecture Improvements (Future)

| #   | Task                                          | Effort | Impact | Notes                              |
| --- | --------------------------------------------- | ------ | ------ | ---------------------------------- |
| 21  | Consider `github.com/go-playground/validator` | 2h     | Medium | For struct validation              |
| 22  | Consider `github.com/google/uuid`             | 1h     | Low    | Alternative to NanoId for UUIDs    |
| 23  | Consider `github.com/shopspring/decimal`      | 1h     | Medium | Alternative to Cents for precision |
| 24  | Extract errors to dedicated package           | 2h     | Medium | Better error handling API          |
| 25  | Add generic validation interface              | 3h     | High   | `Validatable[T]` pattern           |

---

## Type Model Improvements

### Current Architecture Strengths

1. **Phantom Types**: Excellent type safety via `ID[B, V]` pattern
2. **Immutability**: `With*` methods return copies
3. **JSON/SQL Integration**: Custom marshaling for all types
4. **Error Handling**: Constructors return errors, `Must*` panics

### Potential Improvements

| Area        | Current                | Improvement                   | Benefit                      |
| ----------- | ---------------------- | ----------------------------- | ---------------------------- |
| Error Types | Dynamic `errors.New()` | Static sentinel errors        | `errors.Is()` support        |
| Validation  | Per-type methods       | `Validatable[T]` interface    | Unified validation API       |
| Money       | `Cents` + `Money`      | Consider `decimal.Decimal`    | Better precision for finance |
| Time        | `Timestamp` wrapper    | Consider `time.Time` directly | Simpler, less indirection    |
| Collections | Manual slices          | Consider generic collections  | `Set[T]`, `OrderedSet[T]`    |

---

## Library Recommendations

### Currently Using

| Library                      | Purpose           | Verdict |
| ---------------------------- | ----------------- | ------- |
| `github.com/abice/go-enum`   | Enum generation   | Keep    |
| `github.com/bojanz/currency` | ISO 4217 handling | Keep    |
| `github.com/sixafter/nanoid` | ID generation     | Keep    |

### Consider Adding

| Library                                  | Purpose           | Effort | Benefit                |
| ---------------------------------------- | ----------------- | ------ | ---------------------- |
| `github.com/go-playground/validator/v10` | Struct validation | Medium | Declarative validation |
| `github.com/shopspring/decimal`          | Precise decimals  | Medium | Financial precision    |
| `github.com/rotisserie/eris`             | Error wrapping    | Low    | Better error chains    |
| `github.com/google/uuid`                 | UUID generation   | Low    | Standard UUID support  |

---

## Top Question I Cannot Answer

**Question:** Should we disable `forcetypeassert` in golangci-lint config for the `id.go` type-switch-guarded assertions?

**Context:** The 7 type assertions flagged in `id.go` are actually safe because they're inside type-switch branches that have already verified the type. However, the linter can't prove this statically.

**Options:**

1. Add `//nolint:forcetypeassert` comments to each occurrence
2. Disable `forcetypeassert` globally (too broad)
3. Add `forcetypeassert` to exclude rules for `id.go` only
4. Refactor to use checked assertions (adds verbosity, no runtime benefit)

**Recommendation:** Option 1 - Add targeted `//nolint:forcetypeassert` comments with explanation.

---

## Test Coverage

```
Coverage: 82.5% of statements
Status: PASS
Race Detection: Enabled
```

---

## Git Status

```
Modified:
  cbt_test.go           (12 changes, tag assertions)
  datapoint_temporal.go (13 changes, withCorrection helper)
  datapoint_test.go     (419 changes, assertion helpers)
  id.go                 (15 changes, reset() helper)
  id_jsonv2.go          (3 changes, uses reset())
  sql_test.go           (640 changes, assertion helpers)

New Files:
  test_helpers_test.go  (10 helper functions)
  docs/status/2026-02-25_10-17_duplication-analysis-and-quality-report.md

Net Change: -566 lines (834 removed, 268 added)
```

---

## Next Actions

1. **Commit current changes** with detailed message
2. **Push to remote**
3. **Await user decision** on:
   - Static error variable implementation
   - `forcetypeassert` handling strategy
   - Additional refactoring scope

---

_Generated by Crush AI Assistant_
