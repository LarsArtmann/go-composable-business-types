package datapoint

import (
	"encoding/json"
	"fmt"
	"maps"
)

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

// Tags returns all reference tags.
func (r Reference[T]) Tags() map[string]string {
	if r.tags == nil {
		return nil
	}
	return maps.Clone(r.tags)
}

// Tag returns a specific tag value.
func (r Reference[T]) Tag(key string) string {
	if r.tags == nil {
		return ""
	}
	return r.tags[key]
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
	if r.tags == nil {
		r.tags = make(map[string]string)
	}
	r.tags[key] = value
	return r
}

// jsonReference is the JSON representation of Reference.
type jsonReference struct {
	ID       string            `json:"id"`
	Relation string            `json:"relation"`
	Version  int               `json:"version,omitempty"`
	Tags     map[string]string `json:"tags,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (r Reference[T]) MarshalJSON() ([]byte, error) {
	var idStr string
	if s, ok := any(r.id).(interface{ String() string }); ok {
		idStr = s.String()
	} else if s, ok := any(r.id).(string); ok {
		idStr = s
	} else {
		return nil, fmt.Errorf(
			"reference: ID type %T with value %q does not support String() or string conversion",
			r.id,
			idStr,
		)
	}
	b, err := json.Marshal(jsonReference{
		ID:       idStr,
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
	var raw jsonReference
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("unmarshal reference: invalid JSON %q: %w", string(data), err)
	}
	r.relation = raw.Relation
	r.version = raw.Version
	r.tags = raw.Tags

	// Try to unmarshal ID based on type
	var zero T
	switch any(zero).(type) {
	case string:
		id, ok := any(raw.ID).(T)
		if !ok {
			return fmt.Errorf(
				"unmarshal reference: cannot convert ID %q to target type %T",
				raw.ID,
				zero,
			)
		}
		r.id = id
	default:
		// For other types, the ID type must support text unmarshaling
		if u, ok := any(&r.id).(interface{ UnmarshalText([]byte) error }); ok {
			if err := u.UnmarshalText([]byte(raw.ID)); err != nil {
				return fmt.Errorf(
					"unmarshal reference: unmarshal id %q to type %T: %w",
					raw.ID,
					zero,
					err,
				)
			}
		}
	}
	return nil
}
