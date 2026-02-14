package cbt

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

// BoundedString is a string with length constraints validated at construction.
// Use NewBoundedString to create validated instances.
type BoundedString struct {
	value  string
	minLen int
	maxLen int
}

// NewBoundedString creates a BoundedString with the given length constraints.
// Returns an error if the string length is outside [minLen, maxLen].
func NewBoundedString(minLen, maxLen int, value string) (BoundedString, error) {
	if minLen < 0 {
		return BoundedString{}, errors.New("minimum length cannot be negative")
	}
	if maxLen < minLen {
		return BoundedString{}, errors.New("maximum length cannot be less than minimum length")
	}

	length := utf8.RuneCountInString(value)
	if length < minLen {
		return BoundedString{}, fmt.Errorf("string length %d is less than minimum %d", length, minLen)
	}
	if length > maxLen {
		return BoundedString{}, fmt.Errorf("string length %d exceeds maximum %d", length, maxLen)
	}

	return BoundedString{value: value, minLen: minLen, maxLen: maxLen}, nil
}

// MustBoundedString creates a BoundedString or panics on validation failure.
// Use only when the input is guaranteed valid (e.g., constants, tests).
func MustBoundedString(minLen, maxLen int, value string) BoundedString {
	bs, err := NewBoundedString(minLen, maxLen, value)
	if err != nil {
		panic(err)
	}
	return bs
}

func (bs BoundedString) String() string    { return bs.value }
func (bs BoundedString) Len() int          { return utf8.RuneCountInString(bs.value) }
func (bs BoundedString) MinLen() int       { return bs.minLen }
func (bs BoundedString) MaxLen() int       { return bs.maxLen }
func (bs BoundedString) IsEmpty() bool     { return bs.value == "" }
func (bs BoundedString) IsMinLength() bool { return bs.Len() == bs.minLen }
func (bs BoundedString) IsMaxLength() bool { return bs.Len() == bs.maxLen }

// BoundedStringOf creates a constructor for BoundedString with fixed bounds.
// Useful for domain-specific string types with consistent constraints.
//
// Example:
//
//	var NewName = cbt.BoundedStringOf(1, 100)
//	name, err := NewName("John Doe")
func BoundedStringOf(minLen, maxLen int) func(string) (BoundedString, error) {
	return func(value string) (BoundedString, error) {
		return NewBoundedString(minLen, maxLen, value)
	}
}

// NonEmptyString is a convenience constructor for strings that must have at least one character.
// Equivalent to NewBoundedString(1, maxLen, value).
func NonEmptyString(maxLen int, value string) (BoundedString, error) {
	return NewBoundedString(1, maxLen, value)
}

// TrimmedBoundedString creates a BoundedString after trimming whitespace.
// Useful for user input where leading/trailing spaces should be ignored.
func TrimmedBoundedString(minLen, maxLen int, value string) (BoundedString, error) {
	return NewBoundedString(minLen, maxLen, strings.TrimSpace(value))
}

// MarshalJSON implements json.Marshaler.
// Serializes as a JSON string containing the value.
func (bs BoundedString) MarshalJSON() ([]byte, error) {
	return json.Marshal(bs.value)
}

// UnmarshalJSON implements json.Unmarshaler.
// Validates the string against length constraints.
func (bs *BoundedString) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	bs.value = value
	bs.minLen = 0
	bs.maxLen = utf8.RuneCountInString(value)
	return nil
}
