# Project: go-composable-business-types

A Go library of strongly typed, composable base values for business applications.

## Project Goal

Build the best possible data types for Lars Artmann's Golang applications, leveraging superb existing libraries where appropriate. **Backwards compatibility is NOT a concern** — we prioritize clean, optimal design over legacy support.

## Multi-Module Workspace

This library uses a Go workspace (`go.work`) with 6 semi-independent sub-modules for dependency isolation:

| Module    | Directory      | External Deps     | Purpose                                                                                                                             |
| --------- | -------------- | ----------------- | ----------------------------------------------------------------------------------------------------------------------------------- |
| Root      | `./`           | go-branded-id     | Base types: enums, validate, pkg/errors, scanutil, testutil, version, bounded, importance, tag, types, temporal, actor, projectcore |
| nanoid    | `./nanoid/`    | sixafter/nanoid   | URL-safe unique identifiers                                                                                                         |
| locale    | `./locale/`    | golang.org/x/text | BCP 47 language tags                                                                                                                |
| money     | `./money/`     | bojanz/currency   | ISO 4217 currency                                                                                                                   |
| datapoint | `./datapoint/` | —                 | Self-contained data units with audit trail                                                                                          |
| examples  | `./examples/`  | —                 | Usage examples                                                                                                                      |

Consumer import paths are identical to the single-module era. The split only affects dependency isolation.

## Build & Test Commands

```bash
# Build all modules (via workspace)
go build ./...

# Test all modules (per-module due to workspace structure)
go test -race ./...
for mod in nanoid locale money datapoint examples; do (cd $mod && go test -race ./...); done

# Test a specific module
go test -race ./nanoid/...

# Run tests with coverage
go test -race -coverprofile=coverage.out ./...

# Tidy all modules
go mod tidy  # root
for mod in nanoid locale money datapoint examples; do (cd $mod && go mod tidy); done

# Sync workspace
go work sync

# Lint
golangci-lint run --fix

# Generate enum code (in root module)
go generate ./...
```

## Dependencies

Root module:

- `github.com/abice/go-enum` - Enum code generation (`//go:generate go-enum`)
- `github.com/larsartmann/go-branded-id` - Branded phantom-type identifiers
- `github.com/stretchr/testify` - Test assertions (banned per policy, should be replaced with ginkgo/gomega)

Sub-modules:

- `github.com/sixafter/nanoid` - FIPS-140 compatible NanoID (nanoid module only)
- `golang.org/x/text` - BCP 47 locale support (locale module only)
- `github.com/bojanz/currency` - ISO 4217 currency (money module only)

## Package Structure (Go 1.26 Multi-Module Selective Imports)

This library uses Go workspace mode with 6 sub-modules for dependency isolation.
Consumer import paths are unchanged from the single-module era.

```
Root Module (./) — zero heavy external deps
├── actor/              # ActorChain[T], ActorEntry[T] - audit trail tracking
├── bounded/            # BoundedString - length-validated strings
├── enums/              # ActorKind, Priority, Status, Trigger, CauseKind enums
│   └── enum_enum.go    # Generated enum code (do not edit)
├── importance/         # Importance - 0-100 priority classification
├── pkg/errors/         # Centralized sentinel and structured errors
├── projectcore/        # ProjectCore - composite project metadata
├── scanutil/           # SQL Scanner/Valuer helpers
├── tag/                # Tag - validated string labels
├── temporal/           # Bitemporal - valid/recorded time tracking
├── testutil/           # Generic test helpers
├── types/              # Email, URL, Percentage, Cents, Timestamp, Duration
├── validate/           # Validator interface for self-validating types
└── version/            # Build version info from runtime/debug

Sub-Modules (own go.mod)
├── nanoid/             # NanoID - URL-safe unique identifiers [sixafter/nanoid]
├── locale/             # Locale - BCP 47 language tags [golang.org/x/text]
├── money/              # Money - ISO 4217 currency [bojanz/currency]
├── datapoint/          # DataPoint[T] - self-contained data with audit trail
│   ├── datapoint.go    # DataPoint[T] main type
│   ├── context.go      # Execution context
│   ├── reference.go    # Reference[T] entity references
│   └── cause.go        # Cause[T] causal relationships
└── examples/           # Usage examples
    ├── basic/          # Basic usage example
    └── datapoint/      # DataPoint usage example
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

- **Tag format:** Root: SemVer (`v0.5.0`). Sub-modules: `nanoid/v0.5.0`, `locale/v0.5.0`, etc.
- **Release command:** `just release 0.5.0` (creates all tags + pushes)
- **GitHub Actions:** CI runs on push to master (test, lint, security, generate, benchmark)
- **Release workflow:** triggers on `v*` and `*/v*` tag pushes
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
- Modularization docs: `docs/modularization/` (PROPOSAL.md, DEPENDENCY_GRAPH.md, EXECUTION_PLAN.md)
- `go mod tidy` in sub-modules requires `replace` directives (present in each sub-module go.mod) to resolve the root module locally. The published root v0.4.0 still contains all packages, creating ambiguous imports without `replace`.
- **Updated:** 2026-05-22
