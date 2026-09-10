# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v1.3.0] — 2026-09-08

### Added
- `crumb fail <id>` — mark a task as failed (mirrors `done`/`cancel`).
- `crumb version` — print the crumb version.
- `crumb --version` — Cobra built-in version flag.
- `crumb edit <id> <text...>` — edit a task's text by ID.
- `CHANGELOG.md` — version history document.
- README sync for new commands (`fail`, `version`, `edit`, `delete`, `clear clear`).

### Changed
- Task IDs are now 4-char hex (16^4 = 65536 space), up from 3-char.
- `WriteData` is truly atomic: writes to a temp file then `os.Rename` over the target.

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
