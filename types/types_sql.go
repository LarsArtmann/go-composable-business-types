package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"github.com/larsartmann/go-composable-business-types/scanutil"
	"github.com/larsartmann/go-composable-business-types/validate"
)

// scanStringType implements sql.Scanner for string-based types.
// Handles nil receiver check, empty string (zero value), and type-specific parsing.
func scanStringType[T ~string](ptr *T, name string, src any, parse func(string) (T, error)) error {
	if ptr == nil {
		return fmt.Errorf("%s: scan: receiver is nil", name)
	}

	err := scanutil.ScanString(src, func(value string) error {
		if value == "" {
			*ptr = *new(T)

			return nil
		}

		parsed, parseErr := parse(value)
		if parseErr != nil {
			return fmt.Errorf("%s: invalid value %q: %w", name, value, parseErr)
		}

		*ptr = parsed

		return nil
	})
	if err != nil {
		return fmt.Errorf("%s: scan: %w", name, err)
	}

	return nil
}

// valueStringType implements driver.Valuer for string-based types.
// Returns nil for zero values, otherwise the string value.
func valueStringType[T ~string](v T, name string) (driver.Value, error) {
	val, err := scanutil.NullableValue(string(v))
	if err != nil {
		return nil, fmt.Errorf("%s: value: %w", name, err)
	}

	return val, nil
}

// valueInt64Type implements driver.Valuer for int64-based types.
// Returns nil for zero values, otherwise the int64 value.
func valueInt64Type(v int64, name string) (driver.Value, error) {
	val, err := scanutil.Int64Value(v)
	if err != nil {
		return nil, fmt.Errorf("%s: value: %w", name, err)
	}

	return val, nil
}

// Scan implements sql.Scanner for Email.
// Supports string and []byte sources. Empty string/nil results in zero value.
func (e *Email) Scan(src any) error {
	return scanStringType(e, "email", src, NewEmail)
}

// Value implements driver.Valuer for Email.
// Returns nil for empty Email, otherwise the string value.
func (e Email) Value() (driver.Value, error) {
	return valueStringType(e, "email")
}

// Scan implements sql.Scanner for URL.
// Supports string and []byte sources. Empty string/nil results in zero value.
func (u *URL) Scan(src any) error {
	return scanStringType(u, "url", src, NewURL)
}

// Value implements driver.Valuer for URL.
// Returns nil for empty URL, otherwise the string value.
func (u URL) Value() (driver.Value, error) {
	return valueStringType(u, "url")
}

// scanInt64Type implements sql.Scanner for int64-based types.
// Handles nil receiver check, int64 scanning, and type-specific conversion.
func scanInt64Type[T any](ptr *T, name string, src any, convert func(int64) T) error {
	if ptr == nil {
		return fmt.Errorf("%s: scan: receiver is nil", name)
	}

	err := scanutil.ScanInt64(src, func(v int64) error {
		*ptr = convert(v)

		return nil
	})
	if err != nil {
		return fmt.Errorf("%s: scan: %w", name, err)
	}

	return nil
}

// Value implements driver.Valuer for Cents.
func (c Cents) Value() (driver.Value, error) {
	return int64(c), nil
}

// Scan implements sql.Scanner for Cents.
// Supports int64, float64, and []byte sources.
func (c *Cents) Scan(src any) error {
	return scanInt64Type(c, "cents", src, func(v int64) Cents {
		return Cents(v)
	})
}

// Scan implements sql.Scanner for Timestamp.
// Supports time.Time, string (RFC3339), and []byte sources.
func (t *Timestamp) Scan(src any) error {
	if t == nil {
		return errors.New("timestamp: scan: receiver is nil")
	}

	switch v := src.(type) {
	case nil:
		t.Time = time.Time{}

		return nil
	case time.Time:
		t.Time = v

		return nil
	case string:
		parsed, err := time.Parse(time.RFC3339Nano, v)
		if err != nil {
			return fmt.Errorf("timestamp: cannot parse %q from string: %w", v, err)
		}

		t.Time = parsed

		return nil
	case []byte:
		parsed, err := time.Parse(time.RFC3339Nano, string(v))
		if err != nil {
			return fmt.Errorf("timestamp: cannot parse %q from []byte: %w", string(v), err)
		}

		t.Time = parsed

		return nil
	default:
		return fmt.Errorf("timestamp: cannot scan value (got %T)", src)
	}
}

// Value implements driver.Valuer for Timestamp.
func (t Timestamp) Value() (driver.Value, error) {
	if t.Time.IsZero() {
		return nil, nil
	}

	return t.Time, nil
}

// Compile-time interface assertions to ensure types implement validate.Validator.
var (
	_ validate.Validator = Email("")
	_ validate.Validator = URL("")
	_ validate.Validator = Cents(0)
	_ validate.Validator = Percentage(0)
)
