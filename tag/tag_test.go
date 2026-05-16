package tag

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    Tag
		wantErr bool
	}{
		{name: "simple", input: "go", want: "go"},
		{name: "hyphenated", input: "my-project", want: "my-project"},
		{name: "uppercase", input: "Go", want: "Go"},
		{name: "mixed", input: "My-Project-123", want: "My-Project-123"},
		{name: "numbers_only", input: "123", want: "123"},
		{name: "empty", input: "", wantErr: true},
		{name: "underscore", input: "my_tag", wantErr: true},
		{name: "space", input: "my tag", wantErr: true},
		{name: "special_char", input: "tag!", wantErr: true},
		{name: "dot", input: "tag.name", wantErr: true},
		{name: "too_long", input: strings.Repeat("a", 51), wantErr: true},
		{name: "max_length", input: strings.Repeat("a", 50), want: Tag(strings.Repeat("a", 50))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := New(tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestMust(t *testing.T) {
	t.Parallel()

	assert.Equal(t, Tag("go"), Must("go"))
	assert.Panics(t, func() { Must("") })
	assert.Panics(t, func() { Must("invalid_tag") })
}

func TestNewTags(t *testing.T) {
	t.Parallel()

	t.Run("valid batch", func(t *testing.T) {
		t.Parallel()

		tags, err := NewTags("go", "rust", "My-Project")
		require.NoError(t, err)
		assert.Equal(t, []Tag{"go", "rust", "My-Project"}, tags)
	})

	t.Run("invalid in batch", func(t *testing.T) {
		t.Parallel()

		_, err := NewTags("go", "invalid tag", "rust")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "index 1")
	})
}

func TestString(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "go", Tag("go").String())
}

func TestIsZero(t *testing.T) {
	t.Parallel()

	assert.True(t, Tag("").IsZero())
	assert.False(t, Tag("go").IsZero())
}

func TestIsValid(t *testing.T) {
	t.Parallel()

	assert.True(t, Tag("go").IsValid())
	assert.False(t, Tag("").IsValid())
	assert.False(t, Tag("has space").IsValid())
}

func TestValidate(t *testing.T) {
	t.Parallel()

	require.NoError(t, Tag("go").Validate())
	require.Error(t, Tag("").Validate())
	assert.Error(t, Tag("invalid_tag").Validate())
}

func TestJSONRoundTrip(t *testing.T) {
	t.Parallel()

	original := Tag("my-project")

	data, err := json.Marshal(original)
	require.NoError(t, err)
	assert.JSONEq(t, `"my-project"`, string(data))

	var got Tag

	err = json.Unmarshal(data, &got)
	require.NoError(t, err)
	assert.Equal(t, original, got)
}

func TestJSONInStruct(t *testing.T) {
	t.Parallel()

	type Project struct {
		Tags []Tag `json:"tags"`
	}

	p := Project{Tags: []Tag{"go", "Rust", "My-Project"}}

	data, err := json.Marshal(p)
	require.NoError(t, err)
	assert.JSONEq(t, `{"tags":["go","Rust","My-Project"]}`, string(data))

	var got Project

	err = json.Unmarshal(data, &got)
	require.NoError(t, err)
	assert.Equal(t, p.Tags, got.Tags)
}

func TestJSONUnmarshalInvalid(t *testing.T) {
	t.Parallel()

	var got Tag

	err := json.Unmarshal([]byte(`123`), &got)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid JSON")
}

func TestScanValue(t *testing.T) {
	t.Parallel()

	t.Run("string source", func(t *testing.T) {
		t.Parallel()

		var got Tag

		err := got.Scan("go")
		require.NoError(t, err)
		assert.Equal(t, Tag("go"), got)
	})

	t.Run("nil source", func(t *testing.T) {
		t.Parallel()

		var got Tag

		err := got.Scan(nil)
		require.NoError(t, err)
		assert.Equal(t, Tag(""), got)
	})

	t.Run("value", func(t *testing.T) {
		t.Parallel()

		v, err := Tag("go").Value()
		require.NoError(t, err)
		assert.Equal(t, "go", v)
	})

	t.Run("zero value returns nil", func(t *testing.T) {
		t.Parallel()

		v, err := Tag("").Value()
		require.NoError(t, err)
		assert.Nil(t, v)
	})
}

func TestNewTagsFromString(t *testing.T) {
	t.Parallel()

	t.Run("valid batch", func(t *testing.T) {
		t.Parallel()

		tags, err := NewTagsFromString("go", "rust", "My-Project")
		require.NoError(t, err)
		assert.Equal(t, Tags{"go", "rust", "My-Project"}, tags)
	})

	t.Run("empty input", func(t *testing.T) {
		t.Parallel()

		tags, err := NewTagsFromString()
		require.NoError(t, err)
		assert.Equal(t, Tags{}, tags)
	})

	t.Run("invalid in batch", func(t *testing.T) {
		t.Parallel()

		_, err := NewTagsFromString("go", "invalid tag", "rust")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "index 1")
	})
}

func TestTagsStrings(t *testing.T) {
	t.Parallel()

	tags := Tags{"go", "rust", "My-Project"}
	strings := tags.Strings()
	assert.Equal(t, []string{"go", "rust", "My-Project"}, strings)
}

func TestTagsIsEmpty(t *testing.T) {
	t.Parallel()

	assert.True(t, Tags{}.IsEmpty())
	assert.True(t, Tags(nil).IsEmpty())
	assert.False(t, Tags{"go"}.IsEmpty())
}

func TestTagsContains(t *testing.T) {
	t.Parallel()

	tags := Tags{"go", "rust"}

	assert.True(t, tags.Contains("go"))
	assert.True(t, tags.Contains("rust"))
	assert.False(t, tags.Contains("python"))
	assert.False(t, Tags{}.Contains("go"))
}
