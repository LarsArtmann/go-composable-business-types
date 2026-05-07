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

func TestActorKindIsValid(t *testing.T) {
	t.Parallel()

	if !ActorKindUser.IsValid() {
		t.Error("ActorKindUser should be valid")
	}

	if !ActorKindService.IsValid() {
		t.Error("ActorKindService should be valid")
	}

	if ActorKind(99).IsValid() {
		t.Error("ActorKind(99) should not be valid")
	}
}

func TestActorKindMarshalText(t *testing.T) {
	t.Parallel()

	data, err := ActorKindUser.MarshalText()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(data) != "User" {
		t.Errorf("expected User, got %s", string(data))
	}

	var got ActorKind
	if err := got.UnmarshalText([]byte("Bot")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != ActorKindBot {
		t.Errorf("expected Bot, got %v", got)
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

func TestPriorityIsValid(t *testing.T) {
	t.Parallel()

	if !PriorityLow.IsValid() {
		t.Error("PriorityLow should be valid")
	}

	if !PriorityCritical.IsValid() {
		t.Error("PriorityCritical should be valid")
	}

	if Priority(99).IsValid() {
		t.Error("Priority(99) should not be valid")
	}
}
