package cbt

import "testing"

func TestId(t *testing.T) {
	id := NewId("user-123")
	if id.Value() != "user-123" {
		t.Errorf("expected user-123, got %s", id.Value())
	}
	if id.IsZero() {
		t.Error("expected non-zero id")
	}
}

func TestIdInt(t *testing.T) {
	id := NewId(42)
	if id.Value() != 42 {
		t.Errorf("expected 42, got %d", id.Value())
	}
}

func TestActor(t *testing.T) {
	chain := NewActorChain(UserActor(NewId("u-1"), "Alice")).
		Append(ServiceActor(NewId("api-gateway"), "API Gateway")).
		Append(ServiceActor(NewId("order-svc"), "Order Service"))

	if chain.IsEmpty() {
		t.Error("expected non-empty chain")
	}
	if chain.Origin().Kind != ActorKindUser {
		t.Error("expected origin to be user")
	}
	if chain.Origin().Name != "Alice" {
		t.Error("expected origin name to be Alice")
	}
	if chain.Current().Kind != ActorKindService {
		t.Error("expected current to be service")
	}
	if !chain.HasKind(ActorKindUser) {
		t.Error("expected chain to have user")
	}
	if len(chain.ByKind(ActorKindService)) != 2 {
		t.Error("expected 2 services in chain")
	}
}

func TestPercentage(t *testing.T) {
	p := NewPercentage(50)
	if p.Float64() != 0.5 {
		t.Errorf("expected 0.5, got %f", p.Float64())
	}

	clamped := NewPercentage(150)
	if clamped != 100 {
		t.Errorf("expected 100, got %d", clamped)
	}
}

func TestCents(t *testing.T) {
	c := NewCents(1099)
	if c.Int64() != 1099 {
		t.Errorf("expected 1099, got %d", c.Int64())
	}
	if c.Float64() != 10.99 {
		t.Errorf("expected 10.99, got %f", c.Float64())
	}
}

func TestBoundedString(t *testing.T) {
	bs, err := NewBoundedString(1, 10, "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bs.String() != "hello" {
		t.Errorf("expected hello, got %s", bs.String())
	}
	if bs.Len() != 5 {
		t.Errorf("expected length 5, got %d", bs.Len())
	}
	if bs.MinLen() != 1 || bs.MaxLen() != 10 {
		t.Errorf("expected bounds [1,10], got [%d,%d]", bs.MinLen(), bs.MaxLen())
	}
}

func TestBoundedStringValidation(t *testing.T) {
	_, err := NewBoundedString(1, 10, "")
	if err == nil {
		t.Error("expected error for empty string with min 1")
	}

	_, err = NewBoundedString(1, 5, "too long string")
	if err == nil {
		t.Error("expected error for string exceeding max")
	}

	_, err = NewBoundedString(5, 3, "test")
	if err == nil {
		t.Error("expected error for invalid bounds")
	}
}

func TestBoundedStringHelperMethods(t *testing.T) {
	bs := MustBoundedString(2, 10, "test")
	if bs.IsEmpty() {
		t.Error("expected non-empty")
	}

	short := MustBoundedString(2, 10, "ab")
	if !short.IsMinLength() {
		t.Error("expected to be exactly at min length")
	}

	maxed := MustBoundedString(2, 4, "abcd")
	if !maxed.IsMaxLength() {
		t.Error("expected to be at max length")
	}
}

func TestBoundedStringOf(t *testing.T) {
	NewName := BoundedStringOf(1, 50)
	name, err := NewName("John")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name.String() != "John" {
		t.Errorf("expected John, got %s", name.String())
	}
}

func TestTrimmedBoundedString(t *testing.T) {
	bs, err := TrimmedBoundedString(1, 10, "  hello  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bs.String() != "hello" {
		t.Errorf("expected trimmed 'hello', got '%s'", bs.String())
	}
}

func TestNonEmptyString(t *testing.T) {
	_, err := NonEmptyString(10, "")
	if err == nil {
		t.Error("expected error for empty string")
	}

	bs, err := NonEmptyString(10, "x")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bs.String() != "x" {
		t.Errorf("expected x, got %s", bs.String())
	}
}
