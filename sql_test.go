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
		assertError(t, id.Scan(123), "invalid type")
	})

	t.Run("too short", func(t *testing.T) {
		var id NanoId
		assertError(t, id.Scan("short"), "too short NanoId")
	})

	t.Run("invalid characters", func(t *testing.T) {
		var id NanoId
		assertError(t, id.Scan("invalid!@#$%^&*()"), "invalid characters")
	})
}

func TestNanoId_Value(t *testing.T) {
	t.Run("zero returns nil", func(t *testing.T) {
		id := NanoId{}
		val, err := id.Value()
		assertNoError(t, err)
		assertNil(t, val)
	})

	t.Run("non-zero returns string", func(t *testing.T) {
		id := NewNanoId()
		val, err := id.Value()
		assertNoError(t, err)
		assertEqual(t, val.(string), id.String())
	})
}

// ================================================================================
// Email SQL Tests
// ================================================================================

func TestEmail_Scan(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var e Email
		assertNoError(t, e.Scan(nil))
		assertZero(t, e)
	})

	t.Run("string valid", func(t *testing.T) {
		var e Email
		assertNoError(t, e.Scan("test@example.com"))
		assertStringEquals(t, e.String(), "test@example.com")
	})

	t.Run("string empty", func(t *testing.T) {
		var e Email
		assertNoError(t, e.Scan(""))
		assertZero(t, e)
	})

	t.Run("[]byte valid", func(t *testing.T) {
		var e Email
		assertNoError(t, e.Scan([]byte("test@example.com")))
		assertStringEquals(t, e.String(), "test@example.com")
	})

	t.Run("[]byte empty", func(t *testing.T) {
		var e Email
		assertNoError(t, e.Scan([]byte{}))
		assertZero(t, e)
	})

	t.Run("invalid type", func(t *testing.T) {
		var e Email
		assertError(t, e.Scan(123), "invalid type")
	})

	t.Run("invalid email", func(t *testing.T) {
		var e Email
		assertError(t, e.Scan("not-an-email"), "invalid email")
	})
}

func TestEmail_Value(t *testing.T) {
	t.Run("zero returns nil", func(t *testing.T) {
		e := Email("")
		val, err := e.Value()
		assertNoError(t, err)
		assertNil(t, val)
	})

	t.Run("non-zero returns string", func(t *testing.T) {
		e := MustParseEmail("test@example.com")
		val, err := e.Value()
		assertNoError(t, err)
		assertEqual(t, val.(string), "test@example.com")
	})
}

// ================================================================================
// URL SQL Tests
// ================================================================================

func TestURL_Scan(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var u URL
		assertNoError(t, u.Scan(nil))
		assertZero(t, u)
	})

	t.Run("string valid http", func(t *testing.T) {
		var u URL
		assertNoError(t, u.Scan("http://example.com"))
		assertStringEquals(t, u.String(), "http://example.com")
	})

	t.Run("string valid https", func(t *testing.T) {
		var u URL
		assertNoError(t, u.Scan("https://example.com/path"))
		assertStringEquals(t, u.String(), "https://example.com/path")
	})

	t.Run("string empty", func(t *testing.T) {
		var u URL
		assertNoError(t, u.Scan(""))
		assertZero(t, u)
	})

	t.Run("[]byte valid", func(t *testing.T) {
		var u URL
		assertNoError(t, u.Scan([]byte("https://example.com")))
		assertStringEquals(t, u.String(), "https://example.com")
	})

	t.Run("[]byte empty", func(t *testing.T) {
		var u URL
		assertNoError(t, u.Scan([]byte{}))
		assertZero(t, u)
	})

	t.Run("invalid type", func(t *testing.T) {
		var u URL
		assertError(t, u.Scan(123), "invalid type")
	})

	t.Run("invalid URL no scheme", func(t *testing.T) {
		var u URL
		assertError(t, u.Scan("example.com"), "URL without scheme")
	})

	t.Run("invalid URL ftp scheme", func(t *testing.T) {
		var u URL
		assertError(t, u.Scan("ftp://example.com"), "non-http(s) scheme")
	})
}

func TestURL_Value(t *testing.T) {
	t.Run("zero returns nil", func(t *testing.T) {
		u := URL("")
		val, err := u.Value()
		assertNoError(t, err)
		assertNil(t, val)
	})

	t.Run("non-zero returns string", func(t *testing.T) {
		u := MustParseURL("https://example.com")
		val, err := u.Value()
		assertNoError(t, err)
		assertEqual(t, val.(string), "https://example.com")
	})
}

// ================================================================================
// Cents SQL Tests
// ================================================================================

func TestCents_Scan(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var c Cents
		assertNoError(t, c.Scan(nil))
		assertZero(t, c)
	})

	t.Run("int64", func(t *testing.T) {
		var c Cents
		assertNoError(t, c.Scan(int64(12345)))
		assertEqual(t, c.Int64(), int64(12345))
	})

	t.Run("float64", func(t *testing.T) {
		var c Cents
		assertNoError(t, c.Scan(float64(12345.67)))
		assertEqual(t, c.Int64(), int64(12345))
	})

	t.Run("[]byte valid", func(t *testing.T) {
		var c Cents
		assertNoError(t, c.Scan([]byte("999")))
		assertEqual(t, c.Int64(), int64(999))
	})

	t.Run("[]byte invalid", func(t *testing.T) {
		var c Cents
		assertError(t, c.Scan([]byte("not-a-number")), "invalid number")
	})

	t.Run("invalid type", func(t *testing.T) {
		var c Cents
		assertError(t, c.Scan("string"), "invalid type")
	})

	t.Run("negative value", func(t *testing.T) {
		var c Cents
		assertNoError(t, c.Scan(int64(-500)))
		assertEqual(t, c.Int64(), int64(-500))
	})
}

func TestCents_Value(t *testing.T) {
	t.Run("returns int64", func(t *testing.T) {
		c := NewCents(12345)
		val, err := c.Value()
		assertNoError(t, err)
		assertEqual(t, val.(int64), int64(12345))
	})

	t.Run("zero returns int64(0)", func(t *testing.T) {
		c := Cents(0)
		val, err := c.Value()
		assertNoError(t, err)
		assertEqual(t, val.(int64), int64(0))
	})
}

// ================================================================================
// Timestamp SQL Tests
// ================================================================================

func TestTimestamp_Scan(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var ts Timestamp
		assertNoError(t, ts.Scan(nil))
		assertZero(t, ts)
	})

	t.Run("time.Time", func(t *testing.T) {
		var ts Timestamp
		now := time.Now()
		assertNoError(t, ts.Scan(now))
		if !ts.Time.Equal(now) {
			t.Errorf("expected %v, got %v", now, ts.Time)
		}
	})

	t.Run("string RFC3339", func(t *testing.T) {
		var ts Timestamp
		now := time.Now().UTC()
		timeStr := now.Format(time.RFC3339Nano)
		assertNoError(t, ts.Scan(timeStr))
		expected, _ := time.Parse(time.RFC3339Nano, timeStr)
		if !ts.Time.Equal(expected) {
			t.Errorf("expected %v, got %v", expected, ts.Time)
		}
	})

	t.Run("[]byte RFC3339", func(t *testing.T) {
		var ts Timestamp
		now := time.Now().UTC()
		timeStr := now.Format(time.RFC3339Nano)
		assertNoError(t, ts.Scan([]byte(timeStr)))
		expected, _ := time.Parse(time.RFC3339Nano, timeStr)
		if !ts.Time.Equal(expected) {
			t.Errorf("expected %v, got %v", expected, ts.Time)
		}
	})

	t.Run("string invalid format", func(t *testing.T) {
		var ts Timestamp
		assertError(t, ts.Scan("not-a-time"), "invalid time format")
	})

	t.Run("[]byte invalid format", func(t *testing.T) {
		var ts Timestamp
		assertError(t, ts.Scan([]byte("not-a-time")), "invalid time format")
	})

	t.Run("invalid type", func(t *testing.T) {
		var ts Timestamp
		assertError(t, ts.Scan(123), "invalid type")
	})
}

func TestTimestamp_Value(t *testing.T) {
	t.Run("zero returns nil", func(t *testing.T) {
		ts := Timestamp{}
		val, err := ts.Value()
		assertNoError(t, err)
		assertNil(t, val)
	})

	t.Run("non-zero returns time.Time", func(t *testing.T) {
		now := time.Now()
		ts := NewTimestamp(now)
		val, err := ts.Value()
		assertNoError(t, err)
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
		assertNoError(t, bs.Scan(nil))
		assertZero(t, bs)
	})

	t.Run("string", func(t *testing.T) {
		var bs BoundedString
		assertNoError(t, bs.Scan("hello world"))
		assertStringEquals(t, bs.String(), "hello world")
	})

	t.Run("string empty", func(t *testing.T) {
		var bs BoundedString
		assertNoError(t, bs.Scan(""))
		assertZero(t, bs)
	})

	t.Run("[]byte", func(t *testing.T) {
		var bs BoundedString
		assertNoError(t, bs.Scan([]byte("hello")))
		assertStringEquals(t, bs.String(), "hello")
	})

	t.Run("[]byte empty", func(t *testing.T) {
		var bs BoundedString
		assertNoError(t, bs.Scan([]byte{}))
		assertZero(t, bs)
	})

	t.Run("invalid type", func(t *testing.T) {
		var bs BoundedString
		assertError(t, bs.Scan(123), "invalid type")
	})

	t.Run("sets bounds correctly", func(t *testing.T) {
		var bs BoundedString
		_ = bs.Scan("hello")
		assertEqual(t, bs.MinLen(), 0)
		assertEqual(t, bs.MaxLen(), 5)
	})
}

func TestBoundedString_Value(t *testing.T) {
	t.Run("empty returns nil", func(t *testing.T) {
		bs := BoundedString{}
		val, err := bs.Value()
		assertNoError(t, err)
		assertNil(t, val)
	})

	t.Run("non-empty returns string", func(t *testing.T) {
		bs := MustBoundedString(0, 100, "hello")
		val, err := bs.Value()
		assertNoError(t, err)
		assertEqual(t, val.(string), "hello")
	})
}

// ================================================================================
// ID[B, V] SQL Tests
// ================================================================================

type testBrand struct{}

func TestID_Scan_StringBrand(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var id ID[testBrand, string]
		assertNoError(t, id.Scan(nil))
		assertZero(t, id)
	})

	t.Run("string", func(t *testing.T) {
		var id ID[testBrand, string]
		assertNoError(t, id.Scan("user-123"))
		assertEqual(t, id.Get(), "user-123")
	})

	t.Run("[]byte", func(t *testing.T) {
		var id ID[testBrand, string]
		assertNoError(t, id.Scan([]byte("user-456")))
		assertEqual(t, id.Get(), "user-456")
	})

	t.Run("invalid type", func(t *testing.T) {
		var id ID[testBrand, string]
		assertError(t, id.Scan(123), "invalid type")
	})
}

func TestID_Scan_Int64Brand(t *testing.T) {
	type intBrand struct{}

	t.Run("nil", func(t *testing.T) {
		var id ID[intBrand, int64]
		assertNoError(t, id.Scan(nil))
		assertZero(t, id)
	})

	t.Run("int64", func(t *testing.T) {
		var id ID[intBrand, int64]
		assertNoError(t, id.Scan(int64(12345)))
		assertEqual(t, id.Get(), int64(12345))
	})

	t.Run("float64", func(t *testing.T) {
		var id ID[intBrand, int64]
		assertNoError(t, id.Scan(float64(12345.67)))
		assertEqual(t, id.Get(), int64(12345))
	})

	t.Run("invalid type", func(t *testing.T) {
		var id ID[intBrand, int64]
		assertError(t, id.Scan("string"), "invalid type")
	})
}

func TestID_Value_StringBrand(t *testing.T) {
	t.Run("zero returns nil", func(t *testing.T) {
		id := ID[testBrand, string]{}
		val, err := id.Value()
		assertNoError(t, err)
		assertNil(t, val)
	})

	t.Run("non-zero returns string", func(t *testing.T) {
		id := NewID[testBrand, string]("user-123")
		val, err := id.Value()
		assertNoError(t, err)
		assertEqual(t, val.(string), "user-123")
	})
}

func TestID_Value_Int64Brand(t *testing.T) {
	type intBrand struct{}

	t.Run("zero returns nil", func(t *testing.T) {
		id := ID[intBrand, int64]{}
		val, err := id.Value()
		assertNoError(t, err)
		assertNil(t, val)
	})

	t.Run("non-zero returns int64", func(t *testing.T) {
		id := NewID[intBrand, int64](int64(12345))
		val, err := id.Value()
		assertNoError(t, err)
		assertEqual(t, val.(int64), int64(12345))
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
		assertNoError(t, l.Scan(nil))
		assertZero(t, l)
	})

	t.Run("string valid", func(t *testing.T) {
		var l Locale
		assertNoError(t, l.Scan("en-US"))
		assertStringEquals(t, l.String(), "en-US")
	})

	t.Run("string empty", func(t *testing.T) {
		var l Locale
		assertNoError(t, l.Scan(""))
		assertZero(t, l)
	})

	t.Run("[]byte valid", func(t *testing.T) {
		var l Locale
		assertNoError(t, l.Scan([]byte("de-DE")))
		assertStringEquals(t, l.String(), "de-DE")
	})

	t.Run("[]byte empty", func(t *testing.T) {
		var l Locale
		assertNoError(t, l.Scan([]byte{}))
		assertZero(t, l)
	})

	t.Run("invalid type", func(t *testing.T) {
		var l Locale
		assertError(t, l.Scan(123), "invalid type")
	})

	t.Run("invalid locale string", func(t *testing.T) {
		var l Locale
		assertError(t, l.Scan("invalid-locale-format-!!!"), "invalid locale")
	})
}

func TestLocale_Value(t *testing.T) {
	t.Run("zero returns nil", func(t *testing.T) {
		l := Locale{}
		val, err := l.Value()
		assertNoError(t, err)
		assertNil(t, val)
	})

	t.Run("non-zero returns string", func(t *testing.T) {
		l := MustParseLocale("en-US")
		val, err := l.Value()
		assertNoError(t, err)
		assertEqual(t, val.(string), "en-US")
	})
}

// ================================================================================
// Duration SQL Tests
// ================================================================================

func TestDuration_Scan(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var d Duration
		assertNoError(t, d.Scan(nil))
		assertZero(t, d)
	})

	t.Run("int64 nanoseconds", func(t *testing.T) {
		var d Duration
		assertNoError(t, d.Scan(int64(5*time.Second)))
		assertEqual(t, d.Duration, 5*time.Second)
	})

	t.Run("float64 nanoseconds", func(t *testing.T) {
		var d Duration
		assertNoError(t, d.Scan(float64(3*time.Second)))
		assertEqual(t, d.Duration, 3*time.Second)
	})

	t.Run("string valid", func(t *testing.T) {
		var d Duration
		assertNoError(t, d.Scan("1h30m"))
		assertEqual(t, d.Duration, 90*time.Minute)
	})

	t.Run("string empty", func(t *testing.T) {
		var d Duration
		assertNoError(t, d.Scan(""))
		assertZero(t, d)
	})

	t.Run("[]byte valid", func(t *testing.T) {
		var d Duration
		assertNoError(t, d.Scan([]byte("500ms")))
		assertEqual(t, d.Duration, 500*time.Millisecond)
	})

	t.Run("[]byte empty", func(t *testing.T) {
		var d Duration
		assertNoError(t, d.Scan([]byte{}))
		assertZero(t, d)
	})

	t.Run("invalid type", func(t *testing.T) {
		var d Duration
		assertError(t, d.Scan("invalid"), "invalid duration string")
	})

	t.Run("invalid source type", func(t *testing.T) {
		var d Duration
		assertError(t, d.Scan([]int{1, 2, 3}), "invalid type")
	})
}

func TestDuration_Value(t *testing.T) {
	t.Run("zero returns nil", func(t *testing.T) {
		d := Duration{}
		val, err := d.Value()
		assertNoError(t, err)
		assertNil(t, val)
	})

	t.Run("non-zero returns int64 nanoseconds", func(t *testing.T) {
		d := NewDuration(5 * time.Second)
		val, err := d.Value()
		assertNoError(t, err)
		assertEqual(t, val.(int64), int64(5*time.Second))
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
