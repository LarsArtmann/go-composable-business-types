package scanutil

import (
	"testing"
)

func TestScanString(t *testing.T) {
	t.Parallel()

	tests := []scanTestCase[string]{
		{name: "nil", src: nil, wantValue: "", wantErr: false},
		{name: "empty string", src: "", wantValue: "", wantErr: false},
		{name: "non-empty string", src: "hello", wantValue: "hello", wantErr: false},
		{name: "byte slice", src: []byte("world"), wantValue: "world", wantErr: false},
	}

	scanTests(t, "ScanString", tests, ScanString)
}

func TestScanInvalidType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		typeName   string
		invalidSrc any
		scanFn     func(fn func() error) error
	}{
		{
			typeName:   "string",
			invalidSrc: 12345,
			scanFn: func(fn func() error) error {
				return ScanString(fn, func(_ string) error { return nil })
			},
		},
		{
			typeName:   "int64",
			invalidSrc: "not-a-number",
			scanFn: func(fn func() error) error {
				return ScanInt64(fn, func(_ int64) error { return nil })
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.typeName, func(t *testing.T) {
			t.Parallel()
			scanInvalidType(t, tt.typeName, tt.invalidSrc, tt.scanFn)
		})
	}
}

func TestScanInt64(t *testing.T) {
	t.Parallel()

	tests := []scanTestCase[int64]{
		{name: "nil", src: nil, wantValue: 0, wantErr: false},
		{name: "int64", src: int64(42), wantValue: 42, wantErr: false},
		{name: "int", src: int(100), wantValue: 100, wantErr: false},
		{name: "float64", src: float64(200.9), wantValue: 200, wantErr: false},
		{name: "byte slice", src: []byte("999"), wantValue: 999, wantErr: false},
		{name: "empty byte slice", src: []byte{}, wantValue: 0, wantErr: false},
	}

	scanTests(t, "ScanInt64", tests, ScanInt64)
}

type scanTestCase[T any] struct {
	name      string
	src       any
	wantValue T
	wantErr   bool
}

func scanTests[T comparable, F func(src any, setValue func(T) error) error](
	t *testing.T,
	name string,
	tests []scanTestCase[T],
	scanFn F,
) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var got T

			err := scanFn(tt.src, func(v T) error {
				got = v

				return nil
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("%s() error = %v, wantErr %v", name, err, tt.wantErr)

				return
			}

			if got != tt.wantValue {
				t.Errorf("%s() got = %v, want %v", name, got, tt.wantValue)
			}
		})
	}
}

func scanInvalidType(
	t *testing.T,
	typeName string,
	invalidSrc any,
	scanFn func(func() error) error,
) {
	err := scanFn(func() error {
		t.Error("callback should not be called for invalid type")

		return nil
	})
	if err == nil {
		t.Errorf("expected error scanning %T into %s", invalidSrc, typeName)
	}
}

func TestNullableValue(t *testing.T) {
	t.Parallel()

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
		testNullableString(t, tt.name, tt.input, tt.wantNil, tt.wantVal, NullableValue)
	}
}

func TestInt64Value(t *testing.T) {
	t.Parallel()

	got, err := Int64Value(12345)
	if err != nil {
		t.Errorf("Int64Value() error = %v", err)
	}

	if got != int64(12345) {
		t.Errorf("Int64Value() got = %v, want 12345", got)
	}
}

func TestZeroAsNullValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   int64
		wantNil bool
	}{
		{"zero", 0, true},
		{"non-zero", 42, false},
	}

	for _, testCase := range tests {
		testNullable(t, testCase.name, testCase.input, testCase.wantNil, ZeroAsNullValue)
	}
}

func testNullable[T comparable](
	t *testing.T,
	name string,
	input T,
	wantNil bool,
	fn func(T) (any, error),
) {
	t.Run(name, func(t *testing.T) {
		t.Parallel()

		got, err := fn(input)
		if err != nil {
			t.Errorf("error = %v", err)

			return
		}

		if wantNil && got != nil {
			t.Errorf("got nil = false, want nil = true")
		}

		if !wantNil && got != any(input) {
			t.Errorf("got = %v, want %v", got, input)
		}
	})
}

func testNullableString(
	t *testing.T,
	name, input string,
	wantNil bool,
	wantVal string,
	fn func(string) (any, error),
) {
	t.Run(name, func(t *testing.T) {
		t.Parallel()

		got, err := fn(input)
		if err != nil {
			t.Errorf("error = %v", err)

			return
		}

		if wantNil && got != nil {
			t.Errorf("got nil = false, want nil = true")
		}

		if !wantNil && got != wantVal {
			t.Errorf("got = %v, want %v", got, wantVal)
		}
	})
}
