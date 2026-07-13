package types

import (
	"encoding/json/v2"
	"fmt"
)

// MarshalJSON marshals a value to JSON using the provided name.
// This is a helper to reduce boilerplate for simple value types.
func MarshalJSON[T any](name string, value T) ([]byte, error) {
	b, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("%s: marshal JSON: %w", name, err)
	}

	return b, nil
}

// UnmarshalJSON unmarshals JSON data into a target value using the provided name.
// This is a helper to reduce boilerplate for simple value types.
func UnmarshalJSON[T any](name string, data []byte, target *T) error {
	err := json.Unmarshal(data, target)
	if err != nil {
		return fmt.Errorf("%s: invalid JSON %q: %w", name, string(data), err)
	}

	return nil
}
