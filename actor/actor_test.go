package actor

import (
	"testing"

	"github.com/larsartmann/go-composable-business-types/enums"
	"github.com/larsartmann/go-composable-business-types/id"
)

// newTestChain creates a test actor chain with a user and two services.
func newTestChain(userName string) ActorChain[string] {
	return NewActorChain(UserActor(id.NewID[struct{}, string]("u-1"), userName)).
		Append(ServiceActor(id.NewID[struct{}, string]("svc-1"), "Service 1")).
		Append(ServiceActor(id.NewID[struct{}, string]("svc-2"), "Service 2"))
}

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

	chain := newTestChain("Alice")

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

	chain := newTestChain("User")

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

	if actor.ID != userID {
		t.Error("ID mismatch")
	}

	if actor.Name != "John Doe" {
		t.Errorf("expected 'John Doe', got %s", actor.Name)
	}
}

func TestActorKind(t *testing.T) {
	t.Parallel()

	botID := id.NewID[struct{}, string]("bot-1")
	serviceID := id.NewID[struct{}, string]("svc-1")

	tests := []struct {
		name     string
		actor    ActorEntry[string]
		expected enums.ActorKind
	}{
		{
			name:     "Bot",
			actor:    BotActor(botID, "GitHub Bot"),
			expected: enums.ActorKindBot,
		},
		{
			name:     "System",
			actor:    SystemActor[string](),
			expected: enums.ActorKindSystem,
		},
		{
			name:     "Service",
			actor:    ServiceActor(serviceID, "Order Service"),
			expected: enums.ActorKindService,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.actor.Kind != tt.expected {
				t.Errorf("expected %v kind, got %v", tt.expected, tt.actor.Kind)
			}
		})
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

func TestActorChainAll(t *testing.T) {
	t.Parallel()

	chain := newTestChain("Alice")

	indices := make([]int, 0, len(chain))

	names := make([]string, 0, len(chain))
	for i, e := range chain.All() {
		indices = append(indices, i)
		names = append(names, e.Name)
	}

	if len(indices) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(indices))
	}

	for i, idx := range indices {
		if idx != i {
			t.Errorf("expected index %d, got %d", i, idx)
		}
	}

	expected := []string{"Alice", "Service 1", "Service 2"}
	for i, name := range expected {
		if names[i] != name {
			t.Errorf("entry[%d]: expected %q, got %q", i, name, names[i])
		}
	}
}

func TestActorChainIterationBreak(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		iter func(chain ActorChain[string]) func(func(ActorEntry[string], int) bool)
	}{
		{"All", func(chain ActorChain[string]) func(func(ActorEntry[string], int) bool) {
			return func(yield func(ActorEntry[string], int) bool) {
				for i, e := range chain.All() {
					if !yield(e, i) {
						return
					}
				}
			}
		}},
		{"Entries", func(chain ActorChain[string]) func(func(ActorEntry[string], int) bool) {
			return func(yield func(ActorEntry[string], int) bool) {
				for e := range chain.Entries() {
					if !yield(e, 0) {
						return
					}
				}
			}
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			chain := newTestChain("Alice")

			var count int
			for range tt.iter(chain) {
				count++
				break
			}

			if count != 1 {
				t.Errorf("expected break after 1 iteration, got %d", count)
			}
		})
	}
}

func TestActorChainEntries(t *testing.T) {
	t.Parallel()

	chain := newTestChain("Alice")

	names := make([]string, 0, len(chain))
	for e := range chain.Entries() {
		names = append(names, e.Name)
	}

	expected := []string{"Alice", "Service 1", "Service 2"}
	if len(names) != len(expected) {
		t.Fatalf("expected %d entries, got %d", len(expected), len(names))
	}

	for i, name := range expected {
		if names[i] != name {
			t.Errorf("entry[%d]: expected %q, got %q", i, name, names[i])
		}
	}
}

func TestActorChainAllEmpty(t *testing.T) {
	t.Parallel()

	var (
		chain ActorChain[string]
		count int
	)

	for range chain.All() {
		count++
	}

	if count != 0 {
		t.Errorf("expected 0 iterations on empty chain, got %d", count)
	}
}
