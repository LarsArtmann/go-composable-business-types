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

func TestNewIDNumeric(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		brand    any
		value    any
		expected any
	}{
		{"int64", Int64Brand{}, int64(42), int64(42)},
		{"int32", Int32Brand{}, int32(42), int32(42)},
		{"uint64", Uint64Brand{}, uint64(42), uint64(42)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch v := tt.value.(type) {
			case int64:
				id := NewID[Int64Brand, int64](v)
				if id.Get() != tt.expected.(int64) {
					t.Errorf("expected %d, got %d", tt.expected, id.Get())
				}
			case int32:
				id := NewID[Int32Brand, int32](v)
				if id.Get() != tt.expected.(int32) {
					t.Errorf("expected %d, got %d", tt.expected, id.Get())
				}
			case uint64:
				id := NewID[Uint64Brand, uint64](v)
				if id.Get() != tt.expected.(uint64) {
					t.Errorf("expected %d, got %d", tt.expected, id.Get())
				}
			}
		})
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

	tests := []struct {
		name     string
		a, b     int64
		expected int
	}{
		{"less", 100, 200, -1},
		{"equal", 100, 100, 0},
		{"greater", 200, 100, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			idA := NewID[Int64Brand, int64](tt.a)
			idB := NewID[Int64Brand, int64](tt.b)

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

func TestIDCompareUint64(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		a, b     uint64
		expected int
	}{
		{"less", 100, 200, -1},
		{"equal", 100, 100, 0},
		{"greater", 200, 100, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			idA := NewID[Uint64Brand, uint64](tt.a)
			idB := NewID[Uint64Brand, uint64](tt.b)

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

	tests := []struct {
		name     string
		brand    func(v any) any
		value    any
		expected any
	}{
		{
			"max int64",
			func(v any) any { return NewID[Int64Brand, int64](v.(int64)) },
			int64(math.MaxInt64),
			int64(math.MaxInt64),
		},
		{
			"min int64",
			func(v any) any { return NewID[Int64Brand, int64](v.(int64)) },
			int64(math.MinInt64),
			int64(math.MinInt64),
		},
		{
			"max uint64",
			func(v any) any { return NewID[Uint64Brand, uint64](v.(uint64)) },
			uint64(math.MaxUint64),
			uint64(math.MaxUint64),
		},
		{"empty string", func(v any) any { return NewID[StringBrand](v.(string)) }, "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id := tt.brand(tt.value)
			switch v := id.(type) {
			case ID[Int64Brand, int64]:
				if v.Get() != tt.expected.(int64) {
					t.Errorf("expected %d, got %d", tt.expected, v.Get())
				}
			case ID[Uint64Brand, uint64]:
				if v.Get() != tt.expected.(uint64) {
					t.Errorf("expected %v, got %v", tt.expected, v.Get())
				}
			case ID[StringBrand, string]:
				if v.Get() != tt.expected.(string) {
					t.Errorf("expected %s, got %s", tt.expected, v.Get())
				}

				if !v.IsZero() {
					t.Error("empty string should be zero")
				}
			}
		})
	}
}
