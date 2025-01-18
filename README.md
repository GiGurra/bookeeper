# Bookeeper

A basic CLI mod manager for Baldur's Gate 3 on Linux/Steam. This is just an experiment/prototype/weekend hack I wanted
when I switched over from windows on my game machine.

## Features

- Manages individual mods (activation/deactivation)
- Manages mod profiles (save/load different mod combinations)
- Downloads and install BG3 Script Extender
- Support for switching between multiple versions of the same mod
- Full CLI autocompletion (bash/zsh/fish/powershell)

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

### Mod storage

Mods are stored in `~/.local/share/bookeeper/downloaded_mods` and are not automatically downloaded. You need to download
a mod zip file, and then run `bookeeper mods make-available <mod-zip-path>` to make it available for activation.
NOTE: This assumes the zip file comes with a proper `info.json` file. If not, you need to create one manually.

Downloaded mods are stored with one directory per mod version, like so:

```
downloaded_mods
└── mod-name
    ├── version-1
    │   ├── something.pak
    │   └── info.json
    └── version-2
        ├── somethingelse.pak
        └── info.json
```


### Mod activation

Mods are activated by name and version, by running `bookeeper mods activate <mod-name> <version>`. This will create a
symlinks the bg3 mod directory to all the `.pak` files in the mod version directory under `downloaded_mods`.

### CLI help

Experiment yourself. You can navigate the command tree and use `--help` on each level.

Below are some examples:

```bash
> bookeeper --help
Very basic cli mod manager for Baldur's Gate 3

Usage:
  bookeeper [command]

Available Commands:
  bg3se       operations related to bg3se (BG3 Script Extender)
  completion  Generate the autocompletion script for the specified shell
  get         get specific path or info
  help        Help about any command
  mods        operations on mods
  profiles    operations on profiles
  status      print mod status

Flags:
  -v, --verbose                         (env: VERBOSE) (default false)
  -s, --steam-path string               (env: STEAM_PATH) (default "${HOME}/.steam/steam")
  -u, --user-data-path string           (env: USER_DATA_PATH) (default "${SteamPath}/userdata/[0]")
  -m, --mods-install-dir string         (env: MODS_INSTALL_DIR) (default "${SteamPath}/steamapps/compatdata/1086940/pfx/drive_c/users/steamuser/AppData/Local/Larian Studios/Baldur's Gate 3/Mods")
      --mod-settings-lsx-path string    (env: MOD_SETTINGS_LSX_PATH) (default "${SteamPath}/steamapps/compatdata/1086940/pfx/drive_c/users/steamuser/AppData/Local/Larian Studios/Baldur's Gate 3/PlayerProfiles/Public/modsettings.lsx")
  -h, --help                           help for bookeeper

Use "bookeeper [command] --help" for more information about a command.
````

Showing current status:

```bash
> bookeeper status
.
├── bookeeper paths
│   ├── [bookeeper      ]  /home/johkjo/.local/share/bookeeper
│   ├── [downloaded mods]  /home/johkjo/.local/share/bookeeper/downloaded_mods
│   └── [profiles       ]  /home/johkjo/.local/share/bookeeper/profiles
├── bg3 paths
│   ├── [bg3 install dir]  /home/johkjo/.steam/steam/steamapps/common/Baldurs Gate 3
│   │   ├── [bin]  /home/johkjo/.steam/steam/steamapps/common/Baldurs Gate 3/bin
│   │   └── bg3se status
│   │       ├── [installed]  true
│   │       └── [dll path ]  /home/johkjo/.steam/steam/steamapps/common/Baldurs Gate 3/bin/DWrite.dll
│   └── compatdata
│       ├── [mod dir        ]  /home/johkjo/.steam/steam/steamapps/compatdata/1086940/pfx/drive_c/users/steamuser/AppData/Local/Larian Studios/Baldur's Gate 3/Mods
│       └── [modsettings.lsx]  /home/johkjo/.steam/steam/steamapps/compatdata/1086940/pfx/drive_c/users/steamuser/AppData/Local/Larian Studios/Baldur's Gate 3/PlayerProfiles/Public/modsettings.lsx
├── active mods
│   ├── [GustavDev                                  ]  28ac9ce2-2aba-8cda-b3b5-6e922f71b6b8, v 36028797018963968
│   ├── [Party Limit Begone                         ]  1d6c4030-67b9-4b0a-b3ab-caf6dd73d1af, v 72902018968059904
│   ├── [5eSpells                                   ]  fb5f528d-4d48-4bf2-a668-2274d3cfba96, v 1
│   ├── [UnlockLevelCurve                           ]  d903677e-f24b-48ec-ab20-98dcc116a371, v 72057594037927962
│   ├── [UnlockLevelCurve_Patch_XP_x0.5             ]  e53ae4b5-a922-47ef-b69d-d55c5745a65b, v 72057594037927960
│   └── [UnlockLevelCurve_Patch_5eSpells_Improvement]  4d00bc91-f3cb-430a-b86f-a59a6af2171e, v 72057594037927960
├── available profiles
│   ├── no_mods
│   └── op_andersson
│       ├── [Party Limit Begone]  1d6c4030-67b9-4b0a-b3ab-caf6dd73d1af, v 72902018968059904
│       ├── [5eSpells]  fb5f528d-4d48-4bf2-a668-2274d3cfba96, v 1
│       ├── [UnlockLevelCurve]  d903677e-f24b-48ec-ab20-98dcc116a371, v 72057594037927962
│       ├── [UnlockLevelCurve_Patch_XP_x0.5]  e53ae4b5-a922-47ef-b69d-d55c5745a65b, v 72057594037927960
│       └── [UnlockLevelCurve_Patch_5eSpells_Improvement]  4d00bc91-f3cb-430a-b86f-a59a6af2171e, v 72057594037927960
└── available mods
    ├── [5eSpells                                   ]  fb5f528d-4d48-4bf2-a668-2274d3cfba96, v 1
    ├── [Party Limit Begone                         ]  1d6c4030-67b9-4b0a-b3ab-caf6dd73d1af, v 72902018968059904
    ├── [UnlockLevelCurve                           ]  d903677e-f24b-48ec-ab20-98dcc116a371, v 72057594037927962
    ├── [UnlockLevelCurve_Patch_5eSpells_Improvement]  4d00bc91-f3cb-430a-b86f-a59a6af2171e, v 72057594037927960
    └── [UnlockLevelCurve_Patch_XP_x0.5             ]  e53ae4b5-a922-47ef-b69d-d55c5745a65b, v 72057594037927960
```

### Load order

The only way currently to control mod load order is to deactivate and reactivate mods in the desired order,
Sorry :D, this might be improved later, if I (or someone else) continues working on this.

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

## Advanced

### Fixing mods with bad `info.json` files

NOTE: Some mods do come with a `info.json` file, but with incorrectly named keys. This will also not work.
For example `5eSpells` comes with a `json.json` file with bad keys. You need to rename them to the proper key names.

The proper `info.json` format is:
```json
{
  "Mods": [
    {
      "Author": "Charis",
      "Name": "UnlockLevelCurve_Patch_XP_x0.5",
      "Folder": "UnlockLevelCurve_Patch_XP_x0.5",
      "Version": "72057594037927960",
      "Description": "Halves the XP requirement",
      "UUID": "e53ae4b5-a922-47ef-b69d-d55c5745a65b",
      "Created": "2024-07-22T20:17:59.1988778+02:00",
      "Dependencies": [],
      "Group": "dafe83c2-97b3-4b05-9aef-2e1cc2e1de98"
    }
  ],
  "MD5": "87758161b02ba6eb90ac2a6c92cd746f"
}
```

Required fields are:
- `Name`: The name of the mod
- `Folder`: I dont really know what this is, but I think it has to be correct :D
- `Version`: The version of the mod, but can be something as simple as 1
- `UUID`: A unique identifier for the mod (must match what is encoded in the .pak, I think)

Example of fixed `5eSpells` `info.json`:
```json
{
  "mods": [
    {
      "Name": "5eSpells",
      "UUID": "fb5f528d-4d48-4bf2-a668-2274d3cfba96",
      "Folder": "5eSpells",
      "version": "1",
      "MD5": ""
    }
  ]
}
```
## License

MIT
