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

type Importance uint8

const (
	None     Importance = 0
	VeryLow  Importance = 20
	Low      Importance = 40
	Medium   Importance = 50
	High     Importance = 70
	VeryHigh Importance = 90
	Max      Importance = 100
)

func New(v uint8) (Importance, error) {
	if v > maxValue {
		return 0, fmt.Errorf("importance: value %d exceeds maximum %d", v, maxValue)
	}

	return Importance(v), nil
}

func Must(v uint8) Importance {
	i, err := New(v)
	if err != nil {
		panic(err)
	}

	return i
}

func Parse(s string) (Importance, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "none", "0":
		return None, nil
	case "very-low", "verylow":
		return VeryLow, nil
	case "low":
		return Low, nil
	case "medium", "mid":
		return Medium, nil
	case "high":
		return High, nil
	case "very-high", "veryhigh":
		return VeryHigh, nil
	case "max", "100":
		return Max, nil
	default:
		return 0, fmt.Errorf("importance: unknown classification %q", s)
	}
}

func (i Importance) String() string {
	switch {
	case i == 0:
		return "none"
	case i <= 20:
		return "very-low"
	case i <= 40:
		return "low"
	case i <= 60:
		return "medium"
	case i <= 80:
		return "high"
	default:
		return "very-high"
	}
}

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

func (i Importance) IsVeryLow() bool  { return i > 0 && i <= 20 }
func (i Importance) IsLow() bool      { return i >= 21 && i <= 40 }
func (i Importance) IsMedium() bool   { return i >= 41 && i <= 60 }
func (i Importance) IsHigh() bool     { return i >= 61 && i <= 80 }
func (i Importance) IsVeryHigh() bool { return i >= 81 }

func (i Importance) IsValid() bool    { return i <= maxValue }
func (i Importance) IsZero() bool     { return i == 0 }
func (i Importance) Percent() float64 { return float64(i) / 100.0 }

func (i Importance) Compare(other Importance) int {
	if i < other {
		return -1
	}

	if i > other {
		return 1
	}

	return 0
}

func (i Importance) IsNone() bool    { return i == None }
func (i Importance) IsDefault() bool { return i == Medium }

func (i Importance) Clamp() Importance {
	if i < None {
		return None
	}

	if i > Max {
		return Max
	}

	return i
}

func (i Importance) Add(other Importance) Importance {
	result := i + other
	if result > Max || result < i {
		return Max
	}

	return result
}

func (i Importance) Sub(other Importance) Importance {
	if other > i {
		return None
	}

	return i - other
}

func (i Importance) Max(other Importance) Importance {
	if i > other {
		return i
	}

	return other
}

func (i Importance) Min(other Importance) Importance {
	if i < other {
		return i
	}

	return other
}

func (i Importance) PercentString() string {
	return fmt.Sprintf("%d%%", i)
}

func (i Importance) Validate() error {
	if i > maxValue {
		return fmt.Errorf("importance: value %d exceeds maximum %d", i, maxValue)
	}

	return nil
}

func (i Importance) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(uint8(i))
	if err != nil {
		return nil, fmt.Errorf("importance: marshal JSON: %w", err)
	}

	return b, nil
}

func (i *Importance) UnmarshalJSON(data []byte) error {
	var v uint8

	err := json.Unmarshal(data, &v)
	if err != nil {
		return fmt.Errorf("importance: invalid JSON %q: %w", string(data), err)
	}

	*i = Importance(v)

	return nil
}

func (i *Importance) Scan(src any) error {
	if i == nil {
		return errors.New("importance: scan: receiver is nil")
	}

	return scanutil.ScanInt64(src, func(v int64) error {
		*i = Importance(v) //nolint:gosec // G115: int64 to uint8 for Importance (0-100 range)

		return nil
	})
}

func (i Importance) Value() (driver.Value, error) {
	return scanutil.ZeroAsNullValue(int64(i))
}

var _ validate.Validator = Importance(0)
