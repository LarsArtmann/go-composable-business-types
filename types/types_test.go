package types

import (
	"errors"
	"testing"

	pkgerrors "github.com/larsartmann/go-composable-business-types/pkg/errors"
)

// validationCase represents a test case for string validation.
type validationCase struct {
	name    string
	input   string
	wantErr bool
}

// testValidation runs table-driven validation tests for types with String() method.
func testValidation[T interface{ String() string }](t *testing.T, tests []validationCase, constructor func(string) (T, error)) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := constructor(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if got.String() != tt.input {
				t.Errorf("expected %s, got %s", tt.input, got.String())
			}
		})
	}
}

func TestEmailParts(t *testing.T) {
	t.Parallel()
	email, _ := NewEmail("user@example.com")
	if email.LocalPart() != "user" {
		t.Errorf("expected local part 'user', got %s", email.LocalPart())
	}
	if email.Domain() != "example.com" {
		t.Errorf("expected domain 'example.com', got %s", email.Domain())
	}
}

func TestValidation(t *testing.T) {
	t.Parallel()

	// Email validation tests - inline table
	for _, tt := range []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid simple", "test@example.com", false},
		{"valid with dots", "first.last@example.com", false},
		{"valid with plus", "user+tag@example.com", false},
		{"empty", "", true},
		{"no at", "testexample.com", true},
		{"no domain", "test@", true},
		{"no local", "@example.com", true},
	} {
		t.Run("Email/"+tt.name, func(t *testing.T) {
			t.Parallel()
			email, err := NewEmail(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if email.String() != tt.input {
				t.Errorf("expected %s, got %s", tt.input, email.String())
			}
		})
	}

	// URL validation tests - using helper
	testValidation(t, []validationCase{
		{"valid https", "https://example.com", false},
		{"valid http", "http://example.com", false},
		{"valid with path", "https://example.com/path", false},
		{"empty", "", true},
		{"no scheme", "example.com", true},
		{"ftp not allowed", "ftp://example.com", true},
		{"no host", "https:///path", true},
	}, NewURL)
}

func TestParseEmailError(t *testing.T) {
	t.Parallel()
	_, err := NewEmail("invalid-email")
	if err == nil {
		t.Error("expected error for invalid email")
	}
	if !errors.Is(err, pkgerrors.ErrInvalidEmail) {
		t.Errorf("expected ErrInvalidEmail, got %v", err)
	}
}

func TestParseEmail(t *testing.T) {
	t.Parallel()
	email, err := NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if email.String() != "test@example.com" {
		t.Errorf("expected test@example.com, got %s", email.String())
	}
}

func TestEmailNormalize(t *testing.T) {
	t.Parallel()
	email, _ := NewEmail("User@Example.COM")
	normalized := email.Normalize()
	if normalized.String() != "User@example.com" {
		t.Errorf("expected User@example.com, got %s", normalized.String())
	}
}

func TestEmailIsZero(t *testing.T) {
	t.Parallel()
	var zero Email
	if !zero.IsZero() {
		t.Error("zero Email should be zero")
	}

	email, _ := NewEmail("test@example.com")
	if email.IsZero() {
		t.Error("non-zero Email should not be zero")
	}
}

func TestURLParts(t *testing.T) {
	t.Parallel()
	url, _ := NewURL("https://example.com:8080/path/to/resource")
	if url.Scheme() != "https" {
		t.Errorf("expected scheme https, got %s", url.Scheme())
	}
	if url.Host() != "example.com:8080" {
		t.Errorf("expected host example.com:8080, got %s", url.Host())
	}
	if url.Path() != "/path/to/resource" {
		t.Errorf("expected path /path/to/resource, got %s", url.Path())
	}
}

func TestURLIsZero(t *testing.T) {
	t.Parallel()
	var zero URL
	if !zero.IsZero() {
		t.Error("zero URL should be zero")
	}

	url, _ := NewURL("https://example.com")
	if url.IsZero() {
		t.Error("non-zero URL should not be zero")
	}
}

func TestParseURLError(t *testing.T) {
	t.Parallel()
	_, err := NewURL("not-a-valid-url")
	if err == nil {
		t.Error("expected error for invalid URL")
	}
	if !errors.Is(err, pkgerrors.ErrInvalidURL) {
		t.Errorf("expected ErrInvalidURL, got %v", err)
	}
}

func TestParseURL(t *testing.T) {
	t.Parallel()
	url, err := NewURL("https://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if url.String() != "https://example.com" {
		t.Errorf("expected https://example.com, got %s", url.String())
	}
}
