// Package id provides branded, strongly-typed identifiers that prevent mixing
// different entity IDs at compile time using phantom types.
//
// # Type Safety
//
// The ID type uses phantom typing (brand types) to create distinct identifier
// types that cannot be accidentally mixed:
//
//	type UserBrand struct{}
//	type OrderBrand struct{}
//	type UserID = ID[UserBrand, string]
//	type OrderID = ID[OrderBrand, int64]
//
//	func ProcessUser(id UserID) { ... }
//	func ProcessOrder(id OrderID) { ... }
//
//	userID := NewID[UserBrand]("user-123")
//	orderID := NewID[OrderBrand, int64](456)
//
//	ProcessUser(userID)   // OK
//	ProcessUser(orderID)  // Compile error: type mismatch
//
// # Serialization
//
// ID supports multiple serialization formats:
//   - JSON: string-based IDs serialize as strings, numeric IDs as numbers
//   - Text (XML/TOML): string-based IDs only
//   - SQL: string, int64, int32, uint64 types supported
//   - Binary: efficient binary representation
//   - Gob: Go-specific encoding
package id

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"strconv"
)

// ID is a branded, strongly-typed identifier that prevents mixing different entity IDs.
// B is the brand (phantom type for distinctness), V is the value type.
//
// The zero value represents an unset/empty ID and serializes to null in JSON.
type ID[B any, V comparable] struct{ value V }

// NewID creates a new branded identifier from the given value.
func NewID[B any, V comparable](v V) ID[B, V] { return ID[B, V]{value: v} }

// Get returns the underlying value.
func (id ID[B, V]) Get() V { return id.value }

// IsZero returns true if the ID has its zero value.
func (id ID[B, V]) IsZero() bool {
	var zero V
	return id.value == zero
}

// Reset sets the ID to its zero value.
func (id *ID[B, V]) Reset() {
	var zero V
	*id = ID[B, V]{value: zero}
}

// Equal returns true if this ID equals the other ID.
// Both IDs must have the same brand and value type.
func (id ID[B, V]) Equal(other ID[B, V]) bool {
	return id.value == other.value
}

// Compare returns -1 if id < other, 0 if equal, 1 if id > other.
// Requires V to be ordered. Panics if V is not an ordered type.
func (id ID[B, V]) Compare(other ID[B, V]) int {
	switch any(id.value).(type) {
	case int:
		a := any(id.value).(int)
		b := any(other.value).(int)
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	case int8:
		a := any(id.value).(int8)
		b := any(other.value).(int8)
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	case int16:
		a := any(id.value).(int16)
		b := any(other.value).(int16)
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	case int32:
		a := any(id.value).(int32)
		b := any(other.value).(int32)
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	case int64:
		a := any(id.value).(int64)
		b := any(other.value).(int64)
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	case uint:
		a := any(id.value).(uint)
		b := any(other.value).(uint)
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	case uint8:
		a := any(id.value).(uint8)
		b := any(other.value).(uint8)
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	case uint16:
		a := any(id.value).(uint16)
		b := any(other.value).(uint16)
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	case uint32:
		a := any(id.value).(uint32)
		b := any(other.value).(uint32)
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	case uint64:
		a := any(id.value).(uint64)
		b := any(other.value).(uint64)
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	case string:
		a := any(id.value).(string)
		b := any(other.value).(string)
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	case float32:
		a := any(id.value).(float32)
		b := any(other.value).(float32)
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	case float64:
		a := any(id.value).(float64)
		b := any(other.value).(float64)
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	default:
		panic(fmt.Sprintf("id: Compare called on non-ordered type %T", id.value))
	}
}

// Or returns the ID if not zero, otherwise returns the provided default.
func (id ID[B, V]) Or(defaultValue ID[B, V]) ID[B, V] {
	if id.IsZero() {
		return defaultValue
	}
	return id
}

// String returns a string representation of the value.
func (id ID[B, V]) String() string {
	switch v := any(id.value).(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	default:
		return fmt.Sprintf("%v", id.value)
	}
}

// GoString implements fmt.GoStringer for debugging.
func (id ID[B, V]) GoString() string { return id.String() }

// Format implements fmt.Formatter for custom formatting.
// Supports %s (string), %d (decimal), %v (default), %#v (GoString), %q (quoted).
func (id ID[B, V]) Format(f fmt.State, verb rune) {
	switch verb {
	case 's':
		_, _ = fmt.Fprint(f, id.String())
	case 'd':
		switch v := any(id.value).(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			_, _ = fmt.Fprintf(f, "%d", v)
		default:
			_, _ = fmt.Fprintf(f, "%%!d(type=%T)", id.value)
		}
	case 'q':
		_, _ = fmt.Fprintf(f, "%q", id.String())
	case 'v':
		if f.Flag('#') {
			_, _ = fmt.Fprintf(f, "id(%s)", id.String())
		} else {
			_, _ = fmt.Fprint(f, id.String())
		}
	default:
		_, _ = fmt.Fprintf(f, "%%!%c(type=%T)", verb, id.value)
	}
}

// MarshalJSON implements json.Marshaler for proper null handling.
// String-based IDs serialize as JSON strings, numeric IDs as JSON numbers.
// Zero values serialize to JSON null.
func (id ID[B, V]) MarshalJSON() ([]byte, error) {
	if id.IsZero() {
		return []byte("null"), nil
	}

	switch v := any(id.value).(type) {
	case string:
		return json.Marshal(v)
	case int:
		return json.Marshal(v)
	case int8:
		return json.Marshal(v)
	case int16:
		return json.Marshal(v)
	case int32:
		return json.Marshal(v)
	case int64:
		return json.Marshal(v)
	case uint:
		return json.Marshal(v)
	case uint8:
		return json.Marshal(v)
	case uint16:
		return json.Marshal(v)
	case uint32:
		return json.Marshal(v)
	case uint64:
		return json.Marshal(v)
	default:
		return json.Marshal(fmt.Sprintf("%v", v))
	}
}

// UnmarshalJSON implements json.Unmarshaler for JSON deserialization.
// Supports null, strings, and numeric values based on the underlying type V.
func (id *ID[B, V]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		id.Reset()
		return nil
	}

	var zero V
	switch any(zero).(type) {
	case string:
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return fmt.Errorf("id: cannot unmarshal %s into %T", string(data), zero)
		}
		*id = ID[B, V]{value: any(s).(V)}
		return nil

	case int:
		var n int
		if err := json.Unmarshal(data, &n); err != nil {
			return fmt.Errorf("id: cannot unmarshal %s into %T", string(data), zero)
		}
		*id = ID[B, V]{value: any(n).(V)}
		return nil

	case int8:
		var n int8
		if err := json.Unmarshal(data, &n); err != nil {
			return fmt.Errorf("id: cannot unmarshal %s into %T", string(data), zero)
		}
		*id = ID[B, V]{value: any(n).(V)}
		return nil

	case int16:
		var n int16
		if err := json.Unmarshal(data, &n); err != nil {
			return fmt.Errorf("id: cannot unmarshal %s into %T", string(data), zero)
		}
		*id = ID[B, V]{value: any(n).(V)}
		return nil

	case int32:
		var n int32
		if err := json.Unmarshal(data, &n); err != nil {
			return fmt.Errorf("id: cannot unmarshal %s into %T", string(data), zero)
		}
		*id = ID[B, V]{value: any(n).(V)}
		return nil

	case int64:
		var n int64
		if err := json.Unmarshal(data, &n); err != nil {
			return fmt.Errorf("id: cannot unmarshal %s into %T", string(data), zero)
		}
		*id = ID[B, V]{value: any(n).(V)}
		return nil

	case uint:
		var n uint
		if err := json.Unmarshal(data, &n); err != nil {
			return fmt.Errorf("id: cannot unmarshal %s into %T", string(data), zero)
		}
		*id = ID[B, V]{value: any(n).(V)}
		return nil

	case uint8:
		var n uint8
		if err := json.Unmarshal(data, &n); err != nil {
			return fmt.Errorf("id: cannot unmarshal %s into %T", string(data), zero)
		}
		*id = ID[B, V]{value: any(n).(V)}
		return nil

	case uint16:
		var n uint16
		if err := json.Unmarshal(data, &n); err != nil {
			return fmt.Errorf("id: cannot unmarshal %s into %T", string(data), zero)
		}
		*id = ID[B, V]{value: any(n).(V)}
		return nil

	case uint32:
		var n uint32
		if err := json.Unmarshal(data, &n); err != nil {
			return fmt.Errorf("id: cannot unmarshal %s into %T", string(data), zero)
		}
		*id = ID[B, V]{value: any(n).(V)}
		return nil

	case uint64:
		var n uint64
		if err := json.Unmarshal(data, &n); err != nil {
			return fmt.Errorf("id: cannot unmarshal %s into %T", string(data), zero)
		}
		*id = ID[B, V]{value: any(n).(V)}
		return nil

	default:
		return fmt.Errorf("id: unsupported type %T for JSON unmarshaling", zero)
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

// UnmarshalText implements encoding.TextUnmarshaler for text-based decoding (e.g., XML, TOML).
// Note: This only works for string-based IDs. For other types, use JSON.
func (id *ID[B, V]) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		id.Reset()
		return nil
	}

	var zero V
	switch any(zero).(type) {
	case string:
		*id = ID[B, V]{value: any(string(data)).(V)}
		return nil
	case int:
		n, err := strconv.Atoi(string(data))
		if err != nil {
			return fmt.Errorf("id: cannot parse %q as int: %w", data, err)
		}
		*id = ID[B, V]{value: any(n).(V)}
		return nil
	case int64:
		n, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return fmt.Errorf("id: cannot parse %q as int64: %w", data, err)
		}
		*id = ID[B, V]{value: any(n).(V)}
		return nil
	case uint64:
		n, err := strconv.ParseUint(string(data), 10, 64)
		if err != nil {
			return fmt.Errorf("id: cannot parse %q as uint64: %w", data, err)
		}
		*id = ID[B, V]{value: any(n).(V)}
		return nil
	default:
		return fmt.Errorf("id: cannot unmarshal text into %T (only string and numeric IDs supported)", zero)
	}
}

// MarshalBinary implements encoding.BinaryMarshaler for binary encoding.
func (id ID[B, V]) MarshalBinary() ([]byte, error) {
	if id.IsZero() {
		return nil, nil
	}

	var buf bytes.Buffer
	switch v := any(id.value).(type) {
	case string:
		return []byte(v), nil
	case int:
		if err := binary.Write(&buf, binary.LittleEndian, int64(v)); err != nil {
			return nil, fmt.Errorf("id: failed to marshal int: %w", err)
		}
		return buf.Bytes(), nil
	case int8:
		return []byte{byte(v)}, nil
	case int16:
		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, uint16(v))
		return b, nil
	case int32:
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(v))
		return b, nil
	case int64:
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(v))
		return b, nil
	case uint:
		if err := binary.Write(&buf, binary.LittleEndian, uint64(v)); err != nil {
			return nil, fmt.Errorf("id: failed to marshal uint: %w", err)
		}
		return buf.Bytes(), nil
	case uint8:
		return []byte{v}, nil
	case uint16:
		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b, v)
		return b, nil
	case uint32:
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, v)
		return b, nil
	case uint64:
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, v)
		return b, nil
	default:
		return nil, fmt.Errorf("id: unsupported type %T for binary marshaling", id.value)
	}
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler for binary decoding.
func (id *ID[B, V]) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		id.Reset()
		return nil
	}

	var zero V
	switch any(zero).(type) {
	case string:
		*id = ID[B, V]{value: any(string(data)).(V)}
		return nil
	case int:
		if len(data) < 8 {
			return fmt.Errorf("id: insufficient data for int: got %d bytes, want 8", len(data))
		}
		n := int(binary.LittleEndian.Uint64(data))
		*id = ID[B, V]{value: any(n).(V)}
		return nil
	case int8:
		if len(data) < 1 {
			return fmt.Errorf("id: insufficient data for int8")
		}
		*id = ID[B, V]{value: any(int8(data[0])).(V)}
		return nil
	case int16:
		if len(data) < 2 {
			return fmt.Errorf("id: insufficient data for int16: got %d bytes, want 2", len(data))
		}
		n := int16(binary.LittleEndian.Uint16(data))
		*id = ID[B, V]{value: any(n).(V)}
		return nil
	case int32:
		if len(data) < 4 {
			return fmt.Errorf("id: insufficient data for int32: got %d bytes, want 4", len(data))
		}
		n := int32(binary.LittleEndian.Uint32(data))
		*id = ID[B, V]{value: any(n).(V)}
		return nil
	case int64:
		if len(data) < 8 {
			return fmt.Errorf("id: insufficient data for int64: got %d bytes, want 8", len(data))
		}
		n := int64(binary.LittleEndian.Uint64(data))
		*id = ID[B, V]{value: any(n).(V)}
		return nil
	case uint:
		if len(data) < 8 {
			return fmt.Errorf("id: insufficient data for uint: got %d bytes, want 8", len(data))
		}
		n := uint(binary.LittleEndian.Uint64(data))
		*id = ID[B, V]{value: any(n).(V)}
		return nil
	case uint8:
		if len(data) < 1 {
			return fmt.Errorf("id: insufficient data for uint8")
		}
		*id = ID[B, V]{value: any(data[0]).(V)}
		return nil
	case uint16:
		if len(data) < 2 {
			return fmt.Errorf("id: insufficient data for uint16: got %d bytes, want 2", len(data))
		}
		n := binary.LittleEndian.Uint16(data)
		*id = ID[B, V]{value: any(n).(V)}
		return nil
	case uint32:
		if len(data) < 4 {
			return fmt.Errorf("id: insufficient data for uint32: got %d bytes, want 4", len(data))
		}
		n := binary.LittleEndian.Uint32(data)
		*id = ID[B, V]{value: any(n).(V)}
		return nil
	case uint64:
		if len(data) < 8 {
			return fmt.Errorf("id: insufficient data for uint64: got %d bytes, want 8", len(data))
		}
		n := binary.LittleEndian.Uint64(data)
		*id = ID[B, V]{value: any(n).(V)}
		return nil
	default:
		return fmt.Errorf("id: unsupported type %T for binary unmarshaling", zero)
	}
}

// GobEncode implements gob.GobEncoder for Go-specific encoding.
func (id ID[B, V]) GobEncode() ([]byte, error) {
	return id.MarshalBinary()
}

// GobDecode implements gob.GobDecoder for Go-specific decoding.
func (id *ID[B, V]) GobDecode(data []byte) error {
	return id.UnmarshalBinary(data)
}

// Scan implements sql.Scanner for database deserialization.
// Supports string, []byte, int64, int, float64, and nil sources based on the underlying value type V.
func (id *ID[B, V]) Scan(src any) error {
	if src == nil {
		id.Reset()
		return nil
	}

	var zero V
	switch any(zero).(type) {
	case string:
		switch v := src.(type) {
		case string:
			*id = ID[B, V]{value: any(v).(V)}
			return nil
		case []byte:
			*id = ID[B, V]{value: any(string(v)).(V)}
			return nil
		default:
			return fmt.Errorf("id: cannot scan %T into string-based ID", src)
		}

	case int:
		switch v := src.(type) {
		case int64:
			*id = ID[B, V]{value: any(int(v)).(V)}
			return nil
		case int:
			*id = ID[B, V]{value: any(v).(V)}
			return nil
		case float64:
			*id = ID[B, V]{value: any(int(v)).(V)}
			return nil
		default:
			return fmt.Errorf("id: cannot scan %T into int-based ID", src)
		}

	case int32:
		switch v := src.(type) {
		case int64:
			*id = ID[B, V]{value: any(int32(v)).(V)}
			return nil
		case int:
			*id = ID[B, V]{value: any(int32(v)).(V)}
			return nil
		case float64:
			*id = ID[B, V]{value: any(int32(v)).(V)}
			return nil
		default:
			return fmt.Errorf("id: cannot scan %T into int32-based ID", src)
		}

	case int64:
		switch v := src.(type) {
		case int64:
			*id = ID[B, V]{value: any(v).(V)}
			return nil
		case int:
			*id = ID[B, V]{value: any(int64(v)).(V)}
			return nil
		case float64:
			*id = ID[B, V]{value: any(int64(v)).(V)}
			return nil
		default:
			return fmt.Errorf("id: cannot scan %T into int64-based ID", src)
		}

	case uint:
		switch v := src.(type) {
		case int64:
			*id = ID[B, V]{value: any(uint(v)).(V)}
			return nil
		case int:
			*id = ID[B, V]{value: any(uint(v)).(V)}
			return nil
		case float64:
			*id = ID[B, V]{value: any(uint(v)).(V)}
			return nil
		default:
			return fmt.Errorf("id: cannot scan %T into uint-based ID", src)
		}

	case uint32:
		switch v := src.(type) {
		case int64:
			*id = ID[B, V]{value: any(uint32(v)).(V)}
			return nil
		case int:
			*id = ID[B, V]{value: any(uint32(v)).(V)}
			return nil
		case float64:
			*id = ID[B, V]{value: any(uint32(v)).(V)}
			return nil
		default:
			return fmt.Errorf("id: cannot scan %T into uint32-based ID", src)
		}

	case uint64:
		switch v := src.(type) {
		case int64:
			*id = ID[B, V]{value: any(uint64(v)).(V)}
			return nil
		case int:
			*id = ID[B, V]{value: any(uint64(v)).(V)}
			return nil
		case float64:
			*id = ID[B, V]{value: any(uint64(v)).(V)}
			return nil
		default:
			return fmt.Errorf("id: cannot scan %T into uint64-based ID", src)
		}

	default:
		return fmt.Errorf("id: unsupported target type %T for SQL scanning", zero)
	}
}

// Value implements driver.Valuer for database serialization.
// Returns nil for zero values, otherwise the underlying value.
func (id ID[B, V]) Value() (driver.Value, error) {
	if id.IsZero() {
		return nil, nil
	}

	switch v := any(id.value).(type) {
	case string:
		return v, nil
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case uint:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		return int64(v), nil
	default:
		return nil, fmt.Errorf("id: unsupported type %T for SQL value", id.value)
	}
}

// Compile-time interface assertions
var (
	_ fmt.Stringer               = ID[struct{}, string]{}
	_ fmt.GoStringer             = ID[struct{}, string]{}
	_ fmt.Formatter              = ID[struct{}, string]{}
	_ json.Marshaler             = ID[struct{}, string]{}
	_ json.Unmarshaler           = (*ID[struct{}, string])(nil)
	_ json.Marshaler             = ID[struct{}, int64]{}
	_ json.Unmarshaler           = (*ID[struct{}, int64])(nil)
	_ json.Marshaler             = ID[struct{}, int32]{}
	_ json.Unmarshaler           = (*ID[struct{}, int32])(nil)
	_ json.Marshaler             = ID[struct{}, uint64]{}
	_ json.Unmarshaler           = (*ID[struct{}, uint64])(nil)
	_ encoding.TextMarshaler     = ID[struct{}, string]{}
	_ encoding.TextUnmarshaler   = (*ID[struct{}, string])(nil)
	_ encoding.TextMarshaler     = ID[struct{}, int64]{}
	_ encoding.TextUnmarshaler   = (*ID[struct{}, int64])(nil)
	_ encoding.BinaryMarshaler   = ID[struct{}, string]{}
	_ encoding.BinaryUnmarshaler = (*ID[struct{}, string])(nil)
	_ encoding.BinaryMarshaler   = ID[struct{}, int64]{}
	_ encoding.BinaryUnmarshaler = (*ID[struct{}, int64])(nil)
	_ gob.GobEncoder             = ID[struct{}, string]{}
	_ gob.GobDecoder             = (*ID[struct{}, string])(nil)
	_ sql.Scanner                = (*ID[struct{}, string])(nil)
	_ sql.Scanner                = (*ID[struct{}, int64])(nil)
	_ sql.Scanner                = (*ID[struct{}, int32])(nil)
	_ sql.Scanner                = (*ID[struct{}, uint64])(nil)
	_ driver.Valuer              = ID[struct{}, string]{}
	_ driver.Valuer              = ID[struct{}, int64]{}
	_ driver.Valuer              = ID[struct{}, int32]{}
	_ driver.Valuer              = ID[struct{}, uint64]{}
)
