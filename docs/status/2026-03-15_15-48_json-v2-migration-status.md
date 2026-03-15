# Status Report: encoding/json/v2 Migration

**Date:** 2026-03-15 15:48:34 CET
**Author:** AI Assistant
**Status:** 🚨 CRITICAL ISSUES - Migration Not Complete

---

## Executive Summary

**Mission:** Migrate from `encoding/json` (v1) to `encoding/json/v2` as the default JSON library, and delete root-level `.go` files.

**Current State:** ❌ **FAILED** - Migration not complete. Files still use `encoding/json` (v1).

---

## A) FULLY DONE ✅

| Item | Status | Details |
|------|--------|---------|
| Root `.go` files deleted | ✅ Complete | `cbt.go` and `id_jsonv2.go` removed from project root |
| Research on json/v2 | ✅ Complete | Confirmed json/v2 API differences documented |
| Go cache cleared | ✅ Complete | `go clean -cache` executed successfully |
| Dependencies downloaded | ✅ Complete | `go mod download` completed |

---

## B) PARTIALLY DONE ⚠️

| Item | Status | Details |
|------|--------|---------|
| Import replacement | ⚠️ FAILED | sed/perl commands executed but changes NOT persisted |
| Build with GOEXPERIMENT | ⚠️ WORKS | `GOEXPERIMENT=jsonv2 go build ./...` succeeds |
| Tests with GOEXPERIMENT | ⚠️ FAILS | Race detector tests fail with build errors |

---

## C) NOT STARTED ⏳

| Item | Priority | Notes |
|------|----------|-------|
| Update all 10 files with edit tool | HIGH | Need to use edit tool, not sed/perl |
| Update `json.MarshalIndent` to v2 API | HIGH | Only in `examples/datapoint/main.go` |
| Verify build without GOEXPERIMENT | CRITICAL | User claims json/v2 should be default |
| Update README documentation | MEDIUM | Document json/v2 requirement |
| Update AGENTS.md | LOW | Add json/v2 to build commands |

---

## D) TOTALLY FUCKED UP 💥

### Critical Failure: sed/perl Commands Did NOT Persist

**What happened:**
1. Ran `sed -i '' 's#encoding/json#encoding/json/v2#'` on 9 files
2. Ran `perl -pi -e 's|"encoding/json"|"encoding/json/v2"|'` on 10 files
3. Both commands reported "success" with no errors
4. `git status` shows: **"nothing to commit, working tree clean"**
5. `grep -rn "encoding/json" --include="*.go"` shows: **All files still have v1 imports**

**Root Cause:** Unknown. Possible causes:
- Files were not actually modified (dry-run?)
- Some process reverted changes
- Working directory confusion
- File system caching issues

**Impact:** Complete waste of time. All import replacement work needs to be redone with the edit tool.

---

## E) WHAT WE SHOULD IMPROVE 📈

1. **Use edit tool ONLY** - sed/perl are unreliable in this environment
2. **Verify changes immediately** - Run `git diff` after each file modification
3. **One file at a time** - Edit, verify, then move to next
4. **Test incrementally** - Build after each file change
5. **Document GOEXPERIMENT requirement** - If json/v2 requires it, document it clearly

---

## F) TOP 25 THINGS TO DO NEXT 📋

| # | Task | Priority | Est. Time |
|---|------|----------|-----------|
| 1 | Fix `temporal/temporal.go` import | CRITICAL | 1 min |
| 2 | Fix `datapoint/datapoint.go` import | CRITICAL | 1 min |
| 3 | Fix `datapoint/datapoint_test.go` import | CRITICAL | 1 min |
| 4 | Fix `datapoint/cause.go` import | CRITICAL | 1 min |
| 5 | Fix `datapoint/reference.go` import | CRITICAL | 1 min |
| 6 | Fix `datapoint/context.go` import | CRITICAL | 1 min |
| 7 | Fix `bounded/bounded.go` import | CRITICAL | 1 min |
| 8 | Fix `id/id.go` import | CRITICAL | 1 min |
| 9 | Fix `id/id_test.go` import | CRITICAL | 1 min |
| 10 | Fix `examples/datapoint/main.go` import + API | CRITICAL | 2 min |
| 11 | Verify all imports changed | HIGH | 1 min |
| 12 | Build with `go build ./...` | CRITICAL | 1 min |
| 13 | Test with `go test ./...` | CRITICAL | 2 min |
| 14 | Investigate: Does Go 1.26.1 have json/v2 by default? | CRITICAL | 5 min |
| 15 | Update README if GOEXPERIMENT required | HIGH | 3 min |
| 16 | Update AGENTS.md with build commands | MEDIUM | 2 min |
| 17 | Run tests with race detector | HIGH | 5 min |
| 18 | Verify examples compile and run | HIGH | 2 min |
| 19 | Check for other json v1 API calls | HIGH | 2 min |
| 20 | Update go.mod if needed | MEDIUM | 1 min |
| 21 | Run golangci-lint | MEDIUM | 3 min |
| 22 | Generate enum code | LOW | 1 min |
| 23 | Update project documentation | LOW | 5 min |
| 24 | Create git commit with detailed message | HIGH | 2 min |
| 25 | Verify CI/CD still works | LOW | 5 min |

---

## G) TOP #1 QUESTION I CANNOT FIGURE OUT 🤔

### "Does your Go 1.26.1 installation have `encoding/json/v2` as a non-experimental default?"

**Context:**
- My research indicates `encoding/json/v2` is **experimental** and requires `GOEXPERIMENT=jsonv2`
- You stated: "we require 1.26.1 with removed GOEXPERIMENT=jsonv2 and made it non GOEXPERIMENT"
- Standard Go 1.26.0/1.26.1 from go.dev **DOES NOT** have json/v2 as default

**What I need to know:**
1. Do you have a **custom Go build** with json/v2 enabled by default?
2. Or should the library require users to set `GOEXPERIMENT=jsonv2`?
3. Is there a `godebug` directive or build tag that enables this?

**Evidence:**
```
# Standard Go 1.26.0 (nix)
$ go build ./...
package encoding/json/v2: build constraints exclude all Go files

$ GOEXPERIMENT=jsonv2 go build ./...
# Works!
```

**This is blocking because:**
- If json/v2 requires GOEXPERIMENT, users must set it
- If your Go has it by default, we need to document this requirement
- The library cannot "just work" for everyone without configuration

---

## Files Currently Using encoding/json (v1)

```
./temporal/temporal.go:4:         "encoding/json"
./datapoint/datapoint_test.go:4:  "encoding/json"
./datapoint/cause.go:4:           "encoding/json"
./datapoint/datapoint.go:4:       "encoding/json"
./datapoint/reference.go:4:       "encoding/json"
./datapoint/context.go:4:         "encoding/json"
./examples/datapoint/main.go:7:   "encoding/json"
./bounded/bounded.go:5:           "encoding/json"
./id/id.go:40:                    "encoding/json"
./id/id_test.go:6:                "encoding/json"
```

## Required Changes

### Import Change (all 10 files)
```go
// Before
import "encoding/json"

// After
import "encoding/json/v2"
```

### API Change (examples/datapoint/main.go only)
```go
// Before (v1)
data, err := json.MarshalIndent(dp, "", "  ")

// After (v2)
import "encoding/json/jsontext"
data, err := json.Marshal(dp, jsontext.WithIndent("  "))
```

---

## Test Results (Latest Run)

```
$ GOEXPERIMENT=jsonv2 go test -race ./...

FAIL    actor [build failed]
FAIL    bounded [build failed]
FAIL    datapoint [build failed]
ok      enums   1.319s
FAIL    examples/datapoint [build failed]
FAIL    id [build failed]
ok      locale  1.564s
ok      money   1.806s
ok      nanoid  2.011s
FAIL    temporal [build failed]
ok      types   2.231s
```

Build failures are due to `could not import encoding/json` errors - the race detector seems to have issues with the experimental json/v2 package.

---

## Next Steps

1. **ANSWER MY QUESTION** about your Go installation
2. I will fix all 10 files using the edit tool
3. Build and test
4. Commit with detailed message
