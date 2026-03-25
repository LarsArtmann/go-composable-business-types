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
	"fmt"
	"strconv"
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
		return fmt.Errorf("cannot scan %T into string value", src)
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
			return fmt.Errorf("cannot scan %T (len=%d) into int64 value: %w", src, len(v), err)
		}
		return setValue(parsed)
	default:
		return fmt.Errorf("cannot scan %T into int64 value", src)
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
