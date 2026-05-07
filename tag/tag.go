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
	"encoding/json"
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

var validPattern = regexp.MustCompile(`^[A-Za-z0-9-]+$`)

type Tag string

func New(s string) (Tag, error) {
	if s == "" {
		return "", errors.New("tag: cannot be empty")
	}

	if utf8.RuneCountInString(s) > maxLen {
		return "", fmt.Errorf(
			"tag: length %d exceeds maximum %d",
			utf8.RuneCountInString(s),
			maxLen,
		)
	}

	if !validPattern.MatchString(s) {
		return "", fmt.Errorf(
			"tag: %q contains invalid characters (allowed: A-Z, a-z, 0-9, hyphen)",
			s,
		)
	}

	return Tag(s), nil
}

func Must(s string) Tag {
	t, err := New(s)
	if err != nil {
		panic(err)
	}

	return t
}

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

func (t Tag) String() string { return string(t) }

func (t Tag) IsZero() bool { return t == "" }

func (t Tag) IsValid() bool {
	if t == "" {
		return false
	}

	if utf8.RuneCountInString(string(t)) > maxLen {
		return false
	}

	return validPattern.MatchString(string(t))
}

func (t Tag) Validate() error {
	if t == "" {
		return errors.New("tag: cannot be empty")
	}

	if utf8.RuneCountInString(string(t)) > maxLen {
		return fmt.Errorf(
			"tag: length %d exceeds maximum %d",
			utf8.RuneCountInString(string(t)),
			maxLen,
		)
	}

	if !validPattern.MatchString(string(t)) {
		return fmt.Errorf("tag: %q contains invalid characters", string(t))
	}

	return nil
}

func (t Tag) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(string(t))
	if err != nil {
		return nil, fmt.Errorf("tag: marshal JSON: %w", err)
	}

	return b, nil
}

func (t *Tag) UnmarshalJSON(data []byte) error {
	var s string

	err := json.Unmarshal(data, &s)
	if err != nil {
		return fmt.Errorf("tag: invalid JSON %q: %w", string(data), err)
	}

	*t = Tag(s)

	return nil
}

func (t *Tag) Scan(src any) error {
	if t == nil {
		return errors.New("tag: scan: receiver is nil")
	}

	return scanutil.ScanString(src, func(v string) error {
		*t = Tag(v)

		return nil
	})
}

func (t Tag) Value() (driver.Value, error) {
	return scanutil.NullableValue(string(t))
}

var _ validate.Validator = Tag("")

type Tags []Tag

func NewTagsFromString(ss ...string) (Tags, error) {
	tags := make(Tags, 0, len(ss))
	for _, s := range ss {
		t, err := New(s)
		if err != nil {
			return nil, fmt.Errorf("tag at index %d: %w", len(tags), err)
		}

		tags = append(tags, t)
	}

	return tags, nil
}

func (ts Tags) Strings() []string {
	result := make([]string, len(ts))
	for i, t := range ts {
		result[i] = t.String()
	}

	return result
}

func (ts Tags) IsEmpty() bool {
	return len(ts) == 0
}

func (ts Tags) Contains(t Tag) bool {
	return slices.Contains(ts, t)
}
