# 0012 — Unified `.yofiles` Repository Structure

## Status

Accepted

## Context and Problem Statement

The dotfiles repo was designed as a standalone directory, with `dots.path` in config pointing directly to it. As yo grows to cover more domains (app installation, migrations, etc.), each domain would need its own path config. This is fragmented and doesn't reflect that all managed state belongs to a single coherent repo.

## Decision

Replace the per-domain `dots.path` config with a single top-level `yo.path` pointing to a unified `~/.yofiles` repository. Each domain owns a subdirectory within it.

### Repository structure

```
~/.yofiles/
  dots/
    global/
      .zshrc
      .config/nvim/init.lua
    work-laptop/
      .zshrc
  # future domains sit alongside dots/
```

### Config

```yaml
yo:
  path: ~/.yofiles
```

### Responsibility split

- `yo init` — creates `~/.yofiles`, runs `git init`, writes config
- `yo dots init` — scaffolds `dots/global/` and `dots/<hostname>/` within `~/.yofiles`

The internal path for dots operations is derived as `<yo.path>/dots/` — no separate config key needed.

### Migration

ADR-0003 and ADR-0005 are superseded by this ADR. Existing setups using a standalone dotfiles repo can be migrated by moving its contents into `~/.yofiles/dots/`.

## Consequences

- Single source of truth for all yo-managed state.
- New domains require no new config keys — they claim a subdirectory under `yo.path`.
- `yo init` becomes the single bootstrap step for the entire tool, not just config.
- The `dots.path` config key is removed; `yo.path` replaces it.
