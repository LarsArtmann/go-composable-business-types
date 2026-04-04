package errors

import (
	"errors"
	"fmt"
	"testing"
)

// =============================================================================
// Structured Error Tests
// =============================================================================

func testStructuredError(err, original error, want string, checkFn func()) {
	if !errors.Is(err, original) {
		panic("expected errors.Is() to find original error")
	}

	if got := err.Error(); got != want {
		panic(fmt.Sprintf("Error() = %q, want %q", got, want))
	}

	checkFn()
}

func TestUnmarshalError(t *testing.T) {
	t.Parallel()

	original := errors.New("parse failed")
	err := &UnmarshalError{
		Type:  "JSON",
		Input: `{invalid}`,
		Err:   original,
	}
	testStructuredError(err, original, "unmarshal JSON: {invalid}: parse failed", func() {
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
	})
}

func TestValidationError(t *testing.T) {
	t.Parallel()

	original := errors.New("invalid format")
	err := &ValidationError{
		Field: "email",
		Value: "not-an-email",
		Err:   original,
	}
	testStructuredError(
		err,
		original,
		"validation failed for email=not-an-email: invalid format",
		func() {
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
		},
	)
}

func TestRangeError(t *testing.T) {
	t.Parallel()
	testRangeError(t, 5, 10, 100, false, "value 5 below minimum 10")
	testRangeError(t, 150, 10, 100, true, "value 150 exceeds maximum 100")
}

func testRangeError(t *testing.T, value, min, max int, outOfRange bool, want string) {
	t.Helper()

	err := &RangeError{
		Value:      value,
		Min:        min,
		Max:        max,
		OutOfRange: outOfRange,
	}
	if got := err.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

// testAsErrorHelper tests that an As*Error function correctly extracts
// the expected error type and returns false for non-matching errors.
func testAsErrorHelper[E any](
	t *testing.T,
	err error,
	asFn func(error) (E, bool),
	checkFn func(E),
	fnName string,
) {
	t.Helper()

	got, ok := asFn(err)
	if !ok {
		t.Fatalf("expected %s to succeed", fnName)
	}

	checkFn(got)

	_, ok = asFn(errors.New("other"))
	if ok {
		t.Fatalf("%s should return false for non-matching error", fnName)
	}
}

func TestAsUnmarshalError(t *testing.T) {
	t.Parallel()

	err := &UnmarshalError{Type: "JSON", Input: `{bad}`, Err: errors.New("fail")}

	testAsErrorHelper(t, err, AsUnmarshalError, func(got *UnmarshalError) {
		if got.Type != "JSON" {
			t.Errorf("Type = %q, want JSON", got.Type)
		}
	}, "AsUnmarshalError")
}

func TestAsValidationError(t *testing.T) {
	t.Parallel()

	err := &ValidationError{Field: "email", Value: "x", Err: errors.New("fail")}

	testAsErrorHelper(t, err, AsValidationError, func(got *ValidationError) {
		if got.Field != "email" {
			t.Errorf("Field = %q, want email", got.Field)
		}
	}, "AsValidationError")
}

func TestAsRangeError(t *testing.T) {
	t.Parallel()

	err := &RangeError{Value: 5, Min: 1, Max: 10, OutOfRange: true}

	testAsErrorHelper(t, err, AsRangeError, func(got *RangeError) {
		if got.Value != 5 {
			t.Errorf("Value = %v, want 5", got.Value)
		}
	}, "AsRangeError")
}

func TestAsScanError(t *testing.T) {
	t.Parallel()

	err := &ScanError{SourceType: "int64", TargetType: "string", Err: errors.New("fail")}

	testAsErrorHelper(t, err, AsScanError, func(got *ScanError) {
		if got.SourceType != "int64" {
			t.Errorf("SourceType = %q, want int64", got.SourceType)
		}
	}, "AsScanError")
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
