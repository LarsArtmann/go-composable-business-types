package enums

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
	tests := []struct {
		input   string
		want    Status
		wantErr bool
	}{
		{"Draft", StatusDraft, false},
		{"Active", StatusActive, false},
		{"Paused", StatusPaused, false},
		{"Archived", StatusArchived, false},
		{"Deleted", StatusDeleted, false},
		{"Invalid", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			got, err := ParseStatus(tt.input)
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

func TestStatusIsValid(t *testing.T) {
	t.Parallel()
	validStatuses := []Status{
		StatusDraft,
		StatusActive,
		StatusPaused,
		StatusArchived,
		StatusDeleted,
	}
	for _, s := range validStatuses {
		if !s.IsValid() {
			t.Errorf("Status %v should be valid", s)
		}
	}
	invalid := Status(99)
	if invalid.IsValid() {
		t.Error("Status(99) should not be valid")
	}
}

func TestTrigger(t *testing.T) {
	t.Parallel()
	// Test all triggers
	tests := []struct {
		trigger  Trigger
		expected string
	}{
		{TriggerManual, "Manual"},
		{TriggerScheduled, "Scheduled"},
		{TriggerWebhook, "Webhook"},
		{TriggerImport, "Import"},
		{TriggerMigration, "Migration"},
		{TriggerSystem, "System"},
		{TriggerCorrection, "Correction"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			t.Parallel()
			if tt.trigger.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, tt.trigger.String())
			}
		})
	}
}

func TestParseTrigger(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input   string
		want    Trigger
		wantErr bool
	}{
		{"Manual", TriggerManual, false},
		{"Scheduled", TriggerScheduled, false},
		{"Webhook", TriggerWebhook, false},
		{"Import", TriggerImport, false},
		{"Migration", TriggerMigration, false},
		{"System", TriggerSystem, false},
		{"Correction", TriggerCorrection, false},
		{"Invalid", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			got, err := ParseTrigger(tt.input)
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

func TestTriggerIsValid(t *testing.T) {
	t.Parallel()
	validTriggers := []Trigger{
		TriggerManual, TriggerScheduled, TriggerWebhook, TriggerImport,
		TriggerMigration, TriggerSystem, TriggerCorrection,
	}
	for _, tr := range validTriggers {
		if !tr.IsValid() {
			t.Errorf("Trigger %v should be valid", tr)
		}
	}
	invalid := Trigger(99)
	if invalid.IsValid() {
		t.Error("Trigger(99) should not be valid")
	}
}
