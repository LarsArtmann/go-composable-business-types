# go-composable-business-types

Strongly typed, composable base values for safer, more interoperable Go projects.

## Types

| Type | Purpose |
|------|---------|
| `Id[T]` | Type-safe identifier wrapper |
| `ActorChain[T]` | Ordered chain of actors (User → Service → Service) |
| `Email` | Email address string |
| `URL` | URL string |
| `Percentage` | 0-100 value with float conversion |
| `Cents` | Monetary amount in smallest unit (no float errors) |
| `Timestamp` | Domain-wrapped time.Time |
| `Duration` | Domain-wrapped time.Duration |

## Enums (generated)

| Enum | Values |
|------|--------|
| `ActorKind` | User, Bot, System, Service |
| `Currency` | USD, EUR, GBP, JPY, CHF, CAD, AUD, NZD, CNY, INR |
| `Locale` | en_US, en_GB, de_DE, fr_FR, es_ES, it_IT, ja_JP, zh_CN |
| `Priority` | Low, Medium, High, Critical |
| `Status` | Draft, Active, Paused, Archived, Deleted |

## Usage

```go
import "github.com/larsartmann/go-composable-business-types/cbt"

// Type-safe IDs - can't mix UserId with OrderId
type UserId = cbt.Id[string]
type OrderId = cbt.Id[int]

userId := cbt.NewId("user-123")
orderId := cbt.NewId(42)

// Actor chain for audit trails
chain := cbt.NewActorChain(cbt.UserActor(userId, "Alice")).
    Append(cbt.ServiceActor(cbt.NewId("api-gateway"))).
    Append(cbt.ServiceActor(cbt.NewId("order-svc")))

origin := chain.Origin()   // User Alice
current := chain.Current() // order-svc Service

// Money without float errors
price := cbt.NewCents(1099) // $10.99
fmt.Println(price.Float64()) // 10.99

// Percentage
tax := cbt.NewPercentage(8) // 8%
fmt.Println(tax.Float64()) // 0.08
```

## Generate

```bash
go generate ./...
```
