# Component Analysis: go-composable-business-types

> Analysis of components that could be split into standalone reusable libraries/SDKs.
> Last Updated: 2026-02-28

---

## Executive Summary

This project contains **14 distinct components** that fall into three categories:

1. **Keep as-is** - Already thin wrappers or no better alternatives exist
2. **Consider extracting** - Valuable standalone abstractions
3. **Already optimal** - Using best-in-class libraries directly

---

## Component Analysis

### 1. Branded ID (`ID[B, V]`)

**Location:** `id.go`

**Description:** Phantom/branded type pattern for type-safe identifiers. Prevents mixing `UserID` with `OrderID` at compile time.

```go
type ID[B any, V comparable] struct{ value V }
type UserID = ID[UserBrand, string]
type OrderID = ID[OrderBrand, int64]
```

**Alternatives:**

| Library                  | URL                            | Features                  | Trade-offs                  |
| ------------------------ | ------------------------------ | ------------------------- | --------------------------- |
| Custom phantom types     | N/A                            | Simple, zero deps         | Must implement yourself     |
| `github.com/google/uuid` | https://github.com/google/uuid | RFC 4122 UUIDs            | No branding, just UUID type |
| `github.com/oklog/ulid`  | https://github.com/oklog/ulid  | Sortable IDs              | No branding, just ULID type |
| `github.com/rs/xid`      | https://github.com/rs/xid      | Globally unique, sortable | No branding                 |

**Verdict: KEEP AS-IS**

- No popular library provides phantom branding in Go
- The pattern is simple (one generic struct + ~150 lines for serialization)
- Value comes from the _pattern_, not from complex implementation

**Potential Standalone Library?** **Maybe** - Could be published as `github.com/larsartmann/go-branded-id` but limited value since it's so simple.

---

### 2. NanoId

**Location:** `nanoid.go`

**Description:** URL-safe, cryptographically random identifier (default 21 chars = 126 bits entropy).

**Current Dependency:** `github.com/sixafter/nanoid` (FIPS-140 compatible, high-performance)

**Alternatives:**

| Library                 | Stars | Features                      | Trade-offs      |
| ----------------------- | ----- | ----------------------------- | --------------- |
| `sixafter/nanoid`       | ~100  | FIPS-140, high-perf, URL-safe | Current choice  |
| `matoous/go-nanoid`     | ~800  | Simple, popular               | Not FIPS-140    |
| `aidarkhanov/nanoid-go` | ~150  | Minimal                       | Less maintained |

**Verdict: KEEP AS-IS**

- Already using best-in-class library
- The wrapper adds: validation (8-256 char range), SQL Scanner/Valuer, JSON marshaling
- Wrapper value: ~50 lines of validation + interfaces

**Potential Standalone Library?** **No** - Too thin, just validation around existing lib.

---

### 3. Bitemporal

**Location:** `datapoint_temporal.go`

**Description:** Bitemporal tracking with `validFrom`, `validUntil`, `recorded`, and `correction` flag. Enables point-in-time queries and historical corrections.

```go
type Bitemporal struct {
    validFrom  Timestamp  // When fact became true in real world
    validUntil Timestamp  // When fact ceased to be true (zero = still valid)
    recorded   Timestamp  // When we recorded it in system
    correction bool       // Is this a correction?
}
```

**Alternatives:**

| Library                                     | Features               | Trade-offs                                   |
| ------------------------------------------- | ---------------------- | -------------------------------------------- |
| **None exist for Go**                       | -                      | -                                            |
| `looplab/eventhorizon`                      | Full CQRS/ES framework | Overkill, infrastructure required            |
| `ThreeDotsLabs/watermill`                   | Event streaming        | Different use case                           |
| Database-level (PostgreSQL temporal_tables) | SQL AS OF queries      | DB-specific, not portable                    |
| `temporalio/sdk-go`                         | Workflow engine        | Full workflow engine, not just temporal data |

**Verdict: KEEP AS-IS - UNIQUE VALUE**

- **No Go-native bitemporal library exists**
- This is genuinely novel in the Go ecosystem
- Simple, portable, database-agnostic

**Potential Standalone Library?** **YES** - Strong candidate for extraction. Could be `github.com/larsartmann/go-bitemporal` with:

- Core `Bitemporal` struct
- SQL query helpers for AS OF queries
- Repository patterns for common temporal queries
- Optional integrations (GORM, sqlc)

---

### 4. DataPoint[T]

**Location:** `datapoint.go`, `datapoint_*.go`

**Description:** Self-contained data unit with complete audit trail. Combines:

- NanoId for uniqueness
- Bitemporal for time tracking
- ActorEntry for who caused it
- Trigger for what caused it
- Context for execution environment
- References to related entities
- Causal chain (Cause[T])
- Tags for metadata
- Version for optimistic concurrency

**Alternatives:**

| Library                   | Features        | Trade-offs                     |
| ------------------------- | --------------- | ------------------------------ |
| `looplab/eventhorizon`    | Full CQRS/ES    | Heavy, requires infrastructure |
| `ThreeDotsLabs/watermill` | Event streaming | Different abstraction level    |
| Custom event sourcing     | Flexible        | Build from scratch             |

**Verdict: KEEP AS-IS - UNIQUE VALUE**

- Combines multiple patterns into a cohesive abstraction
- "Event sourcing light" without infrastructure requirements
- Novel composition of: bitemporal + actor tracking + causality

**Potential Standalone Library?** **YES** - But only as part of a larger "event audit" package. Consider combining with Bitemporal.

---

### 5. ActorChain[T] / ActorEntry[T]

**Location:** `actor.go`

**Description:** Ordered chain of actors for audit trails. Tracks progression: User → API Gateway → Order Service → Database.

```go
type ActorEntry[T comparable] struct {
    Kind ActorKind  // User, Bot, System, Service
    Id   ID[struct{}, T]
    Name string     // Optional human-readable name
}
type ActorChain[T] []ActorEntry[T]
```

**Alternatives:**

| Library              | Features            | Trade-offs                               |
| -------------------- | ------------------- | ---------------------------------------- |
| OpenTelemetry traces | Distributed tracing | Infrastructure required, different scope |
| Custom audit logs    | Flexible            | No structured chain                      |
| **No equivalent**    | -                   | -                                        |

**Verdict: KEEP AS-IS - UNIQUE VALUE**

- Simple but powerful pattern for service-to-service call chains
- Complements OpenTelemetry (not a replacement)
- Novel in Go ecosystem

**Potential Standalone Library?** **Maybe** - Could be extracted but works best with the rest of the types.

---

### 6. Context (Execution Context)

**Location:** `datapoint_context.go`

**Description:** Captures execution context: environment, session, request ID, source, tags.

```go
type Context struct {
    environment string            // "production", "staging", etc.
    session     string            // Session/correlation ID
    request     string            // Request ID for tracing
    source      string            // Source system/service
    tags        map[string]string // Additional metadata
}
```

**Alternatives:**

| Approach                  | Features         | Trade-offs                   |
| ------------------------- | ---------------- | ---------------------------- |
| `context.Context` values  | Standard library | Unstructured, no type safety |
| OpenTelemetry baggage     | Distributed      | Infrastructure required      |
| Structured logging fields | Common           | Tied to logging              |

**Verdict: KEEP AS-IS**

- Type-safe, serializable alternative to `context.Context` values
- Useful for persisting context with data

**Potential Standalone Library?** **No** - Too simple, better as part of DataPoint.

---

### 7. Reference[T] / Cause[T]

**Location:** `datapoint_ref.go`, `datapoint_cause.go`

**Description:**

- `Reference[T]`: Type-safe reference to another entity with relationship metadata
- `Cause[T]`: Causal chain tracking for building audit/lineage graphs

**Alternatives:**

| Library              | Features              | Trade-offs                |
| -------------------- | --------------------- | ------------------------- |
| Foreign keys in DB   | Referential integrity | DB-specific, not portable |
| GraphQL connections  | Graph traversal       | Different domain          |
| **No Go equivalent** | -                     | -                         |

**Verdict: KEEP AS-IS - UNIQUE VALUE**

- Simple but powerful patterns for entity relationships and causality
- Novel composition for Go

**Potential Standalone Library?** **No** - Works best as part of DataPoint ecosystem.

---

### 8. BoundedString

**Location:** `bounded.go`

**Description:** String with validated length constraints.

```go
func NewBoundedString(minLen, maxLen int, value string) (BoundedString, error)
var NewProductName = cbt.BoundedStringOf(1, 200) // Factory pattern
```

**Alternatives:**

| Library                              | Features              | Trade-offs             |
| ------------------------------------ | --------------------- | ---------------------- |
| `github.com/go-playground/validator` | Struct tag validation | Runtime, not type-safe |
| `github.com/asaskevich/govalidator`  | String validators     | Not type-safe          |
| Custom types                         | Type-safe             | Build yourself         |

**Verdict: KEEP AS-IS**

- Type-safe at construction time
- Factory pattern (`BoundedStringOf`) is elegant for domain types
- SQL Scanner/Valuer built-in

**Potential Standalone Library?** **Maybe** - Could be `github.com/larsartmann/go-bounded-types` with BoundedString, BoundedInt, etc.

---

### 9. Email

**Location:** `common.go`

**Description:** Validated email address with RFC 5322 parsing.

**Alternatives:**

| Library                              | Features         | Trade-offs            |
| ------------------------------------ | ---------------- | --------------------- |
| `net/mail` (stdlib)                  | RFC 5322 parsing | No validation wrapper |
| `github.com/go-playground/validator` | Tag-based        | Runtime validation    |
| `github.com/asaskevich/govalidator`  | `IsEmail()`      | Not type-safe         |

**Verdict: KEEP AS-IS**

- Thin wrapper around `net/mail` with type safety
- Adds: `LocalPart()`, `Domain()`, `Normalize()`, SQL interfaces
- ~80 lines, mostly interface implementations

**Potential Standalone Library?** **No** - Too thin.

---

### 10. URL

**Location:** `common.go`

**Description:** Validated URL with http/https scheme requirement.

**Alternatives:**

| Library            | Features    | Trade-offs            |
| ------------------ | ----------- | --------------------- |
| `net/url` (stdlib) | URL parsing | No validation wrapper |
| Custom validation  | Flexible    | Build yourself        |

**Verdict: KEEP AS-IS**

- Thin wrapper with http/https restriction
- ~60 lines

**Potential Standalone Library?** **No** - Too thin.

---

### 11. Money / Currency

**Location:** `money.go`

**Description:** Wrapper around `github.com/bojanz/currency` for ISO 4217 money handling.

**Current Dependency:** `github.com/bojanz/currency` (370+ locales, full ISO 4217)

**Alternatives:**

| Library              | Stars | Features                           | Trade-offs                         |
| -------------------- | ----- | ---------------------------------- | ---------------------------------- |
| `bojanz/currency`    | ~400  | ISO 4217, 370+ locales, formatting | **Current choice**                 |
| `Rhymond/go-money`   | ~1.6k | Money arithmetic, currency         | Different API, less locale support |
| `shopspring/decimal` | ~6k   | High-precision decimal             | Not currency-specific              |

**Verdict: KEEP AS-IS**

- Already using best-in-class for internationalization
- The wrapper adds: convenience functions, Locale integration
- Consider using directly instead of wrapper

**Potential Standalone Library?** **No** - Already a thin wrapper around excellent library.

---

### 12. Cents

**Location:** `common.go`

**Description:** Monetary amount in smallest currency unit (prevents float errors).

```go
type Cents int64
func NewCents(v int64) Cents
func (c Cents) Float64() float64 // Returns dollars (e.g., 1099 → 10.99)
```

**Alternatives:**

| Library              | Features       | Trade-offs                 |
| -------------------- | -------------- | -------------------------- |
| `shopspring/decimal` | High precision | Heavier, more complex      |
| `bojanz/currency`    | Full currency  | Heavier for simple cases   |
| `int64` directly     | Zero deps      | No type safety, no methods |

**Verdict: KEEP AS-IS**

- Simple, zero-dependency for fixed-point arithmetic
- Good for USD/EUR where 2 decimal places suffice
- Complements Money for different use cases

**Potential Standalone Library?** **No** - Too simple.

---

### 13. Percentage

**Location:** `common.go`

**Description:** 0-100 value with clamping and float conversion.

```go
type Percentage uint8
func NewPercentage(v uint8) Percentage // Clamps >100 to 100
```

**Verdict: KEEP AS-IS**

- Very simple (~30 lines)
- Useful for UI/business logic

**Potential Standalone Library?** **No** - Too simple.

---

### 14. Timestamp / Duration

**Location:** `common.go`

**Description:** Domain wrappers around `time.Time` and `time.Duration`.

```go
type Timestamp struct{ time.Time }
type Duration struct{ time.Duration }
```

**Verdict: KEEP AS-IS**

- Provides domain clarity and consistent `IsZero()` pattern
- SQL interfaces built-in

**Potential Standalone Library?** **No** - Too simple.

---

### 15. Enums (ActorKind, Priority, Status, Trigger)

**Location:** `enum.go`, `enum_enum.go` (generated)

**Description:** Generated enums using `github.com/abice/go-enum`.

**Current Dependency:** `github.com/abice/go-enum` (generates Marshal, Names, Values, MustParse, SQL)

**Alternatives:**

| Library             | Features            | Trade-offs         |
| ------------------- | ------------------- | ------------------ |
| `abice/go-enum`     | Code gen, SQL, JSON | **Current choice** |
| `alvaroloes/enumer` | Similar features    | Less maintained    |
| `go-enum/enum`      | Different approach  | Less popular       |
| Hand-written        | Full control        | Boilerplate        |

**Verdict: KEEP AS-IS**

- Already using good code generator
- Enums are domain-specific, not reusable

**Potential Standalone Library?** **No** - Domain-specific.

---

### 16. Locale

**Location:** `locale.go`

**Description:** BCP 47 language tag wrapper around `golang.org/x/text/language`.

**Current Dependency:** `golang.org/x/text/language`

**Verdict: KEEP AS-IS**

- Thin wrapper with SQL interfaces and convenience constants
- ~130 lines, mostly interface implementations

**Potential Standalone Library?** **No** - Too thin.

---

## Extraction Recommendations

### Strong Candidates for Standalone Libraries

| Component                               | Proposed Library | Rationale                                     |
| --------------------------------------- | ---------------- | --------------------------------------------- |
| **Bitemporal**                          | `go-bitemporal`  | No Go equivalent, genuinely useful standalone |
| **DataPoint + Bitemporal + ActorChain** | `go-event-audit` | Novel composition, "event sourcing light"     |

### Consider Combining

If extracting, consider a single library:

```
github.com/larsartmann/go-event-audit
├── bitemporal.go      # Bitemporal struct
├── actor.go           # ActorChain, ActorEntry
├── context.go         # Execution context
├── reference.go       # Reference[T]
├── cause.go           # Cause[T]
├── datapoint.go       # DataPoint[T] (combines all above)
└── README.md
```

**Benefits:**

- Coherent "audit trail" package
- Bitemporal works standalone
- DataPoint combines everything for full power

---

## Keep in Current Form

| Component               | Reason                            |
| ----------------------- | --------------------------------- |
| `ID[B, V]`              | Simple pattern, no library value  |
| `NanoId`                | Thin wrapper around excellent lib |
| `BoundedString`         | Could extract but works well here |
| `Email`, `URL`          | Thin wrappers                     |
| `Money`                 | Already wraps best-in-class       |
| `Cents`, `Percentage`   | Too simple                        |
| `Timestamp`, `Duration` | Too simple                        |
| `Locale`                | Thin wrapper                      |
| Enums                   | Domain-specific                   |

---

## Summary Matrix

| Component               | Lines | Unique Value | Extract? | Reason                       |
| ----------------------- | ----- | ------------ | -------- | ---------------------------- |
| `ID[B,V]`               | ~160  | Medium       | Maybe    | Pattern is simple but useful |
| `NanoId`                | ~150  | Low          | No       | Thin wrapper                 |
| **`Bitemporal`**        | ~120  | **High**     | **Yes**  | No Go equivalent             |
| **`DataPoint[T]`**      | ~290  | **High**     | **Yes**  | Novel composition            |
| **`ActorChain`**        | ~75   | **High**     | **Yes**  | Novel pattern                |
| `Context`               | ~140  | Medium       | No       | Works with DataPoint         |
| `Reference`, `Cause`    | ~240  | Medium       | No       | Works with DataPoint         |
| `BoundedString`         | ~130  | Medium       | Maybe    | Useful standalone            |
| `Email`, `URL`          | ~180  | Low          | No       | Thin wrappers                |
| `Money`                 | ~70   | Low          | No       | Already wraps lib            |
| `Cents`, `Percentage`   | ~80   | Low          | No       | Too simple                   |
| `Timestamp`, `Duration` | ~100  | Low          | No       | Too simple                   |
| `Locale`                | ~130  | Low          | No       | Thin wrapper                 |
| Enums                   | ~30   | Low          | No       | Domain-specific              |

---

## Final Recommendation

1. **Keep the library as-is** - The composition is valuable
2. **If extraction is desired**, create a single focused library:
   - `github.com/larsartmann/go-bitemporal` - Just the Bitemporal struct + helpers
3. **Document the patterns** - The real value is in the _design patterns_, not the code itself

The current monolithic approach works well because:

- All types work together (DataPoint depends on Bitemporal, ActorEntry, etc.)
- Single import for all business types
- Consistent patterns (`IsZero()`, `With*` methods, SQL interfaces)
- Zero external dependencies beyond well-chosen libs

---

_Generated: 2026-02-28_
