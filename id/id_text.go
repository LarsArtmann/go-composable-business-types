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
		return parseSignedIntegerID(id, data, strconv.ParseInt, "int64")
	case uint64:
		return parseUnsignedIntegerID(id, data, strconv.ParseUint, "uint64")
	default:
		return fmt.Errorf(
			"id: cannot unmarshal text into %T (only string and numeric IDs supported, got data=%q)",
			zero,
			string(data),
		)
	}
}

func parseSignedIntegerID[B any, V comparable](
	id *ID[B, V],
	data []byte,
	parse func(string, int, int) (int64, error),
	typeName string,
) error {
	n, err := parse(string(data), 10, 64)
	if err != nil {
		return fmt.Errorf("id: cannot parse %q as %s: %w", data, typeName, err)
	}

	v, ok := any(n).(V)
	if !ok {
		return errors.New("id: internal error: type assertion failed for " + typeName)
	}

	*id = ID[B, V]{value: v}

	return nil
}

func parseUnsignedIntegerID[B any, V comparable](
	id *ID[B, V],
	data []byte,
	parse func(string, int, int) (uint64, error),
	typeName string,
) error {
	n, err := parse(string(data), 10, 64)
	if err != nil {
		return fmt.Errorf("id: cannot parse %q as %s: %w", data, typeName, err)
	}

	v, ok := any(n).(V)
	if !ok {
		return errors.New("id: internal error: type assertion failed for " + typeName)
	}

	*id = ID[B, V]{value: v}

	return nil
}

// Compile-time interface assertions for text marshaling.
var (
	_ encoding.TextMarshaler   = ID[struct{}, string]{value: ""}
	_ encoding.TextUnmarshaler = (*ID[struct{}, string])(nil)
	_ encoding.TextMarshaler   = ID[struct{}, int64]{value: 0}
	_ encoding.TextUnmarshaler = (*ID[struct{}, int64])(nil)
)
