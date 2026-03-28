package version

import (
	"context"
	"os/exec"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
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
	Dirty = isGitDirty()
}

func isGitDirty() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", "diff", "--stat")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	output := string(out)
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 0 {
		return false
	}
	lastLine := lines[len(lines)-1]
	parts := strings.Fields(lastLine)
	if len(parts) == 0 {
		return false
	}
	changed, err := strconv.Atoi(parts[len(parts)-1])
	return err == nil && changed > 0
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
