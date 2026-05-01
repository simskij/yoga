# yo

Personal kitchen-sink CLI. Currently does one thing: manages dotfiles across machines with optional encryption.

## Install

```sh
go install github.com/simskij/yo/cmd/yo@latest
```

Or build from source:

```sh
just build   # output: ./bin/yo
just install # installs to $GOPATH/bin
```

## Getting started

```sh
yo init
```

Creates `~/.config/yo/config.yaml` and initialises a git repo at `~/.yofiles`. That repo is where everything lives — dots go under `~/.yofiles/dots/`, and future domains will get their own subdirectories.

Then scaffold the dotfiles structure:

```sh
yo dots init
```

Creates `~/.yofiles/dots/global/` and `~/.yofiles/dots/<hostname>/`. Push `~/.yofiles` to a remote and you have a portable setup.

## How dotfiles work

Files live in two layers:

- `global/` — applied to every machine
- `<hostname>/` — machine-specific, takes precedence over global

Each layer mirrors `$HOME`. A file at `global/.config/git/config` gets symlinked to `~/.config/git/config`. Machine-specific files shadow global ones at the same path.

**Add a file:**

```sh
yo dots add ~/.zshrc                     # goes into global by default
yo dots add --machine ~/.zshrc           # goes into <hostname> layer
```

The file is moved into the repo and a symlink is left in its place. If something already exists at the destination, it's backed up as `.bak`.

**Apply on a new machine:**

```sh
yo dots apply
```

Symlinks everything into `$HOME`. Safe to re-run — already-linked files are skipped.

You can apply a subset by path:

```sh
yo dots apply .config/git
```

**Check status:**

```sh
yo dots status
```

```
✓ /home/you/.zshrc
✓ /home/you/.config/git/config
~ /home/you/.ssh/config      ← exists but not linked (run apply)
✗ /home/you/.config/foo/bar  ← missing entirely
```

**See what's drifted:**

```sh
yo dots diff
```

Shows a unified diff for any file that exists at the destination but doesn't match the repo copy.

**Stop managing a file:**

```sh
yo dots remove ~/.zshrc
```

Copies the file back in place and removes it from the repo. If the file exists in multiple layers, yo prompts you to choose which one to remove from.

## Encryption

For secrets — SSH configs, tokens, API keys — yo uses [age](https://age-encryption.org) to encrypt files before storing them in the repo. Encrypted files are copied into place on apply rather than symlinked.

```sh
yo dots add --encrypt ~/.ssh/config
```

The file is encrypted and stored as `~/.yofiles/dots/global/.ssh/config.age`. The original is removed. On `yo dots apply`, it's decrypted and written to `~/.ssh/config`.

To edit an encrypted file without manually handling temp files:

```sh
yo dots edit ~/.ssh/config
```

Decrypts to a temp file, opens `$EDITOR`, then re-encrypts. The plaintext never persists to disk after the editor closes.

**Key discovery:** yo looks for identities in this order:

1. Path set in config (`dots.encryption.identity`)
2. `~/.ssh/id_ed25519`
3. `~/.ssh/id_rsa`

Native age keys (`AGE-SECRET-KEY-...`) and SSH keys (ed25519, RSA) are both supported.

## Configuration

`~/.config/yo/config.yaml`:

```yaml
yo:
  path: ~/.yofiles          # where the repo lives

dots:
  encryption:
    identity: ~/.age/key    # optional; falls back to ~/.ssh/id_ed25519
```

## Targeting a specific machine

Most `yo dots` commands accept `--machine` to override the detected hostname:

```sh
yo dots status --machine work-laptop
yo dots apply --machine work-laptop
yo dots add --machine work-laptop ~/.zshrc
```

Without the flag, yo uses `os.Hostname()`.
