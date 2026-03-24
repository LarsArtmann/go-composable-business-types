// Package validate provides interfaces for common validation patterns.
//
// This package defines the Validator interface for types that can validate
// their own values, enabling consistent validation patterns across the codebase.
//
// Example usage:
//
//	type Email string
//
//	func (e Email) Validate() error {
//		if e == "" {
//			return ErrEmpty
//		}
//		if !emailRegex.MatchString(string(e)) {
//			return ErrInvalid
//		}
//		return nil
//	}
//
//	func Process(v validate.Validatable) error {
//		return v.Validate()
//	}
package validate

// Validator defines the interface for types that can validate themselves.
// Implement this interface to enable consistent validation patterns.
type Validator interface {
	// Validate returns an error if the value is invalid.
	// Returns nil if the value is valid.
	Validate() error
}

// Validatable is a type alias for Validator for ergonomic usage.
//
// Deprecated: Use Validator directly.
type Validatable = Validator

// IsValid checks if a value is valid without returning the error.
// Use this when you only need to know if a value is valid.
func IsValid(v Validator) bool {
	return v.Validate() == nil
}
