package types

import (
	"database/sql/driver"
	"fmt"
)

// Percentage represents a value from 0-100.
type Percentage uint8

// Percentage constants for common values.
const (
	percentageDivisor  = 100 // Used for float64 conversion
	percentageMaxValue = 100 // Maximum valid percentage
)

// NewPercentage creates a new Percentage, clamping values above 100 to 100.
func NewPercentage(v uint8) Percentage {
	if v > percentageMaxValue {
		return percentageMaxValue
	}

	return Percentage(v)
}

// Float64 returns the percentage as a float64 between 0 and 1.
func (p Percentage) Float64() float64 { return float64(p) / percentageDivisor } //nolint:divbyzero // percentageDivisor is a constant 100, never zero

// String returns the percentage as a human-readable string (e.g., "50%").
func (p Percentage) String() string { return fmt.Sprintf("%d%%", p) }

// IsZero returns true if the percentage is 0.
func (p Percentage) IsZero() bool { return p == 0 }

// IsMin returns true if the percentage is 0 (minimum value).
func (p Percentage) IsMin() bool { return p == 0 }

// IsMax returns true if the percentage is 100 (maximum value).
func (p Percentage) IsMax() bool { return p == percentageMaxValue }

// CompareOrdered defines the interface for types that support comparison.
type CompareOrdered interface {
	~int8 | ~int16 | ~int32 | ~int64 | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64
}

// compare compares two ordered values. Returns -1 if a < b, 0 if equal, 1 if a > b.
func compare[T CompareOrdered](a, b T) int {
	if a < b {
		return -1
	}

	if a > b {
		return 1
	}

	return 0
}

// Compare returns -1 if p < other, 0 if equal, 1 if p > other.
func (p Percentage) Compare(other Percentage) int {
	return compare(p, other)
}

// MarshalJSON implements json.Marshaler.
func (p Percentage) MarshalJSON() ([]byte, error) {
	return MarshalJSON("percentage", uint8(p))
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *Percentage) UnmarshalJSON(data []byte) error {
	var v uint8

	err := UnmarshalJSON("percentage", data, &v)
	if err != nil {
		return err
	}

	*p = Percentage(v)

	return nil
}

// Scan implements sql.Scanner for Percentage.
// Supports int64 and uint8 sources.
func (p *Percentage) Scan(src any) error {
	return scanInt64Type(p, "percentage", src, func(v int64) Percentage {
		return Percentage(v) //nolint:gosec // G115: int64 to uint8 for Percentage (0-100 range)
	})
}

// Value implements driver.Valuer for Percentage.
func (p Percentage) Value() (driver.Value, error) {
	return valueInt64Type(int64(p), "percentage")
}

// Validate implements validate.Validator for Percentage.
func (p Percentage) Validate() error { return nil }

// Cents represents monetary amounts in smallest currency unit (prevents float errors).
type Cents int64

// Cents conversion constant.
const centsDivisor = 100 // Used for float64 conversion

// NewCents creates a new Cents value from an int64.
func NewCents(v int64) Cents { return Cents(v) }

// Int64 returns the cents value as an int64.
func (c Cents) Int64() int64 { return int64(c) }

// Float64 returns the cents value as a float64 (divided by 100).
func (c Cents) Float64() float64 { return float64(c) / centsDivisor } //nolint:divbyzero // centsDivisor is a constant 100, never zero

// Validate implements validate.Validator for Cents.
// Cents are always valid since they can represent any integer value.
func (c Cents) Validate() error { return nil }

// Add returns the sum of two Cents values.
func (c Cents) Add(other Cents) Cents { return c + other }

// Sub returns the difference of two Cents values.
func (c Cents) Sub(other Cents) Cents { return c - other }

// Mul multiplies Cents by an int64 factor.
func (c Cents) Mul(n int64) Cents { return c * Cents(n) }

// Div divides Cents by an int64 divisor.
func (c Cents) Div(n int64) Cents { return c / Cents(n) }

// Abs returns the absolute value of Cents.
func (c Cents) Abs() Cents {
	if c < 0 {
		return -c
	}

	return c
}

// Sign returns -1 for negative, 0 for zero, 1 for positive.
func (c Cents) Sign() int {
	if c < 0 {
		return -1
	}

	if c > 0 {
		return 1
	}

	return 0
}

// String returns the cents as a human-readable currency string (e.g., "$12.34").
func (c Cents) String() string { return fmt.Sprintf("$%.2f", float64(c)/centsDivisor) } //nolint:divbyzero // centsDivisor is a constant 100, never zero

// IsZero returns true if Cents is zero.
func (c Cents) IsZero() bool { return c == 0 }

// IsPositive returns true if Cents is greater than zero.
func (c Cents) IsPositive() bool { return c > 0 }

// IsNegative returns true if Cents is less than zero.
func (c Cents) IsNegative() bool { return c < 0 }

// Compare returns -1 if c < other, 0 if equal, 1 if c > other.
func (c Cents) Compare(other Cents) int {
	return compare(c, other)
}
