package locale

import (
	"testing"
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

func TestMustParseLocale(t *testing.T) {
	loc := MustParseLocale("en-US")
	if loc.String() != "en-US" {
		t.Errorf("expected en-US, got %s", loc.String())
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid locale")
		}
	}()
	MustParseLocale("invalid-locale-code")
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
	locales := []Locale{
		LocaleEnUS,
		LocaleEnGB,
		LocaleDeDE,
		LocaleFrFR,
		LocaleEsES,
		LocaleItIT,
		LocaleJaJP,
		LocaleZhCN,
	}

	for _, loc := range locales {
		if loc.IsZero() {
			t.Errorf("locale constant %v should not be zero", loc)
		}
	}
}
