# Crumb Dashboard Redesign Plan

## Requirements Summary

| Requirement | Decision |
|-------------|----------|
| Layout | Full task list + Ideas section + Notes section (no Next button) |
| Content | Ideas: 3 most recent; Notes: 3 most recent (including latest added) |
| Visual | Professional aesthetic, cohesive color palette, color-coded status, color accents |
| Output | 3 distinct design versions for review |

---

## Current State

- `cmd/dashboard.go`: ASCII table (wide) + stacked (narrow) — uses `| - +` only
- `helpers/output.go`: ANSI colors (Green, Yellow, Red, Cyan, Gray, Bold) + `FormatStatus()` returning `[+]/[x]/[ ]`
- `store/json.go`: `Task{ID, Text, Status}`, `CrumbData{Next, Tasks, Ideas, Notes, Done}`

---

## Design Version 1: Compact/Dense (Single-Column, Status Dots)

**Target**: Narrow terminals, maximum information density

```
╭────────────────────────────────────────────╮
│  Crumb Dashboard                           │
├────────────────────────────────────────────┤
│  TASKS                                     │
│  ● #a1b  fix login redirect loop      [ ]  │
│  ● #c3d  deploy to staging            [+]  │
│  ● #e5f  refactor auth                [x]  │
│  ● #g7h  write tests                  [ ]  │
├────────────────────────────────────────────┤
│  IDEAS (3 recent)                          │
│  • rate limiter middleware                 │
│  • jwt rewrite                             │
│  • dashboard redesign                      │
├────────────────────────────────────────────┤
│  NOTES (3 recent, latest first)            │
│  • pipe stderr to file                     │
│  • jwt rotation pattern                    │
│  • meeting notes                           │
╰────────────────────────────────────────────╯
```

**Colors**: 
- Status: Green `[+]`, Yellow `[ ]`, Red `[x]` (via `FormatStatus`)
- Section headers: Cyan/Bold
- Task bullets: Gray
- Ideas/Notes bullets: Dim

**Implementation**: Extend `renderStacked` with color accents, limit Ideas/Notes to 3.

---

## Design Version 2: Panel/Box Style (Unicode Box Drawing)

**Target**: Wide terminals (≥80 cols), visual separation with panels

```
┌──────────────────────────────────────────────────────────────┐
│  Crumb Dashboard                                              │
├──────────────────────────────────────────────────────────────┤
│  TASKS                                    │  IDEAS (3)        │
│  ┌────────────────────────────────────┐  │  ┌─────────────┐  │
│  │ #a1b  fix login redirect loop [ ]  │  │  │ rate limit  │  │
│  │ #c3d  deploy to staging       [+]  │  │  │ jwt rewrite │  │
│  │ #e5f  refactor auth           [x]  │  │  │ dashboard   │  │
│  │ #g7h  write tests             [ ]  │  │  └─────────────┘  │
│  └────────────────────────────────────┘  │                   │
│                                           │  NOTES (3)        │
│                                           │  ┌─────────────┐  │
│                                           │  │ pipe stderr │  │
│                                           │  │ jwt rotation│  │
│                                           │  │ meeting     │  │
│                                           │  └─────────────┘  │
└──────────────────────────────────────────────────────────────┘
```

**Colors**:
- Panel borders: Cyan
- Section titles: Bold + Cyan
- Status badges: Green/Yellow/Red (colored)
- Task text: White
- Ideas/Notes: Gray

**Implementation**: New `renderPanel` function using `┌─┐│└┘├┤┬┴` box chars. Two-column: Tasks (left, 60%) + Ideas/Notes stacked (right, 40%). Auto-fallback to Version 1 when width < 80.

---

## Design Version 3: Minimal Modern (Clean ASCII, Breathing Room)

**Target**: Modern aesthetic, generous whitespace, subtle color

```
┌ Crumb Dashboard ────────────────────────────┐
│                                             │
│  TASKS                                       │
│  ────────────────────────────────────────   │
│  #a1b  fix login redirect loop       [ ]    │
│  #c3d  deploy to staging             [+]    │
│  #e5f  refactor auth                 [x]    │
│  #g7h  write tests                   [ ]    │
│                                             │
│  IDEAS (3)              NOTES (3)           │
│  ─────────────          ─────────────       │
│  rate limit             pipe stderr         │
│  jwt rewrite            jwt rotation        │
│  dashboard              meeting notes       │
│                                             │
└─────────────────────────────────────────────┘
```

**Colors**:
- Border: Dim Gray
- Section headers: Bold Cyan
- Dividers: Dim Gray
- Status: Green `[+]`, Yellow `[ ]`, Red `[x]`
- Task ID: Gray
- Task text: White
- Ideas/Notes: Light Gray

**Implementation**: New `renderMinimal` function. Full-width border, generous padding. Two-column Ideas/Notes at bottom. Auto-fallback to Version 1 when width < 70.

---

## Shared Implementation Details

### Content Filtering (all versions)

```go
// Get 3 most recent ideas (Ideas are appended, so last 3)
recentIdeas := data.Ideas
if len(recentIdeas) > 3 {
    recentIdeas = recentIdeas[len(recentIdeas)-3:]
}
// Reverse for display (newest first)
for i, j := 0, len(recentIdeas)-1; i < j; i, j = i+1, j-1 {
    recentIdeas[i], recentIdeas[j] = recentIdeas[j], recentIdeas[i]
}

// Get 3 most recent notes (Notes are appended, so last 3)
recentNotes := data.Notes
if len(recentNotes) > 3 {
    recentNotes = recentNotes[len(recentNotes)-3:]
}
// Reverse for display (newest first, latest added shown first)
for i, j := 0, len(recentNotes)-1; i < j; i, j = i+1, j-1 {
    recentNotes[i], recentNotes[j] = recentNotes[j], recentNotes[i]
}
```

### Color Palette (extend `helpers/output.go`)

```go
const (
    // Existing
    Reset  = "\033[0m"
    Bold   = "\033[1m"
    Green  = "\033[32m"
    Yellow = "\033[33m"
    Red    = "\033[31m"
    Cyan   = "\033[36m"
    Gray   = "\033[90m"
    
    // New for dashboard
    White      = "\033[97m"
    LightGray  = "\033[37m"
    DarkGray   = "\033[90m"
    Blue       = "\033[34m"
    Magenta    = "\033[35m"
)
```

Helper functions:
- `Header(s string)` → Bold + Cyan
- `Section(s string)` → Bold + Cyan  
- `Muted(s string)` → Gray/DarkGray
- `Accent(s string)` → Blue/Magenta for highlights

### Width Detection & Fallback

| Version | Min Width | Fallback |
|---------|-----------|----------|
| Panel/Box | 80 | Compact/Dense |
| Minimal Modern | 70 | Compact/Dense |
| Compact/Dense | 40 | (base) |

---

## File Changes

| File | Change |
|------|--------|
| `cmd/dashboard.go` | Add `renderPanel`, `renderMinimal`; update `renderDashboard` to dispatch by width/version flag |
| `helpers/output.go` | Add color constants + helper functions (Header, Section, Muted, Accent) |
| `cmd/root.go` | Add `--design` flag (compact|panel|minimal) to select version for testing |

---

## Validation Checklist

- [ ] All 3 versions render without errors at 40, 60, 80, 120, 200 cols
- [ ] Ideas limited to 3 most recent (newest first)
- [ ] Notes limited to 3 most recent (newest first, latest added visible)
- [ ] Status colors consistent: Green `[+]`, Yellow `[ ]`, Red `[x]`
- [ ] Professional color palette applied (no clashing colors)
- [ ] No Next button/navigation element
- [ ] Full task list visible in all versions
- [ ] `go build ./...` and `go vet ./...` pass
- [ ] `--design` flag works for manual testing

---

## Rollout

1. Implement Version 1 (Compact/Dense) as enhanced `renderStacked`
2. Implement Version 2 (Panel/Box) as new `renderPanel` 
3. Implement Version 3 (Minimal Modern) as new `renderMinimal`
4. Add `--design` flag to `RootCmd` for version selection
5. Update `renderDashboard` to auto-select by width or use flag
6. Test all widths, empty states, color rendering