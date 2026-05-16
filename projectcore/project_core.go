// Package projectcore provides a composite type for project metadata.
//
// ProjectCore aggregates name, path, programming languages, importance,
// and tags into a single validated struct with JSON support.
//
// Basic usage:
//
//	core := projectcore.New("my-app", "/src/my-app",
//	    programminglanguage.NewLanguages(programminglanguage.New("go")),
//	    projectcore.WithImportance(importance.High))
package projectcore

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/larsartmann/go-composable-business-types/importance"
	"github.com/larsartmann/go-composable-business-types/programminglanguage"
	"github.com/larsartmann/go-composable-business-types/tag"
	"github.com/larsartmann/go-composable-business-types/validate"
)

var (
	errProjectNil     = errors.New("projectcore: nil")
	errProjectNameReq = errors.New("projectcore: name is required")
	errProjectPathReq = errors.New("projectcore: path is required")
)

// ProjectCore aggregates project metadata: name, path, languages, importance, and tags.
type ProjectCore struct {
	Name       string                        `json:"name"`
	Path       string                        `json:"path"`
	Languages  programminglanguage.Languages `json:"languages"`
	Importance importance.Importance         `json:"importance"`
	Tags       []tag.Tag                     `json:"tags"`
}

// Option configures a ProjectCore during construction.
type Option func(*ProjectCore)

// WithImportance sets the importance on a ProjectCore.
func WithImportance(i importance.Importance) Option {
	return func(p *ProjectCore) { p.Importance = i }
}

// WithTags sets the tags on a ProjectCore.
func WithTags(tags ...tag.Tag) Option {
	return func(p *ProjectCore) { p.Tags = tags }
}

// New creates a ProjectCore with the given name, path, languages, and optional configuration.
func New(name, path string, langs programminglanguage.Languages, opts ...Option) *ProjectCore {
	p := &ProjectCore{
		Name:      name,
		Path:      path,
		Languages: langs,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// IsZero reports whether the ProjectCore is nil or has no name, path, or languages.
func (p *ProjectCore) IsZero() bool {
	return p == nil || (p.Name == "" && p.Path == "" && p.Languages.IsZero())
}

// Validate checks that name and path are present and all nested fields are valid.
func (p *ProjectCore) Validate() error {
	if p == nil {
		return errProjectNil
	}

	if p.Name == "" {
		return errProjectNameReq
	}

	if p.Path == "" {
		return errProjectPathReq
	}

	err := p.Importance.Validate()
	if err != nil {
		return fmt.Errorf("projectcore: importance: %w", err)
	}

	for i, t := range p.Tags {
		err := t.Validate()
		if err != nil {
			return fmt.Errorf("projectcore: tag[%d]: %w", i, err)
		}
	}

	return nil
}

// MarshalJSON encodes the ProjectCore as JSON.
func (p *ProjectCore) MarshalJSON() ([]byte, error) {
	type Alias ProjectCore

	return json.Marshal((*Alias)(p)) //nolint:wrapcheck // internal type alias marshaling
}

// UnmarshalJSON decodes JSON into the ProjectCore.
func (p *ProjectCore) UnmarshalJSON(data []byte) error {
	type Alias ProjectCore

	var alias Alias

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return fmt.Errorf("projectcore: unmarshal: %w", err)
	}

	*p = ProjectCore(alias)

	return nil
}

var _ validate.Validator = (*ProjectCore)(nil)
