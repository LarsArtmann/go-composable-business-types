package types

import "testing"

// testCompare runs table-driven comparison tests for types with Compare method.
func testCompare[T any](t *testing.T, tests []struct {
	a, b     T
	expected int
}, compare func(a, b T) int,
) {
	t.Helper()

	for _, tt := range tests {
		result := compare(tt.a, tt.b)
		if result != tt.expected {
			t.Errorf("Compare(%v, %v) = %d, expected %d", tt.a, tt.b, result, tt.expected)
		}
	}
}

func TestNewPercentage(t *testing.T) {
	t.Parallel()

	p := NewPercentage(50)
	if p.Float64() != 0.5 {
		t.Errorf("expected 0.5, got %f", p.Float64())
	}

	clamped := NewPercentage(150)
	if clamped != 100 {
		t.Errorf("expected 100, got %d", clamped)
	}
}

func TestPercentageHelpers(t *testing.T) {
	t.Parallel()

	zero := NewPercentage(0)
	half := NewPercentage(50)
	full := NewPercentage(100)

	if !zero.IsZero() || !zero.IsMin() || zero.IsMax() {
		t.Error("percentage 0 helpers failed")
	}

	if half.IsZero() || half.IsMin() || half.IsMax() {
		t.Error("percentage 50 helpers failed")
	}

	if full.IsZero() || full.IsMin() || !full.IsMax() {
		t.Error("percentage 100 helpers failed")
	}
}

// testNumericCompare runs comparison tests for numeric types.
func testNumericCompare[T any](
	t *testing.T,
	name string,
	newVal func(int64) T,
	compare func(T, T) int,
) {
	t.Run(name, func(t *testing.T) {
		t.Parallel()
		testCompare(t, []struct {
			a, b     T
			expected int
		}{
			{newVal(50), newVal(50), 0},
			{newVal(30), newVal(50), -1},
			{newVal(70), newVal(50), 1},
		}, compare)
	})
}

func TestPercentageCompare(t *testing.T) {
	testNumericCompare(
		t,
		"Percentage",
		func(v int64) Percentage { return NewPercentage(uint8(v)) },
		Percentage.Compare,
	)
}

func TestPercentageJSON(t *testing.T) {
	t.Parallel()
	// Test MarshalJSON
	p := NewPercentage(50)

	data, err := p.MarshalJSON()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if string(data) != "50" {
		t.Errorf("expected 50, got %s", string(data))
	}

	// Test UnmarshalJSON
	var p2 Percentage
	if err := p2.UnmarshalJSON([]byte("75")); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if p2 != 75 {
		t.Errorf("expected 75, got %d", p2)
	}

	// Test Round-trip
	var p3 Percentage
	if err := p3.UnmarshalJSON(data); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if p3 != p {
		t.Errorf("round-trip failed: expected %d, got %d", p, p3)
	}
}

func TestPercentageSQL(t *testing.T) {
	t.Parallel()
	// Test Value
	p := NewPercentage(75)

	val, err := p.Value()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if val != int64(75) {
		t.Errorf("expected 75, got %v", val)
	}

	// Test Scan with int64
	var p2 Percentage
	if err := p2.Scan(int64(50)); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if p2 != 50 {
		t.Errorf("expected 50, got %d", p2)
	}
}

func TestNewCents(t *testing.T) {
	t.Parallel()

	c := NewCents(1099)
	if c.Int64() != 1099 {
		t.Errorf("expected 1099, got %d", c.Int64())
	}

	if c.Float64() != 10.99 {
		t.Errorf("expected 10.99, got %f", c.Float64())
	}
}

func TestCentsMath(t *testing.T) {
	t.Parallel()

	a := NewCents(100)
	b := NewCents(50)

	if a.Add(b) != 150 {
		t.Error("Add failed")
	}

	if a.Sub(b) != 50 {
		t.Error("Sub failed")
	}

	if a.Mul(2) != 200 {
		t.Error("Mul failed")
	}

	if a.Div(2) != 50 {
		t.Error("Div failed")
	}

	// Abs
	neg := NewCents(-100)
	if neg.Abs() != 100 {
		t.Error("Abs failed")
	}

	// Sign
	if NewCents(-100).Sign() != -1 {
		t.Error("Sign negative failed")
	}

	if NewCents(0).Sign() != 0 {
		t.Error("Sign zero failed")
	}

	if NewCents(100).Sign() != 1 {
		t.Error("Sign positive failed")
	}
}

func TestCentsCompare(t *testing.T) {
	testNumericCompare(t, "Cents", NewCents, Cents.Compare)
}

func TestPercentageString(t *testing.T) {
	t.Parallel()

	p := NewPercentage(50)
	if p.String() != "50%" {
		t.Errorf("expected 50%%, got %s", p.String())
	}

	zero := NewPercentage(0)
	if zero.String() != "0%" {
		t.Errorf("expected 0%%, got %s", zero.String())
	}
}

func TestPercentageValidate(t *testing.T) {
	t.Parallel()

	p := NewPercentage(50)
	err := p.Validate()
	if err != nil {
		t.Errorf("valid percentage should validate: %v", err)
	}
}

func TestCentsHelpers(t *testing.T) {
	t.Parallel()

	positive := NewCents(100)
	negative := NewCents(-100)
	zero := NewCents(0)

	if !positive.IsPositive() {
		t.Error("100 cents should be positive")
	}

	if positive.IsNegative() {
		t.Error("100 cents should not be negative")
	}

	if !positive.IsZero() == false {
		t.Error("100 cents should not be zero")
	}

	if negative.IsPositive() {
		t.Error("-100 cents should not be positive")
	}

	if !negative.IsNegative() {
		t.Error("-100 cents should be negative")
	}

	if !zero.IsZero() {
		t.Error("0 cents should be zero")
	}

	if zero.IsPositive() {
		t.Error("0 cents should not be positive")
	}

	if zero.IsNegative() {
		t.Error("0 cents should not be negative")
	}
}

func TestCentsString(t *testing.T) {
	t.Parallel()

	c := NewCents(1234)
	if c.String() != "$12.34" {
		t.Errorf("expected $12.34, got %s", c.String())
	}

	negative := NewCents(-500)
	if negative.String() != "$-5.00" {
		t.Errorf("expected $-5.00, got %s", negative.String())
	}
}

func TestCentsValidate(t *testing.T) {
	t.Parallel()

	c := NewCents(100)
	err := c.Validate()
	if err != nil {
		t.Errorf("valid cents should validate: %v", err)
	}
}
