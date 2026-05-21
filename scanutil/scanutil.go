// Package scanutil provides reusable utilities for sql.Scanner and driver.Valuer implementations.
//
// This package reduces boilerplate by providing generic helpers for common patterns:
//   - Nullable string scanning (nil, string, []byte → *string)
//   - Nullable value returning (nil for empty strings)
//
// Example usage:
//
//	func (s *MyString) Scan(src any) error {
//		return scanutil.ScanString(src, func(v string) error {
//			// Custom validation
//			if v == "" {
//				*s = MyString{}
//				return nil
//			}
//			*s = MyString{value: v}
//			return nil
//		})
//	}
package scanutil

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	errCannotScanString = errors.New("scanutil: cannot scan into string value")
	errCannotScanInt64  = errors.New("scanutil: cannot scan into int64 value")
)

// ScanString scans src into a string callback, supporting nil, string, and []byte.
// This is the most common Scan pattern across the codebase.
//
// Usage:
//
//	err := scanutil.ScanString(src, func(v string) error {
//		// Custom logic
//		return nil
//	})
func ScanString(src any, setValue func(string) error) error {
	switch v := src.(type) {
	case nil:
		return setValue("")
	case string:
		return setValue(v)
	case []byte:
		return setValue(string(v))
	default:
		return fmt.Errorf("%w: got %T", errCannotScanString, src)
	}
}

// ScanInt64 scans src into an int64 callback, supporting nil, int64, int, float64, and []byte.
//
// Usage:
//
//	err := scanutil.ScanInt64(src, func(v int64) error {
//		// Custom logic
//		return nil
//	})
func ScanInt64(src any, setValue func(int64) error) error {
	if src == nil {
		return setValue(0)
	}

	switch v := src.(type) {
	case int64:
		return setValue(v)
	case int:
		return setValue(int64(v))
	case float64:
		return setValue(int64(v))
	case []byte:
		parsed, err := parseIntFromBytes(v)
		if err != nil {
			return fmt.Errorf("%w: got %T (len=%d): %w", errCannotScanInt64, src, len(v), err)
		}

		return setValue(parsed)
	default:
		return fmt.Errorf("%w: got %T", errCannotScanInt64, src)
	}
}

// parseIntFromBytes parses an int64 from a byte slice.
func parseIntFromBytes(data []byte) (int64, error) {
	if len(data) == 0 {
		return 0, nil
	}

	v, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("scanutil: parse int from bytes: %w", err)
	}

	return v, nil
}

// NullableValue returns nil for empty strings, otherwise returns the string as driver.Value.
// Use this in Value() implementations to get SQL NULL for empty values.
func NullableValue(v string) (any, error) {
	if v == "" {
		return nil, nil
	}

	return v, nil
}

// NullableValueWithError wraps NullableValue with error context for driver.Valuer implementations.
func NullableValueWithError(v, errorPrefix string) (any, error) {
	val, err := NullableValue(v)
	if err != nil {
		return nil, fmt.Errorf("%s: value: %w", errorPrefix, err)
	}

	return val, nil
}

// NonNullableValue returns the string as-is (for types that don't support null in SQL).
func NonNullableValue(v string) (any, error) {
	return v, nil
}

// Int64Value returns int64 as driver.Value.
func Int64Value(v int64) (any, error) {
	return v, nil
}

// ZeroAsNullValue returns nil for zero values, otherwise returns the int64.
func ZeroAsNullValue(v int64) (any, error) {
	if v == 0 {
		return nil, nil
	}

	return v, nil
}

var errNilPtr = errors.New("scanutil: value pointer is nil")

// ScanEnum provides a generic sql.Scanner implementation for iota-based enum types.
// parseFunc is the generated ParseX function (e.g. ParseActorKind).
//
// Usage:
//
//	func (x *ActorKind) Scan(value interface{}) error {
//	    return scanutil.ScanEnum(x, value, ParseActorKind)
//	}
func ScanEnum[T ~uint8](ptr *T, src any, parseFunc func(string) (T, error)) error {
	if src == nil {
		*ptr = T(0)

		return nil
	}

	var err error

	switch v := src.(type) {
	case int64:
		*ptr = T(v)
	case string:
		*ptr, err = parseFunc(v)
	case []byte:
		*ptr, err = parseFunc(string(v))
	case T:
		*ptr = v
	case int:
		*ptr = T(v)
	case *T:
		if v == nil {
			return errNilPtr
		}

		*ptr = *v
	case uint:
		*ptr = T(v)
	case uint64:
		*ptr = T(v)
	case *int:
		if v == nil {
			return errNilPtr
		}

		*ptr = T(*v)
	case *int64:
		if v == nil {
			return errNilPtr
		}

		*ptr = T(*v)
	case float64:
		*ptr = T(v)
	case *float64:
		if v == nil {
			return errNilPtr
		}

		*ptr = T(*v)
	case *uint:
		if v == nil {
			return errNilPtr
		}

		*ptr = T(*v)
	case *uint64:
		if v == nil {
			return errNilPtr
		}

		*ptr = T(*v)
	case *string:
		if v == nil {
			return errNilPtr
		}

		*ptr, err = parseFunc(*v)
	default:
		return fmt.Errorf("scanutil: cannot scan %T into %T", src, *ptr)
	}

	return err
}
