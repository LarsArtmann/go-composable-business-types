package nanoid

import (
	"testing"

	"github.com/larsartmann/go-composable-business-types/testutil"
)

func TestNew(t *testing.T) {
	t.Parallel()

	id := New()
	if id.IsZero() {
		t.Error("expected non-zero NanoID")
	}

	if len(id.String()) != DefaultNanoIDLength {
		t.Errorf("expected length %d, got %d", DefaultNanoIDLength, len(id.String()))
	}
}

func TestNewWithLength(t *testing.T) {
	t.Parallel()

	id := NewWithLength(10)
	if len(id.String()) != 10 {
		t.Errorf("expected length 10, got %d", len(id.String()))
	}
}

func TestParse(t *testing.T) {
	t.Parallel()

	testutil.RunParseTests(t, "NanoID", []testutil.ParseTestCase[NanoID]{
		{Name: "valid", Input: "V1StGXR8_Z5jdHi6B-myT", WantErr: false},
		{Name: "empty", Input: "", WantErr: true},
		{Name: "too short", Input: "abc", WantErr: true},
		{Name: "too long", Input: string(make([]byte, 257)), WantErr: true},
		{Name: "invalid chars", Input: "hello@world!", WantErr: true},
	}, Parse)
}

func TestParseError(t *testing.T) {
	testutil.RunParseErrorTest(t, "NanoID", Parse)
}

func TestNanoIDIsZero(t *testing.T) {
	t.Parallel()

	var zero NanoID
	if !zero.IsZero() {
		t.Error("expected zero NanoID to be zero")
	}

	nonZero := New()
	if nonZero.IsZero() {
		t.Error("expected non-zero NanoID to not be zero")
	}
}

func TestNanoIDJSON(t *testing.T) {
	t.Parallel()

	id, _ := Parse("V1StGXR8_Z5jdHi6B-myT")

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
