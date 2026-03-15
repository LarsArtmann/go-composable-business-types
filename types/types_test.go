package types

import (
	"testing"
	"time"
)

func TestNewPercentage(t *testing.T) {
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
	zero := NewPercentage(0)
	half := NewPercentage(50)
	full := NewPercentage(100)

	if !zero.IsZero() || !zero.IsMin() || full.IsMax() {
		t.Error("percentage 0 helpers failed")
	}
	if half.IsZero() || half.IsMin() || full.IsMax() {
		t.Error("percentage 50 helpers failed")
	}
	if full.IsZero() || full.IsMin() || !full.IsMax() {
		t.Error("percentage 100 helpers failed")
	}
}

func TestNewCents(t *testing.T) {
	c := NewCents(1099)
	if c.Int64() != 1099 {
		t.Errorf("expected 1099, got %d", c.Int64())
	}
	if c.Float64() != 10.99 {
		t.Errorf("expected 10.99, got %f", c.Float64())
	}
}

func TestCentsMath(t *testing.T) {
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

func TestTimestamp(t *testing.T) {
	now := time.Now()
	ts := NewTimestamp(now)

	if !ts.Time.Equal(now) {
		t.Error("Timestamp should store time")
	}

	// IsZero
	var zeroTs Timestamp
	if !zeroTs.IsZero() {
		t.Error("zero timestamp should be zero")
	}

	// Now should not be zero
	nowTs := Now()
	if nowTs.IsZero() {
		t.Error("Now() should not be zero")
	}
}

func TestDuration(t *testing.T) {
	d := NewDuration(time.Hour)
	if d.Duration != time.Hour {
		t.Errorf("expected 1 hour, got %v", d.Duration)
	}

	var zeroDur Duration
	if !zeroDur.IsZero() {
		t.Error("zero duration should be zero")
	}
}

func TestEmail(t *testing.T) {
	tests := []struct {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
}

func TestEmailParts(t *testing.T) {
	email, _ := NewEmail("user@example.com")
	if email.LocalPart() != "user" {
		t.Errorf("expected local part 'user', got %s", email.LocalPart())
	}
	if email.Domain() != "example.com" {
		t.Errorf("expected domain 'example.com', got %s", email.Domain())
	}
}

func TestURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid https", "https://example.com", false},
		{"valid http", "http://example.com", false},
		{"valid with path", "https://example.com/path", false},
		{"empty", "", true},
		{"no scheme", "example.com", true},
		{"ftp not allowed", "ftp://example.com", true},
		{"no host", "https:///path", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := NewURL(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if url.String() != tt.input {
				t.Errorf("expected %s, got %s", tt.input, url.String())
			}
		})
	}
}
