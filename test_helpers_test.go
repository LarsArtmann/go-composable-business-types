package cbt

import (
	"testing"
)

// assertError fails if err is nil.
func assertError(t *testing.T, err error, msg string) {
	t.Helper()
	if err == nil {
		t.Errorf("expected error: %s", msg)
	}
}

// assertNoError fails if err is not nil.
func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// assertTag fails if tags[key] != expected.
func assertTag(t *testing.T, tags map[string]string, key, expected string) {
	t.Helper()
	if tags[key] != expected {
		t.Errorf("tag %s mismatch: expected %q, got %q", key, expected, tags[key])
	}
}

// assertJSONEquals fails if string(data) != expected.
func assertJSONEquals(t *testing.T, data []byte, expected string) {
	t.Helper()
	if string(data) != expected {
		t.Errorf("JSON mismatch: expected %s, got %s", expected, string(data))
	}
}

// assertStringEquals fails if got != expected.
func assertStringEquals[T ~string](t *testing.T, got, expected T) {
	t.Helper()
	if got != expected {
		t.Errorf("string mismatch: expected %q, got %q", expected, got)
	}
}

// assertNil fails if val is not nil.
func assertNil(t *testing.T, val any) {
	t.Helper()
	if val != nil {
		t.Errorf("expected nil, got %v", val)
	}
}

// assertNotNil fails if val is nil.
func assertNotNil(t *testing.T, val any) {
	t.Helper()
	if val == nil {
		t.Error("expected non-nil value")
	}
}

// assertEqual fails if got != expected.
func assertEqual[T comparable](t *testing.T, got, expected T) {
	t.Helper()
	if got != expected {
		t.Errorf("mismatch: expected %v, got %v", expected, got)
	}
}

// assertZero fails if v.IsZero() returns false.
func assertZero[T interface{ IsZero() bool }](t *testing.T, v T) {
	t.Helper()
	if !v.IsZero() {
		t.Errorf("expected zero value, got %v", v)
	}
}

// assertNotZero fails if v.IsZero() returns true.
func assertNotZero[T interface{ IsZero() bool }](t *testing.T, v T) {
	t.Helper()
	if v.IsZero() {
		t.Error("expected non-zero value")
	}
}

// assertTrue fails if cond is false.
func assertTrue(t *testing.T, cond bool, msg string) {
	t.Helper()
	if !cond {
		t.Error(msg)
	}
}

// assertFalse fails if cond is true.
func assertFalse(t *testing.T, cond bool, msg string) {
	t.Helper()
	if cond {
		t.Error(msg)
	}
}
