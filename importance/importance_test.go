package importance

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   uint8
		want    Importance
		wantErr bool
	}{
		{name: "zero", input: 0, want: None},
		{name: "fifty", input: 50, want: Medium},
		{name: "hundred", input: 100, want: Max},
		{name: "overflows", input: 101, wantErr: true},
		{name: "max_uint8", input: 255, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := New(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "exceeds maximum")
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestMust(t *testing.T) {
	t.Parallel()

	assert.Equal(t, Medium, Must(50))
	assert.Panics(t, func() { Must(101) })
}

func TestParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    Importance
		wantErr bool
	}{
		{name: "none", input: "none", want: None},
		{name: "very-low", input: "very-low", want: VeryLow},
		{name: "very_low_underscore", input: "very_low", wantErr: true},
		{name: "low", input: "low", want: Low},
		{name: "medium", input: "medium", want: Medium},
		{name: "mid", input: "mid", want: Medium},
		{name: "high", input: "high", want: High},
		{name: "very-high", input: "very-high", want: VeryHigh},
		{name: "max", input: "max", want: Max},
		{name: "100", input: "100", want: Max},
		{name: "0", input: "0", want: None},
		{name: "uppercase", input: "HIGH", want: High},
		{name: "spaced", input: " medium ", want: Medium},
		{name: "unknown", input: "xyz", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := Parse(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input Importance
		want  string
	}{
		{0, "none"},
		{1, "very-low"},
		{20, "very-low"},
		{21, "low"},
		{40, "low"},
		{41, "medium"},
		{50, "medium"},
		{60, "medium"},
		{61, "high"},
		{80, "high"},
		{81, "very-high"},
		{100, "very-high"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("importance_%d", tt.input), func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, tt.input.String())
		})
	}
}

func TestClassification(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "None", Importance(0).Classification())
	assert.Equal(t, "Very Low", Importance(15).Classification())
	assert.Equal(t, "Low", Importance(30).Classification())
	assert.Equal(t, "Medium", Importance(55).Classification())
	assert.Equal(t, "High", Importance(70).Classification())
	assert.Equal(t, "Very High", Importance(90).Classification())
}

func TestClassificationMethods(t *testing.T) {
	t.Parallel()

	assert.True(t, Importance(15).IsVeryLow())
	assert.False(t, Importance(21).IsVeryLow())

	assert.True(t, Importance(30).IsLow())
	assert.False(t, Importance(41).IsLow())

	assert.True(t, Importance(55).IsMedium())
	assert.False(t, Importance(61).IsMedium())

	assert.True(t, Importance(70).IsHigh())
	assert.False(t, Importance(81).IsHigh())

	assert.True(t, Importance(90).IsVeryHigh())
	assert.False(t, Importance(80).IsVeryHigh())
}

func TestIsValid(t *testing.T) {
	t.Parallel()

	assert.True(t, Importance(0).IsValid())
	assert.True(t, Importance(50).IsValid())
	assert.True(t, Importance(100).IsValid())
	assert.False(t, Importance(101).IsValid())
	assert.False(t, Importance(255).IsValid())
}

func TestIsZero(t *testing.T) {
	t.Parallel()

	assert.True(t, Importance(0).IsZero())
	assert.False(t, Importance(1).IsZero())
}

func TestPercent(t *testing.T) {
	t.Parallel()

	assert.InDelta(t, 0.0, Importance(0).Percent(), 0.001)
	assert.InDelta(t, 0.5, Importance(50).Percent(), 0.001)
	assert.InDelta(t, 1.0, Importance(100).Percent(), 0.001)
}

func TestCompare(t *testing.T) {
	t.Parallel()

	assert.Equal(t, -1, Importance(10).Compare(Importance(20)))
	assert.Equal(t, 0, Importance(50).Compare(Importance(50)))
	assert.Equal(t, 1, Importance(80).Compare(Importance(20)))
}

func TestValidate(t *testing.T) {
	t.Parallel()

	assert.NoError(t, Importance(50).Validate())
	assert.NoError(t, Importance(0).Validate())
	assert.Error(t, Importance(101).Validate())
}

func TestJSONRoundTrip(t *testing.T) {
	t.Parallel()

	tests := []Importance{0, 20, 50, 80, 100}

	for _, imp := range tests {
		t.Run(fmt.Sprintf("importance_%d", imp), func(t *testing.T) {
			t.Parallel()

			data, err := json.Marshal(imp)
			require.NoError(t, err)

			var got Importance
			err = json.Unmarshal(data, &got)
			require.NoError(t, err)
			assert.Equal(t, imp, got)
		})
	}
}

func TestJSONUnmarshalInvalid(t *testing.T) {
	t.Parallel()

	var got Importance
	err := json.Unmarshal([]byte(`"not-a-number"`), &got)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid JSON")
}

func TestJSONUnmarshalOverflow(t *testing.T) {
	t.Parallel()

	var got Importance
	err := json.Unmarshal([]byte(`200`), &got)
	require.NoError(t, err) // unmarshal just reads uint8, validation is separate
	assert.False(t, got.IsValid())
}

func TestJSONInStruct(t *testing.T) {
	t.Parallel()

	type Project struct {
		Name       string     `json:"name"`
		Importance Importance `json:"importance"`
	}

	p := Project{Name: "test", Importance: High}

	data, err := json.Marshal(p)
	require.NoError(t, err)
	assert.JSONEq(t, `{"name":"test","importance":70}`, string(data))

	var got Project
	err = json.Unmarshal(data, &got)
	require.NoError(t, err)
	assert.Equal(t, High, got.Importance)
}

func TestScanValue(t *testing.T) {
	t.Parallel()

	t.Run("int64 source", func(t *testing.T) {
		t.Parallel()

		var got Importance
		err := got.Scan(int64(50))
		require.NoError(t, err)
		assert.Equal(t, Medium, got)
	})

	t.Run("nil source", func(t *testing.T) {
		t.Parallel()

		var got Importance
		err := got.Scan(nil)
		require.NoError(t, err)
		assert.Equal(t, None, got)
	})

	t.Run("value", func(t *testing.T) {
		t.Parallel()

		v, err := Medium.Value()
		require.NoError(t, err)
		assert.Equal(t, int64(50), v)
	})

	t.Run("zero value", func(t *testing.T) {
		t.Parallel()

		v, err := None.Value()
		require.NoError(t, err)
		assert.Nil(t, v)
	})
}

func TestScanNilReceiver(t *testing.T) {
	t.Parallel()

	var got *Importance
	err := got.Scan(int64(50))
	require.Error(t, err)
	assert.True(t, errors.Is(err, err) || err != nil)
}

func TestConstants(t *testing.T) {
	t.Parallel()

	assert.Equal(t, Importance(0), None)
	assert.Equal(t, Importance(20), VeryLow)
	assert.Equal(t, Importance(40), Low)
	assert.Equal(t, Importance(50), Medium)
	assert.Equal(t, Importance(70), High)
	assert.Equal(t, Importance(90), VeryHigh)
	assert.Equal(t, Importance(100), Max)
}
