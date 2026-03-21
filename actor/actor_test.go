package actor

import (
	"testing"

	"github.com/larsartmann/go-composable-business-types/enums"
	"github.com/larsartmann/go-composable-business-types/id"
)

func TestNewActorChain(t *testing.T) {
	t.Parallel()
	userID := id.NewID[struct{}, string]("user-1")
	entry := UserActor(userID, "Alice")
	chain := NewActorChain(entry)

	if chain.IsZero() {
		t.Error("chain should not be zero")
	}
	if len(chain) != 1 {
		t.Errorf("expected length 1, got %d", len(chain))
	}
}

func TestActorChainOriginAndCurrent(t *testing.T) {
	t.Parallel()
	chain := NewActorChain(UserActor(id.NewID[struct{}, string]("u-1"), "Alice")).
		Append(ServiceActor(id.NewID[struct{}, string]("svc-1"), "Service 1")).
		Append(ServiceActor(id.NewID[struct{}, string]("svc-2"), "Service 2"))

	if chain.Origin().Name != "Alice" {
		t.Errorf("expected origin Alice, got %s", chain.Origin().Name)
	}
	if chain.Current().Name != "Service 2" {
		t.Errorf("expected current 'Service 2', got %s", chain.Current().Name)
	}
}

func TestActorChainHasKind(t *testing.T) {
	t.Parallel()
	chain := NewActorChain(UserActor(id.NewID[struct{}, string]("u-1"), "Alice")).
		Append(ServiceActor(id.NewID[struct{}, string]("svc-1"), "Service 1"))

	if !chain.HasKind(enums.ActorKindUser) {
		t.Error("chain should have user")
	}
	if !chain.HasKind(enums.ActorKindService) {
		t.Error("chain should have service")
	}
	if chain.HasKind(enums.ActorKindBot) {
		t.Error("chain should not have bot")
	}
}

func TestActorChainByKind(t *testing.T) {
	t.Parallel()
	chain := NewActorChain(UserActor(id.NewID[struct{}, string]("u-1"), "User")).
		Append(ServiceActor(id.NewID[struct{}, string]("svc-1"), "Service 1")).
		Append(ServiceActor(id.NewID[struct{}, string]("svc-2"), "Service 2"))

	services := chain.ByKind(enums.ActorKindService)
	if len(services) != 2 {
		t.Errorf("expected 2 services, got %d", len(services))
	}

	users := chain.ByKind(enums.ActorKindUser)
	if len(users) != 1 {
		t.Errorf("expected 1 user, got %d", len(users))
	}
}

func TestUserActor(t *testing.T) {
	t.Parallel()
	userID := id.NewID[struct{}, string]("user-1")
	actor := UserActor(userID, "John Doe")

	if actor.Kind != enums.ActorKindUser {
		t.Errorf("expected User kind, got %v", actor.Kind)
	}
	if actor.Id != userID {
		t.Error("ID mismatch")
	}
	if actor.Name != "John Doe" {
		t.Errorf("expected 'John Doe', got %s", actor.Name)
	}
}

func TestBotActor(t *testing.T) {
	t.Parallel()
	botID := id.NewID[struct{}, string]("bot-1")
	actor := BotActor(botID, "GitHub Bot")

	if actor.Kind != enums.ActorKindBot {
		t.Errorf("expected Bot kind, got %v", actor.Kind)
	}
}

func TestSystemActor(t *testing.T) {
	t.Parallel()
	actor := SystemActor[string]()

	if actor.Kind != enums.ActorKindSystem {
		t.Errorf("expected System kind, got %v", actor.Kind)
	}
}

func TestServiceActor(t *testing.T) {
	t.Parallel()
	serviceID := id.NewID[struct{}, string]("svc-1")
	actor := ServiceActor(serviceID, "Order Service")

	if actor.Kind != enums.ActorKindService {
		t.Errorf("expected Service kind, got %v", actor.Kind)
	}
}

func TestActorEntryIsZero(t *testing.T) {
	t.Parallel()
	var zero ActorEntry[string]
	if !zero.IsZero() {
		t.Error("zero ActorEntry should be zero")
	}

	nonZero := UserActor(id.NewID[struct{}, string]("u-1"), "User")
	if nonZero.IsZero() {
		t.Error("non-zero ActorEntry should not be zero")
	}
}

func TestActorEntryOptionalName(t *testing.T) {
	t.Parallel()
	// Actor without name
	actor := UserActor(id.NewID[struct{}, string]("u-1"))
	if actor.Name != "" {
		t.Errorf("expected empty name, got %s", actor.Name)
	}

	// Actor with name
	actorWithName := UserActor(id.NewID[struct{}, string]("u-2"), "John")
	if actorWithName.Name != "John" {
		t.Errorf("expected 'John', got %s", actorWithName.Name)
	}
}
