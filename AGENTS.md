# Project: go-composable-business-types

A Go library of strongly typed, composable base values for business applications.

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
- `github.com/matoous/go-nanoid` - NanoId generation

## Project Structure

```
.
├── actor.go            # ActorChain[T], ActorEntry[T] - audit trail tracking
├── bounded.go          # BoundedString - length-validated strings
├── common.go           # Email, URL, Percentage, Cents, Timestamp, Duration
├── datapoint.go        # DataPoint[T] - self-contained data with audit trail
├── datapoint_*.go      # DataPoint supporting types (Cause, Context, Reference, Temporal)
├── enum.go             # ActorKind, Locale, Priority, Status, Trigger enums
├── enum_enum.go        # Generated enum code (do not edit)
├── money.go            # Money wrapper around bojanz/currency
└── nanoid.go           # NanoId type
```

## Code Conventions

- **Phantom types**: Strong type safety with branded types (NanoId, Email, etc.)
- **Functional patterns**: Immutable value types with `With*` methods returning copies
- **JSON serialization**: Custom MarshalJSON/UnmarshalJSON for type-safe serialization
- **Error handling**: Constructors return errors; `Must*` variants panic on invalid input

## Notes

- `enum_enum.go` is auto-generated - do not edit manually
- Run `go generate ./...` after modifying `enum.go`
