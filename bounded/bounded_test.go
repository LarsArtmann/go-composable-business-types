package bounded

import (
	"testing"
)

func TestNewBoundedString(t *testing.T) {
	tests := []struct {
		name    string
		minLen  int
		maxLen  int
		value   string
		wantErr bool
	}{
		{"valid", 1, 100, "hello", false},
		{"at min", 5, 10, "hello", false},
		{"at max", 1, 5, "hello", false},
		{"empty when min 0", 0, 10, "", false},
		{"too short", 10, 20, "hi", true},
		{"too long", 1, 3, "hello", true},
		{"negative min", -1, 10, "test", true},
		{"max less than min", 10, 5, "test", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs, err := NewBoundedString(tt.minLen, tt.maxLen, tt.value)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if bs.String() != tt.value {
				t.Errorf("expected %s, got %s", tt.value, bs.String())
			}
		})
	}
}

func TestBoundedStringLen(t *testing.T) {
	bs, _ := NewBoundedString(1, 100, "hello")
	if bs.Len() != 5 {
		t.Errorf("expected length 5, got %d", bs.Len())
	}

	// Unicode test
	bs2, _ := NewBoundedString(1, 100, "héllo")
	if bs2.Len() != 5 {
		t.Errorf("expected length 5 for unicode string, got %d", bs2.Len())
	}
}

func TestBoundedStringOf(t *testing.T) {
	NewName := BoundedStringOf(1, 100)
	name, err := NewName("John Doe")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if name.String() != "John Doe" {
		t.Errorf("expected 'John Doe', got %s", name.String())
	}

	// Test with invalid name
	_, err = NewName("")
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestNonEmptyString(t *testing.T) {
	s, err := NonEmptyString(100, "hello")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if s.String() != "hello" {
		t.Errorf("expected 'hello', got %s", s.String())
	}

	// Empty string should fail
	_, err = NonEmptyString(100, "")
	if err == nil {
		t.Error("expected error for empty string")
	}
}

func TestTrimmedBoundedString(t *testing.T) {
	s, err := TrimmedBoundedString(1, 100, "  hello  ")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if s.String() != "hello" {
		t.Errorf("expected 'hello' (trimmed), got %s", s.String())
	}
}

func TestBoundedStringIsZero(t *testing.T) {
	var zero BoundedString
	if !zero.IsZero() {
		t.Error("expected zero BoundedString to be zero")
	}

	nonZero, _ := NewBoundedString(1, 100, "test")
	if nonZero.IsZero() {
		t.Error("expected non-zero BoundedString to not be zero")
	}
}

func TestBoundedStringBounds(t *testing.T) {
	bs, _ := NewBoundedString(5, 10, "hello")
	if bs.MinLen() != 5 {
		t.Errorf("expected MinLen 5, got %d", bs.MinLen())
	}
	if bs.MaxLen() != 10 {
		t.Errorf("expected MaxLen 10, got %d", bs.MaxLen())
	}
	// "hello" is length 5, so it's at min but not at max
	if !bs.IsMinLength() {
		t.Error("expected to be at min")
	}
	if bs.IsMaxLength() {
		t.Error("should not be at max")
	}
}
