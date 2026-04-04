package testutil

import "testing"

type ParseTester interface {
	String() string
}

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

type ParseTestCase[T any] struct {
	Name    string
	Input   string
	WantErr bool
}

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

func RunParseErrorTest[T ParseTester](
	t *testing.T,
	typeName string,
	parse func(string) (T, error),
) {
	t.Parallel()
	RunParseTest(t, typeName, "", parse, true)
}

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

type ZeroChecker interface {
	IsZero() bool
}

func RunIsZeroTest[T ZeroChecker](t *testing.T, makeNonZero func() (T, error)) {
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

type PartAccessor[T any] struct {
	Name     string
	Get      func(T) string
	Expected string
}

func RunPartsTest[T any](t *testing.T, val T, accessors []PartAccessor[T]) {
	t.Helper()
	for _, accessor := range accessors {
		got := accessor.Get(val)
		if got != accessor.Expected {
			t.Errorf("expected %s '%s', got '%s'", accessor.Name, accessor.Expected, got)
		}
	}
}
