# Shared Types Migration — Completion Status

**Updated:** 2026-05-04  
**Status:** All phases complete ✅

---

## Summary

Successfully unified split-brain type systems across four Go projects by extending `go-composable-business-types` (CBT) with shared `Importance`, `Tag`, and `Language` types.

## Phases Completed

### Phase 0: CBT Pre-requisites ✅
- Created `importance/`, `tag/`, `programminglanguage/`, `projectcore/` packages
- Added ADR-001, flake.nix, fixed pre-existing enums test
- **Commit:** `383dd60`

### Phase 2: project-discovery-sdk ✅
- Replaced `Language string` → `programminglanguage.Language`
- Replaced `Languages []string` → `programminglanguage.Languages`
- Replaced `Importance int` → `importance.Importance`
- Replaced `Tags []string` → `[]tag.Tag`
- Removed SDK's `Language` branded type
- Updated all filters, events, tests
- **Commit:** `8b44eb8`

### Phase 3: project-meta ✅
- `Importance` now type alias to `importance.Importance` (uint8)
- `Tag` now type alias to `cbtTag.Tag`
- Updated error messages and test expectations
- Pre-existing failures in cliutil/storage (unrelated)
- **Commit:** `bc51367`

### Phase 4: PMA ✅
- `Importance` changed from `int32` to `uint8` (0-100)
- Fixed SDK compilation breaks (adapted to new Language/Languages types)
- Updated all `int32` → `int` casts for importance parameters
- All 27 test packages pass
- **Commit:** `9593dc5c`

### Phase 5: PDG ✅
- Added local replace directives for SDK + CBT + go-branded-id
- PDG only uses `project.Path` from SDK — no type changes needed
- All tests pass
- **Commit:** `0d9e1d5`

## Key Decisions

| Decision | Rationale |
|----------|-----------|
| PMA `Language` kept as own `uint8` enum | Includes project-type categories (Web, CLI, Library) not in shared type |
| PMA `Importance` changed to `uint8` but not type alias | Can't add methods to external type alias in Go |
| Tag regex: `^[A-Za-z0-9-]+$` | Uppercase allowed, underscores dropped |
| `DiscoverOptions.Languages` stays `[]string` | Backward compatible input filter |
| Local replace directives | Private repos, no GOPROXY access |

## Verification

| Project | Build | Tests |
|---------|-------|-------|
| go-composable-business-types | ✅ | ✅ (GOEXPERIMENT=jsonv2) |
| project-discovery-sdk | ✅ | ✅ |
| project-meta | ✅ | ✅ (2 pre-existing failures in cliutil/storage) |
| projects-management-automation | ✅ | ✅ (all 27 packages) |
| project-dependency-graph | ✅ | ✅ |
