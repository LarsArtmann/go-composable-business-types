// Example: DataPoint with audit trail
//
// This example demonstrates using DataPoint for complete audit trails.
package main

import (
	"encoding/json"
	"fmt"

	"github.com/larsartmann/go-composable-business-types/actor"
	"github.com/larsartmann/go-composable-business-types/datapoint"
	"github.com/larsartmann/go-composable-business-types/enums"
	"github.com/larsartmann/go-composable-business-types/id"
	"github.com/larsartmann/go-composable-business-types/nanoid"
)

// Order represents a business domain type
type Order struct {
	OrderID    string `json:"order_id"`
	CustomerID string `json:"customer_id"`
	Amount     int64  `json:"amount"`
}

// Example values for demonstration purposes
const (
	exampleOrderAmount = 10000 // $100.00 in cents
	exampleRefVersion  = 5
)

func main() {
	// Create an actor (user who initiated the action)
	userID := id.NewID[struct{}, string]("user-123")
	actorEntry := actor.UserActor(userID, "Alice")

	// Create execution context
	ctx := datapoint.NewContext().
		WithEnvironment("production").
		WithSource("order-service").
		WithSession("session-abc-123").
		WithTag("region", "us-east-1")

	// Create the order payload
	order := Order{
		OrderID:    "ORD-2024-001",
		CustomerID: "CUST-456",
		Amount:     exampleOrderAmount,
	}

	// Create DataPoint with complete audit trail
	dp := datapoint.NewDataPoint(order, actorEntry).
		WithTrigger(enums.TriggerWebhook).
		WithReason("Customer checkout via web").
		WithContext(ctx).
		WithTag("priority", "high").
		WithTag("channel", "web").
		WithVersion(1)

	// Add references to related entities
	customerRef := datapoint.NewReference("CUST-456", "customer").
		WithVersion(exampleRefVersion)
	dp = dp.WithReference(customerRef)

	// Add causal chain (this order was triggered by a cart event)
	causeID := nanoid.NewNanoID()
	cartNanoID, _ := nanoid.ParseNanoID("cart-123")
	trace := []nanoid.NanoID{cartNanoID}
	cause := datapoint.NewCauseEvent[string](causeID, "cart-checked-out", trace...)
	dp = dp.WithCause(cause)

	// Display DataPoint info
	fmt.Println("=== DataPoint with Audit Trail ===")
	fmt.Printf("ID: %s\n", dp.ID().String())
	fmt.Printf("Payload: %+v\n", dp.Payload())
	fmt.Printf("Actor: %s (kind: %s)\n", dp.Actor().Name, dp.Actor().Kind.String())
	fmt.Printf("Trigger: %s\n", dp.Trigger().String())
	fmt.Printf("Reason: %s\n", dp.Reason())
	fmt.Printf("Environment: %s\n", dp.Context().Environment())
	fmt.Printf("Tags: %v\n", dp.Tags())
	fmt.Printf("References: %d\n", len(dp.References()))
	fmt.Printf("Causes: %d\n", len(dp.Causes()))

	// Serialize to JSON with indentation
	data, err := json.MarshalIndent(dp, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println("\n=== JSON Representation ===")
	fmt.Println(string(data))

	// Deserialize from JSON
	var parsed datapoint.DataPoint[Order]
	if err := json.Unmarshal(data, &parsed); err != nil {
		panic(err)
	}

	fmt.Println("\n=== Parsed DataPoint ===")
	fmt.Printf("Payload: %+v\n", parsed.Payload())
	fmt.Printf("Actor: %s\n", parsed.Actor().Name)
}
