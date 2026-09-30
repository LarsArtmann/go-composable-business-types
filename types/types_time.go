package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

var (
	errDurationScanNil    = errors.New("duration: scan: receiver is nil")
	errDurationCannotScan = errors.New("duration: cannot scan non-numeric/string value")
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
	return compare(d.Duration, other.Duration)
}

// setFrom parses a textual duration ("300ms", "1h30m") and stores it.
// Empty text resets the duration to zero. source names the value's origin
// ("string", "[]byte", "JSON") for error context.
func (d *Duration) setFrom(value, source string) error {
	if value == "" {
		d.Duration = 0

		return nil
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fmt.Errorf("duration: cannot parse %q from %s: %w", value, source, err)
	}

	d.Duration = parsed

	return nil
}

// Scan implements sql.Scanner for Duration.
// Supports int64 (nanoseconds), float64, string (parseable duration), and []byte sources.
func (d *Duration) Scan(src any) error {
	if d == nil {
		return errDurationScanNil
	}

	switch v := src.(type) {
	case nil:
		return d.setFrom("", "SQL NULL")
	case int64:
		d.Duration = time.Duration(v)

		return nil
	case float64:
		d.Duration = time.Duration(int64(v))

		return nil
	case string:
		return d.setFrom(v, "string")
	case []byte:
		return d.setFrom(string(v), "[]byte")
	default:
		return fmt.Errorf("%w: got %T", errDurationCannotScan, src)
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
	return MarshalJSON("duration", d.String())
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *Duration) UnmarshalJSON(data []byte) error {
	var s string

	if err := UnmarshalJSON("duration", data, &s); err != nil {
		return err
	}

	return d.setFrom(s, "JSON")
}
