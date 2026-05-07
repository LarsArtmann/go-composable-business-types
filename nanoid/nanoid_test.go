package nanoid

import (
	"fmt"
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

func TestNanoIDGoString(t *testing.T) {
	t.Parallel()

	nid, _ := Parse("V1StGXR8_Z5jdHi6B-myT")
	if nid.GoString() != "V1StGXR8_Z5jdHi6B-myT" {
		t.Errorf("expected V1StGXR8_Z5jdHi6B-myT, got %s", nid.GoString())
	}

	var zero NanoID
	if zero.GoString() != "" {
		t.Errorf("expected empty string for zero NanoID GoString, got %q", zero.GoString())
	}
}

func TestNanoIDMarshalTextZero(t *testing.T) {
	t.Parallel()

	var zero NanoID

	data, err := zero.MarshalText()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if data != nil {
		t.Errorf("expected nil for zero NanoID MarshalText, got %v", data)
	}
}

func TestNanoIDUnmarshalTextEmpty(t *testing.T) {
	t.Parallel()

	var id NanoID
	err := id.UnmarshalText([]byte{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !id.IsZero() {
		t.Error("expected zero NanoID after unmarshaling empty data")
	}
}

func TestNanoIDUnmarshalTextError(t *testing.T) {
	t.Parallel()

	var id NanoID
	err := id.UnmarshalText([]byte("bad!"))
	if err == nil {
		t.Error("expected error for invalid NanoID")
	}
}

func TestNanoIDScan(t *testing.T) {
	t.Parallel()

	// Valid string
	var id1 NanoID
	err := id1.Scan("V1StGXR8_Z5jdHi6B-myT")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if id1.String() != "V1StGXR8_Z5jdHi6B-myT" {
		t.Errorf("expected V1StGXR8_Z5jdHi6B-myT, got %s", id1.String())
	}

	// Valid []byte
	var id2 NanoID
	err := id2.Scan([]byte("V1StGXR8_Z5jdHi6B-myT"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if id2.String() != "V1StGXR8_Z5jdHi6B-myT" {
		t.Errorf("expected V1StGXR8_Z5jdHi6B-myT, got %s", id2.String())
	}

	// Empty string
	var id3 NanoID
	err := id3.Scan("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !id3.IsZero() {
		t.Error("expected zero NanoID after scanning empty string")
	}

	// nil
	var id4 NanoID
	err := id4.Scan(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !id4.IsZero() {
		t.Error("expected zero NanoID after scanning nil")
	}

	// Invalid type
	var id5 NanoID
	err := id5.Scan(123)
	if err == nil {
		t.Error("expected error for invalid type")
	}

	// Invalid NanoID value
	var id6 NanoID
	err := id6.Scan("bad!")
	if err == nil {
		t.Error("expected error for invalid NanoID value")
	}
}

func TestNanoIDScanNilReceiver(t *testing.T) {
	t.Parallel()

	var id *NanoID
	err := id.Scan("V1StGXR8_Z5jdHi6B-myT")
	if err == nil {
		t.Error("expected error when scanning into nil receiver")
	}
}

func TestNanoIDValue(t *testing.T) {
	t.Parallel()

	// Non-zero value
	nid, _ := Parse("V1StGXR8_Z5jdHi6B-myT")

	val, err := nid.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if val != "V1StGXR8_Z5jdHi6B-myT" {
		t.Errorf("expected V1StGXR8_Z5jdHi6B-myT, got %v", val)
	}

	// Zero value
	var zero NanoID

	val, err = zero.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if val != nil {
		t.Errorf("expected nil for zero NanoID Value, got %v", val)
	}
}

func TestNanoIDGoStringFormat(t *testing.T) {
	t.Parallel()

	nid, _ := Parse("V1StGXR8_Z5jdHi6B-myT")

	gs := fmt.Sprintf("%#v", nid)
	if gs != "V1StGXR8_Z5jdHi6B-myT" {
		t.Errorf("expected GoString to match, got %q", gs)
	}
}
