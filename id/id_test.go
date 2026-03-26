package id

import (
	"fmt"
	"math"
	"sort"
	"testing"
)

type (
	StringBrand struct{}
	Int64Brand  struct{}
	Int32Brand  struct{}
	Uint64Brand struct{}
)

func TestNewID(t *testing.T) {
	t.Parallel()
	id := NewID[StringBrand]("user-123")
	if id.Get() != "user-123" {
		t.Errorf("expected user-123, got %s", id.Get())
	}
	if id.IsZero() {
		t.Error("expected non-zero id")
	}
}

func TestNewIDInt64(t *testing.T) {
	t.Parallel()
	id := NewID[Int64Brand, int64](42)
	if id.Get() != 42 {
		t.Errorf("expected 42, got %d", id.Get())
	}
}

func TestNewIDInt32(t *testing.T) {
	t.Parallel()
	id := NewID[Int32Brand, int32](42)
	if id.Get() != 42 {
		t.Errorf("expected 42, got %d", id.Get())
	}
}

func TestNewIDUint64(t *testing.T) {
	t.Parallel()
	id := NewID[Uint64Brand, uint64](42)
	if id.Get() != 42 {
		t.Errorf("expected 42, got %d", id.Get())
	}
}

func TestIDIsZero(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	id := NewID[StringBrand]("test")
	id.Reset()
	if !id.IsZero() {
		t.Error("expected zero ID after Reset")
	}
}

func TestIDEqual(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
			t.Parallel()
			idA := NewID[Int64Brand, int](tt.a)
			idB := NewID[Int64Brand, int](tt.b)
			result, err := idA.Compare(idB)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestIDCompareString(t *testing.T) {
	t.Parallel()
	idA := NewID[StringBrand]("a")
	idB := NewID[StringBrand]("b")
	idC := NewID[StringBrand]("a")

	cmp, err := idA.Compare(idB)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmp != -1 {
		t.Error("expected 'a' < 'b'")
	}
	cmp, err = idA.Compare(idC)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmp != 0 {
		t.Error("expected 'a' == 'a'")
	}
	cmp, err = idB.Compare(idA)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmp != 1 {
		t.Error("expected 'b' > 'a'")
	}
}

func TestIDCompareInt64(t *testing.T) {
	t.Parallel()
	idA := NewID[Int64Brand, int64](100)
	idB := NewID[Int64Brand, int64](200)

	cmp, err := idA.Compare(idB)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmp != -1 {
		t.Error("expected 100 < 200")
	}
}

func TestIDCompareUint64(t *testing.T) {
	t.Parallel()
	idA := NewID[Uint64Brand, uint64](100)
	idB := NewID[Uint64Brand, uint64](200)

	cmp, err := idA.Compare(idB)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmp != -1 {
		t.Error("expected 100 < 200")
	}
}

func TestIDOr(t *testing.T) {
	t.Parallel()
	t.Run("non-zero returns self", func(t *testing.T) {
		t.Parallel()
		id := NewID[StringBrand]("test")
		defaultID := NewID[StringBrand]("default")
		result := id.Or(defaultID)
		if result.Get() != "test" {
			t.Errorf("expected test, got %s", result.Get())
		}
	})

	t.Run("zero returns default", func(t *testing.T) {
		t.Parallel()
		var id ID[StringBrand, string]
		defaultID := NewID[StringBrand]("default")
		result := id.Or(defaultID)
		if result.Get() != "default" {
			t.Errorf("expected default, got %s", result.Get())
		}
	})
}

func TestIDString(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
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
	t.Parallel()
	id := NewID[StringBrand]("test-id")
	if id.GoString() != "test-id" {
		t.Errorf("expected test-id, got %s", id.GoString())
	}
}

func TestIDFormat(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			got := fmt.Sprintf(tt.format, id)
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestIDTypeSafety(t *testing.T) {
	t.Parallel()
	type UserBrand struct{}
	type OrderBrand struct{}

	userID := NewID[UserBrand]("user-123")
	orderID := NewID[OrderBrand]("order-456")

	_ = userID.Get()
	_ = orderID.Get()
}

func TestIDSorting(t *testing.T) {
	t.Parallel()
	ids := []ID[Int64Brand, int64]{
		NewID[Int64Brand, int64](3),
		NewID[Int64Brand, int64](1),
		NewID[Int64Brand, int64](2),
	}

	sort.Slice(ids, func(i, j int) bool {
		cmp, err := ids[i].Compare(ids[j])
		if err != nil {
			panic(err)
		}
		return cmp < 0
	})

	expected := []int64{1, 2, 3}
	for i, id := range ids {
		if id.Get() != expected[i] {
			t.Errorf("position %d: expected %d, got %d", i, expected[i], id.Get())
		}
	}
}

func TestIDEdgeCases(t *testing.T) {
	t.Parallel()
	t.Run("max int64", func(t *testing.T) {
		t.Parallel()
		id := NewID[Int64Brand, int64](math.MaxInt64)
		if id.Get() != math.MaxInt64 {
			t.Errorf("expected %d, got %d", math.MaxInt64, id.Get())
		}
	})

	t.Run("min int64", func(t *testing.T) {
		t.Parallel()
		id := NewID[Int64Brand, int64](math.MinInt64)
		if id.Get() != math.MinInt64 {
			t.Errorf("expected %d, got %d", math.MinInt64, id.Get())
		}
	})

	t.Run("max uint64", func(t *testing.T) {
		t.Parallel()
		id := NewID[Uint64Brand, uint64](math.MaxUint64)
		if id.Get() != math.MaxUint64 {
			t.Errorf("expected %v, got %v", uint64(math.MaxUint64), id.Get())
		}
	})

	t.Run("empty string", func(t *testing.T) {
		t.Parallel()
		id := NewID[StringBrand]("")
		if !id.IsZero() {
			t.Error("empty string should be zero")
		}
	})
}
