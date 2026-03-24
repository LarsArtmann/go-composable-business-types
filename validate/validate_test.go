package validate

import (
	"errors"
	"testing"
)

// mockValidator is a test helper that implements Validator.
type mockValidator struct {
	err error
}

func (m mockValidator) Validate() error { return m.err }

func TestIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		validator Validator
		want      bool
	}{
		{
			name:      "valid value returns true",
			validator: mockValidator{err: nil},
			want:      true,
		},
		{
			name:      "invalid value returns false",
			validator: mockValidator{err: errors.New("validation failed")},
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := IsValid(tt.validator); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidatorInterface(t *testing.T) {
	t.Parallel()

	// Ensure mockValidator implements Validator
	var _ Validator = mockValidator{}

	// Test that the interface is satisfied
	var v Validator = mockValidator{err: nil}
	if err := v.Validate(); err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	v = mockValidator{err: errors.New("test error")}
	if err := v.Validate(); err == nil {
		t.Error("expected error, got nil")
	}
}
