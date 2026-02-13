package cbt

import (
	"crypto/rand"
	"errors"
)

// NanoId is a URL-safe, unique identifier (default 21 characters).
// Inspired by https://github.com/ai/nanoid - shorter and more readable than UUID.
type NanoId struct{ value string }

const (
	// nanoIdAlphabet contains URL-safe characters used for NanoId generation.
	nanoIdAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
	// DefaultNanoIdLength is the default length for new NanoIds (21 chars = 126 bits of entropy).
	DefaultNanoIdLength = 21
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
// Panics if length < 1 (programming error).
func NewNanoIdWithLength(length int) NanoId {
	if length < 1 {
		panic("nanoid: length must be at least 1")
	}

	bytes := make([]byte, length)
	// crypto/rand.Read is guaranteed to not return an error on supported platforms
	// and will panic if the system's PRNG fails. This is the desired behavior
	// for identifier generation - we cannot continue without randomness.
	_, _ = rand.Read(bytes)

	for i, b := range bytes {
		bytes[i] = nanoIdAlphabet[b%byte(len(nanoIdAlphabet))]
	}

	return NanoId{value: string(bytes)}
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

// IsEmpty returns true if the NanoId has no value.
func (id NanoId) IsEmpty() bool { return id.value == "" }

// GoString implements fmt.GoStringer for debugging.
func (id NanoId) GoString() string { return id.value }

// MarshalText implements encoding.TextMarshaler for JSON serialization.
func (id NanoId) MarshalText() ([]byte, error) {
	if id.IsEmpty() {
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
