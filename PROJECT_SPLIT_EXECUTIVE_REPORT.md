# Project Split Analysis: go-composable-business-types

## Executive Summary

Go-composable-business-types is a cohesive library providing strongly-typed, composable business types including IDs, DataPoint with audit trails, money handling, validated strings, and enums. **Splitting is NOT RECOMMENDED** - the types are intentionally designed to compose together, with DataPoint as the centerpiece depending on most other types, making separation impractical.

## Current Architecture

Single `cbt` package with ~20 Go files organized by type. Flat structure with all types in one directory. Minimal dependencies (bojanz/currency, sixafter/nanoid, golang.org/x/text). Types include: branded IDs, NanoId, ActorChain, DataPoint[T], Bitemporal, Context, Reference, Cause, BoundedString, Email, URL, Percentage, Cents, Money, Timestamp, Duration, and generated enums (ActorKind, Priority, Status, Trigger, Locale).

## Split Recommendation: NOT RECOMMENDED

### Functional Areas Identified (but should remain unified)

| Area              | Types                                                        | Dependencies                 | Split Viability                  |
| ----------------- | ------------------------------------------------------------ | ---------------------------- | -------------------------------- |
| Identity          | ID[B,V], NanoId                                              | None                         | Low - core types used everywhere |
| Audit/Provenance  | DataPoint, ActorChain, Context, Reference, Cause, Bitemporal | NanoId, ID, Timestamp, enums | Very Low - highly interconnected |
| Financial         | Money, Cents, Percentage                                     | currency lib, Locale         | Medium - could standalone        |
| Validated Strings | BoundedString, Email, URL                                    | None                         | Medium - could standalone        |
| Time              | Timestamp, Duration                                          | None                         | Low - thin wrappers              |
| Enums             | ActorKind, Priority, Status, Trigger, Locale                 | x/text (Locale only)         | Low - used across types          |

### Why Splitting Would Harm the Project

- **DataPoint depends on**: NanoId, ActorEntry, Bitemporal, Context, Reference, Cause, Trigger - splitting creates circular dependencies
- **Intentional composability**: The library's value is in types working together seamlessly
- **Appropriate size**: ~1500-2000 lines across ~20 files is ideal for a single Go module
- **Simple import experience**: Single `import cbt` vs. managing multiple imports/versions
- **Testing complexity**: Integration between types would require multi-module test setups

### Minor Benefits of Splitting (outweighed by costs)

- Financial types could evolve independently for users only needing money handling
- Smaller dependency tree for users only needing validated strings
- Independent versioning per domain

## Implementation Path

N/A - Splitting not recommended.

## Conclusion

**NOT RECOMMENDED** (High Confidence). This is a well-designed, appropriately-sized library where cohesion is a feature, not a problem. The types are explicitly designed to compose (e.g., `DataPoint[T]` with Actor, Reference, Cause, Bitemporal). Splitting would fragment the user experience, create dependency hell, and undermine the library's core value proposition of providing composable business types that work together out of the box.
