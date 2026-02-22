package cbt

import (
	"database/sql/driver"
	"errors"

	"github.com/sixafter/nanoid"
)

// NanoId is a URL-safe, unique identifier (default 21 characters).
// Wrapped around https://github.com/sixafter/nanoid - FIPS-140 compatible, high-performance.
type NanoId struct{ value string }

const (
	// DefaultNanoIdLength is the default length for new NanoIds (21 chars = 126 bits of entropy).
	DefaultNanoIdLength = nanoid.DefaultLength
)

var (
	ErrNanoIdEmpty    = errors.New("nanoid: cannot be empty")
	ErrNanoIdTooShort = errors.New("nanoid: minimum length is 8 characters")
	ErrNanoIdTooLong  = errors.New("nanoid: maximum length is 256 characters")
	ErrNanoIdInvalid  = errors.New("nanoid: contains invalid characters")
)

// NewNanoId generates a new random NanoId with the default length (21 characters).
func NewNanoId() NanoId {
	return NewNanoIdWithLength(DefaultNanoIdLength)
}

// NewNanoIdWithLength generates a new random NanoId with a custom length.
// Panics if generation fails (should never happen with valid length).
func NewNanoIdWithLength(length int) NanoId {
	return NanoId{value: string(nanoid.MustWithLength(length))}
}

// ParseNanoId validates and creates a NanoId from a string.
// Returns an error if the string is empty, too short (<8), too long (>256),
// or contains characters outside the URL-safe alphabet.
func ParseNanoId(s string) (NanoId, error) {
	if s == "" {
		return NanoId{}, ErrNanoIdEmpty
	}
	if len(s) < 8 {
		return NanoId{}, ErrNanoIdTooShort
	}
	if len(s) > 256 {
		return NanoId{}, ErrNanoIdTooLong
	}

	for _, r := range s {
		if !isNanoIdChar(r) {
			return NanoId{}, ErrNanoIdInvalid
		}
	}

	return NanoId{value: s}, nil
}

// MustParseNanoId is like ParseNanoId but panics on error.
// Use only for hardcoded strings known at compile time.
func MustParseNanoId(s string) NanoId {
	id, err := ParseNanoId(s)
	if err != nil {
		panic(err)
	}
	return id
}

// String returns the string representation of the NanoId.
func (id NanoId) String() string { return id.value }

// IsZero returns true if the NanoId has no value.
func (id NanoId) IsZero() bool { return id.value == "" }

// GoString implements fmt.GoStringer for debugging.
func (id NanoId) GoString() string { return id.value }

// MarshalText implements encoding.TextMarshaler for JSON serialization.
func (id NanoId) MarshalText() ([]byte, error) {
	if id.IsZero() {
		return nil, nil
	}
	return []byte(id.value), nil
}

// UnmarshalText implements encoding.TextUnmarshaler for JSON deserialization.
func (id *NanoId) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		*id = NanoId{}
		return nil
	}

	parsed, err := ParseNanoId(string(data))
	if err != nil {
		return err
	}
	*id = parsed
	return nil
}

// isNanoIdChar checks if a rune is a valid NanoId character.
func isNanoIdChar(r rune) bool {
	return (r >= 'A' && r <= 'Z') ||
		(r >= 'a' && r <= 'z') ||
		(r >= '0' && r <= '9') ||
		r == '-' || r == '_'
}

// Scan implements sql.Scanner for database deserialization.
// Supports string and []byte sources. Empty string/nil results in zero value.
func (id *NanoId) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*id = NanoId{}
		return nil
	case string:
		if v == "" {
			*id = NanoId{}
			return nil
		}
		parsed, err := ParseNanoId(v)
		if err != nil {
			return err
		}
		*id = parsed
		return nil
	case []byte:
		if len(v) == 0 {
			*id = NanoId{}
			return nil
		}
		parsed, err := ParseNanoId(string(v))
		if err != nil {
			return err
		}
		*id = parsed
		return nil
	default:
		return errors.New("nanoid: cannot scan non-string/[]byte value")
	}
}

// Value implements driver.Valuer for database serialization.
// Returns nil for empty NanoId, otherwise the string value.
func (id NanoId) Value() (driver.Value, error) {
	if id.IsZero() {
		return nil, nil
	}
	return id.value, nil
}
