# MASTER EXECUTION PLAN — Shared Types Ecosystem Migration

**Created:** 2026-05-04 11:01  
**Total tasks:** 68  
**Max task duration:** 12 minutes  
**Sorted by:** Impact × Customer Value / Effort (highest first)  
**Legend:** ✅ = Done | ⬜ = Not started | 🔶 = Partial

---

## TASK TABLE

| # | Status | Task | Project | Impact | Effort | Phase |
|---|---|---|---|---|---|---|
| **PHASE 0: PRE-REQUISITES** |||||||
| 1 | ✅ | Create `importance/` package with type, constructors, classification, JSON, SQL, tests | cbt | HIGH | 12m | 0 |
| 2 | ✅ | Create `tag/` package with regex validation `^[A-Za-z0-9-]+$`, JSON, SQL, tests | cbt | HIGH | 12m | 0 |
| 3 | ✅ | Create `programminglanguage/` package: Language branded type, Languages slice, Normalize, tests | cbt | HIGH | 12m | 0 |
| 4 | ✅ | Create `projectcore/` package: ProjectCore struct, options, Validate, JSON, tests | cbt | MED | 12m | 0 |
| 5 | ✅ | Add sentinel errors for importance/tag/projectcore to `pkg/errors/` | cbt | MED | 5m | 0 |
| 6 | ⬜ | Write ADR-001: Decision to add project ecosystem types to go-composable-business-types | cbt | MED | 12m | 0 |
| 7 | ⬜ | Create `flake.nix` for go-composable-business-types (nix build + test) | cbt | MED | 12m | 0 |
| 8 | ⬜ | Fix pre-existing `enums/` test failure in go-composable-business-types | cbt | LOW | 12m | 0 |
| **PHASE 2: WIRE INTO project-discovery-sdk** |||||||
| 9 | ⬜ | Add `go-composable-business-types` to SDK go.mod with local replace directive | sdk | MED | 5m | 2 |
| 10 | ⬜ | Replace `sdk.Project.Language string` with `programminglanguage.Language` in types.go | sdk | HIGH | 12m | 2 |
| 11 | ⬜ | Replace `sdk.Project.Languages []string` with `programminglanguage.Languages` in types.go | sdk | HIGH | 10m | 2 |
| 12 | ⬜ | Replace `sdk.Project.Importance int` with `importance.Importance` in types.go | sdk | HIGH | 10m | 2 |
| 13 | ⬜ | Replace `sdk.Project.Tags []string` with `[]tag.Tag` in types.go | sdk | HIGH | 10m | 2 |
| 14 | ⬜ | Update `sdk.Project.Equals()` to work with new types (Language, Importance comparison) | sdk | HIGH | 12m | 2 |
| 15 | ⬜ | Update `sdk.Project.changedFields()` snapshot diffing for new types | sdk | HIGH | 12m | 2 |
| 16 | ⬜ | Replace SDK `detection/NormalizeLanguage` with call to `programminglanguage.Normalize` | sdk | MED | 12m | 2 |
| 17 | ⬜ | Remove `enryLanguageAliases` map from `detection/registry.go` (now in programminglanguage) | sdk | MED | 8m | 2 |
| 18 | ⬜ | Delete `Language` branded type from SDK `branded.go` → use `programminglanguage.Language` | sdk | MED | 10m | 2 |
| 19 | ⬜ | Delete `ProjectName`/`ProjectPath` branded types from SDK `branded.go` (keep in project-meta only) | sdk | MED | 8m | 2 |
| 20 | ⬜ | Update `BrandedProject` struct + `AsBranded()`/`AsProject()` converters for new types | sdk | MED | 12m | 2 |
| 21 | ⬜ | Update `FilterByLanguage` to compare `programminglanguage.Language` instead of strings | sdk | HIGH | 10m | 2 |
| 22 | ⬜ | Update `FilterByImportance` to compare `importance.Importance` instead of raw int | sdk | MED | 8m | 2 |
| 23 | ⬜ | Update `FilterByTags` to compare `[]tag.Tag` instead of `[]string` | sdk | MED | 8m | 2 |
| 24 | ⬜ | Update `FilterByNameGlob` and remaining filters for type compatibility | sdk | MED | 10m | 2 |
| 25 | ⬜ | Update `Event` struct: ensure `*Project` embed works with new field types | sdk | HIGH | 8m | 2 |
| 26 | ⬜ | Update all hook/event call sites that access `e.Project.Language`, `e.Project.Tags` etc. | sdk | MED | 12m | 2 |
| 27 | ⬜ | Update `sdk.go`: `collectLanguages()`, `defaultOptions()` for new Language type | sdk | MED | 12m | 2 |
| 28 | ⬜ | Update `client.go`: `WatchObservable()` snapshot diffing for new types | sdk | MED | 12m | 2 |
| 29 | ⬜ | Update `options.go`: `DiscoverOptions.Languages []string` → keep as string (backward compat) | sdk | MED | 8m | 2 |
| 30 | ⬜ | Update SDK `sdk_test.go`: fix type assertions for Language/Importance/Tags | sdk | MED | 12m | 2 |
| 31 | ⬜ | Update SDK `sdk_extra_test.go`: fix type assertions | sdk | MED | 10m | 2 |
| 32 | ⬜ | Update SDK `filter_test.go`: fix filter tests for new types | sdk | MED | 12m | 2 |
| 33 | ⬜ | Update SDK `identifier_test.go`, `snapshot_test.go`, `events_test.go` for new types | sdk | MED | 12m | 2 |
| 34 | ⬜ | Update SDK `detection/enry_test.go`, `detection/registry_test.go` for Normalize change | sdk | MED | 10m | 2 |
| 35 | ⬜ | Update SDK `discovery_bdd_test.go` for new Language/Importance/Tag types | sdk | MED | 12m | 2 |
| 36 | ⬜ | Run full SDK test suite, fix remaining failures, verify green | sdk | MED | 12m | 2 |
| 37 | ⬜ | SDK: run `go mod tidy`, verify build, commit + push Phase 2 | sdk | HIGH | 5m | 2 |
| **PHASE 3: WIRE INTO project-meta** |||||||
| 38 | ⬜ | Add `go-composable-business-types` to project-meta go.mod with local replace | meta | MED | 5m | 3 |
| 39 | ⬜ | Replace `Importance int32` type definition with `importance.Importance` in `domain/metadata.go` | meta | HIGH | 12m | 3 |
| 40 | ⬜ | Update `NewImportance(int)` → `importance.New(uint8)`, update all callers | meta | HIGH | 12m | 3 |
| 41 | ⬜ | Update `NewImportanceFromString` → delegate to `importance.Parse` | meta | MED | 10m | 3 |
| 42 | ⬜ | Update `Metadata.SetImportance` and `GetImportance` for new type | meta | MED | 10m | 3 |
| 43 | ⬜ | Replace `Tag string` type definition with `tag.Tag` in `domain/metadata.go` | meta | HIGH | 12m | 3 |
| 44 | ⬜ | Update `NewTag()`, `Tag.Validate()` → delegate to `tag.New()`, `tag.Validate()` | meta | HIGH | 12m | 3 |
| 45 | ⬜ | Update `Metadata.AddTags/RemoveTags/SetTags/HasTag` for `[]tag.Tag` | meta | MED | 12m | 3 |
| 46 | ⬜ | Add backward-compat YAML/JSON unmarshaling: accept both old int32 and new uint8 Importance | meta | MED | 12m | 3 |
| 47 | ⬜ | Update `metadata_file.go`: ReadByPath/WriteByPath for new types | meta | MED | 10m | 3 |
| 48 | ⬜ | Update `api/metadata.go`: all API methods returning Importance/Tag | meta | MED | 12m | 3 |
| 49 | ⬜ | Update `manager/manager.go`: filter/per-project methods for new types | meta | MED | 12m | 3 |
| 50 | ⬜ | Update all project-meta tests: fix type assertions for Importance uint8, Tag regex | meta | MED | 12m | 3 |
| 51 | ⬜ | Implement `doctor --fix-tags` command: `strings.ReplaceAll(tag, "_", "-")` | meta | MED | 12m | 3 |
| 52 | ⬜ | Update CLI commands (`cmd_list.go`, `cmd_set.go`, etc.) for new types | meta | MED | 12m | 3 |
| 53 | ⬜ | Run full project-meta test suite, verify green | meta | MED | 12m | 3 |
| 54 | ⬜ | Meta: run `go mod tidy`, verify build, commit + push Phase 3 | meta | HIGH | 5m | 3 |
| **PHASE 4: WIRE INTO projects-management-automation** |||||||
| 55 | ⬜ | Add `go-composable-business-types` to PMA go.mod + go.work | pma | MED | 5m | 4 |
| 56 | ⬜ | Replace `domain.Importance int32` with `importance.Importance uint8` in `domain/importance.go` | pma | HIGH | 12m | 4 |
| 57 | ⬜ | Update `domain/project.go`: Project struct Importance field + constructors | pma | HIGH | 12m | 4 |
| 58 | ⬜ | Fix Importance type leaks: `api/list.go` (raw int), `cmdguard_meta_sync.go` (raw int32) | pma | HIGH | 12m | 4 |
| 59 | ⬜ | Fix Importance type leaks: `metadata_registry.go` (int32 yaml), `cli_output.go` (int) | pma | MED | 10m | 4 |
| 60 | ⬜ | Replace `types.Tag` with `tag.Tag` in `types/strong_types.go` | pma | HIGH | 12m | 4 |
| 61 | ⬜ | Update all Tag references across PMA (search for `types.Tag`, `types.NewTag`, `types.TagsToStrings`) | pma | HIGH | 12m | 4 |
| 62 | ⬜ | DELETE `internal/infrastructure/meta/adapter.go` — no type conversion needed | pma | HIGH | 5m | 4 |
| 63 | ⬜ | Update `meta/service.go`: remove adapter dependencies, use shared types directly | pma | MED | 12m | 4 |
| 64 | ⬜ | Replace `types.Language uint8` enum with `programminglanguage.Language` branded string | pma | HIGH | 12m | 4 |
| 65 | ⬜ | Convert `types/language.go` switch statements to dispatch maps (IsCompiled, IsInterpreted) | pma | MED | 12m | 4 |
| 66 | ⬜ | Convert `dependency_manager.go` switch to `map[Language]depCommand` dispatch | pma | MED | 12m | 4 |
| 67 | ⬜ | Convert `filesystem_project_repository.go` switch to dispatch map | pma | MED | 10m | 4 |
| 68 | ⬜ | Extract `LanguageWeb`/`LanguageMobile`/`LanguageCLI` etc. → `ProjectType` in PMA | pma | MED | 12m | 4 |
| 69 | ⬜ | Update `sdk_discoverer.go`: remove `sdkProjectToDomain` conversion (same types now) | pma | HIGH | 12m | 4 |
| 70 | ⬜ | Deduplicate `sdkProjectToDomain` in `api/list.go` (delete second copy, fix case bug) | pma | HIGH | 8m | 4 |
| 71 | ⬜ | Update all PMA BDD tests: fix Language/Tag/Importance type references | pma | MED | 12m | 4 |
| 72 | ⬜ | Update PMA unit tests: fix Language enum references → branded string | pma | MED | 12m | 4 |
| 73 | ⬜ | Update PMA test utilities: `project_helpers.go`, `test_factory.go` for new types | pma | MED | 12m | 4 |
| 74 | ⬜ | Run full PMA test suite, fix failures, verify green | pma | MED | 12m | 4 |
| 75 | ⬜ | PMA: run `go mod tidy`, verify build, commit + push Phase 4 | pma | HIGH | 5m | 4 |
| **PHASE 5: WIRE INTO project-dependency-graph** |||||||
| 76 | ⬜ | Add `go-composable-business-types` to PDG go.mod | pdg | MED | 5m | 5 |
| 77 | ⬜ | Replace hardcoded `[]string{"go"}` with `programminglanguage.New("go")` in `discover.go` | pdg | MED | 8m | 5 |
| 78 | ⬜ | Run full PDG test suite, verify green | pdg | MED | 10m | 5 |
| 79 | ⬜ | PDG: commit + push Phase 5 | pdg | MED | 5m | 5 |
| **CROSS-CUTTING** |||||||
| 80 | ⬜ | Add `go-composable-business-types` to PMA's `go.work` (cross-project compile safety) | pma | MED | 5m | X |
| 81 | ⬜ | Unify `go-filewatcher` version across SDK (v0.1.0) and PMA (v0.2.0) | all | LOW | 12m | X |
| 82 | ⬜ | Add basic CI workflow to project-discovery-sdk (go test ./...) | sdk | MED | 12m | X |
| 83 | ⬜ | Add basic CI workflow to project-meta (go test ./...) | meta | MED | 12m | X |
| 84 | ⬜ | Add basic CI workflow to project-dependency-graph (go test ./...) | pdg | MED | 12m | X |
| 85 | ⬜ | Update `LIBRARY_GUIDE.md` with final state of all new types | docs | LOW | 10m | X |

---

## PROGRESS SUMMARY

| Category | Count | Details |
|---|---|---|
| ✅ Fully Done | 5 | Tasks 1-5: importance, tag, programminglanguage, projectcore, errors |
| ⬜ Not Started | 80 | Tasks 6-85 |
| 🔶 Partially Done | 0 | — |
| 💥 Fucked Up | 0 | — |
| **TOTAL** | **85** | |

## PHASE COMPLETION

| Phase | Tasks | Done | % | Status |
|---|---|---|---|---|
| Phase 0: Pre-requisites | 8 | 5 | 63% | Missing: ADR, flake.nix, enums fix |
| Phase 2: project-discovery-sdk | 29 | 0 | 0% | Not started |
| Phase 3: project-meta | 17 | 0 | 0% | Not started |
| Phase 4: PMA | 21 | 0 | 0% | Not started |
| Phase 5: PDG | 4 | 0 | 0% | Not started |
| Cross-cutting | 6 | 0 | 0% | Not started |
| **TOTAL** | **85** | **5** | **6%** | |

## CRITICAL PATH

```
Phase 0 (remaining) → Phase 2 (SDK) → Phase 3 (meta) → Phase 4 (PMA) → Phase 5 (PDG) → Cross-cutting
                        ↑               ↑                 ↑
                        └─── unblocks ──┘─── unblocks ───┘
```

Phase 2 MUST complete before Phase 3 or 4 can start (both consume SDK's types).
Phase 3 and Phase 4 can partially overlap but PMA's `meta/adapter.go` depends on Phase 3 completing first.

## ESTIMATED REMAINING EFFORT

| Phase | Tasks Remaining | Est. Time |
|---|---|---|
| Phase 0 | 3 tasks | ~36 min |
| Phase 2 | 29 tasks | ~5.5 hrs |
| Phase 3 | 17 tasks | ~3.2 hrs |
| Phase 4 | 21 tasks | ~4 hrs |
| Phase 5 | 4 tasks | ~45 min |
| Cross-cutting | 6 tasks | ~1.2 hrs |
| **TOTAL** | **80 tasks** | **~15 hrs** |

---

_Generated with Crush — 2026-05-04_
