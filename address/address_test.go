package address

import (
	"encoding/json/v2"
	"testing"

	"github.com/larsartmann/go-composable-business-types/validate"
)

func TestNew(t *testing.T) {
	addr := New("123 Main St", "Berlin", "10115", "de", WithLine2("Apt 4"), WithState("BE"))

	if addr.Line1 != "123 Main St" {
		t.Errorf("Line1 = %q, want %q", addr.Line1, "123 Main St")
	}

	if addr.City != "Berlin" {
		t.Errorf("City = %q, want %q", addr.City, "Berlin")
	}

	if addr.PostalCode != "10115" {
		t.Errorf("PostalCode = %q, want %q", addr.PostalCode, "10115")
	}

	if addr.CountryCode != "DE" {
		t.Errorf("CountryCode = %q, want %q", addr.CountryCode, "DE")
	}

	if addr.Line2 != "Apt 4" {
		t.Errorf("Line2 = %q, want %q", addr.Line2, "Apt 4")
	}

	if addr.State != "BE" {
		t.Errorf("State = %q, want %q", addr.State, "BE")
	}
}

func TestAddress_Validate(t *testing.T) {
	tests := []struct {
		name    string
		address *Address
		wantErr bool
	}{
		{
			name: "valid full address",
			address: New(
				"123 Main St",
				"Berlin",
				"10115",
				"DE",
				WithLine2("Apt 4"),
				WithState("BE"),
			),
			wantErr: false,
		},
		{
			name:    "valid minimal address",
			address: New("123 Main St", "Berlin", "10115", "DE"),
			wantErr: false,
		},
		{
			name:    "nil address",
			address: nil,
			wantErr: true,
		},
		{
			name:    "missing line1",
			address: New("", "Berlin", "10115", "DE"),
			wantErr: true,
		},
		{
			name:    "missing city",
			address: New("123 Main St", "", "10115", "DE"),
			wantErr: true,
		},
		{
			name:    "missing postal code",
			address: New("123 Main St", "Berlin", "", "DE"),
			wantErr: true,
		},
		{
			name:    "missing country code",
			address: New("123 Main St", "Berlin", "10115", ""),
			wantErr: true,
		},
		{
			name:    "country code too short",
			address: New("123 Main St", "Berlin", "10115", "D"),
			wantErr: true,
		},
		{
			name:    "country code too long",
			address: New("123 Main St", "Berlin", "10115", "DEU"),
			wantErr: true,
		},
		{
			name:    "country code non-alpha",
			address: New("123 Main St", "Berlin", "10115", "D3"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.address.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddress_IsZero(t *testing.T) {
	if !(*Address)(nil).IsZero() {
		t.Error("nil Address should be zero")
	}

	if !New("", "", "", "").IsZero() {
		t.Error("empty Address should be zero")
	}

	if New("123 Main St", "Berlin", "10115", "DE").IsZero() {
		t.Error("valid Address should not be zero")
	}
}

func TestAddress_Format(t *testing.T) {
	addr := New(
		"123 Main St",
		"Berlin",
		"10115",
		"DE",
		WithLine2("Apt 4"),
		WithState("BE"),
	)

	formatted := addr.Format()

	want := "123 Main St\nApt 4\nBerlin, BE 10115\nDE"
	if formatted != want {
		t.Errorf("Format() = %q, want %q", formatted, want)
	}
}

func TestAddress_JSON(t *testing.T) {
	addr := New(
		"123 Main St",
		"Berlin",
		"10115",
		"DE",
		WithLine2("Apt 4"),
		WithState("BE"),
	)

	data, err := json.Marshal(addr)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	var decoded Address
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if decoded.Line1 != addr.Line1 {
		t.Errorf("Line1 = %q, want %q", decoded.Line1, addr.Line1)
	}

	if decoded.City != addr.City {
		t.Errorf("City = %q, want %q", decoded.City, addr.City)
	}
}

func TestAddress_ValidatorInterface(t *testing.T) {
	var _ validate.Validator = (*Address)(nil)

	addr := New("123 Main St", "Berlin", "10115", "DE")
	if !validate.IsValid(addr) {
		t.Error("valid address should satisfy Validator interface")
	}
}
