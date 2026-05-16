package version

import (
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	if Version == "" {
		t.Error("Version should not be empty")
	}
}

func TestModulePath(t *testing.T) {
	if ModulePath == "" {
		t.Error("ModulePath should not be empty")
	}

	if !strings.Contains(ModulePath, "go-composable-business-types") {
		t.Error("ModulePath should contain module name")
	}
}

func TestString(t *testing.T) {
	s := String()
	if s == "" {
		t.Error("String() should not be empty")
	}

	if !strings.Contains(s, Version) {
		t.Error("String() should contain Version")
	}
}

func TestStringWithRevision(t *testing.T) {
	testWithGlobal(t, &Revision, "abc1234567890", "abc1234", "shortened revision")
}

func TestStringWithDate(t *testing.T) {
	testWithGlobal(t, &Date, "2026-03-27T10:00:00Z", "2026-03-27", "date")
}

func TestStringDirty(t *testing.T) {
	testWithGlobal(t, &Dirty, true, "dirty", "'dirty' when Dirty is true")
}

// testWithGlobal saves the original value of a global variable, sets a new value,
// runs a test using String(), then restores the original value.
func testWithGlobal[T any](
	t *testing.T,
	global *T,
	newValue T,
	expectedSubstring, description string,
) {
	original := *global
	*global = newValue

	defer func() { *global = original }()

	s := String()
	if !strings.Contains(s, expectedSubstring) {
		t.Errorf("String() should contain %s", description)
	}
}
