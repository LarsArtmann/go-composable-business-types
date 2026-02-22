package cbt

import (
	"encoding/json"
	"maps"
)

// Context captures the execution context in which a DataPoint was created.
// This includes environment, session, request, and other metadata.
type Context struct {
	environment string            // e.g., "production", "staging", "development"
	session     string            // Session or correlation ID (optional)
	request     string            // Request ID for tracing (optional)
	source      string            // Source system/service name
	tags        map[string]string // Additional key-value metadata
}

// NewContext creates a new Context with the given source system.
func NewContext(source string) Context {
	return Context{
		source: source,
		tags:   make(map[string]string),
	}
}

// Environment returns the environment name (e.g., "production").
func (c Context) Environment() string { return c.environment }

// Session returns the session or correlation ID.
func (c Context) Session() string { return c.session }

// Request returns the request ID for tracing.
func (c Context) Request() string { return c.request }

// Source returns the source system/service name.
func (c Context) Source() string { return c.source }

// IsZero returns true if this is the zero value (all fields empty).
func (c Context) IsZero() bool {
	return c.source == "" && c.environment == "" && c.session == "" && c.request == "" && len(c.tags) == 0
}

// Tags returns a copy of the metadata tags.
func (c Context) Tags() map[string]string {
	if c.tags == nil {
		return nil
	}
	result := make(map[string]string, len(c.tags))
	maps.Copy(result, c.tags)
	return result
}

// Tag returns the value for a specific tag key.
func (c Context) Tag(key string) (string, bool) {
	if c.tags == nil {
		return "", false
	}
	v, ok := c.tags[key]
	return v, ok
}

// WithEnvironment returns a copy with environment set.
func (c Context) WithEnvironment(env string) Context {
	c.environment = env
	return c
}

// WithSession returns a copy with session/correlation ID set.
func (c Context) WithSession(session string) Context {
	c.session = session
	return c
}

// WithRequest returns a copy with request ID set.
func (c Context) WithRequest(request string) Context {
	c.request = request
	return c
}

// WithSource returns a copy with source system set.
func (c Context) WithSource(source string) Context {
	c.source = source
	return c
}

// WithTag returns a copy with an additional tag.
func (c Context) WithTag(key, value string) Context {
	if c.tags == nil {
		c.tags = make(map[string]string)
	} else {
		// Copy the map to avoid mutating the original
		newTags := make(map[string]string, len(c.tags)+1)
		maps.Copy(newTags, c.tags)
		c.tags = newTags
	}
	c.tags[key] = value
	return c
}

// WithTags returns a copy with multiple tags added.
func (c Context) WithTags(tags map[string]string) Context {
	result := c
	for k, v := range tags {
		result = result.WithTag(k, v)
	}
	return result
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
	return json.Marshal(jsonContext{
		Environment: c.environment,
		Session:     c.session,
		Request:     c.request,
		Source:      c.source,
		Tags:        c.tags,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (c *Context) UnmarshalJSON(data []byte) error {
	var raw jsonContext
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	c.environment = raw.Environment
	c.session = raw.Session
	c.request = raw.Request
	c.source = raw.Source
	c.tags = raw.Tags
	return nil
}
