# Bookeeper

A basic CLI mod manager for Baldur's Gate 3 on Linux/Steam. This is just an experiment/prototype/weekend hack I wanted
when I switched over from windows on my game machine.

## Features

- Manages individual mods (activation/deactivation)
- Manages mod profiles (save/load different mod combinations)
- Downloads and installs BG3 Script Extender
- Support for switching between multiple versions of the same mod
- CLI autocompletion (bash/zsh/fish/powershell)

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

When you run `bookeeper mods make-available <mod-zip-path>`, the zip file is unpacked and extracted to the
`bookeeper` mod storage directory, with one directory per mod version, like so:

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

Mods are activated by name and version, by running `bookeeper mods activate <mod-name> <version>`. This will create
symlinks in the bg3 mod directory to all the `.pak` files in the mod version directory under `downloaded_mods`.

For example like this:

```bash
> ll $(bookeeper get bg3-mod-dir)
total 20K
lrwxrwxrwx 1 johkjo johkjo  75 jan 19 00:15 5eSpells.pak -> /home/johkjo/.local/share/bookeeper/downloaded_mods/5eSpells/1/5eSpells.pak
lrwxrwxrwx 1 johkjo johkjo 109 jan 19 00:15 PartyLimitBegone.pak -> '/home/johkjo/.local/share/bookeeper/downloaded_mods/Party Limit Begone/72902018968059904/PartyLimitBegone.pak'
lrwxrwxrwx 1 johkjo johkjo 107 jan 19 00:15 UnlockLevelCurve.pak -> /home/johkjo/.local/share/bookeeper/downloaded_mods/UnlockLevelCurve/72057594037927962/UnlockLevelCurve.pak
lrwxrwxrwx 1 johkjo johkjo 161 jan 19 00:15 UnlockLevelCurve_Patch_5eSpells_Improvement.pak -> /home/johkjo/.local/share/bookeeper/downloaded_mods/UnlockLevelCurve_Patch_5eSpells_Improvement/72057594037927960/UnlockLevelCurve_Patch_5eSpells_Improvement.pak
lrwxrwxrwx 1 johkjo johkjo 135 jan 19 00:15 UnlockLevelCurve_Patch_XP_x0.5.pak -> /home/johkjo/.local/share/bookeeper/downloaded_mods/UnlockLevelCurve_Patch_XP_x0.5/72057594037927960/UnlockLevelCurve_Patch_XP_x0.5.pak
```

During the activation and deactivation processes, `bookeeper` will also adjust the `modsettings.lsx` for BG3 to load the
mod during next startup.

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

```
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
├── [active profile]  op_andersson
├── available profiles
│   ├── no_mods
│   └── op_andersson
│       ├── [Party Limit Begone                         ]  1d6c4030-67b9-4b0a-b3ab-caf6dd73d1af, v 72902018968059904
│       ├── [5eSpells                                   ]  fb5f528d-4d48-4bf2-a668-2274d3cfba96, v 1
│       ├── [UnlockLevelCurve                           ]  d903677e-f24b-48ec-ab20-98dcc116a371, v 72057594037927962
│       ├── [UnlockLevelCurve_Patch_XP_x0.5             ]  e53ae4b5-a922-47ef-b69d-d55c5745a65b, v 72057594037927960
│       └── [UnlockLevelCurve_Patch_5eSpells_Improvement]  4d00bc91-f3cb-430a-b86f-a59a6af2171e, v 72057594037927960
└── available mods
    ├── [5eSpells                                   ]  fb5f528d-4d48-4bf2-a668-2274d3cfba96, v 1
    ├── [Party Limit Begone                         ]  1d6c4030-67b9-4b0a-b3ab-caf6dd73d1af, v 72902018968059904
    ├── [UnlockLevelCurve                           ]  d903677e-f24b-48ec-ab20-98dcc116a371, v 72057594037927962
    ├── [UnlockLevelCurve_Patch_5eSpells_Improvement]  4d00bc91-f3cb-430a-b86f-a59a6af2171e, v 72057594037927960
    └── [UnlockLevelCurve_Patch_XP_x0.5             ]  e53ae4b5-a922-47ef-b69d-d55c5745a65b, v 72057594037927960
```

### Command tree

The command tree is as follows (as of the time of writing this readme):

```
> bookeeper print-cmd-tree 
[bookeeper]  Very basic cli mod manager for Baldur's Gate 3
├── [bg3se]  operations related to bg3se (BG3 Script Extender)
│   ├── [install]  download (from github) and install the latest version of bg3se
│   └── [status ]  show status of bg3se (BG3 Script Extender)
├── [completion]  Generate the autocompletion script for the specified shell
│   ├── [bash      ]  Generate the autocompletion script for bash
│   ├── [fish      ]  Generate the autocompletion script for fish
│   ├── [powershell]  Generate the autocompletion script for powershell
│   └── [zsh       ]  Generate the autocompletion script for zsh
├── [get]  get specific path or info
│   ├── [active-profile               ]  get value of active-profile
│   ├── [all                          ]  get value of all
│   ├── [bg3-bin-dir                  ]  get value of bg3-bin-dir
│   ├── [bg3-dir                      ]  get value of bg3-dir
│   ├── [bg3-mod-dir                  ]  get value of bg3-mod-dir
│   ├── [bg3-modsettings-path         ]  get value of bg3-modsettings-path
│   ├── [bookeeper-dir                ]  get value of bookeeper-dir
│   ├── [bookeeper-downloaded-mods-dir]  get value of bookeeper-downloaded-mods-dir
│   └── [bookeeper-profiles-dir       ]  get value of bookeeper-profiles-dir
├── [help]  Help about any command
├── [mods]  operations on mods
│   ├── [activate        ]  activate a specific mod
│   ├── [deactivate      ]  deactivate a specific mod
│   ├── [deactivate-all  ]  deactivate all active mods
│   ├── [list            ]  list active and available mods
│   ├── [list-active     ]  list active mods
│   ├── [list-available  ]  list available mods
│   ├── [make-available  ]  make a new mod available
│   ├── [make-unavailable]  make a mod unavailable
│   └── [status          ]  print mod status
├── [print-cmd-tree]  print the command tree
├── [profiles]  operations on profiles
│   ├── [deactivate-all]  deactivates all active mods, i.e. any profile
│   ├── [delete        ]  delete a profile
│   ├── [list          ]  status/list of profiles
│   ├── [load          ]  load and activate a profile's mods
│   ├── [save          ]  save current active mods to profile
│   └── [status        ]  status/list of profiles
└── [status]  prints full status (on 'everything')
```

### Load order

The only way currently to control mod load order is to deactivate and reactivate mods in the desired order,
Sorry :D, this might be improved later, if I (or someone else) continues working on this.

The easiest solution right now is to just create your load order once, with
`bookeeper mods activate <mod-name> <version>`, and then save it to a profile with
`bookeeper profiles save <profile-name>`. Then you can just load that profile whenever.

You can also just manually reorder the data in the stored profile json files in `~/.local/share/bookeeper/profiles`.

## Limitations

- Early prototype/experiment
- Limited error handling
- No mod dependency management
- No mod version conflict detection
- Manual mod file downloading required
- Only tested on Linux/Steam
- No support for multiple mods in a single zip file
- No support for reading/analyzing mod .pak files.

## Advanced

### Fixing mods with bad `info.json` files

NOTE: Some mods do come with a `info.json` file, but with incorrectly named keys. This will also not work.
For example `5eSpells` comes with a `json.json` file with bad keys. You need to rename them to the proper key names.

The proper (I think :D, what I found in most mods) `info.json` format is:

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

Required fields (I guess?) are:

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
      "Version": "1"
    }
  ]
}
```

The original `info.json` file for `5eSpells`:

```json
{
  "mods": [
    {
      "modName": "5eSpells",
      "UUID": "fb5f528d-4d48-4bf2-a668-2274d3cfba96",
      "folderName": "5eSpells",
      "version": "1",
      "MD5": ""
    }
  ]
}
```

## License

MIT
