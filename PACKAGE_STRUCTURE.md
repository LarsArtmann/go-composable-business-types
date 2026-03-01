# Package Structure for Selective Imports

This document describes the new package structure that enables selective imports while keeping everything in a single Go module.

## Import Pattern

Users can now import only what they need:

```go
// Import just the ID package
import "github.com/larsartmann/go-composable-business-types/id"

// Import just the NanoId package
import "github.com/larsartmann/go-composable-business-types/nanoid"

// Import specific types
import "github.com/larsartmann/go-composable-business-types/types"

// Import the full package (backward compatible)
import cbt "github.com/larsartmann/go-composable-business-types"
```

## Package Organization

| Package     | Path         | Description                                        | Dependencies                          |
| ----------- | ------------ | -------------------------------------------------- | ------------------------------------- |
| `id`        | `/id`        | Branded identifiers (ID[B,V])                      | None                                  |
| `nanoid`    | `/nanoid`    | URL-safe unique IDs                                | None                                  |
| `types`     | `/types`     | Email, URL, Percentage, Cents, Timestamp, Duration | None                                  |
| `enums`     | `/enums`     | ActorKind, Priority, Status, Trigger               | None                                  |
| `bounded`   | `/bounded`   | Length-validated strings                           | None                                  |
| `locale`    | `/locale`    | BCP 47 language tags                               | None                                  |
| `temporal`  | `/temporal`  | Bitemporal tracking                                | types                                 |
| `actor`     | `/actor`     | Actor chains for audit trails                      | id, enums                             |
| `money`     | `/money`     | ISO 4217 currency handling                         | locale                                |
| `datapoint` | `/datapoint` | Complete audit trail data unit                     | nanoid, actor, temporal, enums, types |

## Root Package

The root package (`github.com/larsartmann/go-composable-business-types`) re-exports all subpackages for backward compatibility and convenience:

```go
// Users can still use the monolithic import
import cbt "github.com/larsartmann/go-composable-business-types"

// And access all types through it
type MyID = cbt.ID[MyBrand, string]
```

## Benefits

1. **Selective Imports**: Only compile what you use
2. **No Breaking Changes**: Existing code continues to work
3. **Single Module**: No complex multi-module setup with replace directives
4. **Go 1.26 Compatible**: Uses standard Go module features
5. **Clear Dependencies**: Each package declares its own dependencies

## Migration Guide

### Before (still works)

```go
import cbt "github.com/larsartmann/go-composable-business-types"

func main() {
    id := cbt.NewNanoId()
    email, _ := cbt.NewEmail("test@example.com")
}
```

### After (new option)

```go
import (
    "github.com/larsartmann/go-composable-business-types/nanoid"
    "github.com/larsartmann/go-composable-business-types/types"
)

func main() {
    id := nanoid.NewNanoId()
    email, _ := types.NewEmail("test@example.com")
}
```
