# 0004 — Dotfiles Symlink Strategy

## Status

Accepted

## Context and Problem Statement

Files from the dotfiles repo need to be installed into `$HOME`. Changes to dotfiles should be immediately reflected without a sync step. Existing files must not be silently destroyed.

## Considered Options

- Copy files into place
- Symlink files into place, skip if target exists
- Symlink files into place, back up existing files

## Decision

Symlink files into place. If a target already exists (regular file or a symlink pointing elsewhere), rename it to `<name>.bak` before creating the symlink. `apply` is idempotent — re-running is a no-op for correctly linked files.

## Consequences

- Edits to files under `$HOME` are immediately reflected in the dotfiles repo with no sync needed.
- No existing config is ever silently lost — it is preserved as a `.bak` file.
- Copies would require an explicit sync step back to the repo; rejected for this reason.
