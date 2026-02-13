package cbt

import (
	"github.com/bojanz/currency"
)

// Money wraps currency.Amount for domain clarity and type safety.
type Money = currency.Amount

// NewMoney creates a monetary amount from a numeric string and currency code.
func NewMoney(amount, currencyCode string) (Money, error) {
	return currency.NewAmount(amount, currencyCode)
}

// NewMoneyFromCents creates a monetary amount from minor units (cents).
func NewMoneyFromCents(cents int64, currencyCode string) (Money, error) {
	return currency.NewAmountFromInt64(cents, currencyCode)
}

// IsValidCurrency checks if the currency code is valid ISO 4217.
func IsValidCurrency(code string) bool {
	return currency.IsValid(code)
}

// CurrencyDigits returns the default fraction digits for a currency (e.g., USD=2, JPY=0).
func CurrencyDigits(code string) (uint8, bool) {
	return currency.GetDigits(code)
}

// CurrencySymbol returns the symbol for a currency in the given locale.
func CurrencySymbol(code, locale string) (string, bool) {
	return currency.GetSymbol(code, currency.NewLocale(locale))
}

// CurrencySymbolForLocale returns the symbol for a currency using the Locale type.
func CurrencySymbolForLocale(code string, locale Locale) (string, bool) {
	return currency.GetSymbol(code, currency.NewLocale(locale.String()))
}

// AllCurrencyCodes returns all valid ISO 4217 currency codes.
func AllCurrencyCodes() []string {
	return currency.GetCurrencyCodes()
}

// FormatMoney formats a monetary amount for the given locale.
func FormatMoney(money Money, locale string) string {
	formatter := currency.NewFormatter(currency.NewLocale(locale))
	return formatter.Format(money)
}

// FormatMoneyForLocale formats a monetary amount using the Locale type.
func FormatMoneyForLocale(money Money, locale Locale) string {
	formatter := currency.NewFormatter(currency.NewLocale(locale.String()))
	return formatter.Format(money)
}

// MoneyFormatter wraps currency.Formatter for reuse.
type MoneyFormatter = currency.Formatter

// NewMoneyFormatter creates a locale-aware money formatter.
func NewMoneyFormatter(locale string) *MoneyFormatter {
	return currency.NewFormatter(currency.NewLocale(locale))
}

// NewMoneyFormatterForLocale creates a money formatter using the Locale type.
func NewMoneyFormatterForLocale(locale Locale) *MoneyFormatter {
	return currency.NewFormatter(currency.NewLocale(locale.String()))
}
