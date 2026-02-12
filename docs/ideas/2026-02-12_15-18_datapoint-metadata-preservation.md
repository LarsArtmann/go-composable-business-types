# DataPoint: Preserving All Metadata in Modern Applications

**Date:** 2026-02-12 15:18
**Status:** Idea / Design Proposal
**Context:** `go-composable-business-types`
**Inspired by:** "100 Things I hate in modern Software Development"

---

## Direct Pain Points Addressed

From "100 Things I hate in modern Software Development":

| # | Pain Point | How DataPoint Solves It |
|---|------------|------------------------|
| 1 | No Event Sourcing | DataPoint IS event-sourcing inspired - every change is a fact |
| 2 | Easy to do wrong thing | Strong types make invalid states unrepresentable |
| 25 | Metadata management | Core problem this design solves |
| 26 | Long unreadable UUIDs | Use `NanoId` - shorter, readable, same uniqueness |
| 27 | Errors not isolated/recovered | `Cause` captures full context for debugging |
| 36 | Not learning from mistakes | Full causal chain enables analysis |
| 41 | No event log on data access | Every DataPoint IS an auditable event |
| 44 | MEGA files (+1000 lines) | Implementation kept under 250 lines per file |

---

## The Problem

Modern applications lose enormous amounts of valuable metadata by default:

| What We Lose | Why It Matters |
|--------------|----------------|
| **Why** data was created | No causal chain, can't explain "how did we get here" |
| **Who** initiated changes | Lost intent, context, and accountability |
| **When** (both dimensions) | Only current timestamp, no "occurred vs recorded" |
| **Relationships** | Only explicit FKs, no semantic relationship metadata |
| **Context** | Session, environment, trace info discarded |
| **Intent** | Business reasoning not captured |

**Result:** Systems that can't explain themselves, can't be audited properly, and lose institutional knowledge.

---

## Proposed Solution: DataPoint

A general-purpose data type inspired by event sourcing that preserves **ALL** relationships and metadata at the application layer.

### Core Design

```go
// Cause captures WHY this data exists - the causal chain
type Cause[T comparable] struct {
    trigger   Trigger              // What initiated this
    reason    string               // Human/business reason
    parents   []Id[DataPoint[T]]   // Direct causal ancestors
    intent    string               // What the actor intended
}

type Trigger uint8
const (
    TriggerUserAction Trigger = iota
    TriggerScheduled
    TriggerDerived      // Computed from other data
    TriggerIntegration  // External system
    TriggerCorrection   // Fixing previous data
)

// Context captures the operational environment
type Context struct {
    correlation   Id[any]        // Cross-service trace
    session       Id[any]        // User session
    environment   string         // prod/staging/etc
    service       string         // Originating service
    trace         []string       // Full trace path
}

// Bitemporal: TWO time dimensions (critical for audit/corrections)
type Bitemporal struct {
    occurred    Timestamp  // When the FACT happened in reality
    recorded    Timestamp  // When we LEARNED about it
    validFrom   Timestamp  // When it becomes valid
    validUntil  *Timestamp // When it stops being valid (nil = forever)
}

// DataPoint: The complete unit of truth
type DataPoint[T comparable] struct {
    id          Id[DataPoint[T]]
    payload     T                    // The actual data
    actor       Actor[T]             // Who created it
    cause       Cause[T]             // Why it exists
    context     Context              // Where/how it was created
    temporal    Bitemporal           // When (both dimensions)
    references  []Reference[T]       // All relationships
    tags        map[string]string    // Arbitrary metadata
    version     uint64               // For optimistic concurrency
}

// Reference: Typed relationship to other entities
type Reference[T comparable] struct {
    kind    string     // "parent", "depends_on", "related", etc.
    target  Id[any]    // The referenced entity
    meta    string     // Additional context about the relationship
}
```

---

## Key Design Decisions

### 1. Bitemporal (Not Single Timestamp)

```
occurred  ─── When it happened in the real world
recorded  ─── When your system learned about it
validFrom ─── When it becomes valid
validUntil ── When it stops being valid (nil = forever)
```

**Enables:**
- Backdated corrections
- Late-arriving data handling
- Point-in-time queries ("what did we believe on date X?")
- Proper audit trails

### 2. Causal Chain (Not Just "Created By")

```
parents ─── Which data points directly caused this one
trigger ─── What type of event initiated this
reason  ─── Human-readable business reason
intent  ─── What the actor was trying to accomplish
```

**Enables:**
- Impact analysis ("what will break if I change this?")
- Root cause debugging
- Compliance ("explain this number")
- Process understanding

### 3. Rich References (Not Just FK Columns)

```go
Reference{
    kind:   "fulfillment_center",
    target: centerId,
    meta:   "origin",  // Why this relationship exists
}
```

**Enables:**
- Graph queries
- Semantic relationship understanding
- Impact analysis
- No more "why is this related to that?"

### 4. Full Context Capture

```
correlation ─── Cross-service request tracing
session     ─── User session context
environment ─── prod/staging/dev
service     ─── Which service created this
trace       ─── Full distributed trace path
```

**Enables:**
- Cross-system debugging
- Environment-aware queries
- Compliance and audit

---

## Usage Example

```go
// A user updates an order status
orderStatus := DataPoint[OrderStatus]{
    id:      NewId[DataPoint[OrderStatus]](uuid.New()),
    payload: OrderStatus{OrderID: orderId, Status: "shipped"},
    actor:   NewActorUser(userId),
    cause: Cause[OrderStatus]{
        trigger: TriggerUserAction,
        reason:  "Customer requested expedited shipping",
        parents: []Id[DataPoint[OrderStatus]]{previousStatus.Id()},
        intent:  "Fulfill order",
    },
    context: Context{
        correlation: correlationId,
        session:     sessionId,
        service:     "order-service",
    },
    temporal: Bitemporal{
        occurred:  Now(),
        recorded:  Now(),
        validFrom: Now(),
    },
    references: []Reference[OrderStatus]{
        {kind: "order", target: orderId},
        {kind: "fulfillment_center", target: centerId, meta: "origin"},
    },
}
```

---

## Comparison: Before vs After

### Traditional Approach

```go
type OrderStatus struct {
    OrderID string
    Status  string
    // ... maybe a CreatedAt and UpdatedAt
    // That's it. Everything else is lost.
}
```

### With DataPoint

```go
// ALL context preserved:
// - Who changed it (actor)
// - Why they changed it (cause)
// - When it happened AND when we knew (bitemporal)
// - What else was involved (references)
// - Full operational context (context)
```

---

## Implementation Variations

### 1. Minimal Version
- Just `id`, `payload`, `actor`, `occurred`, `recorded`
- Good starting point, can evolve

### 2. With Versioning
- Built-in event schema versioning
- Enables safe schema evolution
- Backward-compatible deserialization

### 3. With Aggregation
- `DerivedDataPoint` type that references source data
- Computations preserve their inputs
- Reproducible results

---

## Design Constraints (Non-Negotiable)

From project standards:

- **File size:** All implementation files < 250 lines
- **Package structure:** Clear module boundaries, no mixing Kernel/Plugin with Domain
- **Test setup:** Use native Go test or onsi/ginkgo
- **Id type:** Prefer `NanoId` over UUID (shorter, more readable)
- **Error handling:** Isolated, recoverable, well-communicated to users

---

## Open Questions

1. **Storage efficiency** - How to store without exploding size?
   - Consider: Columnar storage, reference compression, archival strategies

2. **Query patterns** - How to query causal chains efficiently?
   - Consider: Graph databases, materialized paths, closure tables

3. **Privacy/GDPR** - Some metadata might need deletion
   - Consider: Selective redaction, anonymization patterns

4. **Schema evolution** - How to handle `payload` type changes over time?
   - Consider: Versioned schemas, migration functions

5. **Deterministic Simulation Testing** (#42) - Can we replay causal chains for testing?
   - Consider: Capture enough context to enable full replay

---

## Implementation Plan

**File structure (each < 250 lines):**

```
cbt/
├── datapoint.go        # Core DataPoint[T] type
├── datapoint_cause.go  # Cause, Trigger types
├── datapoint_context.go # Context, trace info
├── datapoint_temporal.go # Bitemporal handling
├── datapoint_ref.go    # Reference types
├── nanoid.go           # NanoId implementation (not UUID!)
└── datapoint_test.go   # Tests using native go test
```

**Next Steps:**

- [ ] Implement `nanoid.go` - prefer over UUID
- [ ] Implement `datapoint_temporal.go` - Bitemporal first (foundation)
- [ ] Implement `datapoint_cause.go` - Causal chain
- [ ] Implement `datapoint_context.go` - Operational context
- [ ] Implement `datapoint_ref.go` - Rich references
- [ ] Implement `datapoint.go` - Core type composing all above
- [ ] Add native Go tests
- [ ] Add JSON serialization (preserving all metadata)
- [ ] Explore storage adapters (sqlite for small, swap to ScyllaDB for scale)

---

## References

- Event Sourcing patterns
- Bitemporal data modeling
- Martin Fowler: "EventStorming"
- VAT (Value Added Tax) bitemporal requirements
- Compliance audit trail requirements
