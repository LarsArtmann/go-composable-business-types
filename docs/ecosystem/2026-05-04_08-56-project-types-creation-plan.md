# Plan: Shared Types for Project Ecosystem — Execution Plan v3

**Date:** 2026-05-04 | **Status:** Planning (after deep reflection v3)
**Critical pivot:** Don't create `project-types`. Extend `go-composable-business-types` instead.

---

## Why Not a New Library?

`go-composable-business-types` already provides:

| What we need | Already exists | Where |
|---|---|---|
| `Importance uint8` (0-100, bounded) | `Percentage uint8` (0-100, bounded) | `types/types_numeric.go` |
| `Tag string` (regex validated, bounded length) | `BoundedString` (length validated) | `bounded/bounded.go` |
| `validate.Validator` interface | Already exists | `validate/validate.go` |
| Sentinel + structured errors | Already exists | `pkg/errors/` |
| JSON/SQL serialization | Already exists per type | `scanutil/`, `types/json.go` |
| `go-branded-id` integration | Already a dependency | `go.mod` |

Creating `project-types` would duplicate ALL of this infrastructure. Instead, add `Importance`, `Tag`, `Language`, `ProjectCore` as new packages/types in `go-composable-business-types`.

Selective imports are by design — consumers only pull what they need.

---

## What Already Exists (Reuse)

```
go-composable-business-types/
├── validate/          ✅ Validator interface — all types implement this
├── pkg/errors/        ✅ Sentinel + structured errors — add new sentinels
├── scanutil/          ✅ SQL Scan/Value helpers — reuse for all types
├── bounded/           ✅ BoundedString — foundation for Tag
├── types/             ✅ Percentage pattern — foundation for Importance
│   └── json.go        ✅ JSON marshal/unmarshal helpers — reuse
├── locale/            ✅ Human languages (BCP 47) — NOT programming languages
└── nanoid/            ❌ Not needed
```

---

## What to Add

### 1. `importance/` — New package

Follows `Percentage` pattern but with **error on overflow** (not clamping) and **classification levels**.

```go
package importance

type Importance uint8

const (
    None     Importance = 0
    VeryLow  Importance = 20
    Low      Importance = 40
    Medium   Importance = 50
    High     Importance = 70
    VeryHigh Importance = 90
    Max      Importance = 100
)

func New(v uint8) (Importance, error)       // error if > 100
func Must(v uint8) Importance               // panic on invalid
func Parse(s string) (Importance, error)    // from "very-low", "medium", etc.
func (i Importance) String() string         // "very-low", "low", "medium", "high", "very-high", "none"
func (i Importance) Classification() string // "Very Low", "Low", etc. (human-readable)
func (i Importance) IsVeryLow() bool        // <= 20
func (i Importance) IsLow() bool            // 21-40
func (i Importance) IsMedium() bool         // 41-60
func (i Importance) IsHigh() bool           // 61-80
func (i Importance) IsVeryHigh() bool       // >= 81
func (i Importance) IsValid() bool          // 0-100
func (i Importance) IsZero() bool
func (i Importance) Percent() float64       // float64(i) / 100
func (i Importance) Validate() error
func (i Importance) Compare(other Importance) int
// JSON marshal/unmarshal (as integer)
// SQL Scan/Value (via scanutil)
```

**Why new package, not in `types/`?** Classification methods (`IsVeryLow`, `IsLow`, etc.) are domain-specific to project management. `types.Percentage` is a generic numeric primitive. Different abstractions, different packages. Follows existing pattern where each domain concept gets its own package.

### 2. `tag/` — New package

Follows `BoundedString` pattern with **regex validation**.

```go
package tag

type Tag string  // named type, not struct (like Email)

// Validation: ^[A-Za-z0-9-]+$
// Max length: 50

func New(s string) (Tag, error)           // validates regex + length
func Must(s string) Tag                    // panic on invalid
func NewTags(ss ...string) ([]Tag, error) // batch construction
func (t Tag) String() string
func (t Tag) IsZero() bool
func (t Tag) Validate() error
func (t Tag) IsValid() bool
// JSON marshal/unmarshal (as string)
// SQL Scan/Value (via scanutil)
```

### 3. `programminglanguage/` — New package

Branded string via `go-branded-id`. No enum. go-enry is the authority.

```go
package programminglanguage

type Language = id.ID[languageBrand, string]

func New(s string) Language                           // wraps, no validation
func FromSlice(ss []string) []Language                // batch wrap
func Normalize(s string) string                       // alias normalization (from SDK's registry)
func StringToCanonicalAliases() map[string][]string   // for reference

type Languages []Language
func NewLanguages(langs ...Language) Languages
func (ls Languages) Primary() Language
func (ls Languages) Has(lang Language) bool
func (ls Languages) IsGo() bool
func (ls Languages) Strings() []string
func (ls Languages) IsZero() bool
func (ls Languages) Validate() error
```

**Why `programminglanguage` not `language`?** Avoid collision with `locale.Locale` (human languages) and Go's builtin `language` package from `golang.org/x/text`. Name is long but unambiguous. Aliased at import: `import pl "github.com/larsartmann/go-composable-business-types/programminglanguage"`.

### 4. `projectcore/` — New package

Composite struct using the types above.

```go
package projectcore

type ProjectCore struct {
    Name       string                            // display name, no validation
    Path       string                            // filesystem path, no validation
    Languages  programminglanguage.Languages     `json:"languages"`
    Importance importance.Importance             `json:"importance"`
    Tags       []tag.Tag                         `json:"tags"`
}

func New(name, path string, langs programminglanguage.Languages, opts ...Option) *ProjectCore
func (p *ProjectCore) Validate() error
func (p *ProjectCore) IsZero() bool
// Functional options: WithImportance, WithTags
// JSON marshal/unmarshal
```

---

## Updated Architecture Layer Map

```
┌───────────────────────────────────────────────────────────────┐
│  CLI / Presentation                                           │
│  cmdguard · go-output · go-commit                             │
├───────────────────────────────────────────────────────────────┤
│  Workflow / Orchestration                                     │
│  universal-workflow · ActaFlow · go-filewatcher               │
├───────────────────────────────────────────────────────────────┤
│  Domain / Business Logic                                      │
│  go-business-rules · go-composable-business-types             │  ← Importance, Tag, Language, ProjectCore live here
│    ├── importance/                                             │
│    ├── tag/                                                    │
│    ├── programminglanguage/                                    │
│    └── projectcore/                                            │
├───────────────────────────────────────────────────────────────┤
│  Domain / Discovery                                            │
│  project-discovery-sdk · project-dependency-graph              │
├───────────────────────────────────────────────────────────────┤
│  Infrastructure / Config / Storage                             │
│  smart-configs · project-meta · go-branded-id                  │
└───────────────────────────────────────────────────────────────┘
```

---

## Execution Steps — Sorted by Impact / Effort

Each step is a self-contained git commit.

### Phase 1: Add types to go-composable-business-types (Foundation)

| # | Step | Impact | Effort | Self-contained? |
|---|------|--------|--------|-----------------|
| 1.1 | Add `importance/` package: type, constructor, classification, JSON, tests | HIGH | 1.5hr | ✅ |
| 1.2 | Add `tag/` package: type, constructor, regex validation, JSON, tests | HIGH | 1hr | ✅ |
| 1.3 | Add `programminglanguage/` package: Language type, Languages slice, normalization map, tests | HIGH | 1.5hr | ✅ |
| 1.4 | Add `projectcore/` package: ProjectCore struct, options, validation, tests | MED | 1hr | ✅ |
| 1.5 | Add sentinel errors to `pkg/errors/` for all new types | MED | 20min | ✅ |
| 1.6 | Update README + LIBRARY_GUIDE.md with new types | LOW | 30min | ✅ |
| 1.7 | Full test suite + `nix flake check` | MED | 30min | ✅ |

### Phase 2: Wire into project-discovery-sdk (Lowest risk)

| # | Step | Impact | Effort | Self-contained? |
|---|------|--------|--------|-----------------|
| 2.1 | Add `go-composable-business-types` to SDK go.mod (replace → local) | MED | 10min | ✅ |
| 2.2 | Replace SDK `Project.Language string` + `Languages []string` with `programminglanguage.Language` + `Languages` | HIGH | 1hr | ✅ |
| 2.3 | Replace SDK `Project.Importance int` with `importance.Importance` | HIGH | 30min | ✅ |
| 2.4 | Replace SDK `Project.Tags []string` with `[]tag.Tag` | HIGH | 30min | ✅ |
| 2.5 | Remove SDK's `detection/NormalizeLanguage` → use `programminglanguage.Normalize` | MED | 30min | ✅ |
| 2.6 | Remove SDK's `branded.go` Language/ProjectName/ProjectPath types → use `programminglanguage.Language` | MED | 30min | ✅ |
| 2.7 | Update SDK tests | MED | 1hr | ✅ |

### Phase 3: Wire into project-meta

| # | Step | Impact | Effort | Self-contained? |
|---|------|--------|--------|-----------------|
| 3.1 | Add `go-composable-business-types` to project-meta go.mod | MED | 10min | ✅ |
| 3.2 | Replace project-meta's `Importance int32` with `importance.Importance uint8` | HIGH | 1hr | ✅ |
| 3.3 | Replace project-meta's `Tag` with `tag.Tag` | HIGH | 30min | ✅ |
| 3.4 | Add backward-compat YAML unmarshaling (accept both int32 and uint8) | MED | 30min | ✅ |
| 3.5 | Update all project-meta tests | MED | 1hr | ✅ |

### Phase 4: Wire into projects-management-automation (Highest risk)

| # | Step | Impact | Effort | Self-contained? |
|---|------|--------|--------|-----------------|
| 4.1 | Add `go-composable-business-types` to PMA go.mod + go.work | MED | 10min | ✅ |
| 4.2 | Replace PMA's `Importance int32` with `importance.Importance uint8` | HIGH | 1hr | ✅ |
| 4.3 | Replace PMA's `types.Tag` with `tag.Tag` | HIGH | 30min | ✅ |
| 4.4 | Replace PMA's `types.Language uint8` enum with `programminglanguage.Language` | HIGH | 2hr | ✅ |
| 4.5 | Convert language switch statements to dispatch maps | MED | 1hr | ✅ |
| 4.6 | Delete `meta/adapter.go` — no type conversion needed | HIGH | 15min | ✅ |
| 4.7 | Fix Importance type leaks (raw int/int32) | MED | 30min | ✅ |
| 4.8 | Update all PMA tests | MED | 2hr | ✅ |

### Phase 5: Wire into project-dependency-graph

| # | Step | Impact | Effort | Self-contained? |
|---|------|--------|--------|-----------------|
| 5.1 | Add `go-composable-business-types` to PDG go.mod | MED | 10min | ✅ |
| 5.2 | Use `programminglanguage.Language` for filtering | LOW | 15min | ✅ |
| 5.3 | Full test verification | LOW | 30min | ✅ |

---

## What I Got Wrong in v1 and v2

| Mistake | Correction |
|---|---|
| Create new `project-types` library | Extend `go-composable-business-types` — it already has the infrastructure |
| Use govalid for validation | Use the existing `validate.Validator` interface + constructor pattern |
| `ProjectName`/`ProjectPath` in kernel | Keep in project-meta — they're storage identifiers, not shared |
| `GitStatus` in kernel | Keep in PMA — only PMA uses it |
| `ProjectType` in kernel | Keep in PMA — it's a PMA-specific classification |
| `Languages` as `[]Language` with set library | `[]Language` with `slices.Contains` — no external dependency needed |
| Tag regex change needs migration tool | Migration is simple: `strings.ReplaceAll(tag, "_", "-")` — add to `doctor --fix-tags` |

---

## govalid Decision

**Drop govalid.** Reasoning:

1. `go-composable-business-types` already has a validation pattern (`validate.Validator` interface + constructors)
2. govalid can't validate `id.ID` struct fields (sees struct, not string)
3. govalid can't do regex (Tag validation)
4. The value was 2 `lte` checks that constructors already cover
5. Consistency > novelty — follow the existing pattern in the library

---

## Tag Migration

Old regex: `^[a-z0-9_-]+$` → New regex: `^[A-Za-z0-9-]+$`

Migration: `strings.ReplaceAll(tag, "_", "-")` — deterministic, reversible.
Add to project-meta's `doctor --fix-tags` command.

---

## Dependencies After Change

```
go-composable-business-types depends on:
  - go-branded-id (already)
  - golang.org/x/text (already, for locale)
  - No NEW dependencies

project-discovery-sdk depends on:
  - go-composable-business-types (NEW, replaces some internal types)
  - go-enry (already)
  - go-branded-id (already, now transitive via business-types)

project-meta depends on:
  - go-composable-business-types (NEW, replaces internal Importance/Tag)

projects-management-automation depends on:
  - go-composable-business-types (NEW, replaces internal Importance/Tag/Language)
  - project-discovery-sdk (already)
  - project-meta (already)

project-dependency-graph depends on:
  - go-composable-business-types (NEW, for Language type)
  - project-discovery-sdk (already)
```

Dependency direction is clean — no cycles, all point downward.
