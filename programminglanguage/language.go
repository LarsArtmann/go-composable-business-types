package programminglanguage

import (
	"slices"
	"strings"

	id "github.com/larsartmann/go-branded-id"
	"github.com/larsartmann/go-composable-business-types/validate"
)

type languageBrand struct{}

type Language = id.ID[languageBrand, string]

var aliases = map[string]string{
	"golang":    "go",
	"js":        "javascript",
	"mjs":       "javascript",
	"ts":        "typescript",
	"tsx":       "typescript",
	"py":        "python",
	"rs":        "rust",
	"jvm":       "java",
	"c#":        "csharp",
	"cs":        "csharp",
	"rb":        "ruby",
	"gdscript":  "godot",
	"terraform": "hcl",
	"tf":        "hcl",
	"c++":       "cpp",
	"bash":      "shell",
	"sh":        "shell",
	"zsh":       "shell",
	"fish":      "shell",
	"ps1":       "powershell",
	"nodejs":    "node",
	"typescript-jsx": "typescript",
}

func New(s string) Language {
	return id.NewID[languageBrand, string](Normalize(s))
}

func FromSlice(ss []string) []Language {
	result := make([]Language, len(ss))
	for i, s := range ss {
		result[i] = New(s)
	}

	return result
}

func Normalize(s string) string {
	lower := strings.ToLower(strings.TrimSpace(s))
	if alias, ok := aliases[lower]; ok {
		return alias
	}

	return lower
}

func AliasMap() map[string]string {
	cp := make(map[string]string, len(aliases))
	for k, v := range aliases {
		cp[k] = v
	}

	return cp
}

type Languages []Language

func NewLanguages(langs ...Language) Languages {
	return Languages(langs)
}

func (ls Languages) Primary() Language {
	if len(ls) == 0 {
		return Language{}
	}

	return ls[0]
}

func (ls Languages) Has(lang Language) bool {
	return slices.ContainsFunc(ls, func(l Language) bool {
		return l.Equal(lang)
	})
}

func (ls Languages) Is(lang Language) bool {
	return ls.Has(lang)
}

func (ls Languages) IsGo() bool {
	return ls.Has(New("go"))
}

func (ls Languages) Strings() []string {
	result := make([]string, len(ls))
	for i, l := range ls {
		result[i] = l.Get()
	}

	return result
}

func (ls Languages) IsZero() bool {
	return len(ls) == 0
}

func (ls Languages) Validate() error {
	return nil
}

var _ validate.Validator = Languages{}
