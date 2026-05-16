package projectcore

import (
	"encoding/json"
	"testing"

	"github.com/larsartmann/go-composable-business-types/importance"
	"github.com/larsartmann/go-composable-business-types/tag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	langs := []string{"go"}
	p := New("myproject", "/home/user/myproject", langs)

	assert.Equal(t, "myproject", p.Name)
	assert.Equal(t, "/home/user/myproject", p.Path)
	assert.Equal(t, "go", p.PrimaryLanguage())
	assert.True(t, p.Importance.IsZero())
	assert.Empty(t, p.Tags)
}

func TestNewWithOptions(t *testing.T) {
	t.Parallel()

	langs := []string{"go"}
	imp := importance.Must(70)
	t1, _ := tag.New("backend")
	t2, _ := tag.New("production")

	p := New("myproject", "/path", langs,
		WithImportance(imp),
		WithTags(t1, t2),
	)

	assert.Equal(t, imp, p.Importance)
	require.Len(t, p.Tags, 2)
	assert.Equal(t, "backend", p.Tags[0].String())
	assert.Equal(t, "production", p.Tags[1].String())
}

func TestIsZero(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()

		var p *ProjectCore
		assert.True(t, p.IsZero())
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		p := &ProjectCore{}
		assert.True(t, p.IsZero())
	})

	t.Run("with name", func(t *testing.T) {
		t.Parallel()

		p := &ProjectCore{Name: "test"}
		assert.False(t, p.IsZero())
	})
}

func TestValidate(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()

		var p *ProjectCore
		require.Error(t, p.Validate())
	})

	t.Run("empty name", func(t *testing.T) {
		t.Parallel()

		p := &ProjectCore{Path: "/path"}
		require.Error(t, p.Validate())
		assert.Contains(t, p.Validate().Error(), "name is required")
	})

	t.Run("empty path", func(t *testing.T) {
		t.Parallel()

		p := &ProjectCore{Name: "test"}
		require.Error(t, p.Validate())
		assert.Contains(t, p.Validate().Error(), "path is required")
	})

	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		p := New("test", "/path", []string{"go"})
		require.NoError(t, p.Validate())
	})
}

func TestJSONRoundTrip(t *testing.T) {
	t.Parallel()

	p := New(
		"myproject",
		"/home/user/myproject",
		[]string{"go", "python"},
		WithImportance(importance.Must(75)),
		WithTags(tag.Must("backend"), tag.Must("production")),
	)

	data, err := json.Marshal(p)
	require.NoError(t, err)

	var got ProjectCore

	err = json.Unmarshal(data, &got)
	require.NoError(t, err)

	assert.Equal(t, "myproject", got.Name)
	assert.Equal(t, "/home/user/myproject", got.Path)
	assert.Equal(t, []string{"go", "python"}, got.Languages)
	assert.Equal(t, importance.Importance(75), got.Importance)
	require.Len(t, got.Tags, 2)
	assert.Equal(t, "backend", got.Tags[0].String())
	assert.Equal(t, "production", got.Tags[1].String())
}

func TestJSONOutput(t *testing.T) {
	t.Parallel()

	p := New("test", "/path",
		[]string{"go"},
		WithImportance(importance.High),
	)

	data, err := json.Marshal(p)
	require.NoError(t, err)
	assert.JSONEq(
		t,
		`{"name":"test","path":"/path","languages":["go"],"importance":70,"tags":null}`,
		string(data),
	)
}
