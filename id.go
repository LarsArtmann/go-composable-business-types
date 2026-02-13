package cbt

import (
	"encoding"
	"encoding/json"
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

// MarshalJSON implements json.Marshaler for proper null handling.
// Zero values serialize to JSON null, non-zero values serialize as JSON strings.
func (id ID[B, V]) MarshalJSON() ([]byte, error) {
	if id.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(id.String())
}

// UnmarshalJSON implements json.Unmarshaler for JSON deserialization.
func (id *ID[B, V]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		var zero V
		*id = ID[B, V]{value: zero}
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("id: cannot unmarshal %s into %T", string(data), *id)
	}

	var zero V
	switch any(zero).(type) {
	case string:
		*id = ID[B, V]{value: any(s).(V)}
		return nil
	default:
		return fmt.Errorf("id: cannot unmarshal string into %T (only string-based IDs supported)", zero)
	}
}

// MarshalText implements encoding.TextMarshaler for text-based encoding (e.g., XML, TOML).
// For JSON, prefer the json.Marshaler implementation which handles null properly.
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
	_ fmt.Stringer             = ID[struct{}, string]{}
	_ fmt.GoStringer           = ID[struct{}, string]{}
	_ json.Marshaler           = ID[struct{}, string]{}
	_ json.Unmarshaler         = (*ID[struct{}, string])(nil)
	_ encoding.TextMarshaler   = ID[struct{}, string]{}
	_ encoding.TextUnmarshaler = (*ID[struct{}, string])(nil)
)
