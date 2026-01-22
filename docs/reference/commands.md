# Command Reference

Full reference for all bookeeper commands.

## Global Flags

| Flag | Env Var | Default | Description |
|------|---------|---------|-------------|
| `-v, --verbose` | `VERBOSE` | false | Verbose output |
| `-s, --steam-path` | `STEAM_PATH` | `~/.steam/steam` | Steam installation path |
| `-u, --user-data-path` | `USER_DATA_PATH` | (auto) | Steam userdata path |
| `-m, --mods-install-dir` | `MODS_INSTALL_DIR` | (auto) | BG3 mods directory |
| `--mod-settings-lsx-path` | `MOD_SETTINGS_LSX_PATH` | (auto) | modsettings.lsx path |

## bookeeper status

Show comprehensive status.

```bash
bookeeper status
```

Displays:

- Bookeeper paths
- BG3 paths and bg3se status
- Active mods
- Current profile
- Available profiles
- Available mods

## bookeeper mods

### mods activate

Activate a mod.

```bash
bookeeper mods activate <mod-name> <version>
```

Creates symlinks and updates `modsettings.lsx`.

### mods deactivate

Deactivate a mod.

```bash
bookeeper mods deactivate <mod-name> <version>
```

Removes symlinks and updates `modsettings.lsx`.

### mods deactivate-all

Deactivate all mods (except GustavDev).

```bash
bookeeper mods deactivate-all
```

### mods make-available

Import a mod from a zip file.

```bash
bookeeper mods make-available <path-to-zip>
```

Extracts mod to bookeeper storage.

### mods make-unavailable

Remove a mod from bookeeper storage.

```bash
bookeeper mods make-unavailable <mod-name> <version>
```

Deletes mod files from storage.

### mods list

List all mods (active and available).

```bash
bookeeper mods list
```

### mods list-active

List only active mods.

```bash
bookeeper mods list-active
```

### mods list-available

List only available (imported) mods.

```bash
bookeeper mods list-available
```

### mods status

Show mod status (alias for list).

```bash
bookeeper mods status
```

## bookeeper profiles

### profiles save

Save current active mods as a profile.

```bash
bookeeper profiles save <profile-name>
```

### profiles load

Load and activate a profile's mods.

```bash
bookeeper profiles load <profile-name>
```

### profiles delete

Delete a saved profile.

```bash
bookeeper profiles delete <profile-name>
```

### profiles list

List all profiles.

```bash
bookeeper profiles list
```

### profiles status

Show profile status (alias for list).

```bash
bookeeper profiles status
```

### profiles deactivate-all

Deactivate all mods (clear profile).

```bash
bookeeper profiles deactivate-all
```

## bookeeper bg3se

### bg3se install

Download and install BG3 Script Extender.

```bash
bookeeper bg3se install
```

Fetches latest release from GitHub.

### bg3se status

Show bg3se installation status.

```bash
bookeeper bg3se status
```

## bookeeper get

Get specific paths or info.

### get all

Show all paths.

```bash
bookeeper get all
```

### get bookeeper-dir

```bash
bookeeper get bookeeper-dir
# ~/.local/share/bookeeper
```

### get bookeeper-downloaded-mods-dir

```bash
bookeeper get bookeeper-downloaded-mods-dir
# ~/.local/share/bookeeper/downloaded_mods
```

### get bookeeper-profiles-dir

```bash
bookeeper get bookeeper-profiles-dir
# ~/.local/share/bookeeper/profiles
```

### get bg3-dir

```bash
bookeeper get bg3-dir
# ~/.steam/steam/steamapps/common/Baldurs Gate 3
```

### get bg3-bin-dir

```bash
bookeeper get bg3-bin-dir
# ~/.steam/steam/steamapps/common/Baldurs Gate 3/bin
```

### get bg3-mod-dir

```bash
bookeeper get bg3-mod-dir
# .../compatdata/1086940/pfx/drive_c/.../Mods
```

### get bg3-modsettings-path

```bash
bookeeper get bg3-modsettings-path
# .../compatdata/1086940/pfx/drive_c/.../modsettings.lsx
```

### get active-profile

```bash
bookeeper get active-profile
# my-profile (or empty if none)
```

## bookeeper completion

Generate shell completion scripts.

```bash
bookeeper completion bash
bookeeper completion zsh
bookeeper completion fish
bookeeper completion powershell
```

## bookeeper print-cmd-tree

Print the full command hierarchy.

```bash
bookeeper print-cmd-tree
```

## bookeeper help

Get help on any command.

```bash
bookeeper help
bookeeper help mods
bookeeper mods --help
bookeeper mods activate --help
```
