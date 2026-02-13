package cbt

import (
	"encoding/json"
	"time"
)

// Bitemporal captures both valid time (when the fact was true in the real world)
// and transaction time (when we recorded it in the system).
// This enables point-in-time queries and corrections of historical data.
type Bitemporal struct {
	validFrom  Timestamp // When this fact became true in the real world
	validUntil Timestamp // When this fact ceased to be true (zero = still valid)
	recorded   Timestamp // When we recorded this fact in the system
	correction bool      // Is this a correction of a previous record?
}

// NewBitemporal creates a new Bitemporal with valid time starting now.
// The fact is considered valid indefinitely (validUntil is zero).
func NewBitemporal(recorded Timestamp) Bitemporal {
	return Bitemporal{
		validFrom:  recorded,
		validUntil: Timestamp{},
		recorded:   recorded,
		correction: false,
	}
}

// NewBitemporalWithRange creates a Bitemporal with explicit valid time range.
// If validUntil is zero, the fact is valid indefinitely.
func NewBitemporalWithRange(validFrom, validUntil, recorded Timestamp) Bitemporal {
	return Bitemporal{
		validFrom:  validFrom,
		validUntil: validUntil,
		recorded:   recorded,
		correction: false,
	}
}

// NewCorrection creates a Bitemporal that marks this as a correction.
func NewCorrection(validFrom, validUntil, recorded Timestamp) Bitemporal {
	return Bitemporal{
		validFrom:  validFrom,
		validUntil: validUntil,
		recorded:   recorded,
		correction: true,
	}
}

// ValidFrom returns when this fact became true in the real world.
func (b Bitemporal) ValidFrom() Timestamp { return b.validFrom }

// ValidUntil returns when this fact ceased to be true (zero = still valid).
func (b Bitemporal) ValidUntil() Timestamp { return b.validUntil }

// Recorded returns when this fact was recorded in the system.
func (b Bitemporal) Recorded() Timestamp { return b.recorded }

// IsCorrection returns true if this is a correction of a previous record.
func (b Bitemporal) IsCorrection() bool { return b.correction }

// IsValidAt checks if this fact was valid at the given point in time.
func (b Bitemporal) IsValidAt(t Timestamp) bool {
	if t.Before(b.validFrom.Time) {
		return false
	}
	if b.validUntil.IsZero() {
		return true
	}
	return t.Before(b.validUntil.Time) || t.Equal(b.validUntil.Time)
}

// IsCurrentlyValid checks if this fact is currently valid.
func (b Bitemporal) IsCurrentlyValid() bool {
	return b.IsValidAt(Now())
}

// WithValidUntil returns a copy with validUntil set.
func (b Bitemporal) WithValidUntil(until Timestamp) Bitemporal {
	b.validUntil = until
	return b
}

// jsonBitemporal is the JSON representation of Bitemporal.
type jsonBitemporal struct {
	ValidFrom  time.Time `json:"validFrom"`
	ValidUntil time.Time `json:"validUntil"`
	Recorded   time.Time `json:"recorded"`
	Correction bool      `json:"correction,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (b Bitemporal) MarshalJSON() ([]byte, error) {
	return json.Marshal(jsonBitemporal{
		ValidFrom:  b.validFrom.Time,
		ValidUntil: b.validUntil.Time,
		Recorded:   b.recorded.Time,
		Correction: b.correction,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *Bitemporal) UnmarshalJSON(data []byte) error {
	var raw jsonBitemporal
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	b.validFrom = NewTimestamp(raw.ValidFrom)
	b.validUntil = NewTimestamp(raw.ValidUntil)
	b.recorded = NewTimestamp(raw.Recorded)
	b.correction = raw.Correction
	return nil
}
