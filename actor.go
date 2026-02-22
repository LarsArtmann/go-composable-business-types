package cbt

import "slices"

// ActorEntry represents a single actor in an actor chain.
type ActorEntry[T comparable] struct {
	Kind ActorKind
	Id   ID[struct{}, T]
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
func (c ActorChain[T]) ByKind(kind ActorKind) []ActorEntry[T] {
	return slices.DeleteFunc(slices.Clone(c), func(e ActorEntry[T]) bool { return e.Kind != kind })
}

// HasKind checks if any actor in chain is of given kind.
func (c ActorChain[T]) HasKind(kind ActorKind) bool {
	for _, e := range c {
		if e.Kind == kind {
			return true
		}
	}
	return false
}

// Constructor helpers

// UserActor creates an actor entry for a human user.
func UserActor[T comparable](id ID[struct{}, T], name ...string) ActorEntry[T] {
	return newActorEntry(ActorKindUser, id, name...)
}

// BotActor creates an actor entry for an automated bot.
func BotActor[T comparable](id ID[struct{}, T], name ...string) ActorEntry[T] {
	return newActorEntry(ActorKindBot, id, name...)
}

// SystemActor creates an actor entry for system-initiated actions.
func SystemActor[T comparable]() ActorEntry[T] {
	return ActorEntry[T]{Kind: ActorKindSystem}
}

// ServiceActor creates an actor entry for a service-to-service call.
func ServiceActor[T comparable](id ID[struct{}, T], name ...string) ActorEntry[T] {
	return newActorEntry(ActorKindService, id, name...)
}

// newActorEntry is a helper to create ActorEntry with optional name.
func newActorEntry[T comparable](kind ActorKind, id ID[struct{}, T], name ...string) ActorEntry[T] {
	n := ""
	if len(name) > 0 {
		n = name[0]
	}
	return ActorEntry[T]{Kind: kind, Id: id, Name: n}
}
