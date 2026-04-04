package types

import (
	"testing"
)

// SQL Scanner/Valuer tests for Email, URL, and Cents.
func TestEmailSQL(t *testing.T) {
	t.Parallel()
	// Test Value
	email, _ := NewEmail("test@example.com")

	val, err := email.Value()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if val != "test@example.com" {
		t.Errorf("expected test@example.com, got %v", val)
	}

	// Test Value for zero
	var zero Email

	val, err = zero.Value()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if val != nil {
		t.Errorf("expected nil, got %v", val)
	}

	// Test Scan with string
	var e Email
	if err := e.Scan("scanned@example.com"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if e.String() != "scanned@example.com" {
		t.Errorf("expected scanned@example.com, got %s", e.String())
	}

	// Test Scan with []byte
	var e2 Email
	if err := e2.Scan([]byte("byte@example.com")); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if e2.String() != "byte@example.com" {
		t.Errorf("expected byte@example.com, got %s", e2.String())
	}

	// Test Scan with nil
	var e3 Email

	e3, _ = NewEmail("test@example.com")
	if err := e3.Scan(nil); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !e3.IsZero() {
		t.Error("expected zero value after scanning nil")
	}

	// Test Scan with empty string
	var e4 Email
	if err := e4.Scan(""); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !e4.IsZero() {
		t.Error("expected zero value after scanning empty string")
	}

	// Test Scan with invalid type
	var e5 Email
	if err := e5.Scan(123); err == nil {
		t.Error("expected error for invalid type")
	}

	// Test Scan with invalid email
	var e6 Email
	if err := e6.Scan("not-an-email"); err == nil {
		t.Error("expected error for invalid email")
	}
}

func TestURLSQL(t *testing.T) {
	t.Parallel()
	// Test Value
	url, _ := NewURL("https://example.com")

	val, err := url.Value()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if val != "https://example.com" {
		t.Errorf("expected https://example.com, got %v", val)
	}

	// Test Value for zero
	var zero URL

	val, err = zero.Value()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if val != nil {
		t.Errorf("expected nil, got %v", val)
	}

	// Test Scan with string
	var u URL
	if err := u.Scan("https://scanned.example.com"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if u.String() != "https://scanned.example.com" {
		t.Errorf("expected https://scanned.example.com, got %s", u.String())
	}

	// Test Scan with []byte
	var u2 URL
	if err := u2.Scan([]byte("https://byte.example.com")); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if u2.String() != "https://byte.example.com" {
		t.Errorf("expected https://byte.example.com, got %s", u2.String())
	}

	// Test Scan with nil
	var u3 URL
	if err := u3.Scan(nil); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !u3.IsZero() {
		t.Error("expected zero value after scanning nil")
	}

	// Test Scan with empty string
	var u4 URL
	if err := u4.Scan(""); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !u4.IsZero() {
		t.Error("expected zero value after scanning empty string")
	}

	// Test Scan with invalid type
	var u5 URL
	if err := u5.Scan(123); err == nil {
		t.Error("expected error for invalid type")
	}

	// Test Scan with invalid URL
	var u6 URL
	if err := u6.Scan("not-a-url"); err == nil {
		t.Error("expected error for invalid URL")
	}
}

func TestCentsSQL(t *testing.T) {
	t.Parallel()
	// Test Value
	c := NewCents(1099)

	val, err := c.Value()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if val != int64(1099) {
		t.Errorf("expected 1099, got %v", val)
	}

	// Test Scan with int64
	var c2 Cents
	if err := c2.Scan(int64(500)); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if c2 != 500 {
		t.Errorf("expected 500, got %d", c2)
	}

	// Test Scan with float64
	var c3 Cents
	if err := c3.Scan(float64(750)); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if c3 != 750 {
		t.Errorf("expected 750, got %d", c3)
	}

	// Test Scan with []byte
	var c4 Cents
	if err := c4.Scan([]byte("999")); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if c4 != 999 {
		t.Errorf("expected 999, got %d", c4)
	}

	// Test Scan with nil
	c5 := NewCents(100)
	if err := c5.Scan(nil); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if c5 != 0 {
		t.Errorf("expected 0, got %d", c5)
	}

	// Test Scan with invalid type
	var c6 Cents
	if err := c6.Scan("not-a-number"); err == nil {
		t.Error("expected error for invalid type")
	}

	// Test Scan with invalid []byte
	var c7 Cents
	if err := c7.Scan([]byte("not-a-number")); err == nil {
		t.Error("expected error for invalid []byte")
	}
}
