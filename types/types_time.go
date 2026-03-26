package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Timestamp wraps time.Time for domain clarity.
type Timestamp struct{ time.Time }

// NewTimestamp creates a new Timestamp from a time.Time.
func NewTimestamp(t time.Time) Timestamp { return Timestamp{Time: t} }

// Now returns the current time as a Timestamp.
func Now() Timestamp { return Timestamp{Time: time.Now()} }

// Before returns true if this timestamp is before the given time.
func (t Timestamp) Before(other time.Time) bool {
	return t.Time.Before(other)
}

// After returns true if this timestamp is after the given time.
func (t Timestamp) After(other time.Time) bool {
	return t.Time.After(other)
}

// IsZero returns true if the timestamp is the zero time.
func (t Timestamp) IsZero() bool {
	return t.Time.IsZero()
}

// Compare returns -1 if t < other, 0 if equal, 1 if t > other.
func (t Timestamp) Compare(other Timestamp) int {
	return t.Time.Compare(other.Time)
}

// Duration wraps time.Duration for domain clarity.
type Duration struct{ time.Duration }

// NewDuration creates a new Duration from a time.Duration.
func NewDuration(d time.Duration) Duration { return Duration{Duration: d} }

// IsZero returns true if the duration is zero.
func (d Duration) IsZero() bool { return d.Duration == 0 }

// Compare returns -1 if d < other, 0 if equal, 1 if d > other.
func (d Duration) Compare(other Duration) int {
	if d.Duration < other.Duration {
		return -1
	}
	if d.Duration > other.Duration {
		return 1
	}
	return 0
}

// Scan implements sql.Scanner for Duration.
// Supports int64 (nanoseconds), float64, string (parseable duration), and []byte sources.
func (d *Duration) Scan(src any) error {
	if d == nil {
		return errors.New("duration: scan: receiver is nil")
	}
	switch v := src.(type) {
	case nil:
		d.Duration = 0
		return nil
	case int64:
		d.Duration = time.Duration(v)
		return nil
	case float64:
		d.Duration = time.Duration(int64(v))
		return nil
	case string:
		if v == "" {
			d.Duration = 0
			return nil
		}
		parsed, err := time.ParseDuration(v)
		if err != nil {
			return fmt.Errorf("duration: cannot parse %q from string: %w", v, err)
		}
		d.Duration = parsed
		return nil
	case []byte:
		if len(v) == 0 {
			d.Duration = 0
			return nil
		}
		parsed, err := time.ParseDuration(string(v))
		if err != nil {
			return fmt.Errorf("duration: cannot parse %q from []byte: %w", string(v), err)
		}
		d.Duration = parsed
		return nil
	default:
		return fmt.Errorf("duration: cannot scan non-numeric/string value (got %T)", src)
	}
}

// Value implements driver.Valuer for Duration.
// Returns nil for zero duration, otherwise nanoseconds as int64.
func (d Duration) Value() (driver.Value, error) {
	if d.Duration == 0 {
		return nil, nil
	}
	return int64(d.Duration), nil
}

// MarshalJSON implements json.Marshaler.
func (d Duration) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(d.String())
	if err != nil {
		return nil, fmt.Errorf("duration: marshal JSON: %w", err)
	}
	return b, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *Duration) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("duration: invalid JSON %q: %w", string(data), err)
	}
	if s == "" {
		d.Duration = 0
		return nil
	}
	parsed, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("duration: cannot parse %q: %w", s, err)
	}
	d.Duration = parsed
	return nil
}
