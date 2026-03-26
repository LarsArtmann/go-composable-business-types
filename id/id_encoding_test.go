package id

import (
	"bytes"
	"encoding/gob"
	"testing"
)

func TestIDText(t *testing.T) {
	t.Parallel()
	t.Run("marshal non-zero string", func(t *testing.T) {
		t.Parallel()
		id := NewID[StringBrand]("test-id")
		data, err := id.MarshalText()
		if err != nil {
			t.Fatalf("MarshalText failed: %v", err)
		}
		if string(data) != "test-id" {
			t.Errorf("expected test-id, got %s", string(data))
		}
	})

	t.Run("marshal zero string", func(t *testing.T) {
		t.Parallel()
		var id ID[StringBrand, string]
		data, err := id.MarshalText()
		if err != nil {
			t.Fatalf("MarshalText failed: %v", err)
		}
		if data != nil {
			t.Errorf("expected nil, got %s", string(data))
		}
	})

	t.Run("marshal int64", func(t *testing.T) {
		t.Parallel()
		id := NewID[Int64Brand, int64](42)
		data, err := id.MarshalText()
		if err != nil {
			t.Fatalf("MarshalText failed: %v", err)
		}
		if string(data) != "42" {
			t.Errorf("expected 42, got %s", string(data))
		}
	})

	t.Run("unmarshal valid string", func(t *testing.T) {
		t.Parallel()
		testUnmarshalTextRoundTrip[StringBrand, string](t, "test-id", "test-id")
	})

	t.Run("unmarshal empty", func(t *testing.T) {
		t.Parallel()
		var id ID[StringBrand, string]
		err := id.UnmarshalText([]byte{})
		if err != nil {
			t.Fatalf("UnmarshalText failed: %v", err)
		}
		if !id.IsZero() {
			t.Error("expected zero ID")
		}
	})

	t.Run("unmarshal int64", func(t *testing.T) {
		t.Parallel()
		testUnmarshalTextRoundTrip[Int64Brand, int64](t, "42", 42)
	})

	t.Run("unmarshal uint64", func(t *testing.T) {
		t.Parallel()
		testUnmarshalTextRoundTrip[Uint64Brand, uint64](t, "42", 42)
	})
}

func testUnmarshalTextRoundTrip[B any, V comparable](t *testing.T, input string, expected V) {
	t.Helper()
	var id ID[B, V]
	err := id.UnmarshalText([]byte(input))
	if err != nil {
		t.Fatalf("UnmarshalText failed: %v", err)
	}
	if id.Get() != expected {
		t.Errorf("expected %v, got %v", expected, id.Get())
	}
}

func testBinaryRoundTrip[B any, V comparable](t *testing.T, value V) {
	t.Helper()
	original := NewID[B, V](value)
	data, err := original.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}

	var restored ID[B, V]
	err = restored.UnmarshalBinary(data)
	if err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}
	if restored.Get() != original.Get() {
		t.Errorf("expected %v, got %v", original.Get(), restored.Get())
	}
}

func TestIDBinary(t *testing.T) {
	t.Parallel()
	t.Run("string ID", func(t *testing.T) {
		t.Parallel()
		testBinaryRoundTrip[StringBrand, string](t, "test-id")
	})

	t.Run("int64 ID", func(t *testing.T) {
		t.Parallel()
		testBinaryRoundTrip[Int64Brand, int64](t, 42)
	})

	t.Run("int32 ID", func(t *testing.T) {
		t.Parallel()
		testBinaryRoundTrip[Int32Brand, int32](t, 42)
	})

	t.Run("uint64 ID", func(t *testing.T) {
		t.Parallel()
		testBinaryRoundTrip[Uint64Brand, uint64](t, 42)
	})

	t.Run("zero ID", func(t *testing.T) {
		t.Parallel()
		var original ID[StringBrand, string]
		data, err := original.MarshalBinary()
		if err != nil {
			t.Fatalf("MarshalBinary failed: %v", err)
		}
		if data != nil {
			t.Errorf("expected nil, got %v", data)
		}

		var restored ID[StringBrand, string]
		err = restored.UnmarshalBinary(nil)
		if err != nil {
			t.Fatalf("UnmarshalBinary failed: %v", err)
		}
		if !restored.IsZero() {
			t.Error("expected zero ID")
		}
	})
}

func TestIDGob(t *testing.T) {
	t.Parallel()
	t.Run("string ID", func(t *testing.T) {
		t.Parallel()
		original := NewID[StringBrand]("test-id")
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(original)
		if err != nil {
			t.Fatalf("GobEncode failed: %v", err)
		}

		var restored ID[StringBrand, string]
		dec := gob.NewDecoder(&buf)
		err = dec.Decode(&restored)
		if err != nil {
			t.Fatalf("GobDecode failed: %v", err)
		}
		if restored.Get() != original.Get() {
			t.Errorf("expected %s, got %s", original.Get(), restored.Get())
		}
	})

	t.Run("int64 ID", func(t *testing.T) {
		t.Parallel()
		original := NewID[Int64Brand, int64](42)
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(original)
		if err != nil {
			t.Fatalf("GobEncode failed: %v", err)
		}

		var restored ID[Int64Brand, int64]
		dec := gob.NewDecoder(&buf)
		err = dec.Decode(&restored)
		if err != nil {
			t.Fatalf("GobDecode failed: %v", err)
		}
		if restored.Get() != original.Get() {
			t.Errorf("expected %d, got %d", original.Get(), restored.Get())
		}
	})
}
