package enums

import (
	"testing"
)

func TestActorKind(t *testing.T) {
	t.Parallel()
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

	// Test all ActorKind values
	tests := []struct {
		kind     ActorKind
		expected string
	}{
		{ActorKindUser, "User"},
		{ActorKindBot, "Bot"},
		{ActorKindSystem, "System"},
		{ActorKindService, "Service"},
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

func TestParseActorKind(t *testing.T) {
	t.Parallel()
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
		{"", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
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

func TestParseActorKindError(t *testing.T) {
	t.Parallel()
	_, err := ParseActorKind("Invalid")
	if err == nil {
		t.Error("expected error for invalid ActorKind")
	}
}

func TestActorKindIsValid(t *testing.T) {
	t.Parallel()
	if !ActorKindUser.IsValid() {
		t.Error("ActorKindUser should be valid")
	}
	if !ActorKindService.IsValid() {
		t.Error("ActorKindService should be valid")
	}
	// Invalid value
	invalid := ActorKind(99)
	if invalid.IsValid() {
		t.Error("ActorKind(99) should not be valid")
	}
}

func TestPriority(t *testing.T) {
	t.Parallel()
	// Test constants
	_ = PriorityLow
	_ = PriorityMedium
	_ = PriorityHigh
	_ = PriorityCritical

	// Test ordering
	if PriorityLow >= PriorityHigh {
		t.Error("Low should be less than High")
	}

	// Test String() for all values
	tests := []struct {
		priority Priority
		expected string
	}{
		{PriorityLow, "Low"},
		{PriorityMedium, "Medium"},
		{PriorityHigh, "High"},
		{PriorityCritical, "Critical"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			t.Parallel()
			if tt.priority.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, tt.priority.String())
			}
		})
	}
}

func TestParsePriority(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input   string
		want    Priority
		wantErr bool
	}{
		{"Low", PriorityLow, false},
		{"Medium", PriorityMedium, false},
		{"High", PriorityHigh, false},
		{"Critical", PriorityCritical, false},
		{"Invalid", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			got, err := ParsePriority(tt.input)
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

func TestPriorityIsValid(t *testing.T) {
	t.Parallel()
	if !PriorityLow.IsValid() {
		t.Error("PriorityLow should be valid")
	}
	if !PriorityCritical.IsValid() {
		t.Error("PriorityCritical should be valid")
	}
	invalid := Priority(99)
	if invalid.IsValid() {
		t.Error("Priority(99) should not be valid")
	}
}
