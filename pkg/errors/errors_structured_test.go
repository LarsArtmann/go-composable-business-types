package errors

import (
	"errors"
	"fmt"
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

func TestAsUnmarshalError(t *testing.T) {
	t.Parallel()
	err := &UnmarshalError{Type: "JSON", Input: `{bad}`, Err: errors.New("fail")}
	got, ok := AsUnmarshalError(err)
	if !ok {
		t.Fatal("expected AsUnmarshalError to succeed")
	}
	if got.Type != "JSON" {
		t.Errorf("Type = %q, want JSON", got.Type)
	}
	_, ok = AsUnmarshalError(errors.New("other"))
	if ok {
		t.Error("expected AsUnmarshalError to return false for non-matching error")
	}
}

func TestAsValidationError(t *testing.T) {
	t.Parallel()
	err := &ValidationError{Field: "email", Value: "x", Err: errors.New("fail")}
	got, ok := AsValidationError(err)
	if !ok {
		t.Fatal("expected AsValidationError to succeed")
	}
	if got.Field != "email" {
		t.Errorf("Field = %q, want email", got.Field)
	}
	_, ok = AsValidationError(errors.New("other"))
	if ok {
		t.Error("expected AsValidationError to return false for non-matching error")
	}
}

func TestAsRangeError(t *testing.T) {
	t.Parallel()
	err := &RangeError{Value: 5, Min: 1, Max: 10, OutOfRange: true}
	got, ok := AsRangeError(err)
	if !ok {
		t.Fatal("expected AsRangeError to succeed")
	}
	if got.Value != 5 {
		t.Errorf("Value = %v, want 5", got.Value)
	}
	_, ok = AsRangeError(errors.New("other"))
	if ok {
		t.Error("expected AsRangeError to return false for non-matching error")
	}
}

func TestAsScanError(t *testing.T) {
	t.Parallel()
	err := &ScanError{SourceType: "int64", TargetType: "string", Err: errors.New("fail")}
	got, ok := AsScanError(err)
	if !ok {
		t.Fatal("expected AsScanError to succeed")
	}
	if got.SourceType != "int64" {
		t.Errorf("SourceType = %q, want int64", got.SourceType)
	}
	_, ok = AsScanError(errors.New("other"))
	if ok {
		t.Error("expected AsScanError to return false for non-matching error")
	}
}

func TestAsTypeWrappedErrors(t *testing.T) {
	t.Parallel()
	inner := &ValidationError{Field: "name", Value: "", Err: errors.New("empty")}
	wrapped := fmt.Errorf("processing: %w", inner)
	got, ok := AsValidationError(wrapped)
	if !ok {
		t.Fatal("expected AsValidationError to find wrapped error")
	}
	if got.Field != "name" {
		t.Errorf("Field = %q, want name", got.Field)
	}
}
