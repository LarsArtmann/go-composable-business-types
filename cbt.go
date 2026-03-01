// Package cbt provides composable business types for Go applications.
// This is the root package that provides convenient access to all subpackages.
// For selective imports, use the subpackages directly:
//   - github.com/larsartmann/go-composable-business-types/id
//   - github.com/larsartmann/go-composable-business-types/nanoid
//   - github.com/larsartmann/go-composable-business-types/types
//   - github.com/larsartmann/go-composable-business-types/enums
//   - github.com/larsartmann/go-composable-business-types/bounded
//   - github.com/larsartmann/go-composable-business-types/locale
//   - github.com/larsartmann/go-composable-business-types/temporal
//   - github.com/larsartmann/go-composable-business-types/actor
//   - github.com/larsartmann/go-composable-business-types/money
//
// Note: Generic types (ID, ActorEntry, ActorChain) must be imported from
// their respective subpackages with type parameters.
package cbt

// Subpackage imports for selective use:
// Use these directly for fine-grained control over dependencies.
import (
	_ "github.com/larsartmann/go-composable-business-types/actor"
	_ "github.com/larsartmann/go-composable-business-types/bounded"
	_ "github.com/larsartmann/go-composable-business-types/enums"
	_ "github.com/larsartmann/go-composable-business-types/id"
	_ "github.com/larsartmann/go-composable-business-types/locale"
	_ "github.com/larsartmann/go-composable-business-types/money"
	_ "github.com/larsartmann/go-composable-business-types/nanoid"
	_ "github.com/larsartmann/go-composable-business-types/temporal"
	_ "github.com/larsartmann/go-composable-business-types/types"
)
