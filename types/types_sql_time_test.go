package types

import (
	"testing"
	"time"
)

// SQL Scanner/Valuer tests for Timestamp and Duration
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
