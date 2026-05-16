# Project: go-composable-business-types

A Go library of strongly typed, composable base values for business applications.

## Project Goal

Build the best possible data types for Lars Artmann's Golang applications, leveraging superb existing libraries where appropriate. **Backwards compatibility is NOT a concern** — we prioritize clean, optimal design over legacy support.

## Build & Test Commands

```bash
# Generate enum code
go generate ./...

# Build
go build ./...

# Run tests
go test -race ./...

# Run tests with coverage
go test -race -coverprofile=coverage.out ./...

# Lint
golangci-lint run --fix
```

## Dependencies

- `github.com/abice/go-enum` - Enum code generation (`//go:generate go-enum`)
- `github.com/bojanz/currency` - ISO 4217 currency handling
- `github.com/larsartmann/go-branded-id` - Branded phantom-type identifiers (extracted from this project, published as separate module)
- `github.com/sixafter/nanoid` - FIPS-140 compatible, high-performance NanoID generation
- `golang.org/x/text` - BCP 47 locale/language support

## Package Structure (Go 1.26 Selective Imports)

This library uses a single Go module with subpackages for selective imports:

```
.
├── actor/              # ActorChain[T], ActorEntry[T] - audit trail tracking
├── bounded/            # BoundedString - length-validated strings
├── datapoint/          # DataPoint[T] - self-contained data with audit trail
│   ├── datapoint.go    # DataPoint[T] main type
│   ├── context.go      # Execution context
│   ├── reference.go    # Reference[T] entity references
│   └── cause.go        # Cause[T] causal relationships
├── enums/              # ActorKind, Priority, Status, Trigger, CauseKind enums
│   └── enum_enum.go    # Generated enum code (do not edit)
├── importance/         # Importance - 0-100 priority classification
├── locale/             # Locale - BCP 47 language tags
├── money/              # Money - ISO 4217 currency wrapper
├── nanoid/             # NanoID - URL-safe unique identifiers
├── projectcore/        # ProjectCore - composite project metadata
├── tag/                # Tag - validated string labels
├── temporal/           # Bitemporal - valid/recorded time tracking
├── types/              # Email, URL, Percentage, Cents, Timestamp, Duration
├── validate/           # Validator interface for self-validating types
├── version/            # Build version info from runtime/debug
├── pkg/errors/         # Centralized sentinel and structured errors
├── scanutil/           # SQL Scanner/Valuer helpers
└── testutil/           # Generic test helpers
```

**Note:** Branded ID types (`ID[B, V]`) live in the separate module [`go-branded-id`](https://github.com/larsartmann/go-branded-id).

### Selective Import Examples

```go
// Import just what you need
import "github.com/larsartmann/go-composable-business-types/nanoid"
import "github.com/larsartmann/go-composable-business-types/types"

func main() {
    id := nanoid.New()
    email, _ := types.NewEmail("test@example.com")
}
```

```go
// Import generic types with type parameters
import (
    "github.com/larsartmann/go-branded-id"
    "github.com/larsartmann/go-composable-business-types/actor"
)

type UserBrand struct{}
type UserID = id.ID[UserBrand, string]

func main() {
    userID := id.NewID[UserBrand, string]("user-123")
    actorEntry := actor.UserActor(userID, "John Doe")
}
```

```go
// Import enums
import "github.com/larsartmann/go-composable-business-types/enums"

func main() {
    kind := enums.ActorKindUser
    trigger := enums.TriggerManual
}
```

```go
// Import DataPoint for complete audit trails
import (
    "github.com/larsartmann/go-composable-business-types/datapoint"
    "github.com/larsartmann/go-composable-business-types/actor"
    "github.com/larsartmann/go-composable-business-types/enums"
    "github.com/larsartmann/go-branded-id"
)

func main() {
    userID := id.NewID[struct{}, string]("user-123")
    actorEntry := actor.UserActor(userID, "John Doe")

    // Create a DataPoint with complete audit trail
    dp := datapoint.NewDataPoint("order-123", actorEntry).
        WithTrigger(enums.TriggerWebhook).
        WithReason("Customer checkout").
        WithTag("priority", "high")

    // Add references and causes
    ref := datapoint.NewReference("customer-456", "customer")
    dp = dp.WithReference(ref)
}
```

## Code Conventions

- **Phantom types**: Strong type safety with branded types (NanoID, Email, etc.)
- **Functional patterns**: Immutable value types with `With*` methods returning copies
- **JSON serialization**: Custom MarshalJSON/UnmarshalJSON for type-safe serialization
- **Error handling**: Constructors return errors; `Must*` variants panic on invalid input

## Release & CI

- **Tag format:** SemVer (`v0.4.0`, `v0.5.0`, etc.)
- **Release command:** `just release 0.5.0` (creates tag + pushes)
- **GitHub Actions:** CI runs on push to master (test, lint, security, generate, benchmark)
- **Release workflow:** triggers on `v*` tag push — runs tests + lint + creates GitHub Release with git-cliff changelog
- **CI known issue:** GitHub Actions billing is currently failing — all runs fail with billing/spending limit error. This is an account-level issue, not a code problem.

## Notes

- `enums/enum_enum.go` is auto-generated - do not edit manually
- Run `go generate ./...` after modifying `enums/enums.go`
- Generic types (ActorEntry, ActorChain) must be imported from subpackages
- ID types live in `github.com/larsartmann/go-branded-id` (separate module)
- `programminglanguage/` was removed — use [`go-enry`](https://github.com/go-enry/go-enry) for language detection, plain `[]string` in `projectcore.ProjectCore.Languages`
- License: MIT (fixed from PROPRIETARY in 2026-05-07)
- Test coverage: 86.6% overall
- Repo is transitioning from private to public
- **Updated:** 2026-05-07
