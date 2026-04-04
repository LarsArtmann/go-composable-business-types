package enums_test

import (
	"testing"
)

func TestStatus(t *testing.T) {
	t.Parallel()
	// Test all statuses exist
	statuses := []Status{StatusDraft, StatusActive, StatusPaused, StatusArchived, StatusDeleted}
	expected := []string{"Draft", "Active", "Paused", "Archived", "Deleted"}

	for i, s := range statuses {
		if s.String() != expected[i] {
			t.Errorf("Status %v should have string %s, got %s", s, expected[i], s.String())
		}
	}
}

func TestParseStatus(t *testing.T) {
	t.Parallel()

	testParse(t, ParseStatus, []enumParseCase[Status]{
		{"Draft", StatusDraft, false},
		{"Active", StatusActive, false},
		{"Paused", StatusPaused, false},
		{"Archived", StatusArchived, false},
		{"Deleted", StatusDeleted, false},
		{"Invalid", 0, true},
		{"", 0, true},
	})
}

func TestStatusIsValid(t *testing.T) {
	t.Parallel()
	testEnumIsValid(
		t,
		[]Status{StatusDraft, StatusActive, StatusPaused, StatusArchived, StatusDeleted},
		Status(99),
	)
}

func TestTrigger(t *testing.T) {
	t.Parallel()
	testEnumString(t, []enumStringCase[Trigger]{
		{TriggerManual, "Manual"},
		{TriggerScheduled, "Scheduled"},
		{TriggerWebhook, "Webhook"},
		{TriggerImport, "Import"},
		{TriggerMigration, "Migration"},
		{TriggerSystem, "System"},
		{TriggerCorrection, "Correction"},
	})
}

func TestParseTrigger(t *testing.T) {
	t.Parallel()

	testParse(t, ParseTrigger, []enumParseCase[Trigger]{
		{"Manual", TriggerManual, false},
		{"Scheduled", TriggerScheduled, false},
		{"Webhook", TriggerWebhook, false},
		{"Import", TriggerImport, false},
		{"Migration", TriggerMigration, false},
		{"System", TriggerSystem, false},
		{"Correction", TriggerCorrection, false},
		{"Invalid", 0, true},
		{"", 0, true},
	})
}

func TestTriggerIsValid(t *testing.T) {
	t.Parallel()
	testEnumIsValid(
		t,
		[]Trigger{
			TriggerManual,
			TriggerScheduled,
			TriggerWebhook,
			TriggerImport,
			TriggerMigration,
			TriggerSystem,
			TriggerCorrection,
		},
		Trigger(99),
	)
}
