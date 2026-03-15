package id

import (
	"encoding/json"
	"testing"
)

type (
	StringBrand struct{}
	Int64Brand  struct{}
)

func TestNewID(t *testing.T) {
	id := NewID[StringBrand]("user-123")
	if id.Get() != "user-123" {
		t.Errorf("expected user-123, got %s", id.Get())
	}
	if id.IsZero() {
		t.Error("expected non-zero id")
	}
}

func TestNewIDInt64(t *testing.T) {
	id := NewID[Int64Brand, int64](42)
	if id.Get() != 42 {
		t.Errorf("expected 42, got %d", id.Get())
	}
}

func TestIDIsZero(t *testing.T) {
	var zeroID ID[StringBrand, string]
	if !zeroID.IsZero() {
		t.Error("expected zero ID to be zero")
	}

	nonZeroID := NewID[StringBrand]("test")
	if nonZeroID.IsZero() {
		t.Error("expected non-zero ID to not be zero")
	}
}

func TestIDString(t *testing.T) {
	id := NewID[StringBrand]("test-id")
	if id.String() != "test-id" {
		t.Errorf("expected test-id, got %s", id.String())
	}
}

func TestIDGoString(t *testing.T) {
	id := NewID[StringBrand]("test-id")
	if id.GoString() != "test-id" {
		t.Errorf("expected test-id, got %s", id.GoString())
	}
}

func TestIDJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    ID[StringBrand, string]
		expected string
	}{
		{"non-zero", NewID[StringBrand]("abc123"), `"abc123"`},
		{"zero", ID[StringBrand, string]{}, "null"},
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
	t.Run("valid string", func(t *testing.T) {
		var id ID[StringBrand, string]
		err := id.UnmarshalJSON([]byte(`"test-id"`))
		if err != nil {
			t.Fatalf("UnmarshalJSON failed: %v", err)
		}
		if id.Get() != "test-id" {
			t.Errorf("expected test-id, got %s", id.Get())
		}
	})

	t.Run("null", func(t *testing.T) {
		var id ID[StringBrand, string]
		err := id.UnmarshalJSON([]byte("null"))
		if err != nil {
			t.Fatalf("UnmarshalJSON null failed: %v", err)
		}
		if !id.IsZero() {
			t.Error("expected zero ID after unmarshaling null")
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		var id ID[StringBrand, string]
		err := id.UnmarshalJSON([]byte(`invalid`))
		if err == nil {
			t.Error("expected error for invalid JSON")
		}
	})

	t.Run("number into string ID", func(t *testing.T) {
		var id ID[StringBrand, string]
		err := id.UnmarshalJSON([]byte(`123`))
		if err == nil {
			t.Error("expected error for number into string ID")
		}
	})
}

func TestIDText(t *testing.T) {
	t.Run("marshal non-zero", func(t *testing.T) {
		id := NewID[StringBrand]("test-id")
		data, err := id.MarshalText()
		if err != nil {
			t.Fatalf("MarshalText failed: %v", err)
		}
		if string(data) != "test-id" {
			t.Errorf("expected test-id, got %s", string(data))
		}
	})

	t.Run("marshal zero", func(t *testing.T) {
		var id ID[StringBrand, string]
		data, err := id.MarshalText()
		if err != nil {
			t.Fatalf("MarshalText failed: %v", err)
		}
		if data != nil {
			t.Errorf("expected nil, got %s", string(data))
		}
	})

	t.Run("unmarshal valid", func(t *testing.T) {
		var id ID[StringBrand, string]
		err := id.UnmarshalText([]byte("test-id"))
		if err != nil {
			t.Fatalf("UnmarshalText failed: %v", err)
		}
		if id.Get() != "test-id" {
			t.Errorf("expected test-id, got %s", id.Get())
		}
	})

	t.Run("unmarshal empty", func(t *testing.T) {
		var id ID[StringBrand, string]
		err := id.UnmarshalText([]byte{})
		if err != nil {
			t.Fatalf("UnmarshalText failed: %v", err)
		}
		if !id.IsZero() {
			t.Error("expected zero ID after unmarshaling empty")
		}
	})
}

func TestIDScan(t *testing.T) {
	t.Run("string ID from string", func(t *testing.T) {
		var id ID[StringBrand, string]
		err := id.Scan("test-id")
		if err != nil {
			t.Fatalf("Scan failed: %v", err)
		}
		if id.Get() != "test-id" {
			t.Errorf("expected test-id, got %s", id.Get())
		}
	})

	t.Run("string ID from []byte", func(t *testing.T) {
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
		var id ID[StringBrand, string]
		err := id.Scan(nil)
		if err != nil {
			t.Fatalf("Scan failed: %v", err)
		}
		if !id.IsZero() {
			t.Error("expected zero ID after scanning nil")
		}
	})

	t.Run("string ID from invalid type", func(t *testing.T) {
		var id ID[StringBrand, string]
		err := id.Scan(123)
		if err == nil {
			t.Error("expected error for invalid type")
		}
	})

	t.Run("int64 ID from int64", func(t *testing.T) {
		var id ID[Int64Brand, int64]
		err := id.Scan(int64(42))
		if err != nil {
			t.Fatalf("Scan failed: %v", err)
		}
		if id.Get() != 42 {
			t.Errorf("expected 42, got %d", id.Get())
		}
	})

	t.Run("int64 ID from int", func(t *testing.T) {
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
		var id ID[Int64Brand, int64]
		err := id.Scan(nil)
		if err != nil {
			t.Fatalf("Scan failed: %v", err)
		}
		if !id.IsZero() {
			t.Error("expected zero ID after scanning nil")
		}
	})

	t.Run("int64 ID from invalid type", func(t *testing.T) {
		var id ID[Int64Brand, int64]
		err := id.Scan("not-a-number")
		if err == nil {
			t.Error("expected error for invalid type")
		}
	})
}

func TestIDValue(t *testing.T) {
	t.Run("string ID non-zero", func(t *testing.T) {
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
		var id ID[Int64Brand, int64]
		val, err := id.Value()
		if err != nil {
			t.Fatalf("Value failed: %v", err)
		}
		if val != nil {
			t.Errorf("expected nil, got %v", val)
		}
	})
}

func TestIDJSONRoundTrip(t *testing.T) {
	t.Run("string ID", func(t *testing.T) {
		original := NewID[StringBrand]("test-id")
		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		var restored ID[StringBrand, string]
		err = json.Unmarshal(data, &restored)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		if restored.Get() != original.Get() {
			t.Errorf("expected %s, got %s", original.Get(), restored.Get())
		}
	})

	t.Run("zero string ID", func(t *testing.T) {
		var original ID[StringBrand, string]
		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}
		if string(data) != "null" {
			t.Errorf("expected null, got %s", string(data))
		}

		var restored ID[StringBrand, string]
		err = json.Unmarshal(data, &restored)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		if !restored.IsZero() {
			t.Error("expected zero ID")
		}
	})
}

func TestIDTextRoundTrip(t *testing.T) {
	original := NewID[StringBrand]("test-id")
	data, err := original.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText failed: %v", err)
	}

	var restored ID[StringBrand, string]
	err = restored.UnmarshalText(data)
	if err != nil {
		t.Fatalf("UnmarshalText failed: %v", err)
	}
	if restored.Get() != original.Get() {
		t.Errorf("expected %s, got %s", original.Get(), restored.Get())
	}
}

func TestIDTypeSafety(t *testing.T) {
	type UserBrand struct{}
	type OrderBrand struct{}

	userID := NewID[UserBrand]("user-123")
	orderID := NewID[OrderBrand]("order-456")

	// These should compile without issues, demonstrating type safety
	_ = userID.Get()
	_ = orderID.Get()

	// This would not compile (uncomment to verify):
	// _ = userID == orderID // compile error: mismatched types
}
