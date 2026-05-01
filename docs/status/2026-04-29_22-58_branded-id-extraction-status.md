# Status Report: go-branded-id Extraction

**Date:** 2026-04-29_22-58
**Session:** Extract `id/` package into standalone `go-branded-id` library

---

## A) FULLY DONE

1. **Created `/home/lars/projects/go-branded-id/`** — new standalone Go module
   - Module: `github.com/larsartmann/go-branded-id`
   - Go 1.26, zero external dependencies
   - All 5 source files (`id.go`, `id_json.go`, `id_binary.go`, `id_text.go`, `id_sql.go`)
   - All 5 test files + assert helper + bench/fuzz/example test
   - `.gitignore`, `CHANGELOG.md`
   - Two commits already on disk (initial + git-town config)
   - `go test -race ./...` — **PASS**
   - `go vet ./...` — **CLEAN**

2. **Updated `go-composable-business-types` imports** in 4 consumer files:
   - `actor/actor.go` — `id` import → `github.com/larsartmann/go-branded-id`
   - `actor/actor_test.go` — same
   - `datapoint/datapoint_test.go` — same
   - `examples/datapoint/main.go` — same

3. **Added `replace` directive** in `go.mod` for local development
   - `replace github.com/larsartmann/go-branded-id => ../go-branded-id`

4. **Updated `AGENTS.md`** — dependencies section, package structure, import examples, notes

5. **Verified tests pass** in both repos after extraction
   - `actor/`, `datapoint/` all green
   - `enums/` failures are **pre-existing** (unrelated to extraction)

---

## B) PARTIALLY DONE

1. **`id/` directory NOT removed from `go-composable-business-types`**
   - `git rm -rf id/` was run but the deletion did NOT persist — files still on disk and tracked
   - The old `id/` package coexists with the new import paths, creating **ambiguity**
   - **MUST be completed before committing**

2. **`id/id.go` has local modifications** (pre-existing, not from this session)
   - Sentinel error variables added (`errStringAssert`, `errUnsupportedBinary`, etc.)
   - Parameter rename `v` → `value` in `NewID`
   - These changes exist in the old project but were **NOT copied to `go-branded-id`**
   - The new repo has the older version without these improvements

3. **`go mod tidy` ran once** but may need re-run after `id/` removal

---

## C) NOT STARTED

1. **LICENSE file** for `go-branded-id` — missing entirely
2. **README.md** for `go-branded-id` — missing (the old `id/README.md` was not migrated)
3. **CI/CD** for `go-branded-id` — no GitHub Actions, no release workflow
4. **Git tag / version** for `go-branded-id` — no v0.1.0 tag
5. **Sync sentinel errors** from old `id/id.go` to new repo
6. **`.golangci.yml`** for `go-branded-id` — no lint config
7. **`go mod tidy`** in old project after removing `id/` directory
8. **Update `docs/planning/go-branded-id-extraction-analysis.md`** — mark as completed
9. **Remove `id/README.md`** from old project
10. **Verify `go-sum` is clean** after `go mod tidy`

---

## D) TOTALLY FUCKED UP

1. **`git rm` did not stick** — The `git rm -rf id/` ran successfully (saw output) but the files are still present and tracked. Likely the shell session didn't persist state between tool calls, or `git mod tidy` re-added awareness. **Root cause: git rm needs to be part of a single committed transaction with the go.mod/go.sum changes.**

2. **Drift between old and new repo** — The old `id/id.go` has ~89 lines of changes (sentinel errors, parameter renames) that were NOT in the new repo. We copied from the **committed** version, not the **working** version.

---

## E) WHAT WE SHOULD IMPROVE

### Architecture

1. **The extraction plan (2026-04-05) recommended extracting `id/` + `nanoid/` + `scanutil/` + `pkg/errors/`** — we only did `id/`. NanoID is a natural companion since `ID[Brand, NanoID]` is a common pattern.
2. **No type-level composition between `id.ID` and `nanoid.NanoID`** — `NanoID` is a standalone struct `{ value string }` that doesn't use `id.ID[B,V]`. We could make `NanoID` a valid value type for `ID` more ergonomically.
3. **The `id/` package name in the standalone lib is fine** — `id.NewID[UserBrand]("user-123")` is clean. No rename needed.

### Process

4. **Should have committed the `git rm` immediately** in the same shell session as the import updates
5. **Should have diff-checked** the old working tree against the new repo files to catch the sentinel error drift
6. **Should have added LICENSE before any commit** — proprietary license from parent project needs to be replicated

### Code Quality

7. **The sentinel error refactoring in old `id/id.go` is incomplete** — only `id.go` was modified, but `id_binary.go`, `id_sql.go`, `id_text.go` still use `fmt.Errorf` for dynamic errors. The new repo doesn't have any of this.
8. **Linter warnings** in the id package (varnamelen, funlen, err113, etc.) — should be addressed in the new clean repo

---

## F) TOP 25 THINGS WE SHOULD GET DONE NEXT

Sorted by **Impact × Ease** (highest first):

| #   | Task                                                                                        | Impact   | Effort | Why                                    |
| --- | ------------------------------------------------------------------------------------------- | -------- | ------ | -------------------------------------- |
| 1   | Remove `id/` from old project (git rm + commit)                                             | Critical | 5min   | Can't ship with duplicate packages     |
| 2   | Sync sentinel error changes from old → new repo                                             | High     | 15min  | Drift will cause confusion             |
| 3   | Add LICENSE to go-branded-id                                                                | High     | 2min   | Can't publish without it               |
| 4   | Migrate id/README.md → go-branded-id README.md                                              | High     | 10min  | pkg.go.dev needs it                    |
| 5   | Run `go mod tidy` in old project after id/ removal                                          | Critical | 2min   | Module won't build without it          |
| 6   | Commit all changes in old project                                                           | Critical | 5min   | Uncommitted work = lost work           |
| 7   | Commit synced changes in go-branded-id                                                      | High     | 5min   | Get both repos in clean state          |
| 8   | Add `.golangci.yml` to go-branded-id                                                        | Medium   | 10min  | Lint consistency                       |
| 9   | Fix linter warnings in go-branded-id (varnamelen, funlen, err113)                           | Medium   | 30min  | Clean slate, clean code                |
| 10  | Add GitHub Actions CI for go-branded-id                                                     | Medium   | 20min  | Automated testing                      |
| 11  | Tag go-branded-id as v0.1.0                                                                 | Medium   | 2min   | Versioned dependency                   |
| 12  | Push both repos to GitHub                                                                   | High     | 5min   | Backup + collaboration                 |
| 13  | Update extraction analysis doc as completed                                                 | Low      | 5min   | Documentation hygiene                  |
| 14  | Add `justfile` to go-branded-id                                                             | Low      | 10min  | Consistent build commands              |
| 15  | Consider extracting `scanutil/` into go-branded-id                                          | Medium   | 30min  | Natural dependency for ID SQL scanning |
| 16  | Consider extracting `pkg/errors/` into go-branded-id                                        | Medium   | 30min  | Shared error infrastructure            |
| 17  | Consider extracting `nanoid/` into go-branded-id                                            | Medium   | 1hr    | Natural companion: `ID[Brand, NanoID]` |
| 18  | Add `encoding.TextMarshaler`/`TextUnmarshaler` interface assertions for more types          | Low      | 10min  | Compile-time safety                    |
| 19  | Fix pre-existing enums test failures                                                        | Medium   | 1hr    | Unrelated but embarrassing             |
| 20  | Add Go doc examples to go-branded-id README                                                 | Low      | 15min  | Adoption                               |
| 21  | Add dependabot config for go-branded-id                                                     | Low      | 5min   | Security                               |
| 22  | Consider `ID[B,V]` implementing `encoding.BinaryMarshaler` for custom V types via interface | Low      | 20min  | Extensibility                          |
| 23  | Add release workflow (goreleaser or similar) to go-branded-id                               | Low      | 30min  | Automated releases                     |
| 24  | Update PROJECT_SPLIT_EXECUTIVE_REPORT.md to reflect partial extraction done                 | Low      | 10min  | Accuracy                               |
| 25  | Evaluate if `bounded.BoundedString` could be a valid ID value type                          | Low      | 15min  | Type composition exploration           |

---

## G) TOP #1 QUESTION

**Should we follow the original extraction plan and also move `nanoid/`, `scanutil/`, and `pkg/errors/` into `go-branded-id` in this session, or keep `go-branded-id` focused on just `ID[B,V]`?**

Context:

- The 2026-04-05 plan recommended the fuller extraction
- `nanoid/` depends on `scanutil/` and `pkg/errors/` — these would need to come along
- `datapoint/` uses `nanoid.NanoID` as its ID type — importing from `go-branded-id/nanoid` vs keeping it local
- User explicitly said "not a fan of re-exports" — so we'd update imports directly
- This is a significantly larger scope but makes `go-branded-id` a complete "identity" library

---

## Current Git State

### go-composable-business-types

- Branch: `master`
- 7 modified files, **NOT staged**, **NOT committed**
- `id/` directory still exists on disk and in git
- Tests pass (except pre-existing enums failures)

### go-branded-id

- Branch: `master`
- 2 commits on disk
- Working tree clean
- Missing: LICENSE, README, CI
- Tests pass, vet clean
