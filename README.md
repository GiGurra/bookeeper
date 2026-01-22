# Bookeeper

A CLI mod manager for Baldur's Gate 3 on Linux/Steam.

[![CI Status](https://github.com/GiGurra/bookeeper/actions/workflows/ci.yml/badge.svg)](https://github.com/GiGurra/bookeeper/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/GiGurra/bookeeper)](https://goreportcard.com/report/github.com/GiGurra/bookeeper)
[![Docs](https://img.shields.io/badge/docs-GitHub%20Pages-blue)](https://gigurra.github.io/bookeeper/)

> **Note**: This is an **experiment/prototype/weekend hack** I built when switching to Linux on my gaming machine. It works for my use case but comes with limitations. Use at your own risk!

## What It Does

- Activate/deactivate mods via symlinks
- Save and load mod profiles (different mod combinations)
- Install BG3 Script Extender (bg3se)
- Support multiple versions of the same mod
- Shell autocompletion (bash/zsh/fish/powershell)

## Quick Start

```bash
# Install
go install github.com/GiGurra/bookeeper@latest

# See current status
bookeeper status

# Make a downloaded mod available
bookeeper mods make-available ~/Downloads/SomeMod.zip

# Activate a mod
bookeeper mods activate "ModName" "version"

# Save your mod setup as a profile
bookeeper profiles save my-profile

# Later, restore that profile
bookeeper profiles load my-profile
```

## How It Works

1. **Mod Storage**: Mods live in `~/.local/share/bookeeper/downloaded_mods/`
2. **Activation**: Creates symlinks from your mods to BG3's mod directory
3. **modsettings.lsx**: Automatically updated so BG3 loads the mods
4. **Profiles**: JSON files storing which mods are active

```
~/.local/share/bookeeper/
├── downloaded_mods/
│   └── ModName/
│       └── version/
│           ├── mod.pak
│           └── info.json
└── profiles/
    └── my-profile/
        └── profile.json
```

## Key Commands

| Command | Description |
|---------|-------------|
| `bookeeper status` | Show everything (paths, active mods, profiles) |
| `bookeeper mods activate <name> <version>` | Activate a mod |
| `bookeeper mods deactivate <name> <version>` | Deactivate a mod |
| `bookeeper mods make-available <zip>` | Import a mod from zip |
| `bookeeper profiles save <name>` | Save current mods as profile |
| `bookeeper profiles load <name>` | Load a profile |
| `bookeeper bg3se install` | Install BG3 Script Extender |

## Limitations

This is a prototype with known limitations:

- No automatic mod downloading (manual zip download required)
- No mod dependency management
- No conflict detection
- Limited error handling
- Only tested on Linux/Steam
- Load order controlled by activation order only
- Mods must have valid `info.json` files

## Documentation

Full documentation at **[gigurra.github.io/bookeeper](https://gigurra.github.io/bookeeper/)**

- [Installation](https://gigurra.github.io/bookeeper/guide/installation/)
- [Basic Usage](https://gigurra.github.io/bookeeper/guide/usage/)
- [Profiles](https://gigurra.github.io/bookeeper/guide/profiles/)
- [Troubleshooting](https://gigurra.github.io/bookeeper/guide/troubleshooting/)
- [Command Reference](https://gigurra.github.io/bookeeper/reference/commands/)

## License

MIT
