package id

import (
	"encoding/json"
	"fmt"
)

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
