package nanoid

import (
	"testing"
)

func TestNewNanoID(t *testing.T) {
	t.Parallel()
	id := NewNanoID()
	if id.IsZero() {
		t.Error("expected non-zero NanoID")
	}
	if len(id.String()) != DefaultNanoIDLength {
		t.Errorf("expected length %d, got %d", DefaultNanoIDLength, len(id.String()))
	}
}

func TestNewNanoIDWithLength(t *testing.T) {
	t.Parallel()
	id := NewNanoIDWithLength(10)
	if len(id.String()) != 10 {
		t.Errorf("expected length 10, got %d", len(id.String()))
	}
}

func TestParseNanoID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid", "V1StGXR8_Z5jdHi6B-myT", false},
		{"empty", "", true},
		{"too short", "abc", true},
		{"too long", string(make([]byte, 257)), true},
		{"invalid chars", "hello@world!", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			id, err := ParseNanoID(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if id.String() != tt.input {
				t.Errorf("expected %s, got %s", tt.input, id.String())
			}
		})
	}
}

func TestParseNanoIDError(t *testing.T) {
	t.Parallel()
	_, err := ParseNanoID("invalid")
	if err == nil {
		t.Error("expected error for invalid NanoID")
	}
}

func TestNanoIDIsZero(t *testing.T) {
	t.Parallel()
	var zero NanoID
	if !zero.IsZero() {
		t.Error("expected zero NanoID to be zero")
	}

	nonZero := NewNanoID()
	if nonZero.IsZero() {
		t.Error("expected non-zero NanoID to not be zero")
	}
}

func TestNanoIDJSON(t *testing.T) {
	t.Parallel()
	id, _ := ParseNanoID("V1StGXR8_Z5jdHi6B-myT")

	data, err := id.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText failed: %v", err)
	}
	if string(data) != "V1StGXR8_Z5jdHi6B-myT" {
		t.Errorf("expected V1StGXR8_Z5jdHi6B-myT, got %s", string(data))
	}

	// Test unmarshal
	var parsed NanoID
	err = parsed.UnmarshalText([]byte("V1StGXR8_Z5jdHi6B-myT"))
	if err != nil {
		t.Fatalf("UnmarshalText failed: %v", err)
	}
	if parsed.String() != "V1StGXR8_Z5jdHi6B-myT" {
		t.Errorf("expected V1StGXR8_Z5jdHi6B-myT, got %s", parsed.String())
	}
}
