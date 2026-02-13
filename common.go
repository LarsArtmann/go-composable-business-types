package cbt

import (
	"errors"
	"net/mail"
	"net/url"
	"regexp"
	"time"
)

var (
	// emailRegex matches most common valid email formats.
	// Not exhaustive RFC 5322 compliance, but practical validation.
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

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

func (e Email) String() string { return string(e) }

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

// Percentage represents a value from 0-100.
type Percentage uint8

func NewPercentage(v uint8) Percentage {
	if v > 100 {
		return 100
	}
	return Percentage(v)
}
func (p Percentage) Float64() float64 { return float64(p) / 100 }

// Cents represents monetary amounts in smallest currency unit (prevents float errors).
type Cents int64

func NewCents(v int64) Cents     { return Cents(v) }
func (c Cents) Int64() int64     { return int64(c) }
func (c Cents) Float64() float64 { return float64(c) / 100 }

func (c Cents) Add(other Cents) Cents      { return c + other }
func (c Cents) Sub(other Cents) Cents      { return c - other }
func (c Cents) Mul(n int64) Cents          { return c * Cents(n) }
func (c Cents) Div(n int64) Cents          { return c / Cents(n) }
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
func (c Cents) IsZero() bool { return c == 0 }
func (c Cents) IsPositive() bool { return c > 0 }
func (c Cents) IsNegative() bool { return c < 0 }

// Timestamp wraps time.Time for domain clarity.
type Timestamp struct{ time.Time }

func NewTimestamp(t time.Time) Timestamp { return Timestamp{Time: t} }
func Now() Timestamp                     { return Timestamp{Time: time.Now()} }

// Duration wraps time.Duration for domain clarity.
type Duration struct{ time.Duration }

func NewDuration(d time.Duration) Duration { return Duration{Duration: d} }
