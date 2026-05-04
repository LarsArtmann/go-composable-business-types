# Post-Migration Improvement Plan

**Created:** 2026-05-04
**Status:** Prioritized by Impact / Effort

---

## CRITICAL FINDINGS (What was forgotten)

### F1: PMA's `types.Tag` was NEVER migrated to shared `tag.Tag`
- PMA still has its own `Tag string` with `^[a-z0-9_-]+$` regex (allows underscores, forces lowercase)
- Shared CBT `Tag` uses `^[A-Za-z0-9-]+$` (no underscores, allows uppercase)
- The meta adapter silently corrupts tags during conversion: `My-Tag` → `my-tag`, `my_tag` → `my-tag`
- **17 usage sites** in PMA reference `types.Tag`

### F2: Adapter still needed for Tag conversion
- `ConvertTag`, `ConvertTags`, `ConvertTagToMeta`, `ConvertTagsToMeta` in `meta/adapter.go`
- If PMA used shared `tag.Tag`, these become identity operations and can be deleted

### F3: CBT Importance missing arithmetic methods
- PMA needs `Add`, `Sub`, `Clamp`, `Max`, `Min`, `IsDefault` that CBT doesn't provide
- These should be added to CBT so PMA can eventually become a type alias

---

## PLAN — Sorted by Impact × Value / Effort

### Tier 1: Fix Type Split-Brain (HIGH impact, MED effort)

| # | Task | Effort | Impact |
|---|------|--------|--------|
| T1.1 | Add `Add`, `Sub`, `Clamp`, `Max`, `Min`, `IsDefault` to CBT Importance | 20m | HIGH |
| T1.2 | Add `Tags` slice type + `Strings`/`IsEmpty` helpers to CBT tag package | 10m | MED |
| T1.3 | Replace PMA `types.Tag` → CBT `tag.Tag` (17 sites) | 30m | HIGH |
| T1.4 | Delete tag conversion functions from meta/adapter.go | 10m | MED |
| T1.5 | Fix adapter.go comment (claims identical regex — it's wrong) | 2m | LOW |

### Tier 2: Cleanup & Consistency (MED impact, LOW effort)

| # | Task | Effort | Impact |
|---|------|--------|--------|
| T2.1 | Clean remaining `int32` refs in PMA test code | 10m | MED |
| T2.2 | Use `Classification()` for human-readable Importance in CLI output | 10m | MED |
| T2.3 | Fix PMA docs that still say `type Importance int32` | 5m | LOW |
| T2.4 | Add `go.work` to project-meta if missing | 5m | LOW |

### Tier 3: Architecture Improvements (MED impact, MED effort)

| # | Task | Effort | Impact |
|---|------|--------|--------|
| T3.1 | Consider extracting `ProjectType` from PMA's Language enum | 30m | MED |
| T3.2 | Consider auto-generating PMA Language enum maps | 15m | LOW |
| T3.3 | Add `Percent() string` to CBT Importance (match PMA's interface) | 5m | LOW |

---

## Execution Order

**Phase A (CBT enrichment):** T1.1 → T1.2 → T3.3
**Phase B (PMA Tag migration):** T1.3 → T1.4 → T1.5
**Phase C (Cleanup):** T2.1 → T2.2 → T2.3 → T2.4
**Phase D (Architecture):** T3.1 → T3.2 (optional, defer)
