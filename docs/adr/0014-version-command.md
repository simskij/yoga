# ADR-0014: Version command

Status: Accepted

## Context

No way to check which version of yo is installed.

## Decision

Add `yo version` that prints the current version.

The version string is embedded at build time via `-ldflags "-X github.com/simskij/yo/cmd/yo.version=<tag>"`. For untagged builds (local `go build`, `go install` from HEAD) the default value `dev` is used.

The release workflow is updated to pass the tag name as the version.

## Consequences

- `yo version` prints the installed version or `dev` for local builds.
- No runtime reflection or `debug.ReadBuildInfo` — the value is always deterministic.
