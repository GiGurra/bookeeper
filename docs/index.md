# Bookeeper

**CLI mod manager for Baldur's Gate 3 on Linux/Steam**

[![CI Status](https://github.com/GiGurra/bookeeper/actions/workflows/ci.yml/badge.svg)](https://github.com/GiGurra/bookeeper/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/GiGurra/bookeeper)](https://goreportcard.com/report/github.com/GiGurra/bookeeper)

!!! warning "Prototype / Weekend Hack"
    This is an **experiment/prototype/weekend hack** I built when switching to Linux on my gaming machine. It works for my use case but has limitations. Use at your own risk!

## What is Bookeeper?

Bookeeper is a command-line tool for managing Baldur's Gate 3 mods on Linux with Steam. It handles:

- **Mod activation/deactivation** via symlinks
- **Mod profiles** for switching between different mod sets
- **BG3 Script Extender** installation
- **Multiple mod versions** side by side

## Why?

When I switched from Windows to Linux on my gaming PC, I needed a way to manage BG3 mods. Existing mod managers were Windows-focused, so I built this quick hack to solve my problem.

## Quick Example

```bash
# Install bookeeper
go install github.com/GiGurra/bookeeper@latest

# Check current status
bookeeper status

# Import a mod you downloaded
bookeeper mods make-available ~/Downloads/CoolMod.zip

# Activate it
bookeeper mods activate "CoolMod" "1.0"

# Save this setup as a profile
bookeeper profiles save my-playthrough

# Later, switch back to this mod set
bookeeper profiles load my-playthrough
```

## How It Works

```
You download mod.zip
        ↓
bookeeper mods make-available mod.zip
        ↓
Mod extracted to ~/.local/share/bookeeper/downloaded_mods/
        ↓
bookeeper mods activate "ModName" "version"
        ↓
Symlink created in BG3's Mods folder
        ↓
modsettings.lsx updated
        ↓
BG3 loads your mod!
```

## Limitations

Be aware of what this prototype **doesn't** do:

- No automatic mod downloading (you download zips manually)
- No mod dependency management
- No version conflict detection
- Limited error handling
- Only tested on Linux/Steam (not Windows/GOG)
- Load order only via activation order
- Requires mods to have valid `info.json`

## Getting Started

- [Installation](guide/installation.md) - How to install bookeeper
- [Basic Usage](guide/usage.md) - Common workflows
- [Profiles](guide/profiles.md) - Save and load mod sets
- [Troubleshooting](guide/troubleshooting.md) - Common issues

## Reference

- [Commands](reference/commands.md) - Full command reference
- [Configuration](reference/configuration.md) - Paths and settings
- [File Formats](reference/file-formats.md) - info.json and profile.json
