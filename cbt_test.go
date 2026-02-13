package cbt

import (
	"encoding/json"
	"testing"
)

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

// Branded ID tests

type UserBrand struct{}
type OrderBrand struct{}

type UserID = ID[UserBrand, string]
type OrderID = ID[OrderBrand, int64]

func TestBrandedID_String(t *testing.T) {
	uid := NewID[UserBrand, string]("user-123")
	if uid.Value() != "user-123" {
		t.Errorf("expected user-123, got %s", uid.Value())
	}
	if uid.String() != "user-123" {
		t.Errorf("expected String() = user-123, got %s", uid.String())
	}
	if uid.GoString() != "user-123" {
		t.Errorf("expected GoString() = user-123, got %s", uid.GoString())
	}
}

func TestBrandedID_Int64(t *testing.T) {
	oid := NewID[OrderBrand, int64](42)
	if oid.Value() != 42 {
		t.Errorf("expected 42, got %d", oid.Value())
	}
	if oid.String() != "42" {
		t.Errorf("expected String() = 42, got %s", oid.String())
	}
}

func TestBrandedID_IsZero(t *testing.T) {
	uid := NewID[UserBrand, string]("user-123")
	if uid.IsZero() {
		t.Error("expected non-zero ID to return false")
	}

	var zeroUserID UserID
	if !zeroUserID.IsZero() {
		t.Error("expected zero ID to return true")
	}

	var zeroOrderID OrderID
	if !zeroOrderID.IsZero() {
		t.Error("expected zero int64 ID to return true")
	}
}

func TestBrandedID_JSON_String(t *testing.T) {
	uid := NewID[UserBrand, string]("user-abc")

	// Marshal
	data, err := json.Marshal(uid)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	if string(data) != `"user-abc"` {
		t.Errorf("expected JSON \"user-abc\", got %s", string(data))
	}

	// Unmarshal
	var parsed UserID
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if parsed.Value() != "user-abc" {
		t.Errorf("expected user-abc, got %s", parsed.Value())
	}
}

func TestBrandedID_JSON_Zero(t *testing.T) {
	var uid UserID

	// Zero value should marshal to null
	data, err := json.Marshal(uid)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	if string(data) != "null" {
		t.Errorf("expected JSON null for zero value, got %s", string(data))
	}
}

func TestBrandedID_JSON_EmptyString(t *testing.T) {
	data := []byte(`""`)
	var uid UserID
	if err := json.Unmarshal(data, &uid); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if uid.Value() != "" {
		t.Errorf("expected empty string, got %s", uid.Value())
	}
}

func TestBrandedID_BackwardsCompatibility(t *testing.T) {
	// Id[T] should still work as before
	id := NewId("legacy-id")
	if id.Value() != "legacy-id" {
		t.Errorf("expected legacy-id, got %s", id.Value())
	}

	// Id[T] should be assignable from ID[struct{}, T]
	var id2 = NewID[struct{}, string]("test")
	if id2.Value() != "test" {
		t.Errorf("expected test, got %s", id2.Value())
	}
}

func TestBrandedID_JSON_Int64_Marshal(t *testing.T) {
	// int64 IDs serialize as strings (by design)
	oid := NewID[OrderBrand, int64](42)

	data, err := json.Marshal(oid)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	// int64 serializes as string "42" (not number 42)
	if string(data) != `"42"` {
		t.Errorf("expected JSON \"42\", got %s", string(data))
	}
}

func TestBrandedID_JSON_Int64_Unmarshal_Error(t *testing.T) {
	// Unmarshaling into int64 ID should return error (documented limitation)
	data := []byte(`"42"`)
	var oid OrderID
	err := json.Unmarshal(data, &oid)
	if err == nil {
		t.Error("expected error when unmarshaling string into int64 ID")
	}
}

func TestBrandedID_Comparability(t *testing.T) {
	// IDs with same value should be equal
	uid1 := NewID[UserBrand, string]("user-123")
	uid2 := NewID[UserBrand, string]("user-123")
	if uid1 != uid2 {
		t.Error("expected equal IDs to be equal")
	}

	// IDs with different values should not be equal
	uid3 := NewID[UserBrand, string]("user-456")
	if uid1 == uid3 {
		t.Error("expected different IDs to not be equal")
	}
}

func TestBrandedID_JSON_Null_Unmarshal(t *testing.T) {
	data := []byte("null")
	var uid UserID
	if err := json.Unmarshal(data, &uid); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if !uid.IsZero() {
		t.Error("expected null to unmarshal to zero value")
	}
}
