package datapoint

import (
	"slices"

	"github.com/larsartmann/go-composable-business-types/enums"
	"github.com/larsartmann/go-composable-business-types/nanoid"
	"github.com/larsartmann/go-composable-business-types/types"
)

// Cause represents a causal relationship to another DataPoint or event.
type Cause[T comparable] struct {
	id     nanoid.NanoID
	kind   enums.CauseKind
	effect string
	trace  []nanoid.NanoID
}

// NewCause creates a new Cause with the given parameters.
func NewCause[T comparable](
	id nanoid.NanoID,
	kind enums.CauseKind,
	effect string,
	trace []nanoid.NanoID,
) Cause[T] {
	return Cause[T]{
		id:     id,
		kind:   kind,
		effect: effect,
		trace:  slices.Clone(trace),
	}
}

// NewCauseDirect creates a Cause for a direct causal relationship.
func NewCauseDirect[T comparable](id nanoid.NanoID) Cause[T] {
	return Cause[T]{
		id:     id,
		kind:   enums.CauseKindDirect,
		effect: "caused",
		trace:  nil,
	}
}

// NewCauseCommand creates a Cause for a command-triggered relationship.
func NewCauseCommand[T comparable](id nanoid.NanoID, command string) Cause[T] {
	return Cause[T]{
		id:     id,
		kind:   enums.CauseKindCommand,
		effect: command,
		trace:  nil,
	}
}

// NewCauseEvent creates a Cause for an event-triggered relationship.
func NewCauseEvent[T comparable](id nanoid.NanoID, event string, trace ...nanoid.NanoID) Cause[T] {
	return Cause[T]{
		id:     id,
		kind:   enums.CauseKindEvent,
		effect: event,
		trace:  slices.Clone(trace),
	}
}

// ID returns the cause ID.
func (c Cause[T]) ID() nanoid.NanoID { return c.id }

// Kind returns the cause kind.
func (c Cause[T]) Kind() enums.CauseKind { return c.kind }

// Effect returns the effect description.
func (c Cause[T]) Effect() string { return c.effect }

// Trace returns the causal trace chain.
func (c Cause[T]) Trace() []nanoid.NanoID {
	if c.trace == nil {
		return nil
	}

	return slices.Clone(c.trace)
}

// IsZero returns true if this is the zero value.
func (c Cause[T]) IsZero() bool {
	return c.id.IsZero() && c.kind == 0 && c.effect == "" && len(c.trace) == 0
}

// jsonCause is the JSON representation of Cause.
type jsonCause struct {
	ID     nanoid.NanoID   `json:"id"`
	Kind   enums.CauseKind `json:"kind"`
	Effect string          `json:"effect"`
	Trace  []nanoid.NanoID `json:"trace,omitzero"`
}

// MarshalJSON implements json.Marshaler.
func (c Cause[T]) MarshalJSON() ([]byte, error) {
	return types.MarshalJSON("cause", jsonCause{
		ID:     c.id,
		Kind:   c.kind,
		Effect: c.effect,
		Trace:  c.trace,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (c *Cause[T]) UnmarshalJSON(data []byte) error {
	var raw jsonCause

	err := types.UnmarshalJSON("cause", data, &raw)
	if err != nil {
		return err
	}

	c.kind = raw.Kind
	c.effect = raw.Effect

	// Parse ID
	c.id = raw.ID

	// Parse trace
	c.trace = slices.Clone(raw.Trace)

	return nil
}
