package enums

import (
	"testing"

	"github.com/larsartmann/go-composable-business-types/testutil"
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
	testEnumString(t, []enumStringCase[ActorKind]{
		{ActorKindUser, "User"},
		{ActorKindBot, "Bot"},
		{ActorKindSystem, "System"},
		{ActorKindService, "Service"},
	})
}

func TestParseActorKind(t *testing.T) {
	t.Parallel()

	testParse(t, ParseActorKind, []enumParseCase[ActorKind]{
		{"User", ActorKindUser, false},
		{"Bot", ActorKindBot, false},
		{"System", ActorKindSystem, false},
		{"Service", ActorKindService, false},
		{"Invalid", 0, true},
		{"", 0, true},
	})
}

func TestParseActorKindError(t *testing.T) {
	t.Parallel()
	testutil.RunParseTest(t, "ActorKind", "Invalid", ParseActorKind, true)
}

func TestActorKindIsValid(t *testing.T) {
	t.Parallel()
	testEnumIsValid(t, []ActorKind{ActorKindUser, ActorKindService}, ActorKind(99))
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
	testEnumString(t, []enumStringCase[Priority]{
		{PriorityLow, "Low"},
		{PriorityMedium, "Medium"},
		{PriorityHigh, "High"},
		{PriorityCritical, "Critical"},
	})
}

func TestParsePriority(t *testing.T) {
	t.Parallel()

	testParse(t, ParsePriority, []enumParseCase[Priority]{
		{"Low", PriorityLow, false},
		{"Medium", PriorityMedium, false},
		{"High", PriorityHigh, false},
		{"Critical", PriorityCritical, false},
		{"Invalid", 0, true},
		{"", 0, true},
	})
}

func TestPriorityIsValid(t *testing.T) {
	t.Parallel()

	testEnumIsValid(t, []Priority{PriorityLow, PriorityCritical}, Priority(99))
}
