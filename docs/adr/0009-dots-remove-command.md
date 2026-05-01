# 0009 — `yo dots remove` Command

## Status

Accepted

## Context and Problem Statement

The inverse of `yo dots add` is needed: stop managing a file, remove it from the dotfiles repo, and restore it as a regular file in place.

## Decision

Add a `yo dots remove <file>` command that:

1. Resolves which layer(s) the file exists in (`global/`, `<hostname>/`, or both)
2. If the file exists in only one layer, removes it from that layer
3. If the file exists in both layers, prompts the user to choose which layer to remove from
4. Removes the repo copy from the chosen layer
5. Replaces the symlink at the original location with the file contents

Errors out if:
- The file at the given path is not a symlink managed by `yo` (i.e. does not point into the dotfiles repo)

## Consequences

- Unambiguous cases (file in one layer only) require no user input.
- Ambiguous cases (file in both layers) surface an explicit choice rather than defaulting silently, avoiding accidental loss of a layer.
- No `--layer` flag needed — the prompt handles the ambiguity at the moment it matters.
