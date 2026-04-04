package bounded_test

import (
	"encoding/json"
	"testing"
)

func TestNewBoundedString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		minLen  uint
		maxLen  uint
		value   string
		wantErr bool
	}{
		{"valid", 1, 100, "hello", false},
		{"at min", 5, 10, "hello", false},
		{"at max", 1, 5, "hello", false},
		{"empty when min 0", 0, 10, "", false},
		{"too short", 10, 20, "hi", true},
		{"too long", 1, 3, "hello", true},
		{"negative min", 4294967295, 10, "test", true},
		{"max less than min", 10, 5, "test", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			bs, err := NewBoundedString(tt.minLen, tt.maxLen, tt.value)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error")
				}

				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if bs.String() != tt.value {
				t.Errorf("expected %s, got %s", tt.value, bs.String())
			}
		})
	}
}

func TestBoundedStringLen(t *testing.T) {
	t.Parallel()

	bs, _ := NewBoundedString(1, 100, "hello")
	if bs.Len() != 5 {
		t.Errorf("expected length 5, got %d", bs.Len())
	}

	// Unicode test
	bs2, _ := NewBoundedString(1, 100, "héllo")
	if bs2.Len() != 5 {
		t.Errorf("expected length 5 for unicode string, got %d", bs2.Len())
	}
}

func TestBoundedStringOf(t *testing.T) {
	t.Parallel()

	NewName := BoundedStringOf(1, 100)

	name, err := NewName("John Doe")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if name.String() != "John Doe" {
		t.Errorf("expected 'John Doe', got %s", name.String())
	}

	// Test with invalid name
	_, err = NewName("")
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestNonEmptyString(t *testing.T) {
	t.Parallel()

	s, err := NonEmptyString(100, "hello")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if s.String() != "hello" {
		t.Errorf("expected 'hello', got %s", s.String())
	}

	// Empty string should fail
	_, err = NonEmptyString(100, "")
	if err == nil {
		t.Error("expected error for empty string")
	}
}

func TestTrimmedBoundedString(t *testing.T) {
	t.Parallel()

	s, err := TrimmedBoundedString(1, 100, "  hello  ")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if s.String() != "hello" {
		t.Errorf("expected 'hello' (trimmed), got %s", s.String())
	}
}

func TestBoundedStringIsZero(t *testing.T) {
	t.Parallel()

	var zero BoundedString
	if !zero.IsZero() {
		t.Error("expected zero BoundedString to be zero")
	}

	nonZero, _ := NewBoundedString(1, 100, "test")
	if nonZero.IsZero() {
		t.Error("expected non-zero BoundedString to not be zero")
	}
}

func TestBoundedStringBounds(t *testing.T) {
	t.Parallel()

	bs, _ := NewBoundedString(5, 10, "hello")
	if bs.MinLen() != 5 {
		t.Errorf("expected MinLen 5, got %d", bs.MinLen())
	}

	if bs.MaxLen() != 10 {
		t.Errorf("expected MaxLen 10, got %d", bs.MaxLen())
	}
	// "hello" is length 5, so it's at min but not at max
	if !bs.IsMinLength() {
		t.Error("expected to be at min")
	}

	if bs.IsMaxLength() {
		t.Error("should not be at max")
	}
}

func TestBoundedStringMarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		bs      BoundedString
		want    string
		wantErr bool
	}{
		{"simple", mustBoundedString(1, 100, "hello"), `"hello"`, false},
		{"empty", mustBoundedString(0, 100, ""), `""`, false},
		{"unicode", mustBoundedString(1, 100, "héllo"), `"héllo"`, false},
		{"with quotes", mustBoundedString(1, 100, `say "hi"`), `"say \"hi\""`, false},
		{"special chars", mustBoundedString(1, 100, "line1\nline2"), `"line1\nline2"`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := json.Marshal(tt.bs)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error")
				}

				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)

				return
			}

			if string(got) != tt.want {
				t.Errorf("expected %s, got %s", tt.want, string(got))
			}
		})
	}
}

func TestBoundedStringUnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"simple", `"hello"`, "hello", false},
		{"empty", `""`, "", false},
		{"unicode", `"héllo"`, "héllo", false},
		{"with space", `"hello world"`, "hello world", false},
		{"invalid json", `not json`, "", true},
		{"number", `123`, "", true},
		{"null", `null`, "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var bs BoundedString

			err := json.Unmarshal([]byte(tt.input), &bs)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error")
				}

				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)

				return
			}

			if bs.String() != tt.want {
				t.Errorf("expected %q, got %q", tt.want, bs.String())
			}
		})
	}
}

func TestBoundedStringJSONRoundTrip(t *testing.T) {
	t.Parallel()

	original := mustBoundedString(1, 100, "test value")

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var restored BoundedString
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if original.String() != restored.String() {
		t.Errorf("expected %q, got %q", original.String(), restored.String())
	}
}

func TestBoundedStringScan(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		src     any
		want    string
		wantErr bool
	}{
		{"string", "hello", "hello", false},
		{"empty string", "", "", false},
		{"bytes", []byte("world"), "world", false},
		{"empty bytes", []byte(""), "", false},
		{"nil", nil, "", false},
		{"int", 123, "", true},
		{"float", 1.5, "", true},
		{"bool", true, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var bs BoundedString

			err := bs.Scan(tt.src)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error")
				}

				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)

				return
			}

			if bs.String() != tt.want {
				t.Errorf("expected %q, got %q", tt.want, bs.String())
			}
		})
	}
}

func TestBoundedStringScanNilReceiver(t *testing.T) {
	t.Parallel()

	var bs *BoundedString

	err := bs.Scan("test")
	if err == nil {
		t.Error("expected error for nil receiver")
	}
}

func TestBoundedStringValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		bs      BoundedString
		want    any
		wantErr bool
	}{
		{"non-empty", mustBoundedString(1, 100, "hello"), "hello", false},
		{"empty", mustBoundedString(0, 100, ""), nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.bs.Value()
			if tt.wantErr {
				if err == nil {
					t.Error("expected error")
				}

				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)

				return
			}

			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestBoundedStringUnicodeLength(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		want  uint
	}{
		{"ascii", "hello", 5},
		{"latin with accent", "héllo", 5},
		{"emoji", "👍👍👍", 3},
		{"mixed", "a👍b", 3},
		{"chinese", "你好世界", 4},
		{"japanese", "こんにちは", 5},
		{"arabic", "مرحبا", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			bs, err := NewBoundedString(0, 100, tt.value)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if bs.Len() != tt.want {
				t.Errorf("expected length %d, got %d", tt.want, bs.Len())
			}
		})
	}
}

func TestBoundedStringEdgeCases(t *testing.T) {
	t.Parallel()
	t.Run("min equals max", func(t *testing.T) {
		t.Parallel()

		bs, err := NewBoundedString(5, 5, "hello")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !bs.IsMinLength() || !bs.IsMaxLength() {
			t.Error("expected to be at both min and max")
		}
	})

	t.Run("zero bounds", func(t *testing.T) {
		t.Parallel()

		bs, err := NewBoundedString(0, 0, "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !bs.IsZero() {
			t.Error("expected zero value")
		}
	})

	t.Run("large max", func(t *testing.T) {
		t.Parallel()

		large := make([]byte, 10000)
		for i := range large {
			large[i] = 'a'
		}

		bs, err := NewBoundedString(0, 20000, string(large))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if bs.Len() != 10000 {
			t.Errorf("expected length 10000, got %d", bs.Len())
		}
	})
}

func TestBoundedStringIsMinMaxLength(t *testing.T) {
	t.Parallel()
	t.Run("at min", func(t *testing.T) {
		t.Parallel()

		bs, _ := NewBoundedString(3, 10, "abc")
		if !bs.IsMinLength() {
			t.Error("expected to be at min length")
		}

		if bs.IsMaxLength() {
			t.Error("should not be at max length")
		}
	})

	t.Run("at max", func(t *testing.T) {
		t.Parallel()

		bs, _ := NewBoundedString(1, 5, "abcde")
		if bs.IsMinLength() {
			t.Error("should not be at min length")
		}

		if !bs.IsMaxLength() {
			t.Error("expected to be at max length")
		}
	})

	t.Run("in middle", func(t *testing.T) {
		t.Parallel()

		bs, _ := NewBoundedString(1, 10, "abc")
		if bs.IsMinLength() {
			t.Error("should not be at min length")
		}

		if bs.IsMaxLength() {
			t.Error("should not be at max length")
		}
	})
}

func mustBoundedString(minLen, _ uint, value string) BoundedString {
	bs, err := NewBoundedString(minLen, 100, value)
	if err != nil {
		panic(err)
	}

	return bs
}
