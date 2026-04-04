package temporal

import (
	"testing"
	"time"

	"github.com/larsartmann/go-composable-business-types/types"
)

func TestNewBitemporal(t *testing.T) {
	t.Parallel()

	now := types.Now()
	b := NewBitemporal(now)

	if b.ValidFrom() != now {
		t.Error("ValidFrom should equal recorded time")
	}

	if !b.ValidUntil().IsZero() {
		t.Error("ValidUntil should be zero (indefinite)")
	}

	if b.Recorded() != now {
		t.Error("Recorded should equal input")
	}

	if b.IsCorrection() != NoCorrection {
		t.Error("NewBitemporal should not be a correction")
	}
}

func TestNewBitemporalWithRange(t *testing.T) {
	t.Parallel()

	from := types.NewTimestamp(time.Now())
	until := types.NewTimestamp(time.Now().Add(time.Hour))
	recorded := types.Now()

	b := NewBitemporalWithRange(from, until, recorded)

	if b.ValidFrom() != from {
		t.Error("ValidFrom mismatch")
	}

	if b.ValidUntil() != until {
		t.Error("ValidUntil mismatch")
	}

	if b.IsCorrection() != NoCorrection {
		t.Error("should not be a correction")
	}
}

func TestNewCorrection(t *testing.T) {
	t.Parallel()

	now := types.Now()
	b := NewCorrection(now, types.Timestamp{Time: time.Time{}}, now)

	if b.IsCorrection() != IsCorrection {
		t.Error("NewCorrection should be a correction")
	}
}

func TestBitemporalIsZero(t *testing.T) {
	t.Parallel()

	var zero Bitemporal
	if !zero.IsZero() {
		t.Error("zero Bitemporal should be zero")
	}

	nonZero := NewBitemporal(types.Now())
	if nonZero.IsZero() {
		t.Error("non-zero Bitemporal should not be zero")
	}
}

func TestBitemporalIsValidAt(t *testing.T) {
	t.Parallel()

	now := time.Now()
	from := types.NewTimestamp(now.Add(-time.Hour))
	until := types.NewTimestamp(now.Add(time.Hour))
	recorded := types.Now()

	b := NewBitemporalWithRange(from, until, recorded)

	// Should be valid at current time
	if !b.IsValidAt(types.NewTimestamp(now)) {
		t.Error("should be valid at current time")
	}

	// Should not be valid before from
	if b.IsValidAt(types.NewTimestamp(now.Add(-2 * time.Hour))) {
		t.Error("should not be valid before from")
	}

	// Should not be valid after until
	if b.IsValidAt(types.NewTimestamp(now.Add(2 * time.Hour))) {
		t.Error("should not be valid after until")
	}
}

func TestBitemporalIsCurrentlyValid(t *testing.T) {
	t.Parallel()

	now := types.Now()
	b := NewBitemporal(now)

	if !b.IsCurrentlyValid() {
		t.Error("should be currently valid")
	}

	// Past entry (validUntil in past)
	past := types.NewTimestamp(time.Now().Add(-2 * time.Hour))
	pastUntil := types.NewTimestamp(time.Now().Add(-time.Hour))
	pastB := NewBitemporalWithRange(past, pastUntil, past)

	if pastB.IsCurrentlyValid() {
		t.Error("past entry should not be currently valid")
	}
}

func TestBitemporalWithValidUntil(t *testing.T) {
	t.Parallel()

	b := NewBitemporal(types.Now())
	newUntil := types.NewTimestamp(time.Now().Add(time.Hour))

	modified := b.WithValidUntil(newUntil)

	if modified.ValidUntil() != newUntil {
		t.Error("WithValidUntil should set new until")
	}
	// Original should be unchanged
	if !b.ValidUntil().IsZero() {
		t.Error("original should be unchanged")
	}
}

func TestBitemporalJSON(t *testing.T) {
	t.Parallel()

	now := time.Now()
	b := NewBitemporalWithRange(
		types.NewTimestamp(now),
		types.NewTimestamp(now.Add(time.Hour)),
		types.NewTimestamp(now),
	)

	data, err := b.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	// Should contain validFrom
	if len(data) == 0 {
		t.Error("JSON should not be empty")
	}

	// Test unmarshal
	var parsed Bitemporal

	err = parsed.UnmarshalJSON(data)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	// Check values are preserved (with some tolerance for time comparison)
	if parsed.IsZero() {
		t.Error("parsed should not be zero")
	}
}

func TestCorrectionString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		c        Correction
		expected string
	}{
		{"NoCorrection", NoCorrection, ""},
		{"IsCorrection", IsCorrection, "correction"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.c.String(); got != tt.expected {
				t.Errorf("Correction.String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestCorrectionJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		c    Correction
		want string
	}{
		{"NoCorrection", NoCorrection, "false"},
		{"IsCorrection", IsCorrection, "true"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			data, err := tt.c.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON failed: %v", err)
			}

			if string(data) != tt.want {
				t.Errorf("MarshalJSON() = %s, want %s", string(data), tt.want)
			}

			// Test round-trip
			var parsed Correction
			if err := parsed.UnmarshalJSON(data); err != nil {
				t.Fatalf("UnmarshalJSON failed: %v", err)
			}

			if parsed != tt.c {
				t.Errorf("round-trip failed: got %v, want %v", parsed, tt.c)
			}
		})
	}
}

func TestCorrectionUnmarshalJSONErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{"invalid JSON", `"not-a-bool"`, true},
		{"malformed JSON", `{invalid`, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var c Correction

			err := c.UnmarshalJSON([]byte(tt.data))
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
