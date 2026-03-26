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
	tests := []struct {
		kind     CauseKind
		expected string
	}{
		{CauseKindDirect, "Direct"},
		{CauseKindCommand, "Command"},
		{CauseKindEvent, "Event"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			t.Parallel()
			if tt.kind.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, tt.kind.String())
			}
		})
	}
}

func TestParseCauseKind(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input   string
		want    CauseKind
		wantErr bool
	}{
		{"Direct", CauseKindDirect, false},
		{"Command", CauseKindCommand, false},
		{"Event", CauseKindEvent, false},
		{"Invalid", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			got, err := ParseCauseKind(tt.input)
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

func TestCauseKindIsValid(t *testing.T) {
	t.Parallel()
	if !CauseKindDirect.IsValid() {
		t.Error("CauseKindDirect should be valid")
	}
	if !CauseKindEvent.IsValid() {
		t.Error("CauseKindEvent should be valid")
	}
	invalid := CauseKind(99)
	if invalid.IsValid() {
		t.Error("CauseKind(99) should not be valid")
	}
}
