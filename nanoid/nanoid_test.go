package nanoid

import (
	"testing"
)

func TestNewNanoId(t *testing.T) {
	id := NewNanoId()
	if id.IsZero() {
		t.Error("expected non-zero NanoId")
	}
	if len(id.String()) != DefaultNanoIdLength {
		t.Errorf("expected length %d, got %d", DefaultNanoIdLength, len(id.String()))
	}
}

func TestNewNanoIdWithLength(t *testing.T) {
	id := NewNanoIdWithLength(10)
	if len(id.String()) != 10 {
		t.Errorf("expected length 10, got %d", len(id.String()))
	}
}

func TestParseNanoId(t *testing.T) {
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
			id, err := ParseNanoId(tt.input)
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

func TestMustParseNanoId(t *testing.T) {
	id := MustParseNanoId("V1StGXR8_Z5jdHi6B-myT")
	if id.String() != "V1StGXR8_Z5jdHi6B-myT" {
		t.Errorf("expected V1StGXR8_Z5jdHi6B-myT, got %s", id.String())
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid NanoId")
		}
	}()
	MustParseNanoId("invalid")
}

func TestNanoIdIsZero(t *testing.T) {
	var zero NanoId
	if !zero.IsZero() {
		t.Error("expected zero NanoId to be zero")
	}

	nonZero := NewNanoId()
	if nonZero.IsZero() {
		t.Error("expected non-zero NanoId to not be zero")
	}
}

func TestNanoIdJSON(t *testing.T) {
	id := MustParseNanoId("V1StGXR8_Z5jdHi6B-myT")

	data, err := id.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText failed: %v", err)
	}
	if string(data) != "V1StGXR8_Z5jdHi6B-myT" {
		t.Errorf("expected V1StGXR8_Z5jdHi6B-myT, got %s", string(data))
	}

	// Test unmarshal
	var parsed NanoId
	err = parsed.UnmarshalText([]byte("V1StGXR8_Z5jdHi6B-myT"))
	if err != nil {
		t.Fatalf("UnmarshalText failed: %v", err)
	}
	if parsed.String() != "V1StGXR8_Z5jdHi6B-myT" {
		t.Errorf("expected V1StGXR8_Z5jdHi6B-myT, got %s", parsed.String())
	}
}
