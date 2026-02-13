package cbt

import (
	"encoding/json"
)

// Cause[T] represents a causal relationship between data points.
// It enables building a causal chain showing why a DataPoint was created.
type Cause[T comparable] struct {
	id     NanoId   // ID of the causing DataPoint
	kind   string   // Category of cause (e.g., "command", "event", "timeout")
	effect string   // What effect this cause had (e.g., "created", "updated", "deleted")
	trace  []NanoId // Optional trace of intermediate causes
}

// NewCause creates a new Cause referencing a DataPoint by ID.
func NewCause(id NanoId, kind, effect string) Cause[string] {
	return Cause[string]{
		id:     id,
		kind:   kind,
		effect: effect,
		trace:  nil,
	}
}

// NewCauseDirect creates a Cause for a directly caused DataPoint.
func NewCauseDirect(id NanoId) Cause[string] {
	return NewCause(id, "direct", "caused")
}

// NewCauseCommand creates a Cause for a command-triggered DataPoint.
func NewCauseCommand(id NanoId, effect string) Cause[string] {
	return NewCause(id, "command", effect)
}

// NewCauseEvent creates a Cause for an event-triggered DataPoint.
func NewCauseEvent(id NanoId, effect string) Cause[string] {
	return NewCause(id, "event", effect)
}

// Id returns the ID of the causing DataPoint.
func (c Cause[T]) Id() NanoId { return c.id }

// Kind returns the category of cause.
func (c Cause[T]) Kind() string { return c.kind }

// Effect returns what effect this cause had.
func (c Cause[T]) Effect() string { return c.effect }

// Trace returns the causal trace (intermediate causes).
func (c Cause[T]) Trace() []NanoId { return c.trace }

// HasTrace returns true if there are intermediate causes.
func (c Cause[T]) HasTrace() bool { return len(c.trace) > 0 }

// WithTrace returns a copy with the causal trace set.
func (c Cause[T]) WithTrace(trace []NanoId) Cause[T] {
	c.trace = make([]NanoId, len(trace))
	copy(c.trace, trace)
	return c
}

// AppendTrace returns a copy with an ID appended to the trace.
func (c Cause[T]) AppendTrace(id NanoId) Cause[T] {
	c.trace = append(c.trace, id)
	return c
}

// jsonCause is the JSON representation of Cause.
type jsonCause struct {
	Id     string   `json:"id"`
	Kind   string   `json:"kind"`
	Effect string   `json:"effect"`
	Trace  []string `json:"trace,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (c Cause[T]) MarshalJSON() ([]byte, error) {
	trace := make([]string, len(c.trace))
	for i, id := range c.trace {
		trace[i] = id.String()
	}
	return json.Marshal(jsonCause{
		Id:     c.id.String(),
		Kind:   c.kind,
		Effect: c.effect,
		Trace:  trace,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (c *Cause[T]) UnmarshalJSON(data []byte) error {
	var raw jsonCause
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	id, err := ParseNanoId(raw.Id)
	if err != nil {
		return err
	}

	c.id = id
	c.kind = raw.Kind
	c.effect = raw.Effect

	if len(raw.Trace) > 0 {
		c.trace = make([]NanoId, len(raw.Trace))
		for i, s := range raw.Trace {
			traceId, err := ParseNanoId(s)
			if err != nil {
				return err
			}
			c.trace[i] = traceId
		}
	}

	return nil
}
