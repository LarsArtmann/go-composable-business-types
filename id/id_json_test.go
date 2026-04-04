package id

import (
	"encoding/json"
	"fmt"
	"testing"
)

// assertUnmarshalError verifies that JSON unmarshaling fails with an error.
func assertUnmarshalError[B any, V comparable](t *testing.T, data string) {
	t.Helper()

	var id ID[B, V]

	err := json.Unmarshal([]byte(data), &id)
	if err == nil {
		t.Error("expected error for unmarshaling")
	}
}

func testMarshalZeroID[B any, V comparable](t *testing.T) {
	t.Helper()

	var id ID[B, V]

	data, err := json.Marshal(id)
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	if string(data) != "null" {
		t.Errorf("expected null, got %s", string(data))
	}
}

func testMarshalNonZeroID[B any, V comparable](t *testing.T, value V) {
	t.Helper()

	id := NewID[B, V](value)

	data, err := json.Marshal(id)
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	expected := fmt.Sprintf("%v", value)
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}
}

func testMarshalNonZeroStringID[B any](t *testing.T, value string) {
	t.Helper()

	id := NewID[B, string](value)

	data, err := json.Marshal(id)
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	expected := fmt.Sprintf("%q", value)
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}
}

func testUnmarshalNonZeroID[B any, V comparable](t *testing.T, jsonData string, expected V) {
	t.Helper()

	var id ID[B, V]

	err := json.Unmarshal([]byte(jsonData), &id)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	if id.Get() != expected {
		t.Errorf("expected %v, got %v", expected, id.Get())
	}
}

func TestIDJSON(t *testing.T) {
	t.Parallel()
	t.Run("string ID non-zero", func(t *testing.T) {
		t.Parallel()
		testMarshalNonZeroStringID[StringBrand](t, "abc123")
	})

	t.Run("string ID zero", func(t *testing.T) {
		t.Parallel()

		var id ID[StringBrand, string]

		data, err := json.Marshal(id)
		if err != nil {
			t.Fatalf("MarshalJSON failed: %v", err)
		}

		if string(data) != "null" {
			t.Errorf("expected null, got %s", string(data))
		}
	})

	t.Run("int64 ID non-zero", func(t *testing.T) {
		t.Parallel()
		testMarshalNonZeroID[Int64Brand, int64](t, 42)
	})

	t.Run("int64 ID zero", func(t *testing.T) {
		t.Parallel()
		testMarshalZeroID[Int64Brand, int64](t)
	})

	t.Run("int32 ID non-zero", func(t *testing.T) {
		t.Parallel()
		testMarshalNonZeroID[Int32Brand, int32](t, 42)
	})

	t.Run("uint64 ID non-zero", func(t *testing.T) {
		t.Parallel()
		testMarshalNonZeroID[Uint64Brand, uint64](t, 42)
	})
}

func TestIDUnmarshalJSON(t *testing.T) {
	t.Parallel()
	t.Run("string ID from string", func(t *testing.T) {
		t.Parallel()

		var id ID[StringBrand, string]

		err := json.Unmarshal([]byte(`"test-id"`), &id)
		if err != nil {
			t.Fatalf("UnmarshalJSON failed: %v", err)
		}

		if id.Get() != "test-id" {
			t.Errorf("expected test-id, got %s", id.Get())
		}
	})

	t.Run("string ID from null", func(t *testing.T) {
		t.Parallel()

		var id ID[StringBrand, string]

		err := json.Unmarshal([]byte("null"), &id)
		if err != nil {
			t.Fatalf("UnmarshalJSON failed: %v", err)
		}

		if !id.IsZero() {
			t.Error("expected zero ID")
		}
	})

	t.Run("int64 ID from number", func(t *testing.T) {
		t.Parallel()
		testUnmarshalNonZeroID[Int64Brand, int64](t, "42", 42)
	})

	t.Run("int32 ID from number", func(t *testing.T) {
		t.Parallel()
		testUnmarshalNonZeroID[Int32Brand, int32](t, "42", 42)
	})

	t.Run("uint64 ID from number", func(t *testing.T) {
		t.Parallel()
		testUnmarshalNonZeroID[Uint64Brand, uint64](t, "42", 42)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		t.Parallel()
		assertUnmarshalError[StringBrand, string](t, `invalid`)
	})

	t.Run("number into string ID", func(t *testing.T) {
		t.Parallel()
		assertUnmarshalError[StringBrand, string](t, "123")
	})

	t.Run("string into int64 ID", func(t *testing.T) {
		t.Parallel()
		assertUnmarshalError[Int64Brand, int64](t, `"not-a-number"`)
	})
}

func testJSONRoundTrip[B any, V comparable](t *testing.T, value V) {
	t.Helper()

	original := NewID[B, V](value)

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var restored ID[B, V]

	err = json.Unmarshal(data, &restored)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if restored.Get() != original.Get() {
		t.Errorf("expected %v, got %v", original.Get(), restored.Get())
	}
}

func TestIDJSONRoundTrip(t *testing.T) {
	t.Parallel()
	t.Run("string ID", func(t *testing.T) {
		t.Parallel()
		testJSONRoundTrip[StringBrand, string](t, "test-id")
	})

	t.Run("int64 ID", func(t *testing.T) {
		t.Parallel()
		testJSONRoundTrip[Int64Brand, int64](t, 42)
	})

	t.Run("int32 ID", func(t *testing.T) {
		t.Parallel()
		testJSONRoundTrip[Int32Brand, int32](t, 42)
	})

	t.Run("uint64 ID", func(t *testing.T) {
		t.Parallel()
		testJSONRoundTrip[Uint64Brand, uint64](t, 42)
	})
}
