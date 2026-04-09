// Package errors provides centralized domain-specific errors for the library.
package errors

import (
	"errors"
	"testing"
)

// testSentinelErrors runs table-driven tests for sentinel error identity.
func testSentinelErrors(t *testing.T, tests []struct {
	name string
	err  error
	want error
},
) {
	t.Helper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if !errors.Is(tt.err, tt.want) {
				t.Errorf("errors.Is() = false, want true")
			}
		})
	}
}

func testSentinelErrorSet(t *testing.T, errs ...error) {
	t.Helper()
	t.Parallel()

	tests := make([]struct {
		name string
		err  error
		want error
	}, len(errs))
	for i, err := range errs {
		tests[i] = struct {
			name string
			err  error
			want error
		}{
			name: err.Error(),
			err:  err,
			want: err,
		}
	}

	testSentinelErrors(t, tests)
}

// =============================================================================
// Sentinel Error Tests
// =============================================================================

func TestEmailSentinels(t *testing.T) {
	testSentinelErrorSet(t, ErrInvalidEmail, ErrEmailEmpty)
}

func TestURLSentinels(t *testing.T) {
	testSentinelErrorSet(t, ErrInvalidURL, ErrURLEmpty, ErrURLScheme, ErrURLHost)
}

func TestBoundedStringSentinels(t *testing.T) {
	testSentinelErrorSet(
		t,
		ErrBoundedStringMinLength,
		ErrBoundedStringMaxLength,
		ErrBoundedStringMinNegative,
		ErrBoundedStringMaxLessThanMin,
	)
}

func TestNanoIDSentinels(t *testing.T) {
	testSentinelErrorSet(t, ErrNanoIDEmpty, ErrNanoIDTooShort, ErrNanoIDTooLong, ErrNanoIDInvalid)
}

func TestIDSentinels(t *testing.T) {
	testSentinelErrorSet(t, ErrIDInvalid, ErrIDTypeNotSupported, ErrIDInsufficientData)
}
