package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"github.com/larsartmann/go-composable-business-types/scanutil"
	"github.com/larsartmann/go-composable-business-types/validate"
)

// Scan implements sql.Scanner for Email.
// Supports string and []byte sources. Empty string/nil results in zero value.
func (e *Email) Scan(src any) error {
	if e == nil {
		return errors.New("email: scan: receiver is nil")
	}
	err := scanutil.ScanString(src, func(v string) error {
		if v == "" {
			*e = ""
			return nil
		}
		parsed, err := NewEmail(v)
		if err != nil {
			return fmt.Errorf("email: invalid value %q: %w", v, err)
		}
		*e = parsed
		return nil
	})
	if err != nil {
		return fmt.Errorf("email: scan: %w", err)
	}
	return nil
}

// Value implements driver.Valuer for Email.
// Returns nil for empty Email, otherwise the string value.
func (e Email) Value() (driver.Value, error) {
	v, err := scanutil.NullableValue(string(e))
	if err != nil {
		return nil, fmt.Errorf("email: value: %w", err)
	}
	return v, nil
}

// Scan implements sql.Scanner for URL.
// Supports string and []byte sources. Empty string/nil results in zero value.
func (u *URL) Scan(src any) error {
	if u == nil {
		return errors.New("url: scan: receiver is nil")
	}
	err := scanutil.ScanString(src, func(v string) error {
		if v == "" {
			*u = ""
			return nil
		}
		parsed, err := NewURL(v)
		if err != nil {
			return fmt.Errorf("url: invalid value %q: %w", v, err)
		}
		*u = parsed
		return nil
	})
	if err != nil {
		return fmt.Errorf("url: scan: %w", err)
	}
	return nil
}

// Value implements driver.Valuer for URL.
// Returns nil for empty URL, otherwise the string value.
func (u URL) Value() (driver.Value, error) {
	v, err := scanutil.NullableValue(string(u))
	if err != nil {
		return nil, fmt.Errorf("url: value: %w", err)
	}
	return v, nil
}

// Scan implements sql.Scanner for Cents.
// Supports int64, float64, and []byte sources.
func (c *Cents) Scan(src any) error {
	if c == nil {
		return errors.New("cents: scan: receiver is nil")
	}
	err := scanutil.ScanInt64(src, func(v int64) error {
		*c = Cents(v)
		return nil
	})
	if err != nil {
		return fmt.Errorf("cents: scan: %w", err)
	}
	return nil
}

// Value implements driver.Valuer for Cents.
func (c Cents) Value() (driver.Value, error) {
	return int64(c), nil
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
