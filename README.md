# go-composable-business-types

Strongly typed, composable base values for safer, more interoperable Go projects.

## Installation

```bash
go get github.com/larsartmann/go-composable-business-types
```

## Types

| Type | Purpose |
|------|---------|
| `Id[T]` | Type-safe identifier wrapper - prevents mixing different entity IDs |
| `ActorChain[T]` | Ordered chain of actors (User → Service → Service) for audit trails |
| `BoundedString` | String with validated length constraints |
| `Email` | Email address string |
| `URL` | URL string |
| `Percentage` | 0-100 value with float conversion (clamps overflow to 100) |
| `Cents` | Monetary amount in smallest unit (no float errors) |
| `Timestamp` | Domain-wrapped time.Time |
| `Duration` | Domain-wrapped time.Duration |
| `Money` | ISO 4217 currency via `github.com/bojanz/currency` |

## Enums (generated)

| Enum | Values |
|------|--------|
| `ActorKind` | User, Bot, System, Service |
| `Locale` | en_US, en_GB, de_DE, fr_FR, es_ES, it_IT, ja_JP, zh_CN |
| `Priority` | Low, Medium, High, Critical |
| `Status` | Draft, Active, Paused, Archived, Deleted |

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
