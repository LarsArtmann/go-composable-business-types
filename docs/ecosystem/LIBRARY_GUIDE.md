# Project Library Guide — When to Use What

> A decision guide for knowing which library to reach for and where to start reading.

**Last Updated:** 2026-05-03

---

## Quick Decision Matrix

| If you need...                                                           | Use this                                                      |
| ------------------------------------------------------------------------ | ------------------------------------------------------------- |
| Event sourcing / CQRS architecture                                       | [go-cqrs-lite](#go-cqrs-lite)                                 |
| Strongly typed business data with audit trails                           | [go-composable-business-types](#go-composable-business-types) |
| Shared project types (Importance, Tag, ProgrammingLanguage, ProjectCore) | [go-composable-business-types](#go-composable-business-types) |
| Unified static analysis findings / SARIF                                 | [go-finding](#go-finding)                                     |
| Severity-aware business rule validation                                  | [go-business-rules](#go-business-rules)                       |
| Type-safe CLI construction with DI                                       | [cmdguard](#cmdguard)                                         |
| File system watching with debouncing                                     | [go-filewatcher](#go-filewatcher)                             |
| AI-generated git commit messages                                         | [go-commit](#go-commit)                                       |
| Smart configuration with actionable error messages                       | [smart-configs](#smart-configs)                               |
| In-memory actor model (Erlang-style)                                     | [ActaFlow](#actaflow)                                         |
| Dependency-aware workflow orchestration                                  | [universal-workflow](#universal-workflow)                     |
| Multi-format structured data output                                      | [go-output](#go-output)                                       |
| Offline-first app framework with CRDT sync                               | [go-localfirst](#go-localfirst)                               |
| Sync external APIs into local SQLite                                     | [go-localsync](#go-localsync)                                 |

---

## Architecture Layer Map

```
┌───────────────────────────────────────────────────────────────┐
│  CLI / Presentation                                           │
│  cmdguard · go-output · go-commit (CLI)                       │
├───────────────────────────────────────────────────────────────┤
│  Workflow / Orchestration                                     │
│  universal-workflow · ActaFlow · go-filewatcher               │
├───────────────────────────────────────────────────────────────┤
│  Domain / Business Logic                                      │
│  go-business-rules · go-composable-business-types             │
│    ├── importance · tag · programminglanguage · projectcore   │
├───────────────────────────────────────────────────────────────┤
│  Domain / Discovery                                            │
│  project-discovery-sdk · project-dependency-graph              │
├───────────────────────────────────────────────────────────────┤
│  Data / Sync / Storage                                        │
│  go-cqrs-lite · go-localfirst · go-localsync · go-finding     │
├───────────────────────────────────────────────────────────────┤
│  Infrastructure / Config                                      │
│  smart-configs                                                │
└───────────────────────────────────────────────────────────────┘
```

---

## Combinations That Work Well Together

| Pattern                     | Libraries                                                     |
| --------------------------- | ------------------------------------------------------------- |
| Full CLI app                | cmdguard + smart-configs + go-output + go-business-rules      |
| Event-sourced microservice  | go-cqrs-lite + go-composable-business-types + smart-configs   |
| Offline-first app with sync | go-localfirst + go-localsync + go-cqrs-lite                   |
| Static analysis pipeline    | go-finding + go-output + go-filewatcher                       |
| Complex workflow engine     | universal-workflow + go-composable-business-types + go-output |
| Developer tooling CLI       | cmdguard + go-filewatcher + go-commit + go-output             |

---

## Library Details

---

### go-cqrs-lite

**Path:** `/home/lars/projects/go-cqrs-lite/`
**Module:** `github.com/LarsArtmann/go-cqrs-lite`

#### When to Consider

- You need **event sourcing** — append-only event streams, optimistic concurrency, state reconstruction from history
- You need **CQRS** — strict separation of write (commands) and read (queries) paths
- You're doing **Domain-Driven Design** with aggregate roots and strongly-typed IDs
- You need a **transactional outbox** for reliable event publishing in microservices
- You want **auto-generated API documentation** (AsyncAPI, EventCatalog, D2) from Go struct types
- You prefer **functional domain modeling** via the `Decider[State]` pattern
- You need **projections** with replay and live subscription

#### Key Abstractions

| Type                     | Purpose                                                        |
| ------------------------ | -------------------------------------------------------------- |
| `command.Dispatcher`     | Type-safe command routing with middleware                      |
| `query.Dispatcher`       | Type-safe query routing with `DispatchTyped[T]`                |
| `event.Store`            | Append-only event storage with optimistic concurrency          |
| `event.Bus`              | Publisher/Subscriber split with middleware                     |
| `aggregate.Root`         | Event-sourced aggregate with `RecordEvent` / `LoadFromHistory` |
| `decider.Decider[State]` | Functional approach: `Initial` + `Fold` + `Repository.Execute` |
| `id.Of[T]`               | Branded ULID-backed IDs — compile-time prevention of ID mixing |
| `projection.Runner`      | Replay + live subscription with checkpoints                    |

#### Start Reading

| File                         | Why                                              |
| ---------------------------- | ------------------------------------------------ |
| `README.md`                  | Full overview, quick start, architecture diagram |
| `FEATURES.md`                | Honest feature inventory with maturity matrix    |
| `core/event/event.go`        | Central `Event` interface and `Core` struct      |
| `core/event/store.go`        | `Store` interface — the persistence contract     |
| `core/command/dispatcher.go` | Command dispatcher                               |
| `core/decider/decider.go`    | Functional `Decider[State]` pattern              |
| `example/user/main.go`       | Minimal end-to-end demo                          |

---

### go-composable-business-types

**Path:** `/home/lars/projects/go-composable-business-types/`
**Module:** `github.com/artmann/go-composable-business-types`

#### When to Consider

- You need **bitemporal data** tracking (when facts were true in real world vs. when recorded) — financial, legal, medical, compliance systems
- You want **compile-time type-safe IDs** that prevent mixing `UserID` with `OrderID`
- You need **full audit trails baked into your data model** — every `DataPoint[T]` carries who, when, why, and what caused it
- You need **validated domain primitives** — `Email`, `URL`, `Percentage`, `Cents`, `Money`
- You want **float-free money arithmetic** — `Cents` (int64) and `Money` (ISO 4217)
- You need **actor chain tracking** for distributed systems (User → API Gateway → Service)
- You want **causal chain** tracking for data lineage
- You need **shared project ecosystem types** — `Importance`, `Tag`, `ProgrammingLanguage`, `ProjectCore` for `project-discovery-sdk`, `project-meta`, `projects-management-automation`, `project-dependency-graph`

#### Key Abstractions

| Type                           | Purpose                                                                                                      |
| ------------------------------ | ------------------------------------------------------------------------------------------------------------ |
| `DataPoint[T]`                 | Generic wrapper with NanoID, bitemporal tracking, actor, trigger, context, references, causes, tags, version |
| `Bitemporal`                   | Two time axes: valid-time and transaction-time with point-in-time queries                                    |
| `ActorChain[T]`                | Ordered chain of actors for full audit trail                                                                 |
| `ID[B, V]`                     | Branded/phantom IDs — compile-time safety                                                                    |
| `Cents` / `Money`              | Float-free monetary arithmetic                                                                               |
| `Email` / `URL` / `Percentage` | Validated domain primitives                                                                                  |
| `Importance`                   | `uint8` 0–100 with classification levels (VeryLow/Low/Medium/High/VeryHigh)                                  |
| `Tag`                          | Validated string `^[A-Za-z0-9-]+$`, max 50 chars                                                             |
| `Language`                     | Branded string via go-branded-id, go-enry authority, no closed enum                                          |
| `ProjectCore`                  | Shared base entity: Name, Path, Languages, Importance, Tags                                                  |

#### Start Reading

| File                              | Why                                               |
| --------------------------------- | ------------------------------------------------- |
| `README.md`                       | Full usage guide with examples                    |
| `PARTS.md`                        | Component analysis — rates each type's uniqueness |
| `importance/importance.go`        | `Importance` — bounded uint8 with classification  |
| `tag/tag.go`                      | `Tag` — regex-validated identifier                |
| `programminglanguage/language.go` | `Language` — branded string, go-enry compatible   |
| `projectcore/project_core.go`     | `ProjectCore` — shared base entity                |
| `datapoint/datapoint.go`          | `DataPoint[T]` — the main abstraction             |
| `temporal/temporal.go`            | `Bitemporal` — bitemporal tracking                |
| `actor/actor.go`                  | `ActorChain[T]`, `ActorEntry[T]`                  |
| `types/types.go`                  | `Email`, `URL` domain primitives                  |
| `types/types_numeric.go`          | `Percentage`, `Cents` with arithmetic             |

---

### go-finding

**Path:** `/home/lars/projects/go-finding/`
**Module:** `github.com/larsartmann/go-finding`

#### When to Consider

- You're **building a static analysis tool in Go** and need a standard output model
- You need to **merge reports from multiple linters** with deduplication and cross-tool correlation
- You want **SARIF 2.1.0** output for GitHub Code Scanning / Azure DevOps integration
- You need **LSP diagnostic** conversion for IDE extensions
- You want an **automated detect → triage → fix → verify** pipeline with conflict detection
- You're wrapping existing `go/analysis` tools into a unified model

#### Key Abstractions

| Type                 | Purpose                                                                              |
| -------------------- | ------------------------------------------------------------------------------------ |
| `Finding`            | Unified finding model: ID, rule, severity, position, fix strategy, before/after code |
| `Report`             | Thread-safe container for findings with SARIF/JSON export                            |
| `pipeline.Detector`  | Interface for pluggable analysis tools                                               |
| `pipeline.Pipeline`  | Orchestrates detect → triage → fix → verify loop                                     |
| `pipeline.FixEngine` | Pure in-memory fix application with conflict detection                               |
| `LSPDiagnostic`      | Bidirectional LSP conversion                                                         |

#### Start Reading

| File                        | Why                                             |
| --------------------------- | ----------------------------------------------- |
| `README.md`                 | Overview, quick start, API examples             |
| `FEATURES.md`               | Definitive feature reference                    |
| `docs/USAGE_GUIDE.md`       | Comprehensive usage with code examples          |
| `finding.go`                | The `Finding` struct — the heart of the library |
| `pipeline/pipeline.go`      | Pipeline orchestration and `Detector` interface |
| `examples/pipeline/main.go` | Pipeline with custom Detector                   |

---

### go-business-rules

**Path:** `/home/lars/projects/go-business-rules/`
**Module:** `github.com/artmann/businessrules`

#### When to Consider

- You need **severity-aware validation** — not all failures are equal (block on critical, warn on suspicious, advise on style)
- You want to **complement structural validators** (field format checks) with **business rule checks** (domain logic)
- You need **JSON-serializable validation results** for API responses
- You want **mergeable results** from multiple validation sources
- You need **composite combinators** — `All()`, `Any()`, `When()` for conditional rules
- You want **20+ pre-built rules** with a fluent builder API

#### Key Abstractions

| Type                    | Purpose                                                |
| ----------------------- | ------------------------------------------------------ |
| `Rule` (interface)      | `Name()`, `Check() error`, `Severity()`, `Message()`   |
| `Severity`              | 4-level: Info → Warning → Error → Critical             |
| `ValidatorBuilder`      | Fluent builder: `AddRule().AddRules().Build()`         |
| `ValidationResultError` | Rich result with filtering by severity, merge, forEach |

#### Start Reading

| File                    | Why                                               |
| ----------------------- | ------------------------------------------------- |
| `README.md`             | Full API reference and philosophy                 |
| `rule.go`               | Core `Rule` interface and `NewRule()` constructor |
| `severity.go`           | The `Severity` type — foundational                |
| `validator.go`          | `ValidatorBuilder` — the entry point              |
| `validation_result.go`  | `ValidationResultError` — the rich result         |
| `builders_composite.go` | `All`, `Any`, `When` combinators                  |
| `example_test.go`       | Runnable examples                                 |

---

### cmdguard

**Path:** `/home/lars/projects/cmdguard/`
**Module:** `github.com/LarsArtmann/cmdguard` (v2 at `pkg/cmdguard/v2/`)

#### When to Consider

- You're building a **production Go CLI** and want type-safe, validated construction
- You want **typed flags from structs** with struct tags instead of scattered `cmd.Flags().String(...)` calls
- You need **dependency injection** (database, logger, HTTP client managed per-command)
- You want **environment variable fallbacks** with prefix support (12-factor apps)
- You need **graceful shutdown**, **health checks**, and **lifecycle hooks**
- You want **12 output formats** (JSON, YAML, table, CSV, D2, etc.) built in
- You're tired of Cobra footguns (panics, stringly-typed flags, no validation)

#### Key Abstractions

| Type                   | Purpose                                                   |
| ---------------------- | --------------------------------------------------------- |
| `CLI[T]`               | Main entry point. `T` = app config type                   |
| `Command[T, F]`        | Type-safe command with config `T` and flags `F`           |
| `Scope`                | DI scope with `Provide`/`Invoke`, health checks, shutdown |
| `FlagRegistry`         | Struct-tag-driven flag registration, parsing, validation  |
| `TypeHandler`          | Extensible dispatch for custom flag types                 |
| `Middleware[T]`        | Chainable cross-cutting concerns (timing, recovery)       |
| `BranchingFlowContext` | Tree-shaped context through the command hierarchy         |

#### Start Reading

| File                            | Why                                 |
| ------------------------------- | ----------------------------------- |
| `README.md`                     | Full API reference, quick start     |
| `docs/QUICKSTART.md`            | 5-minute getting started guide      |
| `docs/CLI_DESIGN_PRINCIPLES.md` | Design philosophy                   |
| `pkg/cmdguard/v2/cli.go`        | `CLI[T]` — the heart of the library |
| `pkg/cmdguard/v2/command.go`    | `Command[T,F]`, options, validation |
| `pkg/cmdguard/v2/scope.go`      | DI `Scope`                          |
| `examples/typed/main.go`        | Full example with all features      |

---

### go-filewatcher

**Path:** `/home/lars/projects/go-filewatcher/`

#### When to Consider

- You need to **react to file system changes** without raw fsnotify boilerplate
- You want **automatic recursive watching** (including newly created directories)
- You need **smart debouncing** — global (coalesce all events) or per-path (each file independently)
- You want **15+ composable filters** with AND/OR/NOT combinators
- You need to **exclude auto-generated code** (sqlc, protobuf, templ, mockgen)
- You want **middleware chains** for logging, recovery, rate limiting, metrics
- You're building build tools, hot-reload servers, code generation pipelines, or log processors

#### Key Abstractions

| Type                                     | Purpose                                                               |
| ---------------------------------------- | --------------------------------------------------------------------- |
| `Watcher`                                | Central struct with config, state, debouncer — implements `io.Closer` |
| `Filter`                                 | `func(Event) bool` — composable via `FilterAnd/Or/Not`                |
| `Middleware`                             | `func(Handler) Handler` — wraps handlers for cross-cutting concerns   |
| `Debouncer` / `GlobalDebouncer`          | Per-path or global debouncing strategies                              |
| `EventPath` / `RootPath` / `DebounceKey` | Phantom-typed paths — compile-time safety                             |

#### Start Reading

| File                                | Why                                             |
| ----------------------------------- | ----------------------------------------------- |
| `README.md`                         | Full usage docs and API reference               |
| `ARCHITECTURE.md`                   | Layer diagram, concurrency model, state machine |
| `watcher.go`                        | `Watcher` struct, `New()`, `Watch()`, `Close()` |
| `filter.go`                         | All 15 filters + composition combinators        |
| `middleware.go`                     | 10 built-in middleware                          |
| `examples/basic/main.go`            | Minimal usage                                   |
| `examples/filter-generated/main.go` | Excluding auto-generated code                   |

---

### go-commit

**Path:** `/home/lars/projects/go-commit/`

#### When to Consider

- You want **AI-generated git commit messages** with automatic fallback (Groq → OpenAI → heuristic)
- You need **branch-aware** commit scope/type/ticket extraction from branch names
- You want **branch hierarchy management** with parent-child relationships and cycle detection
- You need **branch sync** with merge, rebase, compress, or ff-only strategies
- You want **AI-powered merge conflict resolution** with confidence scoring
- You're building **CI/CD pipelines** that need automated, well-formatted commits

#### Key Abstractions

| Type                 | Purpose                                                            |
| -------------------- | ------------------------------------------------------------------ |
| `Commit`             | Main orchestrator — holds git operations + provider + config       |
| `providers.Provider` | `GenerateCommitMessage(ctx, CommitRequest) (string, error)`        |
| `providers.Chain`    | Tries providers in order until one succeeds                        |
| `git.Operations`     | Composite interface: status, staging, committing, pushing, syncing |
| `branch.Hierarchy`   | Parent/child relationships with cycle detection                    |
| `conflict.Resolver`  | AI-powered conflict resolution with confidence scoring             |

#### Start Reading

| File                               | Why                                          |
| ---------------------------------- | -------------------------------------------- |
| `README.md`                        | Full usage guide                             |
| `docs/ARCHITECTURE.md`             | Architecture decisions and data flow         |
| `pkg/commit/commit.go`             | Core `Commit` struct — main entry point      |
| `pkg/commit/providers/provider.go` | `Provider` interface, `Chain`, `CommitStyle` |
| `pkg/commit/git/git.go`            | `Operations` interface                       |
| `pkg/commit/conflict/resolver.go`  | AI conflict resolution                       |
| `pkg/commit/branch/hierarchy.go`   | Branch tree management                       |

---

### smart-configs

**Path:** `/home/lars/projects/smart-configs/`

#### When to Consider

- You're tired of `"DATABASE_URL not found"` errors with no context — you want **actionable error messages** with copy-paste fix commands
- You need **multi-source config resolution**: CLI args → env → `.env` → CLI tools → config files → cache → defaults
- You want **CI/CD-aware suggestions** — auto-detects GitHub Actions, Docker, K8s and tailors guidance
- You need **service-specific resolvers** (Turso DB URLs, Groq API keys, GitHub tokens via `gh` CLI)
- You want **type-safe generics** — `Get[int]("PORT")`, `Get[bool]("DEBUG")`
- You need **tool detection** with installation guidance (brew, curl, nix)

#### Key Abstractions

| Type                  | Purpose                                                                      |
| --------------------- | ---------------------------------------------------------------------------- |
| `Config`              | Main entry point — resolver + context + formatters + cache                   |
| `Source` (interface)  | `Resolve(ctx, key) → (string, bool, error)`                                  |
| `ExecutionContext`    | Full runtime picture: CI, cloud, container, git, runtime, user               |
| `Suggestion`          | Actionable fix: title, command, URL, priority, platform, confidence          |
| `SmartError`          | Composable error with aspects (key, message, suggestions, security warnings) |
| `SuggestionGenerator` | Pluggable interface for service-specific suggestions                         |

#### Start Reading

| File                                         | Why                                        |
| -------------------------------------------- | ------------------------------------------ |
| `README.md`                                  | Full feature overview, quick start         |
| `QUICK_START.md`                             | 5-minute getting started guide             |
| `AGENTS.md`                                  | Architecture, patterns, gotchas            |
| `main.go`                                    | All exported types, functions, and options |
| `internal/configs/smartconfig.go`            | `Config` struct and `New()`                |
| `internal/suggestions/suggestions_engine.go` | Suggestion dispatch engine                 |
| `docs/services/turso.md`                     | Turso integration example                  |
| `examples/usage/complete_example.go`         | Comprehensive usage demo                   |

---

### ActaFlow

**Path:** `/home/lars/projects/ActaFlow/`

#### When to Consider

- You need **high-concurrency in-memory processing** — trading engines, game servers, real-time analytics, chat systems
- You want the **Erlang/OTP actor model** (supervision, isolation, message passing) in Go
- You need **privacy-preserving architecture** — GDPR/CCPA-aware contextual IDs preventing cross-domain correlation
- You want **schema-first development** — define domain in TypeSpec, get auto-generated Go types
- You need **zero-trust security** — HMAC message signing, trust levels, RBAC policies, audit logging
- You want **"impossible states unrepresentable"** — typed enums prevent invalid state transitions

**Note:** Single-node, in-memory only. No persistence, no distributed runtime. ~80% complete.

#### Key Abstractions

| Type                  | Purpose                                                                      |
| --------------------- | ---------------------------------------------------------------------------- |
| `ActorSystem`         | System lifecycle, health, metrics                                            |
| `ActorRef`            | Actor interaction: `Send`, `Tell`, `Ask`, lifecycle hooks                    |
| `PureActor[S]`        | Generic typed actor with state `S`, message handling, supervision directives |
| `Context`             | Actor execution context with security, tracing, child spawning               |
| `ActorState`          | Lifecycle enum with `CanTransitionTo()` state machine                        |
| `SupervisionStrategy` | OneForOne, OneForAll, RestForOne, Escalate                                   |
| `SecurityContext`     | AuthN/AuthZ with trust levels and permissions                                |

#### Start Reading

| File                                            | Why                                         |
| ----------------------------------------------- | ------------------------------------------- |
| `README.md`                                     | Project overview, quick start, status       |
| `HOW_TO_USE.md`                                 | Step-by-step usage guide with code examples |
| `USAGE.md`                                      | Comprehensive API reference                 |
| `docs/what-we-are-not.md`                       | Clear limitations                           |
| `pkg/interfaces/interfaces.go`                  | All public interfaces                       |
| `internal/actor/interfaces.go`                  | Core actor interfaces and `PureActor[S]`    |
| `docs/architecture/supervision-architecture.md` | Supervision tree design                     |
| `docs/security/SECURITY_GUIDE.md`               | Security guide                              |

---

### universal-workflow

**Path:** `/home/lars/projects/universal-workflow/`
**Module:** `github.com/LarsArtmann/universal-workflow`

#### When to Consider

- You have **multi-step processes with dependencies** (e.g., order: validate → pay + reserve → confirm)
- You need **parallel execution** of independent tasks with automatic synchronization
- You want **type-safe workflow definitions** with compile-time-checked IDs
- You need **retry logic**, timeouts, and error propagation across a pipeline
- You want **conditional branching** — success/failure paths, custom conditions
- You need **DAG visualization** — Mermaid, DOT, JSON export
- You want **distributed tracing** via OpenTelemetry

**Note:** In-process library. No external server. No durable state for long-running sagas.

#### Key Abstractions

| Type                 | Purpose                                                          |
| -------------------- | ---------------------------------------------------------------- |
| `UnifiedWorkflow`    | Main workflow — holds activities, state machine, execution mode  |
| `ActivityHandler`    | `func(ctx ActivityContext) (*ActivityResult, error)`             |
| `ActivityBuilder`    | Fluent: `.WithName().DependsOn().WithRetryPolicy()`              |
| `DependencyResolver` | Kahn's algorithm for topological sort                            |
| `ParallelExecutor`   | Concurrent execution within dependency levels                    |
| `RetryInterceptor`   | Exponential backoff with jitter                                  |
| `EventPublisher`     | Pub/sub for workflow events with filtered subscriptions          |
| `BranchPolicy`       | Conditional branching: Always, IfSuccess, IfFailure, IfCondition |

#### Start Reading

| File                                  | Why                                              |
| ------------------------------------- | ------------------------------------------------ |
| `README.md`                           | Full overview, quick start                       |
| `examples/order-processing/main.go`   | Complete working example using all features      |
| `pkg/workflow/workflow.go`            | `UnifiedWorkflow` — the heart                    |
| `pkg/workflow/workflow_types.go`      | `Activity`, `ActivityHandler`, `ActivityBuilder` |
| `pkg/workflow/dependency_resolver.go` | Topological sort implementation                  |
| `pkg/workflow/parallel_executor.go`   | Parallel execution                               |
| `pkg/workflow/branch_execution.go`    | Conditional branching                            |

---

### go-output

**Path:** `/home/lars/projects/go-output/`

#### When to Consider

- You need **12 output formats** from one data model — JSON, CSV, Markdown, HTML, Mermaid, D2, DOT, YAML, XML, TSV, ASCII tree, styled tables
- You want **type-safe format enums** with `Parse()`, `IsValid()`, `AllowedValues()`
- You need **branded IDs** to prevent mixing node IDs across diagram types at compile time
- You want a **pluggable renderer registry** — register custom formats at runtime
- You need **streaming output** for large datasets via `StreamingRenderer`
- You want **color that respects CI/NO_COLOR** automatically
- You're building a CLI that needs `--format` and `--sort-by` flags

#### Key Abstractions

| Type               | Purpose                                                     |
| ------------------ | ----------------------------------------------------------- |
| `Renderer`         | Core interface: `Render() (string, error)`                  |
| `TableRenderer`    | Extends `Renderer` with `SetHeaders` / `AddRow`             |
| `GraphRenderer`    | Extends `Renderer` with `SetNodes` / `SetEdges`             |
| `Format`           | 12 type-safe format constants                               |
| `TableData`        | Universal table data container (headers + rows)             |
| `D2Diagram`        | Full D2 diagram builder (nodes, edges, SQL tables, classes) |
| `TreeNode`         | Hierarchical tree with metadata                             |
| `BrandedID[Brand]` | Phantom-typed string IDs for compile-time safety            |

#### Start Reading

| File                     | Why                                                     |
| ------------------------ | ------------------------------------------------------- |
| `README.md`              | Full API tour with examples                             |
| `format.go`              | Core interfaces, `Format` enum, `TableData`, `TreeNode` |
| `graph.go`               | `GraphRenderer`, `GraphNode`, `GraphEdge`               |
| `d2_render.go`           | D2 diagram builder                                      |
| `registry.go`            | Pluggable renderer registry                             |
| `examples/basic/main.go` | All 12 formats demonstrated                             |
| `examples/d2/main.go`    | Advanced D2 example                                     |

---

### go-localfirst

**Path:** `/home/lars/projects/go-localfirst/`

#### When to Consider

- You need an app that **works offline** and syncs when connectivity returns
- You need **multi-device/multi-user data sync** without a central authority
- You want **conflict-free merging** of concurrent edits via CRDTs
- You need a **Go backend with event sourcing** for full audit trail and state replay
- You want an **embedded database** (CockroachDB Pebble) — no external DB dependency
- You need **vector clocks** and **LWW conflict resolution** primitives (import `pkg/sync` alone)

#### Key Abstractions

| Type                           | Purpose                                                    |
| ------------------------------ | ---------------------------------------------------------- |
| `pkg/sync.VectorClock`         | Causal ordering: `Increment`, `Merge`, `Compare`, `Clone`  |
| `pkg/sync.Operation[T]`        | Generic typed sync operation with payload and vector clock |
| `pkg/sync.ConflictResolver[T]` | Pluggable conflict resolution interface                    |
| `pkg/sync.LWWResolver[T]`      | Last-Write-Wins: vector clock → timestamp → tiebreaker     |
| `Manager`                      | WebSocket sync orchestrator with peer management           |
| `PebbleStore`                  | Embedded KV store implementing `TodoRepository`            |
| `aggregate.Todo`               | Event-sourced aggregate (CQRS via go-cqrs-lite)            |

#### Start Reading

| File                                           | Why                                      |
| ---------------------------------------------- | ---------------------------------------- |
| `README.md`                                    | Project overview, quick start            |
| `ARCHITECTURE.md`                              | Full architecture with Mermaid diagrams  |
| `pkg/sync/doc.go`                              | Reusable sync SDK overview with examples |
| `pkg/sync/vectorclock.go`                      | `VectorClock` — causal ordering          |
| `pkg/sync/conflict.go`                         | `ConflictResolver[T]`, `LWWResolver[T]`  |
| `internal/sync/manager.go`                     | WebSocket sync manager                   |
| `docs/adr/001-conflict-resolution-strategy.md` | Conflict resolution ADR                  |

---

### go-localsync

**Path:** `/home/lars/projects/go-localsync/`

#### When to Consider

- You need to **sync data from external APIs into local SQLite** — GitHub, Jira, GitLab, etc.
- You want **offline-first dashboards** that aggregate data from multiple sources
- You need **sub-millisecond query latency** against cached API data
- You want **incremental sync** — only fetch items newer than latest stored
- You need **conflict-aware sync** for multi-source scenarios (LWW conflict detection)
- You want **rate limiting** and **retry with backoff** built in
- You need **full-fidelity storage** — complete original JSON payloads preserved
- You want a **custom sync pipeline** where you control the provider, schedule, and storage

#### Key Abstractions

| Type                       | Purpose                                                                         |
| -------------------------- | ------------------------------------------------------------------------------- |
| `provider.Provider`        | Interface: `Name()`, `Fetch()`, `FetchAll()`, `GetRateLimit()`                  |
| `storage.Storage`          | Composed of `Reader` (12 methods) + `Writer` + `Close`                          |
| `sync.Syncer`              | Full and incremental sync engine                                                |
| `sync.ConflictAwareSyncer` | Wraps `Syncer` with LWW conflict resolution                                     |
| `provider.Item`            | Universal sync unit with full JSON payload                                      |
| Branded IDs                | `ItemID`, `ExternalID`, `ProviderID`, `ActorID`, `RepoID` — compile-time safety |

#### Start Reading

| File                               | Why                                                  |
| ---------------------------------- | ---------------------------------------------------- |
| `README.md`                        | Complete overview, quick start, architecture diagram |
| `ROADMAP.md`                       | Project direction and completed work                 |
| `pkg/provider/provider.go`         | `Provider` interface, `Item` struct                  |
| `pkg/storage/interface.go`         | `Storage`/`Reader`/`Writer` interfaces               |
| `pkg/sync/sync.go`                 | `Syncer` — full and incremental sync                 |
| `pkg/sync/conflict_aware.go`       | LWW conflict resolution                              |
| `pkg/providers/github/client.go`   | Reference `Provider` implementation                  |
| `cmd/examples/github-sync/main.go` | Complete CLI example                                 |

---

_Generated with research from all 13 project repositories._
