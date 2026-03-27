package id

import (
	"encoding/json"
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

func TestIDJSON(t *testing.T) {
	t.Parallel()
	t.Run("string ID non-zero", func(t *testing.T) {
		t.Parallel()
		id := NewID[StringBrand]("abc123")
		data, err := json.Marshal(id)
		if err != nil {
			t.Fatalf("MarshalJSON failed: %v", err)
		}
		if string(data) != `"abc123"` {
			t.Errorf("expected \"abc123\", got %s", string(data))
		}
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
		id := NewID[Int64Brand, int64](42)
		data, err := json.Marshal(id)
		if err != nil {
			t.Fatalf("MarshalJSON failed: %v", err)
		}
		if string(data) != "42" {
			t.Errorf("expected 42, got %s", string(data))
		}
	})

	t.Run("int64 ID zero", func(t *testing.T) {
		t.Parallel()
		var id ID[Int64Brand, int64]
		data, err := json.Marshal(id)
		if err != nil {
			t.Fatalf("MarshalJSON failed: %v", err)
		}
		if string(data) != "null" {
			t.Errorf("expected null, got %s", string(data))
		}
	})

	t.Run("int32 ID non-zero", func(t *testing.T) {
		t.Parallel()
		id := NewID[Int32Brand, int32](42)
		data, err := json.Marshal(id)
		if err != nil {
			t.Fatalf("MarshalJSON failed: %v", err)
		}
		if string(data) != "42" {
			t.Errorf("expected 42, got %s", string(data))
		}
	})

	t.Run("uint64 ID non-zero", func(t *testing.T) {
		t.Parallel()
		id := NewID[Uint64Brand, uint64](42)
		data, err := json.Marshal(id)
		if err != nil {
			t.Fatalf("MarshalJSON failed: %v", err)
		}
		if string(data) != "42" {
			t.Errorf("expected 42, got %s", string(data))
		}
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
		var id ID[Int64Brand, int64]
		err := json.Unmarshal([]byte("42"), &id)
		if err != nil {
			t.Fatalf("UnmarshalJSON failed: %v", err)
		}
		if id.Get() != 42 {
			t.Errorf("expected 42, got %d", id.Get())
		}
	})

	t.Run("int32 ID from number", func(t *testing.T) {
		t.Parallel()
		var id ID[Int32Brand, int32]
		err := json.Unmarshal([]byte("42"), &id)
		if err != nil {
			t.Fatalf("UnmarshalJSON failed: %v", err)
		}
		if id.Get() != 42 {
			t.Errorf("expected 42, got %d", id.Get())
		}
	})

	t.Run("uint64 ID from number", func(t *testing.T) {
		t.Parallel()
		var id ID[Uint64Brand, uint64]
		err := json.Unmarshal([]byte("42"), &id)
		if err != nil {
			t.Fatalf("UnmarshalJSON failed: %v", err)
		}
		if id.Get() != 42 {
			t.Errorf("expected 42, got %d", id.Get())
		}
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
