# Configuration

Bookeeper uses sensible defaults but can be configured for non-standard setups.

## Default Paths

### Bookeeper Storage

```
~/.local/share/bookeeper/
├── downloaded_mods/     # Extracted mod files
└── profiles/            # Saved profiles
```

### BG3 Paths (Linux/Steam)

| Path | Default Location |
|------|-----------------|
| Steam | `~/.steam/steam` |
| BG3 Install | `{Steam}/steamapps/common/Baldurs Gate 3` |
| BG3 Bin | `{BG3}/bin` |
| Mods | `{Steam}/steamapps/compatdata/1086940/pfx/drive_c/users/steamuser/AppData/Local/Larian Studios/Baldur's Gate 3/Mods` |
| modsettings.lsx | `{Steam}/steamapps/compatdata/1086940/pfx/drive_c/users/steamuser/AppData/Local/Larian Studios/Baldur's Gate 3/PlayerProfiles/Public/modsettings.lsx` |

## Customization

### Via Command Line

```bash
bookeeper --steam-path /custom/steam/path status
```

### Via Environment Variables

```bash
export STEAM_PATH=/custom/steam/path
export MODS_INSTALL_DIR=/custom/mods/path
bookeeper status
```

### Available Options

| Flag | Env Var | Description |
|------|---------|-------------|
| `--steam-path` | `STEAM_PATH` | Steam installation root |
| `--user-data-path` | `USER_DATA_PATH` | Steam userdata path |
| `--mods-install-dir` | `MODS_INSTALL_DIR` | BG3 mods directory |
| `--mod-settings-lsx-path` | `MOD_SETTINGS_LSX_PATH` | modsettings.lsx location |

## Path Templates

Paths support variable expansion:

| Variable | Expands To |
|----------|-----------|
| `${HOME}` | Home directory |
| `${SteamPath}` | Value of `--steam-path` |

Example:

```bash
--mods-install-dir '${SteamPath}/custom/location'
```

## Non-Standard Setups

### Custom Steam Location

```bash
export STEAM_PATH=/opt/games/steam
bookeeper status
```

### Flatpak Steam

Steam installed via Flatpak uses different paths:

```bash
export STEAM_PATH=~/.var/app/com.valvesoftware.Steam/.steam/steam
bookeeper status
```

### Wolf Game Streaming

For Wolf containers:

```bash
# Compile with static linking
CGO_ENABLED=0 GOBIN=/etc/wolf/12345678/Steam go install github.com/GiGurra/bookeeper@latest

# Then use from within the container
```

## Verify Configuration

Check that bookeeper finds everything:

```bash
bookeeper get all
```

Or see the full status:

```bash
bookeeper status
```

If paths are wrong, the status output will show incorrect locations.
