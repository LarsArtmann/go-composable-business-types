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
	"cmp"
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// Byte sizes for binary marshaling of integer types
const (
	byteSizeInt16 = 2 // size of int16 and uint16 in bytes
	byteSizeInt32 = 4 // size of int32 and uint32 in bytes
	byteSizeInt64 = 8 // size of int64, uint, and uint64 in bytes
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

// ErrNotOrdered is returned when Compare is called on an ID with a non-ordered value type.
var ErrNotOrdered = errors.New("id: Compare requires an ordered type (int, uint, float, or string)")

// Compare returns -1 if id < other, 0 if equal, 1 if id > other.
// Returns ErrNotOrdered if V is not an ordered type.
func (id ID[B, V]) Compare(other ID[B, V]) (int, error) {
	switch a := any(id.value).(type) {
	case int:
		return cmp.Compare(a, any(other.value).(int)), nil
	case int8:
		return cmp.Compare(a, any(other.value).(int8)), nil
	case int16:
		return cmp.Compare(a, any(other.value).(int16)), nil
	case int32:
		return cmp.Compare(a, any(other.value).(int32)), nil
	case int64:
		return cmp.Compare(a, any(other.value).(int64)), nil
	case uint:
		return cmp.Compare(a, any(other.value).(uint)), nil
	case uint8:
		return cmp.Compare(a, any(other.value).(uint8)), nil
	case uint16:
		return cmp.Compare(a, any(other.value).(uint16)), nil
	case uint32:
		return cmp.Compare(a, any(other.value).(uint32)), nil
	case uint64:
		return cmp.Compare(a, any(other.value).(uint64)), nil
	case string:
		return cmp.Compare(a, any(other.value).(string)), nil
	case float32:
		return cmp.Compare(a, any(other.value).(float32)), nil
	case float64:
		return cmp.Compare(a, any(other.value).(float64)), nil
	default:
		return 0, ErrNotOrdered
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
		b, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("id: marshal string JSON: %w", err)
		}
		return b, nil
	case int:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("id: marshal int JSON: %w", err)
		}
		return b, nil
	case int8:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("id: marshal int8 JSON: %w", err)
		}
		return b, nil
	case int16:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("id: marshal int16 JSON: %w", err)
		}
		return b, nil
	case int32:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("id: marshal int32 JSON: %w", err)
		}
		return b, nil
	case int64:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("id: marshal int64 JSON: %w", err)
		}
		return b, nil
	case uint:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("id: marshal uint JSON: %w", err)
		}
		return b, nil
	case uint8:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("id: marshal uint8 JSON: %w", err)
		}
		return b, nil
	case uint16:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("id: marshal uint16 JSON: %w", err)
		}
		return b, nil
	case uint32:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("id: marshal uint32 JSON: %w", err)
		}
		return b, nil
	case uint64:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("id: marshal uint64 JSON: %w", err)
		}
		return b, nil
	default:
		b, err := json.Marshal(fmt.Sprintf("%v", v))
		if err != nil {
			return nil, fmt.Errorf("id: marshal default JSON: %w", err)
		}
		return b, nil
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
		return fmt.Errorf(
			"id: unsupported type %T for JSON unmarshaling (data=%q)",
			zero,
			string(data),
		)
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
		return fmt.Errorf(
			"id: cannot unmarshal text into %T (only string and numeric IDs supported, got data=%q)",
			zero,
			string(data),
		)
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
		return []byte{byte(v)}, nil //nolint:gosec // G115: int8 to byte is safe for serialization
	case int16:
		b := make([]byte, byteSizeInt16)
		binary.LittleEndian.PutUint16(
			b,
			uint16(v),
		)
		return b, nil
	case int32:
		b := make([]byte, byteSizeInt32)
		binary.LittleEndian.PutUint32(
			b,
			uint32(v),
		)
		return b, nil
	case int64:
		b := make([]byte, byteSizeInt64)
		binary.LittleEndian.PutUint64(
			b,
			uint64(v),
		)
		return b, nil
	case uint:
		if err := binary.Write(&buf, binary.LittleEndian, uint64(v)); err != nil {
			return nil, fmt.Errorf("id: failed to marshal uint: %w", err)
		}
		return buf.Bytes(), nil
	case uint8:
		return []byte{v}, nil
	case uint16:
		b := make([]byte, byteSizeInt16)
		binary.LittleEndian.PutUint16(b, v)
		return b, nil
	case uint32:
		b := make([]byte, byteSizeInt32)
		binary.LittleEndian.PutUint32(b, v)
		return b, nil
	case uint64:
		b := make([]byte, byteSizeInt64)
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
		if len(data) < byteSizeInt64 {
			return fmt.Errorf(
				"id: insufficient data for int: got %d bytes, want %d (data=%x, targetType=%T)",
				len(data),
				byteSizeInt64,
				data,
				zero,
			)
		}
		n := int(
			binary.LittleEndian.Uint64(data),
		)
		*id = ID[B, V]{value: any(n).(V)}
		return nil
	case int8:
		if len(data) < 1 {
			return fmt.Errorf(
				"id: insufficient data for int8: got %d bytes, want 1 (data=%x, targetType=%T)",
				len(data),
				data,
				zero,
			)
		}
		*id = ID[B, V]{value: any(int8(data[0])).(V)}
		return nil
	case int16:
		if len(data) < byteSizeInt16 {
			return fmt.Errorf(
				"id: insufficient data for int16: got %d bytes, want %d (data=%x, targetType=%T)",
				len(data),
				byteSizeInt16,
				data,
				zero,
			)
		}
		n := int16(
			binary.LittleEndian.Uint16(data),
		)
		*id = ID[B, V]{value: any(n).(V)}
		return nil
	case int32:
		if len(data) < byteSizeInt32 {
			return fmt.Errorf(
				"id: insufficient data for int32: got %d bytes, want %d (data=%x, targetType=%T)",
				len(data),
				byteSizeInt32,
				data,
				zero,
			)
		}
		n := int32(
			binary.LittleEndian.Uint32(data),
		)
		*id = ID[B, V]{value: any(n).(V)}
		return nil
	case int64:
		if len(data) < byteSizeInt64 {
			return fmt.Errorf(
				"id: insufficient data for int64: got %d bytes, want %d (data=%x, targetType=%T)",
				len(data),
				byteSizeInt64,
				data,
				zero,
			)
		}
		n := int64(
			binary.LittleEndian.Uint64(data),
		)
		*id = ID[B, V]{value: any(n).(V)}
		return nil
	case uint:
		if len(data) < byteSizeInt64 {
			return fmt.Errorf(
				"id: insufficient data for uint: got %d bytes, want %d (data=%x, targetType=%T)",
				len(data),
				byteSizeInt64,
				data,
				zero,
			)
		}
		n := uint(binary.LittleEndian.Uint64(data))
		*id = ID[B, V]{value: any(n).(V)}
		return nil
	case uint8:
		if len(data) < 1 {
			return fmt.Errorf(
				"id: insufficient data for uint8: got %d bytes, want 1 (data=%x, targetType=%T)",
				len(data),
				data,
				zero,
			)
		}
		*id = ID[B, V]{value: any(data[0]).(V)}
		return nil
	case uint16:
		if len(data) < byteSizeInt16 {
			return fmt.Errorf(
				"id: insufficient data for uint16: got %d bytes, want %d (data=%x, targetType=%T)",
				len(data),
				byteSizeInt16,
				data,
				zero,
			)
		}
		n := binary.LittleEndian.Uint16(data)
		*id = ID[B, V]{value: any(n).(V)}
		return nil
	case uint32:
		if len(data) < byteSizeInt32 {
			return fmt.Errorf(
				"id: insufficient data for uint32: got %d bytes, want %d (data=%x, targetType=%T)",
				len(data),
				byteSizeInt32,
				data,
				zero,
			)
		}
		n := binary.LittleEndian.Uint32(data)
		*id = ID[B, V]{value: any(n).(V)}
		return nil
	case uint64:
		if len(data) < byteSizeInt64 {
			return fmt.Errorf(
				"id: insufficient data for uint64: got %d bytes, want %d (data=%x, targetType=%T)",
				len(data),
				byteSizeInt64,
				data,
				zero,
			)
		}
		n := binary.LittleEndian.Uint64(data)
		*id = ID[B, V]{value: any(n).(V)}
		return nil
	default:
		return fmt.Errorf("id: unsupported type %T for binary unmarshaling (data=%x)", zero, data)
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
	if id == nil {
		return errors.New("id: scan: receiver is nil")
	}
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
			return fmt.Errorf("id: cannot scan %T into string-based ID (src=%T)", src, src)
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
			return fmt.Errorf("id: cannot scan %T into int-based ID (targetType=%T)", src, zero)
		}

	case int32:
		switch v := src.(type) {
		case int64:
			*id = ID[B, V]{
				value: any(int32(v)).(V),
			}
			return nil
		case int:
			*id = ID[B, V]{
				value: any(int32(v)).(V),
			}
			return nil
		case float64:
			*id = ID[B, V]{value: any(int32(v)).(V)}
			return nil
		default:
			return fmt.Errorf("id: cannot scan %T into int32-based ID (targetType=%T)", src, zero)
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
			return fmt.Errorf("id: cannot scan %T into int64-based ID (targetType=%T)", src, zero)
		}

	case uint:
		switch v := src.(type) {
		case int64:
			*id = ID[B, V]{
				value: any(uint(v)).(V),
			}
			return nil
		case int:
			*id = ID[B, V]{
				value: any(uint(v)).(V),
			}
			return nil
		case float64:
			*id = ID[B, V]{value: any(uint(v)).(V)}
			return nil
		default:
			return fmt.Errorf("id: cannot scan %T into uint-based ID (targetType=%T)", src, zero)
		}

	case uint32:
		switch v := src.(type) {
		case int64:
			*id = ID[B, V]{
				value: any(uint32(v)).(V),
			}
			return nil
		case int:
			*id = ID[B, V]{
				value: any(uint32(v)).(V),
			}
			return nil
		case float64:
			*id = ID[B, V]{value: any(uint32(v)).(V)}
			return nil
		default:
			return fmt.Errorf("id: cannot scan %T into uint32-based ID (targetType=%T)", src, zero)
		}

	case uint64:
		switch v := src.(type) {
		case int64:
			*id = ID[B, V]{
				value: any(uint64(v)).(V),
			}
			return nil
		case int:
			*id = ID[B, V]{
				value: any(uint64(v)).(V),
			}
			return nil
		case float64:
			*id = ID[B, V]{value: any(uint64(v)).(V)}
			return nil
		default:
			return fmt.Errorf("id: cannot scan %T into uint64-based ID (targetType=%T)", src, zero)
		}

	default:
		return fmt.Errorf("id: unsupported target type %T for SQL scanning (src=%T)", zero, src)
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
		return int64(v), nil //nolint:gosec // G115: uint to int64 for SQL value
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		return int64(v), nil //nolint:gosec // G115: uint64 to int64 for SQL value
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
