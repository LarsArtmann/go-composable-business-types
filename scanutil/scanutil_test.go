package scanutil

import (
	"testing"
)

func TestScanString(t *testing.T) {
	tests := []struct {
		name      string
		src       any
		wantValue string
		wantErr   bool
	}{
		{"nil", nil, "", false},
		{"empty string", "", "", false},
		{"non-empty string", "hello", "hello", false},
		{"byte slice", []byte("world"), "world", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string
			err := ScanString(tt.src, func(v string) error {
				got = v
				return nil
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("ScanString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantValue {
				t.Errorf("ScanString() got = %v, want %v", got, tt.wantValue)
			}
		})
	}
}

func TestScanString_InvalidType(t *testing.T) {
	err := ScanString(12345, func(v string) error {
		t.Error("callback should not be called for invalid type")
		return nil
	})
	if err == nil {
		t.Error("expected error for invalid type")
	}
}

func TestScanInt64(t *testing.T) {
	tests := []struct {
		name      string
		src       any
		wantValue int64
		wantErr   bool
	}{
		{"nil", nil, 0, false},
		{"int64", int64(42), 42, false},
		{"int", int(100), 100, false},
		{"float64", float64(200.9), 200, false},
		{"byte slice", []byte("999"), 999, false},
		{"empty byte slice", []byte{}, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got int64
			err := ScanInt64(tt.src, func(v int64) error {
				got = v
				return nil
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("ScanInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantValue {
				t.Errorf("ScanInt64() got = %v, want %v", got, tt.wantValue)
			}
		})
	}
}

func TestScanInt64_InvalidType(t *testing.T) {
	err := ScanInt64("not-a-number", func(v int64) error {
		t.Error("callback should not be called for invalid type")
		return nil
	})
	if err == nil {
		t.Error("expected error for invalid type")
	}
}

func TestNullableValue(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantNil bool
		wantVal string
	}{
		{"empty string", "", true, ""},
		{"non-empty string", "hello", false, "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NullableValue(tt.input)
			if err != nil {
				t.Errorf("NullableValue() error = %v", err)
				return
			}
			if tt.wantNil && got != nil {
				t.Errorf("NullableValue() got nil = false, want nil = true")
			}
			if !tt.wantNil && got != tt.wantVal {
				t.Errorf("NullableValue() got = %v, want %v", got, tt.wantVal)
			}
		})
	}
}

func TestInt64Value(t *testing.T) {
	got, err := Int64Value(12345)
	if err != nil {
		t.Errorf("Int64Value() error = %v", err)
	}
	if got != int64(12345) {
		t.Errorf("Int64Value() got = %v, want 12345", got)
	}
}

func TestZeroAsNullValue(t *testing.T) {
	tests := []struct {
		name    string
		input   int64
		wantNil bool
	}{
		{"zero", 0, true},
		{"non-zero", 42, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ZeroAsNullValue(tt.input)
			if err != nil {
				t.Errorf("ZeroAsNullValue() error = %v", err)
				return
			}
			if tt.wantNil && got != nil {
				t.Errorf("ZeroAsNullValue() got nil = false, want nil = true")
			}
			if !tt.wantNil && got != tt.input {
				t.Errorf("ZeroAsNullValue() got = %v, want %v", got, tt.input)
			}
		})
	}
}
