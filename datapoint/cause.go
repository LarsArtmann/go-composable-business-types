package datapoint

import (
	"encoding/json"
	"fmt"

	"github.com/larsartmann/go-composable-business-types/nanoid"
)

// Cause represents a causal relationship to another DataPoint or event.
type Cause[T comparable] struct {
	id     nanoid.NanoId
	kind   string
	effect string
	trace  []nanoid.NanoId
}

// NewCause creates a new Cause with the given parameters.
func NewCause[T comparable](id nanoid.NanoId, kind, effect string, trace []nanoid.NanoId) Cause[T] {
	t := make([]nanoid.NanoId, len(trace))
	copy(t, trace)
	return Cause[T]{
		id:     id,
		kind:   kind,
		effect: effect,
		trace:  t,
	}
}

// NewCauseDirect creates a Cause for a direct causal relationship.
func NewCauseDirect[T comparable](id nanoid.NanoId) Cause[T] {
	return Cause[T]{
		id:     id,
		kind:   "direct",
		effect: "caused",
		trace:  nil,
	}
}

// NewCauseCommand creates a Cause for a command-triggered relationship.
func NewCauseCommand[T comparable](id nanoid.NanoId, command string) Cause[T] {
	return Cause[T]{
		id:     id,
		kind:   "command",
		effect: command,
		trace:  nil,
	}
}

// NewCauseEvent creates a Cause for an event-triggered relationship.
func NewCauseEvent[T comparable](id nanoid.NanoId, event string, trace ...nanoid.NanoId) Cause[T] {
	t := make([]nanoid.NanoId, len(trace))
	copy(t, trace)
	return Cause[T]{
		id:     id,
		kind:   "event",
		effect: event,
		trace:  t,
	}
}

// Id returns the cause ID.
func (c Cause[T]) Id() nanoid.NanoId { return c.id }

// Kind returns the cause kind (e.g., "direct", "command", "event").
func (c Cause[T]) Kind() string { return c.kind }

// Effect returns the effect description.
func (c Cause[T]) Effect() string { return c.effect }

// Trace returns the causal trace chain.
func (c Cause[T]) Trace() []nanoid.NanoId {
	if c.trace == nil {
		return nil
	}
	result := make([]nanoid.NanoId, len(c.trace))
	copy(result, c.trace)
	return result
}

// IsZero returns true if this is the zero value.
func (c Cause[T]) IsZero() bool {
	return c.id.IsZero() && c.kind == "" && c.effect == "" && len(c.trace) == 0
}

// jsonCause is the JSON representation of Cause.
type jsonCause struct {
	ID     string   `json:"id"`
	Kind   string   `json:"kind"`
	Effect string   `json:"effect"`
	Trace  []string `json:"trace,omitempty"`
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
	id, err := nanoid.ParseNanoId(raw.ID)
	if err != nil {
		return fmt.Errorf("unmarshal cause: parse id %q: %w", raw.ID, err)
	}
	c.id = id

	// Parse trace
	c.trace = make([]nanoid.NanoId, len(raw.Trace))
	for i, t := range raw.Trace {
		parsed, err := nanoid.ParseNanoId(t)
		if err != nil {
			return fmt.Errorf("unmarshal cause: parse trace[%d] %q from JSON %q: %w", i, t, string(data), err)
		}
		c.trace[i] = parsed
	}
	return nil
}
