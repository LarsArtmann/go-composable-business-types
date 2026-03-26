package errors

import (
	"errors"
	"testing"
)

// =============================================================================
// Wrap Helper Tests
// =============================================================================

func TestWrapMalformed(t *testing.T) {
	t.Parallel()
	t.Run("wraps error", func(t *testing.T) {
		t.Parallel()
		inner := errors.New("invalid syntax")
		err := WrapMalformed(inner, "JSON", `{bad}`)

		if !errors.Is(err, ErrMalformedInput) {
			t.Error("expected errors.Is(err, ErrMalformedInput) = true")
		}
	})

	t.Run("nil returns nil", func(t *testing.T) {
		t.Parallel()
		err := WrapMalformed(nil, "JSON", `{}`)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})
}

func TestWrapInvalid(t *testing.T) {
	t.Parallel()
	t.Run("wraps error", func(t *testing.T) {
		t.Parallel()
		inner := errors.New("invalid value")
		err := WrapInvalid(inner, "email", "test@")

		var target *ValidationError
		if !errors.As(err, &target) {
			t.Error("expected errors.As() to succeed")
		}
		if target.Field != "email" {
			t.Errorf("Field = %q, want email", target.Field)
		}
	})

	t.Run("nil returns nil", func(t *testing.T) {
		t.Parallel()
		err := WrapInvalid(nil, "field", "value")
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})
}

func TestWrapRange(t *testing.T) {
	t.Parallel()
	err := WrapRange(150, 10, 100, true)

	var target *RangeError
	if !errors.As(err, &target) {
		t.Error("expected errors.As() to succeed")
	}
	if target.Value != 150 {
		t.Errorf("Value = %v, want 150", target.Value)
	}
	if target.Min != 10 {
		t.Errorf("Min = %v, want 10", target.Min)
	}
	if target.Max != 100 {
		t.Errorf("Max = %v, want 100", target.Max)
	}
}

func TestWrapScan(t *testing.T) {
	t.Parallel()
	t.Run("wraps error", func(t *testing.T) {
		t.Parallel()
		inner := errors.New("conversion failed")
		err := WrapScan(inner, "[]byte", "string")

		var target *ScanError
		if !errors.As(err, &target) {
			t.Error("expected errors.As() to succeed")
		}
		if target.SourceType != "[]byte" {
			t.Errorf("SourceType = %q, want []byte", target.SourceType)
		}
		if target.TargetType != "string" {
			t.Errorf("TargetType = %q, want string", target.TargetType)
		}
	})

	t.Run("nil returns nil", func(t *testing.T) {
		t.Parallel()
		err := WrapScan(nil, "int", "string")
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})
}

func TestWrapUnmarshal(t *testing.T) {
	t.Parallel()
	t.Run("wraps error", func(t *testing.T) {
		t.Parallel()
		inner := errors.New("unexpected EOF")
		err := WrapUnmarshal(inner, "JSON", `{incomplete`)

		var target *UnmarshalError
		if !errors.As(err, &target) {
			t.Error("expected errors.As() to succeed")
		}
		if target.Type != "JSON" {
			t.Errorf("Type = %q, want JSON", target.Type)
		}
		if target.Input != `{incomplete` {
			t.Errorf("Input = %q, want {incomplete", target.Input)
		}
	})

	t.Run("nil returns nil", func(t *testing.T) {
		t.Parallel()
		err := WrapUnmarshal(nil, "XML", `<tag>`)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})
}
