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

// TestAppendText tests AppendText for all enum types.
func TestAppendText(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		value any
	}{
		{"ActorKindUser", ActorKindUser},
		{"PriorityHigh", PriorityHigh},
		{"StatusActive", StatusActive},
		{"TriggerWebhook", TriggerWebhook},
		{"CauseKindDirect", CauseKindDirect},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			switch v := tt.value.(type) {
			case ActorKind:
				testutil.RunAppendTextTest(t, tt.name, func(ActorKind) ([]byte, error) { return v.AppendText(nil) }, v)
			case Priority:
				testutil.RunAppendTextTest(t, tt.name, func(Priority) ([]byte, error) { return v.AppendText(nil) }, v)
			case Status:
				testutil.RunAppendTextTest(t, tt.name, func(Status) ([]byte, error) { return v.AppendText(nil) }, v)
			case Trigger:
				testutil.RunAppendTextTest(t, tt.name, func(Trigger) ([]byte, error) { return v.AppendText(nil) }, v)
			case CauseKind:
				testutil.RunAppendTextTest(t, tt.name, func(CauseKind) ([]byte, error) { return v.AppendText(nil) }, v)
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

// convertAndTestScan converts []any cases to []scanTestCase[T] and runs tests.
func convertAndTestScan[T comparable](t *testing.T, cases []any, scanFunc func(*T, any) error) {
	converted := make([]scanTestCase[T], len(cases))
	for i, c := range cases {
		converted[i] = c.(scanTestCase[T])
	}
	testScanAllTypes(t, converted, scanFunc)
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
			scanTestCase[T]{"int64", int64(c.intVal), c.want},
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

	tests := []struct {
		typeName string
		cases    []any
		scanFunc any
	}{
		{
			"ActorKind",
			makeScanTestCases([]enumScanTestCase[ActorKind]{
				{1, "System", ActorKindSystem},
				{2, "Service", ActorKindService},
				{0, "User", ActorKindUser},
				{3, "Bot", ActorKindBot},
			}),
			(*ActorKind).Scan,
		},
		{
			"Priority",
			makeScanTestCases([]enumScanTestCase[Priority]{
				{2, "Critical", PriorityCritical},
				{0, "Low", PriorityLow},
				{1, "Medium", PriorityMedium},
				{3, "High", PriorityHigh},
			}),
			(*Priority).Scan,
		},
		{
			"Status",
			makeScanTestCases([]enumScanTestCase[Status]{
				{1, "Archived", StatusArchived},
				{4, "Deleted", StatusDeleted},
				{0, "Draft", StatusDraft},
				{2, "Paused", StatusPaused},
				{3, "Active", StatusActive},
			}),
			(*Status).Scan,
		},
		{
			"Trigger",
			makeScanTestCases([]enumScanTestCase[Trigger]{
				{2, "Correction", TriggerCorrection},
				{0, "Import", TriggerImport},
				{5, "System", TriggerSystem},
				{6, "Webhook", TriggerWebhook},
			}),
			(*Trigger).Scan,
		},
		{
			"CauseKind",
			makeScanTestCases([]enumScanTestCase[CauseKind]{
				{1, "Event", CauseKindEvent},
				{0, "Direct", CauseKindDirect},
				{2, "Command", CauseKindCommand},
			}),
			(*CauseKind).Scan,
		},
	}

	for _, tt := range tests {
		t.Run(tt.typeName, func(t *testing.T) {
			switch f := tt.scanFunc.(type) {
			case func(*ActorKind, any) error:
				convertAndTestScan(t, tt.cases, f)
			case func(*Priority, any) error:
				convertAndTestScan(t, tt.cases, f)
			case func(*Status, any) error:
				convertAndTestScan(t, tt.cases, f)
			case func(*Trigger, any) error:
				convertAndTestScan(t, tt.cases, f)
			case func(*CauseKind, any) error:
				convertAndTestScan(t, tt.cases, f)
			}
		})
	}
}
