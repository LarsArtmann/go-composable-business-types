# TODO List

> Short-term, actionable, bounded work items, verified against the actual code.
> For long-term vision and unrefined ideas, use ROADMAP.md.
> Items are ranked by impact. Status is verified, not assumed.

## Status legend

| Status      | Meaning                                                   |
| ----------- | --------------------------------------------------------- |
| TODO        | Not started. Needs doing.                                 |
| IN_PROGRESS | Actively being worked on.                                 |
| BLOCKED     | Cannot proceed, external dependency or decision needed.   |
| DONE        | Completed. Remove from this list and log in CHANGELOG.md. |

## High Impact

| Task                                                | Status | Impact | Effort | Evidence                                                                                                         |
| --------------------------------------------------- | ------ | ------ | ------ | ---------------------------------------------------------------------------------------------------------------- |
| Fix projectcore JSON v2 test failure                | TODO   | High   | 30min  | `projectcore/project_core_test.go:144`, nil vs empty slice for tags after json v2 migration                      |
| Resolve go.mod version vs encoding/json v2 mismatch | TODO   | High   | 30min  | `go.mod` says `go 1.26.4` but code imports `encoding/json/v2`; requires `GOEXPERIMENT=jsonv2` or bump to Go 1.27 |

## Medium Impact

| Task                                             | Status | Impact | Effort | Evidence                                                                                                                         |
| ------------------------------------------------ | ------ | ------ | ------ | -------------------------------------------------------------------------------------------------------------------------------- |
| Replace testify with ginkgo/gomega in 3 packages | TODO   | Med    | 4h     | `projectcore/project_core_test.go`, `tag/tag_test.go`, `importance/importance_test.go` all import testify (banned per AGENTS.md) |
| Add tests for testutil package                   | TODO   | Med    | 2h     | `testutil/parse.go` has 0% test coverage                                                                                         |

## Low Impact

| Task                                                     | Status | Impact | Effort | Evidence                                                                     |
| -------------------------------------------------------- | ------ | ------ | ------ | ---------------------------------------------------------------------------- |
| Fill in docs/DOMAIN_LANGUAGE.md with actual domain terms | TODO   | Low    | 1h     | `docs/DOMAIN_LANGUAGE.md` is still a blank template with placeholder content |
