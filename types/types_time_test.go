package types

import (
	"testing"
	"time"
)

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

func TestTimestampBeforeAfter(t *testing.T) {
	t.Parallel()

	earlier := NewTimestamp(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	later := NewTimestamp(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC))

	if !earlier.Before(later.Time) {
		t.Error("earlier should be before later")
	}

	if earlier.Before(earlier.Time) {
		t.Error("same time should not be before")
	}

	if !later.After(earlier.Time) {
		t.Error("later should be after earlier")
	}

	if later.After(later.Time) {
		t.Error("same time should not be after")
	}
}

func TestDurationUnmarshalJSONErrors(t *testing.T) {
	t.Parallel()

	t.Run("non-string JSON", func(t *testing.T) {
		t.Parallel()

		var dur Duration

		err := dur.UnmarshalJSON([]byte(`123`))
		if err == nil {
			t.Error("expected error for non-string JSON")
		}
	})

	t.Run("invalid duration string", func(t *testing.T) {
		t.Parallel()

		var dur Duration

		err := dur.UnmarshalJSON([]byte(`"not-a-duration"`))
		if err == nil {
			t.Error("expected error for invalid duration string")
		}
	})
}
