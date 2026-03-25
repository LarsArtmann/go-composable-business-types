package datapoint

import (
	"encoding/json"
	"fmt"
	"maps"
)

// Context represents the execution context for a DataPoint.
type Context struct {
	environment string
	session     string
	request     string
	source      string
	tags        map[string]string
}

// NewContext creates a new Context with default values.
func NewContext() Context {
	return Context{
		tags: make(map[string]string),
	}
}

// Environment returns the environment name (e.g., "production", "staging").
func (c Context) Environment() string { return c.environment }

// Session returns the session identifier.
func (c Context) Session() string { return c.session }

// Request returns the request identifier.
func (c Context) Request() string { return c.request }

// Source returns the source system identifier.
func (c Context) Source() string { return c.source }

// Tags returns all context tags.
func (c Context) Tags() map[string]string {
	if c.tags == nil {
		return nil
	}
	return maps.Clone(c.tags)
}

// Tag returns a specific tag value.
func (c Context) Tag(key string) string {
	if c.tags == nil {
		return ""
	}
	return c.tags[key]
}

// IsZero returns true if this is the zero value.
func (c Context) IsZero() bool {
	return c.environment == "" && c.session == "" && c.request == "" && c.source == "" &&
		len(c.tags) == 0
}

// WithEnvironment returns a copy with environment set.
func (c Context) WithEnvironment(env string) Context {
	c.environment = env
	return c
}

// WithSession returns a copy with session set.
func (c Context) WithSession(session string) Context {
	c.session = session
	return c
}

// WithRequest returns a copy with request set.
func (c Context) WithRequest(request string) Context {
	c.request = request
	return c
}

// WithSource returns a copy with source set.
func (c Context) WithSource(source string) Context {
	c.source = source
	return c
}

// WithTag returns a copy with a single tag added.
func (c Context) WithTag(key, value string) Context {
	if c.tags == nil {
		c.tags = make(map[string]string)
	}
	c.tags[key] = value
	return c
}

// WithTags returns a copy with multiple tags merged.
func (c Context) WithTags(tags map[string]string) Context {
	if len(tags) == 0 {
		return c
	}
	if c.tags == nil {
		c.tags = make(map[string]string)
	}
	maps.Copy(c.tags, tags)
	return c
}

// jsonContext is the JSON representation of Context.
type jsonContext struct {
	Environment string            `json:"environment,omitempty"`
	Session     string            `json:"session,omitempty"`
	Request     string            `json:"request,omitempty"`
	Source      string            `json:"source,omitempty"`
	Tags        map[string]string `json:"tags,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (c Context) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(jsonContext{
		Environment: c.environment,
		Session:     c.session,
		Request:     c.request,
		Source:      c.source,
		Tags:        c.tags,
	})
	if err != nil {
		return nil, fmt.Errorf("context: marshal JSON: %w", err)
	}
	return b, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (c *Context) UnmarshalJSON(data []byte) error {
	var raw jsonContext
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("unmarshal context: invalid JSON %q: %w", string(data), err)
	}
	c.environment = raw.Environment
	c.session = raw.Session
	c.request = raw.Request
	c.source = raw.Source
	c.tags = raw.Tags
	return nil
}
