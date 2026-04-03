package enums

import (
	"database/sql/driver"
	"slices"
	"testing"
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

// enumStringCase represents a test case for enum String() method.
type enumStringCase[T interface{ String() string }] struct {
	value    T
	expected string
}

// testEnumString runs table-driven String() tests for enum types.
func testEnumString[T interface{ String() string }](t *testing.T, tests []enumStringCase[T]) {
	t.Helper()

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			t.Parallel()

			if tt.value.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, tt.value.String())
			}
		})
	}
}

// enumParseCase represents a test case for enum Parse() function.
type enumParseCase[T any] struct {
	input   string
	want    T
	wantErr bool
}

// testParse runs table-driven Parse() tests for enum types.
func testParse[T comparable](t *testing.T, parse func(string) (T, error), tests []enumParseCase[T]) {
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

// enumIsValidCase represents a test case for enum IsValid() method.
type enumIsValidCase[T any] struct {
	value    T
	expected bool
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
	akValues := ActorKindValues()
	if len(akValues) != 4 {
		t.Errorf("expected 4 ActorKind values, got %d", len(akValues))
	}

	// Test PriorityValues
	pValues := PriorityValues()
	if len(pValues) != 4 {
		t.Errorf("expected 4 Priority values, got %d", len(pValues))
	}

	// Test StatusValues
	sValues := StatusValues()
	if len(sValues) != 5 {
		t.Errorf("expected 5 Status values, got %d", len(sValues))
	}

	// Test TriggerValues
	tValues := TriggerValues()
	if len(tValues) != 7 {
		t.Errorf("expected 7 Trigger values, got %d", len(tValues))
	}

	// Test CauseKindValues
	ckValues := CauseKindValues()
	if len(ckValues) != 3 {
		t.Errorf("expected 3 CauseKind values, got %d", len(ckValues))
	}
}

func TestEnumNames(t *testing.T) {
	t.Parallel()
	// Test ActorKindNames
	akNames := ActorKindNames()
	if len(akNames) != 4 {
		t.Errorf("expected 4 ActorKind names, got %d", len(akNames))
	}

	if !slices.Contains(akNames, "User") {
		t.Error("ActorKindNames should include 'User'")
	}

	// Test PriorityNames
	pNames := PriorityNames()
	if len(pNames) != 4 {
		t.Errorf("expected 4 Priority names, got %d", len(pNames))
	}

	// Test StatusNames
	sNames := StatusNames()
	if len(sNames) != 5 {
		t.Errorf("expected 5 Status names, got %d", len(sNames))
	}

	// Test TriggerNames
	tNames := TriggerNames()
	if len(tNames) != 7 {
		t.Errorf("expected 7 Trigger names, got %d", len(tNames))
	}

	// Test CauseKindNames
	ckNames := CauseKindNames()
	if len(ckNames) != 3 {
		t.Errorf("expected 3 CauseKind names, got %d", len(ckNames))
	}

	if !slices.Contains(ckNames, "Direct") {
		t.Error("CauseKindNames should include 'Direct'")
	}
}

func TestCauseKind(t *testing.T) {
	t.Parallel()
	testEnumString(t, []enumStringCase[CauseKind]{
		{CauseKindDirect, "Direct"},
		{CauseKindCommand, "Command"},
		{CauseKindEvent, "Event"},
	})
}

func TestParseCauseKind(t *testing.T) {
	t.Parallel()

	testParse(t, ParseCauseKind, []enumParseCase[CauseKind]{
		{"Direct", CauseKindDirect, false},
		{"Command", CauseKindCommand, false},
		{"Event", CauseKindEvent, false},
		{"Invalid", 0, true},
		{"", 0, true},
	})
}

func TestCauseKindIsValid(t *testing.T) {
	t.Parallel()

	testEnumIsValid(t, []CauseKind{CauseKindDirect, CauseKindEvent}, CauseKind(99))
}
