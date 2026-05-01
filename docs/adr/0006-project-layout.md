# 0006 — Go Project Layout

## Status

Accepted

## Context and Problem Statement

The project will grow to cover multiple domains beyond dotfiles. The layout needs to accommodate that without requiring restructuring later.

## Decision

Standard Go project layout:

```
cmd/yo/         — main entrypoint and command wiring
internal/       — all domain logic, not importable externally
  config/       — config load/save
  dots/         — dotfiles logic
  ui/           — shared Bubbletea/Lipgloss primitives
docs/
  adr/          — architecture decision records
  backlog.md    — deferred feature ideas
justfile        — build, install, lint, tidy
```

New domains are added as packages under `internal/` with a corresponding command file in `cmd/yo/`.

## Consequences

- Clear separation between CLI wiring and domain logic.
- `internal/` enforces that domain packages are not accidentally imported by external tools.
- Scales cleanly as new domains (app installation, migrations, etc.) are added.
