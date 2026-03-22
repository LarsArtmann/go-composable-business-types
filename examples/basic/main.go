// Example: Basic selective imports
//
// This example demonstrates importing only the types you need.
package main

import (
	"fmt"
	"log"

	"github.com/larsartmann/go-composable-business-types/nanoid"
	"github.com/larsartmann/go-composable-business-types/types"
)

// Example values for demonstration purposes
const (
	exampleCents      = 1099 // $10.99 in cents
	examplePercentage = 75   // 75%
)

func main() {
	// Generate a unique NanoID
	id := nanoid.NewNanoID()
	fmt.Println("Generated NanoID:", id.String())

	// Parse and validate an email
	email, err := types.NewEmail("user@example.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Email:", email.String())
	fmt.Println("Local part:", email.LocalPart())
	fmt.Println("Domain:", email.Domain())

	// Create monetary amounts
	cents := types.NewCents(exampleCents)
	fmt.Printf("Cents: %s (%.2f)\n", cents.String(), cents.Float64())

	// Create percentage
	pct := types.NewPercentage(examplePercentage)
	fmt.Printf("Percentage: %s (%.2f)\n", pct.String(), pct.Float64())

	// Parse URL
	url, err := types.NewURL("https://example.com/api")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("URL:", url.String())
	fmt.Println("Host:", url.Host())
}
