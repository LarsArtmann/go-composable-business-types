// Package temporal provides bitemporal time tracking.
//
// Bitemporal captures both valid time (when the fact was true in the real world)
// and transaction time (when it was recorded in the system). This enables
// point-in-time queries and corrections of historical data.
//
// Basic usage:
//
//	from := types.NewTimestamp(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
//	until := types.NewTimestamp(time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC))
//	recorded := types.Now()
//	temp := temporal.NewBitemporalWithRange(from, until, recorded)
package temporal

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/larsartmann/go-composable-business-types/types"
)

// Correction represents whether a Bitemporal record is a correction.
// This type-safe approach provides named constants while maintaining JSON compatibility.
type Correction bool

const (
	// NoCorrection indicates this is not a correction record.
	NoCorrection Correction = false
	// IsCorrection indicates this is a correction of a previous record.
	IsCorrection Correction = true
)

// String returns "correction" for true, empty string for false.
func (c Correction) String() string {
	if c {
		return "correction"
	}

	return ""
}

// MarshalJSON implements json.Marshaler for Correction.
func (c Correction) MarshalJSON() ([]byte, error) {
	return types.MarshalJSON("correction", bool(c))
}

// UnmarshalJSON implements json.Unmarshaler for Correction.
func (c *Correction) UnmarshalJSON(data []byte) error {
	var b bool

	err := types.UnmarshalJSON("correction", data, &b)
	if err != nil {
		return err
	}

	*c = Correction(b)

	return nil
}

// Bitemporal captures both valid time (when the fact was true in the real world)
// and transaction time (when we recorded it in the system).
// This enables point-in-time queries and corrections of historical data.
type Bitemporal struct {
	validFrom  types.Timestamp // When this fact became true in the real world
	validUntil types.Timestamp // When this fact ceased to be true (zero = still valid)
	recorded   types.Timestamp // When we recorded this fact in the system
	correction Correction      // Is this a correction of a previous record?
}

// NewBitemporal creates a new Bitemporal with valid time starting now.
// The fact is considered valid indefinitely (validUntil is zero).
func NewBitemporal(recorded types.Timestamp) Bitemporal {
	return Bitemporal{
		validFrom:  recorded,
		validUntil: types.Timestamp{Time: time.Time{}},
		recorded:   recorded,
		correction: NoCorrection,
	}
}

// NewBitemporalWithRange creates a Bitemporal with explicit valid time range.
// If validUntil is zero, the fact is valid indefinitely.
func NewBitemporalWithRange(validFrom, validUntil, recorded types.Timestamp) Bitemporal {
	return NewCorrection(validFrom, validUntil, recorded).withCorrection(NoCorrection)
}

// NewCorrection creates a Bitemporal that marks this as a correction.
func NewCorrection(validFrom, validUntil, recorded types.Timestamp) Bitemporal {
	return Bitemporal{
		validFrom:  validFrom,
		validUntil: validUntil,
		recorded:   recorded,
		correction: IsCorrection,
	}
}

// withCorrection returns a copy with the correction flag set.
func (b Bitemporal) withCorrection(c Correction) Bitemporal {
	b.correction = c

	return b
}

// ValidFrom returns when this fact became true in the real world.
func (b Bitemporal) ValidFrom() types.Timestamp { return b.validFrom }

// ValidUntil returns when this fact ceased to be true (zero = still valid).
func (b Bitemporal) ValidUntil() types.Timestamp { return b.validUntil }

// Recorded returns when this fact was recorded in the system.
func (b Bitemporal) Recorded() types.Timestamp { return b.recorded }

// IsCorrection returns the Correction type indicating if this is a correction.
func (b Bitemporal) IsCorrection() Correction { return b.correction }

// IsZero returns true if this is the zero value (all timestamps zero).
func (b Bitemporal) IsZero() bool {
	return b.validFrom.IsZero() && b.validUntil.IsZero() && b.recorded.IsZero()
}

// IsValidAt checks if this fact was valid at the given point in time.
func (b Bitemporal) IsValidAt(t types.Timestamp) bool {
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
	return b.IsValidAt(types.Now())
}

// WithValidUntil returns a copy with validUntil set.
func (b Bitemporal) WithValidUntil(until types.Timestamp) Bitemporal {
	b.validUntil = until

	return b
}

// jsonBitemporal is the JSON representation of Bitemporal.
type jsonBitemporal struct {
	ValidFrom  time.Time  `json:"validFrom"`
	ValidUntil time.Time  `json:"validUntil"`
	Recorded   time.Time  `json:"recorded"`
	Correction Correction `json:"correction,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (b Bitemporal) MarshalJSON() ([]byte, error) {
	bj, err := json.Marshal(jsonBitemporal{
		ValidFrom:  b.validFrom.Time,
		ValidUntil: b.validUntil.Time,
		Recorded:   b.recorded.Time,
		Correction: b.correction,
	})
	if err != nil {
		return nil, fmt.Errorf("bitemporal: marshal JSON: %w", err)
	}

	return bj, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *Bitemporal) UnmarshalJSON(data []byte) error {
	var raw jsonBitemporal

	err := json.Unmarshal(data, &raw)
	if err != nil {
		return fmt.Errorf("unmarshal bitemporal: invalid JSON %q: %w", string(data), err)
	}

	b.validFrom = types.NewTimestamp(raw.ValidFrom)
	b.validUntil = types.NewTimestamp(raw.ValidUntil)
	b.recorded = types.NewTimestamp(raw.Recorded)
	b.correction = raw.Correction

	return nil
}
