package cbt

import (
	"encoding"
	"fmt"
)

// ID is a branded, strongly-typed identifier that prevents mixing different entity IDs.
// B is the brand (phantom type for distinctness), V is the value type.
//
// Example:
//
//	type UserBrand struct{}
//	type OrderBrand struct{}
//	type UserID = ID[UserBrand, string]
//	type OrderID = ID[OrderBrand, int64]
type ID[B any, V comparable] struct{ value V }

// NewID creates a new branded identifier.
func NewID[B any, V comparable](v V) ID[B, V] { return ID[B, V]{value: v} }

// Value returns the underlying value.
func (id ID[B, V]) Value() V { return id.value }

// IsZero returns true if the ID has its zero value.
func (id ID[B, V]) IsZero() bool { var zero V; return id.value == zero }

// String returns a string representation of the value.
func (id ID[B, V]) String() string { return fmt.Sprintf("%v", id.value) }

// GoString implements fmt.GoStringer for debugging.
func (id ID[B, V]) GoString() string { return fmt.Sprintf("%v", id.value) }

// MarshalText implements encoding.TextMarshaler for JSON serialization.
func (id ID[B, V]) MarshalText() ([]byte, error) {
	if id.IsZero() {
		return nil, nil
	}
	return []byte(id.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler for JSON deserialization.
// Note: This only works for string-based IDs. For other types, implement a custom unmarshaler.
func (id *ID[B, V]) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		var zero V
		*id = ID[B, V]{value: zero}
		return nil
	}

	var zero V
	switch any(zero).(type) {
	case string:
		if v, ok := any(string(data)).(V); ok {
			*id = ID[B, V]{value: v}
			return nil
		}
	}
	return fmt.Errorf("id: cannot unmarshal into %T", zero)
}

// Compile-time interface assertions
var (
	_ fmt.Stringer           = ID[struct{}, string]{}
	_ fmt.GoStringer         = ID[struct{}, string]{}
	_ encoding.TextMarshaler = ID[struct{}, string]{}
	_ encoding.TextUnmarshaler = (*ID[struct{}, string])(nil)
)

// Id is an unbranded identifier for backwards compatibility.
// For new code, prefer ID[B, V] with explicit branding.
type Id[T comparable] = ID[struct{}, T]

// NewId creates a new unbranded identifier (backwards compatibility).
func NewId[T comparable](v T) Id[T] { return NewID[struct{}, T](v) }
