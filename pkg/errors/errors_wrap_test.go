package errors

import (
	"errors"
	"testing"
)

// nilReturnsNilHelper tests that a wrap function returns nil when given a nil error.
func nilReturnsNilHelper(t *testing.T, name string, wrapFunc func() error) {
	t.Run(name, func(t *testing.T) {
		t.Parallel()

		err := wrapFunc()
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})
}

const (
	jsonType        = "JSON"
	incompleteInput = `{incomplete`
)

func assertFieldValue(t *testing.T, fieldName, actual, want string) {
	if actual != want {
		t.Errorf("%s = %q, want %s", fieldName, actual, want)
	}
}

// =============================================================================
// Wrap Helper Tests
// =============================================================================

func TestWrapMalformed(t *testing.T) {
	t.Parallel()
	t.Run("wraps error", func(t *testing.T) {
		t.Parallel()

		inner := errors.New("invalid syntax")
		err := WrapMalformed(inner, "JSON", `{bad}`)

		if !errors.Is(err, ErrMalformedInput) {
			t.Error("expected errors.Is(err, ErrMalformedInput) = true")
		}
	})

	t.Run("nil returns nil", func(t *testing.T) {
		nilReturnsNilHelper(t, "nil returns nil", func() error {
			return WrapMalformed(nil, "JSON", `{}`)
		})
	})
}

func TestWrapInvalid(t *testing.T) {
	t.Parallel()
	t.Run("wraps error", func(t *testing.T) {
		t.Parallel()

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
		nilReturnsNilHelper(t, "nil returns nil", func() error {
			return WrapInvalid(nil, "field", "value")
		})
	})
}

func TestWrapRange(t *testing.T) {
	t.Parallel()

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

// WrapScanUnmarshalTests holds test data for WrapScan and WrapUnmarshal.
type WrapScanUnmarshalTests struct {
	Name       string
	WrapFunc   func(err error, typeName, input string) error
	ErrorType  string
	TypeField  string
	InputField string
}

var wrapScanUnmarshalTestCases = []WrapScanUnmarshalTests{
	{
		Name:       "WrapScan",
		WrapFunc:   WrapScan,
		ErrorType:  "*ScanError",
		TypeField:  "SourceType",
		InputField: "TargetType",
	},
	{
		Name:       "WrapUnmarshal",
		WrapFunc:   WrapUnmarshal,
		ErrorType:  "*UnmarshalError",
		TypeField:  "Type",
		InputField: "Input",
	},
}

// assertErrorAs checks that err can be asserted to type T using errors.As().
// It logs a test error if the assertion fails and returns the typed error.
func assertErrorAs[T error](t *testing.T, err error) T {
	var target T
	if !errors.As(err, &target) {
		t.Error("expected errors.As() to succeed")
	}

	return target
}

func TestWrapScanUnmarshal(t *testing.T) {
	t.Parallel()

	for _, tc := range wrapScanUnmarshalTestCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Run("wraps error", func(t *testing.T) {
				t.Parallel()

				inner := errors.New("unexpected EOF")
				err := tc.WrapFunc(inner, "JSON", `{incomplete`)

				var target interface{ Error() string }

				switch tc.ErrorType {
				case "*ScanError":
					target = assertErrorAs[*ScanError](t, err)
				case "*UnmarshalError":
					target = assertErrorAs[*UnmarshalError](t, err)
				}

				if target == nil {
					return
				}

				switch e := target.(type) {
				case *ScanError:
					assertFieldValue(t, "SourceType", e.SourceType, jsonType)
					assertFieldValue(t, "TargetType", e.TargetType, incompleteInput)
				case *UnmarshalError:
					assertFieldValue(t, "Type", e.Type, jsonType)
					assertFieldValue(t, "Input", e.Input, incompleteInput)
				}
			})

			nilReturnsNilHelper(t, "nil returns nil", func() error {
				return tc.WrapFunc(nil, "XML", `<tag>`)
			})
		})
	}
}

// Original tests removed - see TestWrapScanUnmarshal for parameterized version
