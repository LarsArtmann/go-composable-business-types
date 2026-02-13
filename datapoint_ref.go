package cbt

import (
	"encoding/json"
	"maps"
)

// Reference[T] is a type-safe reference to another entity by ID.
// It captures the relationship kind for query and traversal purposes.
type Reference[T comparable] struct {
	id       T                 // The entity ID being referenced
	relation string            // Relationship type (e.g., "parent", "owner", "source")
	version  int               // Version of the referenced entity (0 = any/latest)
	tags     map[string]string // Additional metadata about the reference
}

// NewReference creates a new reference to an entity by ID.
func NewReference[T comparable](id T, relation string) Reference[T] {
	return Reference[T]{
		id:       id,
		relation: relation,
		version:  0,
		tags:     nil,
	}
}

// NewReferenceWithVersion creates a reference to a specific version of an entity.
func NewReferenceWithVersion[T comparable](id T, relation string, version int) Reference[T] {
	return Reference[T]{
		id:       id,
		relation: relation,
		version:  version,
		tags:     nil,
	}
}

// Id returns the referenced entity's ID.
func (r Reference[T]) Id() T { return r.id }

// Relation returns the relationship type.
func (r Reference[T]) Relation() string { return r.relation }

// Version returns the referenced entity's version (0 = any/latest).
func (r Reference[T]) Version() int { return r.version }

// Tags returns a copy of the reference metadata tags.
func (r Reference[T]) Tags() map[string]string {
	if r.tags == nil {
		return nil
	}
	result := make(map[string]string, len(r.tags))
	maps.Copy(result, r.tags)
	return result
}

// Tag returns the value for a specific tag key.
func (r Reference[T]) Tag(key string) (string, bool) {
	if r.tags == nil {
		return "", false
	}
	v, ok := r.tags[key]
	return v, ok
}

// WithVersion returns a copy with version set.
func (r Reference[T]) WithVersion(version int) Reference[T] {
	r.version = version
	return r
}

// WithTag returns a copy with an additional tag.
func (r Reference[T]) WithTag(key, value string) Reference[T] {
	if r.tags == nil {
		r.tags = make(map[string]string)
	} else {
		newTags := make(map[string]string, len(r.tags)+1)
		maps.Copy(newTags, r.tags)
		r.tags = newTags
	}
	r.tags[key] = value
	return r
}

// jsonReference is the JSON representation of Reference.
type jsonReference[T comparable] struct {
	Id       T                 `json:"id"`
	Relation string            `json:"relation"`
	Version  int               `json:"version,omitempty"`
	Tags     map[string]string `json:"tags,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (r Reference[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(jsonReference[T]{
		Id:       r.id,
		Relation: r.relation,
		Version:  r.version,
		Tags:     r.tags,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (r *Reference[T]) UnmarshalJSON(data []byte) error {
	var raw jsonReference[T]
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	r.id = raw.Id
	r.relation = raw.Relation
	r.version = raw.Version
	r.tags = raw.Tags
	return nil
}
