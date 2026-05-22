# Execution Plan — Module Split

**Ordered by Pareto impact.** Each step is independently committable and leaves the project in a buildable, testable state.

---

## Step 1: Create go.work (1% → 51% impact)

**What:** Create `go.work` at repo root. Initially just points to root module.

**Why:** Foundation for all subsequent steps. Without go.work, sub-modules can't resolve root imports locally.

**Actions:**
- Create `go.work` with `use .`
- Remove `GOWORK=off` from flake.nix

**Verify:** `go work sync && go build ./... && go test ./...`

**Rollback:** Delete go.work, restore flake.nix GOWORK setting.

---

## Step 2: Extract nanoid module (4% → 55% impact)

**What:** Create `nanoid/go.mod`, making nanoid a separate module.

**Why:** Isolates sixafter/nanoid + crypto deps. Most impactful isolation because nanoid is widely used but has specific crypto deps.

**Actions:**
- Create `nanoid/go.mod`:
  ```
  module github.com/larsartmann/go-composable-business-types/nanoid
  go 1.26.2
  require github.com/larsartmann/go-composable-business-types v0.0.0
  require github.com/sixafter/nanoid v1.64.2
  ```
- Add `./nanoid` to `go.work` use list
- Run `go work sync`
- Remove nanoid, sixafter/nanoid, and crypto deps from root go.mod
- Run `go mod tidy` in both root and nanoid

**Verify:**
- `cd nanoid && go build ./... && go test ./...`
- `go build ./... && go test ./...` (from root with go.work)
- `go mod tidy` changes nothing in both modules

**Rollback:** Delete nanoid/go.mod, restore root go.mod, remove nanoid from go.work.

---

## Step 3: Extract locale module (4% → 59% impact)

**What:** Create `locale/go.mod`, making locale a separate module.

**Why:** Isolates golang.org/x/text.

**Actions:**
- Create `locale/go.mod`:
  ```
  module github.com/larsartmann/go-composable-business-types/locale
  go 1.26.2
  require github.com/larsartmann/go-composable-business-types v0.0.0
  require golang.org/x/text v0.37.0
  ```
- Add `./locale` to `go.work`
- Run `go work sync`
- Remove locale and x/text from root go.mod
- Run `go mod tidy` in root, locale

**Verify:**
- `cd locale && go build ./... && go test ./...`
- Root build + test still pass

**Rollback:** Delete locale/go.mod, restore root go.mod.

---

## Step 4: Extract money module (4% → 63% impact)

**What:** Create `money/go.mod`, making money a separate module.

**Why:** Isolates bojanz/currency (heaviest dep: apd, decimal, pq). Money depends on locale module.

**Actions:**
- Create `money/go.mod`:
  ```
  module github.com/larsartmann/go-composable-business-types/money
  go 1.26.2
  require github.com/larsartmann/go-composable-business-types v0.0.0
  require github.com/larsartmann/go-composable-business-types/locale v0.0.0
  require github.com/bojanz/currency v1.4.4
  ```
- Add `./money` to `go.work`
- Run `go work sync`
- Remove money, bojanz/currency and transitive deps from root go.mod
- Run `go mod tidy` in root, money

**Verify:**
- `cd money && go build ./... && go test ./...`
- Root build + test still pass

**Rollback:** Delete money/go.mod, restore root go.mod.

---

## Step 5: Extract datapoint module (4% → 67% impact)

**What:** Create `datapoint/go.mod`, making datapoint a separate module.

**Why:** Datapoint depends on nanoid module — can't stay in root. Isolates the "composite" type.

**Actions:**
- Create `datapoint/go.mod`:
  ```
  module github.com/larsartmann/go-composable-business-types/datapoint
  go 1.26.2
  require github.com/larsartmann/go-composable-business-types v0.0.0
  require github.com/larsartmann/go-composable-business-types/nanoid v0.0.0
  ```
- Add `./datapoint` to `go.work`
- Run `go work sync`
- Run `go mod tidy` in root, datapoint

**Verify:**
- `cd datapoint && go build ./... && go test ./...`
- Root build + test still pass
- Datapoint doesn't pull in bojanz/currency

**Rollback:** Delete datapoint/go.mod, restore root go.mod.

---

## Step 6: Extract examples module (20% → 80% impact)

**What:** Create `examples/go.mod`.

**Why:** Examples depend on multiple modules. Keeps example deps out of library modules.

**Actions:**
- Create `examples/go.mod`:
  ```
  module github.com/larsartmann/go-composable-business-types/examples
  go 1.26.2
  require github.com/larsartmann/go-composable-business-types v0.0.0
  require github.com/larsartmann/go-composable-business-types/nanoid v0.0.0
  require github.com/larsartmann/go-composable-business-types/datapoint v0.0.0
  ```
- Add `./examples` to `go.work`
- Run `go work sync`
- Run `go mod tidy` in examples

**Verify:**
- `cd examples && go build ./...`
- Root build + test still pass

**Rollback:** Delete examples/go.mod, restore root go.mod.

---

## Step 7: Update flake.nix (20% → 87% impact)

**What:** Update flake.nix to build and test each module independently.

**Why:** CI must verify all modules, not just root.

**Actions:**
- Update build check to build each module: `go build ./...` in root, nanoid, locale, money, datapoint, examples
- Update test check to test each module
- Remove GOWORK=off, ensure GOEXPERIMENT=jsonv2 is set

**Verify:** `nix build` and `nix flake check` pass

**Rollback:** Restore previous flake.nix.

---

## Step 8: Update CI/CD (20% → 93% impact)

**What:** Update GitHub Actions for multi-module.

**Why:** CI must test all modules, not just root.

**Actions:**
- Update test job to iterate over modules
- Consider parallel per-module jobs
- Update release workflow for multi-module tagging

**Verify:** Push to branch, verify CI passes

**Rollback:** Restore previous CI config.

---

## Step 9: Update documentation (remaining → 97% impact)

**What:** Update README.md, AGENTS.md, architecture.d2 to reflect new structure.

**Why:** Documentation must reflect reality.

**Actions:**
- Update README with module boundaries
- Update AGENTS.md with per-module build/test commands
- Update architecture.d2 with module boundaries
- Update justfile release command

**Verify:** Review documentation accuracy

**Rollback:** Restore previous documentation.

---

## Step 10: Final verification (remaining → 100%)

**What:** Full end-to-end verification of the modularized project.

**Actions:**
- `go work sync`
- Build and test every module independently
- Build and test from root with go.work
- Verify `go mod tidy` is clean in all modules
- Verify `go vet ./...` passes everywhere
- Run lint

**Verify:** Everything passes clean.

**Rollback:** N/A (final state)

---

## Dependency Map Between Steps

```
Step 1 (go.work) ← prerequisite for all
  ├── Step 2 (nanoid) ← no prerequisite beyond Step 1
  ├── Step 3 (locale) ← no prerequisite beyond Step 1
  ├── Step 4 (money) ← requires Step 3 (locale module exists)
  ├── Step 5 (datapoint) ← requires Step 2 (nanoid module exists)
  └── Step 6 (examples) ← requires Steps 2, 5 (nanoid + datapoint exist)
Step 7 (flake.nix) ← requires Steps 2-6
Step 8 (CI) ← requires Step 7
Step 9 (docs) ← requires Steps 2-6
Step 10 (final) ← requires all
```
