package enums

import "testing"

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

	// Test String()
	if PriorityLow.String() != "Low" {
		t.Errorf("expected Low, got %s", PriorityLow.String())
	}
}
