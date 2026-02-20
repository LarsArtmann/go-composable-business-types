# Improvement Analysis Report
## go-composable-business-types

**Generated:** 2026-02-20 02:39 CET
**Branch:** master (synced with origin)
**Based on:** Codebase analysis + konfetty library patterns

---

## Executive Summary

Comprehensive analysis identified **8 improvement areas** derived from:
1. Current codebase gaps (missing SQL interfaces, weak Money type)
2. konfetty library learnings (validator pattern, builder pattern, zero-value preservation)

**Top 3 Recommendations:**
1. Add SQL interfaces to all types (NanoId, Email, URL, Cents, BoundedString, Timestamp, ID)
2. Strengthen Money type from alias to wrapper with currency safety
3. Add standard Validator interface for pipeline validation

---

## Current State

| Metric | Value |
|--------|-------|
| Test coverage | 80.1% |
| Tests passing | 117+ with race detector |
| Core types | 16 |
| Enums | 5 (with SQL support) |
| Lines of code | 5,452 |
| Status | Production ready |

### Current Strengths

- ✅ Phantom-typed identifiers (ID[B, V])
- ✅ Self-contained DataPoint with full audit trail
- ✅ FIPS-140 compliant NanoId generation
- ✅ Fluent `With*` methods for immutable updates
- ✅ Enums with SQL Scan/Value support
- ✅ Comprehensive JSON serialization

---

## Identified Gaps

### 1. Missing SQL Interfaces (HIGH PRIORITY)

**Problem:** Most types cannot be persisted to databases.

| Type | Has SQL Scanner | Has SQL Valuer | Impact |
|------|-----------------|----------------|--------|
| `NanoId` | ❌ | ❌ | Cannot store IDs in DB |
| `Email` | ❌ | ❌ | Cannot store emails in DB |
| `URL` | ❌ | ❌ | Cannot store URLs in DB |
| `Cents` | ❌ | ❌ | Cannot store money in DB |
| `BoundedString` | ❌ | ❌ | Cannot store validated strings in DB |
| `Timestamp` | ❌ | ❌ | Cannot store timestamps in DB |
| `ID[B,V]` | ❌ | ❌ | Cannot store branded IDs in DB |
| Enums | ✅ | ✅ | Already implemented |

**Solution:**
```go
// Example for NanoId
func (id *NanoId) Scan(src any) error {
    switch v := src.(type) {
    case string:
        parsed, err := ParseNanoId(v)
        if err != nil {
            return err
        }
        *id = parsed
        return nil
    case []byte:
        return id.Scan(string(v))
    default:
        return fmt.Errorf("cannot scan %T into NanoId", src)
    }
}

func (id NanoId) Value() (driver.Value, error) {
    if id.IsEmpty() {
        return nil, nil
    }
    return id.String(), nil
}
```

**Effort:** 2 hours | **Impact:** HIGH - Enables DB persistence

---

### 2. Weak Money Type (HIGH PRIORITY)

**Problem:** `Money` is just a type alias, not a distinct type.

```go
type Money = currency.Amount  // Alias - no type safety
```

**Issues:**
- Can mix with raw `currency.Amount` from any source
- No compile-time currency safety (USD vs EUR)
- Arithmetic returns raw `currency.Amount`, not `Money`

**Solution:**
```go
type Money struct {
    amount currency.Amount
}

func NewMoney(amount, currency string) (Money, error) {
    amt, err := currency.NewAmount(amount, currency)
    return Money{amount: amt}, err
}

func (m Money) Add(other Money) (Money, error) {
    result, err := m.amount.Add(other.amount)
    return Money{amount: result}, err
}

// Currency-safe operations
func (m Money) SameCurrency(other Money) bool {
    return m.amount.Currency() == other.amount.Currency()
}
```

**Effort:** 2 hours | **Impact:** HIGH - Type safety for financial operations

---

### 3. Missing Validator Interface (MEDIUM - from konfetty)

**Problem:** No standard validation contract across types.

**Solution:**
```go
// Validator is implemented by types that can validate themselves.
type Validator interface {
    Validate() error
}

// Pipeline validation
func ValidateAll(validators ...Validator) error {
    for _, v := range validators {
        if err := v.Validate(); err != nil {
            return err
        }
    }
    return nil
}

// Example implementation
func (e Email) Validate() error {
    if e == "" {
        return ErrInvalidEmail
    }
    return nil
}

func (bs BoundedString) Validate() error {
    length := utf8.RuneCountInString(bs.value)
    if length < bs.minLen || length > bs.maxLen {
        return fmt.Errorf("bounded string validation failed")
    }
    return nil
}
```

**Effort:** 2 hours | **Impact:** MEDIUM - Standardized validation

---

### 4. DataPoint Builder Pattern (MEDIUM - from konfetty)

**Problem:** `NewDataPoint` has many parameters, constructor is verbose.

**Current:**
```go
dp := NewDataPoint(payload, actor, now, now, "reason").
    WithTrigger(TriggerWebhook).
    WithContext(ctx).
    AddTag("source", "api")
```

**Solution:**
```go
dp, err := NewDataPointBuilder[OrderState]().
    WithPayload(order).
    WithActor(actor).
    WithTrigger(TriggerWebhook).
    WithReason("Order created via webhook").
    WithTag("source", "api").
    AddReference(orderRef).
    Build()

// Or with validation:
dp, err := NewDataPointBuilder[OrderState]().
    WithPayload(order).
    WithActor(actor).
    BuildValidated()  // Returns error if required fields missing
```

**Effort:** 3 hours | **Impact:** MEDIUM - Better ergonomics

---

### 5. Comparison Interface (MEDIUM)

**Problem:** No standard way to compare/order types.

**Solution:**
```go
type Comparer[T any] interface {
    Compare(other T) int  // -1, 0, 1
}

func (t Timestamp) Compare(other Timestamp) int {
    return t.Time.Compare(other.Time)
}

func (c Cents) Compare(other Cents) int {
    switch {
    case c < other: return -1
    case c > other: return 1
    default: return 0
    }
}

// Enables sorting
func SortCents(cents []Cents) {
    sort.Slice(cents, func(i, j int) bool {
        return cents[i].Compare(cents[j]) < 0
    })
}
```

**Effort:** 1 hour | **Impact:** MEDIUM - Sorting and comparison

---

### 6. Zero-Value Preservation Merge (LOW - from konfetty)

**Problem:** No way to merge defaults while preserving user values.

**Solution:**
```go
// MergeDefaults merges non-zero values from defaults into target.
// Zero values in target are filled from defaults.
func MergeDefaults[T any](target, defaults T) T {
    // Reflection-based merge
    // Only overwrites zero values in target
}
```

**Use Case:**
```go
defaultConfig := Config{Timeout: 30, Retries: 3}
userConfig := Config{Timeout: 0, Retries: 5}  // Timeout not specified
result := MergeDefaults(userConfig, defaultConfig)
// result.Timeout = 30 (from defaults)
// result.Retries = 5 (user value preserved)
```

**Effort:** 3 hours | **Impact:** LOW - Configuration handling

---

### 7. Provider Interface (LOW - from konfetty)

**Problem:** No lazy-loading pattern for expensive resources.

**Solution:**
```go
type Provider[T any] interface {
    Load() (T, error)
}

type CachedProvider[T any] struct {
    loader func() (T, error)
    cached *T
    mu     sync.RWMutex
}

func (p *CachedProvider[T]) Load() (T, error) {
    p.mu.RLock()
    if p.cached != nil {
        defer p.mu.RUnlock()
        return *p.cached, nil
    }
    p.mu.RUnlock()

    p.mu.Lock()
    defer p.mu.Unlock()

    if p.cached != nil {
        return *p.cached, nil
    }

    v, err := p.loader()
    if err != nil {
        return v, err
    }
    p.cached = &v
    return v, nil
}
```

**Effort:** 2 hours | **Impact:** LOW - Lazy loading

---

### 8. Circular Reference Protection (LOW)

**Problem:** Recursive DataPoint structures could cause infinite loops.

**Solution:**
```go
func (dp DataPoint[T]) equalsVisited(other DataPoint[T], visited map[uintptr]bool) bool {
    addr := reflect.ValueOf(dp).Pointer()
    if visited[addr] {
        return true  // Already comparing this reference
    }
    visited[addr] = true
    // ... recursive comparison with visited tracking
}
```

**Effort:** 2 hours | **Impact:** LOW - Safety edge case

---

## konfetty Patterns Applicable

| Pattern | konfetty Implementation | Applicability Here |
|---------|------------------------|-------------------|
| Fluent Builder API | `WithDefaults().WithTransformer().Build()` | Already using `With*` pattern ✅ |
| Generics + Reflection | Type-safe API, reflection internally | Already using phantom types ✅ |
| Validator Interface | `Validator` interface for pipeline | Add to all types |
| Zero-Value Preservation | Only overwrite zero values | Add merge functionality |
| Provider Interface | `Provider[T]` for lazy loading | Add for expensive resources |
| Circular Reference Protection | `visited` map tracking | Add for recursive structures |

---

## Prioritized Implementation Plan

### Phase 1: Critical (Do First)

| Task | Effort | Files Changed | Impact |
|------|--------|---------------|--------|
| SQL interfaces for all types | 2h | nanoid.go, common.go, bounded.go, id.go | HIGH |
| Stronger Money wrapper | 2h | money.go | HIGH |

### Phase 2: Important

| Task | Effort | Files Changed | Impact |
|------|--------|---------------|--------|
| Validator interface | 2h | New validator.go, all types | MEDIUM |
| DataPoint builder | 3h | datapoint.go | MEDIUM |
| Comparison interface | 1h | common.go, timestamp | MEDIUM |

### Phase 3: Nice to Have

| Task | Effort | Files Changed | Impact |
|------|--------|---------------|--------|
| Merge/defaults | 3h | New merge.go | LOW |
| Provider interface | 2h | New provider.go | LOW |
| Circular reference protection | 2h | datapoint.go | LOW |

---

## Estimated Total Effort

| Phase | Effort | Value Delivered |
|-------|--------|-----------------|
| Phase 1 | 4 hours | 60% of improvement value |
| Phase 2 | 6 hours | 30% of improvement value |
| Phase 3 | 7 hours | 10% of improvement value |
| **Total** | **17 hours** | **100% of identified improvements** |

---

## Quick Wins (Can Do Immediately)

1. **Add SQL interfaces** - 2h, enables DB persistence for all types
2. **Fix gopls warnings** - 30min, remove 18 unnecessary type args in tests
3. **Add Compare methods** - 1h, enables sorting for Timestamp, Cents

---

## Risk Assessment

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| SQL interface edge cases | Medium | Comprehensive tests for null/empty values |
| Money wrapper breaking changes | Low | Keep `currency.Amount` accessible via method |
| Builder pattern complexity | Low | Keep existing constructors, add builder as option |

---

## Conclusion

The go-composable-business-types library is production-ready with solid foundations. The identified improvements focus on:

1. **Completing DB persistence** (SQL interfaces)
2. **Strengthening type safety** (Money wrapper)
3. **Standardizing patterns** (Validator, Comparer interfaces)
4. **Improving ergonomics** (DataPoint builder)

**Recommended Next Step:** Implement SQL interfaces for all types (Phase 1, Task 1). This provides the highest value (DB persistence) for lowest effort (2 hours).

---

*Report generated by Crush CLI Agent*
*Generated: 2026-02-20 02:39 CET*
*Analysis based on 17 source files, 5,452 lines of code*
*Patterns derived from nikoksr/konfetty research*
