package datapoint

import (
	"encoding/json"
	"fmt"

	"github.com/larsartmann/go-composable-business-types/enums"
	"github.com/larsartmann/go-composable-business-types/nanoid"
)

// Cause represents a causal relationship to another DataPoint or event.
type Cause[T comparable] struct {
	id     nanoid.NanoID
	kind   enums.CauseKind
	effect string
	trace  []nanoid.NanoID
}

// NewCause creates a new Cause with the given parameters.
func NewCause[T comparable](id nanoid.NanoID, kind enums.CauseKind, effect string, trace []nanoid.NanoID) Cause[T] {
	t := make([]nanoid.NanoID, len(trace))
	copy(t, trace)
	return Cause[T]{
		id:     id,
		kind:   kind,
		effect: effect,
		trace:  t,
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
	t := make([]nanoid.NanoID, len(trace))
	copy(t, trace)
	return Cause[T]{
		id:     id,
		kind:   enums.CauseKindEvent,
		effect: event,
		trace:  t,
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
	result := make([]nanoid.NanoID, len(c.trace))
	copy(result, c.trace)
	return result
}

// IsZero returns true if this is the zero value.
func (c Cause[T]) IsZero() bool {
	return c.id.IsZero() && c.kind == 0 && c.effect == "" && len(c.trace) == 0
}

// jsonCause is the JSON representation of Cause.
type jsonCause struct {
	ID     string           `json:"id"`
	Kind   enums.CauseKind `json:"kind"`
	Effect string          `json:"effect"`
	Trace  []string        `json:"trace,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (c Cause[T]) MarshalJSON() ([]byte, error) {
	trace := make([]string, len(c.trace))
	for i, t := range c.trace {
		trace[i] = t.String()
	}
	return json.Marshal(jsonCause{
		ID:     c.id.String(),
		Kind:   c.kind,
		Effect: c.effect,
		Trace:  trace,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (c *Cause[T]) UnmarshalJSON(data []byte) error {
	var raw jsonCause
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("unmarshal cause: invalid JSON %q: %w", string(data), err)
	}
	c.kind = raw.Kind
	c.effect = raw.Effect

	// Parse ID
	id, err := nanoid.ParseNanoID(raw.ID)
	if err != nil {
		return fmt.Errorf("unmarshal cause: parse id %q: %w", raw.ID, err)
	}
	c.id = id

	// Parse trace
	c.trace = make([]nanoid.NanoID, len(raw.Trace))
	for i, t := range raw.Trace {
		parsed, err := nanoid.ParseNanoID(t)
		if err != nil {
			return fmt.Errorf(
				"unmarshal cause: parse trace[%d] %q from JSON %q: %w",
				i,
				t,
				string(data),
				err,
			)
		}
		c.trace[i] = parsed
	}
	return nil
}
