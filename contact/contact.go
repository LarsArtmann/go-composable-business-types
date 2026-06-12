// Package contact provides validated contact types.
//
// Contact represents a person or organization's contact information, including
// name, email, phone, website, and physical address.
// It implements JSON marshaling and validate.Validator.
//
// Basic usage:
//
//	c := contact.New(
//	    "Jane Doe",
//	    contact.WithEmail("jane@example.com"),
//	    contact.WithPhone("+49 30 123456"),
//	    contact.WithAddress(address.New("Main Str. 1", "Berlin", "10115", "DE")),
//	)
//
//	if err := c.Validate(); err != nil { ... }
package contact

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/larsartmann/go-composable-business-types/address"
	"github.com/larsartmann/go-composable-business-types/validate"
)

const (
	maxNameLen  = 200
	maxEmailLen = 254
	maxPhoneLen = 50
)

var (
	errContactNil        = errors.New("contact: nil")
	errContactNameReq    = errors.New("contact: name is required")
	errContactNameTooLong = errors.New("contact: name exceeds maximum length")
	errContactEmailTooLong = errors.New("contact: email exceeds maximum length")
	errContactEmailInvalid = errors.New("contact: email is invalid")
	errContactPhoneTooLong = errors.New("contact: phone exceeds maximum length")
	errContactWebsiteInvalid = errors.New("contact: website URL is invalid")
)

// Contact represents a person or organization's contact information.
type Contact struct {
	Name    string           `json:"name"`
	Email   string           `json:"email,omitempty"`
	Phone   string           `json:"phone,omitempty"`
	Website string           `json:"website,omitempty"`
	Address *address.Address `json:"address,omitempty"`
}

// Option configures a Contact during construction.
type Option func(*Contact)

// WithEmail sets the contact email address.
func WithEmail(email string) Option {
	return func(c *Contact) { c.Email = strings.TrimSpace(email) }
}

// WithPhone sets the contact phone number.
func WithPhone(phone string) Option {
	return func(c *Contact) { c.Phone = strings.TrimSpace(phone) }
}

// WithWebsite sets the contact website URL.
func WithWebsite(website string) Option {
	return func(c *Contact) { c.Website = strings.TrimSpace(website) }
}

// WithAddress sets the contact physical address.
func WithAddress(addr *address.Address) Option {
	return func(c *Contact) { c.Address = addr }
}

// New creates a Contact with the required name and optional configuration.
func New(name string, opts ...Option) *Contact {
	c := &Contact{Name: strings.TrimSpace(name)}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// IsZero reports whether the Contact is nil or has no name.
func (c *Contact) IsZero() bool {
	return c == nil || c.Name == ""
}

// Validate checks that the contact is valid.
func (c *Contact) Validate() error {
	if c == nil {
		return errContactNil
	}

	if c.Name == "" {
		return errContactNameReq
	}

	if utf8.RuneCountInString(c.Name) > maxNameLen {
		return fmt.Errorf("%w: %d > %d", errContactNameTooLong, utf8.RuneCountInString(c.Name), maxNameLen)
	}

	if c.Email != "" {
		if err := validateEmail(c.Email); err != nil {
			return err
		}
	}

	if c.Phone != "" && utf8.RuneCountInString(c.Phone) > maxPhoneLen {
		return fmt.Errorf("%w: %d > %d", errContactPhoneTooLong, utf8.RuneCountInString(c.Phone), maxPhoneLen)
	}

	if c.Website != "" {
		if err := validateWebsite(c.Website); err != nil {
			return err
		}
	}

	if c.Address != nil {
		if err := c.Address.Validate(); err != nil {
			return fmt.Errorf("address: %w", err)
		}
	}

	return nil
}

func validateEmail(email string) error {
	if utf8.RuneCountInString(email) > maxEmailLen {
		return fmt.Errorf("%w: %d > %d", errContactEmailTooLong, utf8.RuneCountInString(email), maxEmailLen)
	}

	at := strings.LastIndex(email, "@")
	if at <= 0 || at == len(email)-1 {
		return fmt.Errorf("%w: %q", errContactEmailInvalid, email)
	}

	local := email[:at]
	domain := email[at+1:]
	if local == "" || domain == "" || strings.Contains(domain, "..") || strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return fmt.Errorf("%w: %q", errContactEmailInvalid, email)
	}

	return nil
}

func validateWebsite(website string) error {
	u, err := url.Parse(website)
	if err != nil {
		return fmt.Errorf("%w: %q", errContactWebsiteInvalid, website)
	}

	if u.Scheme != "http" && u.Scheme != "https" || u.Host == "" {
		return fmt.Errorf("%w: %q", errContactWebsiteInvalid, website)
	}

	return nil
}

var _ validate.Validator = (*Contact)(nil)
