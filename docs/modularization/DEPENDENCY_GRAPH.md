# Dependency Graph — go-composable-business-types

## Module Dependency Matrix (Current)

| Module ↓ depends on → | Root | nanoid | locale | money | datapoint | examples |
| --------------------- | ---- | ------ | ------ | ----- | --------- | -------- |
| **Root**              | —    | ✗      | ✗      | ✗     | ✗         | ✗        |
| **nanoid**            | ✓    | —      | ✗      | ✗     | ✗         | ✗        |
| **locale**            | ✓    | ✗      | —      | ✗     | ✗         | ✗        |
| **money**             | ✓    | ✗      | ✓      | —     | ✗         | ✗        |
| **datapoint**         | ✓    | ✓      | ✗      | ✗     | —         | ✗        |
| **examples**          | ✓    | ✓      | ✗      | ✗     | ✓         | —        |

All edges point upward. No cycles. ✓

## Module Details

```
                        ┌──────────────┐
                        │    ROOT      │
                        │  13 packages │
                        │  go-branded-id │
                        │  go-enum (tool) │
                        └──────┬───────┘
                               │
              ┌────────────────┼────────────────┐
              │                │                │
     ┌────────┴───────┐ ┌─────┴──────┐ ┌───────┴──────┐
     │    nanoid      │ │   locale   │ │    money     │
     │ sixafter/nanoid│ │ golang.org/│ │ bojanz/      │
     │                │ │   x/text   │ │  currency    │
     └────────┬───────┘ └─────┬──────┘ └───────┬──────┘
              │               │                │
              │               └───────┬────────┘
              │                       │
     ┌────────┴───────────────────────┘
     │
┌────┴──────────┐
│   datapoint   │
│ go-branded-id │
└────┬──────────┘
     │
┌────┴──────────┐
│   examples    │
│ go-branded-id │
└───────────────┘
```

## Replace Directive Strategy

Sub-modules use `replace` directives to resolve the root module and unpublished sibling modules locally:

| Module    | Replace Directives                                                                        |
| --------- | ----------------------------------------------------------------------------------------- |
| Root      | None                                                                                      |
| nanoid    | `go-composable-business-types => ../`                                                     |
| locale    | `go-composable-business-types => ../`                                                     |
| money     | `go-composable-business-types => ../`, `locale => ../locale`                              |
| datapoint | `go-composable-business-types => ../`, `nanoid => ../nanoid`                              |
| examples  | `go-composable-business-types => ../`, `nanoid => ../nanoid`, `datapoint => ../datapoint` |

**Why replace directives?** The root module v0.4.0 (published) still contains all packages including nanoid, locale, money, datapoint, and examples. Without `replace`, Go would see ambiguous imports — the same package provided by both the published root module and the local sub-module. The `replace` directive forces resolution to the local root module, which correctly excludes sub-module directories (they have their own go.mod).

**Why not pure go.work?** `go.work` with `use` directives works for building but `go mod tidy` in sub-modules fails due to ambiguous imports when resolving against the published root v0.4.0. The `replace` directives in each sub-module's go.mod ensure `go mod tidy` works both with and without the workspace.

## Package-to-Module Assignment

| Package     | Module        | Reason                                                |
| ----------- | ------------- | ----------------------------------------------------- |
| enums       | Root          | Leaf, zero deps                                       |
| validate    | Root          | Leaf, zero deps                                       |
| pkg/errors  | Root          | Leaf, zero deps — sentinel errors for entire library  |
| scanutil    | Root          | Leaf, depends only on pkg/errors                      |
| testutil    | Root          | Leaf, zero deps (testing only)                        |
| version     | Root          | Leaf, zero deps                                       |
| bounded     | Root          | Only depends on root packages                         |
| importance  | Root          | Only depends on root packages                         |
| tag         | Root          | Only depends on root packages                         |
| types       | Root          | Only depends on root packages                         |
| temporal    | Root          | Only depends on root packages (types)                 |
| actor       | Root          | Only depends on root packages (enums) + go-branded-id |
| projectcore | Root          | Only depends on root packages                         |
| nanoid      | **nanoid**    | Isolates sixafter/nanoid + crypto deps                |
| locale      | **locale**    | Isolates golang.org/x/text                            |
| money       | **money**     | Isolates bojanz/currency (heaviest dep)               |
| datapoint   | **datapoint** | Depends on nanoid module (can't be in root)           |
| examples/\* | **examples**  | Depends on multiple modules                           |

## Dependency Isolation After Split

| Consumer needs   | Dependencies after split                                |
| ---------------- | ------------------------------------------------------- |
| Just `enums`     | **go-branded-id only** (via root)                       |
| Just `nanoid`    | **sixafter/nanoid + go-branded-id**                     |
| Just `types`     | **go-branded-id only** (via root)                       |
| Just `locale`    | **golang.org/x/text + go-branded-id**                   |
| Just `money`     | **bojanz/currency + golang.org/x/text + go-branded-id** |
| Just `datapoint` | **everything** (expected — it composes everything)      |
