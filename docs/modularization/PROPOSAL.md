# Modularization Proposal — go-composable-business-types

**Date:** 2026-05-22
**Status:** Draft
**Approach:** Strategic split around external dependency boundaries

---

## 1. Executive Summary

Split the single-module Go library into 6 semi-independent sub-modules, each with its own `go.mod`, coordinated by a root `go.work` file. The primary motivation is **dependency isolation for consumers** — today, importing just `enums` or `types` pulls in `bojanz/currency`, `sixafter/nanoid`, `golang.org/x/text`, and 30+ transitive dependencies. After the split, most packages carry zero external dependencies.

**Consumer import paths do not change.** The split is purely structural — a consumer writing `import "github.com/larsartmann/go-composable-business-types/nanoid"` today writes the exact same import tomorrow.

---

## 2. Current State Analysis

### 2.1 Module Landscape

| Field           | Value                                                              |
| --------------- | ------------------------------------------------------------------ |
| State           | **Monolith** — single go.mod, no go.work                           |
| Module          | `github.com/larsartmann/go-composable-business-types`              |
| Go version      | 1.26.2                                                             |
| Packages        | 19                                                                 |
| Total LOC       | ~12K prod, ~5.5K test                                              |
| Ext deps (prod) | bojanz/currency, sixafter/nanoid, golang.org/x/text, go-branded-id |
| Ext deps (test) | stretchr/testify (banned per policy)                               |

### 2.2 Package Dependency Graph (Production)

```
Layer 0 (leafs)     Layer 1             Layer 2         Layer 3
─────────────────   ─────────────────   ────────────    ────────────
enums               bounded             temporal        projectcore
pkg/errors          importance          money           datapoint
scanutil            locale              actor
validate            nanoid
testutil            tag
version             types
```

### 2.3 Coupling Assessment

- **No circular dependencies** — the graph is a clean DAG
- **No god-packages** — every package has a single focused responsibility
- **datapoint is the heaviest composite** — imports 5 other packages (by design)
- **pkg/errors, scanutil, validate are shared foundations** — zero internal deps, imported by many

### 2.4 Dependency Bloat Today

| Consumer needs   | Current transitive deps pulled in                                 |
| ---------------- | ----------------------------------------------------------------- |
| Just `enums`     | bojanz/currency, sixafter/nanoid, golang.org/x/text, 30+ indirect |
| Just `nanoid`    | bojanz/currency, golang.org/x/text, 30+ indirect                  |
| Just `types`     | bojanz/currency, sixafter/nanoid, golang.org/x/text, 30+ indirect |
| Just `money`     | sixafter/nanoid (unnecessary), 30+ indirect                       |
| Just `datapoint` | everything (expected — it's the composite)                        |

---

## 3. Proposed Module Structure

### 3.1 Module Definitions

| #   | Module Path                    | Directory      | Packages                                                                                                                | Ext Deps (prod)   | Internal Deps           |
| --- | ------------------------------ | -------------- | ----------------------------------------------------------------------------------------------------------------------- | ----------------- | ----------------------- |
| 1   | `go-composable-business-types` | `./`           | enums, validate, pkg/errors, scanutil, testutil, version, bounded, importance, tag, types, temporal, actor, projectcore | go-branded-id     | —                       |
| 2   | `.../nanoid`                   | `./nanoid/`    | nanoid                                                                                                                  | sixafter/nanoid   | root                    |
| 3   | `.../locale`                   | `./locale/`    | locale                                                                                                                  | golang.org/x/text | root                    |
| 4   | `.../money`                    | `./money/`     | money                                                                                                                   | bojanz/currency   | root, locale            |
| 5   | `.../datapoint`                | `./datapoint/` | datapoint                                                                                                               | —                 | root, nanoid            |
| 6   | `.../examples`                 | `./examples/`  | basic, datapoint                                                                                                        | —                 | root, nanoid, datapoint |

### 3.2 Why These Boundaries

- **Root (13 packages):** The base module with zero heavy external deps. Contains all leaf and mid-layer packages. Only external dep is `go-branded-id` (lightweight, same author, no transitive deps). A consumer importing only `enums`, `types`, `validate`, etc. gets a near-zero dependency footprint.
- **nanoid:** Isolates `sixafter/nanoid` + crypto transitive deps (aes-ctr-drbg, prng-chacha, x/crypto). Consumers who don't use NanoID avoid all crypto deps.
- **locale:** Isolates `golang.org/x/text`. Consumers who don't need locale avoid the x/text dep.
- **money:** Isolates `bojanz/currency` (heaviest dep: brings cockroachdb/apd, shopspring/decimal, lib/pq). Only consumers who need Money get this.
- **datapoint:** Must be separate because it depends on nanoid module (can't be in root — would create a cycle).
- **examples:** Depends on multiple modules. Kept separate to avoid polluting library modules with example code deps.

### 3.3 DAG Verification

```
root ← (nothing)
  ↑
  ├── nanoid → root
  ├── locale → root
  ├── money → root, locale
  ├── datapoint → root, nanoid
  └── examples → root, nanoid, datapoint
```

All edges point upward. No cycles. ✓

### 3.4 Dependency Isolation After Split

| Consumer needs   | Dependencies after split                                |
| ---------------- | ------------------------------------------------------- |
| Just `enums`     | **go-branded-id only** (via root)                       |
| Just `nanoid`    | **sixafter/nanoid + go-branded-id**                     |
| Just `types`     | **go-branded-id only** (via root)                       |
| Just `locale`    | **golang.org/x/text + go-branded-id**                   |
| Just `money`     | **bojanz/currency + golang.org/x/text + go-branded-id** |
| Just `datapoint` | **everything** (expected — it composes everything)      |

---

## 4. Replace / Workspace Strategy

**Approach:** `go.work` at repo root + `replace` directives in sub-module go.mod files for local development.

### go.work

```go
// go.work
go 1.26.3

use (
    .
    ./nanoid
    ./locale
    ./money
    ./datapoint
    ./examples
)
```

### Sub-module replace directives

Each sub-module includes `replace` directives for:
1. **Root module** (`go-composable-business-types => ../`) — required because the published root v0.4.0 still contains all packages, creating ambiguous imports with the local sub-module directories
2. **Sibling modules** (e.g., `nanoid => ../nanoid`) — required for unpublished sub-modules that have no versioned tag

**Rationale:**

- `go.work` provides workspace-level module resolution for builds
- `replace` directives in sub-modules ensure `go mod tidy` works correctly — without them, Go resolves against published v0.4.0 and sees ambiguous imports (same package in both root v0.4.0 and local sub-module)
- `replace` directives are ignored by consumers of published modules — they only apply during local development
- Once a new root version is published (without the split-out packages), the root `replace` can be removed
- Once sub-modules get their first versioned tags, the sibling `replace` directives can be removed

---

## 5. Test Dependency Isolation

| Module    | Test deps from other modules                            |
| --------- | ------------------------------------------------------- |
| Root      | testutil (internal), enums (internal) — all within root |
| nanoid    | testutil (from root)                                    |
| locale    | testutil (from root)                                    |
| money     | locale (from locale module)                             |
| datapoint | actor, enums, nanoid, temporal, types (from deps)       |
| examples  | none                                                    |

No bidirectional test dependencies. ✓

---

## 6. Versioning Strategy

**Independent semver per module.**

- Git tag format: `nanoid/v0.4.0`, `locale/v0.4.0`, `money/v0.4.0`, `datapoint/v0.4.0`, `examples/v0.4.0`
- Root module: `v0.4.0`
- All modules initially released at the same version (0.4.0)
- Release script creates all tags atomically

**Rationale:** Independent semver is the standard for Go sub-modules. It allows nanoid (stable) to not bump when datapoint (active development) changes. Initially, all modules share the same version for simplicity.

---

## 7. Migration Strategy

1. Create feature branch `modularize/split-modules`
2. Create `go.work` at root
3. Extract nanoid module (simplest, 1 package)
4. Extract locale module
5. Extract money module (depends on locale)
6. Extract datapoint module (depends on nanoid + root)
7. Extract examples module
8. Update root module go.mod (remove separated packages)
9. Update flake.nix for multi-module builds
10. Update CI/CD for multi-module testing
11. Update documentation (README, AGENTS.md)
12. Verify full build + test

Each step leaves the project in a buildable, testable state.

---

## 8. Risk Assessment

| Risk                                            | Likelihood      | Mitigation                                               |
| ----------------------------------------------- | --------------- | -------------------------------------------------------- |
| Import path confusion                           | Low             | Import paths don't change — verified                     |
| go.work conflicts with CI                       | Medium          | Test CI with GOWORK=off and with go.work                 |
| Version tag complexity                          | Medium          | Automate with release script                             |
| testify in root go.mod                          | High (existing) | Replace with ginkgo/gomega (separate task, not blocking) |
| Consumer confusion about which module to import | Low             | README clearly documents module boundaries               |

---

## 9. Build System Impact

### flake.nix

- Update build check to build each module independently
- Update test check to test each module independently
- Add per-module checks or aggregated root-level checks

### CI/CD (GitHub Actions)

- Update test job to run `go test ./...` in each module directory
- Consider parallel per-module CI jobs
- Update release workflow for multi-module tagging

### justfile

- Update release command to create multiple tags
- Add module-specific build/test commands

---

## 10. Key Decisions

1. **6 modules** — not too many, not too few. Each external dep isolated.
2. **go.work over replace** — cleaner, standard Go practice.
3. **Independent semver** — future-proof, standard for Go sub-modules.
4. **Consumer import paths unchanged** — zero migration for consumers.
5. **datapoint must be separate** — it depends on nanoid module, creating a module-level dep that prevents it from staying in root.

---

## 11. Self-Review (Phase 4)

**Date:** 2026-05-22

- **Split brains:** None — each type in exactly one module. No duplication across module boundaries.
- **Granularity:** 6 modules for 12K LOC is justified. Each has clear purpose and isolated external deps.
- **go.work feasibility:** Verified import paths don't change. go.work `use` directives point to correct directories.
- **Test dep isolation:** testify (banned) stays only in root module. Separated modules are testify-free.
- **Build system:** flake.nix currently has `GOWORK = "off"` — must be removed. GOEXPERIMENT=jsonv2 needed in all module builds.
- **Code generation:** go-enum tool dep stays in root (enums package stays in root). No impact on generation.
- **Banned deps:** testify in 3 root test files (importance, tag, projectcore). Separate task to replace with ginkgo/gomega. Not blocking modularization.
- **CI value:** Marginal speed improvement. Primary value is dependency isolation for consumers.
- **Alternative considered:** Not modularizing at all. Rejected because dependency bloat for consumers is real (30+ transitive deps when only needing enums).
- **Alternative considered:** Fewer modules (root + nanoid only). Rejected because datapoint depends on nanoid, forcing at least root + nanoid + datapoint split. Adding locale and money is minimal extra overhead for maximum dep isolation.
