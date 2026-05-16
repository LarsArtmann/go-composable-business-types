package enums_test

import (
	"fmt"
	"testing"

	"github.com/larsartmann/go-composable-business-types/enums"
	"github.com/larsartmann/go-composable-business-types/testutil"
)

// Test SQL Scanner/Valuer interfaces.
func TestActorKindSQL(t *testing.T) {
	t.Parallel()
	// Test Value
	val, err := enums.ActorKindUser.Value()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if val != "User" {
		t.Errorf("expected User, got %v", val)
	}

	// Test Scan with string
	var ak enums.ActorKind
	if err := ak.Scan("Bot"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if ak != enums.ActorKindBot {
		t.Errorf("expected Bot, got %v", ak)
	}

	// Test Scan with []byte
	var ak2 enums.ActorKind
	if err := ak2.Scan([]byte("System")); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if ak2 != enums.ActorKindSystem {
		t.Errorf("expected System, got %v", ak2)
	}

	// Test Scan with nil
	ak3 := enums.ActorKindUser
	if err := ak3.Scan(nil); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if ak3 != 0 {
		t.Errorf("expected zero value, got %v", ak3)
	}

	// Test Scan with int (supported type - sets raw value)
	var ak4 enums.ActorKind
	if err := ak4.Scan(2); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if ak4 != enums.ActorKindSystem {
		t.Errorf("expected enums.ActorKindSystem(2), got %v", ak4)
	}

	// Test Scan with int64 (supported type)
	var ak5 enums.ActorKind
	if err := ak5.Scan(int64(3)); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if ak5 != enums.ActorKindService {
		t.Errorf("expected enums.ActorKindService(3), got %v", ak5)
	}
}

// enumSQLCase holds test data for SQL Value/Scan tests.
type enumSQLCase[T any] struct {
	value    T
	valueStr string
	scanStr  string
	scanWant T
}

// testEnumSQL tests Value() and Scan() methods for an enum type.
func testEnumSQL[T comparable](
	t *testing.T,
	cases []enumSQLCase[T],
	valueFunc func(T) (string, error),
	scanFunc func(*T, any) error,
) {
	t.Helper()

	for _, tc := range cases {
		t.Run(tc.valueStr, func(t *testing.T) {
			t.Parallel()

			// Test Value
			val, err := valueFunc(tc.value)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if val != tc.valueStr {
				t.Errorf("expected %s, got %v", tc.valueStr, val)
			}

			// Test Scan with string
			var e T
			if err := scanFunc(&e, tc.scanStr); err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if e != tc.scanWant {
				t.Errorf("expected %v, got %v", tc.scanWant, e)
			}
		})
	}
}

// enumValueFunc returns the string representation of a value's MarshalText.
func enumValueFunc[T interface{ MarshalText() ([]byte, error) }](v T) (string, error) {
	b, err := v.MarshalText()

	return string(b), err
}

func TestPrioritySQL(t *testing.T) {
	t.Parallel()
	testEnumSQL(t, []enumSQLCase[enums.Priority]{
		{
			value:    enums.PriorityHigh,
			valueStr: "High",
			scanStr:  "Critical",
			scanWant: enums.PriorityCritical,
		},
	}, enumValueFunc, (*enums.Priority).Scan)
}

func TestStatusSQL(t *testing.T) {
	t.Parallel()
	testEnumSQL(t, []enumSQLCase[enums.Status]{
		{
			value:    enums.StatusActive,
			valueStr: "Active",
			scanStr:  "Archived",
			scanWant: enums.StatusArchived,
		},
	}, enumValueFunc, (*enums.Status).Scan)
}

func TestTriggerSQL(t *testing.T) {
	t.Parallel()
	testEnumSQL(t, []enumSQLCase[enums.Trigger]{
		{
			value:    enums.TriggerWebhook,
			valueStr: "Webhook",
			scanStr:  "Migration",
			scanWant: enums.TriggerMigration,
		},
	}, enumValueFunc, (*enums.Trigger).Scan)
}

// enumMarshalCase holds test data for MarshalText/UnmarshalText tests.
type enumMarshalCase[T any] struct {
	marshalValue  T
	marshalStr    string
	unmarshalStr  string
	unmarshalWant T
}

// testMarshalUnmarshal tests MarshalText() and UnmarshalText() methods for an enum type.
func testMarshalUnmarshal[T any](
	t *testing.T,
	cases []enumMarshalCase[T],
) {
	t.Helper()

	for _, tc := range cases {
		t.Run(tc.marshalStr, func(t *testing.T) {
			t.Parallel()

			data, err := any(&tc.marshalValue).(interface{ MarshalText() ([]byte, error) }).MarshalText()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if string(data) != tc.marshalStr {
				t.Errorf("expected %s, got %s", tc.marshalStr, string(data))
			}

			var e T
			if err := any(&e).(interface{ UnmarshalText([]byte) error }).UnmarshalText(
				[]byte(tc.unmarshalStr),
			); err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if any(e) != any(tc.unmarshalWant) {
				t.Errorf("expected %v, got %v", tc.unmarshalWant, e)
			}
		})
	}
}

// testEnumMarshal is a helper for testing enum MarshalText/UnmarshalText with a single test case.
func testEnumMarshal[T any](
	t *testing.T,
	marshalValue T,
	marshalStr, unmarshalStr string,
	unmarshalWant T,
) {
	t.Parallel()
	testMarshalUnmarshal(t, []enumMarshalCase[T]{
		{
			marshalValue:  marshalValue,
			marshalStr:    marshalStr,
			unmarshalStr:  unmarshalStr,
			unmarshalWant: unmarshalWant,
		},
	})
}

func TestPriorityMarshal(t *testing.T) {
	testEnumMarshal(t, enums.PriorityLow, "Low", "Medium", enums.PriorityMedium)
}

func TestStatusMarshal(t *testing.T) {
	testEnumMarshal(t, enums.StatusDeleted, "Deleted", "Draft", enums.StatusDraft)
}

func TestTriggerMarshal(t *testing.T) {
	testEnumMarshal(t, enums.TriggerCorrection, "Correction", "Scheduled", enums.TriggerScheduled)
}

func TestCauseKindSQL(t *testing.T) {
	t.Parallel()

	val, err := enums.CauseKindDirect.Value()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if val != "Direct" {
		t.Errorf("expected Direct, got %v", val)
	}

	var ck enums.CauseKind
	if err := ck.Scan("Command"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if ck != enums.CauseKindCommand {
		t.Errorf("expected Command, got %v", ck)
	}

	// Test Scan with []byte
	var ck2 enums.CauseKind
	if err := ck2.Scan([]byte("Event")); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if ck2 != enums.CauseKindEvent {
		t.Errorf("expected Event, got %v", ck2)
	}

	// Test Scan with nil
	ck3 := enums.CauseKindDirect
	if err := ck3.Scan(nil); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if ck3 != 0 {
		t.Errorf("expected zero value, got %v", ck3)
	}

	// Test Scan with int
	var ck4 enums.CauseKind
	if err := ck4.Scan(1); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if ck4 != enums.CauseKindCommand {
		t.Errorf("expected enums.CauseKindCommand(1), got %v", ck4)
	}
}

func TestCauseKindMarshal(t *testing.T) {
	testEnumMarshal(t, enums.CauseKindEvent, "Event", "Direct", enums.CauseKindDirect)
}

// TestAppendText tests AppendText for all enum types.
func TestAppendText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value any
	}{
		{"User", enums.ActorKindUser},
		{"High", enums.PriorityHigh},
		{"Active", enums.StatusActive},
		{"Webhook", enums.TriggerWebhook},
		{"Direct", enums.CauseKindDirect},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			switch v := tt.value.(type) {
			case enums.ActorKind:
				testutil.RunAppendTextTestSimple(t, tt.name, &v)
			case enums.Priority:
				testutil.RunAppendTextTestSimple(t, tt.name, &v)
			case enums.Status:
				testutil.RunAppendTextTestSimple(t, tt.name, &v)
			case enums.Trigger:
				testutil.RunAppendTextTestSimple(t, tt.name, &v)
			case enums.CauseKind:
				testutil.RunAppendTextTestSimple(t, tt.name, &v)
			}
		})
	}
}

// testInvalidEnumString verifies that invalid enum values format correctly.
func testInvalidEnumString[T any](t *testing.T, invalidValue T, typeName string, value int) {
	t.Helper()

	expected := fmt.Sprintf("%s(%d)", typeName, value)
	if any(invalidValue).(fmt.Stringer).String() != expected {
		t.Errorf("expected %s, got %s", expected, any(invalidValue).(fmt.Stringer).String())
	}
}

// TestInvalidEnumStrings verifies invalid enum String() formatting.
func TestInvalidEnumStrings(t *testing.T) {
	t.Parallel()

	t.Run("ActorKind", func(t *testing.T) {
		testInvalidEnumString(t, enums.ActorKind(99), "ActorKind", 99)
	})
	t.Run("Priority", func(t *testing.T) {
		testInvalidEnumString(t, enums.Priority(99), "Priority", 99)
	})
	t.Run("Status", func(t *testing.T) {
		testInvalidEnumString(t, enums.Status(99), "Status", 99)
	})
	t.Run("Trigger", func(t *testing.T) {
		testInvalidEnumString(t, enums.Trigger(99), "Trigger", 99)
	})
	t.Run("CauseKind", func(t *testing.T) {
		testInvalidEnumString(t, enums.CauseKind(99), "CauseKind", 99)
	})
}

// Test UnmarshalText error cases.
func TestUnmarshalTextErrors(t *testing.T) {
	testUnmarshalTextErrorsAll(t,
		[]enumUnmarshalTextErrorCase[enums.ActorKind]{{"enums.ActorKind"}},
		[]enumUnmarshalTextErrorCase[enums.Priority]{{"enums.Priority"}},
		[]enumUnmarshalTextErrorCase[enums.Status]{{"enums.Status"}},
		[]enumUnmarshalTextErrorCase[enums.Trigger]{{"enums.Trigger"}},
		[]enumUnmarshalTextErrorCase[enums.CauseKind]{{"enums.CauseKind"}},
	)
}

// Test all Value methods.
func TestValueMethods(t *testing.T) {
	t.Run("ActorKind", func(t *testing.T) {
		t.Parallel()
		testEnumValueOne(t, enums.ActorKindSystem, "System")
	})
	t.Run("Priority", func(t *testing.T) {
		t.Parallel()
		testEnumValueOne(t, enums.PriorityMedium, "Medium")
	})
	t.Run("Status", func(t *testing.T) {
		t.Parallel()
		testEnumValueOne(t, enums.StatusPaused, "Paused")
	})
	t.Run("Trigger", func(t *testing.T) {
		t.Parallel()
		testEnumValueOne(t, enums.TriggerImport, "Import")
	})
	t.Run("CauseKind", func(t *testing.T) {
		t.Parallel()
		testEnumValueOne(t, enums.CauseKindEvent, "Event")
	})
}

// scanTestCase is a generic test case for Scan methods.
type scanTestCase[T any] struct {
	name string
	src  any
	want T
}

// testScanAllTypes is a generic helper for testing Scan methods on enum types.
// The scanFunc parameter allows calling Scan on a pointer to T.
func testScanAllTypes[T comparable](
	t *testing.T,
	tests []scanTestCase[T],
	scanFunc func(*T, any) error,
) {
	t.Helper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var e T

			err := scanFunc(&e, tt.src)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if e != tt.want {
				t.Errorf("expected %v, got %v", tt.want, e)
			}
		})
	}
}

// enumScanTestCase defines input/output for Scan method testing.
// The same enum value can be tested with different input types.
type enumScanTestCase[T any] struct {
	intVal int64
	strVal string
	want   T
}

// makeScanTestCases generates all 8 scanTestCase variants from the provided test cases.
// This eliminates duplicate test case definitions for each enum type.
func makeScanTestCases[T comparable](cases []enumScanTestCase[T]) []scanTestCase[T] {
	result := make([]scanTestCase[T], 0, len(cases)*2+1)
	for _, c := range cases {
		result = append(result,
			scanTestCase[T]{"int64", c.intVal, c.want},
			scanTestCase[T]{"string", c.strVal, c.want},
			scanTestCase[T]{"bytes", []byte(c.strVal), c.want},
			scanTestCase[T]{"int", int(c.intVal), c.want},
			scanTestCase[T]{"uint", uint(c.intVal), c.want},
			scanTestCase[T]{"uint64", uint64(c.intVal), c.want},
			scanTestCase[T]{"float64", float64(c.intVal), c.want},
		)
	}

	return append(result, scanTestCase[T]{"nil", nil, *new(T)})
}

// Test comprehensive Scan types for all enums.
func TestAllEnumScanAllTypes(t *testing.T) {
	t.Parallel()

	t.Run("ActorKind", func(t *testing.T) {
		t.Parallel()
		// ActorKind: User=0, Bot=1, System=2, Service=3
		testScanAllTypes(t, makeScanTestCases([]enumScanTestCase[enums.ActorKind]{
			{0, "User", enums.ActorKindUser},
			{1, "Bot", enums.ActorKindBot},
			{2, "System", enums.ActorKindSystem},
			{3, "Service", enums.ActorKindService},
		}), (*enums.ActorKind).Scan)
	})

	t.Run("Priority", func(t *testing.T) {
		t.Parallel()
		// Priority: Low=0, Medium=1, High=2, Critical=3
		testScanAllTypes(t, makeScanTestCases([]enumScanTestCase[enums.Priority]{
			{0, "Low", enums.PriorityLow},
			{1, "Medium", enums.PriorityMedium},
			{2, "High", enums.PriorityHigh},
			{3, "Critical", enums.PriorityCritical},
		}), (*enums.Priority).Scan)
	})

	t.Run("Status", func(t *testing.T) {
		t.Parallel()
		// Status: Draft=0, Active=1, Paused=2, Archived=3, Deleted=4
		testScanAllTypes(t, makeScanTestCases([]enumScanTestCase[enums.Status]{
			{0, "Draft", enums.StatusDraft},
			{1, "Active", enums.StatusActive},
			{2, "Paused", enums.StatusPaused},
			{3, "Archived", enums.StatusArchived},
			{4, "Deleted", enums.StatusDeleted},
		}), (*enums.Status).Scan)
	})

	t.Run("Trigger", func(t *testing.T) {
		t.Parallel()
		// Trigger: Manual=0, Scheduled=1, Webhook=2, Import=3, Migration=4, System=5, Correction=6
		testScanAllTypes(t, makeScanTestCases([]enumScanTestCase[enums.Trigger]{
			{0, "Manual", enums.TriggerManual},
			{1, "Scheduled", enums.TriggerScheduled},
			{2, "Webhook", enums.TriggerWebhook},
			{3, "Import", enums.TriggerImport},
			{4, "Migration", enums.TriggerMigration},
			{5, "System", enums.TriggerSystem},
			{6, "Correction", enums.TriggerCorrection},
		}), (*enums.Trigger).Scan)
	})

	t.Run("CauseKind", func(t *testing.T) {
		t.Parallel()
		// CauseKind: Direct=0, Command=1, Event=2
		testScanAllTypes(t, makeScanTestCases([]enumScanTestCase[enums.CauseKind]{
			{0, "Direct", enums.CauseKindDirect},
			{1, "Command", enums.CauseKindCommand},
			{2, "Event", enums.CauseKindEvent},
		}), (*enums.CauseKind).Scan)
	})
}

// TestEnumScanPointerTypes tests Scan with pointer and same-type inputs
// for all enum types, covering cases not exercised by makeScanTestCases.
func TestEnumScanPointerTypes(t *testing.T) {
	t.Parallel()

	t.Run("ActorKind", func(t *testing.T) {
		t.Parallel()
		testEnumScanPointerCases(t, enums.ActorKindUser, enums.ActorKindBot, "Bot",
			(*enums.ActorKind).Scan)
	})

	t.Run("CauseKind", func(t *testing.T) {
		t.Parallel()
		testEnumScanPointerCases(t, enums.CauseKindDirect, enums.CauseKindCommand, "Command",
			(*enums.CauseKind).Scan)
	})

	t.Run("Priority", func(t *testing.T) {
		t.Parallel()
		testEnumScanPointerCases(t, enums.PriorityLow, enums.PriorityHigh, "High",
			(*enums.Priority).Scan)
	})

	t.Run("Status", func(t *testing.T) {
		t.Parallel()
		testEnumScanPointerCases(t, enums.StatusDraft, enums.StatusActive, "Active",
			(*enums.Status).Scan)
	})

	t.Run("Trigger", func(t *testing.T) {
		t.Parallel()
		testEnumScanPointerCases(t, enums.TriggerManual, enums.TriggerWebhook, "Webhook",
			(*enums.Trigger).Scan)
	})
}

func testEnumScanPointerCases[T comparable](
	t *testing.T,
	zeroVal, nonZeroVal T,
	validStr string,
	scanFunc func(*T, any) error,
) {
	t.Helper()

	t.Run("same_type", func(t *testing.T) {
		t.Parallel()

		var got T

		err := scanFunc(&got, nonZeroVal)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != nonZeroVal {
			t.Errorf("expected %v, got %v", nonZeroVal, got)
		}
	})

	t.Run("pointer_non_nil", func(t *testing.T) {
		t.Parallel()

		val := nonZeroVal

		var got T

		err := scanFunc(&got, &val)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != nonZeroVal {
			t.Errorf("expected %v, got %v", nonZeroVal, got)
		}
	})

	t.Run("pointer_nil", func(t *testing.T) {
		t.Parallel()

		var got T

		err := scanFunc(&got, (*T)(nil))
		if err == nil {
			t.Error("expected error for nil pointer")
		}
	})

	t.Run("pointer_int_non_nil", func(t *testing.T) {
		t.Parallel()

		v := int(1)

		var got T

		err := scanFunc(&got, &v)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got == zeroVal {
			t.Error("expected non-zero value")
		}
	})

	t.Run("pointer_int_nil", func(t *testing.T) {
		t.Parallel()

		var got T

		err := scanFunc(&got, (*int)(nil))
		if err == nil {
			t.Error("expected error for nil *int")
		}
	})

	t.Run("pointer_int64_non_nil", func(t *testing.T) {
		t.Parallel()

		v := int64(1)

		var got T

		err := scanFunc(&got, &v)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got == zeroVal {
			t.Error("expected non-zero value")
		}
	})

	t.Run("pointer_int64_nil", func(t *testing.T) {
		t.Parallel()

		var got T

		err := scanFunc(&got, (*int64)(nil))
		if err == nil {
			t.Error("expected error for nil *int64")
		}
	})

	t.Run("pointer_float64_non_nil", func(t *testing.T) {
		t.Parallel()

		v := float64(1)

		var got T

		err := scanFunc(&got, &v)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got == zeroVal {
			t.Error("expected non-zero value")
		}
	})

	t.Run("pointer_float64_nil", func(t *testing.T) {
		t.Parallel()

		var got T

		err := scanFunc(&got, (*float64)(nil))
		if err == nil {
			t.Error("expected error for nil *float64")
		}
	})

	t.Run("pointer_uint_non_nil", func(t *testing.T) {
		t.Parallel()

		v := uint(1)

		var got T

		err := scanFunc(&got, &v)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got == zeroVal {
			t.Error("expected non-zero value")
		}
	})

	t.Run("pointer_uint_nil", func(t *testing.T) {
		t.Parallel()

		var got T

		err := scanFunc(&got, (*uint)(nil))
		if err == nil {
			t.Error("expected error for nil *uint")
		}
	})

	t.Run("pointer_uint64_non_nil", func(t *testing.T) {
		t.Parallel()

		v := uint64(1)

		var got T

		err := scanFunc(&got, &v)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got == zeroVal {
			t.Error("expected non-zero value")
		}
	})

	t.Run("pointer_uint64_nil", func(t *testing.T) {
		t.Parallel()

		var got T

		err := scanFunc(&got, (*uint64)(nil))
		if err == nil {
			t.Error("expected error for nil *uint64")
		}
	})

	t.Run("pointer_string_non_nil", func(t *testing.T) {
		t.Parallel()

		v := validStr

		var got T

		err := scanFunc(&got, &v)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got == zeroVal {
			t.Error("expected non-zero value")
		}
	})

	t.Run("pointer_string_nil", func(t *testing.T) {
		t.Parallel()

		var got T

		err := scanFunc(&got, (*string)(nil))
		if err == nil {
			t.Error("expected error for nil *string")
		}
	})
}
