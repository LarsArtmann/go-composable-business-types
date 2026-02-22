package cbt

import (
	"database/sql"
	"database/sql/driver"
	"testing"
	"time"
)

// ================================================================================
// NanoId SQL Tests
// ================================================================================

func TestNanoId_Scan(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var id NanoId
		err := id.Scan(nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !id.IsZero() {
			t.Errorf("expected empty NanoId, got %v", id)
		}
	})

	t.Run("string valid", func(t *testing.T) {
		var id NanoId
		validID := NewNanoId()
		err := id.Scan(validID.String())
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if id.String() != validID.String() {
			t.Errorf("expected %v, got %v", validID, id)
		}
	})

	t.Run("string empty", func(t *testing.T) {
		var id NanoId
		err := id.Scan("")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !id.IsZero() {
			t.Errorf("expected empty NanoId, got %v", id)
		}
	})

	t.Run("[]byte valid", func(t *testing.T) {
		var id NanoId
		validID := NewNanoId()
		err := id.Scan([]byte(validID.String()))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if id.String() != validID.String() {
			t.Errorf("expected %v, got %v", validID, id)
		}
	})

	t.Run("[]byte empty", func(t *testing.T) {
		var id NanoId
		err := id.Scan([]byte{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !id.IsZero() {
			t.Errorf("expected empty NanoId, got %v", id)
		}
	})

	t.Run("invalid type", func(t *testing.T) {
		var id NanoId
		err := id.Scan(123)
		if err == nil {
			t.Error("expected error for invalid type")
		}
	})

	t.Run("too short", func(t *testing.T) {
		var id NanoId
		err := id.Scan("short")
		if err == nil {
			t.Error("expected error for too short NanoId")
		}
	})

	t.Run("invalid characters", func(t *testing.T) {
		var id NanoId
		err := id.Scan("invalid!@#$%^&*()")
		if err == nil {
			t.Error("expected error for invalid characters")
		}
	})
}

func TestNanoId_Value(t *testing.T) {
	t.Run("zero returns nil", func(t *testing.T) {
		id := NanoId{}
		val, err := id.Value()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if val != nil {
			t.Errorf("expected nil, got %v", val)
		}
	})

	t.Run("non-zero returns string", func(t *testing.T) {
		id := NewNanoId()
		val, err := id.Value()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if val != id.String() {
			t.Errorf("expected %v, got %v", id.String(), val)
		}
	})
}

// ================================================================================
// Email SQL Tests
// ================================================================================

func TestEmail_Scan(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var e Email
		err := e.Scan(nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !e.IsZero() {
			t.Errorf("expected zero Email, got %v", e)
		}
	})

	t.Run("string valid", func(t *testing.T) {
		var e Email
		err := e.Scan("test@example.com")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if e.String() != "test@example.com" {
			t.Errorf("expected test@example.com, got %v", e)
		}
	})

	t.Run("string empty", func(t *testing.T) {
		var e Email
		err := e.Scan("")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !e.IsZero() {
			t.Errorf("expected zero Email, got %v", e)
		}
	})

	t.Run("[]byte valid", func(t *testing.T) {
		var e Email
		err := e.Scan([]byte("test@example.com"))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if e.String() != "test@example.com" {
			t.Errorf("expected test@example.com, got %v", e)
		}
	})

	t.Run("[]byte empty", func(t *testing.T) {
		var e Email
		err := e.Scan([]byte{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !e.IsZero() {
			t.Errorf("expected zero Email, got %v", e)
		}
	})

	t.Run("invalid type", func(t *testing.T) {
		var e Email
		err := e.Scan(123)
		if err == nil {
			t.Error("expected error for invalid type")
		}
	})

	t.Run("invalid email", func(t *testing.T) {
		var e Email
		err := e.Scan("not-an-email")
		if err == nil {
			t.Error("expected error for invalid email")
		}
	})
}

func TestEmail_Value(t *testing.T) {
	t.Run("zero returns nil", func(t *testing.T) {
		e := Email("")
		val, err := e.Value()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if val != nil {
			t.Errorf("expected nil, got %v", val)
		}
	})

	t.Run("non-zero returns string", func(t *testing.T) {
		e := MustParseEmail("test@example.com")
		val, err := e.Value()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if val != "test@example.com" {
			t.Errorf("expected test@example.com, got %v", val)
		}
	})
}

// ================================================================================
// URL SQL Tests
// ================================================================================

func TestURL_Scan(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var u URL
		err := u.Scan(nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !u.IsZero() {
			t.Errorf("expected zero URL, got %v", u)
		}
	})

	t.Run("string valid http", func(t *testing.T) {
		var u URL
		err := u.Scan("http://example.com")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if u.String() != "http://example.com" {
			t.Errorf("expected http://example.com, got %v", u)
		}
	})

	t.Run("string valid https", func(t *testing.T) {
		var u URL
		err := u.Scan("https://example.com/path")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if u.String() != "https://example.com/path" {
			t.Errorf("expected https://example.com/path, got %v", u)
		}
	})

	t.Run("string empty", func(t *testing.T) {
		var u URL
		err := u.Scan("")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !u.IsZero() {
			t.Errorf("expected zero URL, got %v", u)
		}
	})

	t.Run("[]byte valid", func(t *testing.T) {
		var u URL
		err := u.Scan([]byte("https://example.com"))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if u.String() != "https://example.com" {
			t.Errorf("expected https://example.com, got %v", u)
		}
	})

	t.Run("[]byte empty", func(t *testing.T) {
		var u URL
		err := u.Scan([]byte{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !u.IsZero() {
			t.Errorf("expected zero URL, got %v", u)
		}
	})

	t.Run("invalid type", func(t *testing.T) {
		var u URL
		err := u.Scan(123)
		if err == nil {
			t.Error("expected error for invalid type")
		}
	})

	t.Run("invalid URL no scheme", func(t *testing.T) {
		var u URL
		err := u.Scan("example.com")
		if err == nil {
			t.Error("expected error for URL without scheme")
		}
	})

	t.Run("invalid URL ftp scheme", func(t *testing.T) {
		var u URL
		err := u.Scan("ftp://example.com")
		if err == nil {
			t.Error("expected error for non-http(s) scheme")
		}
	})
}

func TestURL_Value(t *testing.T) {
	t.Run("zero returns nil", func(t *testing.T) {
		u := URL("")
		val, err := u.Value()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if val != nil {
			t.Errorf("expected nil, got %v", val)
		}
	})

	t.Run("non-zero returns string", func(t *testing.T) {
		u := MustParseURL("https://example.com")
		val, err := u.Value()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if val != "https://example.com" {
			t.Errorf("expected https://example.com, got %v", val)
		}
	})
}

// ================================================================================
// Cents SQL Tests
// ================================================================================

func TestCents_Scan(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var c Cents
		err := c.Scan(nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !c.IsZero() {
			t.Errorf("expected zero Cents, got %v", c)
		}
	})

	t.Run("int64", func(t *testing.T) {
		var c Cents
		err := c.Scan(int64(12345))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if c.Int64() != 12345 {
			t.Errorf("expected 12345, got %v", c.Int64())
		}
	})

	t.Run("float64", func(t *testing.T) {
		var c Cents
		err := c.Scan(float64(12345.67))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if c.Int64() != 12345 {
			t.Errorf("expected 12345, got %v", c.Int64())
		}
	})

	t.Run("[]byte valid", func(t *testing.T) {
		var c Cents
		err := c.Scan([]byte("999"))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if c.Int64() != 999 {
			t.Errorf("expected 999, got %v", c.Int64())
		}
	})

	t.Run("[]byte invalid", func(t *testing.T) {
		var c Cents
		err := c.Scan([]byte("not-a-number"))
		if err == nil {
			t.Error("expected error for invalid number")
		}
	})

	t.Run("invalid type", func(t *testing.T) {
		var c Cents
		err := c.Scan("string")
		if err == nil {
			t.Error("expected error for invalid type")
		}
	})

	t.Run("negative value", func(t *testing.T) {
		var c Cents
		err := c.Scan(int64(-500))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if c.Int64() != -500 {
			t.Errorf("expected -500, got %v", c.Int64())
		}
	})
}

func TestCents_Value(t *testing.T) {
	t.Run("returns int64", func(t *testing.T) {
		c := NewCents(12345)
		val, err := c.Value()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if val != int64(12345) {
			t.Errorf("expected 12345, got %v", val)
		}
	})

	t.Run("zero returns int64(0)", func(t *testing.T) {
		c := Cents(0)
		val, err := c.Value()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if val != int64(0) {
			t.Errorf("expected 0, got %v", val)
		}
	})
}

// ================================================================================
// Timestamp SQL Tests
// ================================================================================

func TestTimestamp_Scan(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var ts Timestamp
		err := ts.Scan(nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !ts.IsZero() {
			t.Errorf("expected zero Timestamp, got %v", ts)
		}
	})

	t.Run("time.Time", func(t *testing.T) {
		var ts Timestamp
		now := time.Now()
		err := ts.Scan(now)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !ts.Time.Equal(now) {
			t.Errorf("expected %v, got %v", now, ts.Time)
		}
	})

	t.Run("string RFC3339", func(t *testing.T) {
		var ts Timestamp
		now := time.Now().UTC()
		timeStr := now.Format(time.RFC3339Nano)
		err := ts.Scan(timeStr)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		// Parse with RFC3339Nano for comparison since we used that in Scan
		expected, _ := time.Parse(time.RFC3339Nano, timeStr)
		if !ts.Time.Equal(expected) {
			t.Errorf("expected %v, got %v", expected, ts.Time)
		}
	})

	t.Run("[]byte RFC3339", func(t *testing.T) {
		var ts Timestamp
		now := time.Now().UTC()
		timeStr := now.Format(time.RFC3339Nano)
		err := ts.Scan([]byte(timeStr))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		expected, _ := time.Parse(time.RFC3339Nano, timeStr)
		if !ts.Time.Equal(expected) {
			t.Errorf("expected %v, got %v", expected, ts.Time)
		}
	})

	t.Run("string invalid format", func(t *testing.T) {
		var ts Timestamp
		err := ts.Scan("not-a-time")
		if err == nil {
			t.Error("expected error for invalid time format")
		}
	})

	t.Run("[]byte invalid format", func(t *testing.T) {
		var ts Timestamp
		err := ts.Scan([]byte("not-a-time"))
		if err == nil {
			t.Error("expected error for invalid time format")
		}
	})

	t.Run("invalid type", func(t *testing.T) {
		var ts Timestamp
		err := ts.Scan(123)
		if err == nil {
			t.Error("expected error for invalid type")
		}
	})
}

func TestTimestamp_Value(t *testing.T) {
	t.Run("zero returns nil", func(t *testing.T) {
		ts := Timestamp{}
		val, err := ts.Value()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if val != nil {
			t.Errorf("expected nil, got %v", val)
		}
	})

	t.Run("non-zero returns time.Time", func(t *testing.T) {
		now := time.Now()
		ts := NewTimestamp(now)
		val, err := ts.Value()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if tm, ok := val.(time.Time); !ok || !tm.Equal(now) {
			t.Errorf("expected %v, got %v", now, val)
		}
	})
}

// ================================================================================
// BoundedString SQL Tests
// ================================================================================

func TestBoundedString_Scan(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var bs BoundedString
		err := bs.Scan(nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !bs.IsZero() {
			t.Errorf("expected empty BoundedString, got %v", bs)
		}
	})

	t.Run("string", func(t *testing.T) {
		var bs BoundedString
		err := bs.Scan("hello world")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if bs.String() != "hello world" {
			t.Errorf("expected 'hello world', got %v", bs.String())
		}
	})

	t.Run("string empty", func(t *testing.T) {
		var bs BoundedString
		err := bs.Scan("")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !bs.IsZero() {
			t.Errorf("expected empty BoundedString, got %v", bs)
		}
	})

	t.Run("[]byte", func(t *testing.T) {
		var bs BoundedString
		err := bs.Scan([]byte("hello"))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if bs.String() != "hello" {
			t.Errorf("expected 'hello', got %v", bs.String())
		}
	})

	t.Run("[]byte empty", func(t *testing.T) {
		var bs BoundedString
		err := bs.Scan([]byte{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !bs.IsZero() {
			t.Errorf("expected empty BoundedString, got %v", bs)
		}
	})

	t.Run("invalid type", func(t *testing.T) {
		var bs BoundedString
		err := bs.Scan(123)
		if err == nil {
			t.Error("expected error for invalid type")
		}
	})

	t.Run("sets bounds correctly", func(t *testing.T) {
		var bs BoundedString
		_ = bs.Scan("hello")
		if bs.MinLen() != 0 {
			t.Errorf("expected min 0, got %d", bs.MinLen())
		}
		if bs.MaxLen() != 5 {
			t.Errorf("expected max 5, got %d", bs.MaxLen())
		}
	})
}

func TestBoundedString_Value(t *testing.T) {
	t.Run("empty returns nil", func(t *testing.T) {
		bs := BoundedString{}
		val, err := bs.Value()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if val != nil {
			t.Errorf("expected nil, got %v", val)
		}
	})

	t.Run("non-empty returns string", func(t *testing.T) {
		bs := MustBoundedString(0, 100, "hello")
		val, err := bs.Value()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if val != "hello" {
			t.Errorf("expected 'hello', got %v", val)
		}
	})
}

// ================================================================================
// ID[B, V] SQL Tests
// ================================================================================

type testBrand struct{}

func TestID_Scan_StringBrand(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var id ID[testBrand, string]
		err := id.Scan(nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !id.IsZero() {
			t.Errorf("expected zero ID, got %v", id)
		}
	})

	t.Run("string", func(t *testing.T) {
		var id ID[testBrand, string]
		err := id.Scan("user-123")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if id.Get() != "user-123" {
			t.Errorf("expected 'user-123', got %v", id.Get())
		}
	})

	t.Run("[]byte", func(t *testing.T) {
		var id ID[testBrand, string]
		err := id.Scan([]byte("user-456"))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if id.Get() != "user-456" {
			t.Errorf("expected 'user-456', got %v", id.Get())
		}
	})

	t.Run("invalid type", func(t *testing.T) {
		var id ID[testBrand, string]
		err := id.Scan(123)
		if err == nil {
			t.Error("expected error for invalid type")
		}
	})
}

func TestID_Scan_Int64Brand(t *testing.T) {
	type intBrand struct{}

	t.Run("nil", func(t *testing.T) {
		var id ID[intBrand, int64]
		err := id.Scan(nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !id.IsZero() {
			t.Errorf("expected zero ID, got %v", id)
		}
	})

	t.Run("int64", func(t *testing.T) {
		var id ID[intBrand, int64]
		err := id.Scan(int64(12345))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if id.Get() != int64(12345) {
			t.Errorf("expected 12345, got %v", id.Get())
		}
	})

	t.Run("float64", func(t *testing.T) {
		var id ID[intBrand, int64]
		err := id.Scan(float64(12345.67))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if id.Get() != int64(12345) {
			t.Errorf("expected 12345, got %v", id.Get())
		}
	})

	t.Run("invalid type", func(t *testing.T) {
		var id ID[intBrand, int64]
		err := id.Scan("string")
		if err == nil {
			t.Error("expected error for invalid type")
		}
	})
}

func TestID_Value_StringBrand(t *testing.T) {
	t.Run("zero returns nil", func(t *testing.T) {
		id := ID[testBrand, string]{}
		val, err := id.Value()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if val != nil {
			t.Errorf("expected nil, got %v", val)
		}
	})

	t.Run("non-zero returns string", func(t *testing.T) {
		id := NewID[testBrand, string]("user-123")
		val, err := id.Value()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if val != "user-123" {
			t.Errorf("expected 'user-123', got %v", val)
		}
	})
}

func TestID_Value_Int64Brand(t *testing.T) {
	type intBrand struct{}

	t.Run("zero returns nil", func(t *testing.T) {
		id := ID[intBrand, int64]{}
		val, err := id.Value()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if val != nil {
			t.Errorf("expected nil, got %v", val)
		}
	})

	t.Run("non-zero returns int64", func(t *testing.T) {
		id := NewID[intBrand, int64](int64(12345))
		val, err := id.Value()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if val != int64(12345) {
			t.Errorf("expected 12345, got %v", val)
		}
	})
}

// ================================================================================
// Interface Compliance Tests
// ================================================================================

func TestSQL_Interfaces(t *testing.T) {
	// Verify all types implement sql.Scanner and driver.Valuer
	var _ driver.Valuer = NanoId{}
	var _ driver.Valuer = Email("")
	var _ driver.Valuer = URL("")
	var _ driver.Valuer = Cents(0)
	var _ driver.Valuer = Timestamp{}
	var _ driver.Valuer = Duration{}
	var _ driver.Valuer = BoundedString{}
	var _ driver.Valuer = Locale{}
	var _ driver.Valuer = ID[testBrand, string]{}
	var _ driver.Valuer = ID[testBrand, int64]{}

	// Scanner requires pointer types
	var _ sql.Scanner = (*NanoId)(nil)
	var _ sql.Scanner = (*Email)(nil)
	var _ sql.Scanner = (*URL)(nil)
	var _ sql.Scanner = (*Cents)(nil)
	var _ sql.Scanner = (*Timestamp)(nil)
	var _ sql.Scanner = (*Duration)(nil)
	var _ sql.Scanner = (*BoundedString)(nil)
	var _ sql.Scanner = (*Locale)(nil)
	var _ sql.Scanner = (*ID[testBrand, string])(nil)
	var _ sql.Scanner = (*ID[testBrand, int64])(nil)
}

// ================================================================================
// Locale SQL Tests
// ================================================================================

func TestLocale_Scan(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var l Locale
		err := l.Scan(nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !l.IsZero() {
			t.Errorf("expected zero Locale, got %v", l)
		}
	})

	t.Run("string valid", func(t *testing.T) {
		var l Locale
		err := l.Scan("en-US")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if l.String() != "en-US" {
			t.Errorf("expected en-US, got %v", l)
		}
	})

	t.Run("string empty", func(t *testing.T) {
		var l Locale
		err := l.Scan("")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !l.IsZero() {
			t.Errorf("expected zero Locale, got %v", l)
		}
	})

	t.Run("[]byte valid", func(t *testing.T) {
		var l Locale
		err := l.Scan([]byte("de-DE"))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if l.String() != "de-DE" {
			t.Errorf("expected de-DE, got %v", l)
		}
	})

	t.Run("[]byte empty", func(t *testing.T) {
		var l Locale
		err := l.Scan([]byte{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !l.IsZero() {
			t.Errorf("expected zero Locale, got %v", l)
		}
	})

	t.Run("invalid type", func(t *testing.T) {
		var l Locale
		err := l.Scan(123)
		if err == nil {
			t.Error("expected error for invalid type")
		}
	})

	t.Run("invalid locale string", func(t *testing.T) {
		var l Locale
		err := l.Scan("invalid-locale-format-!!!")
		if err == nil {
			t.Error("expected error for invalid locale")
		}
	})
}

func TestLocale_Value(t *testing.T) {
	t.Run("zero returns nil", func(t *testing.T) {
		l := Locale{}
		val, err := l.Value()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if val != nil {
			t.Errorf("expected nil, got %v", val)
		}
	})

	t.Run("non-zero returns string", func(t *testing.T) {
		l := MustParseLocale("en-US")
		val, err := l.Value()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if val != "en-US" {
			t.Errorf("expected en-US, got %v", val)
		}
	})
}

// ================================================================================
// Duration SQL Tests
// ================================================================================

func TestDuration_Scan(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var d Duration
		err := d.Scan(nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !d.IsZero() {
			t.Errorf("expected zero Duration, got %v", d)
		}
	})

	t.Run("int64 nanoseconds", func(t *testing.T) {
		var d Duration
		err := d.Scan(int64(5 * time.Second))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if d.Duration != 5*time.Second {
			t.Errorf("expected 5s, got %v", d)
		}
	})

	t.Run("float64 nanoseconds", func(t *testing.T) {
		var d Duration
		err := d.Scan(float64(3 * time.Second))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if d.Duration != 3*time.Second {
			t.Errorf("expected 3s, got %v", d)
		}
	})

	t.Run("string valid", func(t *testing.T) {
		var d Duration
		err := d.Scan("1h30m")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if d.Duration != 90*time.Minute {
			t.Errorf("expected 1h30m, got %v", d)
		}
	})

	t.Run("string empty", func(t *testing.T) {
		var d Duration
		err := d.Scan("")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !d.IsZero() {
			t.Errorf("expected zero Duration, got %v", d)
		}
	})

	t.Run("[]byte valid", func(t *testing.T) {
		var d Duration
		err := d.Scan([]byte("500ms"))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if d.Duration != 500*time.Millisecond {
			t.Errorf("expected 500ms, got %v", d)
		}
	})

	t.Run("[]byte empty", func(t *testing.T) {
		var d Duration
		err := d.Scan([]byte{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !d.IsZero() {
			t.Errorf("expected zero Duration, got %v", d)
		}
	})

	t.Run("invalid type", func(t *testing.T) {
		var d Duration
		err := d.Scan("invalid")
		if err == nil {
			t.Error("expected error for invalid duration string")
		}
	})

	t.Run("invalid source type", func(t *testing.T) {
		var d Duration
		err := d.Scan([]int{1, 2, 3})
		if err == nil {
			t.Error("expected error for invalid type")
		}
	})
}

func TestDuration_Value(t *testing.T) {
	t.Run("zero returns nil", func(t *testing.T) {
		d := Duration{}
		val, err := d.Value()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if val != nil {
			t.Errorf("expected nil, got %v", val)
		}
	})

	t.Run("non-zero returns int64 nanoseconds", func(t *testing.T) {
		d := NewDuration(5 * time.Second)
		val, err := d.Value()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if val != int64(5*time.Second) {
			t.Errorf("expected %d, got %v", int64(5*time.Second), val)
		}
	})
}

// ================================================================================
// Compare and IsZero Method Tests
// ================================================================================

func TestCents_Compare(t *testing.T) {
	t.Run("less than", func(t *testing.T) {
		c1 := NewCents(100)
		c2 := NewCents(200)
		if c1.Compare(c2) != -1 {
			t.Errorf("expected -1, got %d", c1.Compare(c2))
		}
	})

	t.Run("equal", func(t *testing.T) {
		c1 := NewCents(100)
		c2 := NewCents(100)
		if c1.Compare(c2) != 0 {
			t.Errorf("expected 0, got %d", c1.Compare(c2))
		}
	})

	t.Run("greater than", func(t *testing.T) {
		c1 := NewCents(200)
		c2 := NewCents(100)
		if c1.Compare(c2) != 1 {
			t.Errorf("expected 1, got %d", c1.Compare(c2))
		}
	})

	t.Run("negative values", func(t *testing.T) {
		c1 := NewCents(-100)
		c2 := NewCents(100)
		if c1.Compare(c2) != -1 {
			t.Errorf("expected -1 for negative vs positive, got %d", c1.Compare(c2))
		}
	})
}

func TestTimestamp_Compare(t *testing.T) {
	t1 := NewTimestamp(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	t2 := NewTimestamp(time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC))

	t.Run("less than", func(t *testing.T) {
		if t1.Compare(t2) != -1 {
			t.Errorf("expected -1, got %d", t1.Compare(t2))
		}
	})

	t.Run("equal", func(t *testing.T) {
		if t1.Compare(t1) != 0 {
			t.Errorf("expected 0, got %d", t1.Compare(t1))
		}
	})

	t.Run("greater than", func(t *testing.T) {
		if t2.Compare(t1) != 1 {
			t.Errorf("expected 1, got %d", t2.Compare(t1))
		}
	})
}

func TestDuration_IsZero(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		d := Duration{}
		if !d.IsZero() {
			t.Error("expected IsZero to be true")
		}
	})

	t.Run("non-zero", func(t *testing.T) {
		d := NewDuration(5 * time.Second)
		if d.IsZero() {
			t.Error("expected IsZero to be false")
		}
	})
}

func TestDuration_Compare(t *testing.T) {
	d1 := NewDuration(5 * time.Second)
	d2 := NewDuration(10 * time.Second)

	t.Run("less than", func(t *testing.T) {
		if d1.Compare(d2) != -1 {
			t.Errorf("expected -1, got %d", d1.Compare(d2))
		}
	})

	t.Run("equal", func(t *testing.T) {
		if d1.Compare(d1) != 0 {
			t.Errorf("expected 0, got %d", d1.Compare(d1))
		}
	})

	t.Run("greater than", func(t *testing.T) {
		if d2.Compare(d1) != 1 {
			t.Errorf("expected 1, got %d", d2.Compare(d1))
		}
	})
}

// sql.Scanner interface check (using the standard library type indirectly)
type sqlScanner interface {
	Scan(src any) error
}

// ================================================================================
// Money SQL Tests (inherited from currency.Amount)
// ================================================================================

func TestMoney_Scan(t *testing.T) {
	t.Run("valid composite format", func(t *testing.T) {
		var m Money
		// currency.Amount expects PostgreSQL composite format: "(9.99,USD)"
		err := m.Scan("(9.99,USD)")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if m.Number() != "9.99" {
			t.Errorf("expected 9.99, got %s", m.Number())
		}
		if m.CurrencyCode() != "USD" {
			t.Errorf("expected USD, got %s", m.CurrencyCode())
		}
	})

	t.Run("zero value", func(t *testing.T) {
		var m Money
		err := m.Scan("")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		// Empty string results in zero value
	})

	t.Run("invalid type", func(t *testing.T) {
		var m Money
		err := m.Scan(123)
		if err == nil {
			t.Error("expected error for invalid type")
		}
	})
}

func TestMoney_Value(t *testing.T) {
	t.Run("returns composite format", func(t *testing.T) {
		m, err := NewMoney("99.99", "USD")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		val, err := m.Value()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		// Returns PostgreSQL composite format
		if val != "(99.99,USD)" {
			t.Errorf("expected (99.99,USD), got %v", val)
		}
	})
}

func TestMoney_ImplementsInterfaces(t *testing.T) {
	// Money is an alias for currency.Amount, which implements these
	var _ driver.Valuer = Money{}
	var _ sql.Scanner = (*Money)(nil)
}
