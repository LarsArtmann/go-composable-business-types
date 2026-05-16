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

func testAs[E any](t *testing.T, err error, checkFn func(E)) {
	t.Helper()

	var target E
	if !errors.As(err, &target) {
		t.Error("expected errors.As() to succeed")
	} else {
		checkFn(target)
	}
}

func TestUnmarshalError(t *testing.T) {
	t.Parallel()

	original := errors.New("parse failed")
	err := &UnmarshalError{
		Type:         "JSON",
		Input:        `{invalid}`,
		wrappedError: wrappedError{Err: original},
	}
	testStructuredError(err, original, "unmarshal JSON: {invalid}: parse failed", func() {
		testAs[*UnmarshalError](t, err, func(target *UnmarshalError) {
			if target.Type != "JSON" {
				t.Errorf("Type = %q, want JSON", target.Type)
			}

			if target.Input != `{invalid}` {
				t.Errorf("Input = %q, want {invalid}", target.Input)
			}
		})
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
			testAs[*ValidationError](t, err, func(target *ValidationError) {
				if target.Field != "email" {
					t.Errorf("Field = %q, want email", target.Field)
				}

				if target.Value != "not-an-email" {
					t.Errorf("Value = %v, want not-an-email", target.Value)
				}
			})
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

// testAsErrorField tests that an As*Error function extracts the correct field value.
func testAsErrorField(
	t *testing.T,
	err error,
	asFn func(error) (any, bool),
	fnName string,
	fieldName string,
	wantValue any,
	getField func(any) any,
) {
	t.Helper()

	got, ok := asFn(err)
	if !ok {
		t.Fatalf("expected %s to succeed", fnName)
	}

	if gotField := getField(got); gotField != wantValue {
		t.Errorf("%s = %v, want %v", fieldName, gotField, wantValue)
	}

	_, ok = asFn(errors.New("other"))
	if ok {
		t.Fatalf("%s should return false for non-matching error", fnName)
	}
}

func TestAsErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		err       error
		asFn      func(error) (any, bool)
		fnName    string
		fieldName string
		wantValue any
		getField  func(any) any
	}{
		{
			name: "UnmarshalError",
			err: &UnmarshalError{
				Type:         "JSON",
				Input:        `{bad}`,
				wrappedError: wrappedError{Err: errors.New("fail")},
			},
			asFn:      func(err error) (any, bool) { return AsUnmarshalError(err) },
			fnName:    "AsUnmarshalError",
			fieldName: "Type",
			wantValue: "JSON",
			getField:  func(e any) any { return e.(*UnmarshalError).Type },
		},
		{
			name:      "ValidationError",
			err:       &ValidationError{Field: "email", Value: "x", Err: errors.New("fail")},
			asFn:      func(err error) (any, bool) { return AsValidationError(err) },
			fnName:    "AsValidationError",
			fieldName: "Field",
			wantValue: "email",
			getField:  func(e any) any { return e.(*ValidationError).Field },
		},
		{
			name:      "RangeError",
			err:       &RangeError{Value: 5, Min: 1, Max: 10, OutOfRange: true},
			asFn:      func(err error) (any, bool) { return AsRangeError(err) },
			fnName:    "AsRangeError",
			fieldName: "Value",
			wantValue: 5,
			getField:  func(e any) any { return e.(*RangeError).Value },
		},
		{
			name: "ScanError",
			err: &ScanError{
				SourceType:   "int64",
				TargetType:   "string",
				wrappedError: wrappedError{Err: errors.New("fail")},
			},
			asFn:      func(err error) (any, bool) { return AsScanError(err) },
			fnName:    "AsScanError",
			fieldName: "SourceType",
			wantValue: "int64",
			getField:  func(e any) any { return e.(*ScanError).SourceType },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testAsErrorField(t, tt.err, tt.asFn, tt.fnName, tt.fieldName, tt.wantValue, tt.getField)
		})
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
