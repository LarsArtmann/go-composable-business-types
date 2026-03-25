# Status Report: id/ Package Review

**Date:** 2026-03-25 13:40
**Author:** Crush (AI Assistant)
**Scope:** Deep architectural review of `id/` package

---

## Executive Summary

The `id/` package provides branded, strongly-typed identifiers using phantom types. While the concept is excellent, the implementation has **critical architectural issues** that violate project standards and compromise type safety at compile time.

---

## Work Status

### a) FULLY DONE

| Task | Status | Notes |
|------|--------|-------|
| Refactored `Compare` to use `cmp.Compare` | ✅ Complete | Reduced ~130 lines of manual compare functions to ~30 lines using Go 1.21+ stdlib |
| All tests passing | ✅ Complete | 1252 lines of tests, all green |
| Linter passing | ✅ Complete | 0 issues |

### b) PARTIALLY DONE

| Task | Status | Notes |
|------|--------|-------|
| Type safety enforcement | ⚠️ Partial | Phantom types work, but value type constraint is too loose |

### c) NOT STARTED

| Task | Priority | Impact |
|------|----------|--------|
| Split id.go into multiple files | CRITICAL | High |
| Add compile-time type constraint for supported value types | HIGH | High |
| Fix test file unnecessary type arguments | LOW | Low |
| Add `MustNewID` constructor | MEDIUM | Medium |
| Consider removing `Reset()` method | MEDIUM | Medium |

### d) TOTALLY FUCKED UP

| Issue | Severity | Description |
|-------|----------|-------------|
| **File size violation** | CRITICAL | `id.go` is 854 lines (max 350) |
| **Loose type constraint** | HIGH | `V comparable` allows unsupported types, causing runtime errors |
| **Massive code duplication** | HIGH | Same type-switch pattern repeated 10+ times (~600 lines) |
| **Inconsistent mutability** | MEDIUM | `Reset()` mutates, everything else is immutable |

### e) WHAT WE SHOULD IMPROVE

1. **Type Constraint** - Replace `comparable` with a union type that only allows supported types
2. **File Structure** - Split into logical files: `id.go`, `json.go`, `text.go`, `binary.go`, `sql.go`
3. **Code Duplication** - Consider helper functions or accept the duplication as explicit (Go idiom)
4. **Mutability** - Remove or document the exception for `Reset()`

---

## Critical Issues Analysis

### 1. File Size Violation (CRITICAL)

**Current:** 854 lines
**Limit:** 350 lines
**Violation:** 2.4x over limit

The file contains:
- Core type and methods: ~160 lines
- JSON serialization: ~190 lines
- Text serialization: ~50 lines
- Binary serialization: ~210 lines
- SQL serialization: ~165 lines
- Interface assertions: ~30 lines

### 2. Loose Type Constraint (HIGH)

**Current:**
```go
type ID[B any, V comparable] struct{ value V }
```

**Problem:** `comparable` allows ANY comparable type, but we only support:
- `string`, `int`, `int8`, `int16`, `int32`, `int64`
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- `float32`, `float64`

**Result:** Runtime errors for unsupported types instead of compile-time safety.

**Proposed Fix:**
```go
// IDValue is the set of types supported for ID values.
type IDValue interface {
    string | int | int8 | int16 | int32 | int64 |
    uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

type ID[B any, V IDValue] struct{ value V }
```

### 3. Code Duplication Pattern (MEDIUM)

The type-switch pattern is repeated in 10 methods:
- `Compare`, `String`, `Format`
- `MarshalJSON`, `UnmarshalJSON`
- `MarshalText`, `UnmarshalText`
- `MarshalBinary`, `UnmarshalBinary`
- `Scan`, `Value`

This is ~600 lines of similar code. However, in Go, explicit type switches are often preferred over reflection or interface tricks for performance and clarity.

**Decision:** Accept duplication as explicit Go idiom, but split across files.

### 4. Inconsistent Mutability (LOW)

`Reset()` mutates the ID in place, while all other operations are immutable. This breaks functional patterns.

**Options:**
1. Remove `Reset()` - use `var zero ID[B, V]` instead
2. Document the exception
3. Add `WithReset() ID[B, V]` that returns a zero value

---

## Top 25 Things to Get Done Next

### Priority 1: CRITICAL (Do Immediately)

1. **Split id.go into multiple files** - 854 lines → ~150 lines each
   - `id.go` - Core type, NewID, Get, IsZero, Equal, Compare, Or, String, Format
   - `json.go` - MarshalJSON, UnmarshalJSON
   - `text.go` - MarshalText, UnmarshalText
   - `binary.go` - MarshalBinary, UnmarshalBinary, GobEncode, GobDecode
   - `sql.go` - Scan, Value, interface assertions

2. **Add compile-time type constraint** - Replace `comparable` with `IDValue` interface

### Priority 2: HIGH (Do Soon)

3. **Add `MustNewID` constructor** - Panic on invalid input (for tests/cfg)
4. **Consider `Reset()` removal or rename** - Document mutability exception
5. **Review integration with `nanoid` package** - Ensure seamless usage

### Priority 3: MEDIUM (Do Eventually)

6. **Fix test file unnecessary type arguments** - 12 gopls warnings
7. **Add BDD-style tests** - Given/When/Then for critical paths
8. **Add examples for NanoID integration** - Show recommended pattern
9. **Review error messages consistency** - Ensure structured errors
10. **Add `IsZero()` benchmark** - Performance validation

### Priority 4: LOW (Nice to Have)

11. **Consider `ID.parse(string)` method** - For text parsing
12. **Add `ID.fromBytes([]byte)` method** - For binary parsing
13. **Consider `OrderedID` constraint** - IDs that support Compare
14. **Add `ID.MarshalTextJSON()` - Text-based JSON for large IDs**
15. **Review float support necessity** - Are float IDs actually used?

### Priority 5: FUTURE (Consider for Later)

16. **UUID support** - Add if needed
17. **ULID support** - Add if needed
18. **CUID2 support** - Add if needed
19. **Custom validation** - Allow brand types to validate values
20. **Pool support** - For high-performance scenarios

### Priority 6: DOCUMENTATION

21. **Update README with type constraint** - Document supported types
22. **Add architecture diagram** - Show phantom type flow
23. **Add migration guide** - If breaking changes
24. **Add performance benchmarks** - Compare with raw types
25. **Add security considerations** - Discuss ID exposure risks

---

## Proposed File Structure

```
id/
├── id.go           # Core type, NewID, Get, IsZero, Equal, Compare, Or, String, Format (~160 lines)
├── json.go         # MarshalJSON, UnmarshalJSON (~190 lines)
├── text.go         # MarshalText, UnmarshalText (~50 lines)
├── binary.go       # MarshalBinary, UnmarshalBinary, GobEncode, GobDecode (~210 lines)
├── sql.go          # Scan, Value, interface assertions (~165 lines)
├── id_test.go      # Tests (1252 lines - acceptable for tests)
└── README.md       # Documentation
```

---

## Impact Assessment

### Customer Value

1. **Type Safety** - Compile-time errors prevent production bugs
2. **Maintainability** - Smaller files are easier to understand and modify
3. **Performance** - Explicit type switches are fast
4. **Developer Experience** - Clear errors, good documentation

### Risk Assessment

| Change | Risk | Mitigation |
|--------|------|------------|
| File split | LOW | No API changes, just file organization |
| Type constraint | MEDIUM | Breaking change for custom types, but those didn't work anyway |
| Remove Reset() | LOW | Easy to add back if needed |

---

## Non-Obvious Insights

1. **Phantom types have zero runtime cost** - The brand type `B` is erased at compile time
2. **Empty string as zero is intentional** - Allows optional IDs without pointers
3. **Binary format is little-endian** - Matches most modern architectures
4. **SQL returns int64 for all integers** - That's why we don't need int8/int16 in Scan

---

## Execution Plan (Sorted by Impact/Effort)

| Step | Task | Effort | Impact | ROI |
|------|------|--------|--------|-----|
| 1 | Split id.go into files | 30min | HIGH | ⭐⭐⭐⭐⭐ |
| 2 | Add IDValue type constraint | 15min | HIGH | ⭐⭐⭐⭐⭐ |
| 3 | Fix test type arguments | 10min | LOW | ⭐⭐ |
| 4 | Add MustNewID | 5min | MEDIUM | ⭐⭐⭐ |
| 5 | Document Reset() exception | 5min | LOW | ⭐⭐ |
| 6 | Commit and push | 5min | - | - |

**Total Estimated Time:** ~70 minutes

---

## Questions for User

### #1 Question I Cannot Figure Out Myself

**Should we add a `Brand` interface constraint with a `Name()` method?**

Currently:
```go
type ID[B any, V comparable] struct{ value V }
```

Proposed:
```go
// Brand is a phantom type that identifies the entity type of an ID.
type Brand interface {
    Name() string // Returns the entity name for debugging/logging
}

type ID[B Brand, V comparable] struct{ value V }
```

Usage would change from:
```go
type UserBrand struct{}
type UserID = ID[UserBrand, string]
```

To:
```go
type UserBrand struct{}
func (UserBrand) Name() string { return "User" }
type UserID = ID[UserBrand, string]
```

**Benefits:**
- ✅ Runtime introspection of brand name (useful for debugging, logging, error messages)
- ✅ Self-documenting ID types
- ✅ Could enable generic error messages: `"invalid User ID"` instead of `"invalid ID"`

**Trade-offs:**
- ❌ Breaking change - requires adding `Name()` to all existing brand types
- ❌ Tiny overhead (but phantom types are erased, so negligible)
- ❌ Slightly more boilerplate

**Alternative:** Keep `B any` - the phantom type pattern works without any interface

**My recommendation:** Keep `B any` for now. The phantom type pattern is elegant precisely because it requires NO methods. Adding `Name()` adds boilerplate without enough benefit. If we need debugging, we can use reflection or code generation.

---

## Updated Priority List (After User Feedback)

| # | Task | Priority | Status |
|---|------|----------|--------|
| 1 | Split id.go into 5 files | CRITICAL | Pending |
| 2 | Keep `V comparable` | DECIDED | ✅ No change needed |
| 3 | Keep `B any` | DECIDED | ✅ No change needed |
| 4 | Fix test type arguments | LOW | Pending |
| 5 | Add MustNewID constructor | MEDIUM | Pending |

---

*End of Status Report*
