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

func main() {
	// Generate a unique NanoId
	id := nanoid.NewNanoId()
	fmt.Println("Generated NanoId:", id.String())

	// Parse and validate an email
	email, err := types.NewEmail("user@example.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Email:", email.String())
	fmt.Println("Local part:", email.LocalPart())
	fmt.Println("Domain:", email.Domain())

	// Create monetary amounts
	cents := types.NewCents(1099)
	fmt.Printf("Cents: %s (%.2f)\n", cents.String(), cents.Float64())

	// Create percentage
	pct := types.NewPercentage(75)
	fmt.Printf("Percentage: %s (%.2f)\n", pct.String(), pct.Float64())

	// Parse URL
	url, err := types.NewURL("https://example.com/api")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("URL:", url.String())
	fmt.Println("Host:", url.Host())
}
