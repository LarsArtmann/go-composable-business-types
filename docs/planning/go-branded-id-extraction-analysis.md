# PRO/CONTRA Analysis: Extracting `id/` and `nanoid/` to `go-branded-id`

**Date:** 2026-04-05
**Scope:** Move `id/` and `nanoid/` packages from `go-composable-business-types` to a dedicated `go-branded-id` repository

---

## PRO (Arguments For Extraction)

| #   | Benefit                                 | Details                                                                                                                             |
| --- | --------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------- |
| 1   | **Focused Domain**                      | A library dedicated to "branded IDs" has clearer purpose than one of 15 types in "business types". Single-responsibility principle. |
| 2   | **Independent Versioning**              | ID patterns are stable; they change rarely compared to DataPoint, money, or temporal types. Separate release cadence makes sense.   |
| 3   | **Reduced Coupling for ID-Heavy Users** | Applications needing only ID safety (not DataPoint, Money, etc.) can depend on a minimal, focused package.                          |
| 4   | **NanoID is a Natural Companion**       | Both are about creating/handling identifiers. They compose naturally: `ID[Brand, NanoID]`.                                          |
| 5   | **Simpler Dependency Graph**            | `id/` currently has **zero** external dependencies. Clean, extractable unit.                                                        |
| 6   | **Easier Adoption**                     | Users can adopt "branded IDs" without buying into the full "composable business types" philosophy.                                  |
| 7   | **NanoID Already Conceptually Split**   | It has its own README, tests, and error definitions — extraction formalizes this.                                                   |
| 8   | **Go Best Practice**                    | Small, focused modules align with Go's philosophy of minimal dependency trees.                                                      |

---

## CONTRA (Arguments Against Extraction)

| #   | Risk                                       | Details                                                                                                                                  |
| --- | ------------------------------------------ | ---------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | **New Repo Maintenance Overhead**          | Separate CI/CD, releases, security scans, dependabot configs, issue tracking.                                                            |
| 2   | **Version Coordination Complexity**        | When to bump `go-branded-id` vs `go-composable-business-types`? Which depends on which?                                                  |
| 3   | **Internal Dependency Chain**              | `nanoid/` depends on `scanutil/` and `pkg/errors/`. These must either (a) extract too, (b) be duplicated, or (c) `nanoid/` stays in cbt. |
| 4   | **`actor/` Depends on `id/`**              | `ActorEntry[T]` uses `id.ID[struct{}, T]`. Extracting `id/` means actor must import external package.                                    |
| 5   | **`datapoint/` Depends on `nanoid/`**      | `DataPoint.ID()` returns `nanoid.NanoID`. Extracting nanoid means datapoint must import external package.                                |
| 6   | **Two Import Paths to Maintain**           | Documentation, examples, tests need updating in both repos. Risk of stale docs.                                                          |
| 7   | **Circular Dependency Risk**               | If `go-branded-id` ever needs something from cbt, we're stuck with a circular import.                                                    |
| 8   | **PROJECT_SPLIT_EXECUTIVE_REPORT Says No** | Previous analysis concluded splitting creates "dependency hell" and "fragments user experience".                                         |

---

## Key Decision Points

### What to Extract?

```
Minimal extraction: id/ only
├── PRO: Clean, focused, zero deps
└── CON: NanoID stays behind (inconsistent)

Recommended extraction: id/ + nanoid/ + scanutil/ + errors/ (ID/NanoID errors only)
├── PRO: Complete ID ecosystem in one place
├── PRO: scanutil and ID errors are infrastructure for these types
└── CON: Larger scope, more coordination
```

### Dependency Impact

| Package      | Depends on                   | After Extraction                      |
| ------------ | ---------------------------- | ------------------------------------- |
| `actor/`     | `id/`                        | Imports `go-branded-id/id`            |
| `datapoint/` | `nanoid/`, `id/` (via actor) | Imports `go-branded-id/nanoid` + `id` |
| `nanoid/`    | `scanutil/`, `pkg/errors/`   | Either extract these too OR redefine  |

### Previous Analysis Context

The `PROJECT_SPLIT_EXECUTIVE_REPORT.md` (current, outdated) recommended **against** splitting, citing:

> "DataPoint depends on: NanoId, ActorEntry, Bitemporal, Context, Reference, Cause, Trigger — splitting creates circular dependencies"

However, extracting only `id/` + `nanoid/` (the "identity" layer) is a **different scope** than the original "split into 6 separate libraries" proposal. This is a **horizontal slice**, not vertical fragmentation.

---

## Recommendation

**Extract `id/` + `nanoid/` + `scanutil/` + `pkg/errors/` (ID/NanoID errors only) into `go-branded-id`**

### Structure in New Repo

```
go-branded-id/
├── id/              # ID[B,V] branded, phantom-type identifiers
├── nanoid/          # NanoID generator (URL-safe, FIPS-140 compatible)
├── scanutil/        # SQL scan helpers (minimal, reusable)
├── errors/          # Domain errors for ID/NanoID only
├── go.mod
└── README.md
```

### Rationale

1. **Orthogonal Layer**: These packages exist at the "type system" layer, not business logic layer
2. **Natural Pairing**: `ID[Brand, V]` + `NanoID` are identifier primitives that compose well
3. **Infrastructure Belongs Together**: `scanutil` and errors are utilities for these types
4. **Remaining cbt is Business Logic**: `go-composable-business-types` becomes a higher-level orchestration layer

### Compatibility Path

```go
// In go-composable-business-types, temporary re-export for migration:
package cbt

import (
    "github.com/larsartmann/go-branded-id/id"
    "github.com/larsartmann/go-branded-id/nanoid"
)

type ID = id.ID  // Re-export for backwards compat
type NanoID = nanoid.NanoID
```

### Final Verdict

| Aspect                    | Assessment                                            |
| ------------------------- | ----------------------------------------------------- |
| **Technical Feasibility** | High — clean dependency boundaries                    |
| **Maintenance Overhead**  | Moderate — new repo, but stable code                  |
| **User Benefit**          | High — focused library, clearer API                   |
| **Risk**                  | Low-to-Moderate — can maintain compat with re-exports |

**Proceed with extraction** if the maintenance overhead is acceptable and there's a clear ownership/CI plan for the new repository.
