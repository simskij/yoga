# ADR-0013: Dry-run flag for setup and apply commands

Status: Accepted

## Context

Several yo commands make irreversible or hard-to-undo changes to the filesystem:

- `yo init` — writes `~/.config/yo/config.yaml` and initialises a git repo
- `yo dots init` — scaffolds the dots directory structure
- `yo dots apply` — symlinks or copies (encrypted) files into `$HOME`
- `yo dots clone` (ADR-0011, not yet implemented) — clones a remote repo and wires everything up

For all of these, a preview mode is useful: see what would happen before committing, especially on a new machine.

`yo dots add` and `yo dots remove` are also destructive but operate on a single named file where intent is already explicit. The value of dry-run there is lower and they are excluded from this ADR.

## Decision

Add `--dry-run` to `yo init`, `yo dots init`, and `yo dots apply`. When `yo dots clone` is implemented (ADR-0011), it must include `--dry-run` at that time.

When `--dry-run` is set on any of these commands:
- No filesystem changes are made
- Output mirrors what a real run would print, with a `(dry run)` suffix on each line so it is visually distinct
- `yo dots apply --dry-run` also lists files that would be skipped (already correctly linked), since the flag is typically used for a full overview

`--dry-run` is a local flag on each command, not a persistent flag on `dotsCmd` or `rootCmd`.

## Consequences

- Users can safely preview any setup step on a new machine before touching the filesystem
- No change to normal-run output format for any command
- `add` and `remove` remain unchanged
- `clone` is obligated to ship with `--dry-run` when implemented
