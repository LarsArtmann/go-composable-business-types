// Package nanoid provides URL-safe, unique identifiers.
//
// NanoID is a FIPS-140 compatible, high-performance unique ID generator
// wrapped around github.com/sixafter/nanoid. Default length is 21 characters
// (126 bits of entropy).
//
// Basic usage:
//
//	id := nanoid.NewNanoID()  // 21 chars
//	id := nanoid.NewNanoIDWithLength(32)
package nanoid

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/larsartmann/go-composable-business-types/scanutil"
	"github.com/sixafter/nanoid"
)

// NanoID is a URL-safe, unique identifier (default 21 characters).
// Wrapped around https://github.com/sixafter/nanoid - FIPS-140 compatible, high-performance.
type NanoID struct{ value string }

const (
	// DefaultNanoIDLength is the default length for new NanoIDs (21 chars = 126 bits of entropy).
	DefaultNanoIDLength = nanoid.DefaultLength

	// Minimum length constraints
	minNanoIDLength = 8   // ErrNanoIDTooShort
	maxNanoIDLength = 256 // ErrNanoIDTooLong
)

var (
	// ErrNanoIDEmpty is returned when a NanoID is empty.
	ErrNanoIDEmpty    = errors.New("nanoid: cannot be empty")
	ErrNanoIDTooShort = errors.New("nanoid: minimum length is 8 characters")
	ErrNanoIDTooLong  = errors.New("nanoid: maximum length is 256 characters")
	ErrNanoIDInvalid  = errors.New("nanoid: contains invalid characters")
)

// NewNanoID generates a new random NanoID with the default length (21 characters).
func NewNanoID() NanoID {
	return NewNanoIDWithLength(DefaultNanoIDLength)
}

// NewNanoIDWithLength generates a new random NanoID with a custom length.
// Panics if generation fails (should never happen with valid length).
func NewNanoIDWithLength(length int) NanoID {
	return NanoID{value: string(nanoid.MustWithLength(length))}
}

// ParseNanoID validates and creates a NanoID from a string.
// Returns an error if the string is empty, too short (<8), too long (>256),
// or contains characters outside the URL-safe alphabet.
func ParseNanoID(s string) (NanoID, error) {
	if s == "" {
		return NanoID{}, ErrNanoIDEmpty
	}
	if len(s) < minNanoIDLength {
		return NanoID{}, ErrNanoIDTooShort
	}
	if len(s) > maxNanoIDLength {
		return NanoID{}, ErrNanoIDTooLong
	}

	for _, r := range s {
		if !isNanoIDChar(r) {
			return NanoID{}, ErrNanoIDInvalid
		}
	}

	return NanoID{value: s}, nil
}

// String returns the string representation of the NanoID.
func (id NanoID) String() string { return id.value }

// IsZero returns true if the NanoID has no value.
func (id NanoID) IsZero() bool { return id.value == "" }

// GoString implements fmt.GoStringer for debugging.
func (id NanoID) GoString() string { return id.value }

// MarshalText implements encoding.TextMarshaler for JSON serialization.
func (id NanoID) MarshalText() ([]byte, error) {
	if id.IsZero() {
		return nil, nil
	}
	return []byte(id.value), nil
}

// UnmarshalText implements encoding.TextUnmarshaler for JSON deserialization.
func (id *NanoID) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		*id = NanoID{}
		return nil
	}

	parsed, err := ParseNanoID(string(data))
	if err != nil {
		return err
	}
	*id = parsed
	return nil
}

// isNanoIDChar checks if a rune is a valid NanoID character.
func isNanoIDChar(r rune) bool {
	return (r >= 'A' && r <= 'Z') ||
		(r >= 'a' && r <= 'z') ||
		(r >= '0' && r <= '9') ||
		r == '-' || r == '_'
}

// Scan implements sql.Scanner for database deserialization.
// Supports string and []byte sources. Empty string/nil results in zero value.
func (id *NanoID) Scan(src any) error {
	return scanutil.ScanString(src, func(v string) error {
		if v == "" {
			*id = NanoID{}
			return nil
		}
		parsed, err := ParseNanoID(v)
		if err != nil {
			return fmt.Errorf("nanoid: scan %q: %w", v, err)
		}
		*id = parsed
		return nil
	})
}

// Value implements driver.Valuer for database serialization.
// Returns nil for empty NanoID, otherwise the string value.
func (id NanoID) Value() (driver.Value, error) {
	return scanutil.NullableValue(id.value)
}
