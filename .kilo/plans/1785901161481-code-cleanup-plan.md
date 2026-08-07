# Crumb Code Cleanup Plan

## Goal
Refactor existing code to be production-level: consistent error handling, colorized output via helper, switch-based control flow, readable comments, and DRY logic — **without changing user-facing behavior or adding features**. Keep the project recognizable.

---

## Current Issues

| Area | Problem |
|------|---------|
| Error handling | `helpers.HandleError` panics; commands print and return instead of propagating errors |
| Colors | Raw `fmt.Printf` everywhere; no consistent styling |
| Control flow | If/else chains; no switch statements for command scenarios |
| Comments | Almost none; logic not self-documenting |
| Duplication | Every command does `ReadData()` + `WriteData()` manually |
| `touch.go` | Empty file — remove |
| `note` vs `idea` | Inconsistent arg handling (`args[0]` vs `args...`) — keep as-is, just clean up |

---

## Plan (ordered tasks)

### 1. Create color/output helper
- New file: `helpers/output.go`
- Constants: `Reset`, `Bold`, `Green`, `Yellow`, `Red`, `Cyan`, `Gray`
- Functions: `Info()`, `Success()`, `Error()`, `Warn()`, `Dim()` — thin wrappers around `fmt.Print` with colors
- `FormatStatus(status string) string` — badge renderer (done/canceled/pending)

### 2. Fix error handling
- Remove `helpers/handleError.go` (panic-based)
- Change all commands from `Run` to `RunE` returning `error`
- Errors bubble up to Cobra; main prints via `output.Error()`

### 3. Add store transaction helper
- New: `store.Update(fn func(*CrumbData) error) error`
- Handles read → fn → write atomically
- Eliminates duplicate read/write boilerplate in every command

### 4. Refactor each command (`RunE` + switch + output helper)
**`note.go`**: Switch on `len(args)` → 0: list, 1: add, default: error. Use `output.Success/Info`.
**`idea.go`**: Switch on `len(args)` → 0: list, >0: add all. Clean reverse loop comment.
**`task.go`**: Switch on `args[0]` → "clear", "list", "add". `done` subcommand: switch on args.
**`next.go`**: Switch on args → empty=show, "clear"=clear, else=set.
**`root.go`**: Keep minimal — show next focus + task count + recent notes/ideas (read-only, no new flags).

### 5. Add short, readable comments
- Each command: one-line purpose + switch case comments
- Store functions: doc comments
- Helpers: doc comments

### 6. Remove `touch.go`
- Delete empty file

### 7. Verify
- `go build ./...`
- `go vet ./...`
- Manual smoke test: all existing commands work identically

---

## Out of Scope
- New features (search, timestamps on notes/ideas, delete for notes/ideas)
- Data model changes (Task ID format, Done struct, CreatedAt fields)
- Tests
- Config file / flags
- Dashboard as a new feature (root command just shows existing data cleanly)

---

## Rollout
Single PR. Zero behavior change. JSON format unchanged.