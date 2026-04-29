// Package actor provides types for tracking audit trails and actor chains.
//
// ActorChain represents an ordered chain of actors (e.g., User → Service → Service)
// for tracking who performed actions in a system. Each ActorEntry contains
// an ID, name, and kind (User, Bot, System, Service).
//
// Basic usage:
//
//	chain := actor.NewActorChain(actor.UserActor(userID, "Alice")).
//	    Append(actor.ServiceActor(serviceID, "API Gateway"))
//	origin := chain.Origin()   // User Alice
//	current := chain.Current() // API Gateway
package actor

import (
	"iter"
	"slices"

	"github.com/larsartmann/go-composable-business-types/enums"
	"github.com/larsartmann/go-branded-id"
)

// ActorEntry represents a single actor in an actor chain.
//
// Fields:
//   - Kind: The type of actor (User, Bot, System, Service)
//   - ID: The unique identifier for this actor
//   - Name: Optional human-readable name
type ActorEntry[T comparable] struct {
	Kind enums.ActorKind
	ID   id.ID[struct{}, T]
	Name string // optional human-readable name
}

// ActorChain represents an ordered chain of actors from original to current.
// Index 0 = original actor, last index = current actor.
type ActorChain[T comparable] []ActorEntry[T]

// NewActorChain creates a new ActorChain with the first actor entry.
func NewActorChain[T comparable](first ActorEntry[T]) ActorChain[T] {
	return ActorChain[T]{first}
}

// Origin returns the first actor in the chain.
func (c ActorChain[T]) Origin() ActorEntry[T] { return c[0] }

// Current returns the most recent actor in the chain.
func (c ActorChain[T]) Current() ActorEntry[T] { return c[len(c)-1] }

// IsZero returns true if the chain is empty.
func (c ActorChain[T]) IsZero() bool { return len(c) == 0 }

// All returns an iterator over all actor entries in the chain, yielding index-value pairs.
// Follows the Go 1.23+ range-over-func convention.
func (c ActorChain[T]) All() iter.Seq2[int, ActorEntry[T]] {
	return func(yield func(int, ActorEntry[T]) bool) {
		for i, e := range c {
			if !yield(i, e) {
				return
			}
		}
	}
}

// Entries returns an iterator over actor entries (values only).
func (c ActorChain[T]) Entries() iter.Seq[ActorEntry[T]] {
	return func(yield func(ActorEntry[T]) bool) {
		for _, e := range c {
			if !yield(e) {
				return
			}
		}
	}
}

// Append adds a new actor entry to the chain and returns the updated chain.
func (c ActorChain[T]) Append(e ActorEntry[T]) ActorChain[T] { return append(c, e) }

// ByKind returns all actors of a given kind in the chain.
func (c ActorChain[T]) ByKind(kind enums.ActorKind) []ActorEntry[T] {
	return slices.DeleteFunc(slices.Clone(c), func(e ActorEntry[T]) bool { return e.Kind != kind })
}

// HasKind checks if any actor in chain is of given kind.
func (c ActorChain[T]) HasKind(kind enums.ActorKind) bool {
	for _, e := range c {
		if e.Kind == kind {
			return true
		}
	}

	return false
}

// Constructor helpers

// MakeActor creates an actor entry for the given kind.
func MakeActor[T comparable](
	kind enums.ActorKind,
	id id.ID[struct{}, T],
	name ...string,
) ActorEntry[T] {
	n := ""
	if len(name) > 0 {
		n = name[0]
	}

	return ActorEntry[T]{Kind: kind, ID: id, Name: n}
}

// UserActor creates an actor entry for a human user.
func UserActor[T comparable](id id.ID[struct{}, T], name ...string) ActorEntry[T] {
	return MakeActor(enums.ActorKindUser, id, name...)
}

// BotActor creates an actor entry for an automated bot.
func BotActor[T comparable](id id.ID[struct{}, T], name ...string) ActorEntry[T] {
	return MakeActor(enums.ActorKindBot, id, name...)
}

// SystemActor creates an actor entry for system-initiated actions.
func SystemActor[T comparable]() ActorEntry[T] {
	var zeroID id.ID[struct{}, T]

	return ActorEntry[T]{Kind: enums.ActorKindSystem, ID: zeroID, Name: ""}
}

// ServiceActor creates an actor entry for a service-to-service call.
func ServiceActor[T comparable](id id.ID[struct{}, T], name ...string) ActorEntry[T] {
	return MakeActor(enums.ActorKindService, id, name...)
}

// IsZero returns true if this is the zero value.
func (e ActorEntry[T]) IsZero() bool {
	var zeroID id.ID[struct{}, T]

	return e.ID == zeroID && e.Name == ""
}
