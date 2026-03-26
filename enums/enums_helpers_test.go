package enums

import (
	"slices"
	"testing"
)

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
}
