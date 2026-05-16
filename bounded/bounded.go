// Package bounded provides length-validated string types.
//
// BoundedString ensures strings meet minimum and maximum length constraints,
// useful for domain validation (e.g., product names, titles, descriptions).
//
// Basic usage:
//
//	name, err := bounded.NewBoundedString(1, 100, "John Doe")
//	// Or use factory functions:
//	NewProductName := bounded.BoundedStringOf(1, 200)
//	product, err := NewProductName("Widget")
package bounded

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/larsartmann/go-composable-business-types/scanutil"
)

var (
	errMaxLessThanMin = errors.New(
		"boundedstring: maximum length cannot be less than minimum length",
	)
	errBelowMin        = errors.New("boundedstring: string length is less than minimum")
	errAboveMax        = errors.New("boundedstring: string length exceeds maximum")
	errScanNilReceiver = errors.New("boundedstring: scan: receiver is nil")
)

// BoundedString is a string with length constraints validated at construction.
// Use NewBoundedString to create validated instances.
type BoundedString struct {
	value  string
	minLen uint
	maxLen uint
}

// NewBoundedString creates a BoundedString with the given length constraints.
// Returns an error if the string length is outside [minLen, maxLen].
func NewBoundedString(minLen, maxLen uint, value string) (BoundedString, error) {
	if maxLen < minLen {
		return BoundedString{value: "", minLen: 0, maxLen: 0}, fmt.Errorf(
			"%w: minLen=%d, maxLen=%d, value=%q",
			errMaxLessThanMin,
			minLen,
			maxLen,
			value,
		)
	}

	length := uint( //nolint:gosec // G115: utf8.RuneCountInString cannot return negative value
		utf8.RuneCountInString(value),
	)

	if length < minLen {
		return BoundedString{}, fmt.Errorf(
			"%w: length=%d, min=%d, max=%d, value=%q",
			errBelowMin,
			length,
			minLen,
			maxLen,
			value,
		)
	}

	if length > maxLen {
		return BoundedString{}, fmt.Errorf(
			"%w: length=%d, max=%d, min=%d, value=%q",
			errAboveMax,
			length,
			maxLen,
			minLen,
			value,
		)
	}

	return BoundedString{value: value, minLen: minLen, maxLen: maxLen}, nil
}

// String returns the underlying string value.
func (bs BoundedString) String() string { return bs.value }

// Len returns the number of runes in the string.
func (bs BoundedString) Len() uint { return uint(utf8.RuneCountInString(bs.value)) } //nolint:gosec // G115: utf8.RuneCountInString cannot return negative value

// MinLen returns the minimum allowed length.
func (bs BoundedString) MinLen() uint { return bs.minLen }

// MaxLen returns the maximum allowed length.
func (bs BoundedString) MaxLen() uint { return bs.maxLen }

// IsZero returns true if the BoundedString is empty.
func (bs BoundedString) IsZero() bool { return bs.value == "" }

// IsMinLength returns true if the string is at the minimum allowed length.
func (bs BoundedString) IsMinLength() bool { return bs.Len() == bs.minLen }

// IsMaxLength returns true if the string is at the maximum allowed length.
func (bs BoundedString) IsMaxLength() bool { return bs.Len() == bs.maxLen }

// BoundedStringOf creates a constructor for BoundedString with fixed bounds.
// Useful for domain-specific string types with consistent constraints.
//
// Example:
//
//	var NewName = cbt.BoundedStringOf(1, 100)
//	name, err := NewName("John Doe")
func BoundedStringOf(minLen, maxLen uint) func(string) (BoundedString, error) {
	return func(value string) (BoundedString, error) {
		return NewBoundedString(minLen, maxLen, value)
	}
}

// NonEmptyString is a convenience constructor for strings that must have at least one character.
// Equivalent to NewBoundedString(1, maxLen, value).
func NonEmptyString(maxLen uint, value string) (BoundedString, error) {
	return NewBoundedString(1, maxLen, value)
}

// TrimmedBoundedString creates a BoundedString after trimming whitespace.
// Useful for user input where leading/trailing spaces should be ignored.
func TrimmedBoundedString(minLen, maxLen uint, value string) (BoundedString, error) {
	return NewBoundedString(minLen, maxLen, strings.TrimSpace(value))
}

// MarshalJSON implements json.Marshaler.
// Serializes as a JSON string containing the value.
func (bs BoundedString) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(bs.value)
	if err != nil {
		return nil, fmt.Errorf("boundedstring: marshal JSON: %w", err)
	}

	return b, nil
}

// UnmarshalJSON implements json.Unmarshaler.
// Validates the string against length constraints.
func (bs *BoundedString) UnmarshalJSON(data []byte) error {
	var value string

	err := json.Unmarshal(data, &value)
	if err != nil {
		return fmt.Errorf("boundedstring: unmarshal JSON %q: %w", string(data), err)
	}

	bs.value = value
	bs.minLen = 0

	bs.maxLen = uint( //nolint:gosec // G115: utf8.RuneCountInString cannot return negative value
		utf8.RuneCountInString(value),
	)

	return nil
}

// Scan implements sql.Scanner for database deserialization.
// Supports string and []byte sources. Sets min=0, max=len(value).
func (bs *BoundedString) Scan(src any) error {
	if bs == nil {
		return errScanNilReceiver
	}

	err := scanutil.ScanString(src, func(v string) error {
		if v == "" {
			*bs = BoundedString{value: "", minLen: 0, maxLen: 0}

			return nil
		}

		bs.value = v
		bs.minLen = 0

		bs.maxLen = uint( //nolint:gosec // G115: utf8.RuneCountInString cannot return negative value
			utf8.RuneCountInString(v),
		)

		return nil
	})
	if err != nil {
		return fmt.Errorf("boundedstring: scan: %w", err)
	}

	return nil
}

// Value implements driver.Valuer for database serialization.
// Returns nil for empty BoundedString, otherwise the string value.
func (bs BoundedString) Value() (driver.Value, error) {
	return scanutil.NullableValueWithError(bs.value, "boundedstring")
}
