package types

import (
	"errors"
	"testing"
	"time"

	pkgerrors "github.com/larsartmann/go-composable-business-types/pkg/errors"
)

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

func TestTimestamp(t *testing.T) {
	t.Parallel()
	now := time.Now()
	ts := NewTimestamp(now)

	if !ts.Equal(now) {
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
	t.Parallel()
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
	t.Parallel()
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

func TestURL(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
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

func TestPercentageCompare(t *testing.T) {
	t.Parallel()
	tests := []struct {
		a, b     Percentage
		expected int
	}{
		{NewPercentage(50), NewPercentage(50), 0},
		{NewPercentage(30), NewPercentage(50), -1},
		{NewPercentage(70), NewPercentage(50), 1},
	}

	for _, tt := range tests {
		result := tt.a.Compare(tt.b)
		if result != tt.expected {
			t.Errorf("Compare(%d, %d) = %d, expected %d", tt.a, tt.b, result, tt.expected)
		}
	}
}

func TestCentsCompare(t *testing.T) {
	t.Parallel()
	tests := []struct {
		a, b     Cents
		expected int
	}{
		{NewCents(100), NewCents(100), 0},
		{NewCents(50), NewCents(100), -1},
		{NewCents(150), NewCents(100), 1},
	}

	for _, tt := range tests {
		result := tt.a.Compare(tt.b)
		if result != tt.expected {
			t.Errorf("Compare(%d, %d) = %d, expected %d", tt.a, tt.b, result, tt.expected)
		}
	}
}

func TestTimestampCompare(t *testing.T) {
	t.Parallel()
	t1 := NewTimestamp(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	t2 := NewTimestamp(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC))
	t3 := NewTimestamp(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))

	if t1.Compare(t2) != -1 {
		t.Error("t1 should be before t2")
	}
	if t2.Compare(t1) != 1 {
		t.Error("t2 should be after t1")
	}
	if t1.Compare(t3) != 0 {
		t.Error("t1 should equal t3")
	}
}

func TestDurationCompare(t *testing.T) {
	t.Parallel()
	d1 := NewDuration(time.Hour)
	d2 := NewDuration(2 * time.Hour)
	d3 := NewDuration(time.Hour)

	if d1.Compare(d2) != -1 {
		t.Error("d1 should be less than d2")
	}
	if d2.Compare(d1) != 1 {
		t.Error("d2 should be greater than d1")
	}
	if d1.Compare(d3) != 0 {
		t.Error("d1 should equal d3")
	}
}

// SQL Scanner/Valuer tests
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

func TestTimestampSQL(t *testing.T) {
	t.Parallel()
	// Test Value
	ts := NewTimestamp(time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC))
	val, err := ts.Value()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	valTime, ok := val.(time.Time)
	if !ok {
		t.Errorf("expected time.Time, got %T", val)
		return
	}
	if !valTime.Equal(ts.Time) {
		t.Errorf("expected %v, got %v", ts.Time, val)
	}

	// Test Value for zero
	var zero Timestamp
	val, err = zero.Value()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if val != nil {
		t.Errorf("expected nil, got %v", val)
	}

	// Test Scan with time.Time
	var ts2 Timestamp
	inputTime := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	if err := ts2.Scan(inputTime); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !ts2.Equal(inputTime) {
		t.Errorf("expected %v, got %v", inputTime, ts2.Time)
	}

	// Test Scan with string (RFC3339)
	var ts3 Timestamp
	if err := ts3.Scan("2024-03-15T12:00:00Z"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	expected := time.Date(2024, 3, 15, 12, 0, 0, 0, time.UTC)
	if !ts3.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, ts3.Time)
	}

	// Test Scan with []byte (RFC3339)
	var ts4 Timestamp
	if err := ts4.Scan([]byte("2024-09-01T00:00:00Z")); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	expected2 := time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC)
	if !ts4.Equal(expected2) {
		t.Errorf("expected %v, got %v", expected2, ts4.Time)
	}

	// Test Scan with nil
	ts5 := NewTimestamp(time.Now())
	if err := ts5.Scan(nil); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !ts5.Time.IsZero() {
		t.Error("expected zero time after scanning nil")
	}

	// Test Scan with invalid type
	var ts6 Timestamp
	if err := ts6.Scan(123); err == nil {
		t.Error("expected error for invalid type")
	}

	// Test Scan with invalid string
	var ts7 Timestamp
	if err := ts7.Scan("not-a-timestamp"); err == nil {
		t.Error("expected error for invalid timestamp string")
	}
}

func TestDurationSQL(t *testing.T) {
	t.Parallel()
	// Test Value
	d := NewDuration(time.Hour + 30*time.Minute)
	val, err := d.Value()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if val != int64(time.Hour+30*time.Minute) {
		t.Errorf("expected %d, got %v", int64(time.Hour+30*time.Minute), val)
	}

	// Test Value for zero
	var zero Duration
	val, err = zero.Value()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if val != nil {
		t.Errorf("expected nil, got %v", val)
	}

	// Test Scan with int64
	var d2 Duration
	if err := d2.Scan(int64(time.Hour)); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if d2.Duration != time.Hour {
		t.Errorf("expected %v, got %v", time.Hour, d2.Duration)
	}

	// Test Scan with float64
	var d3 Duration
	if err := d3.Scan(float64(2 * time.Hour)); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if d3.Duration != 2*time.Hour {
		t.Errorf("expected %v, got %v", 2*time.Hour, d3.Duration)
	}

	// Test Scan with string
	var d4 Duration
	if err := d4.Scan("30m"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if d4.Duration != 30*time.Minute {
		t.Errorf("expected 30m, got %v", d4.Duration)
	}

	// Test Scan with []byte
	var d5 Duration
	if err := d5.Scan([]byte("1h30m")); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if d5.Duration != time.Hour+30*time.Minute {
		t.Errorf("expected 1h30m, got %v", d5.Duration)
	}

	// Test Scan with nil
	d6 := NewDuration(time.Hour)
	if err := d6.Scan(nil); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if d6.Duration != 0 {
		t.Errorf("expected 0, got %v", d6.Duration)
	}

	// Test Scan with empty string
	var d7 Duration
	if err := d7.Scan(""); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if d7.Duration != 0 {
		t.Errorf("expected 0, got %v", d7.Duration)
	}

	// Test Scan with empty []byte
	var d8 Duration
	if err := d8.Scan([]byte{}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if d8.Duration != 0 {
		t.Errorf("expected 0, got %v", d8.Duration)
	}

	// Test Scan with invalid type
	var d9 Duration
	if err := d9.Scan(struct{}{}); err == nil {
		t.Error("expected error for invalid type")
	}

	// Test Scan with invalid string
	var d10 Duration
	if err := d10.Scan("not-a-duration"); err == nil {
		t.Error("expected error for invalid duration string")
	}
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

func TestDurationJSON(t *testing.T) {
	t.Parallel()
	// Test MarshalJSON
	d := NewDuration(time.Hour + 30*time.Minute)
	data, err := d.MarshalJSON()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	expected := `"1h30m0s"`
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}

	// Test UnmarshalJSON
	var d2 Duration
	if err := d2.UnmarshalJSON([]byte(`"2h15m0s"`)); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if d2.Duration != 2*time.Hour+15*time.Minute {
		t.Errorf("expected 2h15m0s, got %v", d2.Duration)
	}

	// Test Empty
	var d3 Duration
	if err := d3.UnmarshalJSON([]byte(`""`)); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if d3.Duration != 0 {
		t.Errorf("expected 0, got %v", d3.Duration)
	}

	// Test Round-trip
	var d4 Duration
	if err := d4.UnmarshalJSON(data); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if d4.Duration != d.Duration {
		t.Errorf("round-trip failed: expected %v, got %v", d.Duration, d4.Duration)
	}
}
