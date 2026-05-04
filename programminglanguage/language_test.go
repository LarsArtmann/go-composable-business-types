package programminglanguage

import (
	"testing"

	id "github.com/larsartmann/go-branded-id"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	lang := New("Go")
	assert.Equal(t, "go", lang.Get())
}

func TestFromSlice(t *testing.T) {
	t.Parallel()

	langs := FromSlice([]string{"Go", "Python", "Rust"})
	require.Len(t, langs, 3)
	assert.Equal(t, "go", langs[0].Get())
	assert.Equal(t, "python", langs[1].Get())
	assert.Equal(t, "rust", langs[2].Get())
}

func TestNormalize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		want  string
	}{
		{"Go", "go"},
		{"golang", "go"},
		{"js", "javascript"},
		{"TypeScript", "typescript"},
		{"c#", "csharp"},
		{"C++", "cpp"},
		{"bash", "shell"},
		{"nodejs", "node"},
		{"Python", "python"},
		{"unknown-lang", "unknown-lang"},
		{"", ""},
		{"  Go  ", "go"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, Normalize(tt.input))
		})
	}
}

func TestAliasMap(t *testing.T) {
	t.Parallel()

	m := AliasMap()
	assert.Equal(t, "go", m["golang"])
	assert.Equal(t, "csharp", m["c#"])
	// Verify it's a copy — modifying returned map doesn't affect original
	m["newkey"] = "newval"
	_, exists := aliases["newkey"]
	assert.False(t, exists)
}

func TestLanguages_Primary(t *testing.T) {
	t.Parallel()

	t.Run("non-empty", func(t *testing.T) {
		t.Parallel()

		ls := NewLanguages(New("go"), New("python"))
		assert.Equal(t, "go", ls.Primary().Get())
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		ls := NewLanguages()
		assert.True(t, ls.Primary().IsZero())
	})
}

func TestLanguages_Has(t *testing.T) {
	t.Parallel()

	ls := NewLanguages(New("go"), New("python"))
	assert.True(t, ls.Has(New("go")))
	assert.True(t, ls.Has(New("python")))
	assert.False(t, ls.Has(New("rust")))
}

func TestLanguages_Is(t *testing.T) {
	t.Parallel()

	ls := NewLanguages(New("go"))
	assert.True(t, ls.Is(New("go")))
	assert.False(t, ls.Is(New("rust")))
}

func TestLanguages_IsGo(t *testing.T) {
	t.Parallel()

	assert.True(t, NewLanguages(New("go")).IsGo())
	assert.False(t, NewLanguages(New("python")).IsGo())
	assert.False(t, NewLanguages().IsGo())
}

func TestLanguages_Strings(t *testing.T) {
	t.Parallel()

	ls := NewLanguages(New("go"), New("python"))
	assert.Equal(t, []string{"go", "python"}, ls.Strings())
}

func TestLanguages_IsZero(t *testing.T) {
	t.Parallel()

	assert.True(t, NewLanguages().IsZero())
	assert.False(t, NewLanguages(New("go")).IsZero())
}

func TestLanguages_Validate(t *testing.T) {
	t.Parallel()

	assert.NoError(t, NewLanguages().Validate())
	assert.NoError(t, NewLanguages(New("go")).Validate())
}

func TestLanguageCompileTimeDistinctness(t *testing.T) {
	t.Parallel()

	lang := New("go")
	_ = lang

	type otherBrand struct{}

	other := id.NewID[otherBrand, string]("go")

	assert.Equal(t, lang.Get(), other.Get())
	// lang and other have different brands → they are different types at compile time
	// This proves phantom typing works: you cannot accidentally mix Language with another branded ID
	_ = other
}
