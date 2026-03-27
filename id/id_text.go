package id

import (
	"encoding"
	"errors"
	"fmt"
	"strconv"
)

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
		v, ok := any(string(data)).(V)
		if !ok {
			return errors.New("id: internal error: type assertion failed for string")
		}
		*id = ID[B, V]{value: v}
		return nil
	case int:
		n, err := strconv.Atoi(string(data))
		if err != nil {
			return fmt.Errorf("id: cannot parse %q as int: %w", data, err)
		}
		v, ok := any(n).(V)
		if !ok {
			return errors.New("id: internal error: type assertion failed for int")
		}
		*id = ID[B, V]{value: v}
		return nil
	case int64:
		n, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return fmt.Errorf("id: cannot parse %q as int64: %w", data, err)
		}
		v, ok := any(n).(V)
		if !ok {
			return errors.New("id: internal error: type assertion failed for int64")
		}
		*id = ID[B, V]{value: v}
		return nil
	case uint64:
		n, err := strconv.ParseUint(string(data), 10, 64)
		if err != nil {
			return fmt.Errorf("id: cannot parse %q as uint64: %w", data, err)
		}
		v, ok := any(n).(V)
		if !ok {
			return errors.New("id: internal error: type assertion failed for uint64")
		}
		*id = ID[B, V]{value: v}
		return nil
	default:
		return fmt.Errorf(
			"id: cannot unmarshal text into %T (only string and numeric IDs supported, got data=%q)",
			zero,
			string(data),
		)
	}
}

// Compile-time interface assertions for text marshaling
var (
	_ encoding.TextMarshaler   = ID[struct{}, string]{}
	_ encoding.TextUnmarshaler = (*ID[struct{}, string])(nil)
	_ encoding.TextMarshaler   = ID[struct{}, int64]{}
	_ encoding.TextUnmarshaler = (*ID[struct{}, int64])(nil)
)
