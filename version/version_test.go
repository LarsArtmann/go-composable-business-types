package version

import (
	"strings"
	"testing"
)

//nolint:paralleltest // Tests modify global state (Revision, Date, Dirty)
func TestVersion(t *testing.T) {
	if Version == "" {
		t.Error("Version should not be empty")
	}
}

//nolint:paralleltest // Tests modify global state
func TestModulePath(t *testing.T) {
	if ModulePath == "" {
		t.Error("ModulePath should not be empty")
	}
	if !strings.Contains(ModulePath, "go-composable-business-types") {
		t.Error("ModulePath should contain module name")
	}
}

//nolint:paralleltest // Tests modify global state
func TestString(t *testing.T) {
	s := String()
	if s == "" {
		t.Error("String() should not be empty")
	}
	if !strings.Contains(s, Version) {
		t.Error("String() should contain Version")
	}
}

//nolint:paralleltest // Tests modify global state
func TestStringWithRevision(t *testing.T) {
	original := Revision
	Revision = "abc1234567890"
	defer func() { Revision = original }()

	s := String()
	if !strings.Contains(s, "abc1234") {
		t.Error("String() should contain shortened revision")
	}
}

//nolint:paralleltest // Tests modify global state
func TestStringWithDate(t *testing.T) {
	original := Date
	Date = "2026-03-27T10:00:00Z"
	defer func() { Date = original }()

	s := String()
	if !strings.Contains(s, "2026-03-27") {
		t.Error("String() should contain date")
	}
}

//nolint:paralleltest // Tests modify global state
func TestStringDirty(t *testing.T) {
	original := Dirty
	Dirty = true
	defer func() { Dirty = original }()

	s := String()
	if !strings.Contains(s, "dirty") {
		t.Error("String() should contain 'dirty' when Dirty is true")
	}
}
