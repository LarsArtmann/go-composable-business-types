# go-composable-business-types

Strongly typed, composable base values for safer, more interoperable Go projects.

## Installation

```bash
go get github.com/larsartmann/go-composable-business-types
```

## Types

| Type            | Purpose                                                             |
| --------------- | ------------------------------------------------------------------- |
| `Id[T]`         | Type-safe identifier wrapper - prevents mixing different entity IDs |
| `NanoId`        | URL-safe, cryptographically random ID (default 21 chars)            |
| `ActorChain[T]` | Ordered chain of actors (User → Service → Service) for audit trails |
| `DataPoint[T]`  | Self-contained data unit with complete audit trail                  |
| `Bitemporal`    | Bitemporal tracking (validFrom, validUntil, recorded)               |
| `Context`       | Execution context (environment, session, request, source)           |
| `Reference[T]`  | Type-safe reference to another entity with relationship metadata    |
| `Cause[T]`      | Causal chain tracking for building audit/lineage graphs             |
| `BoundedString` | String with validated length constraints                            |
| `Email`         | Email address string                                                |
| `URL`           | URL string                                                          |
| `Percentage`    | 0-100 value with float conversion (clamps overflow to 100)          |
| `Cents`         | Monetary amount in smallest unit (no float errors)                  |
| `Timestamp`     | Domain-wrapped time.Time                                            |
| `Duration`      | Domain-wrapped time.Duration                                        |
| `Money`         | ISO 4217 currency via `github.com/bojanz/currency`                  |

## Enums (generated)

| Enum        | Values                                                            |
| ----------- | ----------------------------------------------------------------- |
| `ActorKind` | User, Bot, System, Service                                        |
| `Locale`    | en_US, en_GB, de_DE, fr_FR, es_ES, it_IT, ja_JP, zh_CN            |
| `Priority`  | Low, Medium, High, Critical                                       |
| `Status`    | Draft, Active, Paused, Archived, Deleted                          |
| `Trigger`   | Manual, Scheduled, Webhook, Import, Migration, System, Correction |

## Usage

```go
import cbt "github.com/larsartmann/go-composable-business-types"

// Type-safe IDs - can't mix UserId with OrderId
type UserId = cbt.Id[string]
type OrderId = cbt.Id[int]

userId := cbt.NewId("user-123")
orderId := cbt.NewId(42)

// Actor chain for audit trails and authorization
chain := cbt.NewActorChain(cbt.UserActor(userId, "Alice")).
    Append(cbt.ServiceActor(cbt.NewId("api-gateway"), "API Gateway")).
    Append(cbt.ServiceActor(cbt.NewId("order-svc"), "Order Service"))

origin := chain.Origin()   // User Alice
current := chain.Current() // Order Service

// BoundedString - validated string lengths
name, err := cbt.NewBoundedString(1, 100, "John Doe")
if err != nil {
    // handle validation error
}

// Factory for domain-specific bounded strings
var NewProductName = cbt.BoundedStringOf(1, 200)
product, err := NewProductName("Widget")

// NonEmpty convenience (min=1)
title, err := cbt.NonEmptyString(50, "  trimmed input  ")

// Trimmed input
clean, err := cbt.TrimmedBoundedString(1, 50, "  hello  ") // "hello"

// Money without float errors (via bojanz/currency)
price := cbt.NewCents(1099) // $10.99
fmt.Println(price.Float64()) // 10.99

// ISO 4217 Money with full currency support
money, err := cbt.NewMoney("99.99", "USD")
money, err := cbt.NewMoneyFromCents(1099, "EUR") // €10.99

// Format for locale
formatted := cbt.FormatMoney(money, "de_DE") // "99,99 €"

// Currency utilities
cbt.IsValidCurrency("USD")     // true
cbt.CurrencyDigits("JPY")      // 0 (no decimal places)
cbt.AllCurrencyCodes()         // all ISO 4217 codes

// Percentage (clamped to 0-100)
tax := cbt.NewPercentage(8)  // 8%
fmt.Println(tax.Float64())   // 0.08
```

## DataPoint[T] - Complete Audit Trail

`DataPoint[T]` is a self-contained unit of data with complete audit trail. Inspired by event sourcing, it preserves ALL relationships and metadata at the application layer, enabling full traceability without external systems.

### Core Features

- **NanoId**: Unique, URL-safe identifier (21 chars by default)
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
type OrderState struct {
    OrderId   string
    Status    string
    Total     int
}

// Create a DataPoint
actor := cbt.UserActor(cbt.NewId("user-1"), "Alice")
dp := cbt.NewDataPointNow(OrderState{
    OrderId: "order-123",
    Status:  "created",
    Total:   9900,
}, actor, "customer placed order")

// Access fields
fmt.Println(dp.Id())           // NanoId (unique)
fmt.Println(dp.Payload())      // OrderState
fmt.Println(dp.Actor().Name)   // "Alice"
fmt.Println(dp.Trigger())      // TriggerManual
fmt.Println(dp.Reason())       // "customer placed order"
```

### With Builder Methods

```go
dp := cbt.NewDataPointNow(payload, actor).
    WithTrigger(cbt.TriggerWebhook).
    WithReason("webhook received from payment provider").
    WithContext(cbt.NewContext("payment-service").
        WithEnvironment("production").
        WithSession("sess-abc")).
    WithVersion(3).
    AddTag("correlation_id", "corr-123").
    AddReference(cbt.NewReference("order-456", "parent")).
    AddCause(cbt.NewCauseCommand(causeId, "approved"))
```

### Bitemporal Tracking

```go
// Create with explicit time range
from := cbt.NewTimestamp(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
until := cbt.NewTimestamp(time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC))
recorded := cbt.Now()

temporal := cbt.NewBitemporalWithRange(from, until, recorded)
dp := cbt.NewDataPointNow(payload, actor).WithTemporal(temporal)

// Check if valid at a point in time
if dp.Temporal().IsValidAt(someTime) {
    // DataPoint was valid at that time
}
```

### References and Causal Chain

```go
// Reference to another entity
ref := cbt.NewReferenceWithVersion("doc-123", "source", 5).
    WithTag("department", "legal")

// Cause tracking
cause := cbt.NewCauseCommand(cbt.NewNanoId(), "created").
    WithTrace([]cbt.NanoId{intermediateId})

dp := cbt.NewDataPointNow(payload, actor).
    AddReference(ref).
    AddCause(cause)
```

### JSON Serialization

```go
// Full JSON support with round-trip
data, _ := json.Marshal(dp)
var parsed cbt.DataPoint[OrderState]
json.Unmarshal(data, &parsed)

// Fields are preserved: id, payload, actor, temporal, trigger,
// reason, context, version, tags, references, causes
```

## Generate

Enums are generated with `go-enum`:

```bash
go generate ./...
```

## Dependencies

- `github.com/abice/go-enum` - Enum code generation
- `github.com/bojanz/currency` - ISO 4217 currency handling with 370+ locales

## License

MIT
