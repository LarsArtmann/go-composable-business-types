# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [Unreleased]

## [0.7.0] - 2026-07-23

> Sub-modules (`nanoid`, `locale`, `money`, `datapoint`, `examples`) ship as **v0.6.0** in this cycle.

### Changed

- Migrated to Go 1.27 `encoding/json/v2` across all packages (requires `GOEXPERIMENT=jsonv2` on Go 1.26.x)
- Centralized JSON marshal/unmarshal through `types.MarshalJSON`/`types.UnmarshalJSON` helpers (`importance`, `tag`, `temporal`, `projectcore`)
- Adopted `slices.Clone` and idiomatic slice construction throughout
- Updated `go-branded-id` to v0.3.2
- Updated all dependency versions across all modules
- Tightened `golangci-lint` configuration and applied gosec G115 compliance fixes

### Added

- `flake.nix` build automation with `build`, `test`, `lint`, and `format` checks (replaces justfile)
- `FEATURES.md` and `TODO_LIST.md`

### Removed

- Deprecated justfile in favor of `flake.nix`
- Generated HTML report assets from repository

## [0.6.0] - 2026-06-12

### Added

- Shared `Address` type (`address/`) with validated postal address fields
- Shared `Contact` type (`contact/`) with name, email, phone, website, address
- `examples/` module with documentation and working programs
- Nix format and build checks

## [0.5.0] - 2026-05-28

### Added

- Modularized monolith into 6 semi-independent sub-modules for dependency isolation
- Generic `ScanEnum[T]` for iota-based enum `sql.Scanner` support (`scanutil/`)
- `docs/DOMAIN_LANGUAGE.md` for ubiquitous language documentation
- `CONTRIBUTING.md`
- `datapoint/` module as standalone sub-module

### Changed

- Upgraded `go-branded-id` v0.1.0 to v0.3.0
- Migrated `golangci-lint` config to v2 format
- Enhanced `bounded.UnmarshalJSON` error messages

### Removed

- `programminglanguage/` package — use [go-enry](https://github.com/go-enry/go-enry) for language detection; plain `[]string` in `projectcore.ProjectCore.Languages`

## [0.4.0] - 2026-05-07

### Changed

- Transitioned to open-source MIT license (fixed from PROPRIETARY)
- Switched release tags from date-based to SemVer format
- Resolved root module dependencies from GitHub (removed local `replace` scaffolding)

### Added

- `POLICY.md` — versioning, breaking changes, deprecation, security, contribution policies
- `SUPPORT.md` — user support guide
- `SECURITY.md` — security policy
- GitHub issue templates and PR template
- Dependabot configuration
- CI security scanning via govulncheck

## [0.3.0] - 2026-05-04

### Added

- `projectcore/` package — composite project metadata type
- `tag/` package — validated string label with `^[A-Za-z0-9-]+$` constraint
- `importance/` package — uint8 0-100 with named classification levels
- `Tags` slice type with `Strings`, `IsEmpty`, `Contains` helpers
- Importance arithmetic, comparison, and helper methods
- Sentinel errors for importance, tag, and projectcore types
- `flake.nix` for reproducible Nix builds
- Ecosystem architecture docs and cross-project planning

## [0.2.0] - 2026-04-30

### Changed

- Extracted `ID[B, V]` into standalone [`go-branded-id`](https://github.com/larsartmann/go-branded-id) library
- Renamed nanoid functions to idiomatic Go names

### Added

- Encoding interface fallbacks for non-primitive ID value types
- Nix flakes migration proposal

## [0.1.0] - 2026-04-05

### Added

- Initial release with core types: enums, types, bounded, actor, temporal, nanoid, locale, money
- `pkg/errors` centralized error definitions
- `validate` Validator interface
- `scanutil` SQL scan helpers
- `version` build metadata
- `testutil` test utilities
