package cbt

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"golang.org/x/text/language"
)

// Locale represents a BCP 47 language tag for internationalization.
// It wraps golang.org/x/text/language.Tag for full BCP 47 compliance.
type Locale struct {
	tag language.Tag
}

// Common locale constants for convenience.
var (
	LocaleEnUS = Locale{tag: language.MustParse("en-US")} // American English
	LocaleEnGB = Locale{tag: language.MustParse("en-GB")} // British English
	LocaleDeDE = Locale{tag: language.MustParse("de-DE")} // German (Germany)
	LocaleFrFR = Locale{tag: language.MustParse("fr-FR")} // French (France)
	LocaleEsES = Locale{tag: language.MustParse("es-ES")} // Spanish (Spain)
	LocaleItIT = Locale{tag: language.MustParse("it-IT")} // Italian (Italy)
	LocaleJaJP = Locale{tag: language.MustParse("ja-JP")} // Japanese (Japan)
	LocaleZhCN = Locale{tag: language.MustParse("zh-CN")} // Simplified Chinese (China)
)

// ParseLocale parses a BCP 47 language tag string.
// Accepts both hyphen (en-US) and underscore (en_US) formats.
func ParseLocale(s string) (Locale, error) {
	tag, err := language.Parse(s)
	if err != nil {
		return Locale{}, fmt.Errorf("parse locale: %w", err)
	}
	return Locale{tag: tag}, nil
}

// MustParseLocale parses a BCP 47 language tag string, panicking on error.
func MustParseLocale(s string) Locale {
	return Locale{tag: language.MustParse(s)}
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
	return l.tag.MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler for JSON/XML deserialization.
func (l *Locale) UnmarshalText(data []byte) error {
	return l.tag.UnmarshalText(data)
}

// Scan implements sql.Scanner for Locale.
// Supports string and []byte sources. Empty string/nil results in zero value.
func (l *Locale) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		l.tag = language.Und
		return nil
	case string:
		if v == "" {
			l.tag = language.Und
			return nil
		}
		parsed, err := ParseLocale(v)
		if err != nil {
			return err
		}
		*l = parsed
		return nil
	case []byte:
		if len(v) == 0 {
			l.tag = language.Und
			return nil
		}
		parsed, err := ParseLocale(string(v))
		if err != nil {
			return err
		}
		*l = parsed
		return nil
	default:
		return errors.New("locale: cannot scan non-string value")
	}
}

// Value implements driver.Valuer for Locale.
// Returns nil for empty/undefined locale, otherwise the BCP 47 string.
func (l Locale) Value() (driver.Value, error) {
	if l.IsZero() {
		return nil, nil
	}
	return l.String(), nil
}
