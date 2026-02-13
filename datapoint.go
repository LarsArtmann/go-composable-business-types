package cbt

import (
	"encoding/json"
	"maps"
	"time"
)

// DataPoint is a self-contained unit of data with complete audit trail.
// Inspired by event sourcing, it preserves ALL relationships and metadata
// at the application layer, enabling full traceability without external systems.
//
// The generic T represents the payload type (e.g., OrderState, CustomerInfo).
// T must be comparable to support equality checks.
type DataPoint[T comparable] struct {
	id         NanoId
	payload    T
	actor      ActorEntry[string]  // Who caused this data point
	temporal   Bitemporal          // Bitemporal tracking (valid/recorded times)
	trigger    Trigger             // What caused this data point
	reason     string              // Human-readable explanation (optional)
	context    Context             // Execution context (environment, source, etc.)
	version    int                 // Version number for optimistic concurrency
	tags       map[string]string   // Additional metadata tags
	references []Reference[string] // References to other entities
	causes     []Cause[string]     // Causal chain (why this was created)
}

// NewDataPoint creates a new DataPoint with explicit timestamps.
// Use NewDataPointNow for automatic timestamping.
func NewDataPoint[T comparable](
	payload T,
	actor ActorEntry[string],
	occurred Timestamp,
	recorded Timestamp,
	reason ...string,
) DataPoint[T] {
	r := ""
	if len(reason) > 0 {
		r = reason[0]
	}
	return DataPoint[T]{
		id:       NewNanoId(),
		payload:  payload,
		actor:    actor,
		temporal: NewBitemporalWithRange(occurred, Timestamp{}, recorded),
		trigger:  TriggerManual,
		reason:   r,
	}
}

// NewDataPointNow creates a new DataPoint with recorded time set to now.
// Occurred time is set equal to recorded time (use NewDataPoint for different values).
func NewDataPointNow[T comparable](
	payload T,
	actor ActorEntry[string],
	reason ...string,
) DataPoint[T] {
	now := Now()
	return NewDataPoint(payload, actor, now, now, reason...)
}

// Id returns the unique identifier of this data point.
func (dp DataPoint[T]) Id() NanoId { return dp.id }

// Payload returns the data contained in this data point.
func (dp DataPoint[T]) Payload() T { return dp.payload }

// Actor returns who caused this data point.
func (dp DataPoint[T]) Actor() ActorEntry[string] { return dp.actor }

// Occurred returns when the event happened (business time).
// Deprecated: Use Temporal().ValidFrom() for full bitemporal support.
func (dp DataPoint[T]) Occurred() Timestamp { return dp.temporal.ValidFrom() }

// Recorded returns when the system recorded this data point.
// Deprecated: Use Temporal().Recorded() for full bitemporal support.
func (dp DataPoint[T]) Recorded() Timestamp { return dp.temporal.Recorded() }

// Temporal returns the bitemporal information.
func (dp DataPoint[T]) Temporal() Bitemporal { return dp.temporal }

// Trigger returns what caused this data point to be created.
func (dp DataPoint[T]) Trigger() Trigger { return dp.trigger }

// Reason returns the human-readable explanation (may be empty).
func (dp DataPoint[T]) Reason() string { return dp.reason }

// Context returns the execution context.
func (dp DataPoint[T]) Context() Context { return dp.context }

// Version returns the version number for optimistic concurrency.
func (dp DataPoint[T]) Version() int { return dp.version }

// Tags returns a copy of the metadata tags.
func (dp DataPoint[T]) Tags() map[string]string {
	if dp.tags == nil {
		return nil
	}
	result := make(map[string]string, len(dp.tags))
	maps.Copy(result, dp.tags)
	return result
}

// Tag returns the value for a specific tag key.
func (dp DataPoint[T]) Tag(key string) (string, bool) {
	v, ok := dp.tags[key]
	return v, ok
}

// References returns the references to other entities.
func (dp DataPoint[T]) References() []Reference[string] {
	if dp.references == nil {
		return nil
	}
	result := make([]Reference[string], len(dp.references))
	copy(result, dp.references)
	return result
}

// Causes returns the causal chain.
func (dp DataPoint[T]) Causes() []Cause[string] {
	if dp.causes == nil {
		return nil
	}
	result := make([]Cause[string], len(dp.causes))
	copy(result, dp.causes)
	return result
}

// WithReason returns a copy with the reason set.
func (dp DataPoint[T]) WithReason(reason string) DataPoint[T] {
	dp.reason = reason
	return dp
}

// WithTrigger returns a copy with the trigger set.
func (dp DataPoint[T]) WithTrigger(trigger Trigger) DataPoint[T] {
	dp.trigger = trigger
	return dp
}

// WithContext returns a copy with the context set.
func (dp DataPoint[T]) WithContext(ctx Context) DataPoint[T] {
	dp.context = ctx
	return dp
}

// WithTemporal returns a copy with the bitemporal information set.
func (dp DataPoint[T]) WithTemporal(temporal Bitemporal) DataPoint[T] {
	dp.temporal = temporal
	return dp
}

// WithVersion returns a copy with the version set.
func (dp DataPoint[T]) WithVersion(version int) DataPoint[T] {
	dp.version = version
	return dp
}

// WithTags returns a copy with the tags set.
func (dp DataPoint[T]) WithTags(tags map[string]string) DataPoint[T] {
	dp.tags = make(map[string]string, len(tags))
	maps.Copy(dp.tags, tags)
	return dp
}

// AddTag returns a copy with an additional tag.
func (dp DataPoint[T]) AddTag(key, value string) DataPoint[T] {
	if dp.tags == nil {
		dp.tags = make(map[string]string)
	}
	dp.tags[key] = value
	return dp
}

// WithReferences returns a copy with the references set.
func (dp DataPoint[T]) WithReferences(refs []Reference[string]) DataPoint[T] {
	dp.references = make([]Reference[string], len(refs))
	copy(dp.references, refs)
	return dp
}

// AddReference returns a copy with an additional reference.
func (dp DataPoint[T]) AddReference(ref Reference[string]) DataPoint[T] {
	dp.references = append(dp.references, ref)
	return dp
}

// WithCauses returns a copy with the causes set.
func (dp DataPoint[T]) WithCauses(causes []Cause[string]) DataPoint[T] {
	dp.causes = make([]Cause[string], len(causes))
	copy(dp.causes, causes)
	return dp
}

// AddCause returns a copy with an additional cause.
func (dp DataPoint[T]) AddCause(cause Cause[string]) DataPoint[T] {
	dp.causes = append(dp.causes, cause)
	return dp
}

// jsonDataPoint is the JSON representation of DataPoint.
type jsonDataPoint[T comparable] struct {
	Id         string              `json:"id"`
	Payload    T                   `json:"payload"`
	Actor      ActorEntry[string]  `json:"actor"`
	Temporal   Bitemporal          `json:"temporal"`
	Trigger    string              `json:"trigger"`
	Reason     string              `json:"reason,omitempty"`
	Context    Context             `json:"context"`
	Version    int                 `json:"version,omitempty"`
	Tags       map[string]string   `json:"tags,omitempty"`
	References []Reference[string] `json:"references,omitempty"`
	Causes     []Cause[string]     `json:"causes,omitempty"`
	// Legacy fields for backward compatibility
	Occurred time.Time `json:"occurred"`
	Recorded time.Time `json:"recorded"`
}

// MarshalJSON implements json.Marshaler.
func (dp DataPoint[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(jsonDataPoint[T]{
		Id:         dp.id.String(),
		Payload:    dp.payload,
		Actor:      dp.actor,
		Temporal:   dp.temporal,
		Trigger:    dp.trigger.String(),
		Reason:     dp.reason,
		Context:    dp.context,
		Version:    dp.version,
		Tags:       dp.tags,
		References: dp.references,
		Causes:     dp.causes,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (dp *DataPoint[T]) UnmarshalJSON(data []byte) error {
	var raw jsonDataPoint[T]
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	id, err := ParseNanoId(raw.Id)
	if err != nil {
		return err
	}

	// Parse trigger, default to Manual if empty or invalid
	trigger, err := ParseTrigger(raw.Trigger)
	if err != nil {
		trigger = TriggerManual
	}

	dp.id = id
	dp.payload = raw.Payload
	dp.actor = raw.Actor
	dp.temporal = raw.Temporal
	dp.trigger = trigger
	dp.reason = raw.Reason
	dp.context = raw.Context
	dp.version = raw.Version
	dp.tags = raw.Tags
	dp.references = raw.References
	dp.causes = raw.Causes

	// Handle legacy format (occurred/recorded instead of temporal)
	if raw.Occurred.IsZero() && !raw.Temporal.ValidFrom().IsZero() {
		// Already using new format
	} else if !raw.Occurred.IsZero() {
		// Legacy format - migrate to temporal
		dp.temporal = NewBitemporalWithRange(
			NewTimestamp(raw.Occurred),
			Timestamp{},
			NewTimestamp(raw.Recorded),
		)
	}

	return nil
}
