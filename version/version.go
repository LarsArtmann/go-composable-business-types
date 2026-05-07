// Package version provides build-time version information extracted from runtime/debug.
//
// It exposes Version, Revision (git commit), Date, and Dirty flag, populated
// automatically at init time from Go's build info.
//
// Basic usage:
//
//	fmt.Println(version.String()) // "development+abc1234+2026-05-07+dirty"
package version

import (
	"runtime/debug"
	"strings"
)

const (
	Version    = "development"
	ModulePath = "github.com/larsartmann/go-composable-business-types"

	shortRevisionLength = 7
)

//nolint:gochecknoinits // init is required to read build info at runtime
func init() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			Revision = setting.Value
		case "vcs.time":
			Date = setting.Value
		case "vcs.modified":
			Dirty = setting.Value == "true"
		}
	}
}

var (
	Revision string
	Date     string
	Dirty    bool
)

func String() string {
	var parts []string

	parts = append(parts, Version)

	if Revision != "" {
		short := Revision
		if len(short) > shortRevisionLength {
			short = short[:shortRevisionLength]
		}

		parts = append(parts, short)
	}

	if Date != "" {
		if iso, _, _ := strings.Cut(Date, "T"); iso != "" {
			parts = append(parts, iso)
		}
	}

	if Dirty {
		parts = append(parts, "dirty")
	}

	return strings.Join(parts, "+")
}
