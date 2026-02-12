package cbt

// Id is a strongly-typed identifier that prevents mixing different entity IDs.
type Id[T comparable] struct{ value T }

func NewId[T comparable](v T) Id[T] { return Id[T]{value: v} }

func (id Id[T]) Value() T         { return id.value }
func (id Id[T]) IsZero() bool     { var zero T; return id.value == zero }
func (id Id[T]) GoString() string { return any(id.value).(string) }
