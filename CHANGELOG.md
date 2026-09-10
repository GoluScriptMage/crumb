# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v1.3.0] — 2026-09-10

### Added
- `crumb timer <minutes> [task]` — live running aesthetic terminal desk clock with 20 centered random ASCII mascots (Sia, Cat, Fox, Bunny, Bear, Frog, Owl, Penguin, Ghost, Slime, Robot, Knight, Wizard, Ninja, Samurai, Pirate, Astronaut, Demon, Reaper, Dragon), clean progress bar, and Ctrl+C background detach.
- `crumb del <id>` & `crumb del all` — intuitive, streamlined task deletion.
- `crumb fail <id>` — mark a task as failed (mirrors `done`).
- Index-based deletion for notes (`crumb note del <num>`) and ideas (`crumb idea del <num>`).
- Subtle monochrome tactile CLI feedback (dim prefixes, bold white typography, zero harsh RGB clown colors).
- `crumb version` & `crumb --version` — print current crumb version.

### Changed
- Converted all read operations (dashboard, task list, note list, idea list) to `store.ReadData()` (zero disk IO on reads).
- Deprecated and removed bloated `next`, `cancel`, `edit`, and `clear clear` commands.
- Task IDs are now 4-char hex (16^4 = 65536 space).
- `WriteData` is truly atomic: writes to a temp file then `os.Rename` over target.

## [v1.2.0] — 2026-08-07

### Added
- Full black-box test suite under `tests/`.
- `store.SetDbPathOverride` for temp DB isolation in tests.
- `crumb delete <id>` command to remove a task permanently.
- Exported `TaskCmd`, `DoneCmd`, `CancelCmd`, `DeleteCmd`, `ClearCmd`, `FormatStatus`, `GenerateShortId`.

## [v1.1.0] — 2026-08-07

### Added
- `crumb task clear` double-keyword safety.
- `crumb clear clear` dedicated clear command.
- Batch `add` support (multiple args to `crumb task`).

### Changed
- Dashboard layout redesigned (compact list, recent notes/ideas).
- Status badges: `done`, `canceled`, `failed`, `pending`.

## [v1.0.0] — 2026-08-07

### Added
- Initial release: `crumb note`, `crumb idea`, `crumb task`, `crumb done`, `crumb cancel`, `crumb next`.
- Atomic `ReadData`/`WriteData` JSON store in `~/.config/crumb/data.json`.
- Dashboard (`crumb` with no args).
- `helpers` color output layer.
