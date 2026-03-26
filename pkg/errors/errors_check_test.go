package errors

import (
	"errors"
	"fmt"
	"testing"
)

// =============================================================================
// Error Checker Tests
// =============================================================================

func TestIsInvalidEmail(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			if got := IsInvalidEmail(tt.err); got != tt.want {
				t.Errorf("IsInvalidEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsInvalidURL(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			if got := IsInvalidURL(tt.err); got != tt.want {
				t.Errorf("IsInvalidURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsBoundedStringError(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			if got := IsBoundedStringError(tt.err); got != tt.want {
				t.Errorf("IsBoundedStringError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNanoIDError(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			if got := IsNanoIDError(tt.err); got != tt.want {
				t.Errorf("IsNanoIDError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsIDError(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			if got := IsIDError(tt.err); got != tt.want {
				t.Errorf("IsIDError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsParseError(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			if got := IsParseError(tt.err); got != tt.want {
				t.Errorf("IsParseError() = %v, want %v", got, tt.want)
			}
		})
	}
}
