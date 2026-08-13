# Crumb Dashboard — Creative Concepts Plan

## Data Structure Recap

```go
type CrumbData struct {
    Next  string   // Current focus
    Tasks []Task   // {ID, Text, Status: pending|done|canceled|failed}
    Ideas []string // Raw ideas
    Notes []string // Raw notes
    Done  []string // Completed task texts (legacy)
}
```

Current dashboard: focus + active task count + recent 3 notes + recent 3 ideas (static, read-only).

---

## Concept 1: Kanban Board (Spatial Layout) — REFINED

**Visual:** Three task columns + dedicated Ideas/Notes panel. No redundant counts.

```
┌─────────────────┬─────────────────┬─────────────────┬─────────────────┐
│   📋 PENDING    │    ✔ DONE       │   ✖ CANCELED    │  💡 IDEAS       │
│  ┌───────────┐  │  ┌───────────┐  │  ┌───────────┐  │  ┌───────────┐  │
│  │ #a1b      │  │  │ #c3d      │  │  │ #e5f      │  │  │ rate lim  │  │
│  │ fix bug   │  │  │ wrote     │  │  │ refactor  │  │  │ middleware│  │
│  │ [1]       │  │  │ tests     │  │  │ auth      │  │  └───────────┘  │
│  └───────────┘  │  └───────────┘  │  └───────────┘  │  ┌───────────┐  │
│  ┌───────────┐  │                 │                 │  │ jwt       │  │
│  │ #g7h      │  │                 │                 │  │ rewrite   │  │
│  │ deploy    │  │                 │                 │  └───────────┘  │
│  │ [2]       │  │                 │                 │                 │
│  └───────────┘  │                 │                 │                 │
└─────────────────┴─────────────────┴─────────────────┴─────────────────┘
┌───────────────────────────────────────────────────────────────────────┐
│  🎯 FOCUS: ship v2 auth refactor                                      │
│  📝 NOTES (5)  •  jwt rotation...  •  pipe stderr...  •  meeting...  │
└───────────────────────────────────────────────────────────────────────┘
```

**Refinements from your feedback:**
- **Ideas column** — shows ideas as cards in the 4th column (same visual language as tasks)
- **Notes bar** — single-line scrolling/truncated list under focus (no separate panel needed)
- **No redundant DONE count** — the ✔ DONE column *is* the count
- **Column widths** — Tasks: ~35 chars each; Ideas: ~25 chars; Notes: remaining width
- **Empty states** — columns show dimmed "(empty)" when no items

**Interaction:** None — **static render only**. Just a visual dashboard to see at a glance. No keybindings, no navigation, no input loops. Future phases can add interactivity.

**Tech:** Pure ANSI box-drawing (`fmt.Printf` with calculated column widths). No deps. Works in any terminal ≥80 cols.

**Layout logic:**
```go
// Terminal width detection (default 120)
termWidth := getTermWidth() // fallback 120

// Column widths
taskColW := 35  // 3 task columns
ideaColW := termWidth - 3*taskColW - 6 // borders/spacers
// Notes bar uses full width below
```

---

## Concept 1b: Compact Kanban (Narrow Terminals)

For terminals <100 cols, stack vertically:

```
┌────────────────────────────────────────────────────────────┐
│  🎯 FOCUS: ship v2 auth refactor                           │
├────────────────────────────────────────────────────────────┤
│  📋 PENDING (2)                                            │
│  ┌──────────────────────────────────────────────────────┐  │
│  │ #a1b  fix login redirect loop              [•]       │  │
│  │ #g7h  deploy to staging                    [•]       │  │
│  └──────────────────────────────────────────────────────┘  │
├────────────────────────────────────────────────────────────┤
│  ✔ DONE (1)                                                │
│  ┌──────────────────────────────────────────────────────┐  │
│  │ #c3d  wrote integration tests              [✔]       │  │
│  └──────────────────────────────────────────────────────┘  │
├────────────────────────────────────────────────────────────┤
│  ✖ CANCELED (1)                                            │
│  ┌──────────────────────────────────────────────────────┐  │
│  │ #e5f  refactor auth                          [✖]       │  │
│  └──────────────────────────────────────────────────────┘  │
├────────────────────────────────────────────────────────────┤
│  💡 IDEAS (2)  📝 NOTES (5)                                │
│  • rate limiter middleware   • jwt rewrite                 │
│  • jwt rotation...  • pipe stderr...  • meeting notes...   │
└────────────────────────────────────────────────────────────┘
```

Auto-detects width; falls back to this layout if `termWidth < 100`.

---

## Concept 2: Chronological Feed (Timeline View)

**Visual:** Reverse-chronological unified stream — all buckets interleaved by time.

```
▼ ▼ ▼  CRUMB FEED  ▼ ▼ ▼

🎯 14:32  Focus set: "ship v2 auth refactor"
📝 14:30  Note: "jwt refresh rotation = store hash, rotate on use"
💡 14:28  Idea: "add rate limiter middleware before auth"
✅ 14:25  Task done: "write integration tests for login" [#a1b]
📝 14:20  Note: "pipe stderr: cmd 2> err.log"
💡 14:15  Idea: "rewrite auth with jwt"
📋 14:10  Task added: "fix login redirect loop" [#c3d]
🎯 13:45  Focus set: "debug payment webhook"
...

[←/→ filter: all | tasks | notes | ideas | focus]
[↑/↓ scroll]  [Enter: act]  [q: quit]
```

**Interaction:** Filter chips at bottom; scroll with `j`/`k`; `Enter` on task → toggle status; `/` search.

**Data needed:** Add `CreatedAt time.Time` to Task, Note, Idea (schema migration).

---

## Concept 3: Focus-Centered Radial (Priority Dashboard)

**Visual:** Focus in center; tasks/notes/ideas as orbiting rings by urgency.

```
                    ┌─────────────────────┐
                    │  🎯 SHIP V2 AUTH    │
                    │     REFACTOR        │
                    └──────────┬──────────┘
                               │
        ┌──────────────────────┼──────────────────────┐
        │                      │                      │
   ┌────▼────┐           ┌────▼────┐           ┌────▼────┐
   │ URGENT  │           │  SOON   │           │ LATER   │
   │ 📋 #a1b │           │ 📋 #c3d │           │ 💡 rate │
   │ fix bug │           │ deploy  │           │ limiter │
   │ 📝 jwt  │           │ 💡 jwt  │           │ 📝 pipe │
   │ rotation│           │ rewrite │           │ stderr  │
   └─────────┘           └─────────┘           └─────────┘

   [1] [2] [3] quick-jump    [Tab] cycle rings    [Space] act on selected
```

**Logic:** Urgency = (status == pending) + (mentioned in focus) + recency.

**Interaction:** Number keys 1-9 jump to item; `Tab` cycles rings; `Space` on task → status menu.

---

## Concept 4: Compact "Today" View (Time-Aware)

**Visual:** Grouped by date; shows only today + yesterday + this week.

```
╭────────────────────────────────────────────────────╮
│  📅 TODAY (Aug 7)              🎯 Focus: auth     │
├────────────────────────────────────────────────────┤
│  📋 Tasks:  3 pending  │ 1 done  │ 0 canceled      │
│    ┌────────────────────────────────────────────┐  │
│    │ [#a1b] fix login redirect loop        [•]  │  │
│    │ [#c3d] deploy to staging                [•]  │
│    │ [#e5f] write tests          ✔ done       │  │
│    └────────────────────────────────────────────┘  │
│  📝 Notes (2):  "jwt rotation...", "pipe stderr"  │
│  💡 Ideas (1):  "rate limiter middleware"         │
├────────────────────────────────────────────────────┤
│  📅 YESTERDAY (Aug 6)                               │
│    ✅ Completed: "auth redesign", "ci fix"         │
│    📝 Notes: "meeting notes: sprint planning"      │
├────────────────────────────────────────────────────┤
│  📅 THIS WEEK                                       │
│    📋 12 tasks  •  8 done  •  4 pending            │
│    💡 5 ideas  •  📝 18 notes                       │
╰────────────────────────────────────────────────────╯
```

**Interaction:** `h`/`l` expand/collapse days; `t` jump to today; `/` search all time.

**Data needed:** `CreatedAt` on all items (migration).

---

## Concept 5: Interactive TUI (Bubble Tea / Charm)

**Visual:** Full-screen navigable app with panels, forms, live updates.

```
┌─ crumb ──────────────────────────────────────────┐
│  🎯 Focus: ship v2 auth refactor    [e]dit [c]lear│
├──────────────────────────────────────────────────┤
│  📋 TASKS (3 pending, 1 done)                    │
│  ▸ [#a1b] fix login redirect loop      [•]       │
│    [#c3d] deploy to staging            [•]       │
│    [#e5f] write tests                  [✔]       │
│    [#g7h] refactor auth                [✖]       │
│  [a]dd  [d]one  [x]cancel  [f]ail  [/]search     │
├──────────────────────────────────────────────────┤
│  📝 NOTES (5)          💡 IDEAS (3)              │
│  ▸ jwt rotation...      ▸ rate limiter           │
│    pipe stderr...        jwt rewrite             │
│    meeting notes...      dashboard redesign      │
│  [n]ew                 [i]dea                    │
└──────────────────────────────────────────────────┘
  Status: Ready  |  ? help  |  q quit
```

**Tech:** `github.com/charmbracelet/bubbletea` + `lipgloss` + `bubbles`.
**Pros:** Rich interaction, forms, live resize, mouse support.
**Cons:** Adds ~2MB binary, external deps, more complex.

---

## Concept 6: Stats + Sparklines (Metric Dashboard)

**Visual:** Numbers-first with ASCII sparklines for trends.

```
╔════════════════════════════════════════════════════╗
║  📊 CRUMB STATS — Last 7 Days                      ║
╠════════════════════════════════════════════════════╣
║  🎯 Focus Changes: 12  ▁▂▃▅▂▁▂▃▅▂▁▂▃▅▂▁▂▃▅▂▁▂▃▅▂  ║
║  📋 Tasks Created: 34  ▂▃▅▂▁▂▃▅▂▁▂▃▅▂▁▂▃▅▂▁▂▃▅▂▁  ║
║  ✔ Tasks Done:     28  ▁▂▃▅▂▁▂▃▅▂▁▂▃▅▂▁▂▃▅▂▁▂▃▅▂  ║
║  📝 Notes:         45  ▃▅▂▁▂▃▅▂▁▂▃▅▂▁▂▃▅▂▁▂▃▅▂▁▂  ║
║  💡 Ideas:         12  ▁▁▂▂▃▃▅▅▂▂▁▁▂▂▃▃▅▅▂▂▁▁▂▂▃  ║
╠════════════════════════════════════════════════════╣
║  🎯 Current Focus: "ship v2 auth refactor"         ║
║  📋 Active Tasks:  3  (2 pending, 1 blocked)       ║
║  ⚡ Velocity:      4.2 tasks/day (↑ 12%)           ║
╠════════════════════════════════════════════════════╣
║  🔥 HOT: "auth" (8 mentions), "deploy" (5)         ║
║  📌 STALE: #g7h "refactor auth" — 5 days pending   ║
╚════════════════════════════════════════════════════╝
```

**Interaction:** `1`/`7`/`30` day windows; `Enter` on stat → filtered list.

**Data needed:** `CreatedAt` + computed aggregates (can compute on read).

---

## Concept 7: Minimal "Command Palette" Dashboard

**Visual:** Single-line prompt with fuzzy results — like `fzf` or VS Code palette.

```
> crumb dashboard
═══════════════════════════════════════════════════
🎯 Focus: ship v2 auth refactor
📋 3 active tasks  •  📝 5 notes  •  💡 3 ideas

Type to filter, ↑/↓ navigate, Enter to act:

> auth
  📋 [#a1b] fix login redirect loop        [• pending]
  📋 [#g7h] refactor auth                  [✖ canceled]
  💡 rewrite auth with jwt
  📝 jwt refresh rotation pattern...
  🎯 Focus: ship v2 auth refactor

> deploy
  📋 [#c3d] deploy to staging              [• pending]
  📝 deploy checklist: staging -> prod
```

**Interaction:** Starts in filter mode; type → fuzzy match across all buckets; `Enter` → action menu (done/cancel/fail/edit).

**Tech:** Pure Go — simple fuzzy matcher (levenshtein or substring), ANSI for highlighting.

---

## Decision Matrix

| Concept | Complexity | Deps | Data Migration | Best For |
|---------|------------|------|----------------|----------|
| 1. Kanban | Low | None | No | Visual task triage |
| 2. Feed | Medium | None | Yes (CreatedAt) | Audit trail, context restore |
| 3. Radial | Medium | None | No | Focus-driven workflow |
| 4. Today | Medium | None | Yes (CreatedAt) | Daily standup, journaling |
| 5. TUI | High | Bubble Tea | Optional | Power users, full app feel |
| 6. Stats | Low | None | Optional (CreatedAt) | Metrics, motivation |
| 7. Palette | Low | None | No | Speed, keyboard-first |

---

## Recommended Path
 
**Phase 1 (now, no deps, no migration):** Implement **Concept 8 (ASCII Table)** as the default dashboard:
- `crumb` → ASCII Table (dynamic columns, width-safe)
- Falls back to stacked single-column layout when width < 60
- `--simple` / `-s` flag → legacy simple dashboard (backward compat)
 
**Phase 2 (later):** Add `CreatedAt` to all types → enables Feed, Today, Stats.
 
**Phase 3 (optional):** Bubble Tea TUI for power users (Concept 5).

---

## Open Questions

1. **Default mode:** Kanban or Palette? (Recommend: Kanban for visual, Palette via flag)
2. **Color theme:** Keep current ANSI or add theme config (dark/light/high-contrast)?
3. **Keybindings:** Vim-style (`j`/`k`/`h`/`l`) or arrow keys? (Recommend: both)
4. **Persist view preference:** Save last-used mode to config file?

---

## Next Steps

1. **Add terminal width detection** in `helpers/output.go` (uses `os.Stdout` + `golang.org/x/term`)
2. **Create `renderKanban(data *store.CrumbData, width int) string`** in `cmd/root.go` (or new `cmd/dashboard.go`)
3. **Implement 4-column layout** with calculated column widths, box-drawing chars
4. **Add Compact layout** for width < 100
5. **Update `RootCmd.RunE`** to call `renderKanban` instead of current dashboard
6. **Add `--simple` flag** to preserve old behavior
7. **Test:** `go run .` (wide), resize terminal narrow, `go run . --simple`

---

## Implementation Details (for the agent)

### New helper: `helpers/output.go`
```go
func GetTermWidth() int {
    w, _, err := term.GetSize(int(os.Stdout.Fd()))
    if err != nil || w < 40 {
        return 120 // default
    }
    return w
}

func Truncate(s string, max int) string {
    if len(s) <= max { return s }
    return s[:max-1] + "…"
}
```

### Kanban renderer (pseudo)
```go
func renderKanban(data *store.CrumbData) string {
    w := GetTermWidth()
    if w < 100 { return renderCompactKanban(data, w) }
    
    taskColW := 35
    ideaColW := w - 3*taskColW - 6 // borders
    
    // Build each column as []string lines
    pendingLines := buildTaskColumn(data.Tasks, "pending", taskColW)
    doneLines    := buildTaskColumn(data.Tasks, "done", taskColW)
    cancelLines  := buildTaskColumn(data.Tasks, "canceled", taskColW)
    ideaLines    := buildIdeaColumn(data.Ideas, ideaColW)
    
    // Zip columns horizontally (pad shorter ones)
    // Render header row, then zipped body, then focus+notes footer
}
```

### Column builder
```go
func buildTaskColumn(tasks []store.Task, status string, width int) []string {
    var lines []string
    header := fmt.Sprintf(" %s %s ", icon(status), strings.ToUpper(status))
    lines = append(lines, boxTop(width), boxHeader(header, width))
    for _, t := range tasks {
        if t.Status == status {
            lines = append(lines, boxRow(fmt.Sprintf("#%s %s", t.ID, Truncate(t.Text, width-8)), width))
        }
    }
    if len(lines) == 2 { lines = append(lines, boxRow("(empty)", width)) }
    lines = append(lines, boxBottom(width))
    return lines
}
```

### Box drawing primitives
```go
const (
    TL = "┌"; TR = "┐"; BL = "└"; BR = "┘"
    H  = "─"; V  = "│"
    T  = "┬"; B  = "┴"; L  = "├"; R  = "┤"; X = "┼"
)

func boxTop(w int) string      { return TL + strings.Repeat(H, w-2) + TR }
func boxBottom(w int) string   { return BL + strings.Repeat(H, w-2) + BR }
func boxHeader(s string, w int) string { return V + " " + padRight(s, w-4) + " " + V }
func boxRow(s string, w int) string    { return V + " " + padRight(s, w-4) + " " + V }
```

### Compact renderer
Single-column stack with full-width boxes, section headers, inline counts.
 
---
 
## Concept 8: ASCII Table Dashboard (Width-Safe)
 
**Visual:** Clean table using only `|` `-` `+` characters. No box-drawing Unicode. Works at any terminal width by dynamically sizing columns.
 
```
+----+---------------------------+--------+----------------------+
| ID | Task                      | Status | Focus                |
+----+---------------------------+--------+----------------------+
| a1b| fix login redirect loop   | [•]    | ship v2 auth refactor|
| c3d| deploy to staging         | [✔]    |                      |
| e5f| refactor auth             | [✖]    |                      |
+----+---------------------------+--------+----------------------+
|    | IDEAS (3)                 |        | NOTES (5)            |
|    | • rate limiter middleware |        | • jwt rotation...    |
|    | • jwt rewrite             |        | • pipe stderr...     |
|    | • dashboard redesign      |        | • meeting notes...   |
+----+---------------------------+--------+----------------------+
```
 
**How it handles narrow widths:**
- Column widths computed from terminal width (`COLUMNS` env / fallback 120)
- Text truncates with `…` when exceeding column width
- At very narrow widths (<60): switches to stacked single-column layout
- No alignment dependencies — each row is independent
 
**Layout logic:**
```
termW := GetTermWidth()           // e.g. 120
focusW := min(40, termW/3)        // focus column
taskW  := termW - focusW - 4      // ID(4) + Task + Status(6) + separators
```
 
**Sections:**
1. **Focus + Tasks table** — ID | Task text | Status badge | Focus (repeated only on first row)
2. **Ideas | Notes side-by-side** — two sub-columns, truncated
3. **Empty states** — `(none)` / `(empty)` in gray
 
**No interaction** — static render only. `--simple` flag preserves legacy view.
 
**Tech:** Pure Go, zero deps. Uses existing `helpers.GetTermWidth()`, `helpers.Truncate()`, `helpers.FormatStatus()`.
 
**Pros:** 
- Works at ANY width (40–200+ cols)
- No Unicode box chars = no font/terminal rendering issues
- Simple to maintain and extend
- Fast render (single pass)
 
**Cons:** Less "visual" than Kanban, but bulletproof.