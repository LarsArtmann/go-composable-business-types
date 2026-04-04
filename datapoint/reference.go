package datapoint

import (
	"encoding/json"
	"fmt"
	"maps"
)

// Taggable is an interface for types that support tagging.
type Taggable interface {
	Tags() map[string]string
	Tag(key string) string
}

// Reference is a type-safe reference to another entity.
type Reference[T comparable] struct {
	id       T
	relation string
	version  int
	tags     map[string]string
}

// NewReference creates a new Reference.
func NewReference[T comparable](id T, relation string) Reference[T] {
	return Reference[T]{
		id:       id,
		relation: relation,
		version:  0,
		tags:     make(map[string]string),
	}
}

// ID returns the referenced entity ID.
func (r Reference[T]) ID() T { return r.id }

// Relation returns the relationship type.
func (r Reference[T]) Relation() string { return r.relation }

// Version returns the expected version of the referenced entity.
func (r Reference[T]) Version() int { return r.version }

// getTags returns all tags from a map, returning nil if the map is nil.
func getTags(tags map[string]string) map[string]string {
	if tags == nil {
		return nil
	}

	return maps.Clone(tags)
}

// getTag returns a specific tag value, returning "" if tags is nil.
func getTag(tags map[string]string, key string) string {
	if tags == nil {
		return ""
	}

	return tags[key]
}

// Tags returns all reference tags.
func (r Reference[T]) Tags() map[string]string {
	return getTags(r.tags)
}

// Tag returns a specific tag value.
func (r Reference[T]) Tag(key string) string {
	return getTag(r.tags, key)
}

// IsZero returns true if this is the zero value.
func (r Reference[T]) IsZero() bool {
	var zero T

	return r.id == zero && r.relation == "" && r.version == 0 && len(r.tags) == 0
}

// WithVersion returns a copy with version set.
func (r Reference[T]) WithVersion(v int) Reference[T] {
	r.version = v

	return r
}

// WithTag returns a copy with a single tag added.
func (r Reference[T]) WithTag(key, value string) Reference[T] {
	addTag(&r.tags, key, value)

	return r
}

// addTag adds a tag to the map, initializing it if nil.
func addTag(tags *map[string]string, key, value string) {
	if *tags == nil {
		*tags = make(map[string]string)
	}

	(*tags)[key] = value
}

// jsonReference is the JSON representation of Reference.
type jsonReference[T comparable] struct {
	ID       T                 `json:"id"`
	Relation string            `json:"relation"`
	Version  int               `json:"version,omitzero"`
	Tags     map[string]string `json:"tags,omitzero"`
}

// MarshalJSON implements json.Marshaler.
func (r Reference[T]) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(jsonReference[T]{
		ID:       r.id,
		Relation: r.relation,
		Version:  r.version,
		Tags:     r.tags,
	})
	if err != nil {
		return nil, fmt.Errorf("reference: marshal JSON: %w", err)
	}

	return b, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (r *Reference[T]) UnmarshalJSON(data []byte) error {
	var raw jsonReference[T]

	err := json.Unmarshal(data, &raw)
	if err != nil {
		return fmt.Errorf("unmarshal reference: invalid JSON %q: %w", string(data), err)
	}

	r.id = raw.ID
	r.relation = raw.Relation
	r.version = raw.Version
	r.tags = raw.Tags

	return nil
}
