package id

import (
	"testing"
)

func TestNewID(t *testing.T) {
	id := NewID[struct{}, string]("user-123")
	if id.Get() != "user-123" {
		t.Errorf("expected user-123, got %s", id.Get())
	}
	if id.IsZero() {
		t.Error("expected non-zero id")
	}
}

func TestNewIDInt(t *testing.T) {
	id := NewID[struct{}, int](42)
	if id.Get() != 42 {
		t.Errorf("expected 42, got %d", id.Get())
	}
}

func TestIDIsZero(t *testing.T) {
	var zeroID ID[struct{}, string]
	if !zeroID.IsZero() {
		t.Error("expected zero ID to be zero")
	}

	nonZeroID := NewID[struct{}, string]("test")
	if nonZeroID.IsZero() {
		t.Error("expected non-zero ID to not be zero")
	}
}

func TestIDString(t *testing.T) {
	id := NewID[struct{}, string]("test-id")
	if id.String() != "test-id" {
		t.Errorf("expected test-id, got %s", id.String())
	}
}

func TestIDJSON(t *testing.T) {
	type TestBrand struct{}

	tests := []struct {
		name     string
		input    ID[TestBrand, string]
		expected string
	}{
		{"non-zero", NewID[TestBrand, string]("abc123"), `"abc123"`},
		{"zero", ID[TestBrand, string]{}, "null"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := tt.input.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON failed: %v", err)
			}
			if string(data) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(data))
			}
		})
	}
}

func TestIDUnmarshalJSON(t *testing.T) {
	type TestBrand struct{}

	var id ID[TestBrand, string]
	err := id.UnmarshalJSON([]byte(`"test-id"`))
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}
	if id.Get() != "test-id" {
		t.Errorf("expected test-id, got %s", id.Get())
	}

	// Test null
	var nullID ID[TestBrand, string]
	err = nullID.UnmarshalJSON([]byte("null"))
	if err != nil {
		t.Fatalf("UnmarshalJSON null failed: %v", err)
	}
	if !nullID.IsZero() {
		t.Error("expected zero ID after unmarshaling null")
	}
}
