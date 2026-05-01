# 0001 — Go as Implementation Language

## Status

Accepted

## Context and Problem Statement

The tool needs to be fast, produce a single distributable binary, and be maintainable long-term as a personal project.

## Considered Options

- Go
- Rust
- Python

## Decision

Go. It compiles to a single static binary, has excellent cross-platform support, and strikes the right balance between low ceremony and type safety for a CLI tool of this scope.

## Consequences

- Single binary output — easy to install and distribute.
- No runtime dependency on Python, Node, etc.
- Rust would offer better performance and memory safety but at significantly higher development cost for diminishing returns at this scale.
