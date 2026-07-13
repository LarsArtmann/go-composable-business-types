package enums_test

import (
	"database/sql/driver"
	"slices"
	"testing"

	"github.com/larsartmann/go-composable-business-types/enums"
	"github.com/larsartmann/go-composable-business-types/testutil"
)

// enumValueCase represents a test case for enum Value() method.
type enumValueCase[T any] struct {
	value    T
	expected string
}

// valuer is an interface for types that implement the driver.Valuer interface.
type valuer interface {
	Value() (driver.Value, error)
}

// testEnumValue runs table-driven Value() tests for enum types.
func testEnumValue[T valuer](t *testing.T, tests []enumValueCase[T]) {
	t.Helper()

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			t.Parallel()

			v, err := tt.value.Value()
			if err != nil {
				t.Errorf("unexpected error: %v", err)

				return
			}

			if v != tt.expected {
				t.Errorf("expected %s, got %v", tt.expected, v)
			}
		})
	}
}

// testEnumValueOne tests a single enum Value() case.
func testEnumValueOne[T valuer](t *testing.T, value T, expected string) {
	t.Helper()
	testEnumValue(t, []enumValueCase[T]{{value, expected}})
}

// enumStringCase represents a test case for enum String() method.
type enumStringCase[T interface{ String() string }] struct {
	value    T
	expected string
}

// testEnumString runs table-driven String() tests for enum types.
func testEnumString[T interface{ String() string }](t *testing.T, tests []enumStringCase[T]) {
	t.Helper()

	cases := make([]testutil.StringCase[T], 0, len(tests))
	for _, tt := range tests {
		cases = append(cases, testutil.StringCase[T]{Value: tt.value, Expected: tt.expected})
	}

	testutil.RunStringTests(t, "", cases)
}

// enumParseCase represents a test case for enum Parse() function.
type enumParseCase[T any] struct {
	input   string
	want    T
	wantErr bool
}

// testParse runs table-driven Parse() tests for enum types.
func testParse[T comparable](
	t *testing.T,
	parse func(string) (T, error),
	tests []enumParseCase[T],
) {
	t.Helper()

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			got, err := parse(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error")
				}

				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

// testEnumIsValid tests that all provided valid values return true for IsValid()
// and that the provided invalid value returns false.
func testEnumIsValid[T interface{ IsValid() bool }](t *testing.T, validValues []T, invalidValue T) {
	t.Helper()

	for _, v := range validValues {
		if !v.IsValid() {
			t.Errorf("expected %v to be valid", v)
		}
	}

	if invalidValue.IsValid() {
		t.Error("invalid value should not be valid")
	}
}

func TestEnumValues(t *testing.T) {
	t.Parallel()
	// Test ActorKindValues
	akValues := enums.ActorKindValues()
	if len(akValues) != 4 {
		t.Errorf("expected 4 ActorKind values, got %d", len(akValues))
	}

	// Test PriorityValues
	pValues := enums.PriorityValues()
	if len(pValues) != 4 {
		t.Errorf("expected 4 Priority values, got %d", len(pValues))
	}

	// Test StatusValues
	sValues := enums.StatusValues()
	if len(sValues) != 5 {
		t.Errorf("expected 5 Status values, got %d", len(sValues))
	}

	// Test TriggerValues
	trigValues := enums.TriggerValues()
	if len(trigValues) != 7 {
		t.Errorf("expected 7 Trigger values, got %d", len(trigValues))
	}

	// Test CauseKindValues
	ckValues := enums.CauseKindValues()
	if len(ckValues) != 3 {
		t.Errorf("expected 3 CauseKind values, got %d", len(ckValues))
	}
}

func TestEnumNames(t *testing.T) {
	t.Parallel()
	// Test ActorKindNames
	akNames := enums.ActorKindNames()
	if len(akNames) != 4 {
		t.Errorf("expected 4 ActorKind names, got %d", len(akNames))
	}

	if !slices.Contains(akNames, "User") {
		t.Error("ActorKindNames should include 'User'")
	}

	// Test PriorityNames
	pNames := enums.PriorityNames()
	if len(pNames) != 4 {
		t.Errorf("expected 4 Priority names, got %d", len(pNames))
	}

	// Test StatusNames
	sNames := enums.StatusNames()
	if len(sNames) != 5 {
		t.Errorf("expected 5 Status names, got %d", len(sNames))
	}

	// Test TriggerNames
	trigNames := enums.TriggerNames()
	if len(trigNames) != 7 {
		t.Errorf("expected 7 Trigger names, got %d", len(trigNames))
	}

	// Test CauseKindNames
	ckNames := enums.CauseKindNames()
	if len(ckNames) != 3 {
		t.Errorf("expected 3 CauseKind names, got %d", len(ckNames))
	}

	if !slices.Contains(ckNames, "Direct") {
		t.Error("CauseKindNames should include 'Direct'")
	}
}

// enumUnmarshalTextErrorCase represents a test case for enum UnmarshalText error.
type enumUnmarshalTextErrorCase[T any] struct {
	name string
}

// testUnmarshalTextError runs table-driven UnmarshalText error tests for enum types.
func testUnmarshalTextError[T any](t *testing.T, tests []enumUnmarshalTextErrorCase[T]) {
	t.Helper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var v T

			err := any(&v).(interface{ UnmarshalText([]byte) error }).UnmarshalText(
				[]byte("Invalid"),
			)
			if err == nil {
				t.Error("expected error for invalid value")
			}
		})
	}
}

// testUnmarshalTextErrorsAll tests UnmarshalText error cases for multiple enum types in one function.
func testUnmarshalTextErrorsAll[T1, T2, T3, T4, T5 any](
	t *testing.T,
	tests1 []enumUnmarshalTextErrorCase[T1],
	tests2 []enumUnmarshalTextErrorCase[T2],
	tests3 []enumUnmarshalTextErrorCase[T3],
	tests4 []enumUnmarshalTextErrorCase[T4],
	tests5 []enumUnmarshalTextErrorCase[T5],
) {
	t.Parallel()
	testUnmarshalTextError(t, tests1)
	testUnmarshalTextError(t, tests2)
	testUnmarshalTextError(t, tests3)
	testUnmarshalTextError(t, tests4)
	testUnmarshalTextError(t, tests5)
}

func TestCauseKind(t *testing.T) {
	t.Parallel()
	testEnumString(t, []enumStringCase[enums.CauseKind]{
		{enums.CauseKindDirect, "Direct"},
		{enums.CauseKindCommand, "Command"},
		{enums.CauseKindEvent, "Event"},
	})
}

func TestParseCauseKind(t *testing.T) {
	t.Parallel()

	testParse(t, enums.ParseCauseKind, []enumParseCase[enums.CauseKind]{
		{"Direct", enums.CauseKindDirect, false},
		{"Command", enums.CauseKindCommand, false},
		{"Event", enums.CauseKindEvent, false},
		{"Invalid", 0, true},
		{"", 0, true},
	})
}

func TestCauseKindIsValid(t *testing.T) {
	t.Parallel()

	testEnumIsValid(
		t,
		[]enums.CauseKind{enums.CauseKindDirect, enums.CauseKindEvent},
		enums.CauseKind(99),
	)
}
