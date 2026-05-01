# 0011 — Self-Upgrade and Remote Install

## Status

Accepted

> Note: implemented after ADR-0012. `--dots-path` references in the original text
> map to `--path` in the implementation. `yo dots clone` clones into the
> `dots/` subdirectory of `yo.path`, not directly into `dots.path`.

## Context and Problem Statement

Installing and updating yo requires either building from source or manually downloading a binary. A self-upgrade command and a curl-pipe installer lower the barrier to entry and make bootstrapping a new machine frictionless.

## Decision

### `yo upgrade`

Fetches the latest release from GitHub (`simskij/yo`) and replaces the running binary in-place using `os.Executable()`.

- Detects current OS and architecture to download the correct asset
- Downloads to a temp file, verifies it is executable, then atomically replaces the current binary
- No installation method tracking — always replaces the binary at `os.Executable()`

### `yo dots clone <repo>`

Clones an existing dotfiles repo into `dots.path` and runs `yo dots apply`.

- Errors if `dots.path` already exists and is non-empty
- Accepts any valid git URL

### `install.sh`

A shell script hosted at `https://raw.githubusercontent.com/simskij/yo/main/install.sh`.

**Basic usage** (interactive):
```sh
curl -sL https://raw.githubusercontent.com/simskij/yo/main/install.sh | sh
```
Installs the binary to `~/.local/bin/yo`, then prints instructions to run `yo init`.

**One-shot usage** (non-interactive, full bootstrap):
```sh
curl -sL https://raw.githubusercontent.com/simskij/yo/main/install.sh | sh -s -- \
  --repo git@github.com:you/dotfiles.git \
  --dots-path ~/.dotfiles
```
Runs: install binary → `yo init --dots-path <path>` (non-interactive) → `yo dots clone <repo>` → `yo dots apply`.

**Flags**:
- `--repo <url>` — dotfiles repo to clone; triggers one-shot mode
- `--dots-path <path>` — override dotfiles path (default: `~/.dotfiles`)
- `--bin-dir <path>` — override binary install directory (default: `~/.local/bin`)

### `yo init` non-interactive flag

To support one-shot install, `yo init` gains a `--dots-path` flag that skips the interactive prompt and writes config directly.

## Consequences

- New machines can be fully bootstrapped with a single curl command.
- `yo upgrade` keeps the tool current without needing a package manager.
- `yo dots clone` completes the symmetry with `yo dots init` for existing repos.
- `yo init --dots-path` enables scripted/non-interactive config creation.
