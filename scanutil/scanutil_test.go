package scanutil

import (
	"fmt"
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

func TestNullableValueWithError(t *testing.T) {
	t.Parallel()

	t.Run("non-empty string", func(t *testing.T) {
		t.Parallel()

		got, err := NullableValueWithError("hello", "test")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if got != "hello" {
			t.Errorf("got = %v, want hello", got)
		}
	})

	t.Run("empty string returns nil", func(t *testing.T) {
		t.Parallel()

		got, err := NullableValueWithError("", "test")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if got != nil {
			t.Errorf("got = %v, want nil", got)
		}
	})
}

func TestNonNullableValue(t *testing.T) {
	t.Parallel()

	t.Run("non-empty string", func(t *testing.T) {
		t.Parallel()

		got, err := NonNullableValue("hello")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if got != "hello" {
			t.Errorf("got = %v, want hello", got)
		}
	})

	t.Run("empty string still returns value", func(t *testing.T) {
		t.Parallel()

		got, err := NonNullableValue("")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if got != "" {
			t.Errorf("got = %v, want empty string", got)
		}
	})
}

func TestScanInt64InvalidByteSlice(t *testing.T) {
	t.Parallel()

	err := ScanInt64([]byte("not-a-number"), func(_ int64) error { return nil })
	if err == nil {
		t.Error("expected error for invalid byte slice")
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

func TestScanEnum(t *testing.T) {
	t.Parallel()

	type testEnum uint8

	const (
		testEnumA testEnum = iota
		testEnumB
		testEnumC
	)

	parseTestEnum := func(s string) (testEnum, error) {
		switch s {
		case "A":
			return testEnumA, nil
		case "B":
			return testEnumB, nil
		case "C":
			return testEnumC, nil
		default:
			return 0, fmt.Errorf("unknown testEnum: %s", s)
		}
	}

	tests := []struct {
		name    string
		src     any
		want    testEnum
		wantErr bool
	}{
		{name: "nil", src: nil, want: testEnumA},
		{name: "int64", src: int64(1), want: testEnumB},
		{name: "string", src: "B", want: testEnumB},
		{name: "[]byte", src: []byte("C"), want: testEnumC},
		{name: "enum value", src: testEnumC, want: testEnumC},
		{name: "int", src: int(2), want: testEnumC},
		{name: "uint", src: uint(1), want: testEnumB},
		{name: "uint64", src: uint64(0), want: testEnumA},
		{name: "float64", src: float64(1), want: testEnumB},
		{name: "invalid string", src: "Z", wantErr: true},
		{name: "invalid type", src: struct{}{}, wantErr: true},
		{name: "nil *int", src: (*int)(nil), wantErr: true},
		{name: "nil *string", src: (*string)(nil), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var got testEnum

			err := ScanEnum(&got, tt.src, parseTestEnum)
			if (err != nil) != tt.wantErr {
				t.Errorf("ScanEnum() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr && got != tt.want {
				t.Errorf("ScanEnum() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScanEnumPointer(t *testing.T) {
	t.Parallel()

	type testEnum uint8

	const testEnumB testEnum = 1

	parseTestEnum := func(s string) (testEnum, error) {
		return testEnumB, nil
	}

	t.Run("*enum value", func(t *testing.T) {
		t.Parallel()

		val := testEnumB

		var got testEnum

		err := ScanEnum(&got, &val, parseTestEnum)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if got != testEnumB {
			t.Errorf("got = %v, want %v", got, testEnumB)
		}
	})

	t.Run("nil *enum", func(t *testing.T) {
		t.Parallel()

		var got testEnum

		err := ScanEnum(&got, (*testEnum)(nil), parseTestEnum)
		if err == nil {
			t.Error("expected error for nil pointer")
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
