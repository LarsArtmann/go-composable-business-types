# Project: go-composable-business-types

A Go library of strongly typed, composable base values for business applications.

## Project Goal

Build the best possible data types for Lars Artmann's Golang applications, leveraging superb existing libraries where appropriate. **Backwards compatibility is NOT a concern** — we prioritize clean, optimal design over legacy support.

## Multi-Module Workspace

This library uses a Go workspace (`go.work`) with 6 semi-independent sub-modules for dependency isolation:

| Module    | Directory      | External Deps         | Purpose                                                                                                                                               |
| --------- | -------------- | --------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- |
| Root      | `./`           | go-branded-id         | Base types: enums, validate, pkg/errors, scanutil, testutil, version, bounded, importance, tag, types, temporal, actor, projectcore, address, contact |
| nanoid    | `./nanoid/`    | sixafter/nanoid       | URL-safe unique identifiers                                                                                                                           |
| locale    | `./locale/`    | golang.org/x/text     | BCP 47 language tags                                                                                                                                  |
| money     | `./money/`     | bojanz/currency       | ISO 4217 currency                                                                                                                                     |
| datapoint | `./datapoint/` | go-branded-id, nanoid | Self-contained data units with audit trail                                                                                                            |
| examples  | `./examples/`  | —                     | Usage examples                                                                                                                                        |

Consumer import paths are identical to the single-module era. The split only affects dependency isolation.

## Build & Test Commands

All build automation lives in `flake.nix`. Use Nix flakes, never Makefile.

| Task         | Command                    |
| ------------ | -------------------------- |
| Dev shell    | `nix develop`              |
| Build        | `nix build .#check-build`  |
| Test         | `nix build .#check-test`   |
| Lint         | `nix build .#check-lint`   |
| Format check | `nix build .#check-format` |

Manual (non-Nix) equivalents require `GOEXPERIMENT=jsonv2` (see Gotchas):

```bash
GOEXPERIMENT=jsonv2 go build ./...
GOEXPERIMENT=jsonv2 go test -race ./...
golangci-lint run ./...
go generate ./...
```

Each sub-module must be built/tested separately:

```bash
for d in nanoid locale money datapoint examples; do (cd $d && GOEXPERIMENT=jsonv2 go test -race ./...); done
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
// Import enums
import "github.com/larsartmann/go-composable-business-types/enums"

func main() {
    kind := enums.ActorKindUser
    trigger := enums.TriggerManual
}
```

## Code Conventions

- **Phantom types**: Strong type safety with branded types (NanoID, Email, etc.)
- **Functional patterns**: Immutable value types with `With*` methods returning copies
- **JSON serialization**: Custom MarshalJSON/UnmarshalJSON for type-safe serialization
- **Error handling**: Constructors return errors; `Must*` variants panic on invalid input

## Release & CI

- **Tag format:** Root: SemVer (`v0.6.0`). Sub-modules: `nanoid/v0.6.0`, `locale/v0.6.0`, etc.
- **Release:** Create tags manually and push (justfile was removed in favor of flake.nix)
- **GitHub Actions:** CI runs on push to master (test, lint, security, generate, benchmark)
- **Release workflow:** triggers on `v*` and `*/v*` tag pushes
- **CI known issue:** GitHub Actions billing is currently failing — all runs fail with billing/spending limit error. This is an account-level issue, not a code problem.

## Notes

- `enums/enums_enum.go` is auto-generated - do not edit manually
- Run `go generate ./...` after modifying `enums/enums.go`
- Generic types (ActorEntry, ActorChain) must be imported from subpackages
- ID types live in `github.com/larsartmann/go-branded-id` (separate module)
- `programminglanguage/` was removed — use [`go-enry`](https://github.com/go-enry/go-enry) for language detection, plain `[]string` in `projectcore.ProjectCore.Languages`
- License: MIT (fixed from PROPRIETARY in 2026-05-07)
- Repo is transitioning from private to public
- Modularization docs: `docs/modularization/` (PROPOSAL.md, DEPENDENCY_GRAPH.md, EXECUTION_PLAN.md)
- `go mod tidy` in sub-modules requires `replace` directives (present in each sub-module go.mod) to resolve the root module locally.
- **Gotcha:** Code uses `encoding/json/v2` which requires `GOEXPERIMENT=jsonv2` on Go 1.26.x, or Go 1.27+. The `go.mod` currently says `go 1.26.4`. Without this env var, builds fail with "build constraints exclude all Go files".
