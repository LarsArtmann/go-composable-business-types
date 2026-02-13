//go:build goexperiment.jsonv2

package cbt

import (
	"testing"

	jsonv2 "encoding/json/v2"
)

func TestBrandedID_JSONv2_Marshal(t *testing.T) {
	uid := NewID[UserBrand, string]("user-v2")

	data, err := jsonv2.Marshal(uid)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	if string(data) != `"user-v2"` {
		t.Errorf("expected JSON \"user-v2\", got %s", string(data))
	}
}

func TestBrandedID_JSONv2_Unmarshal(t *testing.T) {
	data := []byte(`"user-v2-unmarshal"`)

	var parsed UserID
	if err := jsonv2.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if parsed.Value() != "user-v2-unmarshal" {
		t.Errorf("expected user-v2-unmarshal, got %s", parsed.Value())
	}
}

func TestBrandedID_JSONv2_Zero(t *testing.T) {
	var uid UserID

	data, err := jsonv2.Marshal(uid)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	if string(data) != "null" {
		t.Errorf("expected JSON null for zero value, got %s", string(data))
	}
}

func TestBrandedID_JSONv2_Null_Unmarshal(t *testing.T) {
	data := []byte("null")

	var uid UserID
	if err := jsonv2.Unmarshal(data, &uid); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if !uid.IsZero() {
		t.Error("expected null to unmarshal to zero value")
	}
}

func TestBrandedID_JSONv2_EmptyString(t *testing.T) {
	data := []byte(`""`)

	var uid UserID
	if err := jsonv2.Unmarshal(data, &uid); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if uid.Value() != "" {
		t.Errorf("expected empty string, got %s", uid.Value())
	}
}

func TestBrandedID_JSONv2_Struct(t *testing.T) {
	type User struct {
		ID   UserID `json:"id"`
		Name string `json:"name"`
	}

	user := User{ID: NewID[UserBrand, string]("user-struct"), Name: "Alice"}

	data, err := jsonv2.Marshal(user)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	expected := `{"id":"user-struct","name":"Alice"}`
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}

	var parsed User
	if err := jsonv2.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if parsed.ID.Value() != "user-struct" {
		t.Errorf("expected user-struct, got %s", parsed.ID.Value())
	}
}

func TestBrandedID_JSONv2_Struct_ZeroID(t *testing.T) {
	type User struct {
		ID   UserID `json:"id"`
		Name string `json:"name"`
	}

	user := User{Name: "Anonymous"} // ID is zero

	data, err := jsonv2.Marshal(user)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	expected := `{"id":null,"name":"Anonymous"}`
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}
}
