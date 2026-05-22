# Dependency Graph — go-composable-business-types

## Current (Single Module)

```
                    ┌─────────────────────────────────┐
                    │         ROOT MODULE              │
                    │  enums, validate, pkg/errors,    │
                    │  scanutil, testutil, version,    │
                    │  bounded, importance, tag,       │
                    │  types, temporal, actor,         │
                    │  projectcore, nanoid, locale,    │
                    │  money, datapoint, examples      │
                    └─────────────────────────────────┘
                    All packages in one go.mod
                    All deps pulled in for every consumer
```

## Proposed (6 Modules)

```
                        ┌──────────────┐
                        │    ROOT      │
                        │  13 packages │
                        │  go-branded-id │
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
│  (no ext deps)│
└────┬──────────┘
     │
┌────┴──────────┐
│   examples    │
│  (no ext deps)│
└───────────────┘
```

## Module Dependency Matrix

| Module ↓ depends on → | Root | nanoid | locale | money | datapoint | examples |
|---|---|---|---|---|---|---|
| **Root** | — | ✗ | ✗ | ✗ | ✗ | ✗ |
| **nanoid** | ✓ | — | ✗ | ✗ | ✗ | ✗ |
| **locale** | ✓ | ✗ | — | ✗ | ✗ | ✗ |
| **money** | ✓ | ✗ | ✓ | — | ✗ | ✗ |
| **datapoint** | ✓ | ✓ | ✗ | ✗ | — | ✗ |
| **examples** | ✓ | ✓ | ✗ | ✗ | ✓ | — |

## Package-to-Module Assignment

| Package | Current Module | Proposed Module | Reason |
|---|---|---|---|
| enums | root | Root | Leaf, zero deps |
| validate | root | Root | Leaf, zero deps |
| pkg/errors | root | Root | Leaf, zero deps |
| scanutil | root | Root | Leaf, zero deps |
| testutil | root | Root | Leaf, zero deps |
| version | root | Root | Leaf, zero deps |
| bounded | root | Root | Only depends on root packages |
| importance | root | Root | Only depends on root packages |
| tag | root | Root | Only depends on root packages |
| types | root | Root | Only depends on root packages |
| temporal | root | Root | Only depends on root packages (types) |
| actor | root | Root | Only depends on root packages (enums) + go-branded-id |
| projectcore | root | Root | Only depends on root packages |
| nanoid | root | **nanoid** | Isolates sixafter/nanoid + crypto deps |
| locale | root | **locale** | Isolates golang.org/x/text |
| money | root | **money** | Isolates bojanz/currency (heaviest dep) |
| datapoint | root | **datapoint** | Depends on nanoid module (can't be in root) |
| examples/* | root | **examples** | Depends on multiple modules |
