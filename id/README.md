# go-composable-business-types/id

A Go package providing branded, strongly-typed identifiers that prevent mixing different entity IDs at compile time.

## Why?

In Go, regular types like `string` or `int64` provide no compile-time safety:

```go
func GetUser(id string) error { ... }
func GetOrder(id string) error { ... }

userID := "user-123"
orderID := "order-456"

GetOrder(userID)  // Compiles! Runtime bug.
```

With this package, the compiler catches these errors:

```go
type UserBrand struct{}
type OrderBrand struct{}

type UserID = ID[UserBrand, string]
type OrderID = ID[OrderBrand, string]

func GetUser(id UserID) error { ... }
func GetOrder(id OrderID) error { ... }

userID := NewID[UserBrand]("user-123")
orderID := NewID[OrderBrand]("order-456")

GetOrder(userID)  // Compile error: type mismatch
```

## Installation

```bash
go get github.com/larsartmann/go-composable-business-types/id
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/larsartmann/go-composable-business-types/id"
)

type UserBrand struct{}
type OrderBrand struct{}

type UserID = id.ID[UserBrand, string]
type OrderID = id.ID[OrderBrand, string]

func main() {
    userID := id.NewID[UserBrand]("user-123")
    orderID := id.NewID[OrderBrand]("order-456")
    
    fmt.Println(userID.Get())  // user-123
    
    // Type-safe comparison
    otherUserID := id.NewID[UserBrand]("user-123")
    fmt.Println(userID.Equal(otherUserID))  // true
    
    // Zero value check
    var emptyUserID UserID
    fmt.Println(emptyUserID.IsZero())  // true
    
    // Default value with Or
    defaultID := id.NewID[UserBrand]("unknown")
    fmt.Println(emptyUserID.Or(defaultID).Get())  // unknown
}
```

## Supported Value Types

The package supports multiple underlying value types:

| Type | Alias |
|------|-------|
| `string` | `ID[Brand, string]` |
| `int` | `ID[Brand, int]` |
| `int8` | `ID[Brand, int8]` |
| `int16` | `ID[Brand, int16]` |
| `int32` | `ID[Brand, int32]` |
| `int64` | `ID[Brand, int64]` |
| `uint` | `ID[Brand, uint]` |
| `uint8` | `ID[Brand, uint8]` |
| `uint16` | `ID[Brand, uint16]` |
| `uint32` | `ID[Brand, uint32]` |
| `uint64` | `ID[Brand, uint64]` |

The `NewID` function can infer the value type for common cases:

```go
// String IDs - type inferred
userID := NewID[UserBrand]("user-123")

// Numeric IDs - specify type explicitly
orderID := NewID[OrderBrand, int64](12345)
```

## Serialization

The package implements all standard Go interfaces for seamless serialization.

### JSON

```go
id := NewID[UserBrand]("user-123")
data, _ := json.Marshal(id)
fmt.Println(string(data))  // "user-123"

// Zero values serialize to null
var empty UserID
data, _ := json.Marshal(empty)
fmt.Println(string(data))  // null

// Unmarshal works too
var restored UserID
json.Unmarshal([]byte(`"user-123"`), &restored)
fmt.Println(restored.Get())  // user-123
```

### SQL/Database

```go
// Scan from database
var userID UserID
row.Scan(&userID)

// Save to database
_, err := db.Exec("INSERT INTO users (id, name) VALUES (?, ?)", userID, "John")
```

Supported SQL types: `string`, `int64`, `int32`, `int`, `uint64`, `[]byte`.

### Binary

```go
data, _ := id.MarshalBinary()
// ... transmit or store ...
id.UnmarshalBinary(data)
```

### Text (XML, TOML)

```go
data, _ := id.MarshalText()
// ... use in XML/TOML config ...
id.UnmarshalText(data)
```

### Gob (Go binary)

```go
var buf bytes.Buffer
enc := gob.NewEncoder(&buf)
enc.Encode(id)
// ... decode ...
dec := gob.NewDecoder(&buf)
dec.Decode(&restored)
```

## Formatting

The ID type supports custom formatting:

```go
id := NewID[UserBrand, int64](42)

fmt.Sprintf("%s", id)  // "42"
fmt.Sprintf("%d", id)  // "42"
fmt.Sprintf("%q", id)  // "\"42\""
fmt.Sprintf("%v", id)  // "42"
fmt.Sprintf("%#v", id) // "id(42)"
```

## Comparison

For ordered types (`int`, `int64`, `string`, etc.), you can compare IDs:

```go
id1 := NewID[UserBrand, int64](100)
id2 := NewID[UserBrand, int64](200)

id1.Compare(id2)  // -1 (less)
id2.Compare(id1)  //  1 (greater)
id1.Compare(id1)  //  0 (equal)
```

This enables sorting slices of IDs:

```go
sort.Slice(ids, func(i, j int) bool {
    return ids[i].Compare(ids[j]) < 0
})
```

## API Reference

### Constructors

| Function | Description |
|----------|-------------|
| `NewID[Brand](value)` | Create a new ID (type inferred for strings) |
| `NewID[Brand, ValueType](value)` | Create a new ID with explicit type |

### Methods

| Method | Description |
|--------|-------------|
| `Get() V` | Returns the underlying value |
| `IsZero() bool` | True if ID has its zero value |
| `Reset()` | Sets ID to its zero value |
| `Equal(other ID) bool` | True if IDs are equal |
| `Compare(other ID) int` | -1, 0, or 1 for less/equal/greater |
| `Or(default ID) ID` | Returns self if not zero, otherwise default |
| `String() string` | String representation |
| `MarshalJSON() ([]byte, error)` | JSON serialization |
| `UnmarshalJSON([]byte) error` | JSON deserialization |
| `MarshalText() ([]byte, error)` | Text serialization (XML/TOML) |
| `UnmarshalText([]byte) error` | Text deserialization |
| `MarshalBinary() ([]byte, error)` | Binary serialization |
| `UnmarshalBinary([]byte) error` | Binary deserialization |
| `Scan(any) error` | SQL scan |
| `Value() (driver.Value, error)` | SQL value |

## Performance

The package includes benchmarks. Run them with:

```bash
go test -bench=.
```

Typical results on modern hardware:
- `NewID`: ~1-2 ns/op
- `Get`: ~1 ns/op
- `MarshalJSON`: ~50-100 ns/op
- `Scan` (string): ~30-50 ns/op