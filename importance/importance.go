// Package importance provides a 0-100 priority classification with named levels.
//
// Importance wraps a uint8 clamped to 0-100 with named constants (None, VeryLow, Low,
// Medium, High, VeryHigh, Max) and convenience predicates. It supports JSON, SQL, and
// string parsing.
//
// Basic usage:
//
//	imp := importance.New(75)
//	fmt.Println(imp.Classification()) // "High"
//	fmt.Println(imp.Percent())        // 0.75
package importance

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/larsartmann/go-composable-business-types/scanutil"
	"github.com/larsartmann/go-composable-business-types/validate"
)

const (
	maxValue = 100
)

// String representations for parsing and output.
const (
	strNone     = "none"
	strVeryLow  = "very-low"
	strLow      = "low"
	strMedium   = "medium"
	strHigh     = "high"
	strVeryHigh = "very-high"
)

var (
	errExceedsMax      = errors.New("importance: value exceeds maximum")
	errUnknownClass    = errors.New("importance: unknown classification")
	errScanNilReceiver = errors.New("importance: scan: receiver is nil")
)

// Importance represents a priority level from 0 to 100.
type Importance uint8

// Named importance levels.
const (
	None     Importance = 0
	VeryLow  Importance = 20
	Low      Importance = 40
	Medium   Importance = 50
	High     Importance = 70
	VeryHigh Importance = 90
	Max      Importance = 100
)

// New creates an Importance from a uint8 value, returning an error if v exceeds 100.
func New(v uint8) (Importance, error) {
	if v > maxValue {
		return 0, fmt.Errorf("%w: %d > %d", errExceedsMax, v, maxValue)
	}

	return Importance(v), nil
}

// Must creates an Importance or panics if v exceeds 100.
func Must(v uint8) Importance {
	i, err := New(v)
	if err != nil {
		panic(err)
	}

	return i
}

// Parse parses a named importance level string.
func Parse(s string) (Importance, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case strNone, "0":
		return None, nil
	case strVeryLow, "verylow":
		return VeryLow, nil
	case strLow:
		return Low, nil
	case strMedium, "mid":
		return Medium, nil
	case strHigh:
		return High, nil
	case strVeryHigh, "veryhigh":
		return VeryHigh, nil
	case "max", "100":
		return Max, nil
	default:
		return 0, fmt.Errorf("%w: %q", errUnknownClass, s)
	}
}

// String returns the lowercase kebab-case importance level name.
func (i Importance) String() string {
	switch {
	case i == 0:
		return strNone
	case i <= 20:
		return strVeryLow
	case i <= 40:
		return strLow
	case i <= 60:
		return strMedium
	case i <= 80:
		return strHigh
	default:
		return strVeryHigh
	}
}

// Classification returns the human-readable importance level name.
func (i Importance) Classification() string {
	switch {
	case i == 0:
		return "None"
	case i <= 20:
		return "Very Low"
	case i <= 40:
		return "Low"
	case i <= 60:
		return "Medium"
	case i <= 80:
		return "High"
	default:
		return "Very High"
	}
}

// IsVeryLow reports whether the importance is in the VeryLow range.
func (i Importance) IsVeryLow() bool { return i > 0 && i <= 20 }

// IsLow reports whether the importance is in the Low range.
func (i Importance) IsLow() bool { return i >= 21 && i <= 40 }

// IsMedium reports whether the importance is in the Medium range.
func (i Importance) IsMedium() bool { return i >= 41 && i <= 60 }

// IsHigh reports whether the importance is in the High range.
func (i Importance) IsHigh() bool { return i >= 61 && i <= 80 }

// IsVeryHigh reports whether the importance is in the VeryHigh range.
func (i Importance) IsVeryHigh() bool { return i >= 81 }

// IsValid reports whether the importance is within the valid 0-100 range.
func (i Importance) IsValid() bool { return i <= maxValue }

// IsZero reports whether the importance is zero.
func (i Importance) IsZero() bool { return i == 0 }

// Percent returns the importance as a float64 between 0.0 and 1.0.
func (i Importance) Percent() float64 { return float64(i) / 100.0 }

// Compare returns -1, 0, or +1 depending on comparison.
func (i Importance) Compare(other Importance) int {
	if i < other {
		return -1
	}

	if i > other {
		return 1
	}

	return 0
}

// IsNone reports whether the importance equals None.
func (i Importance) IsNone() bool { return i == None }

// IsDefault reports whether the importance equals Medium.
func (i Importance) IsDefault() bool { return i == Medium }

// Clamp returns i clamped to the valid range.
func (i Importance) Clamp() Importance {
	if i < None {
		return None
	}

	if i > Max {
		return Max
	}

	return i
}

// Add returns the sum of i and other, clamped to Max.
func (i Importance) Add(other Importance) Importance {
	result := i + other
	if result > Max || result < i {
		return Max
	}

	return result
}

// Sub returns i minus other, clamped to None.
func (i Importance) Sub(other Importance) Importance {
	if other > i {
		return None
	}

	return i - other
}

// Max returns the larger of i and other.
func (i Importance) Max(other Importance) Importance {
	if i > other {
		return i
	}

	return other
}

// Min returns the smaller of i and other.
func (i Importance) Min(other Importance) Importance {
	if i < other {
		return i
	}

	return other
}

// PercentString returns the importance formatted as a percentage string.
func (i Importance) PercentString() string {
	return fmt.Sprintf("%d%%", i)
}

// Validate returns an error if the importance exceeds 100.
func (i Importance) Validate() error {
	if i > maxValue {
		return fmt.Errorf("%w: %d > %d", errExceedsMax, i, maxValue)
	}

	return nil
}

// MarshalJSON encodes the importance as a JSON number.
func (i Importance) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(uint8(i))
	if err != nil {
		return nil, fmt.Errorf("importance: marshal JSON: %w", err)
	}

	return b, nil
}

// UnmarshalJSON decodes a JSON number into the importance.
func (i *Importance) UnmarshalJSON(data []byte) error {
	var v uint8

	err := json.Unmarshal(data, &v)
	if err != nil {
		return fmt.Errorf("importance: invalid JSON %q: %w", string(data), err)
	}

	*i = Importance(v)

	return nil
}

// Scan implements sql.Scanner for the importance.
func (i *Importance) Scan(src any) error {
	if i == nil {
		return errScanNilReceiver
	}

	return scanutil.ScanInt64(src, func(v int64) error {
		*i = Importance(v) //nolint:gosec // G115: int64 to uint8 for Importance (0-100 range)

		return nil
	})
}

// Value implements driver.Valuer for the importance.
func (i Importance) Value() (driver.Value, error) {
	return scanutil.ZeroAsNullValue(int64(i))
}

var _ validate.Validator = Importance(0)
