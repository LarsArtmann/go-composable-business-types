package errors

import (
	"errors"
	"testing"
)

// =============================================================================
// Structured Error Tests
// =============================================================================

func TestUnmarshalError(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
	t.Run("below minimum", func(t *testing.T) {
		t.Parallel()
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
		t.Parallel()
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
	t.Parallel()
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
