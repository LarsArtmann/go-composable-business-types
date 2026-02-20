package cbt

import (
	"database/sql/driver"
	"errors"
	"net/mail"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// emailRegex matches most common valid email formats.
// Not exhaustive RFC 5322 compliance, but practical validation.
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Email represents a validated email address.
type Email string

// ErrInvalidEmail is returned when an email address fails validation.
var ErrInvalidEmail = errors.New("invalid email address")

// NewEmail creates a validated Email. Returns ErrInvalidEmail if validation fails.
func NewEmail(v string) (Email, error) {
	if v == "" {
		return "", ErrInvalidEmail
	}
	// mail.ParseAddress handles RFC 5322 parsing (including display names)
	addr, err := mail.ParseAddress(v)
	if err != nil {
		return "", ErrInvalidEmail
	}
	// Extract just the email address (strip any display name)
	email := addr.Address
	// Additional format check for common patterns
	if !emailRegex.MatchString(email) {
		return "", ErrInvalidEmail
	}
	return Email(email), nil
}

// MustParseEmail panics if the email is invalid. Use for compile-time known constants.
func MustParseEmail(v string) Email {
	e, err := NewEmail(v)
	if err != nil {
		panic(err)
	}
	return e
}

func (e Email) String() string    { return string(e) }
func (e Email) IsZero() bool      { return e == "" }
func (e Email) LocalPart() string { s, _, _ := e.split(); return s }
func (e Email) Domain() string    { _, d, _ := e.split(); return d }

// Normalize returns an email with normalized case.
// Per RFC 1035, domain names are case-insensitive, so the domain is lowercased.
func (e Email) Normalize() Email {
	local, domain, ok := e.split()
	if !ok {
		return e
	}
	return Email(local + "@" + strings.ToLower(domain))
}

// split returns local part, domain, and whether the split was successful.
func (e Email) split() (local, domain string, ok bool) {
	if e == "" {
		return "", "", false
	}
	idx := -1
	for i := 0; i < len(e); i++ {
		if e[i] == '@' {
			idx = i
			break
		}
	}
	if idx <= 0 || idx >= len(e)-1 {
		return "", "", false
	}
	return string(e[:idx]), string(e[idx+1:]), true
}

// URL represents a validated URL with http or https scheme.
type URL string

// ErrInvalidURL is returned when a URL fails validation.
var ErrInvalidURL = errors.New("invalid URL")

// NewURL creates a validated URL. Returns ErrInvalidURL if validation fails.
// Requires http or https scheme and a valid host.
func NewURL(v string) (URL, error) {
	if v == "" {
		return "", ErrInvalidURL
	}
	parsed, err := url.Parse(v)
	if err != nil {
		return "", ErrInvalidURL
	}
	// Require scheme and host
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", ErrInvalidURL
	}
	if parsed.Host == "" {
		return "", ErrInvalidURL
	}
	return URL(v), nil
}

// MustParseURL panics if the URL is invalid. Use for compile-time known constants.
func MustParseURL(v string) URL {
	u, err := NewURL(v)
	if err != nil {
		panic(err)
	}
	return u
}

func (u URL) String() string { return string(u) }

// IsZero returns true if the URL is empty.
func (u URL) IsZero() bool { return u == "" }

// Parse returns the underlying url.URL. Since URLs are validated at construction,
// this should never fail, but returns error for safety.
func (u URL) Parse() (*url.URL, error) { return url.Parse(string(u)) }

// Scheme returns the URL scheme (http or https). Returns empty string if URL is zero.
func (u URL) Scheme() string {
	if u == "" {
		return ""
	}
	for i := 0; i < len(u); i++ {
		if u[i] == ':' {
			return string(u[:i])
		}
	}
	return ""
}

// Host returns the URL host (e.g., "example.com" or "example.com:8080").
func (u URL) Host() string {
	parsed, _ := u.Parse()
	if parsed == nil {
		return ""
	}
	return parsed.Host
}

// Path returns the URL path (e.g., "/api/v1/users").
func (u URL) Path() string {
	parsed, _ := u.Parse()
	if parsed == nil {
		return ""
	}
	return parsed.Path
}

// Percentage represents a value from 0-100.
type Percentage uint8

func NewPercentage(v uint8) Percentage {
	if v > 100 {
		return 100
	}
	return Percentage(v)
}
func (p Percentage) Float64() float64 { return float64(p) / 100 }

// IsZero returns true if the percentage is 0.
func (p Percentage) IsZero() bool { return p == 0 }

// IsMin returns true if the percentage is 0 (minimum value).
func (p Percentage) IsMin() bool { return p == 0 }

// IsMax returns true if the percentage is 100 (maximum value).
func (p Percentage) IsMax() bool { return p == 100 }

// Cents represents monetary amounts in smallest currency unit (prevents float errors).
type Cents int64

func NewCents(v int64) Cents     { return Cents(v) }
func (c Cents) Int64() int64     { return int64(c) }
func (c Cents) Float64() float64 { return float64(c) / 100 }

func (c Cents) Add(other Cents) Cents { return c + other }
func (c Cents) Sub(other Cents) Cents { return c - other }
func (c Cents) Mul(n int64) Cents     { return c * Cents(n) }
func (c Cents) Div(n int64) Cents     { return c / Cents(n) }
func (c Cents) Abs() Cents {
	if c < 0 {
		return -c
	}
	return c
}

func (c Cents) Sign() int {
	if c < 0 {
		return -1
	}
	if c > 0 {
		return 1
	}
	return 0
}
func (c Cents) IsZero() bool     { return c == 0 }
func (c Cents) IsPositive() bool { return c > 0 }
func (c Cents) IsNegative() bool { return c < 0 }

// Compare returns -1 if c < other, 0 if equal, 1 if c > other.
func (c Cents) Compare(other Cents) int {
	if c < other {
		return -1
	}
	if c > other {
		return 1
	}
	return 0
}

// Timestamp wraps time.Time for domain clarity.
type Timestamp struct{ time.Time }

func NewTimestamp(t time.Time) Timestamp { return Timestamp{Time: t} }
func Now() Timestamp                     { return Timestamp{Time: time.Now()} }

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
			return errors.New("duration: cannot scan value")
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
			return errors.New("duration: cannot scan value")
		}
		d.Duration = parsed
		return nil
	default:
		return errors.New("duration: cannot scan non-numeric/string value")
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

// Scan implements sql.Scanner for Email.
// Supports string and []byte sources. Empty string/nil results in zero value.
func (e *Email) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*e = ""
		return nil
	case string:
		if v == "" {
			*e = ""
			return nil
		}
		parsed, err := NewEmail(v)
		if err != nil {
			return err
		}
		*e = parsed
		return nil
	case []byte:
		if len(v) == 0 {
			*e = ""
			return nil
		}
		parsed, err := NewEmail(string(v))
		if err != nil {
			return err
		}
		*e = parsed
		return nil
	default:
		return ErrInvalidEmail
	}
}

// Value implements driver.Valuer for Email.
// Returns nil for empty Email, otherwise the string value.
func (e Email) Value() (driver.Value, error) {
	if e.IsZero() {
		return nil, nil
	}
	return string(e), nil
}

// Scan implements sql.Scanner for URL.
// Supports string and []byte sources. Empty string/nil results in zero value.
func (u *URL) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*u = ""
		return nil
	case string:
		if v == "" {
			*u = ""
			return nil
		}
		parsed, err := NewURL(v)
		if err != nil {
			return err
		}
		*u = parsed
		return nil
	case []byte:
		if len(v) == 0 {
			*u = ""
			return nil
		}
		parsed, err := NewURL(string(v))
		if err != nil {
			return err
		}
		*u = parsed
		return nil
	default:
		return ErrInvalidURL
	}
}

// Value implements driver.Valuer for URL.
// Returns nil for empty URL, otherwise the string value.
func (u URL) Value() (driver.Value, error) {
	if u.IsZero() {
		return nil, nil
	}
	return string(u), nil
}

// Scan implements sql.Scanner for Cents.
// Supports int64, float64, and []byte sources.
func (c *Cents) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*c = 0
		return nil
	case int64:
		*c = Cents(v)
		return nil
	case float64:
		*c = Cents(int64(v))
		return nil
	case []byte:
		val, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return errors.New("cents: cannot scan value")
		}
		*c = Cents(val)
		return nil
	default:
		return errors.New("cents: cannot scan non-numeric value")
	}
}

// Value implements driver.Valuer for Cents.
func (c Cents) Value() (driver.Value, error) {
	return int64(c), nil
}

// Scan implements sql.Scanner for Timestamp.
// Supports time.Time, string (RFC3339), and []byte sources.
func (t *Timestamp) Scan(src any) error {
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
			return err
		}
		t.Time = parsed
		return nil
	case []byte:
		parsed, err := time.Parse(time.RFC3339Nano, string(v))
		if err != nil {
			return err
		}
		t.Time = parsed
		return nil
	default:
		return errors.New("timestamp: cannot scan value")
	}
}

// Value implements driver.Valuer for Timestamp.
func (t Timestamp) Value() (driver.Value, error) {
	if t.Time.IsZero() {
		return nil, nil
	}
	return t.Time, nil
}
