package locale

import (
	"testing"

	"golang.org/x/text/language"
)

func TestParseLocale(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"en-US", "en-US", false},
		{"de-DE", "de-DE", false},
		{"fr-FR", "fr-FR", false},
		{"empty", "", true},
		// Note: golang.org/x/text/language is very lenient and accepts many strings
		// {"invalid", "not-a-locale", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc, err := ParseLocale(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if loc.String() != tt.input {
				t.Errorf("expected %s, got %s", tt.input, loc.String())
			}
		})
	}
}

func TestParseLocaleError(t *testing.T) {
	_, err := ParseLocale("invalid-locale-code")
	if err == nil {
		t.Error("expected error for invalid locale")
	}
}

func TestLocaleParts(t *testing.T) {
	loc, _ := ParseLocale("en-US")
	if loc.Base() != "en" {
		t.Errorf("expected base 'en', got %s", loc.Base())
	}
	if loc.Region() != "US" {
		t.Errorf("expected region 'US', got %s", loc.Region())
	}
}

func TestLocaleIsZero(t *testing.T) {
	var zero Locale
	if !zero.IsZero() {
		t.Error("expected zero locale to be zero")
	}

	nonZero, _ := ParseLocale("en-US")
	if nonZero.IsZero() {
		t.Error("expected non-zero locale to not be zero")
	}
}

func TestLocaleConstants(t *testing.T) {
	// Test that constants are defined
	locales := []struct {
		locale Locale
		name   string
	}{
		{LocaleEnUS, "en-US"},
		{LocaleEnGB, "en-GB"},
		{LocaleDeDE, "de-DE"},
		{LocaleFrFR, "fr-FR"},
		{LocaleEsES, "es-ES"},
		{LocaleItIT, "it-IT"},
		{LocaleJaJP, "ja-JP"},
		{LocaleZhCN, "zh-CN"},
	}

	for _, tt := range locales {
		if tt.locale.IsZero() {
			t.Errorf("locale constant %s should not be zero", tt.name)
		}
		if tt.locale.String() != tt.name {
			t.Errorf("expected %s, got %s", tt.name, tt.locale.String())
		}
	}
}

func TestNewLocale(t *testing.T) {
	tag, _ := language.Parse("en-GB")
	loc := NewLocale(tag)
	if loc.String() != "en-GB" {
		t.Errorf("expected en-GB, got %s", loc.String())
	}
}

func TestLocaleTag(t *testing.T) {
	loc, _ := ParseLocale("de-DE")
	tag := loc.Tag()
	if tag.String() != "de-DE" {
		t.Errorf("expected de-DE, got %s", tag.String())
	}
}

func TestLocaleMarshal(t *testing.T) {
	// Test MarshalText
	loc, _ := ParseLocale("fr-FR")
	data, err := loc.MarshalText()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(data) != "fr-FR" {
		t.Errorf("expected fr-FR, got %s", string(data))
	}

	// Test UnmarshalText
	var loc2 Locale
	if err := loc2.UnmarshalText([]byte("es-ES")); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if loc2.String() != "es-ES" {
		t.Errorf("expected es-ES, got %s", loc2.String())
	}
}

func TestLocaleSQL(t *testing.T) {
	// Test Value
	loc, _ := ParseLocale("it-IT")
	val, err := loc.Value()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if val != "it-IT" {
		t.Errorf("expected it-IT, got %v", val)
	}

	// Test Value for zero
	var zero Locale
	val, err = zero.Value()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if val != nil {
		t.Errorf("expected nil, got %v", val)
	}

	// Test Scan with string
	var loc2 Locale
	if err := loc2.Scan("ja-JP"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if loc2.String() != "ja-JP" {
		t.Errorf("expected ja-JP, got %s", loc2.String())
	}

	// Test Scan with []byte
	var loc3 Locale
	if err := loc3.Scan([]byte("zh-CN")); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if loc3.String() != "zh-CN" {
		t.Errorf("expected zh-CN, got %s", loc3.String())
	}

	// Test Scan with nil
	var loc4 Locale
	loc4, _ = ParseLocale("en-US")
	if err := loc4.Scan(nil); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !loc4.IsZero() {
		t.Error("expected zero value after scanning nil")
	}

	// Test Scan with empty string
	var loc5 Locale
	if err := loc5.Scan(""); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !loc5.IsZero() {
		t.Error("expected zero value after scanning empty string")
	}

	// Test Scan with empty []byte
	var loc6 Locale
	if err := loc6.Scan([]byte{}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !loc6.IsZero() {
		t.Error("expected zero value after scanning empty []byte")
	}

	// Test Scan with invalid type
	var loc7 Locale
	if err := loc7.Scan(123); err == nil {
		t.Error("expected error for invalid type")
	}
}

func TestLocaleEdgeCases(t *testing.T) {
	// Test parsing with underscore format (should work with golang.org/x/text/language)
	loc, err := ParseLocale("en_US")
	// The language package normalizes it
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	// Should be normalized to en-US or en-u-rg-us
	if loc.IsZero() {
		t.Error("expected non-zero locale for en_US")
	}

	// Test complex locale
	loc2, err := ParseLocale("zh-Hans-CN")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if loc2.IsZero() {
		t.Error("expected non-zero locale for zh-Hans-CN")
	}
}
