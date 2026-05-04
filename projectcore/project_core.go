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

type ProjectCore struct {
	Name       string                        `json:"name"`
	Path       string                        `json:"path"`
	Languages  programminglanguage.Languages `json:"languages"`
	Importance importance.Importance         `json:"importance"`
	Tags       []tag.Tag                     `json:"tags"`
}

type Option func(*ProjectCore)

func WithImportance(i importance.Importance) Option {
	return func(p *ProjectCore) { p.Importance = i }
}

func WithTags(tags ...tag.Tag) Option {
	return func(p *ProjectCore) { p.Tags = tags }
}

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

func (p *ProjectCore) IsZero() bool {
	return p == nil || (p.Name == "" && p.Path == "" && p.Languages.IsZero())
}

func (p *ProjectCore) Validate() error {
	if p == nil {
		return errors.New("projectcore: nil")
	}

	if p.Name == "" {
		return errors.New("projectcore: name is required")
	}

	if p.Path == "" {
		return errors.New("projectcore: path is required")
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

func (p *ProjectCore) MarshalJSON() ([]byte, error) {
	type Alias ProjectCore

	return json.Marshal((*Alias)(p))
}

func (p *ProjectCore) UnmarshalJSON(data []byte) error {
	type Alias ProjectCore

	var a Alias

	err := json.Unmarshal(data, &a)
	if err != nil {
		return fmt.Errorf("projectcore: unmarshal: %w", err)
	}

	*p = ProjectCore(a)

	return nil
}

var _ validate.Validator = (*ProjectCore)(nil)
