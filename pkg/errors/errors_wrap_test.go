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
					var scanErr *ScanError
					if !errors.As(err, &scanErr) {
						t.Error("expected errors.As() to succeed")
					}

					target = scanErr
				case "*UnmarshalError":
					var unmarshalErr *UnmarshalError
					if !errors.As(err, &unmarshalErr) {
						t.Error("expected errors.As() to succeed")
					}

					target = unmarshalErr
				}

				if target == nil {
					return
				}

				switch e := target.(type) {
				case *ScanError:
					if e.SourceType != "JSON" {
						t.Errorf("SourceType = %q, want JSON", e.SourceType)
					}

					if e.TargetType != `{incomplete` {
						t.Errorf("TargetType = %q, want {incomplete", e.TargetType)
					}
				case *UnmarshalError:
					if e.Type != "JSON" {
						t.Errorf("Type = %q, want JSON", e.Type)
					}

					if e.Input != `{incomplete` {
						t.Errorf("Input = %q, want {incomplete", e.Input)
					}
				}
			})

			nilReturnsNilHelper(t, "nil returns nil", func() error {
				return tc.WrapFunc(nil, "XML", `<tag>`)
			})
		})
	}
}

// Original tests removed - see TestWrapScanUnmarshal for parameterized version
