package contact

import (
	"encoding/json"
	"testing"

	"github.com/larsartmann/go-composable-business-types/address"
	"github.com/larsartmann/go-composable-business-types/validate"
)

func TestNew(t *testing.T) {
	addr := address.New("Main Str. 1", "Berlin", "10115", "DE")
	c := New(
		"Jane Doe ",
		WithEmail("jane@example.com"),
		WithPhone("+49 30 123456"),
		WithWebsite("https://example.com"),
		WithAddress(addr),
	)

	if c.Name != "Jane Doe" {
		t.Errorf("Name = %q, want %q", c.Name, "Jane Doe")
	}

	if c.Email != "jane@example.com" {
		t.Errorf("Email = %q, want %q", c.Email, "jane@example.com")
	}

	if c.Phone != "+49 30 123456" {
		t.Errorf("Phone = %q, want %q", c.Phone, "+49 30 123456")
	}

	if c.Website != "https://example.com" {
		t.Errorf("Website = %q, want %q", c.Website, "https://example.com")
	}

	if c.Address != addr {
		t.Error("Address pointer mismatch")
	}
}

func TestContact_Validate(t *testing.T) {
	tests := []struct {
		name    string
		contact *Contact
		wantErr bool
	}{
		{
			name:    "valid minimal contact",
			contact: New("Jane Doe"),
			wantErr: false,
		},
		{
			name: "valid full contact",
			contact: New(
				"Jane Doe",
				WithEmail("jane@example.com"),
				WithPhone("+49 30 123456"),
				WithWebsite("https://example.com"),
				WithAddress(address.New("Main Str. 1", "Berlin", "10115", "DE")),
			),
			wantErr: false,
		},
		{
			name:    "nil contact",
			contact: nil,
			wantErr: true,
		},
		{
			name:    "missing name",
			contact: New(""),
			wantErr: true,
		},
		{
			name:    "invalid email",
			contact: New("Jane Doe", WithEmail("not-an-email")),
			wantErr: true,
		},
		{
			name:    "invalid website scheme",
			contact: New("Jane Doe", WithWebsite("ftp://example.com")),
			wantErr: true,
		},
		{
			name:    "invalid address",
			contact: New("Jane Doe", WithAddress(address.New("", "", "", ""))),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.contact.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContact_IsZero(t *testing.T) {
	if !(*Contact)(nil).IsZero() {
		t.Error("nil Contact should be zero")
	}

	if !New("").IsZero() {
		t.Error("empty Contact should be zero")
	}

	if New("Jane Doe").IsZero() {
		t.Error("Contact with name should not be zero")
	}
}

func TestContact_JSON(t *testing.T) {
	c := New(
		"Jane Doe",
		WithEmail("jane@example.com"),
		WithAddress(address.New("Main Str. 1", "Berlin", "10115", "DE")),
	)

	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	var decoded Contact
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if decoded.Name != c.Name {
		t.Errorf("Name = %q, want %q", decoded.Name, c.Name)
	}

	if decoded.Email != c.Email {
		t.Errorf("Email = %q, want %q", decoded.Email, c.Email)
	}
}

func TestContact_ValidatorInterface(t *testing.T) {
	var _ validate.Validator = (*Contact)(nil)

	c := New("Jane Doe", WithEmail("jane@example.com"))
	if !validate.IsValid(c) {
		t.Error("valid contact should satisfy Validator interface")
	}
}
