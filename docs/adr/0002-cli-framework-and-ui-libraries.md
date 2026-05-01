# 0002 — CLI Framework and UI Libraries

## Status

Accepted

## Context and Problem Statement

The tool needs command routing, subcommand namespacing, flag parsing, and polished terminal output including interactive prompts and colored status lines.

## Considered Options

- Cobra + plain fmt
- Cobra + Bubbletea + Lipgloss
- urfave/cli + Bubbletea + Lipgloss

## Decision

Cobra for command routing and flag parsing. Bubbletea for interactive prompts and all terminal output. Lipgloss for styling.

## Consequences

- Consistent, polished output across all commands.
- Cobra is the de facto standard for Go CLIs — well documented and familiar.
- Bubbletea's Elm-style model adds some structure overhead for simple output cases, but the visual consistency is worth it.
