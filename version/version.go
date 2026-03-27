package version

import (
	"runtime/debug"
	"strings"
)

const (
	Version    = "development"
	ModulePath = "github.com/larsartmann/go-composable-business-types"
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
		if len(short) > 7 {
			short = short[:7]
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
