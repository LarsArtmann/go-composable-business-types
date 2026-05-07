# go-composable-business-types

**Strongly typed, composable business types for Go — featuring DataPoint[T], a self-contained data unit with complete audit trail built-in.**

## What is this?

A Go library of type-safe base values designed for business applications. It combines:

- **Branded/phantom types** (`ID[B, V]`) — prevent mixing `UserID` with `OrderID` at compile time (via [`go-branded-id`](https://github.com/larsartmann/go-branded-id))
- **Audit trail primitives** — `ActorChain[T]`, `Context`, `Reference[T]`, `Cause[T]` for traceability
- **Bitemporal tracking** — separate valid time from recorded time
- **Domain primitives** — `BoundedString`, `Money`, `Email`, `Percentage`, `Cents`, `Timestamp`, `Duration`
- **Self-contained data units** — `DataPoint[T]` wraps any payload with full metadata (who, when, why, what caused it, references, tags)

Selective imports let you use only what you need.

## Installation

```bash
go get github.com/larsartmann/go-composable-business-types

# For branded IDs (used by actor, datapoint, etc.)
go get github.com/larsartmann/go-branded-id
```

## Types

| Type            | Package            | Purpose                                                              |
| --------------- | ------------------ | -------------------------------------------------------------------- |
| `NanoID`        | `nanoid/`          | URL-safe, cryptographically random ID (default 21 chars)             |
| `ActorChain[T]` | `actor/`           | Ordered chain of actors (User → Service → Service) for audit trails  |
| `DataPoint[T]`  | `datapoint/`       | Self-contained data unit with complete audit trail                   |
| `Bitemporal`    | `temporal/`        | Bitemporal tracking (validFrom, validUntil, recorded)                |
| `Context`       | `datapoint/`       | Execution context (environment, session, request, source)            |
| `Reference[T]`  | `datapoint/`       | Type-safe reference to another entity with relationship metadata     |
| `Cause[T]`      | `datapoint/`       | Causal chain tracking for building audit/lineage graphs              |
| `BoundedString` | `bounded/`         | String with validated length constraints                             |
| `Money`         | `money/`           | ISO 4217 currency via `github.com/bojanz/currency`                   |
| `Locale`        | `locale/`          | BCP 47 language tag for internationalization                         |
| `Email`         | `types/`           | Email address string                                                 |
| `URL`           | `types/`           | URL string                                                           |
| `Percentage`    | `types/`           | 0-100 value with float conversion (clamps overflow to 100)           |
| `Cents`         | `types/`           | Monetary amount in smallest unit (no float errors)                   |
| `Timestamp`     | `types/`           | Domain-wrapped time.Time                                             |
| `Duration`      | `types/`           | Domain-wrapped time.Duration                                         |
| `Importance`    | `importance/`      | Priority classification (0-100) with named levels                    |
| `Tag`           | `tag/`             | Validated string label with alphanumeric+hyphen constraint           |

### External Types (separate module)

| Type       | Module                                                                  | Purpose                                              |
| ---------- | ----------------------------------------------------------------------- | ---------------------------------------------------- |
| `ID[B, V]` | [`github.com/larsartmann/go-branded-id`](https://github.com/larsartmann/go-branded-id) | Branded, type-safe identifier — prevents mixing IDs |

## Enums (generated)

| Enum        | Values                                                            |
| ----------- | ----------------------------------------------------------------- |
| `ActorKind` | User, Bot, System, Service                                        |
| `Priority`  | Low, Medium, High, Critical                                       |
| `Status`    | Draft, Active, Paused, Archived, Deleted                          |
| `Trigger`   | Manual, Scheduled, Webhook, Import, Migration, System, Correction |
| `CauseKind` | Direct, Command, Event                                            |

## Usage

```go
// Selective imports - import only what you need
import (
    id "github.com/larsartmann/go-branded-id"
    "github.com/larsartmann/go-composable-business-types/actor"
    "github.com/larsartmann/go-composable-business-types/bounded"
    "github.com/larsartmann/go-composable-business-types/money"
    "github.com/larsartmann/go-composable-business-types/types"
)

// Branded IDs - can't mix UserId with OrderId at compile time
type UserBrand struct{}
type OrderBrand struct{}
type UserID = id.ID[UserBrand, string]
type OrderID = id.ID[OrderBrand, int64]

userId := id.NewID[UserBrand, string]("user-123")
orderId := id.NewID[OrderBrand, int64](42)

// Unbranded ID (when you don't need type separation)
type SessionID = id.ID[struct{}, string]
sessionId := id.NewID[struct{}, string]("sess-abc")

// Actor chain for audit trails and authorization
chain := actor.NewActorChain(actor.UserActor(id.NewID[struct{}, string]("user-1"), "Alice")).
    Append(actor.ServiceActor(id.NewID[struct{}, string]("api-gateway"), "API Gateway")).
    Append(actor.ServiceActor(id.NewID[struct{}, string]("order-svc"), "Order Service"))

origin := chain.Origin()   // User Alice
current := chain.Current() // Order Service

// BoundedString - validated string lengths
name, err := bounded.NewBoundedString(1, 100, "John Doe")
if err != nil {
    // handle validation error
}

// Factory for domain-specific bounded strings
var NewProductName = bounded.BoundedStringOf(1, 200)
product, err := NewProductName("Widget")

// NonEmpty convenience (min=1)
title, err := bounded.NonEmptyString(50, "  trimmed input  ")

// Trimmed input
clean, err := bounded.TrimmedBoundedString(1, 50, "  hello  ") // "hello"

// Money without float errors (via bojanz/currency)
price := types.NewCents(1099) // $10.99
fmt.Println(price.Float64()) // 10.99

// ISO 4217 Money with full currency support
usd, err := money.NewMoney("99.99", "USD")
eur, err := money.NewMoneyFromCents(1099, "EUR") // €10.99

// Format for locale
formatted := money.FormatMoney(usd, "de_DE") // "99,99 €"

// Currency utilities
money.IsValidCurrency("USD")     // true
digits, _ := money.CurrencyDigits("JPY") // 0 (no decimal places)
codes := money.AllCurrencyCodes()         // all ISO 4217 codes

// Percentage (clamped to 0-100)
tax := types.NewPercentage(8)  // 8%
fmt.Println(tax.Float64())   // 0.08
```

## DataPoint[T] - Complete Audit Trail

`DataPoint[T]` is a self-contained unit of data with complete audit trail. Inspired by event sourcing, it preserves ALL relationships and metadata at the application layer, enabling full traceability without external systems.

### Core Features

- **NanoID**: Unique, URL-safe identifier (21 chars by default)
- **Bitemporal tracking**: `validFrom`, `validUntil`, `recorded` timestamps
- **Actor tracking**: Who caused this data point (User, Bot, Service, System)
- **Trigger**: What caused this data point (Manual, Scheduled, Webhook, etc.)
- **Context**: Execution environment, session, request, source
- **References**: Type-safe references to related entities
- **Causes**: Causal chain for building audit/lineage graphs
- **Tags**: Additional metadata key-value pairs
- **Version**: Optimistic concurrency support

### Basic Usage

```go
import (
    id "github.com/larsartmann/go-branded-id"
    "github.com/larsartmann/go-composable-business-types/actor"
    "github.com/larsartmann/go-composable-business-types/datapoint"
    "github.com/larsartmann/go-composable-business-types/enums"
)

type OrderState struct {
    OrderId   string
    Status    string
    Total     int
}

// Create a DataPoint
actorEntry := actor.UserActor(id.NewID[struct{}, string]("user-1"), "Alice")
dp := datapoint.NewDataPoint(OrderState{
    OrderId: "order-123",
    Status:  "created",
    Total:   9900,
}, actorEntry).WithReason("customer placed order")

// Access fields
fmt.Println(dp.ID())           // NanoID (unique)
fmt.Println(dp.Payload())      // OrderState
fmt.Println(dp.Actor().Name)   // "Alice"
fmt.Println(dp.Trigger())      // TriggerManual
fmt.Println(dp.Reason())       // "customer placed order"
```

### With Builder Methods

```go
import (
    "github.com/larsartmann/go-composable-business-types/datapoint"
    "github.com/larsartmann/go-composable-business-types/enums"
    "github.com/larsartmann/go-composable-business-types/nanoid"
)

dp := datapoint.NewDataPoint(payload, actorEntry).
    WithTrigger(enums.TriggerWebhook).
    WithReason("webhook received from payment provider").
    WithContext(datapoint.NewContext().
        WithEnvironment("production").
        WithSource("payment-service").
        WithSession("sess-abc")).
    WithVersion(3).
    WithTag("correlation_id", "corr-123").
    WithReference(datapoint.NewReference("order-456", "parent")).
    WithCause(datapoint.NewCauseCommand[string](nanoid.New(), "approved"))
```

### Bitemporal Tracking

```go
import (
    "time"

    "github.com/larsartmann/go-composable-business-types/datapoint"
    "github.com/larsartmann/go-composable-business-types/temporal"
    "github.com/larsartmann/go-composable-business-types/types"
)

// Create with explicit time range
from := types.NewTimestamp(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
until := types.NewTimestamp(time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC))
recorded := types.Now()

temp := temporal.NewBitemporalWithRange(from, until, recorded)
dp := datapoint.NewDataPoint(payload, actorEntry).WithTemporal(temp)

// Check if valid at a point in time
if dp.Temporal().IsValidAt(someTime) {
    // DataPoint was valid at that time
}
```

### References and Causal Chain

```go
import (
    "github.com/larsartmann/go-composable-business-types/datapoint"
    "github.com/larsartmann/go-composable-business-types/nanoid"
)

// Reference to another entity
ref := datapoint.NewReference("doc-123", "source").
    WithVersion(5).
    WithTag("department", "legal")

// Cause tracking (event-triggered)
causeID := nanoid.New()
trace := []nanoid.NanoID{intermediateId}
cause := datapoint.NewCauseEvent[string](causeID, "created", trace...)

dp := datapoint.NewDataPoint(payload, actorEntry).
    WithReference(ref).
    WithCause(cause)
```

### JSON Serialization

```go
import "encoding/json"

// Full JSON support with round-trip
data, _ := json.Marshal(dp)
var parsed datapoint.DataPoint[OrderState]
json.Unmarshal(data, &parsed)

// Fields are preserved: id, payload, actor, temporal, trigger,
// reason, context, version, tags, references, causes
```

## Additional Packages

| Package                  | Purpose                                                      |
| ------------------------ | ------------------------------------------------------------ |
| `validate/`              | `Validator` interface for self-validating types              |
| `pkg/errors/`            | Centralized sentinel and structured error definitions        |
| `scanutil/`              | Helpers for implementing `sql.Scanner` / `driver.Valuer`    |
| `version/`               | Build version info from `runtime/debug.ReadBuildInfo`        |
| `tag/`                   | Validated `Tag` type with length and character constraints   |
| `importance/`            | `Importance` classification (0-100) with named levels        |
| `locale/`                | BCP 47 `Locale` wrapping `golang.org/x/text/language`        |
| `programminglanguage/`   | Programming language normalization and branded ID type        |
| `projectcore/`           | Composite project metadata type (name, path, languages, tags)|

See the [`examples/`](./examples/) directory for complete working programs.

## Generate

Enums are generated with `go-enum`:

```bash
go generate ./...
```

## Dependencies

| Package                                           | Purpose                               |
| ------------------------------------------------- | ------------------------------------- |
| `github.com/larsartmann/go-branded-id`            | Branded/phantom-type identifiers      |
| `github.com/bojanz/currency`                      | ISO 4217 currency handling            |
| `github.com/sixafter/nanoid`                      | FIPS-140 compatible NanoID generation |
| `golang.org/x/text`                               | BCP 47 locale/language support        |
| `github.com/abice/go-enum`                        | Enum code generation (dev only)       |

## Documentation

| Document                       | Description                                                   |
| ------------------------------ | ------------------------------------------------------------- |
| [README.md](./README.md)       | This file — usage and examples                                |
| [POLICY.md](./POLICY.md)       | Library policies: versioning, breaking changes, contributions |
| [SUPPORT.md](./SUPPORT.md)     | Getting help, reporting issues, FAQ                           |
| [SECURITY.md](./SECURITY.md)   | Security policy and vulnerability reporting                   |
| [CHANGELOG.md](./CHANGELOG.md) | Version history and release notes                             |
| [PARTS.md](./PARTS.md)         | Component analysis and architecture                           |

## Contributing

See [POLICY.md](./POLICY.md) for:

- Contribution guidelines
- Breaking changes policy
- Commit message conventions
- Code review process

## Support

- **Issues**: [GitHub Issues](https://github.com/larsartmann/go-composable-business-types/issues)
- **Discussions**: [GitHub Discussions](https://github.com/larsartmann/go-composable-business-types/discussions)
- **Security**: See [SECURITY.md](./SECURITY.md)

## License

MIT
