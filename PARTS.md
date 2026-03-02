# Component Analysis: go-composable-business-types

> Analysis of components that could be split into standalone reusable libraries/SDKs.
> Last Updated: 2026-03-01

---

## Executive Summary

This project contains **9 packages** with **14 distinct types** that fall into three categories:

1. **Keep as-is** - Already thin wrappers or no better alternatives exist
2. **Consider extracting** - Valuable standalone abstractions with unique value
3. **Already optimal** - Using best-in-class libraries directly (per HOW_TO_GOLANG.md)

### Quick Reference

| Package | Types | Unique Value | Extract? | Recommendation |
|---------|-------|--------------|----------|----------------|
| `id/` | `ID[B,V]` | Medium | Maybe | Pattern is simple but useful |
| `nanoid/` | `NanoId` | Low | No | Thin wrapper around `sixafter/nanoid` |
| `temporal/` | `Bitemporal` | **High** | **Yes** | **No Go equivalent exists** |
| `actor/` | `ActorChain[T]`, `ActorEntry[T]` | **High** | **Yes** | Novel audit trail pattern |
| `bounded/` | `BoundedString` | Medium | Maybe | Useful standalone, consider extending |
| `types/` | `Email`, `URL`, `Cents`, `Percentage`, `Timestamp`, `Duration` | Low | No | Thin wrappers, keep together |
| `money/` | `Money` | Low | No | Direct alias to `bojanz/currency` |
| `locale/` | `Locale` | Low | No | Thin wrapper around `x/text/language` |
| `enums/` | `ActorKind`, `Priority`, `Status`, `Trigger` | Low | No | Domain-specific, generated |

---

## Detailed Component Analysis

### 1. Branded ID (`id/`)

**Location:** `id/id.go` (~167 lines)

**Description:** Phantom/branded type pattern for type-safe identifiers. Prevents mixing `UserID` with `OrderID` at compile time.

```go
type ID[B any, V comparable] struct{ value V }
type UserID = ID[UserBrand, string]
type OrderID = ID[OrderBrand, int64]
```

**Features:**
- Zero-cost abstraction at runtime
- Full JSON serialization (null for zero, string for values)
- SQL Scanner/Valuer for database integration
- Text marshaling for XML/TOML
- Works with both `string` and `int64` underlying types

**Alternatives:**

| Library | Stars | Features | Trade-offs |
|---------|-------|----------|------------|
| Custom phantom types | N/A | Simple, zero deps | Must implement yourself |
| `google/uuid` | 5k+ | RFC 4122 UUIDs | No branding, just UUID type |
| `oklog/ulid` | 3k+ | Sortable IDs | No branding, just ULID type |
| `rs/xid` | 4k+ | Globally unique, sortable | No branding |
| `gofrs/uuid` | 1k+ | UUID with more features | No branding |

**Verdict: KEEP AS-IS**

- No popular library provides phantom branding in Go
- The pattern is simple but the serialization interfaces add value (~100 lines)
- Value comes from the _pattern_, not complex implementation

**Potential Standalone Library?** **Maybe** - Could be `github.com/larsartmann/go-branded-id` but limited standalone value.

---

### 2. NanoId (`nanoid/`)

**Location:** `nanoid/nanoid.go` (~151 lines)

**Description:** URL-safe, cryptographically random identifier with validation.

**Current Dependency:** `github.com/sixafter/nanoid` (FIPS-140 compatible, high-performance)

**Features:**
- Validation (8-256 character range)
- URL-safe alphabet enforcement
- SQL Scanner/Valuer
- JSON text marshaling
- `MustParse*` variants for constants

**Alternatives:**

| Library | Stars | Features | Trade-offs |
|---------|-------|----------|------------|
| `sixafter/nanoid` | ~100 | FIPS-140, high-perf | **Current choice** |
| `matoous/go-nanoid` | ~800 | Simple, popular | Not FIPS-140 |
| `aidarkhanov/nanoid-go` | ~150 | Minimal | Less maintained |

**Verdict: KEEP AS-IS**

- Already using best-in-class library per requirements
- The wrapper adds: validation, SQL interfaces, JSON marshaling
- Wrapper value: ~80 lines of validation + interfaces

**Potential Standalone Library?** **No** - Too thin, just validation around existing lib.

---

### 3. Bitemporal (`temporal/`)

**Location:** `temporal/temporal.go` (~121 lines)

**Description:** Bitemporal tracking with `validFrom`, `validUntil`, `recorded`, and `correction` flag. Enables point-in-time queries and historical corrections.

```go
type Bitemporal struct {
    validFrom  Timestamp  // When fact became true in real world
    validUntil Timestamp  // When fact ceased to be true (zero = still valid)
    recorded   Timestamp  // When we recorded it in system
    correction bool       // Is this a correction?
}
```

**Features:**
- Point-in-time validity checking (`IsValidAt`, `IsCurrentlyValid`)
- Immutable `With*` methods for functional updates
- Full JSON serialization
- Correction tracking for audit trails

**Alternatives:**

| Library | Features | Trade-offs |
|---------|----------|------------|
| **None exist for Go** | - | - |
| `looplab/eventhorizon` | Full CQRS/ES framework | Overkill, requires infrastructure |
| `ThreeDotsLabs/watermill` | Event streaming | Different use case |
| PostgreSQL `temporal_tables` | SQL AS OF queries | DB-specific, not portable |
| `temporalio/sdk-go` | Workflow engine | Full workflow engine, not just temporal data |

**Verdict: KEEP AS-IS - UNIQUE VALUE**

- **No Go-native bitemporal library exists**
- This is genuinely novel in the Go ecosystem
- Simple, portable, database-agnostic
- Aligns with event sourcing patterns from HOW_TO_GOLANG.md

**Potential Standalone Library?** **YES - STRONG CANDIDATE**

Could be `github.com/larsartmann/go-bitemporal` with:
- Core `Bitemporal` struct
- SQL query helpers for AS OF queries
- Repository patterns for common temporal queries
- Optional integrations (sqlc helpers)

---

### 4. ActorChain / ActorEntry (`actor/`)

**Location:** `actor/actor.go` (~80 lines)

**Description:** Ordered chain of actors for audit trails. Tracks progression: User → API Gateway → Order Service → Database.

```go
type ActorEntry[T comparable] struct {
    Kind ActorKind  // User, Bot, System, Service
    Id   ID[struct{}, T]
    Name string     // Optional human-readable name
}
type ActorChain[T] []ActorEntry[T]
```

**Features:**
- Generic over ID type (string, int64, UUID, etc.)
- Origin/Current accessors
- Kind filtering (`ByKind`, `HasKind`)
- Constructor helpers (`UserActor`, `BotActor`, `SystemActor`, `ServiceActor`)

**Alternatives:**

| Library | Features | Trade-offs |
|---------|----------|------------|
| OpenTelemetry traces | Distributed tracing | Infrastructure required, different scope |
| Custom audit logs | Flexible | No structured chain |
| **No equivalent** | - | - |

**Verdict: KEEP AS-IS - UNIQUE VALUE**

- Simple but powerful pattern for service-to-service call chains
- Complements OpenTelemetry (not a replacement)
- Novel in Go ecosystem
- Aligns with Domain-Driven Design from HOW_TO_GOLANG.md

**Potential Standalone Library?** **Maybe** - Could be extracted but works best with the rest of the types (ID, enums).

---

### 5. BoundedString (`bounded/`)

**Location:** `bounded/bounded.go` (~132 lines)

**Description:** String with validated length constraints.

```go
func NewBoundedString(minLen, maxLen int, value string) (BoundedString, error)
var NewProductName = bounded.BoundedStringOf(1, 200) // Factory pattern
```

**Features:**
- Min/max length validation at construction
- Factory pattern (`BoundedStringOf`) for domain types
- Convenience constructors (`NonEmptyString`, `TrimmedBoundedString`)
- SQL Scanner/Valuer
- JSON marshaling

**Alternatives:**

| Library | Features | Trade-offs |
|---------|----------|------------|
| `go-playground/validator` | Struct tag validation | Runtime, not type-safe |
| `asaskevich/govalidator` | String validators | Not type-safe |
| Custom types | Type-safe | Build yourself |

**Verdict: KEEP AS-IS**

- Type-safe at construction time
- Factory pattern is elegant for domain types
- SQL interfaces built-in

**Potential Standalone Library?** **Maybe** - Could be `github.com/larsartmann/go-bounded-types` with:
- `BoundedString`
- `BoundedInt[T]` (min/max)
- `BoundedFloat[T]` (min/max)
- `NonEmptyString` as standalone

---

### 6. Common Types (`types/`)

**Location:** `types/types.go` (~486 lines)

**Types:**
- `Email` - RFC 5322 validated email
- `URL` - http/https validated URL
- `Percentage` - 0-100 clamped value
- `Cents` - Fixed-point monetary (smallest unit)
- `Timestamp` - Domain wrapper for `time.Time`
- `Duration` - Domain wrapper for `time.Duration`

**Features per type:**
- SQL Scanner/Valuer for all types
- JSON marshaling where appropriate
- Domain-specific methods (`Email.LocalPart()`, `Email.Domain()`, `Cents.Add()`, etc.)
- `IsZero()` pattern for null handling

**Alternatives:**

| Type | Library | Features | Trade-offs |
|------|---------|----------|------------|
| Email | `net/mail` (stdlib) | RFC 5322 parsing | No validation wrapper |
| Email | `go-playground/validator` | Tag-based | Runtime validation |
| URL | `net/url` (stdlib) | URL parsing | No validation wrapper |
| Money | `bojanz/currency` (used) | ISO 4217, locales | Heavier |
| Money | `shopspring/decimal` | High precision | Not currency-specific |
| Money | `Rhymond/go-money` | Money arithmetic | Different API |

**Verdict: KEEP AS-IS**

- Thin wrappers with type safety
- Consistent `IsZero()` pattern across all types
- SQL interfaces built-in
- `Cents` is particularly useful for fixed-point arithmetic without float errors

**Potential Standalone Library?** **No** - Too thin, better as part of larger collection.

---

### 7. Money (`money/`)

**Location:** `money/money.go` (~69 lines)

**Description:** Direct type alias to `github.com/bojanz/currency`.

**Current Dependency:** `github.com/bojanz/currency` (370+ locales, full ISO 4217)

```go
type Money = currency.Amount  // Direct alias
```

**Features:**
- Convenience constructors (`NewMoney`, `NewMoneyFromCents`)
- Currency validation helpers
- Locale-aware formatting
- Integration with `Locale` type

**Alternatives:**

| Library | Stars | Features | Trade-offs |
|---------|-------|----------|------------|
| `bojanz/currency` | ~400 | ISO 4217, 370+ locales | **Current choice** |
| `Rhymond/go-money` | ~1.6k | Money arithmetic | Different API, less locale support |
| `shopspring/decimal` | ~6k | High-precision decimal | Not currency-specific |

**Verdict: KEEP AS-IS**

- Already using best-in-class per HOW_TO_GOLANG.md philosophy
- The wrapper adds convenience functions and Locale integration
- Consider removing wrapper and using directly for simplicity

**Potential Standalone Library?** **No** - Already a thin wrapper around excellent library.

---

### 8. Locale (`locale/`)

**Location:** `locale/locale.go` (~128 lines)

**Description:** BCP 47 language tag wrapper around `golang.org/x/text/language`.

**Features:**
- Common locale constants (`LocaleEnUS`, `LocaleDeDE`, etc.)
- Hyphen and underscore format support
- Base language and region extraction
- SQL Scanner/Valuer
- JSON text marshaling

**Verdict: KEEP AS-IS**

- Thin wrapper with SQL interfaces and convenience constants
- Integrates well with Money formatting

**Potential Standalone Library?** **No** - Too thin.

---

### 9. Enums (`enums/`)

**Location:** `enums/enums.go` (~29 lines + generated)

**Types:**
- `ActorKind` - User, Bot, System, Service
- `Priority` - Low, Medium, High, Critical
- `Status` - Draft, Active, Paused, Archived, Deleted
- `Trigger` - Manual, Scheduled, Webhook, Import, Migration, System, Correction

**Current Dependency:** `github.com/abice/go-enum` (generates Marshal, Names, Values, MustParse, SQL)

**Verdict: KEEP AS-IS**

- Already using good code generator
- Enums are domain-specific, not reusable
- Generated code provides all needed functionality

**Potential Standalone Library?** **No** - Domain-specific.

---

## Extraction Recommendations

### Strong Candidates for Standalone Libraries

| Component | Proposed Library | Rationale | Lines |
|-----------|------------------|-----------|-------|
| **Bitemporal** | `go-bitemporal` | No Go equivalent, genuinely useful standalone | ~120 |
| **ActorChain** | `go-actor-chain` | Novel pattern for audit trails | ~80 |

### Consider Combining

If extracting, a single focused library may be better:

```
github.com/larsartmann/go-event-audit
├── bitemporal.go      # Bitemporal struct + helpers
├── actor.go           # ActorChain, ActorEntry
├── enums.go           # ActorKind, Trigger (shared)
└── README.md
```

**Benefits:**
- Coherent "audit trail" package
- Bitemporal works standalone
- ActorChain provides call chain tracking
- Combined: full audit capability

---

## Keep in Current Form

| Component | Reason |
|-----------|--------|
| `ID[B,V]` | Simple pattern, no library value |
| `NanoId` | Thin wrapper around excellent lib |
| `BoundedString` | Could extract but works well here |
| `Email`, `URL` | Thin wrappers |
| `Money` | Direct alias to best-in-class |
| `Cents`, `Percentage` | Too simple |
| `Timestamp`, `Duration` | Too simple |
| `Locale` | Thin wrapper |
| Enums | Domain-specific |

---

## Alignment with HOW_TO_GOLANG.md

This library follows the principles from HOW_TO_GOLANG.md:

| Principle | Status | Notes |
|-----------|--------|-------|
| Type Safety First | ✅ | All types prevent invalid states |
| Errors as Values | ✅ | Constructors return errors, `Must*` variants panic |
| Generated over Handwritten | ✅ | Uses `go-enum` for enum generation |
| Best-in-class dependencies | ✅ | `sixafter/nanoid`, `bojanz/currency`, `x/text/language` |
| No magic strings/numbers | ✅ | Constants and typed values throughout |
| Small, focused functions | ✅ | Most functions under 10 lines |

---

## Summary Matrix

| Package | Type | Lines | Unique Value | Extract? | Reason |
|---------|------|-------|--------------|----------|--------|
| `id/` | `ID[B,V]` | ~167 | Medium | Maybe | Pattern is simple but useful |
| `nanoid/` | `NanoId` | ~151 | Low | No | Thin wrapper |
| `temporal/` | `Bitemporal` | ~121 | **High** | **Yes** | **No Go equivalent** |
| `actor/` | `ActorChain` | ~80 | **High** | **Yes** | Novel pattern |
| `bounded/` | `BoundedString` | ~132 | Medium | Maybe | Useful standalone |
| `types/` | Email, URL, etc. | ~486 | Low | No | Thin wrappers |
| `money/` | `Money` | ~69 | Low | No | Direct alias |
| `locale/` | `Locale` | ~128 | Low | No | Thin wrapper |
| `enums/` | Various | ~29+gen | Low | No | Domain-specific |

---

## Final Recommendation

1. **Keep the library as-is** - The composition is valuable, selective imports work well
2. **If extraction is desired**, prioritize:
   - `go-bitemporal` - Standalone temporal data handling (highest value)
   - `go-actor-chain` - Audit trail tracking (medium value)
3. **Document the patterns** - The real value is in the _design patterns_, not just the code

The current modular approach works well because:
- All types work together (ActorChain uses ID and enums)
- Selective imports allow importing only what's needed
- Consistent patterns (`IsZero()`, `With*` methods, SQL interfaces)
- Zero external dependencies beyond well-chosen libs

---

_Generated: 2026-03-01_
