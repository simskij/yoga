# 0008 — `yo dots add` Command

## Status

Accepted

## Context and Problem Statement

Adding a file to the dotfiles repo currently requires manually moving it, creating the correct directory structure, and symlinking it back. This is error-prone and tedious.

## Decision

Add a `yo dots add <file>` command that:

1. Resolves the file's path relative to `$HOME` to determine where it goes in the dotfiles repo
2. Moves the file into the `global/` layer of the dotfiles repo, mirroring its `$HOME`-relative path
3. Symlinks the original location back to the new repo location

Flags:
- `--machine` — place the file in the `<hostname>/` layer instead of `global/`

Errors out if:
- The path is a directory (files only)
- The file already exists in the target layer of the dotfiles repo
- The file is already a symlink managed by `yo` (i.e. already points into the dotfiles repo)

## Consequences

- The common case (add to global) requires no flags, matching the expected workflow.
- Machine-specific additions require an explicit `--machine` flag, making the choice deliberate.
- Erroring on conflicts avoids silent data loss.
- Directory support can be added later.
