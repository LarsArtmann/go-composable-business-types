package enums_test

import (
	"fmt"
	"testing"

	"github.com/larsartmann/go-composable-business-types/testutil"
)

// Test SQL Scanner/Valuer interfaces.
func TestActorKindSQL(t *testing.T) {
	t.Parallel()
	// Test Value
	val, err := ActorKindUser.Value()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if val != "User" {
		t.Errorf("expected User, got %v", val)
	}

	// Test Scan with string
	var ak ActorKind
	if err := ak.Scan("Bot"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if ak != ActorKindBot {
		t.Errorf("expected Bot, got %v", ak)
	}

	// Test Scan with []byte
	var ak2 ActorKind
	if err := ak2.Scan([]byte("System")); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if ak2 != ActorKindSystem {
		t.Errorf("expected System, got %v", ak2)
	}

	// Test Scan with nil
	ak3 := ActorKindUser
	if err := ak3.Scan(nil); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if ak3 != 0 {
		t.Errorf("expected zero value, got %v", ak3)
	}

	// Test Scan with int (supported type - sets raw value)
	var ak4 ActorKind
	if err := ak4.Scan(2); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if ak4 != ActorKindSystem {
		t.Errorf("expected ActorKindSystem(2), got %v", ak4)
	}

	// Test Scan with int64 (supported type)
	var ak5 ActorKind
	if err := ak5.Scan(int64(3)); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if ak5 != ActorKindService {
		t.Errorf("expected ActorKindService(3), got %v", ak5)
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

func TestPrioritySQL(t *testing.T) {
	t.Parallel()
	testEnumSQL(t, []enumSQLCase[Priority]{
		{value: PriorityHigh, valueStr: "High", scanStr: "Critical", scanWant: PriorityCritical},
	}, func(p Priority) (string, error) {
		b, e := p.MarshalText()

		return string(b), e
	}, (*Priority).Scan)
}

func TestStatusSQL(t *testing.T) {
	t.Parallel()
	testEnumSQL(t, []enumSQLCase[Status]{
		{value: StatusActive, valueStr: "Active", scanStr: "Archived", scanWant: StatusArchived},
	}, func(s Status) (string, error) {
		b, e := s.MarshalText()

		return string(b), e
	}, (*Status).Scan)
}

func TestTriggerSQL(t *testing.T) {
	t.Parallel()
	testEnumSQL(t, []enumSQLCase[Trigger]{
		{
			value:    TriggerWebhook,
			valueStr: "Webhook",
			scanStr:  "Migration",
			scanWant: TriggerMigration,
		},
	}, func(tr Trigger) (string, error) {
		b, e := tr.MarshalText()

		return string(b), e
	}, (*Trigger).Scan)
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

func TestPriorityMarshal(t *testing.T) {
	t.Parallel()
	testMarshalUnmarshal(t, []enumMarshalCase[Priority]{
		{
			marshalValue:  PriorityLow,
			marshalStr:    "Low",
			unmarshalStr:  "Medium",
			unmarshalWant: PriorityMedium,
		},
	})
}

func TestStatusMarshal(t *testing.T) {
	t.Parallel()
	testMarshalUnmarshal(t, []enumMarshalCase[Status]{
		{
			marshalValue:  StatusDeleted,
			marshalStr:    "Deleted",
			unmarshalStr:  "Draft",
			unmarshalWant: StatusDraft,
		},
	})
}

func TestTriggerMarshal(t *testing.T) {
	t.Parallel()
	testMarshalUnmarshal(t, []enumMarshalCase[Trigger]{
		{
			marshalValue:  TriggerCorrection,
			marshalStr:    "Correction",
			unmarshalStr:  "Scheduled",
			unmarshalWant: TriggerScheduled,
		},
	})
}

func TestCauseKindSQL(t *testing.T) {
	t.Parallel()

	val, err := CauseKindDirect.Value()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if val != "Direct" {
		t.Errorf("expected Direct, got %v", val)
	}

	var ck CauseKind
	if err := ck.Scan("Command"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if ck != CauseKindCommand {
		t.Errorf("expected Command, got %v", ck)
	}

	// Test Scan with []byte
	var ck2 CauseKind
	if err := ck2.Scan([]byte("Event")); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if ck2 != CauseKindEvent {
		t.Errorf("expected Event, got %v", ck2)
	}

	// Test Scan with nil
	ck3 := CauseKindDirect
	if err := ck3.Scan(nil); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if ck3 != 0 {
		t.Errorf("expected zero value, got %v", ck3)
	}

	// Test Scan with int
	var ck4 CauseKind
	if err := ck4.Scan(1); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if ck4 != CauseKindCommand {
		t.Errorf("expected CauseKindCommand(1), got %v", ck4)
	}
}

func TestCauseKindMarshal(t *testing.T) {
	t.Parallel()
	testMarshalUnmarshal(t, []enumMarshalCase[CauseKind]{
		{
			marshalValue:  CauseKindEvent,
			marshalStr:    "Event",
			unmarshalStr:  "Direct",
			unmarshalWant: CauseKindDirect,
		},
	})
}

func TestActorKindAppendText(t *testing.T) {
	t.Parallel()
	testutil.RunAppendTextTest(
		t,
		"User",
		func(v ActorKind) ([]byte, error) { return v.AppendText(nil) },
		ActorKindUser,
	)
}

func TestPriorityAppendText(t *testing.T) {
	t.Parallel()
	testutil.RunAppendTextTest(
		t,
		"High",
		func(v Priority) ([]byte, error) { return v.AppendText(nil) },
		PriorityHigh,
	)
}

func TestStatusAppendText(t *testing.T) {
	t.Parallel()
	testutil.RunAppendTextTest(
		t,
		"Active",
		func(v Status) ([]byte, error) { return v.AppendText(nil) },
		StatusActive,
	)
}

func TestTriggerAppendText(t *testing.T) {
	t.Parallel()
	testutil.RunAppendTextTest(
		t,
		"Webhook",
		func(v Trigger) ([]byte, error) { return v.AppendText(nil) },
		TriggerWebhook,
	)
}

func TestCauseKindAppendText(t *testing.T) {
	t.Parallel()
	testutil.RunAppendTextTest(
		t,
		"Direct",
		func(v CauseKind) ([]byte, error) { return v.AppendText(nil) },
		CauseKindDirect,
	)
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
		t.Parallel()
		testInvalidEnumString(t, ActorKind(99), "ActorKind", 99)
	})
	t.Run("Priority", func(t *testing.T) {
		t.Parallel()
		testInvalidEnumString(t, Priority(99), "Priority", 99)
	})
	t.Run("Status", func(t *testing.T) {
		t.Parallel()
		testInvalidEnumString(t, Status(99), "Status", 99)
	})
	t.Run("Trigger", func(t *testing.T) {
		t.Parallel()
		testInvalidEnumString(t, Trigger(99), "Trigger", 99)
	})
	t.Run("CauseKind", func(t *testing.T) {
		t.Parallel()
		testInvalidEnumString(t, CauseKind(99), "CauseKind", 99)
	})
}

// Test UnmarshalText error cases.
func TestUnmarshalTextErrors(t *testing.T) {
	testUnmarshalTextErrorsAll(t,
		[]enumUnmarshalTextErrorCase[ActorKind]{{"ActorKind"}},
		[]enumUnmarshalTextErrorCase[Priority]{{"Priority"}},
		[]enumUnmarshalTextErrorCase[Status]{{"Status"}},
		[]enumUnmarshalTextErrorCase[Trigger]{{"Trigger"}},
		[]enumUnmarshalTextErrorCase[CauseKind]{{"CauseKind"}},
	)
}

// Test all Value methods.
func TestValueMethods(t *testing.T) {
	t.Parallel()
	testEnumValue(t, []enumValueCase[ActorKind]{
		{ActorKindSystem, "System"},
	})
	testEnumValue(t, []enumValueCase[Priority]{
		{PriorityMedium, "Medium"},
	})
	testEnumValue(t, []enumValueCase[Status]{
		{StatusPaused, "Paused"},
	})
	testEnumValue(t, []enumValueCase[Trigger]{
		{TriggerImport, "Import"},
	})
	testEnumValue(t, []enumValueCase[CauseKind]{
		{CauseKindEvent, "Event"},
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

// Test comprehensive Scan types for all enums.
func TestAllEnumScanAllTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		typeName string
		cases    []any
		scanFunc any
	}{
		{
			"ActorKind",
			[]any{
				scanTestCase[ActorKind]{"int64", int64(1), ActorKindBot},
				scanTestCase[ActorKind]{"string", "System", ActorKindSystem},
				scanTestCase[ActorKind]{"bytes", []byte("Service"), ActorKindService},
				scanTestCase[ActorKind]{"int", int(0), ActorKindUser},
				scanTestCase[ActorKind]{"uint", uint(2), ActorKindSystem},
				scanTestCase[ActorKind]{"uint64", uint64(3), ActorKindService},
				scanTestCase[ActorKind]{"float64", float64(1), ActorKindBot},
				scanTestCase[ActorKind]{"nil", nil, ActorKind(0)},
			},
			(*ActorKind).Scan,
		},
		{
			"Priority",
			[]any{
				scanTestCase[Priority]{"int64", int64(2), PriorityHigh},
				scanTestCase[Priority]{"string", "Critical", PriorityCritical},
				scanTestCase[Priority]{"bytes", []byte("Low"), PriorityLow},
				scanTestCase[Priority]{"int", int(1), PriorityMedium},
				scanTestCase[Priority]{"uint", uint(0), PriorityLow},
				scanTestCase[Priority]{"uint64", uint64(3), PriorityCritical},
				scanTestCase[Priority]{"float64", float64(2), PriorityHigh},
				scanTestCase[Priority]{"nil", nil, Priority(0)},
			},
			(*Priority).Scan,
		},
		{
			"Status",
			[]any{
				scanTestCase[Status]{"int64", int64(1), StatusActive},
				scanTestCase[Status]{"string", "Archived", StatusArchived},
				scanTestCase[Status]{"bytes", []byte("Deleted"), StatusDeleted},
				scanTestCase[Status]{"int", int(0), StatusDraft},
				scanTestCase[Status]{"uint", uint(2), StatusPaused},
				scanTestCase[Status]{"uint64", uint64(4), StatusDeleted},
				scanTestCase[Status]{"float64", float64(1), StatusActive},
				scanTestCase[Status]{"nil", nil, Status(0)},
			},
			(*Status).Scan,
		},
		{
			"Trigger",
			[]any{
				scanTestCase[Trigger]{"int64", int64(2), TriggerWebhook},
				scanTestCase[Trigger]{"string", "Correction", TriggerCorrection},
				scanTestCase[Trigger]{"bytes", []byte("Import"), TriggerImport},
				scanTestCase[Trigger]{"int", int(0), TriggerManual},
				scanTestCase[Trigger]{"uint", uint(5), TriggerSystem},
				scanTestCase[Trigger]{"uint64", uint64(6), TriggerCorrection},
				scanTestCase[Trigger]{"float64", float64(2), TriggerWebhook},
				scanTestCase[Trigger]{"nil", nil, Trigger(0)},
			},
			(*Trigger).Scan,
		},
		{
			"CauseKind",
			[]any{
				scanTestCase[CauseKind]{"int64", int64(1), CauseKindCommand},
				scanTestCase[CauseKind]{"string", "Event", CauseKindEvent},
				scanTestCase[CauseKind]{"bytes", []byte("Direct"), CauseKindDirect},
				scanTestCase[CauseKind]{"int", int(0), CauseKindDirect},
				scanTestCase[CauseKind]{"uint", uint(2), CauseKindEvent},
				scanTestCase[CauseKind]{"uint64", uint64(1), CauseKindCommand},
				scanTestCase[CauseKind]{"float64", float64(1), CauseKindCommand},
				scanTestCase[CauseKind]{"nil", nil, CauseKind(0)},
			},
			(*CauseKind).Scan,
		},
	}

	for _, tt := range tests {
		t.Run(tt.typeName, func(t *testing.T) {
			switch f := tt.scanFunc.(type) {
			case func(*ActorKind, any) error:
				cases := make([]scanTestCase[ActorKind], len(tt.cases))
				for i, c := range tt.cases {
					cases[i] = c.(scanTestCase[ActorKind])
				}

				testScanAllTypes(t, cases, f)
			case func(*Priority, any) error:
				cases := make([]scanTestCase[Priority], len(tt.cases))
				for i, c := range tt.cases {
					cases[i] = c.(scanTestCase[Priority])
				}

				testScanAllTypes(t, cases, f)
			case func(*Status, any) error:
				cases := make([]scanTestCase[Status], len(tt.cases))
				for i, c := range tt.cases {
					cases[i] = c.(scanTestCase[Status])
				}

				testScanAllTypes(t, cases, f)
			case func(*Trigger, any) error:
				cases := make([]scanTestCase[Trigger], len(tt.cases))
				for i, c := range tt.cases {
					cases[i] = c.(scanTestCase[Trigger])
				}

				testScanAllTypes(t, cases, f)
			case func(*CauseKind, any) error:
				cases := make([]scanTestCase[CauseKind], len(tt.cases))
				for i, c := range tt.cases {
					cases[i] = c.(scanTestCase[CauseKind])
				}

				testScanAllTypes(t, cases, f)
			}
		})
	}
}
