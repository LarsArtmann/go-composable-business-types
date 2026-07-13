// Package tag provides validated string labels with length and character constraints.
//
// Tag values are alphanumeric strings (plus hyphens) up to 50 characters, suitable
// for categorization and filtering. Tags implement JSON, SQL, and validate.Validator.
//
// Basic usage:
//
//	t, err := tag.New("priority-high")
//	tags, err := tag.NewTags("frontend", "urgent")
package tag

import (
	"database/sql/driver"
	"encoding/json/v2"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"unicode/utf8"

	"github.com/larsartmann/go-composable-business-types/scanutil"
	"github.com/larsartmann/go-composable-business-types/validate"
)

const (
	maxLen = 50
)

var (
	validPattern       = regexp.MustCompile(`^[A-Za-z0-9-]+$`)
	errTagEmpty        = errors.New("tag: cannot be empty")
	errTagTooLong      = errors.New("tag: length exceeds maximum")
	errTagInvalidChars = errors.New("tag: contains invalid characters")
	errTagScanNil      = errors.New("tag: scan: receiver is nil")
)

// Tag is a validated string label with length and character constraints.
type Tag string

// New creates a validated Tag.
func New(s string) (Tag, error) {
	if s == "" {
		return "", errTagEmpty
	}

	if utf8.RuneCountInString(s) > maxLen {
		return "", fmt.Errorf("%w: %d > %d", errTagTooLong, utf8.RuneCountInString(s), maxLen)
	}

	if !validPattern.MatchString(s) {
		return "", fmt.Errorf("%w: %q (allowed: A-Z, a-z, 0-9, hyphen)", errTagInvalidChars, s)
	}

	return Tag(s), nil
}

// Must creates a Tag or panics on validation failure.
func Must(s string) Tag {
	t, err := New(s)
	if err != nil {
		panic(err)
	}

	return t
}

// NewTags creates a slice of validated Tags.
func NewTags(ss ...string) ([]Tag, error) {
	tags := make([]Tag, 0, len(ss))
	for _, s := range ss {
		t, err := New(s)
		if err != nil {
			return nil, fmt.Errorf("tag at index %d: %w", len(tags), err)
		}

		tags = append(tags, t)
	}

	return tags, nil
}

// String returns the raw tag string.
func (t Tag) String() string { return string(t) }

// IsZero reports whether the tag is empty.
func (t Tag) IsZero() bool { return t == "" }

// IsValid reports whether the tag satisfies all validation rules.
func (t Tag) IsValid() bool {
	if t == "" {
		return false
	}

	if utf8.RuneCountInString(string(t)) > maxLen {
		return false
	}

	return validPattern.MatchString(string(t))
}

// Validate returns an error if the tag is invalid.
func (t Tag) Validate() error {
	if t == "" {
		return errTagEmpty
	}

	if utf8.RuneCountInString(string(t)) > maxLen {
		return fmt.Errorf("%w: %d > %d", errTagTooLong, utf8.RuneCountInString(string(t)), maxLen)
	}

	if !validPattern.MatchString(string(t)) {
		return fmt.Errorf("%w: %q", errTagInvalidChars, string(t))
	}

	return nil
}

// MarshalJSON encodes the tag as a JSON string.
func (t Tag) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(string(t))
	if err != nil {
		return nil, fmt.Errorf("tag: marshal JSON: %w", err)
	}

	return b, nil
}

// UnmarshalJSON decodes a JSON string into the tag.
func (t *Tag) UnmarshalJSON(data []byte) error {
	var s string

	err := json.Unmarshal(data, &s)
	if err != nil {
		return fmt.Errorf("tag: invalid JSON %q: %w", string(data), err)
	}

	*t = Tag(s)

	return nil
}

// Scan implements sql.Scanner for the tag.
func (t *Tag) Scan(src any) error {
	if t == nil {
		return errTagScanNil
	}

	return scanutil.ScanString(src, func(v string) error {
		*t = Tag(v)

		return nil
	})
}

// Value implements driver.Valuer for the tag.
func (t Tag) Value() (driver.Value, error) {
	return scanutil.NullableValue(string(t))
}

var _ validate.Validator = Tag("")

// Tags is an ordered collection of Tag values.
type Tags []Tag

// NewTagsFromString creates a Tags collection from raw strings.
func NewTagsFromString(ss ...string) (Tags, error) {
	tags, err := NewTags(ss...)
	if err != nil {
		return nil, err
	}

	return Tags(tags), nil
}

// Strings returns the string representation of each tag.
func (ts Tags) Strings() []string {
	result := make([]string, len(ts))
	for i, t := range ts {
		result[i] = t.String()
	}

	return result
}

// IsEmpty reports whether the collection has no tags.
func (ts Tags) IsEmpty() bool {
	return len(ts) == 0
}

// Contains reports whether the collection contains the given tag.
func (ts Tags) Contains(t Tag) bool {
	return slices.Contains(ts, t)
}
