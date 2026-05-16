// Package programminglanguage provides normalized, branded identifiers for programming languages.
//
// It normalizes common aliases (e.g., "golang" → "go", "ts" → "typescript") and wraps
// them in a branded ID type for type safety.
//
// Basic usage:
//
//	lang := programminglanguage.New("golang") // branded ID with value "go"
//	langs := programminglanguage.NewLanguages(lang, programminglanguage.New("rust"))
package programminglanguage

import (
	"maps"
	"slices"
	"strings"

	id "github.com/larsartmann/go-branded-id"
	"github.com/larsartmann/go-composable-business-types/validate"
)

type languageBrand struct{}

// Language is a branded identifier for a normalized programming language name.
type Language = id.ID[languageBrand, string]

var aliases = map[string]string{
	"golang":         "go",
	"js":             "javascript",
	"mjs":            "javascript",
	"ts":             "typescript",
	"tsx":            "typescript",
	"py":             "python",
	"rs":             "rust",
	"jvm":            "java",
	"c#":             "csharp",
	"cs":             "csharp",
	"rb":             "ruby",
	"gdscript":       "godot",
	"terraform":      "hcl",
	"tf":             "hcl",
	"c++":            "cpp",
	"bash":           "shell",
	"sh":             "shell",
	"zsh":            "shell",
	"fish":           "shell",
	"ps1":            "powershell",
	"nodejs":         "node",
	"typescript-jsx": "typescript",
}

// New creates a Language from a raw string, normalizing aliases.
func New(s string) Language {
	return id.NewID[languageBrand, string](Normalize(s))
}

// FromSlice converts raw strings into Languages.
func FromSlice(ss []string) []Language {
	result := make([]Language, len(ss))
	for i, s := range ss {
		result[i] = New(s)
	}

	return result
}

// Normalize lowercases and resolves known aliases.
func Normalize(s string) string {
	lower := strings.ToLower(strings.TrimSpace(s))
	if alias, ok := aliases[lower]; ok {
		return alias
	}

	return lower
}

// AliasMap returns a copy of the internal alias mapping.
func AliasMap() map[string]string {
	cp := make(map[string]string, len(aliases))
	maps.Copy(cp, aliases)

	return cp
}

// Languages is an ordered collection of Language values.
type Languages []Language

// NewLanguages creates a Languages collection.
func NewLanguages(langs ...Language) Languages {
	return Languages(langs)
}

// Primary returns the first language, or a zero Language if empty.
func (ls Languages) Primary() Language {
	if len(ls) == 0 {
		return Language{}
	}

	return ls[0]
}

// Has reports whether the collection contains the given language.
func (ls Languages) Has(lang Language) bool {
	return slices.ContainsFunc(ls, func(l Language) bool {
		return l.Equal(lang)
	})
}

// Is is an alias for Has.
func (ls Languages) Is(lang Language) bool {
	return ls.Has(lang)
}

// IsGo reports whether the collection includes Go.
func (ls Languages) IsGo() bool {
	return ls.Has(New("go"))
}

// Strings returns the raw string values of all languages.
func (ls Languages) Strings() []string {
	result := make([]string, len(ls))
	for i, l := range ls {
		result[i] = l.Get()
	}

	return result
}

// IsZero reports whether the collection is empty.
func (ls Languages) IsZero() bool {
	return len(ls) == 0
}

// Validate implements validate.Validator.
func (ls Languages) Validate() error {
	return nil
}

var _ validate.Validator = Languages{}
