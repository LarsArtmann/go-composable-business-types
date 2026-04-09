package enums_test

import (
	"testing"

	"github.com/larsartmann/go-composable-business-types/enums"
)

func TestStatus(t *testing.T) {
	t.Parallel()
	// Test all statuses exist
	statuses := []enums.Status{enums.StatusDraft, enums.StatusActive, enums.StatusPaused, enums.StatusArchived, enums.StatusDeleted}
	expected := []string{"Draft", "Active", "Paused", "Archived", "Deleted"}

	for i, s := range statuses {
		if s.String() != expected[i] {
			t.Errorf("Status %v should have string %s, got %s", s, expected[i], s.String())
		}
	}
}

func TestParseStatus(t *testing.T) {
	t.Parallel()

	testParse(t, enums.ParseStatus, []enumParseCase[enums.Status]{
		{"Draft", enums.StatusDraft, false},
		{"Active", enums.StatusActive, false},
		{"Paused", enums.StatusPaused, false},
		{"Archived", enums.StatusArchived, false},
		{"Deleted", enums.StatusDeleted, false},
		{"Invalid", 0, true},
		{"", 0, true},
	})
}

func TestStatusIsValid(t *testing.T) {
	t.Parallel()
	testEnumIsValid(
		t,
		[]enums.Status{enums.StatusDraft, enums.StatusActive, enums.StatusPaused, enums.StatusArchived, enums.StatusDeleted},
		enums.Status(99),
	)
}

func TestTrigger(t *testing.T) {
	t.Parallel()
	testEnumString(t, []enumStringCase[enums.Trigger]{
		{enums.TriggerManual, "Manual"},
		{enums.TriggerScheduled, "Scheduled"},
		{enums.TriggerWebhook, "Webhook"},
		{enums.TriggerImport, "Import"},
		{enums.TriggerMigration, "Migration"},
		{enums.TriggerSystem, "System"},
		{enums.TriggerCorrection, "Correction"},
	})
}

func TestParseTrigger(t *testing.T) {
	t.Parallel()

	testParse(t, enums.ParseTrigger, []enumParseCase[enums.Trigger]{
		{"Manual", enums.TriggerManual, false},
		{"Scheduled", enums.TriggerScheduled, false},
		{"Webhook", enums.TriggerWebhook, false},
		{"Import", enums.TriggerImport, false},
		{"Migration", enums.TriggerMigration, false},
		{"System", enums.TriggerSystem, false},
		{"Correction", enums.TriggerCorrection, false},
		{"Invalid", 0, true},
		{"", 0, true},
	})
}

func TestTriggerIsValid(t *testing.T) {
	t.Parallel()
	testEnumIsValid(
		t,
		[]enums.Trigger{
			enums.TriggerManual,
			enums.TriggerScheduled,
			enums.TriggerWebhook,
			enums.TriggerImport,
			enums.TriggerMigration,
			enums.TriggerSystem,
			enums.TriggerCorrection,
		},
		enums.Trigger(99),
	)
}
