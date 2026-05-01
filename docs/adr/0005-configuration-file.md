# 0005 — Configuration File Location and Format

## Status

Accepted

## Context and Problem Statement

The tool needs persistent configuration (at minimum, the path to the dotfiles repo) that survives across invocations and is easy to inspect and edit manually.

## Considered Options

- `~/.yorc` (flat dotfile)
- `~/.config/yo/config.yaml` (XDG)
- Environment variables only

## Decision

`~/.config/yo/config.yaml`, following the XDG Base Directory specification. YAML format. Top-level keys are namespaced by domain to match the CLI structure.

```yaml
dots:
  path: ~/.dotfiles
```

Created by `yo init`, which prompts for `dots.path` and writes the file with defaults.

## Consequences

- Follows XDG convention — keeps `$HOME` clean.
- YAML is human-readable and easy to edit manually.
- Namespaced keys (`dots.path`) scale cleanly as new domains are added.
- Environment variables alone would make the tool harder to use across shells and tools.
