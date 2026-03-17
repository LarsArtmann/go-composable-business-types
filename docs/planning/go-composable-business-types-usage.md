# go-composable-business-types Library Usage Guide

## Executive Summary

This document provides a comprehensive analysis of how the `id` package (and related types) should be used within the `go-composable-business-types` library. The `id` package provides **branded, phantom-type identifiers** that prevent mixing different entity IDs at compile time.

**Key Principle**: The `id` package is designed for **external consumers** of this library to create type-safe identifiers for their domain entities. Internal packages should use the appropriate primitive or specific type that best fits their use case.

---

## 1. The `id` Package: Purpose and Design

### 1.1 Core Concept: Phantom Types for Compile-Time Safety

The `id` package uses Go generics to create distinct identifier types that cannot be accidentally mixed:

```go
// Without id package (unsafe)
func GetUser(id string) error { ... }
func GetOrder(id string) error { ... }
GetOrder("user-123") // Compiles! Runtime bug.

// With id package (safe)
type UserBrand struct{}
type OrderBrand struct{}
type UserID = id.ID[UserBrand, string]
type OrderID = id.ID[OrderBrand, string]

func GetUser(id UserID) error { ... }
func GetOrder(id OrderID) error { ... }
GetOrder(id.NewID[UserBrand]("user-123")) // Compile error: type mismatch
```

### 1.2 API Surface

| Constructor                      | Usage                                 |
| -------------------------------- | ------------------------------------- |
| `NewID[Brand](value)`            | Create ID (type inferred for strings) |
| `NewID[Brand, ValueType](value)` | Create ID with explicit type          |

| Methods              | Purpose                  |
| -------------------- | ------------------------ |
| `Get() V`            | Extract underlying value |
| `IsZero() bool`      | Check for zero value     |
| `Equal(other) bool`  | Type-safe equality       |
| `Compare(other) int` | Ordered comparison       |
| `Or(default) ID`     | Default value handling   |

### 1.3 Serialization Support

The `ID` type implements all standard Go interfaces:

- `json.Marshaler` / `json.Unmarshaler` → `"user-123"` or `null`
- `sql.Scanner` / `driver.Valuer` → Database storage
- `encoding.TextMarshaler` / `TextUnmarshaler` → XML/TOML
- `encoding.BinaryMarshaler` / `BinaryUnmarshaler` → Binary protocols
- `gob.GobEncoder` / `GobDecoder` → Go-specific encoding

---

## 2. Internal Package Usage Patterns

### 2.1 Current State Analysis

| Package                  | Uses `id` Package? | Uses `nanoid` Directly? | Pattern                                    |
| ------------------------ | ------------------ | ----------------------- | ------------------------------------------ |
| `actor`                  | ✅ Yes             | ❌ No                   | `id.ID[struct{}, T]` for generic actor IDs |
| `datapoint`              | ❌ No              | ✅ Yes                  | `nanoid.NanoId` for DataPoint IDs          |
| `datapoint/reference.go` | ❌ No              | ❌ No                   | Generic `T comparable` for reference IDs   |
| `datapoint/cause.go`     | ❌ No              | ✅ Yes                  | `nanoid.NanoId` for cause IDs              |

### 2.2 Analysis: Correct Usage Patterns

#### ✅ `actor` Package: Correct Usage

The `actor` package correctly uses `id.ID[struct{}, T]` because:

- It provides **generic actor entry types** for consumers
- The brand is `struct{}` (unbranded) because actor IDs can be any comparable type
- Consumers can use `id.NewID[struct{}, string]("user-123")` directly

```go
// actor/actor.go
type ActorEntry[T comparable] struct {
    Kind enums.ActorKind
    Id   id.ID[struct{}, T]  // Generic, works with any ID type
    Name string
}

// Usage in consumer code:
userID := id.NewID[struct{}, string]("user-123")
actorEntry := actor.UserActor(userID, "Alice")
```

#### ✅ `datapoint` Package: Correct Usage

The `datapoint` package correctly uses `nanoid.NanoId` directly because:

- DataPoint IDs are **always NanoIds** (internal implementation detail)
- Consumers don't need to brand DataPoint IDs (they're system-generated)
- The `Reference[T]` and `Cause[T]` types use generics for flexibility

```go
// datapoint/datapoint.go
type DataPoint[T comparable] struct {
    id      nanoid.NanoId  // Always a NanoId
    payload T
    // ...
}

// Reference allows any comparable type for external entity IDs
type Reference[T comparable] struct {
    id       T  // Consumer's ID type
    relation string
}
```

---

## 3. Recommended Usage Patterns for Library Consumers

### 3.1 Pattern 1: String-Based Entity IDs (Simple)

For simple cases where string IDs are sufficient:

```go
package myapp

import "github.com/larsartmann/go-composable-business-types/id"

type UserBrand struct{}
type OrderBrand struct{}

type UserID = id.ID[UserBrand, string]
type OrderID = id.ID[OrderBrand, string]

func NewUserID(s string) UserID { return id.NewID[UserBrand](s) }
func NewOrderID(s string) OrderID { return id.NewID[OrderBrand](s) }

func GetUser(id UserID) (*User, error) { ... }
func GetOrder(id OrderID) (*Order, error) { ... }
```

### 3.2 Pattern 2: NanoId-Based Entity IDs (Recommended)

For production systems requiring unique, URL-safe identifiers:

```go
package myapp

import (
    "github.com/larsartmann/go-composable-business-types/id"
    "github.com/larsartmann/go-composable-business-types/nanoid"
)

type UserBrand struct{}
type OrderBrand struct{}

type UserID = id.ID[UserBrand, nanoid.NanoId]
type OrderID = id.ID[OrderBrand, nanoid.NanoId]

// Generate new IDs
func NewUserID() UserID {
    return id.NewID[UserBrand](nanoid.NewNanoId())
}

// Parse existing IDs
func ParseUserID(s string) (UserID, error) {
    nano, err := nanoid.ParseNanoId(s)
    if err != nil {
        return UserID{}, err
    }
    return id.NewID[UserBrand](nano), nil
}
```

### 3.3 Pattern 3: Numeric Entity IDs

For systems using auto-increment or sequence-based IDs:

```go
package myapp

import "github.com/larsartmann/go-composable-business-types/id"

type UserBrand struct{}
type ProductBrand struct{}

type UserID = id.ID[UserBrand, int64]
type ProductID = id.ID[ProductBrand, int64]

func NewUserID(n int64) UserID {
    return id.NewID[UserBrand, int64](n)
}
```

### 3.4 Pattern 4: Integration with DataPoint

When using DataPoint with branded IDs:

```go
package myapp

import (
    "github.com/larsartmann/go-composable-business-types/actor"
    "github.com/larsartmann/go-composable-business-types/datapoint"
    "github.com/larsartmann/go-composable-business-types/id"
)

type UserBrand struct{}
type UserID = id.ID[UserBrand, string]

// Create actor with branded ID
userID := id.NewID[UserBrand]("user-123")
actorEntry := actor.UserActor(
    id.NewID[struct{}, string](userID.Get()), // Extract string for ActorEntry
    "Alice",
)

// Create DataPoint
dp := datapoint.NewDataPoint(orderPayload, actorEntry)

// Reference with branded ID (extract underlying value)
customerRef := datapoint.NewReference(userID.Get(), "customer")
dp = dp.WithReference(customerRef)
```

---

## 4. Anti-Patterns to Avoid

### ❌ Anti-Pattern 1: Using `id.ID` for Internal Implementation Details

```go
// BAD: DataPoint doesn't need branded IDs internally
type DataPoint[T comparable] struct {
    id id.ID[DataPointBrand, nanoid.NanoId]  // Unnecessary overhead
}

// GOOD: Use the underlying type directly
type DataPoint[T comparable] struct {
    id nanoid.NanoId  // Clear, simple, sufficient
}
```

### ❌ Anti-Pattern 2: Requiring Consumers to Use `id` Package

```go
// BAD: Forces consumers to use id package
type Reference[T comparable] struct {
    id id.ID[struct{}, T]  // Unnecessary constraint
}

// GOOD: Accept any comparable type
type Reference[T comparable] struct {
    id T  // Consumer can use string, their own ID type, etc.
}
```

### ❌ Anti-Pattern 3: Multiple Brand Types for Same Entity

```go
// BAD: Confusing, fragments the type system
type UserBrand1 struct{}
type UserBrand2 struct{}
type UserID1 = id.ID[UserBrand1, string]
type UserID2 = id.ID[UserBrand2, string]

// GOOD: Single brand per entity
type UserBrand struct{}
type UserID = id.ID[UserBrand, string]
```

---

## 5. When to Use Each Type

### Decision Matrix

| Use Case                           | Recommended Type                                | Rationale                          |
| ---------------------------------- | ----------------------------------------------- | ---------------------------------- |
| Public entity IDs in your domain   | `id.ID[Brand, V]`                               | Compile-time safety                |
| System-generated internal IDs      | `nanoid.NanoId`                                 | Direct, simple, no branding needed |
| Generic storage (Reference, Cause) | `T comparable`                                  | Maximum flexibility for consumers  |
| Actor IDs in audit trails          | `id.ID[struct{}, T]`                            | Generic but still type-safe        |
| Database primary keys              | `id.ID[Brand, int64]` or `id.ID[Brand, string]` | Depends on DB schema               |
| External API identifiers           | `id.ID[Brand, string]`                          | JSON-friendly, human-readable      |

### Quick Reference

```go
// For your domain entities (consumers)
type MyEntityBrand struct{}
type MyEntityID = id.ID[MyEntityBrand, nanoid.NanoId]

// For internal system IDs (library internals)
type internalStruct struct {
    id nanoid.NanoId  // No branding needed
}

// For generic references (library internals)
type Reference[T comparable] struct {
    id T  // Consumer's choice
}

// For actor entries (library provides)
type ActorEntry[T comparable] struct {
    Id id.ID[struct{}, T]  // Generic but type-safe
}
```

---

## 6. Implementation Guidelines

### 6.1 For Library Maintainers

1. **Internal packages should NOT force `id` usage on consumers**
   - Accept `T comparable` or specific types (string, nanoid.NanoId)
   - Allow consumers to decide if they want branded IDs

2. **Use `id` package for generic types that cross domain boundaries**
   - `ActorEntry` uses `id.ID[struct{}, T]` because it travels across system boundaries
   - The generic `T` allows consumers to use their preferred ID type

3. **Document the expected usage patterns**
   - Examples should show both simple string IDs and NanoId-based IDs
   - Make it clear when branding is beneficial vs. overkill

### 6.2 For Library Consumers

1. **Define brands for your domain entities**

   ```go
   type UserBrand struct{}
   type OrderBrand struct{}
   type ProductBrand struct{}
   ```

2. **Create type aliases for convenience**

   ```go
   type UserID = id.ID[UserBrand, nanoid.NanoId]
   type OrderID = id.ID[OrderBrand, nanoid.NanoId]
   ```

3. **Provide constructors for your ID types**

   ```go
   func NewUserID() UserID { ... }
   func ParseUserID(s string) (UserID, error) { ... }
   ```

4. **Extract underlying values when interfacing with library types**
   ```go
   // id.ID[UserBrand, string] → string for ActorEntry
   actor.UserActor(id.NewID[struct{}, string](userID.Get()), "Name")
   ```

---

## 7. Testing Considerations

### 7.1 Testing ID Types

```go
func TestUserID(t *testing.T) {
    id1 := NewUserID()
    id2 := NewUserID()

    // Equality
    if id1.Equal(id2) {
        t.Error("different IDs should not be equal")
    }

    // Same value
    id3 := id.NewID[UserBrand](id1.Get())
    if !id1.Equal(id3) {
        t.Error("same value should be equal")
    }

    // Zero value
    var zero UserID
    if !zero.IsZero() {
        t.Error("zero value should be zero")
    }
}
```

### 7.2 Testing Integration with Library Types

```go
func TestDataPointWithBrandedID(t *testing.T) {
    userID := NewUserID()

    // Extract for ActorEntry
    actorEntry := actor.UserActor(
        id.NewID[struct{}, string](userID.Get()),
        "Test User",
    )

    dp := datapoint.NewDataPoint(payload, actorEntry).
        WithReference(datapoint.NewReference(userID.Get(), "user"))

    // Verify serialization
    data, err := json.Marshal(dp)
    // ... assertions
}
```

---

## 8. Migration Guide

### From Raw Strings

```go
// Before
func GetUser(id string) (*User, error)

// After
type UserBrand struct{}
type UserID = id.ID[UserBrand, string]
func GetUser(id UserID) (*User, error)
```

### From Plain NanoId

```go
// Before
func NewUser() (nanoid.NanoId, error)

// After
type UserBrand struct{}
type UserID = id.ID[UserBrand, nanoid.NanoId]
func NewUser() (UserID, error)
```

---

## 9. Summary

The `id` package is a **powerful tool for library consumers** to achieve compile-time safety for their domain entity identifiers. It should be used:

- ✅ By **consumers** to brand their entity IDs
- ✅ By **generic library types** (like `ActorEntry`) that cross domain boundaries
- ❌ NOT for internal implementation details that don't need branding
- ❌ NOT to force constraints on consumers

The current architecture correctly separates concerns:

- `actor` uses `id.ID[struct{}, T]` for generic, type-safe actor identification
- `datapoint` uses `nanoid.NanoId` directly for system-generated IDs
- `Reference[T]` and `Cause[T]` use `T comparable` for maximum consumer flexibility

---

## 10. References

- `/id/README.md` - Complete ID package documentation
- `/id/id.go` - Implementation details
- `/examples/basic/main.go` - Basic usage examples
- `/examples/datapoint/main.go` - DataPoint integration examples
- `/AGENTS.md` - Project conventions and build commands
