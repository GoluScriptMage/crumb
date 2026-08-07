# crumb

> A tiny terminal-first scratchpad for developers who think in the CLI.

No GUI. No sync. No nonsense. Just `crumb note "thing"` and it's saved.

---

## Why?

You're in the flow. A thought hits. A task appears. An idea sparks.

```
crumb note "pipe stderr to file: cmd 2> err.log"
crumb idea "rewrite auth with jwt"
crumb task "fix login redirect loop"
crumb next "debug payment webhook"
```

Done. Back to work. Zero friction.

---

## Install

```bash
go install github.com/yourusername/crumb@latest
```

Or clone & build:

```bash
git clone https://github.com/yourusername/crumb
cd crumb
go build -o crumb .
sudo mv crumb /usr/local/bin/
```

---

## Commands

| Command | What it does |
|---------|--------------|
| `crumb note "text"` | Save a note (single arg) |
| `crumb note` | List notes (newest first) |
| `crumb idea "a" "b"` | Save multiple ideas at once |
| `crumb idea` | List ideas (newest first) |
| `crumb task "do thing"` | Add a task |
| `crumb task` | List tasks with status badges |
| `crumb task clear` | Nuke all tasks |
| `crumb done <id>` | Mark task ✔ done |
| `crumb cancel <id>` | Mark task ✖ canceled |
| `crumb fail <id>` | Mark task ✖ failed |
| `crumb next "focus"` | Set your current focus |
| `crumb next` | Show current focus |
| `crumb next clear` | Clear focus |
| `crumb` | Dashboard: focus + tasks + recent notes/ideas |

---

## How It Works (The Diagram)

```
┌─────────────────────────────────────────────────────────────┐
│                        YOU TYPE                               │
│  crumb note "golang slices are views, not copies"            │
└─────────────────────────┬─────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                      store.Update()                           │
│  ┌─────────────┐    ┌──────────────┐    ┌────────────────┐  │
│  │  ReadData() │───▶│  Your Func   │───▶│  WriteData()   │  │
│  │  (JSON)     │    │  (mutates)   │    │  (atomic)      │  │
│  └─────────────┘    └──────────────┘    └────────────────┘  │
└─────────────────────────┬─────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                    ~/.config/crumb/data.json                  │
│  {                                                            │
│    "next": "debug payment webhook",                          │
│    "tasks": [                                                 │
│      {"id": "a1b", "text": "fix login", "status": "done"}    │
│    ],                                                         │
│    "notes": ["golang slices are views..."],                  │
│    "ideas": ["rewrite auth with jwt"],                       │
│    "done": []                                                │
│  }                                                            │
└─────────────────────────────────────────────────────────────┘
```

**Data flow:** `ReadData` → your mutation function runs on a pointer → `WriteData` persists atomically. No manual file handling in commands.

---

## Dashboard (`crumb`)

```
=== Crumb Dashboard ===
  Focus: debug payment webhook
  Tasks: 3 active
  Recent notes:
    • golang slices are views, not copies
    • pipe stderr to file: cmd 2> err.log
  Recent ideas:
    • rewrite auth with jwt
    • add rate limiter to api
```

Instant context. No flags. No subcommands.

---

## Status Badges

Tasks show colored status inline:

```
📋 Active tasks:
--------------------------------------------------
  [1]  fix login redirect loop            [#a1b]  [✔ done]
  [2]  add rate limiter to api            [#c3d]  [• pending]
  [3]  rewrite auth with jwt              [#e5f]  [✖ canceled]
  [4]  debug payment webhook              [#g7h]  [✖ failed]
--------------------------------------------------
```

- `[✔ done]` — finished
- `[• pending]` — active
- `[✖ canceled]` — won't do
- `[✖ failed]` — tried, blocked

---

## Philosophy

- **Append-only by default** — capture fast, organize later
- **Local-first** — JSON in `~/.config/crumb/`, yours forever
- **No accounts, no cloud** — it's a scratchpad, not a SaaS
- **Terminal-native** — works over SSH, in tmux, in CI logs

---

## Hacking

```bash
# Add a command
# 1. Create cmd/yourcmd.go with RunE + store.Update
# 2. Register in init(): RootCmd.AddCommand(yourCmd)

# The pattern:
store.Update(func(data *store.CrumbData) error {
    data.YourField = "value"
    return nil
})
```

---

## License

MIT. Do whatever.

---

*Built with Go, Cobra, and zero dependencies beyond stdlib.*
*AUTHOR: Built by Golu*