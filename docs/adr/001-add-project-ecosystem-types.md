# ADR-001: Add Project Ecosystem Types to go-composable-business-types

**Date:** 2026-05-04
**Status:** Accepted
**Deciders:** Lars Artmann

## Context

The LarsArtmann ecosystem has four Go projects that share overlapping domain types — `Importance`, `Tag`, and `Language` — but each defines them differently:

| Type       | project-discovery-sdk      | project-meta                               | projects-management-automation                |
| ---------- | -------------------------- | ------------------------------------------ | --------------------------------------------- |
| Importance | `int` (unvalidated)        | `int32` (0–100, thresholds 26–49 = "low")  | `int32` (0–100, thresholds 21–40 = "low")     |
| Tag        | N/A (raw `[]string`)       | `string` (`^[a-z0-9_-]+$`, reject invalid) | `string` (`^[a-z0-9_-]+$`, normalize invalid) |
| Language   | `string` (raw enry output) | N/A                                        | `uint8` enum (27 values incl. non-languages)  |

This causes:

- **3 split-brain type systems** with incompatible thresholds and validation
- **4 incompatible `Project` structs** requiring converter layers at every boundary
- **2 bugs**: duplicate `sdkProjectToDomain` with divergent case handling; `BrandedProject` requiring manual field sync
- Data loss at SDK→PMA boundary (8 fields silently dropped)

## Decision

Add four new packages to `go-composable-business-types` (the existing domain types library), NOT a new `project-types` library.

### Packages

| Package                | Type                                           | Rationale                                                                                                            |
| ---------------------- | ---------------------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| `importance/`          | `uint8` (0–100)                                | Impossible to be negative. Unified classification: VeryLow ≤20, Low 21–40, Medium 41–60, High 61–80, VeryHigh 81–100 |
| `tag/`                 | `string` (validated `^[A-Za-z0-9-]+$`, max 50) | Uppercase allowed (was forbidden). Underscores dropped. Migration: `strings.ReplaceAll(tag, "_", "-")`               |
| `programminglanguage/` | branded `string` via go-branded-id             | NOT a closed enum. go-enry is the authority (400+ languages). Zero maintenance when enry adds languages              |
| `projectcore/`         | composite struct                               | Name, Path, Languages, Importance, Tags — the intersection of all four projects                                      |

### Why extend go-composable-business-types?

1. It already has `Percentage uint8`, `BoundedString`, `validate.Validator`, `pkg/errors` with sentinels, JSON/SQL helpers
2. It already depends on `go-branded-id`
3. One less library to maintain, version, and add to go.work
4. Types are domain primitives, not project-specific logic

### Key design choices

- **Language as branded string, not enum**: User explicitly rejected enum because most projects are polyglot and go-enry already handles 400+ languages well
- **govalid dropped**: Can't validate `id.ID` struct fields, can't do regex, only marginal value
- **ProjectName/ProjectPath NOT in shared types**: They're filesystem identifiers with strict validation, fundamentally different from display names

## Consequences

### Positive

- Single source of truth for Importance, Tag, Language across 4 projects
- Eliminates all type conversion layers at project boundaries
- Fixes the duplicate `sdkProjectToDomain` bug
- Type system prevents invalid values at compile time (`uint8` can't be negative)
- Consistent JSON/SQL serialization across ecosystem

### Negative

- All four projects must migrate (85 tasks, ~15 hours)
- Tag regex change (`_` → `-`) requires migration (`doctor --fix-tags`)
- PMA's `LanguageWeb`/`LanguageMobile` etc. need extraction to `ProjectType`
- Backward-compatible YAML/JSON unmarshaling needed in project-meta

### Risks

- Migration is sequential (SDK → meta → PMA → PDG) due to dependency chain
- Breaking change for any external consumers of these projects (internal-only, so acceptable)

## References

- Architecture diagrams: `docs/ecosystem/2026-05-04_08-56-cross-project-*.d2`
- Master execution plan: `docs/ecosystem/2026-05-04_11-01_MASTER_EXECUTION_PLAN.md`
- Migration status: `docs/ecosystem/2026-05-04_11-01_cross-project-ecosystem-shared-types-migration.md`
