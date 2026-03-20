package id

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"testing"
	"unicode/utf8"
)

type (
	StringBrand struct{}
	Int64Brand  struct{}
	Int32Brand  struct{}
	Uint64Brand struct{}
)

// Tests for basic functionality

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

func TestNewIDInt32(t *testing.T) {
	id := NewID[Int32Brand, int32](42)
	if id.Get() != 42 {
		t.Errorf("expected 42, got %d", id.Get())
	}
}

func TestNewIDUint64(t *testing.T) {
	id := NewID[Uint64Brand, uint64](42)
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

func TestIDReset(t *testing.T) {
	id := NewID[StringBrand]("test")
	id.Reset()
	if !id.IsZero() {
		t.Error("expected zero ID after Reset")
	}
}

func TestIDEqual(t *testing.T) {
	id1 := NewID[StringBrand]("test")
	id2 := NewID[StringBrand]("test")
	id3 := NewID[StringBrand]("other")

	if !id1.Equal(id2) {
		t.Error("expected equal IDs")
	}
	if id1.Equal(id3) {
		t.Error("expected unequal IDs")
	}
}

func TestIDCompare(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"less", 1, 2, -1},
		{"equal", 5, 5, 0},
		{"greater", 3, 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idA := NewID[Int64Brand, int](tt.a)
			idB := NewID[Int64Brand, int](tt.b)
			result := idA.Compare(idB)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestIDCompareString(t *testing.T) {
	idA := NewID[StringBrand]("a")
	idB := NewID[StringBrand]("b")
	idC := NewID[StringBrand]("a")

	if idA.Compare(idB) != -1 {
		t.Error("expected 'a' < 'b'")
	}
	if idA.Compare(idC) != 0 {
		t.Error("expected 'a' == 'a'")
	}
	if idB.Compare(idA) != 1 {
		t.Error("expected 'b' > 'a'")
	}
}

func TestIDCompareInt64(t *testing.T) {
	idA := NewID[Int64Brand, int64](100)
	idB := NewID[Int64Brand, int64](200)

	if idA.Compare(idB) != -1 {
		t.Error("expected 100 < 200")
	}
}

func TestIDCompareUint64(t *testing.T) {
	idA := NewID[Uint64Brand, uint64](100)
	idB := NewID[Uint64Brand, uint64](200)

	if idA.Compare(idB) != -1 {
		t.Error("expected 100 < 200")
	}
}

func TestIDOr(t *testing.T) {
	t.Run("non-zero returns self", func(t *testing.T) {
		id := NewID[StringBrand]("test")
		defaultID := NewID[StringBrand]("default")
		result := id.Or(defaultID)
		if result.Get() != "test" {
			t.Errorf("expected test, got %s", result.Get())
		}
	})

	t.Run("zero returns default", func(t *testing.T) {
		var id ID[StringBrand, string]
		defaultID := NewID[StringBrand]("default")
		result := id.Or(defaultID)
		if result.Get() != "default" {
			t.Errorf("expected default, got %s", result.Get())
		}
	})
}

func TestIDString(t *testing.T) {
	tests := []struct {
		name     string
		id       any
		expected string
	}{
		{"string", NewID[StringBrand]("test-id"), "test-id"},
		{"int64", NewID[Int64Brand, int64](42), "42"},
		{"int32", NewID[Int32Brand, int32](42), "42"},
		{"uint64", NewID[Uint64Brand, uint64](42), "42"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string
			switch v := tt.id.(type) {
			case ID[StringBrand, string]:
				got = v.String()
			case ID[Int64Brand, int64]:
				got = v.String()
			case ID[Int32Brand, int32]:
				got = v.String()
			case ID[Uint64Brand, uint64]:
				got = v.String()
			}
			if got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}

func TestIDGoString(t *testing.T) {
	id := NewID[StringBrand]("test-id")
	if id.GoString() != "test-id" {
		t.Errorf("expected test-id, got %s", id.GoString())
	}
}

func TestIDFormat(t *testing.T) {
	id := NewID[Int64Brand, int64](42)

	tests := []struct {
		format   string
		expected string
	}{
		{"%s", "42"},
		{"%d", "42"},
		{"%q", `"42"`},
		{"%v", "42"},
		{"%#v", "id(42)"},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			got := fmt.Sprintf(tt.format, id)
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

// JSON Tests

func TestIDJSON(t *testing.T) {
	t.Run("string ID non-zero", func(t *testing.T) {
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
	t.Run("string ID from string", func(t *testing.T) {
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
		var id ID[StringBrand, string]
		err := json.Unmarshal([]byte(`invalid`), &id)
		if err == nil {
			t.Error("expected error for invalid JSON")
		}
	})

	t.Run("number into string ID", func(t *testing.T) {
		var id ID[StringBrand, string]
		err := json.Unmarshal([]byte("123"), &id)
		if err == nil {
			t.Error("expected error for number into string ID")
		}
	})

	t.Run("string into int64 ID", func(t *testing.T) {
		var id ID[Int64Brand, int64]
		err := json.Unmarshal([]byte(`"not-a-number"`), &id)
		if err == nil {
			t.Error("expected error for string into int64 ID")
		}
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
	t.Run("string ID", func(t *testing.T) {
		testJSONRoundTrip[StringBrand, string](t, "test-id")
	})

	t.Run("int64 ID", func(t *testing.T) {
		testJSONRoundTrip[Int64Brand, int64](t, 42)
	})

	t.Run("int32 ID", func(t *testing.T) {
		testJSONRoundTrip[Int32Brand, int32](t, 42)
	})

	t.Run("uint64 ID", func(t *testing.T) {
		testJSONRoundTrip[Uint64Brand, uint64](t, 42)
	})
}

// Text encoding tests

func TestIDText(t *testing.T) {
	t.Run("marshal non-zero string", func(t *testing.T) {
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
		testUnmarshalTextRoundTrip[StringBrand, string](t, "test-id", "test-id")
	})

	t.Run("unmarshal empty", func(t *testing.T) {
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
		testUnmarshalTextRoundTrip[Int64Brand, int64](t, "42", 42)
	})

	t.Run("unmarshal uint64", func(t *testing.T) {
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

// Binary encoding tests

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
	t.Run("string ID", func(t *testing.T) {
		testBinaryRoundTrip[StringBrand, string](t, "test-id")
	})

	t.Run("int64 ID", func(t *testing.T) {
		testBinaryRoundTrip[Int64Brand, int64](t, 42)
	})

	t.Run("int32 ID", func(t *testing.T) {
		testBinaryRoundTrip[Int32Brand, int32](t, 42)
	})

	t.Run("uint64 ID", func(t *testing.T) {
		testBinaryRoundTrip[Uint64Brand, uint64](t, 42)
	})

	t.Run("zero ID", func(t *testing.T) {
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

// Gob encoding tests

func TestIDGob(t *testing.T) {
	t.Run("string ID", func(t *testing.T) {
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

// SQL Scan/Value tests

func TestIDScan(t *testing.T) {
	t.Run("string ID from string", func(t *testing.T) {
		testScanRoundTrip[StringBrand, string](t, "test-id", "test-id")
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
			t.Error("expected zero ID")
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
		testScanRoundTrip[Int64Brand, int64](t, int64(42), int64(42))
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
			t.Error("expected zero ID")
		}
	})

	t.Run("int64 ID from invalid type", func(t *testing.T) {
		var id ID[Int64Brand, int64]
		err := id.Scan("not-a-number")
		if err == nil {
			t.Error("expected error for invalid type")
		}
	})

	t.Run("int32 ID from int64", func(t *testing.T) {
		testScanRoundTrip[Int32Brand, int32](t, int64(42), int32(42))
	})

	t.Run("uint64 ID from int64", func(t *testing.T) {
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

	t.Run("int32 ID non-zero", func(t *testing.T) {
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

// Type safety test

func TestIDTypeSafety(t *testing.T) {
	type UserBrand struct{}
	type OrderBrand struct{}

	userID := NewID[UserBrand]("user-123")
	orderID := NewID[OrderBrand]("order-456")

	_ = userID.Get()
	_ = orderID.Get()
}

// Sorting test

func TestIDSorting(t *testing.T) {
	ids := []ID[Int64Brand, int64]{
		NewID[Int64Brand, int64](3),
		NewID[Int64Brand, int64](1),
		NewID[Int64Brand, int64](2),
	}

	sort.Slice(ids, func(i, j int) bool {
		return ids[i].Compare(ids[j]) < 0
	})

	expected := []int64{1, 2, 3}
	for i, id := range ids {
		if id.Get() != expected[i] {
			t.Errorf("position %d: expected %d, got %d", i, expected[i], id.Get())
		}
	}
}

// Edge case tests

func TestIDEdgeCases(t *testing.T) {
	t.Run("max int64", func(t *testing.T) {
		id := NewID[Int64Brand, int64](math.MaxInt64)
		if id.Get() != math.MaxInt64 {
			t.Errorf("expected %d, got %d", math.MaxInt64, id.Get())
		}
	})

	t.Run("min int64", func(t *testing.T) {
		id := NewID[Int64Brand, int64](math.MinInt64)
		if id.Get() != math.MinInt64 {
			t.Errorf("expected %d, got %d", math.MinInt64, id.Get())
		}
	})

	t.Run("max uint64", func(t *testing.T) {
		id := NewID[Uint64Brand, uint64](math.MaxUint64)
		if id.Get() != math.MaxUint64 {
			t.Errorf("expected %v, got %v", uint64(math.MaxUint64), id.Get())
		}
	})

	t.Run("empty string", func(t *testing.T) {
		id := NewID[StringBrand]("")
		if !id.IsZero() {
			t.Error("empty string should be zero")
		}
	})
}

// Fuzz tests

func FuzzIDJSONString(f *testing.F) {
	testcases := []string{"test", "hello-world", "123", "", "unicode-日本語"}
	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, orig string) {
		// Skip invalid UTF-8 strings - JSON requires valid UTF-8 and cannot
		// round-trip invalid sequences (they become replacement characters)
		if !utf8.ValidString(orig) {
			t.Skip("skipping invalid UTF-8 string")
		}

		id := NewID[StringBrand](orig)
		data, err := json.Marshal(id)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		var restored ID[StringBrand, string]
		err = json.Unmarshal(data, &restored)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if restored.Get() != orig {
			t.Errorf("expected %q, got %q", orig, restored.Get())
		}
	})
}

func FuzzIDJSONInt64(f *testing.F) {
	testcases := []int64{0, 1, -1, 42, math.MaxInt64, math.MinInt64}
	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, orig int64) {
		id := NewID[Int64Brand, int64](orig)
		data, err := json.Marshal(id)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		var restored ID[Int64Brand, int64]
		err = json.Unmarshal(data, &restored)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if restored.Get() != orig {
			t.Errorf("expected %d, got %d", orig, restored.Get())
		}
	})
}

func FuzzIDBinaryString(f *testing.F) {
	testcases := []string{"test", "hello-world", "123", "unicode-日本語"}
	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, orig string) {
		if orig == "" {
			return // empty string is zero value
		}
		id := NewID[StringBrand](orig)
		data, err := id.MarshalBinary()
		if err != nil {
			t.Fatalf("MarshalBinary failed: %v", err)
		}

		var restored ID[StringBrand, string]
		err = restored.UnmarshalBinary(data)
		if err != nil {
			t.Fatalf("UnmarshalBinary failed: %v", err)
		}

		if restored.Get() != orig {
			t.Errorf("expected %q, got %q", orig, restored.Get())
		}
	})
}

// Benchmarks

func BenchmarkNewID(b *testing.B) {
	for b.Loop() {
		_ = NewID[StringBrand]("test-id")
	}
}

func BenchmarkIDGet(b *testing.B) {
	id := NewID[StringBrand]("test-id")
	for b.Loop() {
		_ = id.Get()
	}
}

func BenchmarkIDString(b *testing.B) {
	id := NewID[StringBrand]("test-id")
	for b.Loop() {
		_ = id.String()
	}
}

func BenchmarkIDStringInt64(b *testing.B) {
	id := NewID[Int64Brand, int64](123456789)
	for b.Loop() {
		_ = id.String()
	}
}

func BenchmarkIDIsZero(b *testing.B) {
	id := NewID[StringBrand]("test-id")
	for b.Loop() {
		_ = id.IsZero()
	}
}

func BenchmarkIDEqual(b *testing.B) {
	id1 := NewID[StringBrand]("test-id")
	id2 := NewID[StringBrand]("test-id")
	for b.Loop() {
		_ = id1.Equal(id2)
	}
}

func BenchmarkIDCompare(b *testing.B) {
	id1 := NewID[Int64Brand, int64](100)
	id2 := NewID[Int64Brand, int64](200)
	for b.Loop() {
		_ = id1.Compare(id2)
	}
}

func BenchmarkIDMarshalJSON(b *testing.B) {
	id := NewID[StringBrand]("test-id-12345")
	for b.Loop() {
		_, _ = id.MarshalJSON()
	}
}

func BenchmarkIDMarshalJSONInt64(b *testing.B) {
	id := NewID[Int64Brand, int64](123456789)
	for b.Loop() {
		_, _ = id.MarshalJSON()
	}
}

func BenchmarkIDUnmarshalJSON(b *testing.B) {
	data := []byte(`"test-id-12345"`)
	for b.Loop() {
		var id ID[StringBrand, string]
		_ = id.UnmarshalJSON(data)
	}
}

func BenchmarkIDUnmarshalJSONInt64(b *testing.B) {
	data := []byte(`123456789`)
	for b.Loop() {
		var id ID[Int64Brand, int64]
		_ = id.UnmarshalJSON(data)
	}
}

func BenchmarkIDMarshalBinary(b *testing.B) {
	id := NewID[Int64Brand, int64](123456789)
	for b.Loop() {
		_, _ = id.MarshalBinary()
	}
}

func BenchmarkIDUnmarshalBinary(b *testing.B) {
	id := NewID[Int64Brand, int64](123456789)
	data, _ := id.MarshalBinary()
	for b.Loop() {
		var restored ID[Int64Brand, int64]
		_ = restored.UnmarshalBinary(data)
	}
}

func BenchmarkIDScan(b *testing.B) {
	for b.Loop() {
		var id ID[StringBrand, string]
		_ = id.Scan("test-id-12345")
	}
}

func BenchmarkIDScanInt64(b *testing.B) {
	for b.Loop() {
		var id ID[Int64Brand, int64]
		_ = id.Scan(int64(123456789))
	}
}

func BenchmarkIDValue(b *testing.B) {
	id := NewID[StringBrand]("test-id-12345")
	for b.Loop() {
		_, _ = id.Value()
	}
}

func BenchmarkIDValueInt64(b *testing.B) {
	id := NewID[Int64Brand, int64](123456789)
	for b.Loop() {
		_, _ = id.Value()
	}
}

func BenchmarkJSONRoundTrip(b *testing.B) {
	id := NewID[StringBrand]("test-id-12345")
	for b.Loop() {
		data, _ := json.Marshal(id)
		var restored ID[StringBrand, string]
		_ = json.Unmarshal(data, &restored)
	}
}

func BenchmarkJSONRoundTripInt64(b *testing.B) {
	id := NewID[Int64Brand, int64](123456789)
	for b.Loop() {
		data, _ := json.Marshal(id)
		var restored ID[Int64Brand, int64]
		_ = json.Unmarshal(data, &restored)
	}
}

// Example functions for godoc

func ExampleNewID() {
	type UserBrand struct{}

	userID := NewID[UserBrand]("user-123")
	fmt.Println(userID.Get())
	// Output: user-123
}

func ExampleID_String() {
	type OrderBrand struct{}

	orderID := NewID[OrderBrand, int64](12345)
	fmt.Println(orderID.String())
	// Output: 12345
}

func ExampleID_Equal() {
	type UserBrand struct{}

	id1 := NewID[UserBrand]("user-123")
	id2 := NewID[UserBrand]("user-123")
	id3 := NewID[UserBrand]("user-456")

	fmt.Println(id1.Equal(id2))
	fmt.Println(id1.Equal(id3))
	// Output:
	// true
	// false
}

func ExampleID_Compare() {
	type OrderBrand struct{}

	id1 := NewID[OrderBrand, int64](100)
	id2 := NewID[OrderBrand, int64](200)
	id3 := NewID[OrderBrand, int64](100)

	fmt.Println(id1.Compare(id2))
	fmt.Println(id2.Compare(id1))
	fmt.Println(id1.Compare(id3))
	// Output:
	// -1
	// 1
	// 0
}

func ExampleID_Or() {
	type UserBrand struct{}

	id := NewID[UserBrand]("user-123")
	defaultID := NewID[UserBrand]("unknown")

	fmt.Println(id.Or(defaultID).Get())

	var zeroID ID[UserBrand, string]
	fmt.Println(zeroID.Or(defaultID).Get())
	// Output:
	// user-123
	// unknown
}

func ExampleID_IsZero() {
	type UserBrand struct{}

	id := NewID[UserBrand]("user-123")
	var zeroID ID[UserBrand, string]

	fmt.Println(id.IsZero())
	fmt.Println(zeroID.IsZero())
	// Output:
	// false
	// true
}

func ExampleID_Reset() {
	type UserBrand struct{}

	id := NewID[UserBrand]("user-123")
	id.Reset()
	fmt.Println(id.IsZero())
	// Output: true
}
