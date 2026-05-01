# 0000 — Use ADRs for Decision Recording

## Status

Accepted

## Context and Problem Statement

As the project grows across multiple domains (dotfiles, app installation, migrations, etc.), architectural and design decisions need to be traceable. Without a record, the reasoning behind choices is lost, making future changes harder to evaluate.

## Decision

We use Markdown Architectural Decision Records (MADR) stored in `docs/adr/`. Every change to the codebase must be preceded by an ADR with status Accepted.

The workflow is:
1. An ADR is written with status Proposed.
2. The author reviews and approves it verbally.
3. Status is updated to Accepted.
4. Implementation begins.

## Consequences

- All significant decisions are documented with their rationale.
- No code is written without prior design sign-off.
- The backlog for deferred decisions lives in `docs/backlog.md`.
