# 0003 — Dotfiles Repository Structure

## Status

Accepted

## Context and Problem Statement

Dotfiles need to work across multiple machines with shared config and machine-specific overrides, without duplicating files or requiring complex tooling.

## Considered Options

- Flat (all files in one dir)
- Grouped by topic (git/, zsh/, nvim/)
- Layered by machine with topic grouping
- Layered by machine, flat mini-$HOME per layer

## Decision

Two layers — `global/` and `<hostname>/` — each a mini-`$HOME`. Files are placed at their natural path relative to `$HOME` with no prefix or metadata files. Machine-specific files take full precedence over global files on a per-file basis.

```
dotfiles/
  global/
    .zshrc
    .config/nvim/init.lua
  work-laptop/
    .zshrc          ← overrides global
```

## Consequences

- No magic naming conventions or metadata files to learn.
- Selective apply is possible by passing a subpath (e.g. `yo dots apply .config/nvim`).
- Full-file precedence avoids fragile line-level merging; merging can be added later via a `--merge` flag if needed.
- The dotfiles repo is a git repo, initialized by `yo dots init`.
