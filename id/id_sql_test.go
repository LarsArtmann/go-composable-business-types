package id

import (
	"testing"
)

func TestIDScan(t *testing.T) {
	t.Parallel()
	t.Run("string ID from string", func(t *testing.T) {
		t.Parallel()
		testScanRoundTrip[StringBrand, string](t, "test-id", "test-id")
	})

	t.Run("string ID from []byte", func(t *testing.T) {
		t.Parallel()
		var id ID[StringBrand, string]
		err := id.Scan([]byte("test-id"))
		if err != nil {
			t.Fatalf("Scan failed: %v", err)
		}
		if id.Get() != "test-id" {
			t.Errorf("expected test-id, got %s", id.Get())
		}
	})

	t.Run("string ID from nil", func(t *testing.T) {
		t.Parallel()
		var id ID[StringBrand, string]
		err := id.Scan(nil)
		if err != nil {
			t.Fatalf("Scan failed: %v", err)
		}
		if !id.IsZero() {
			t.Error("expected zero ID")
		}
	})

	t.Run("string ID from invalid type", func(t *testing.T) {
		t.Parallel()
		var id ID[StringBrand, string]
		err := id.Scan(123)
		if err == nil {
			t.Error("expected error for invalid type")
		}
	})

	t.Run("int64 ID from int64", func(t *testing.T) {
		t.Parallel()
		testScanRoundTrip[Int64Brand, int64](t, int64(42), int64(42))
	})

	t.Run("int64 ID from int", func(t *testing.T) {
		t.Parallel()
		var id ID[Int64Brand, int64]
		err := id.Scan(42)
		if err != nil {
			t.Fatalf("Scan failed: %v", err)
		}
		if id.Get() != 42 {
			t.Errorf("expected 42, got %d", id.Get())
		}
	})

	t.Run("int64 ID from float64", func(t *testing.T) {
		t.Parallel()
		var id ID[Int64Brand, int64]
		err := id.Scan(float64(42))
		if err != nil {
			t.Fatalf("Scan failed: %v", err)
		}
		if id.Get() != 42 {
			t.Errorf("expected 42, got %d", id.Get())
		}
	})

	t.Run("int64 ID from nil", func(t *testing.T) {
		t.Parallel()
		var id ID[Int64Brand, int64]
		err := id.Scan(nil)
		if err != nil {
			t.Fatalf("Scan failed: %v", err)
		}
		if !id.IsZero() {
			t.Error("expected zero ID")
		}
	})

	t.Run("int64 ID from invalid type", func(t *testing.T) {
		t.Parallel()
		var id ID[Int64Brand, int64]
		err := id.Scan("not-a-number")
		if err == nil {
			t.Error("expected error for invalid type")
		}
	})

	t.Run("int32 ID from int64", func(t *testing.T) {
		t.Parallel()
		testScanRoundTrip[Int32Brand, int32](t, int64(42), int32(42))
	})

	t.Run("uint64 ID from int64", func(t *testing.T) {
		t.Parallel()
		testScanRoundTrip[Uint64Brand, uint64](t, int64(42), uint64(42))
	})
}

func testScanRoundTrip[B any, V comparable](t *testing.T, input any, expected V) {
	t.Helper()
	var id ID[B, V]
	err := id.Scan(input)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}
	if id.Get() != expected {
		t.Errorf("expected %v, got %v", expected, id.Get())
	}
}

func TestIDValue(t *testing.T) {
	t.Parallel()
	t.Run("string ID non-zero", func(t *testing.T) {
		t.Parallel()
		id := NewID[StringBrand]("test-id")
		val, err := id.Value()
		if err != nil {
			t.Fatalf("Value failed: %v", err)
		}
		if val != "test-id" {
			t.Errorf("expected test-id, got %v", val)
		}
	})

	t.Run("string ID zero", func(t *testing.T) {
		t.Parallel()
		var id ID[StringBrand, string]
		val, err := id.Value()
		if err != nil {
			t.Fatalf("Value failed: %v", err)
		}
		if val != nil {
			t.Errorf("expected nil, got %v", val)
		}
	})

	t.Run("int64 ID non-zero", func(t *testing.T) {
		t.Parallel()
		id := NewID[Int64Brand, int64](42)
		val, err := id.Value()
		if err != nil {
			t.Fatalf("Value failed: %v", err)
		}
		if val != int64(42) {
			t.Errorf("expected 42, got %v", val)
		}
	})

	t.Run("int64 ID zero", func(t *testing.T) {
		t.Parallel()
		var id ID[Int64Brand, int64]
		val, err := id.Value()
		if err != nil {
			t.Fatalf("Value failed: %v", err)
		}
		if val != nil {
			t.Errorf("expected nil, got %v", val)
		}
	})

	t.Run("int32 ID non-zero", func(t *testing.T) {
		t.Parallel()
		id := NewID[Int32Brand, int32](42)
		val, err := id.Value()
		if err != nil {
			t.Fatalf("Value failed: %v", err)
		}
		if val != int64(42) {
			t.Errorf("expected 42, got %v", val)
		}
	})

	t.Run("uint64 ID non-zero", func(t *testing.T) {
		t.Parallel()
		id := NewID[Uint64Brand, uint64](42)
		val, err := id.Value()
		if err != nil {
			t.Fatalf("Value failed: %v", err)
		}
		if val != int64(42) {
			t.Errorf("expected 42, got %v", val)
		}
	})
}
