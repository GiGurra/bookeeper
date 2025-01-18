# Bookeeper

A basic CLI mod manager for Baldur's Gate 3 on Linux/Steam. This is just an experiment/prototype.

## Features

- Manage mod profiles (save/load different mod combinations)
- Download and install BG3 Script Extender
- Basic mod management (activation/deactivation)
- Support for multiple versions of the same mod

## Installation

Requires Go 1.23.5 or later.

```bash
go install github.com/GiGurra/bookeeper@latest
```

## Usage

```bash
# Show status (active mods, available mods, profiles)
bookeeper status

# Manage mods
bookeeper mods activate <mod-name> <version>
bookeeper mods deactivate <mod-name> <version>
bookeeper mods make-available <mod-zip-path>
bookeeper mods make-unavailable <mod-name> <version>

# Manage profiles
bookeeper profiles save <profile-name>
bookeeper profiles load <profile-name>
bookeeper profiles delete <profile-name>

# Install BG3 Script Extender
bookeeper bg3se install
```

## Configuration

Default paths can be overridden with flags or environment variables:

```bash
  --steam-path string           Steam installation path (default "$HOME/.steam/steam")
  --user-data-path string       Steam user data path (default "${SteamPath}/userdata/[0]")
  --mods-install-dir string     BG3 mods installation directory
  --mod-settings-lsx-path string BG3 modsettings.lsx path
```

## Limitations

- Early prototype/experiment
- Limited error handling
- No mod dependency management
- No mod version conflict detection
- Manual mod file downloading required
- Only tested on Linux/Steam

## License

MIT
