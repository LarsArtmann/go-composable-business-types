// Package testutil provides generic test helpers for parse/string/zero-value testing.
//
// It offers table-driven test runners for types implementing String(), IsZero(),
// and parse constructors — reducing boilerplate in type validation tests.
package testutil

import "testing"

// ParseTester is a constraint for types with a String method.
type ParseTester interface {
	String() string
}

// RunParseTest runs a single parse test.
func RunParseTest[T ParseTester](
	t *testing.T,
	typeName, input string,
	parse func(string) (T, error),
	wantErr bool,
) {
	val, err := parse(input)
	if wantErr {
		if err == nil {
			t.Errorf("expected error for %s", typeName)
		}

		return
	}

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if val.String() != input {
		t.Errorf("expected %s, got %s", input, val.String())
	}
}

// ParseTestCase is a test case for table-driven parse tests.
type ParseTestCase[T any] struct {
	Name    string
	Input   string
	WantErr bool
}

// RunParseTests runs a suite of parse test cases as subtests.
func RunParseTests[T ParseTester](
	t *testing.T,
	typeName string,
	tests []ParseTestCase[T],
	parse func(string) (T, error),
) {
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			RunParseTest(t, typeName, tc.Input, parse, tc.WantErr)
		})
	}
}

// RunParseErrorTest asserts that parsing an empty string returns an error.
func RunParseErrorTest[T ParseTester](
	t *testing.T,
	typeName string,
	parse func(string) (T, error),
) {
	t.Parallel()
	RunParseTest(t, typeName, "", parse, true)
}

// RunAppendTextTest tests appendText output.
func RunAppendTextTest[T any](
	t *testing.T,
	expected string,
	appendText func(T) ([]byte, error),
	val T,
) {
	n, err := appendText(val)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if string(n) != expected {
		t.Errorf("expected %s, got %s", expected, string(n))
	}
}

// AppendTexter is a constraint for types implementing AppendText.
type AppendTexter[T any] interface {
	AppendText(b []byte) ([]byte, error)
}

// RunAppendTextTestSimple tests AppendText on a value directly.
func RunAppendTextTestSimple[T AppendTexter[T]](t *testing.T, expected string, val T) {
	n, err := val.AppendText(nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if string(n) != expected {
		t.Errorf("expected %s, got %s", expected, string(n))
	}
}

// ZeroChecker is a constraint for types reporting zero state.
type ZeroChecker interface {
	IsZero() bool
}

// RunIsZeroTest verifies zero and non-zero value behavior.
func RunIsZeroTest[T ZeroChecker](t *testing.T, makeNonZero func() (T, error)) {
	t.Parallel()

	var zero T
	if !zero.IsZero() {
		t.Error("zero value should be zero")
	}

	nonZero, err := makeNonZero()
	if err != nil {
		t.Fatalf("failed to create non-zero value: %v", err)
	}

	if nonZero.IsZero() {
		t.Error("non-zero value should not be zero")
	}
}

// PartAccessor defines a test case for accessing a sub-part.
type PartAccessor[T any] struct {
	Name     string
	Get      func(T) string
	Expected string
}

// RunPartsTest checks accessor output for a value.
func RunPartsTest[T any](t *testing.T, val T, accessors []PartAccessor[T]) {
	t.Helper()

	for _, accessor := range accessors {
		got := accessor.Get(val)
		if got != accessor.Expected {
			t.Errorf("expected %s '%s', got '%s'", accessor.Name, accessor.Expected, got)
		}
	}
}

// StringCase is a test case for String() output.
type StringCase[T any] struct {
	Value    T
	Expected string
}

// Stringer is a constraint for types implementing String().
type Stringer interface {
	String() string
}

// RunStringTests runs a suite of String() test cases.
func RunStringTests[T Stringer](t *testing.T, _ string, tests []StringCase[T]) {
	t.Helper()

	for _, tc := range tests {
		t.Run(tc.Expected, func(t *testing.T) {
			t.Parallel()

			if tc.Value.String() != tc.Expected {
				t.Errorf("expected %s, got %s", tc.Expected, tc.Value.String())
			}
		})
	}
}
