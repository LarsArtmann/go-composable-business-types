// Package errors provides centralized domain-specific errors for the library.
package errors

import (
	"errors"
	"testing"
)

// =============================================================================
// Sentinel Error Tests
// =============================================================================

func TestEmailSentinels(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		err  error
		want error
	}{
		{"ErrInvalidEmail", ErrInvalidEmail, ErrInvalidEmail},
		{"ErrEmailEmpty", ErrEmailEmpty, ErrEmailEmpty},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if !errors.Is(tt.err, tt.want) {
				t.Errorf("errors.Is() = false, want true")
			}
		})
	}
}

func TestURLSentinels(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		err  error
		want error
	}{
		{"ErrInvalidURL", ErrInvalidURL, ErrInvalidURL},
		{"ErrURLEmpty", ErrURLEmpty, ErrURLEmpty},
		{"ErrURLScheme", ErrURLScheme, ErrURLScheme},
		{"ErrURLHost", ErrURLHost, ErrURLHost},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if !errors.Is(tt.err, tt.want) {
				t.Errorf("errors.Is() = false, want true")
			}
		})
	}
}

func TestBoundedStringSentinels(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		err  error
		want error
	}{
		{"ErrBoundedStringMinLength", ErrBoundedStringMinLength, ErrBoundedStringMinLength},
		{"ErrBoundedStringMaxLength", ErrBoundedStringMaxLength, ErrBoundedStringMaxLength},
		{"ErrBoundedStringMinNegative", ErrBoundedStringMinNegative, ErrBoundedStringMinNegative},
		{
			"ErrBoundedStringMaxLessThanMin",
			ErrBoundedStringMaxLessThanMin,
			ErrBoundedStringMaxLessThanMin,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if !errors.Is(tt.err, tt.want) {
				t.Errorf("errors.Is() = false, want true")
			}
		})
	}
}

func TestNanoIDSentinels(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		err  error
		want error
	}{
		{"ErrNanoIDEmpty", ErrNanoIDEmpty, ErrNanoIDEmpty},
		{"ErrNanoIDTooShort", ErrNanoIDTooShort, ErrNanoIDTooShort},
		{"ErrNanoIDTooLong", ErrNanoIDTooLong, ErrNanoIDTooLong},
		{"ErrNanoIDInvalid", ErrNanoIDInvalid, ErrNanoIDInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if !errors.Is(tt.err, tt.want) {
				t.Errorf("errors.Is() = false, want true")
			}
		})
	}
}

func TestIDSentinels(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		err  error
		want error
	}{
		{"ErrIDInvalid", ErrIDInvalid, ErrIDInvalid},
		{"ErrIDTypeNotSupported", ErrIDTypeNotSupported, ErrIDTypeNotSupported},
		{"ErrIDInsufficientData", ErrIDInsufficientData, ErrIDInsufficientData},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if !errors.Is(tt.err, tt.want) {
				t.Errorf("errors.Is() = false, want true")
			}
		})
	}
}
