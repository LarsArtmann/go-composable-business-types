package locale

import (
	"testing"

	"github.com/larsartmann/go-composable-business-types/testutil"
	"golang.org/x/text/language"
)

func TestParseLocale(t *testing.T) {
	t.Parallel()

	testutil.RunParseTests(t, "Locale", []testutil.ParseTestCase[Locale]{
		{Name: "en-US", Input: "en-US", WantErr: false},
		{Name: "de-DE", Input: "de-DE", WantErr: false},
		{Name: "fr-FR", Input: "fr-FR", WantErr: false},
		{Name: "empty", Input: "", WantErr: true},
	}, ParseLocale)
}

func TestParseLocaleError(t *testing.T) {
	testutil.RunParseErrorTest(t, "Locale", ParseLocale)
}

func TestLocaleParts(t *testing.T) {
	t.Parallel()

	loc, _ := ParseLocale("en-US")
	testutil.RunPartsTest(t, loc, []testutil.PartAccessor[Locale]{
		{Name: "base", Get: func(l Locale) string { return l.Base() }, Expected: "en"},
		{Name: "region", Get: func(l Locale) string { return l.Region() }, Expected: "US"},
	})
}

func TestLocaleIsZero(t *testing.T) {
	testutil.RunIsZeroTest(t, func() (Locale, error) {
		return ParseLocale("en-US")
	})
}

func TestLocaleConstants(t *testing.T) {
	t.Parallel()
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
	}

	testutil.RunStringTests(t, "Locale", []testutil.StringCase[Locale]{
		{Value: LocaleEnUS, Expected: "en-US"},
		{Value: LocaleEnGB, Expected: "en-GB"},
		{Value: LocaleDeDE, Expected: "de-DE"},
		{Value: LocaleFrFR, Expected: "fr-FR"},
		{Value: LocaleEsES, Expected: "es-ES"},
		{Value: LocaleItIT, Expected: "it-IT"},
		{Value: LocaleJaJP, Expected: "ja-JP"},
		{Value: LocaleZhCN, Expected: "zh-CN"},
	})
}

func TestNewLocale(t *testing.T) {
	t.Parallel()

	tag, _ := language.Parse("en-GB")

	loc := NewLocale(tag)
	if loc.String() != "en-GB" {
		t.Errorf("expected en-GB, got %s", loc.String())
	}
}

func TestLocaleTag(t *testing.T) {
	t.Parallel()

	loc, _ := ParseLocale("de-DE")

	tag := loc.Tag()
	if tag.String() != "de-DE" {
		t.Errorf("expected de-DE, got %s", tag.String())
	}
}

func TestLocaleMarshal(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
