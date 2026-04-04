// Package locale provides BCP 47 language tag handling for internationalization.
//
// Locale wraps golang.org/x/text/language.Tag for full BCP 47 compliance.
// It supports parsing, serialization, and SQL storage.
//
// Basic usage:
//
//	loc, err := locale.ParseLocale("en-US")
//	// Or use constants:
//	loc = locale.LocaleDeDE
package locale

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/larsartmann/go-composable-business-types/scanutil"
	"golang.org/x/text/language"
)

// Locale represents a BCP 47 language tag for internationalization.
// It wraps golang.org/x/text/language.Tag for full BCP 47 compliance.
type Locale struct {
	tag language.Tag
}

// Common locale constants for convenience (compile-time validated).
var (
	LocaleEnUS = mustNewLocale("en-US") // American English
	LocaleEnGB = mustNewLocale("en-GB") // British English
	LocaleDeDE = mustNewLocale("de-DE") // German (Germany)
	LocaleFrFR = mustNewLocale("fr-FR") // French (France)
	LocaleEsES = mustNewLocale("es-ES") // Spanish (Spain)
	LocaleItIT = mustNewLocale("it-IT") // Italian (Italy)
	LocaleJaJP = mustNewLocale("ja-JP") // Japanese (Japan)
	LocaleZhCN = mustNewLocale("zh-CN") // Simplified Chinese (China)
)

// mustNewLocale validates locale strings at package initialization.
// These are compile-time constants, so panics here indicate developer error.
func mustNewLocale(s string) Locale {
	loc, err := ParseLocale(s)
	if err != nil {
		panic(fmt.Errorf("invalid locale constant %q: %w", s, err))
	}

	return loc
}

// ParseLocale parses a BCP 47 language tag string.
// Accepts both hyphen (en-US) and underscore (en_US) formats.
func ParseLocale(s string) (Locale, error) {
	tag, err := language.Parse(s)
	if err != nil {
		return Locale{}, fmt.Errorf("parse locale %q: %w", s, err)
	}

	return Locale{tag: tag}, nil
}

// NewLocale creates a Locale from a language.Tag.
func NewLocale(tag language.Tag) Locale {
	return Locale{tag: tag}
}

// Tag returns the underlying language.Tag.
func (l Locale) Tag() language.Tag {
	return l.tag
}

// String returns the BCP 47 representation (e.g., "en-US").
func (l Locale) String() string {
	return l.tag.String()
}

// IsZero returns true if the locale is empty/undefined.
func (l Locale) IsZero() bool {
	return l.tag == language.Und
}

// Base returns the base language (e.g., "en" for "en-US").
func (l Locale) Base() string {
	base, _ := l.tag.Base()

	return base.String()
}

// Region returns the region code (e.g., "US" for "en-US").
// Returns empty string if no region is specified.
func (l Locale) Region() string {
	region, _ := l.tag.Region()

	return region.String()
}

// MarshalText implements encoding.TextMarshaler for JSON/XML serialization.
func (l Locale) MarshalText() ([]byte, error) {
	b, err := l.tag.MarshalText()
	if err != nil {
		return nil, fmt.Errorf("locale: marshal text: %w", err)
	}

	return b, nil
}

// UnmarshalText implements encoding.TextUnmarshaler for JSON/XML deserialization.
func (l *Locale) UnmarshalText(data []byte) error {
	err := l.tag.UnmarshalText(data)
	if err != nil {
		return fmt.Errorf("locale: unmarshal text: %w", err)
	}

	return nil
}

// Scan implements sql.Scanner for Locale.
// Supports string and []byte sources. Empty string/nil results in zero value.
func (l *Locale) Scan(src any) error {
	if l == nil {
		return errors.New("locale: scan: receiver is nil")
	}

	err := scanutil.ScanString(src, func(v string) error {
		if v == "" {
			l.tag = language.Und

			return nil
		}

		parsed, err := ParseLocale(v)
		if err != nil {
			return fmt.Errorf("scan locale: %w", err)
		}

		*l = parsed

		return nil
	})
	if err != nil {
		return fmt.Errorf("locale: scan: %w", err)
	}

	return nil
}

// Value implements driver.Valuer for Locale.
// Returns nil for empty/undefined locale, otherwise the BCP 47 string.
func (l Locale) Value() (driver.Value, error) {
	if l.IsZero() {
		return nil, nil
	}

	return l.String(), nil
}
