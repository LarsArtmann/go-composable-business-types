# Status Report — Shared Types Ecosystem Migration

**Date:** 2026-05-04 17:00  
**Scope:** go-composable-business-types, project-discovery-sdk, project-meta, projects-management-automation, project-dependency-graph  
**Overall Status:** Migration COMPLETE. Post-migration improvements IN PROGRESS.

---

## A) FULLY DONE

### Phase 0: CBT Pre-requisites ✅
- **Commit:** `383dd60` (pushed)
- Created `importance/` package: uint8 type, constructors, classification thresholds (VeryLow ≤20, Low 21-40, Medium 41-60, High 61-80, VeryHigh 81-100), JSON/SQL marshaling
- Created `tag/` package: string type, regex `^[A-Za-z0-9-]+$`, max 50 runes, JSON/SQL
- Created `programminglanguage/` package: branded string via go-branded-id, `Languages` slice with `.Primary()`, `FromSlice()`, `Normalize()`
- Created `projectcore/` package: shared base entity struct
- Added sentinel errors for importance/tag/projectcore
- Wrote ADR-001 documenting the decision
- Created `flake.nix` for CBT (nix build + test)
- Fixed pre-existing `enums/enums_test.go` bugs (wrong intVal→iota mappings, incorrect String() format expectations)

### Phase 2: project-discovery-sdk ✅
- **Commit:** `8b44eb8` (pushed)
- `Language string` → `programminglanguage.Language` (branded string)
- `Languages []string` → `programminglanguage.Languages` (has `.Primary()`, `.Strings()`, `.Equal()`)
- `Importance int` → `importance.Importance` (uint8)
- `Tags []string` → `[]tag.Tag`
- Removed SDK's own `Language` branded type, `LanguageBrand`, `NewLanguage`, `LanguageFrom`, `LanguagesFrom`, `StringsFromLanguages`
- Updated `BrandedProject` + `AsBranded()`/`AsProject()` converters
- Updated all 7 filter functions, Event struct, 12 hook call sites
- Updated all test files (benchmark, branded, filter, snapshot, sdk, extra, example, discovery BDD, mr/*)
- All 4 test packages pass

### Phase 3: project-meta ✅
- **Commits:** `bc51367`, `da188ac` (pushed)
- `Importance int32` → `importance.Importance` (uint8) via type alias
- `Tag` → `cbtTag.Tag` via type alias
- Removed old threshold constants, regex, `normalizeTag`, custom `IsValid`/`IsVeryLow` methods, `MarshalJSON`/`UnmarshalJSON` on Importance
- Added backward-compat `MarshalJSON`/`UnmarshalJSON` on Metadata struct
- Updated test expectations: String format `None`/`Very Low` → `none`/`very-low`, `game_engine` tag now invalid (underscores)
- Updated error message: `lowercase letters, numbers, underscores` → `letters, numbers, and hyphens`
- Updated `FormatImportance` to use `Classification()` for Title Case display: `"50 (Medium)"` instead of `"50 (medium)"`
- 6 test packages pass (2 pre-existing failures in cliutil/storage unrelated)

### Phase 4: PMA (projects-management-automation) ✅
- **Commits:** `9593dc5c`, `a6b38049`, `80564a64`, `a17b8e02` (all pushed)
- Fixed SDK compilation breaks: adapted `sdkProjectToDomain` and `api/list.go` to use `programminglanguage.Languages.Strings()` and `.Language.Get()` instead of raw string access
- Changed `Importance int32` → `Importance uint8` (own type, not alias — can't add methods to external type)
- Updated all `int32` → `int` casts for importance range parameters across: filter.go, service.go, service_interfaces.go, metadata_registry.go, cmdguard_meta_sync.go, language_filter.go, project.go
- **Replaced `types.Tag` with CBT `tag.Tag` type alias** (the critical fix from post-migration audit)
  - Added `NormalizeTag()` for backward-compatible input normalization (lowercase + underscore→hyphen)
  - Updated `NewTag` signature from `(string) Tag` to `(string) (Tag, error)` — returns error on invalid input
  - Added `MustNewTag`, `ParseTags`, `FormatTags` helpers
- Simplified `meta/adapter.go` — Tag conversions are now identity operations since both types are `cbtTag.Tag`
- Fixed docs: `type Importance int32` → `type Importance uint8/int`
- All 26 test packages pass

### Phase 5: PDG (project-dependency-graph) ✅
- **Commit:** `0d9e1d5` (pushed)
- Added local replace directives for SDK, CBT, and go-branded-id
- PDG only uses `project.Path` from SDK — no Language/Tag/Importance types needed
- 1 test package passes

### Post-Migration Improvements ✅
- **Commit:** `e4a7b7a` — Added `Add`, `Sub`, `Clamp`, `Max`, `Min`, `IsDefault`, `IsNone`, `PercentString` to CBT Importance
- **Commit:** `bc6c6a3` — Added `Tags` slice type with `Strings()`, `IsEmpty()`, `Contains()` to CBT tag package
- **Commit:** `8e9de74` — Created `MIGRATION_COMPLETE.md` status document
- **Commit:** `43494c0` — Created `POST_MIGRATION_PLAN.md` improvement plan

---

## B) PARTIALLY DONE

### PMA Importance: Own type, not yet type alias
PMA's `Importance` is `uint8` (matching CBT's underlying type) but is NOT a type alias to `importance.Importance`. This means:
- Two separate `uint8` types with identical behavior exist
- The meta adapter still needs `ConvertImportance()` for the cast
- Can't become an alias until all PMA-specific methods (`Add`, `Sub`, `Clamp`, etc.) exist on CBT's type
- **Progress:** CBT now has all needed methods (`Add`, `Sub`, `Clamp`, `Max`, `Min`, `IsDefault`, `IsNone`). PMA could theoretically become a type alias now, but the `NewImportance(int)` and `NewImportanceFromString` constructors have different signatures than CBT's `New(uint8)` and `Parse(string)`.

### PMA Language enum: Not migrated (intentional)
PMA's `Language uint8` enum includes project-type categories that don't exist in the shared `programminglanguage.Language` branded string:
- `LanguageWeb`, `LanguageMobile`, `LanguageDesktop`, `LanguageCLI`, `LanguageLibrary`, `LanguageFramework`, `LanguageMonoRepo`, `LanguageMicroservice`
- These are **not programming languages** — they're project type classifications
- A proper fix would extract a `ProjectType` enum, but this is a larger refactor
- **Status:** Intentionally deferred. Not a split-brain — PMA's Language serves a different purpose than the shared type.

---

## C) NOT STARTED

| # | Task | Effort | Why |
|---|------|--------|-----|
| 1 | Extract `ProjectType` from PMA's Language enum | 2-4h | Requires touching 50+ switch statements across PMA |
| 2 | Auto-generate PMA Language string maps via `go-enum` or `stringer` | 30m | `stringToLanguage`/`languageStrings` are manually maintained inverse maps |
| 3 | Make PMA `Importance` a true type alias to CBT | 1h | Need to update `NewImportance(int)` callers to use `New(uint8)` + `Parse(string)` |
| 4 | Delete `ConvertImportance`/`ConvertImportanceToMeta` from adapter | 15m | Blocked on PMA Importance type alias |
| 5 | Root `go.work` for all 5 projects | 30m | Risk: could break builds of unrelated projects |
| 6 | CI workflow updates (3 projects) | 1-2h | PMA, project-meta, PDG CI pipelines need CBT SDK reference |
| 7 | `doctor --fix-tags` migration command in project-meta | 1h | Tag regex changed: underscores → hyphens. Existing YAML files may have underscores |
| 8 | Fix pre-existing project-meta `cliutil`/`storage` test failures | 30m | `TestResolveProjectNameFromPath` and `TestDetectProjectRoot` — unrelated to migration |

---

## D) TOTALLY FUCKED UP (Problems Found & Fixed)

### 🔴 D1: PMA `types.Tag` was never migrated (FIXED)
**What happened:** The original Phase 4 migration changed `Importance int32 → uint8` but completely missed `types.Tag`. PMA kept its own `Tag string` with `^[a-z0-9_-]+$` (allows underscores, forces lowercase) while project-meta used CBT's `tag.Tag` with `^[A-Za-z0-9-]+$` (no underscores, allows uppercase).

**Impact:** The `meta/adapter.go` `ConvertTag` function silently corrupted tags:
- `My-Tag` → `my-tag` (lowercased)
- Tags with underscores passed PMA validation but would fail CBT validation
- Comment in adapter.go claimed "identical validation rules" — **factually wrong**

**Fix:** `a6b38049` — Made `types.Tag` a type alias to `cbtTag.Tag`, added `NormalizeTag()` for backward compatibility, updated adapter to identity operation.

### 🟡 D2: CBT Importance was missing arithmetic methods
**What happened:** PMA needed `Add`, `Sub`, `Clamp`, `Max`, `Min`, `IsDefault`, `IsNone` that the shared type didn't have. This forced PMA to keep its own `Importance uint8` type instead of using a type alias.

**Impact:** Two identical `uint8` types with the same behavior. Conversion functions needed in adapter.

**Fix:** `e4a7b7a` — Added all missing methods to CBT Importance. PMA could now theoretically use a type alias.

### 🟡 D3: Adapter comment was factually wrong
**What happened:** `adapter.go:12` said `Both types are string aliases with identical validation rules (^[a-z0-9_-]+$)`. The two types actually had **opposite** rules (one allows uppercase/no underscores, other allows lowercase/underscores).

**Fix:** `80564a64` — Rewrote adapter. Since both types are now the same `cbtTag.Tag`, the conversion is an identity operation.

---

## E) WHAT WE SHOULD IMPROVE

### Architecture
1. **PMA Language enum mixes two concepts** — programming languages AND project types. These should be separate types. `LanguageGo` and `LanguageCLI` are fundamentally different categories.
2. **PMA Importance should become a type alias** to CBT once `NewImportance`/`NewImportanceFromString` are replaced with `New`/`Parse`. This eliminates the adapter entirely.
3. **Tag normalization strategy** — PMA auto-normalizes (lowercase + replace `_`), CBT validates strictly. There should be one canonical path: normalize at input boundary, validate everywhere else.

### Code Quality
4. **`meta/adapter.go` is nearly empty now** — most functions are identity operations. Consider deleting it entirely once PMA Importance becomes an alias.
5. **PMA's `stringsToLanguage` map** has aliases like `"golang" → LanguageGo`, `"js" → LanguageJavaScript`. This could be auto-generated from a single source of truth.
6. **Pre-existing test failures** in project-meta (`cliutil`/`storage`) should be fixed — they mask real bugs.

### Infrastructure
7. **No root `go.work`** — each project has its own workspace. A root workspace would enable cross-project builds but risks side effects.
8. **CI pipelines** don't reference local CBT — they'll break until CBT is published or CI uses go.work.
9. **`doctor --fix-tags`** is needed for the tag regex migration — existing YAML files with `game_engine` tags will fail validation.

---

## F) TOP 25 THINGS TO DO NEXT

Sorted by Impact × Value / Effort (highest first):

| # | Task | Project | Effort | Impact |
|---|------|---------|--------|--------|
| 1 | Make PMA `Importance` a type alias to CBT `importance.Importance` | PMA | 1h | HIGH |
| 2 | Delete `ConvertImportance`/`ConvertImportanceToMeta` from adapter | PMA | 15m | MED |
| 3 | Delete or gut `meta/adapter.go` (all conversions now identity/cast) | PMA | 30m | MED |
| 4 | Fix pre-existing project-meta `cliutil`/`storage` test failures | meta | 30m | MED |
| 5 | Implement `doctor --fix-tags` in project-meta | meta | 1h | HIGH |
| 6 | Add `go.work` to project-meta (if missing) | meta | 5m | LOW |
| 7 | Extract `ProjectType` from PMA Language enum | PMA | 2-4h | HIGH |
| 8 | Auto-generate PMA Language maps with `go-enum` or `stringer` | PMA | 30m | MED |
| 9 | Create root `go.work` for all 5 projects | root | 30m | MED |
| 10 | Update CI workflows for CBT dependency (3 projects) | CI | 1-2h | HIGH |
| 11 | Add `ParseTags` to CBT `tag` package (reusable everywhere) | CBT | 15m | MED |
| 12 | Replace PMA `NewImportance(int)` with CBT `New(uint8)` everywhere | PMA | 30m | MED |
| 13 | Replace PMA `NewImportanceFromString` with CBT `Parse(string)` | PMA | 30m | MED |
| 14 | Add `PercentString()` to CBT Importance (done), use it in PMA | PMA | 10m | LOW |
| 15 | Add integration test: round-trip Importance through all 3 projects | tests | 1h | MED |
| 16 | Add integration test: round-trip Tag through all 3 projects | tests | 1h | MED |
| 17 | Verify CBT `go:generate` works with new importance methods | CBT | 15m | LOW |
| 18 | Update `LIBRARY_GUIDE.md` with Tag/Importance method reference | CBT | 30m | LOW |
| 19 | Add fuzz tests for Tag validation edge cases (unicode, long strings) | CBT | 30m | LOW |
| 20 | Audit all YAML files in PMA testdata for underscore tags | PMA | 15m | MED |
| 21 | Add benchmark: Tag conversion before/after (was O(n) normalization, now identity) | PMA | 15m | LOW |
| 22 | Consider `ProjectType` as a shared type in CBT (for PMA/PDG) | CBT | 2h | MED |
| 23 | Update `CONTEXT.md` in PDG with new dependency info | PDG | 10m | LOW |
| 24 | Clean up `PUBLIC_OR_PRIVATE.md` files in SDK and project-meta | meta/sdk | 5m | LOW |
| 25 | Verify `GOEXPERIMENT=jsonv2` is set in all flake.nix devShells | all | 15m | MED |

---

## G) TOP #1 QUESTION

**Can CBT (`go-composable-business-types`) be published to a private Go module proxy, or should all consumer projects permanently use local `replace` directives?**

This matters because:
- All 4 consumer projects have `replace github.com/larsartmann/go-composable-business-types => /home/lars/projects/go-composable-business-types` in their `go.mod`
- CI pipelines will fail without either a published version or workspace awareness
- The `go.work` approach in PMA works locally but isn't portable
- If there's already a private GOPROXY or GitHub Package Registry set up, I should use `go get` with a version tag instead of local replace

---

## Test Results Summary

| Project | Packages | Pass | Fail | Pre-existing |
|---------|----------|------|------|-------------|
| go-composable-business-types | 17 | 17 | 0 | 0 |
| project-discovery-sdk | 4 | 4 | 0 | 0 |
| project-meta | 8 | 6 | 2 | 2 (cliutil, storage) |
| projects-management-automation | 26 | 26 | 0 | 0 |
| project-dependency-graph | 1 | 1 | 0 | 0 |
| **TOTAL** | **56** | **54** | **2** | **2** |

## Type Alignment Matrix

| Type | CBT | SDK | project-meta | PMA | PDG |
|------|-----|-----|-------------|-----|-----|
| Importance | `uint8` source of truth | `importance.Importance` | `= importance.Importance` (alias) | `uint8` own type | N/A |
| Tag | `string`, `^[A-Za-z0-9-]+$` | `tag.Tag` | `= cbtTag.Tag` (alias) | `= cbtTag.Tag` (alias) | N/A |
| Language | `branded string` | `programminglanguage.Language` | N/A | `uint8` own enum (intentional) | N/A |
| Languages | `[]Language` + `.Primary()` | `programminglanguage.Languages` | N/A | N/A | N/A |

## Git Log (All Projects)

### go-composable-business-types (4 commits this session)
```
43494c0 docs: Add post-migration improvement plan
bc6c6a3 feat(tag): add Tags slice type with Strings, IsEmpty, Contains helpers
e4a7b7a feat(importance): add arithmetic, comparison, and helper methods
8e9de74 docs: Add migration completion status document
383dd60 Phase 0 complete: ADR, flake.nix, fix enums tests
```

### project-discovery-sdk (1 commit this session)
```
8b44eb8 Phase 2: Wire shared types into project-discovery-sdk
```

### project-meta (2 commits this session)
```
da188ac refactor: Use Classification() for human-readable Importance display
bc51367 Phase 3: Wire shared types into project-meta domain
```

### projects-management-automation (4 commits this session)
```
a17b8e02 docs: Update Importance type references from int32 to uint8/int
80564a64 refactor: Simplify meta adapter — Tag conversions are now identity ops
a6b38049 feat: Replace PMA types.Tag with shared CBT tag.Tag type alias
9593dc5c Phase 4: Wire shared types into PMA
```

### project-dependency-graph (1 commit this session)
```
0d9e1d5 Phase 5: Wire local SDK + CBT dependencies into PDG
```

---

_Report generated by Crush AI Assistant_
