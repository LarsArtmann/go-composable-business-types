// Package types provides common domain types for business applications.
//
// Types included: Email, URL, Percentage, Cents, Timestamp, Duration.
// All types include validation, SQL scanning, and JSON serialization support.
//
// Basic usage:
//
//	email, _ := types.NewEmail("user@example.com")
//	url, _ := types.NewURL("https://example.com")
//	cents := types.NewCents(1099) // $10.99
//	pct := types.NewPercentage(50) // 50%
package types

import (
	"fmt"
	"net/mail"
	"net/url"
	"regexp"
	"strings"

	pkgerrors "github.com/larsartmann/go-composable-business-types/pkg/errors"
)

// emailRegex matches most common valid email formats.
// Not exhaustive RFC 5322 compliance, but practical validation.
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Email represents a validated email address.
type Email string

// NewEmail creates a validated Email. Returns ErrInvalidEmail if validation fails.
func NewEmail(v string) (Email, error) {
	if v == "" {
		return "", pkgerrors.ErrInvalidEmail
	}
	// mail.ParseAddress handles RFC 5322 parsing (including display names)
	addr, err := mail.ParseAddress(v)
	if err != nil {
		return "", pkgerrors.ErrInvalidEmail
	}
	// Extract just the email address (strip any display name)
	email := addr.Address
	// Additional format check for common patterns
	if !emailRegex.MatchString(email) {
		return "", pkgerrors.ErrInvalidEmail
	}

	return Email(email), nil
}

// String returns the raw email string.
func (e Email) String() string { return string(e) }

// IsZero reports whether the email is empty.
func (e Email) IsZero() bool { return e == "" }

// LocalPart returns the portion before the "@".
func (e Email) LocalPart() string {
	s, _, _ := e.split()

	return s
}

// Domain returns the portion after the "@".
func (e Email) Domain() string {
	_, d, _ := e.split()

	return d
}

// Validate implements validate.Validator for Email.
func (e Email) Validate() error {
	_, err := NewEmail(string(e))

	return err
}

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
func (e Email) split() (string, string, bool) {
	if e == "" {
		return "", "", false
	}

	s := string(e)

	idx := strings.IndexByte(s, '@')
	if idx <= 0 || idx >= len(s)-1 {
		return "", "", false
	}

	return s[:idx], s[idx+1:], true
}

// URL represents a validated URL with http or https scheme.
type URL string

// NewURL creates a validated URL. Returns ErrInvalidURL if validation fails.
// Requires http or https scheme and a valid host.
func NewURL(v string) (URL, error) {
	if v == "" {
		return "", pkgerrors.ErrInvalidURL
	}

	parsed, err := url.Parse(v)
	if err != nil {
		return "", pkgerrors.ErrInvalidURL
	}
	// Require scheme and host
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", pkgerrors.ErrInvalidURL
	}

	if parsed.Host == "" {
		return "", pkgerrors.ErrInvalidURL
	}

	return URL(v), nil
}

// String returns the raw URL string.
func (u URL) String() string { return string(u) }

// IsZero returns true if the URL is empty.
func (u URL) IsZero() bool { return u == "" }

// Validate implements validate.Validator for URL.
func (u URL) Validate() error {
	_, err := NewURL(string(u))

	return err
}

// Parse returns the underlying url.URL. Since URLs are validated at construction,
// this should never fail, but returns error for safety.
func (u URL) Parse() (*url.URL, error) {
	parsed, err := url.Parse(string(u))
	if err != nil {
		return nil, fmt.Errorf("url: parse %q: %w", string(u), err)
	}

	return parsed, nil
}

// Scheme returns the URL scheme (http or https). Returns empty string if URL is zero.
func (u URL) Scheme() string {
	if u == "" {
		return ""
	}

	s := string(u)

	idx := strings.IndexByte(s, ':')
	if idx <= 0 {
		return ""
	}

	return s[:idx]
}

// Host returns the URL host (e.g., "example.com" or "example.com:8080").
func (u URL) Host() string {
	return u.extractField(func(p *url.URL) string { return p.Host })
}

// Path returns the URL path (e.g., "/api/v1/users").
func (u URL) Path() string {
	return u.extractField(func(p *url.URL) string { return p.Path })
}

func (u URL) extractField(fn func(*url.URL) string) string {
	parsed, err := u.Parse()
	if err != nil || parsed == nil {
		return ""
	}

	return fn(parsed)
}
