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
	"slices"

	"github.com/larsartmann/go-composable-business-types/enums"
	"github.com/larsartmann/go-composable-business-types/id"
)

// ActorEntry represents a single actor in an actor chain.
type ActorEntry[T comparable] struct {
	Kind enums.ActorKind
	Id   id.ID[struct{}, T]
	Name string // optional human-readable name
}

// ActorChain represents an ordered chain of actors from original to current.
// Index 0 = original actor, last index = current actor.
type ActorChain[T comparable] []ActorEntry[T]

func NewActorChain[T comparable](first ActorEntry[T]) ActorChain[T] {
	return ActorChain[T]{first}
}

func (c ActorChain[T]) Origin() ActorEntry[T]                { return c[0] }
func (c ActorChain[T]) Current() ActorEntry[T]               { return c[len(c)-1] }
func (c ActorChain[T]) IsZero() bool                         { return len(c) == 0 }
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

// UserActor creates an actor entry for a human user.
func UserActor[T comparable](id id.ID[struct{}, T], name ...string) ActorEntry[T] {
	return newActorEntry(enums.ActorKindUser, id, name...)
}

// BotActor creates an actor entry for an automated bot.
func BotActor[T comparable](id id.ID[struct{}, T], name ...string) ActorEntry[T] {
	return newActorEntry(enums.ActorKindBot, id, name...)
}

// SystemActor creates an actor entry for system-initiated actions.
func SystemActor[T comparable]() ActorEntry[T] {
	return ActorEntry[T]{Kind: enums.ActorKindSystem}
}

// ServiceActor creates an actor entry for a service-to-service call.
func ServiceActor[T comparable](id id.ID[struct{}, T], name ...string) ActorEntry[T] {
	return newActorEntry(enums.ActorKindService, id, name...)
}

// newActorEntry is a helper to create ActorEntry with optional name.
func newActorEntry[T comparable](kind enums.ActorKind, id id.ID[struct{}, T], name ...string) ActorEntry[T] {
	n := ""
	if len(name) > 0 {
		n = name[0]
	}
	return ActorEntry[T]{Kind: kind, Id: id, Name: n}
}

// IsZero returns true if this is the zero value.
func (e ActorEntry[T]) IsZero() bool {
	var zeroID id.ID[struct{}, T]
	return e.Id == zeroID && e.Name == ""
}
