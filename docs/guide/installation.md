# Installation

## Requirements

- **Go 1.23.5+** (for installation via `go install`)
- **Linux** with Steam
- **Baldur's Gate 3** installed via Steam

## Install via Go

```bash
go install github.com/GiGurra/bookeeper@latest
```

Make sure `$GOPATH/bin` (usually `~/go/bin`) is in your `PATH`.

## Verify Installation

```bash
bookeeper --help
```

You should see the help message with available commands.

## Shell Completion

Bookeeper supports shell completion for bash, zsh, fish, and powershell.

### Bash

```bash
# Add to ~/.bashrc
source <(bookeeper completion bash)
```

### Zsh

```bash
# Add to ~/.zshrc
source <(bookeeper completion zsh)
```

### Fish

```bash
bookeeper completion fish | source
# Or save to ~/.config/fish/completions/bookeeper.fish
bookeeper completion fish > ~/.config/fish/completions/bookeeper.fish
```

## First Run

After installation, run:

```bash
bookeeper status
```

This will:

1. Show you the detected paths
2. Create the bookeeper directory (`~/.local/share/bookeeper/`)
3. Verify BG3 installation is found

If paths are wrong, see [Configuration](../reference/configuration.md).

## Directory Structure

Bookeeper creates:

```
~/.local/share/bookeeper/
├── downloaded_mods/    # Your extracted mods
└── profiles/           # Saved mod profiles
```

## Next Steps

- [Basic Usage](usage.md) - Learn common workflows
- [Configuration](../reference/configuration.md) - Customize paths if needed
