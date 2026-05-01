# yo — Claude Code Instructions

## ADR-first workflow

Every change to the codebase must be preceded by an Architecture Decision Record.

1. Write the ADR in `docs/adr/` with status **Proposed**
2. Wait for explicit user approval
3. Update status to **Accepted**
4. Implement

No code is written before an ADR is approved. This applies to new features, refactors, and any non-trivial structural change.

## Testing (ADR-0007)

- Follow TDD: write a failing test first, then implement, then refactor
- Write both unit and integration tests
- Unit tests: pure logic, run always
- Integration tests: full flows against real filesystem, tagged `//go:build integration`
- Run `just test` and ensure all tests pass between each increment
- Shared fixtures live in `internal/testutil/`

## Project conventions

- New domains: `internal/<domain>/` + `cmd/yo/cmd_<domain>.go`
- Config: `~/.config/yo/config.yaml`, namespaced by domain
- Dotfiles: layered `global/` + `<hostname>/`, symlinked into `$HOME`
- Deferred ideas: `docs/backlog.md`
