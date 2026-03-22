// Package errors provides centralized domain-specific errors for the library.
package errors

import (
	"errors"
	"fmt"
	"testing"
)

// =============================================================================
// Sentinel Error Tests
// =============================================================================

func TestEmailSentinels(t *testing.T) {
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
			if !errors.Is(tt.err, tt.want) {
				t.Errorf("errors.Is() = false, want true")
			}
		})
	}
}

func TestURLSentinels(t *testing.T) {
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
			if !errors.Is(tt.err, tt.want) {
				t.Errorf("errors.Is() = false, want true")
			}
		})
	}
}

func TestBoundedStringSentinels(t *testing.T) {
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
			if !errors.Is(tt.err, tt.want) {
				t.Errorf("errors.Is() = false, want true")
			}
		})
	}
}

func TestNanoIDSentinels(t *testing.T) {
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
			if !errors.Is(tt.err, tt.want) {
				t.Errorf("errors.Is() = false, want true")
			}
		})
	}
}

func TestIDSentinels(t *testing.T) {
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
			if !errors.Is(tt.err, tt.want) {
				t.Errorf("errors.Is() = false, want true")
			}
		})
	}
}

// =============================================================================
// Structured Error Tests
// =============================================================================

func TestUnmarshalError(t *testing.T) {
	original := errors.New("parse failed")
	err := &UnmarshalError{
		Type:  "JSON",
		Input: `{invalid}`,
		Err:   original,
	}

	if !errors.Is(err, original) {
		t.Error("expected errors.Is() to find original error")
	}

	want := "unmarshal JSON: {invalid}: parse failed"
	if got := err.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}

	var target *UnmarshalError
	if !errors.As(err, &target) {
		t.Error("expected errors.As() to succeed")
	} else {
		if target.Type != "JSON" {
			t.Errorf("Type = %q, want JSON", target.Type)
		}
		if target.Input != `{invalid}` {
			t.Errorf("Input = %q, want {invalid}", target.Input)
		}
	}
}

func TestValidationError(t *testing.T) {
	original := errors.New("invalid format")
	err := &ValidationError{
		Field: "email",
		Value: "not-an-email",
		Err:   original,
	}

	if !errors.Is(err, original) {
		t.Error("expected errors.Is() to find original error")
	}

	want := "validation failed for email=not-an-email: invalid format"
	if got := err.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}

	var target *ValidationError
	if !errors.As(err, &target) {
		t.Error("expected errors.As() to succeed")
	} else {
		if target.Field != "email" {
			t.Errorf("Field = %q, want email", target.Field)
		}
		if target.Value != "not-an-email" {
			t.Errorf("Value = %v, want not-an-email", target.Value)
		}
	}
}

func TestRangeError(t *testing.T) {
	t.Run("below minimum", func(t *testing.T) {
		err := &RangeError{
			Value:      5,
			Min:        10,
			Max:        100,
			OutOfRange: false,
		}

		want := "value 5 below minimum 10"
		if got := err.Error(); got != want {
			t.Errorf("Error() = %q, want %q", got, want)
		}
	})

	t.Run("above maximum", func(t *testing.T) {
		err := &RangeError{
			Value:      150,
			Min:        10,
			Max:        100,
			OutOfRange: true,
		}

		want := "value 150 exceeds maximum 100"
		if got := err.Error(); got != want {
			t.Errorf("Error() = %q, want %q", got, want)
		}
	})
}

func TestScanError(t *testing.T) {
	original := errors.New("unsupported type")
	err := &ScanError{
		SourceType: "int64",
		TargetType: "string",
		Err:        original,
	}

	if !errors.Is(err, original) {
		t.Error("expected errors.Is() to find original error")
	}

	want := "cannot scan int64 into string: unsupported type"
	if got := err.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}

	var target *ScanError
	if !errors.As(err, &target) {
		t.Error("expected errors.As() to succeed")
	} else {
		if target.SourceType != "int64" {
			t.Errorf("SourceType = %q, want int64", target.SourceType)
		}
		if target.TargetType != "string" {
			t.Errorf("TargetType = %q, want string", target.TargetType)
		}
	}
}

// =============================================================================
// Wrap Helper Tests
// =============================================================================

func TestWrapMalformed(t *testing.T) {
	t.Run("wraps error", func(t *testing.T) {
		inner := errors.New("invalid syntax")
		err := WrapMalformed(inner, "JSON", `{bad}`)

		if !errors.Is(err, ErrMalformedInput) {
			t.Error("expected errors.Is(err, ErrMalformedInput) = true")
		}
	})

	t.Run("nil returns nil", func(t *testing.T) {
		err := WrapMalformed(nil, "JSON", `{}`)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})
}

func TestWrapInvalid(t *testing.T) {
	t.Run("wraps error", func(t *testing.T) {
		inner := errors.New("invalid value")
		err := WrapInvalid(inner, "email", "test@")

		var target *ValidationError
		if !errors.As(err, &target) {
			t.Error("expected errors.As() to succeed")
		}
		if target.Field != "email" {
			t.Errorf("Field = %q, want email", target.Field)
		}
	})

	t.Run("nil returns nil", func(t *testing.T) {
		err := WrapInvalid(nil, "field", "value")
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})
}

func TestWrapRange(t *testing.T) {
	err := WrapRange(150, 10, 100, true)

	var target *RangeError
	if !errors.As(err, &target) {
		t.Error("expected errors.As() to succeed")
	}
	if target.Value != 150 {
		t.Errorf("Value = %v, want 150", target.Value)
	}
	if target.Min != 10 {
		t.Errorf("Min = %v, want 10", target.Min)
	}
	if target.Max != 100 {
		t.Errorf("Max = %v, want 100", target.Max)
	}
}

func TestWrapScan(t *testing.T) {
	t.Run("wraps error", func(t *testing.T) {
		inner := errors.New("conversion failed")
		err := WrapScan(inner, "[]byte", "string")

		var target *ScanError
		if !errors.As(err, &target) {
			t.Error("expected errors.As() to succeed")
		}
		if target.SourceType != "[]byte" {
			t.Errorf("SourceType = %q, want []byte", target.SourceType)
		}
		if target.TargetType != "string" {
			t.Errorf("TargetType = %q, want string", target.TargetType)
		}
	})

	t.Run("nil returns nil", func(t *testing.T) {
		err := WrapScan(nil, "int", "string")
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})
}

func TestWrapUnmarshal(t *testing.T) {
	t.Run("wraps error", func(t *testing.T) {
		inner := errors.New("unexpected EOF")
		err := WrapUnmarshal(inner, "JSON", `{incomplete`)

		var target *UnmarshalError
		if !errors.As(err, &target) {
			t.Error("expected errors.As() to succeed")
		}
		if target.Type != "JSON" {
			t.Errorf("Type = %q, want JSON", target.Type)
		}
		if target.Input != `{incomplete` {
			t.Errorf("Input = %q, want {incomplete", target.Input)
		}
	})

	t.Run("nil returns nil", func(t *testing.T) {
		err := WrapUnmarshal(nil, "XML", `<tag>`)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})
}

// =============================================================================
// Error Checker Tests
// =============================================================================

func TestIsInvalidEmail(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"ErrInvalidEmail", ErrInvalidEmail, true},
		{"ErrEmailEmpty", ErrEmailEmpty, true},
		{"other error", errors.New("other"), false},
		{"wrapped email error", fmt.Errorf("wrapped: %w", ErrInvalidEmail), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInvalidEmail(tt.err); got != tt.want {
				t.Errorf("IsInvalidEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsInvalidURL(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"ErrInvalidURL", ErrInvalidURL, true},
		{"ErrURLEmpty", ErrURLEmpty, true},
		{"ErrURLScheme", ErrURLScheme, true},
		{"ErrURLHost", ErrURLHost, true},
		{"other error", errors.New("other"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInvalidURL(tt.err); got != tt.want {
				t.Errorf("IsInvalidURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsBoundedStringError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"ErrBoundedStringMinLength", ErrBoundedStringMinLength, true},
		{"ErrBoundedStringMaxLength", ErrBoundedStringMaxLength, true},
		{"ErrBoundedStringMinNegative", ErrBoundedStringMinNegative, true},
		{"ErrBoundedStringMaxLessThanMin", ErrBoundedStringMaxLessThanMin, true},
		{"other error", errors.New("other"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsBoundedStringError(tt.err); got != tt.want {
				t.Errorf("IsBoundedStringError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNanoIDError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"ErrNanoIDEmpty", ErrNanoIDEmpty, true},
		{"ErrNanoIDTooShort", ErrNanoIDTooShort, true},
		{"ErrNanoIDTooLong", ErrNanoIDTooLong, true},
		{"ErrNanoIDInvalid", ErrNanoIDInvalid, true},
		{"other error", errors.New("other"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNanoIDError(tt.err); got != tt.want {
				t.Errorf("IsNanoIDError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsIDError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"ErrIDInvalid", ErrIDInvalid, true},
		{"ErrIDTypeNotSupported", ErrIDTypeNotSupported, true},
		{"ErrIDInsufficientData", ErrIDInsufficientData, true},
		{"other error", errors.New("other"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsIDError(tt.err); got != tt.want {
				t.Errorf("IsIDError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsParseError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"ErrMalformedInput", ErrMalformedInput, true},
		{"ErrUnsupportedType", ErrUnsupportedType, true},
		{"ErrInvalidJSON", ErrInvalidJSON, true},
		{"other error", errors.New("other"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsParseError(tt.err); got != tt.want {
				t.Errorf("IsParseError() = %v, want %v", got, tt.want)
			}
		})
	}
}
