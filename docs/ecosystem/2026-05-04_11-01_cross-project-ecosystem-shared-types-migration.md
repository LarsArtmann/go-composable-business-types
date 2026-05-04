# Cross-Project Ecosystem: Shared Types Migration — Status Report

**Date:** 2026-05-04 11:01  
**Session:** Deep architecture review + Phase 1 implementation  
**Scope:** `go-composable-business-types`, `project-discovery-sdk`, `project-meta`, `project-dependency-graph`, `projects-management-automation`

---

## Executive Summary

Researched and mapped the relationships between four Go projects in the LarsArtmann ecosystem. Discovered **3 split-brain type systems** (Importance, Tag, Language), **zero CI cross-project safety**, and a **type babel** of 4 incompatible "Project" representations. Executed Phase 1: added 4 new packages to `go-composable-business-types` as the shared type kernel. **47 new tests, 5 commits, all pushed.**

---

## A) FULLY DONE ✅

### Research & Analysis (100%)

| Deliverable | Status | Location |
|---|---|---|
| Four-project architecture survey | ✅ Complete | Each project's structure, go.mod, domain model fully mapped |
| Cross-project dependency graph | ✅ Complete | Current dependency direction mapped: SDK→PMA, SDK→PDG, Meta→PMA |
| Type duplication analysis | ✅ Complete | Importance (3 versions), Tag (2 versions), Language (3 representations), Project (4 representations) |
| Existing library audit | ✅ Complete | Discovered `go-composable-business-types` already has Percentage, BoundedString, validate.Validator, pkg/errors |
| go-branded-id capabilities | ✅ Complete | Pure phantom wrapper, no validation, JSON/SQL/Text encoding |
| govalid feasibility | ✅ Complete | Can't validate id.ID structs, can't do regex — **dropped from plan** |
| Backward compatibility analysis | ✅ Complete | YAML/JSON serialization safe for int32→uint8, Tag regex needs migration |
| BDD test impact analysis | ✅ Complete | 3/8 feature files reference Language strings ("go", "node", "python", "typescript") |
| Bug discovery: duplicate sdkProjectToDomain | ✅ Complete | Found case-handling divergence between two copies |
| Bug discovery: BrandedProject sync risk | ✅ Complete | Manual field-by-field mapping can silently drop data |
| Architecture diagrams — current state | ✅ Complete | `docs/architecture-understanding/2026-05-04_08-56-cross-project-current.d2/.svg` |
| Architecture diagrams — ideal state (v3) | ✅ Complete | `docs/architecture-understanding/2026-05-04_08-56-cross-project-ideal.d2/.svg` |
| LIBRARY_GUIDE.md updated | ✅ Complete | Added project types row, updated layer map, updated go-composable-business-types section |
| Plan document (v3) | ✅ Complete | `docs/planning/2026-05-04_08-56-project-types-creation-plan.md` — 35 steps across 5 phases |

### Phase 1: Foundation Types in go-composable-business-types (100%)

| Package | Files | Tests | Commit | Pushed |
|---|---|---|---|---|
| `importance/` | importance.go, importance_test.go | 18 tests PASS | `5b4acf1` | ✅ |
| `tag/` | tag.go, tag_test.go | 11 tests PASS | `be7eb37` | ✅ |
| `programminglanguage/` | language.go, language_test.go | 12 tests PASS | `955a82f` | ✅ |
| `projectcore/` | project_core.go, project_core_test.go | 6 tests PASS | `541d4db` | ✅ |
| `pkg/errors/` | errors.go updated | existing + 3 new helpers PASS | `fb4b025` | ✅ |

**Total: 47 new tests, 5 commits, all green, all pushed to master.**

---

## B) PARTIALLY DONE 🔶

### Plan v3 Execution — 5 of 35 steps complete (14%)

| Phase | Steps | Done | Status |
|---|---|---|---|
| Phase 1: Create types in go-composable-business-types | 7 | 5 | Missing: flake.nix (0.1), ADR-001 (0.2) |
| Phase 2: Wire into project-discovery-sdk | 11 | 0 | Not started |
| Phase 3: Wire into project-meta | 5 | 0 | Not started |
| Phase 4: Wire into projects-management-automation | 8 | 0 | Not started |
| Phase 5: Wire into project-dependency-graph | 3 | 0 | Not started |
| Cross-cutting fixes | 1 | 0 | Not started |

---

## C) NOT STARTED ⬜

| # | Task | Phase | Impact | Effort |
|---|---|---|---|---|
| 0.1 | Create `flake.nix` for go-composable-business-types | Pre | MED | 1hr |
| 0.2 | Write ADR-001: Adding project ecosystem types | Pre | MED | 30min |
| 2.1 | Add go-composable-business-types to SDK go.mod (replace→local) | P2 | MED | 10min |
| 2.2 | Replace SDK Project.Language/Languages with programminglanguage types | P2 | HIGH | 1hr |
| 2.3 | Replace SDK Project.Importance int with importance.Importance | P2 | HIGH | 30min |
| 2.4 | Replace SDK Project.Tags []string with []tag.Tag | P2 | HIGH | 30min |
| 2.5 | Remove SDK's NormalizeLanguage → use programminglanguage.Normalize | P2 | MED | 30min |
| 2.6 | Remove SDK's branded.go types that now live in project-types | P2 | MED | 30min |
| 2.7 | Update SDK tests | P2 | MED | 1hr |
| 2.8 | Update/delete BrandedProject parallel type | P2 | MED | 30min |
| 2.9 | Update all 7 SDK filter functions for new types | P2 | MED | 30min |
| 2.10 | Update Event struct and all 12 event hook call sites | P2 | MED | 30min |
| 2.11 | Deduplicate sdkProjectToDomain (fix case bug) | P2 | MED | 15min |
| 3.1 | Add go-composable-business-types to project-meta go.mod | P3 | MED | 10min |
| 3.2 | Replace project-meta's Importance int32 with importance.Importance uint8 | P3 | HIGH | 1hr |
| 3.3 | Replace project-meta's Tag with tag.Tag | P3 | HIGH | 30min |
| 3.4 | Add backward-compat YAML unmarshaling (int32 + uint8) | P3 | MED | 30min |
| 3.5 | Update all project-meta tests | P3 | MED | 1hr |
| 4.1 | Add go-composable-business-types to PMA go.mod + go.work | P4 | MED | 10min |
| 4.2 | Replace PMA's Importance int32 with importance.Importance uint8 | P4 | HIGH | 1hr |
| 4.3 | Replace PMA's types.Tag with tag.Tag | P4 | HIGH | 30min |
| 4.4 | Replace PMA's types.Language uint8 enum with programminglanguage.Language | P4 | HIGH | 2hr |
| 4.5 | Convert language switch statements to dispatch maps | P4 | MED | 1hr |
| 4.6 | DELETE meta/adapter.go — no type conversion needed | P4 | HIGH | 15min |
| 4.7 | Fix Importance type leaks (raw int/int32 in API, CLI layers) | P4 | MED | 30min |
| 4.8 | Update all PMA tests | P4 | MED | 2hr |
| 5.1 | Add go-composable-business-types to PDG go.mod | P5 | MED | 10min |
| 5.2 | Use programminglanguage.Language for filtering | P5 | LOW | 15min |
| 5.3 | Full PDG test verification | P5 | LOW | 30min |
| 6.1 | Create root-level go.work or add to PMA's workspace | P4 | MED | 15min |

---

## D) TOTALLY FUCKED UP 💥

### 1. Pre-existing enums test failure in go-composable-business-types
`TestAllEnumScanAllTypes` in `enums/` fails: `expected Active, got Archived`. This is NOT caused by our changes — confirmed by testing on clean checkout. **Needs separate fix.**

### 2. Zero CI across all four projects
No GitHub Actions workflows exist in project-discovery-sdk, project-meta, or project-dependency-graph. PMA has workflows but they weren't checked. **A breaking change to the SDK's types won't be caught by any automated pipeline in the consuming projects.**

### 3. No root-level go.work
Only PMA has a `go.work` with local replaces. There's no `/home/lars/projects/go.work`. This means cross-project type changes won't compile-check across all consumers simultaneously.

### 4. go-composable-business-types has no flake.nix
There's a migration proposal document but no actual `flake.nix`. Can't use `nix build` / `nix run .#test` for this project.

---

## E) WHAT WE SHOULD IMPROVE 🔧

### Plan Quality Improvements

1. **Pivot from new library to existing one** — v1 planned a new `project-types` library. v3 correctly identified that `go-composable-business-types` already has all the infrastructure. Saved weeks of duplicated work.

2. **Dropped govalid** — v1 planned to use govalid. Research showed it can't validate `id.ID` structs or do regex. The existing `validate.Validator` pattern in the library is sufficient. Correct call.

3. **Tag regex change needs migration path** — Old `^[a-z0-9_-]+$` → New `^[A-Za-z0-9-]+$`. The plan mentions `doctor --fix-tags` but doesn't implement it. Needs a concrete step in Phase 3.

4. **Importance uint8 not int32** — Corrected in plan. uint8 makes negative values impossible by type. But existing PMA code leaks raw `int32` and `int` in 5+ places (api/list.go, cmdguard_meta_sync.go, metadata_registry.go, cli_output.go, test steps). These must all be found and fixed.

5. **Language enum→branded string leaves dispatch gap** — PMA has ~6 switch statements on Language enum values. The plan says "convert to dispatch maps" but doesn't show the pattern. Need a concrete example before Phase 4.

### Architecture Improvements

6. **ProjectType extraction from Language enum** — PMA's `LanguageWeb`, `LanguageMobile`, `LanguageCLI` etc. are NOT languages. They must become a separate `ProjectType` concept. But where does it live? PMA-only, or in the kernel?

7. **SDK→PMA conversion is lossy** — 8 fields dropped silently (HasGit, Branch, IsDirty, LastCommit, Markers, RemoteURL, Host, Organization, Repository). The plan's `ProjectCore` doesn't address this. Should we add optional enrichment fields?

8. **BrandedProject in SDK is unmaintained** — Parallel type to sdk.Project with manual field mapping. Should be deleted or auto-generated after Phase 2.

9. **No shared context document** — The four projects have no shared architecture document. Should create a `docs/ECOSYSTEM.md` or similar.

10. **SDK event system tightly couples to Project struct** — All 12 event types embed `*Project`. Any field change propagates to every consumer. Consider using interfaces or event-specific DTOs.

---

## F) TOP 25 THINGS TO GET DONE NEXT

Sorted by impact × effort (highest first):

| Priority | Task | Why |
|---|---|---|
| **1** | Phase 2: Wire types into `project-discovery-sdk` | Unblocks all other consumers; SDK is the foundation |
| **2** | Phase 2: Deduplicate `sdkProjectToDomain` (fix case bug) | Active bug causing silent data loss |
| **3** | Phase 2: Update/delete `BrandedProject` parallel type | Prevents silent data loss if fields are forgotten |
| **4** | Phase 2: Update all 7 SDK filter functions | Filters will break on type change |
| **5** | Phase 2: Update Event struct + 12 hook call sites | Events embed *Project, must reflect new types |
| **6** | Phase 3: Wire types into `project-meta` | Metadata authority — must adopt shared Importance/Tag |
| **7** | Phase 3: Add backward-compat YAML unmarshaling | Existing .config/metadata.yaml files must not break |
| **8** | Phase 3: Implement `doctor --fix-tags` migration | Tag regex change requires data migration |
| **9** | Phase 4: Wire types into `projects-management-automation` | Largest consumer, highest risk, highest reward |
| **10** | Phase 4: DELETE `meta/adapter.go` | The entire file becomes unnecessary with shared types |
| **11** | Phase 4: Replace Language uint8 enum with branded string | Biggest conceptual change in PMA |
| **12** | Phase 4: Extract `ProjectType` from Language enum | Clean up the `LanguageWeb`/`LanguageMobile` naming collision |
| **13** | Phase 4: Convert switch statements to dispatch maps | Required for Language string migration |
| **14** | Phase 4: Fix Importance type leaks (raw int/int32) | 5+ files use raw types instead of domain type |
| **15** | Phase 5: Wire types into `project-dependency-graph` | Low effort, completes the ecosystem |
| **16** | Create `flake.nix` for go-composable-business-types | Enables nix-based build/test workflow |
| **17** | Write ADR-001: Project ecosystem types decision | Documents the architectural decision |
| **18** | Fix pre-existing enums test failure | Test suite should be fully green |
| **19** | Create root-level `go.work` or expand PMA's workspace | Compile-time cross-project safety |
| **20** | Add CI workflow to project-discovery-sdk | Automated testing on push |
| **21** | Add CI workflow to project-meta | Automated testing on push |
| **22** | Unify go-filewatcher version (v0.1.0 vs v0.2.0) | Consistency across projects |
| **23** | Make project-dependency-graph importable as library | Split pkg/graph/ from cmd/pdg/ |
| **24** | Connect PDG to project-meta for node enrichment | Graphs show importance/tags |
| **25** | Create shared `docs/ECOSYSTEM.md` | Document inter-project relationships |

---

## G) MY TOP #1 QUESTION I CANNOT FIGURE OUT MYSELF 🤔

**`ProjectCore.Name` and `ProjectCore.Path` — what validation, if any?**

Right now `ProjectCore` uses plain `string` for both. But:

- **PMA's `Name`** allows `"My Cool Project"` (uppercase, spaces, max 255) — it's a display name
- **project-meta's `ProjectName`** requires `^[a-z][a-z0-9._-]*[a-z0-9]$` (lowercase, no spaces, max 100) — it's a filesystem identifier
- **SDK's `Project.Name`** returns the directory name (raw string from `filepath.Base`)

These are three different concepts sharing the same field name. The kernel shouldn't pick a winner, but leaving them as untyped `string` means no compile-time safety at the boundary.

**Should `ProjectCore.Name` be:**
- (a) Plain `string` (current — no validation, each consumer validates as needed)
- (b) A branded type with minimal validation (non-empty, max 255)
- (c) Left out of `ProjectCore` entirely (each project defines its own name type)

This is a design decision that affects all four consumers and I cannot resolve it without your input.

---

## Session Metrics

| Metric | Value |
|---|---|
| Projects analyzed | 4 |
| Files read | ~80+ |
| Split brains found | 3 (Importance, Tag, Language) |
| Bugs discovered | 2 (case-handling divergence, BrandedProject sync risk) |
| New packages created | 4 (importance, tag, programminglanguage, projectcore) |
| New tests written | 47 (all passing) |
| Commits pushed | 5 |
| Plan iterations | 3 (v1: new library, v2: +govalid, v3: extend existing) |
| Diagrams generated | 2 (current state + ideal state v3) |
| Phase completion | Phase 1: 5/7 steps done (71%) |
| Total plan completion | 5/35 steps (14%) |

---

_Generated with Crush — 2026-05-04_
