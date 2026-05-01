# 0010 — Dotfile Encryption with age

## Status

Accepted

## Context and Problem Statement

Some dotfiles contain sensitive data (API keys, tokens, private config) that should not be stored in plaintext in the dotfiles repo, even if the repo is private. Encryption at rest is needed.

## Decision

Add opt-in encryption using [age](https://age-encryption.org) for individual files, enabled via `--encrypt` on `yo dots add`.

### Encryption

When `yo dots add --encrypt <file>` is run, the file is encrypted with age before being written to the dotfiles repo. The repo never contains plaintext. Encrypted files are stored with an `.age` extension.

### Decryption

Encrypted files cannot be symlinked — symlinks point to files, not content. Instead, encrypted files are **copied into place** (decrypted) at `yo dots apply` time. Edits to the destination are not automatically reflected back to the repo. To update an encrypted file, run `yo dots add --encrypt` again.

### Key Resolution

Keys are resolved in this order, skipping any that are unavailable or fail:

1. `~/.ssh/id_ed25519`
2. `~/.ssh/id_rsa`
3. `dots.encryption.identity` in `~/.config/yo/config.yaml`

If no key is found, the operation errors with a clear message. If a key is found but fails (wrong key, unsupported type), it is skipped and the next is tried.

**Note**: SSH agent (`SSH_AUTH_SOCK`) support is not implemented. The `filippo.io/age` library requires direct access to private key material for decryption, which the agent never exposes. Agent support is tracked in `docs/backlog.md`.

### Config

```yaml
dots:
  encryption:
    identity: ~/.config/yo/age.key  # optional override
```

### UX

- `yo dots add --encrypt <file>` — encrypt and add
- `yo dots apply` — decrypts encrypted files into place (copy, not symlink)
- `yo dots edit <file>` — decrypt to a temp file, open `$EDITOR`, re-encrypt back to the repo on save
- `yo dots status` — encrypted files show a distinct status indicator

## Considered Options

- **Symlink to decrypted copy** — decrypt to `~/.local/share/yo/decrypted/`, symlink there. Keeps apply model consistent but adds indirection and leaves plaintext at a predictable path.
- **Copy decrypted into place** (chosen) — simpler, no persistent plaintext outside `$HOME`, consistent with how chezmoi handles encryption.

## Consequences

- Sensitive files can be safely committed to the dotfiles repo.
- Encrypted files do not benefit from the symlink live-edit model — re-add is required to update them.
- Key resolution is automatic and covers the common cases (SSH agent, existing SSH keys) without requiring explicit configuration.
- Multiple key candidates are tried silently; only a total failure is surfaced to the user.
