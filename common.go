package cbt

import "time"

// Email represents a validated email address.
type Email string

func NewEmail(v string) Email { return Email(v) }
func (e Email) String() string { return string(e) }

// URL represents a validated URL.
type URL string

func NewURL(v string) URL { return URL(v) }
func (u URL) String() string { return string(u) }

// Percentage represents a value from 0-100.
type Percentage uint8

func NewPercentage(v uint8) Percentage {
	if v > 100 {
		return 100
	}
	return Percentage(v)
}
func (p Percentage) Float64() float64 { return float64(p) / 100 }

// Cents represents monetary amounts in smallest currency unit (prevents float errors).
type Cents int64

func NewCents(v int64) Cents   { return Cents(v) }
func (c Cents) Int64() int64   { return int64(c) }
func (c Cents) Float64() float64 { return float64(c) / 100 }

// Timestamp wraps time.Time for domain clarity.
type Timestamp struct{ time.Time }

func NewTimestamp(t time.Time) Timestamp { return Timestamp{Time: t} }
func Now() Timestamp                     { return Timestamp{Time: time.Now()} }

// Duration wraps time.Duration for domain clarity.
type Duration struct{ time.Duration }

func NewDuration(d time.Duration) Duration { return Duration{Duration: d} }
