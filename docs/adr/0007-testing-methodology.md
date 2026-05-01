# 0007 — Testing Methodology

## Status

Accepted

## Context and Problem Statement

The project needs a consistent testing approach that catches real bugs without becoming a maintenance burden or slowing down the development loop.

## Decision

Test-driven development (TDD) as internal discipline: write a failing test first, then implement, then refactor.

Both unit and integration tests are written:
- **Unit tests** cover pure logic — path resolution, status detection, layer precedence, config parsing.
- **Integration tests** cover full command flows against a real temporary directory (real filesystem, real symlinks).

Tests live alongside the code they test (`foo_test.go` next to `foo.go`), following standard Go convention.

Integration tests are gated behind a `//go:build integration` tag and do not run by default. The justfile exposes:
- `just test` — unit tests only (fast)
- `just test-integration` — all tests including integration

Shared test fixtures and helpers (temp dotfiles dirs, fake `$HOME`, etc.) live in `internal/testutil/`.

## Consequences

- Fast feedback loop by default — `just test` runs in milliseconds.
- Integration tests exercise real filesystem behavior, catching issues that unit tests with mocks would miss.
- `internal/testutil` prevents fixture duplication as new domains are added.
- TDD as internal discipline means no ADR or user approval needed per test — it is part of the implementation flow.
