package enums

import (
	"slices"
	"testing"
)

func TestActorKind(t *testing.T) {
	// Test constants exist
	_ = ActorKindUser
	_ = ActorKindBot
	_ = ActorKindSystem
	_ = ActorKindService

	// Test String()
	if ActorKindUser.String() != "User" {
		t.Errorf("expected User, got %s", ActorKindUser.String())
	}
	if ActorKindSystem.String() != "System" {
		t.Errorf("expected System, got %s", ActorKindSystem.String())
	}
}

func TestParseActorKind(t *testing.T) {
	tests := []struct {
		input   string
		want    ActorKind
		wantErr bool
	}{
		{"User", ActorKindUser, false},
		{"Bot", ActorKindBot, false},
		{"System", ActorKindSystem, false},
		{"Service", ActorKindService, false},
		{"Invalid", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseActorKind(tt.input)
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

func TestPriority(t *testing.T) {
	// Test constants
	_ = PriorityLow
	_ = PriorityMedium
	_ = PriorityHigh
	_ = PriorityCritical

	// Test ordering
	if PriorityLow >= PriorityHigh {
		t.Error("Low should be less than High")
	}

	// Test String()
	if PriorityCritical.String() != "Critical" {
		t.Errorf("expected Critical, got %s", PriorityCritical.String())
	}
}

func TestStatus(t *testing.T) {
	// Test all statuses exist
	statuses := []Status{StatusDraft, StatusActive, StatusPaused, StatusArchived, StatusDeleted}
	for _, s := range statuses {
		if s.String() == "" {
			t.Errorf("Status %v should have a string representation", s)
		}
	}
}

func TestTrigger(t *testing.T) {
	// Test key triggers
	_ = TriggerManual
	_ = TriggerScheduled
	_ = TriggerSystem

	// Test String()
	if TriggerManual.String() != "Manual" {
		t.Errorf("expected Manual, got %s", TriggerManual.String())
	}
}

func TestEnumValues(t *testing.T) {
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
}

func TestEnumNames(t *testing.T) {
	// Test ActorKindNames
	akNames := ActorKindNames()
	if len(akNames) != 4 {
		t.Errorf("expected 4 ActorKind names, got %d", len(akNames))
	}

	// Check specific names exist
	if !slices.Contains(akNames, "User") {
		t.Error("ActorKindNames should include 'User'")
	}
}
