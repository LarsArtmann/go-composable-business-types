package money

import (
	"testing"

	"github.com/larsartmann/go-composable-business-types/locale"
)

func TestNewMoney(t *testing.T) {
	m, err := NewMoney("10.99", "USD")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify it's a valid amount
	if m.IsZero() {
		t.Error("amount should not be zero")
	}
}

func TestNewMoneyFromCents(t *testing.T) {
	m, err := NewMoneyFromCents(1099, "USD")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Just verify it's a valid amount
	if m.IsZero() {
		t.Error("amount should not be zero")
	}
}

func TestIsValidCurrency(t *testing.T) {
	if !IsValidCurrency("USD") {
		t.Error("USD should be valid")
	}
	if !IsValidCurrency("EUR") {
		t.Error("EUR should be valid")
	}
	if IsValidCurrency("INVALID") {
		t.Error("INVALID should not be valid")
	}
}

func TestCurrencyDigits(t *testing.T) {
	// USD has 2 decimal digits
	digits, ok := CurrencyDigits("USD")
	if !ok {
		t.Error("USD should have digits defined")
	}
	if digits != 2 {
		t.Errorf("USD should have 2 digits, got %d", digits)
	}

	// JPY has 0 decimal digits
	digits, ok = CurrencyDigits("JPY")
	if !ok {
		t.Error("JPY should have digits defined")
	}
	if digits != 0 {
		t.Errorf("JPY should have 0 digits, got %d", digits)
	}
}

func TestCurrencySymbol(t *testing.T) {
	symbol, ok := CurrencySymbol("USD", "en-US")
	if !ok {
		t.Error("should get symbol")
	}
	if symbol == "" {
		t.Error("symbol should not be empty")
	}
}

func TestCurrencySymbolForLocale(t *testing.T) {
	loc := locale.LocaleEnUS
	symbol, ok := CurrencySymbolForLocale("USD", loc)
	if !ok {
		t.Error("should get symbol")
	}
	if symbol == "" {
		t.Error("symbol should not be empty")
	}
}

func TestAllCurrencyCodes(t *testing.T) {
	codes := AllCurrencyCodes()
	if len(codes) == 0 {
		t.Error("should have currency codes")
	}

	// Check for common currencies
	hasUSD := false
	for _, code := range codes {
		if code == "USD" {
			hasUSD = true
			break
		}
	}
	if !hasUSD {
		t.Error("should have USD in currency codes")
	}
}

func TestFormatMoney(t *testing.T) {
	m, _ := NewMoney("10.99", "USD")
	formatted := FormatMoney(m, "en-US")
	if formatted == "" {
		t.Error("formatted should not be empty")
	}
}

func TestFormatMoneyForLocale(t *testing.T) {
	m, _ := NewMoney("10.99", "USD")
	formatted := FormatMoneyForLocale(m, locale.LocaleEnUS)
	if formatted == "" {
		t.Error("formatted should not be empty")
	}
}

func TestNewMoneyFormatter(t *testing.T) {
	formatter := NewMoneyFormatter("en-US")
	if formatter == nil {
		t.Error("formatter should not be nil")
	}
}

func TestNewMoneyFormatterForLocale(t *testing.T) {
	formatter := NewMoneyFormatterForLocale(locale.LocaleEnUS)
	if formatter == nil {
		t.Error("formatter should not be nil")
	}
}
