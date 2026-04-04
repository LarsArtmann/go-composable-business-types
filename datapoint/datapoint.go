// Package datapoint provides self-contained data units with complete audit trails.
//
// DataPoint[T] wraps a payload of type T with comprehensive metadata including
// NanoID, Bitemporal tracking, Actor, Trigger, Context, References, Causes, Tags, and Version.
// Inspired by event sourcing, it preserves all relationships at the application layer.
//
// Basic usage:
//
//	dp := datapoint.NewDataPoint(payload, actorEntry).
//	    WithTrigger(enums.TriggerWebhook).
//	    WithReason("customer checkout").
//	    WithReference(datapoint.NewReference("order-123", "parent"))
package datapoint

import (
	"encoding/json"
	"fmt"
	"iter"

	"github.com/larsartmann/go-composable-business-types/actor"
	"github.com/larsartmann/go-composable-business-types/enums"
	"github.com/larsartmann/go-composable-business-types/nanoid"
	"github.com/larsartmann/go-composable-business-types/temporal"
	"github.com/larsartmann/go-composable-business-types/types"
)

// DataPoint is a self-contained data unit with a complete audit trail.
type DataPoint[T comparable] struct {
	id         nanoid.NanoID
	payload    T
	actor      actor.ActorEntry[string]
	temporal   temporal.Bitemporal
	trigger    enums.Trigger
	reason     string
	context    Context
	version    int
	tags       map[string]string
	references []Reference[string]
	causes     []Cause[string]
}

// NewDataPoint creates a new DataPoint with the given payload and actor.
func NewDataPoint[T comparable](payload T, actorEntry actor.ActorEntry[string]) DataPoint[T] {
	return DataPoint[T]{
		id:         nanoid.NewNanoID(),
		payload:    payload,
		actor:      actorEntry,
		temporal:   temporal.NewBitemporal(types.Now()),
		trigger:    enums.TriggerManual,
		reason:     "",
		context:    NewContext(),
		version:    1,
		tags:       make(map[string]string),
		references: nil,
		causes:     nil,
	}
}

// ID returns the DataPoint ID.
func (d DataPoint[T]) ID() nanoid.NanoID { return d.id }

// Payload returns the data payload.
func (d DataPoint[T]) Payload() T { return d.payload }

// Actor returns the actor entry.
func (d DataPoint[T]) Actor() actor.ActorEntry[string] { return d.actor }

// Temporal returns the bitemporal timestamps.
func (d DataPoint[T]) Temporal() temporal.Bitemporal { return d.temporal }

// Trigger returns the creation trigger.
func (d DataPoint[T]) Trigger() enums.Trigger { return d.trigger }

// Reason returns the reason for creation.
func (d DataPoint[T]) Reason() string { return d.reason }

// Context returns the execution context.
func (d DataPoint[T]) Context() Context { return d.context }

// Version returns the version number.
func (d DataPoint[T]) Version() int { return d.version }

func (d DataPoint[T]) Tags() map[string]string { return getTags(d.tags) }

// Tag returns a specific tag value.
func (d DataPoint[T]) Tag(key string) string { return getTag(d.tags, key) }

// cloneSlice creates a copy of the given slice, returning nil if the input is nil.
func cloneSlice[T any](s []T) []T {
	if s == nil {
		return nil
	}

	result := make([]T, len(s))
	copy(result, s)

	return result
}

// References returns all references.
func (d DataPoint[T]) References() []Reference[string] {
	return cloneSlice(d.references)
}

// Causes returns all causes.
func (d DataPoint[T]) Causes() []Cause[string] {
	return cloneSlice(d.causes)
}

// IsZero returns true if this is the zero value.
func (d DataPoint[T]) IsZero() bool {
	var (
		zeroPayload T
		zeroActor   actor.ActorEntry[string]
	)

	return d.id.IsZero() && d.payload == zeroPayload && d.actor == zeroActor && d.temporal.IsZero()
}

// WithTrigger returns a copy with trigger set.
func (d DataPoint[T]) WithTrigger(t enums.Trigger) DataPoint[T] {
	d.trigger = t

	return d
}

// WithReason returns a copy with reason set.
func (d DataPoint[T]) WithReason(reason string) DataPoint[T] {
	d.reason = reason

	return d
}

// WithContext returns a copy with context set.
func (d DataPoint[T]) WithContext(ctx Context) DataPoint[T] {
	d.context = ctx

	return d
}

// WithVersion returns a copy with version set.
func (d DataPoint[T]) WithVersion(v int) DataPoint[T] {
	d.version = withVersion(v)

	return d
}

// WithTag returns a copy with a single tag added.
func (d DataPoint[T]) WithTag(key, value string) DataPoint[T] {
	d.tags = withTag(d.tags, key, value)

	return d
}

// WithTags returns a copy with multiple tags merged.
func (d DataPoint[T]) WithTags(tags map[string]string) DataPoint[T] {
	d.tags = mergeTags(d.tags, tags)

	return d
}

// WithReference returns a copy with a reference appended.
func (d DataPoint[T]) WithReference(ref Reference[string]) DataPoint[T] {
	d.references = appendToSlice(d.references, ref)

	return d
}

// WithCause returns a copy with a cause appended.
func (d DataPoint[T]) WithCause(cause Cause[string]) DataPoint[T] {
	d.causes = appendToSlice(d.causes, cause)

	return d
}

// appendToSlice appends an item to a slice and returns the new slice.
func appendToSlice[T any](slice []T, item T) []T {
	return append(slice, item)
}

// WithTemporal returns a copy with temporal set.
func (d DataPoint[T]) WithTemporal(t temporal.Bitemporal) DataPoint[T] {
	d.temporal = t

	return d
}

// sliceToSeq converts a slice to an iterator sequence.
func sliceToSeq[T any](slice []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, item := range slice {
			if !yield(item) {
				return
			}
		}
	}
}

// AllReferences returns an iterator over all references.
func (d DataPoint[T]) AllReferences() iter.Seq[Reference[string]] {
	return sliceToSeq(d.references)
}

// AllCauses returns an iterator over all causes.
func (d DataPoint[T]) AllCauses() iter.Seq[Cause[string]] {
	return sliceToSeq(d.causes)
}

// AllTags returns an iterator over all tag key-value pairs.
func (d DataPoint[T]) AllTags() iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		for k, v := range d.tags {
			if !yield(k, v) {
				return
			}
		}
	}
}

// jsonDataPoint is the JSON representation of DataPoint.
type jsonDataPoint[T comparable] struct {
	ID         nanoid.NanoID            `json:"id"`
	Payload    T                        `json:"payload"`
	Actor      actor.ActorEntry[string] `json:"actor"`
	Temporal   temporal.Bitemporal      `json:"temporal"`
	Trigger    enums.Trigger            `json:"trigger"`
	Reason     string                   `json:"reason,omitzero"`
	Context    Context                  `json:"context"`
	Version    int                      `json:"version"`
	Tags       map[string]string        `json:"tags,omitzero"`
	References []Reference[string]      `json:"references,omitzero"`
	Causes     []Cause[string]          `json:"causes,omitzero"`
}

// MarshalJSON implements json.Marshaler.
func (d DataPoint[T]) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(jsonDataPoint[T]{
		ID:         d.id,
		Payload:    d.payload,
		Actor:      d.actor,
		Temporal:   d.temporal,
		Trigger:    d.trigger,
		Reason:     d.reason,
		Context:    d.context,
		Version:    d.version,
		Tags:       d.tags,
		References: d.references,
		Causes:     d.causes,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal datapoint: %w", err)
	}

	return b, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *DataPoint[T]) UnmarshalJSON(data []byte) error {
	var raw jsonDataPoint[T]

	err := json.Unmarshal(data, &raw)
	if err != nil {
		return fmt.Errorf("unmarshal datapoint: invalid JSON %q: %w", string(data), err)
	}

	// Parse ID
	d.id = raw.ID

	d.payload = raw.Payload
	d.actor = raw.Actor
	d.temporal = raw.Temporal
	d.trigger = raw.Trigger
	d.reason = raw.Reason
	d.context = raw.Context
	d.version = raw.Version
	d.tags = raw.Tags
	d.references = raw.References
	d.causes = raw.Causes

	return nil
}
