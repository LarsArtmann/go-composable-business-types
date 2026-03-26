package id

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

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
	_ fmt.Stringer     = ID[struct{}, string]{}
	_ fmt.GoStringer   = ID[struct{}, string]{}
	_ fmt.Formatter    = ID[struct{}, string]{}
	_ json.Marshaler   = ID[struct{}, string]{}
	_ json.Unmarshaler = (*ID[struct{}, string])(nil)
	_ json.Marshaler   = ID[struct{}, int64]{}
	_ json.Unmarshaler = (*ID[struct{}, int64])(nil)
	_ json.Marshaler   = ID[struct{}, int32]{}
	_ json.Unmarshaler = (*ID[struct{}, int32])(nil)
	_ json.Marshaler   = ID[struct{}, uint64]{}
	_ json.Unmarshaler = (*ID[struct{}, uint64])(nil)
	_ sql.Scanner      = (*ID[struct{}, string])(nil)
	_ sql.Scanner      = (*ID[struct{}, int64])(nil)
	_ sql.Scanner      = (*ID[struct{}, int32])(nil)
	_ sql.Scanner      = (*ID[struct{}, uint64])(nil)
	_ driver.Valuer    = ID[struct{}, string]{}
	_ driver.Valuer    = ID[struct{}, int64]{}
	_ driver.Valuer    = ID[struct{}, int32]{}
	_ driver.Valuer    = ID[struct{}, uint64]{}
)
