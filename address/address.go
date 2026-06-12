// Package address provides validated postal address types.
//
// Address represents a physical mailing address with support for street lines,
// city, state/province, postal code, and ISO 3166-1 alpha-2 country code.
// It implements JSON marshaling and validate.Validator.
//
// Basic usage:
//
//	addr := address.New(
//	    "123 Main St",
//	    "Berlin",
//	    "10115",
//	    "DE",
//	    address.WithState("BE"),
//	)
//
//	if err := addr.Validate(); err != nil { ... }
package address

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/larsartmann/go-composable-business-types/validate"
)

const (
	maxLineLen       = 100
	maxCityLen       = 100
	maxStateLen      = 100
	maxPostalCodeLen = 20
	countryCodeLen   = 2
)

var (
	errAddressNil           = errors.New("address: nil")
	errAddressLine1Req      = errors.New("address: line1 is required")
	errAddressLine1TooLong  = errors.New("address: line1 exceeds maximum length")
	errAddressLine2TooLong  = errors.New("address: line2 exceeds maximum length")
	errAddressCityReq       = errors.New("address: city is required")
	errAddressCityTooLong   = errors.New("address: city exceeds maximum length")
	errAddressStateTooLong  = errors.New("address: state exceeds maximum length")
	errAddressPostalReq     = errors.New("address: postal code is required")
	errAddressPostalTooLong = errors.New("address: postal code exceeds maximum length")
	errAddressCountryReq    = errors.New("address: country code is required")
	errAddressCountryLen    = errors.New("address: country code must be 2 characters")
	errAddressCountryAlpha  = errors.New("address: country code must be alphabetic")
)

// Address represents a validated postal address.
type Address struct {
	Line1       string `json:"line1"`
	Line2       string `json:"line2,omitempty"`
	City        string `json:"city"`
	State       string `json:"state,omitempty"`
	PostalCode  string `json:"postalCode"`
	CountryCode string `json:"countryCode"`
}

// Option configures an Address during construction.
type Option func(*Address)

// WithLine2 sets the second address line.
func WithLine2(line2 string) Option {
	return func(a *Address) { a.Line2 = line2 }
}

// WithState sets the state or province.
func WithState(state string) Option {
	return func(a *Address) { a.State = state }
}

// New creates an Address with the required fields and optional configuration.
func New(line1, city, postalCode, countryCode string, opts ...Option) *Address {
	a := &Address{
		Line1:       line1,
		City:        city,
		PostalCode:  postalCode,
		CountryCode: strings.ToUpper(countryCode),
	}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

// IsZero reports whether the Address is nil or has no required fields.
func (a *Address) IsZero() bool {
	return a == nil || (a.Line1 == "" && a.City == "" && a.PostalCode == "" && a.CountryCode == "")
}

// Validate checks that all required fields are present and within limits.
func (a *Address) Validate() error {
	if a == nil {
		return errAddressNil
	}

	if a.Line1 == "" {
		return errAddressLine1Req
	}

	if utf8.RuneCountInString(a.Line1) > maxLineLen {
		return fmt.Errorf("%w: %d > %d", errAddressLine1TooLong, utf8.RuneCountInString(a.Line1), maxLineLen)
	}

	if a.Line2 != "" && utf8.RuneCountInString(a.Line2) > maxLineLen {
		return fmt.Errorf("%w: %d > %d", errAddressLine2TooLong, utf8.RuneCountInString(a.Line2), maxLineLen)
	}

	if a.City == "" {
		return errAddressCityReq
	}

	if utf8.RuneCountInString(a.City) > maxCityLen {
		return fmt.Errorf("%w: %d > %d", errAddressCityTooLong, utf8.RuneCountInString(a.City), maxCityLen)
	}

	if a.State != "" && utf8.RuneCountInString(a.State) > maxStateLen {
		return fmt.Errorf("%w: %d > %d", errAddressStateTooLong, utf8.RuneCountInString(a.State), maxStateLen)
	}

	if a.PostalCode == "" {
		return errAddressPostalReq
	}

	if utf8.RuneCountInString(a.PostalCode) > maxPostalCodeLen {
		return fmt.Errorf(
			"%w: %d > %d",
			errAddressPostalTooLong,
			utf8.RuneCountInString(a.PostalCode),
			maxPostalCodeLen,
		)
	}

	if a.CountryCode == "" {
		return errAddressCountryReq
	}

	if utf8.RuneCountInString(a.CountryCode) != countryCodeLen {
		return fmt.Errorf("%w: %q", errAddressCountryLen, a.CountryCode)
	}

	for _, r := range a.CountryCode {
		if (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') {
			return fmt.Errorf("%w: %q", errAddressCountryAlpha, a.CountryCode)
		}
	}

	return nil
}

// Format returns a human-readable multi-line representation of the address.
func (a *Address) Format() string {
	if a == nil {
		return ""
	}

	var parts []string

	parts = append(parts, a.Line1)

	if a.Line2 != "" {
		parts = append(parts, a.Line2)
	}

	cityLine := a.City
	if a.State != "" {
		cityLine += ", " + a.State
	}

	if a.PostalCode != "" {
		cityLine += " " + a.PostalCode
	}

	parts = append(parts, strings.TrimSpace(cityLine))
	parts = append(parts, strings.ToUpper(a.CountryCode))

	return strings.Join(parts, "\n")
}

var _ validate.Validator = (*Address)(nil)
