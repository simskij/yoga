# ADR-0015: Shell completions

Status: Accepted

## Context

No tab completion makes the CLI slower to use. Bash and zsh are the common shells; nushell because it's interesting.

## Decision

Add `yo completion <shell>` with subcommands for `bash`, `zsh`, and `nushell`.

- **bash / zsh**: use Cobra's built-in generators (`GenBashCompletion`, `GenZshCompletion`).
- **nushell**: Cobra has no built-in support. A custom generator in `internal/completion/` walks the Cobra command tree and produces a nushell `module completions { }` block with `export extern` definitions for every command and its flags.

Cobra's auto-generated `completion` command is disabled (`DisableDefaultCmd = true`) so our command is the single entry point and nushell sits alongside the other shells cleanly.

## Consequences

- `yo completion bash | sudo tee /etc/bash_completion.d/yo`
- `yo completion zsh > ~/.zsh/completions/_yo`
- `yo completion nushell | save ~/.config/nushell/completions/yo.nu` then `use yo.nu *` in `config.nu`
- Fish and PowerShell are dropped — add later if there is demand.
