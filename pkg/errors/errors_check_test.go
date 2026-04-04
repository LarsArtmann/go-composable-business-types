package errors

import (
	"errors"
	"fmt"
	"testing"
)

// testErrorChecker is a table-driven test helper for error checking functions.
func testErrorChecker[Fn func(error) bool](
	t *testing.T,
	name string,
	checker Fn,
	testCases []struct {
		name string
		err  error
		want bool
	},
) {
	t.Helper()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := checker(tc.err); got != tc.want {
				t.Errorf("%s() = %v, want %v", name, got, tc.want)
			}
		})
	}
}

// =============================================================================
// Error Checker Tests
// =============================================================================

func TestIsInvalidEmail(t *testing.T) {
	testErrorChecker(t, "IsInvalidEmail", IsInvalidEmail, []struct {
		name string
		err  error
		want bool
	}{
		{"ErrInvalidEmail", ErrInvalidEmail, true},
		{"ErrEmailEmpty", ErrEmailEmpty, true},
		{"other error", errors.New("other"), false},
		{"wrapped email error", fmt.Errorf("wrapped: %w", ErrInvalidEmail), true},
	})
}

func TestIsInvalidURL(t *testing.T) {
	testErrorChecker(t, "IsInvalidURL", IsInvalidURL, []struct {
		name string
		err  error
		want bool
	}{
		{"ErrInvalidURL", ErrInvalidURL, true},
		{"ErrURLEmpty", ErrURLEmpty, true},
		{"ErrURLScheme", ErrURLScheme, true},
		{"ErrURLHost", ErrURLHost, true},
		{"other error", errors.New("other"), false},
	})
}

func TestIsBoundedStringError(t *testing.T) {
	testErrorChecker(t, "IsBoundedStringError", IsBoundedStringError, []struct {
		name string
		err  error
		want bool
	}{
		{"ErrBoundedStringMinLength", ErrBoundedStringMinLength, true},
		{"ErrBoundedStringMaxLength", ErrBoundedStringMaxLength, true},
		{"ErrBoundedStringMinNegative", ErrBoundedStringMinNegative, true},
		{"ErrBoundedStringMaxLessThanMin", ErrBoundedStringMaxLessThanMin, true},
		{"other error", errors.New("other"), false},
	})
}

func TestIsNanoIDError(t *testing.T) {
	testErrorChecker(t, "IsNanoIDError", IsNanoIDError, []struct {
		name string
		err  error
		want bool
	}{
		{"ErrNanoIDEmpty", ErrNanoIDEmpty, true},
		{"ErrNanoIDTooShort", ErrNanoIDTooShort, true},
		{"ErrNanoIDTooLong", ErrNanoIDTooLong, true},
		{"ErrNanoIDInvalid", ErrNanoIDInvalid, true},
		{"other error", errors.New("other"), false},
	})
}

func TestIsIDError(t *testing.T) {
	testErrorChecker(t, "IsIDError", IsIDError, []struct {
		name string
		err  error
		want bool
	}{
		{"ErrIDInvalid", ErrIDInvalid, true},
		{"ErrIDTypeNotSupported", ErrIDTypeNotSupported, true},
		{"ErrIDInsufficientData", ErrIDInsufficientData, true},
		{"other error", errors.New("other"), false},
	})
}

func TestIsParseError(t *testing.T) {
	testErrorChecker(t, "IsParseError", IsParseError, []struct {
		name string
		err  error
		want bool
	}{
		{"ErrMalformedInput", ErrMalformedInput, true},
		{"ErrUnsupportedType", ErrUnsupportedType, true},
		{"ErrInvalidJSON", ErrInvalidJSON, true},
		{"other error", errors.New("other"), false},
	})
}
