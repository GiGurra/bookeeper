# File Formats

Reference for file formats used by bookeeper.

## info.json

Every mod must have an `info.json` file describing its contents.

### Location

Inside the mod zip, typically at the root or in a subdirectory alongside `.pak` files.

### Format

```json
{
  "Mods": [
    {
      "Author": "ModAuthor",
      "Name": "ModName",
      "Folder": "ModFolder",
      "Version": "1.0.0",
      "Description": "What this mod does",
      "UUID": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
      "Created": "2024-01-01T00:00:00Z",
      "Dependencies": [],
      "Group": "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy"
    }
  ],
  "MD5": "abc123..."
}
```

### Required Fields

| Field | Description |
|-------|-------------|
| `Name` | Display name of the mod |
| `UUID` | Unique identifier (must match pak file) |
| `Folder` | Folder name (usually same as Name) |
| `Version` | Version string (can be simple like "1") |

### Optional Fields

| Field | Description |
|-------|-------------|
| `Author` | Mod creator |
| `Description` | What the mod does |
| `Created` | Creation timestamp |
| `Dependencies` | List of required mods |
| `Group` | Mod group UUID |
| `MD5` | Checksum |

### Common Issues

Some mods use non-standard key names:

| Wrong | Correct |
|-------|---------|
| `modName` | `Name` |
| `folderName` | `Folder` |
| `version` | `Version` |
| `mods` | `Mods` |

If a mod has wrong keys, either:

1. Fix the `info.json` in the zip before importing
2. Fix it in `~/.local/share/bookeeper/downloaded_mods/ModName/version/info.json` after importing

## profile.json

Profiles are stored as JSON files.

### Location

```
~/.local/share/bookeeper/profiles/{profile-name}/profile.json
```

### Format

```json
{
  "name": "my-profile",
  "mods": [
    {
      "name": "ModName1",
      "uuid": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
      "version": "1.0",
      "folder": "ModName1",
      "md5": "abc123..."
    },
    {
      "name": "ModName2",
      "uuid": "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy",
      "version": "2.0",
      "folder": "ModName2",
      "md5": "def456..."
    }
  ]
}
```

### Fields

| Field | Description |
|-------|-------------|
| `name` | Profile name |
| `mods` | Array of mods in load order |

### Mod Entry Fields

| Field | Description |
|-------|-------------|
| `name` | Mod display name |
| `uuid` | Mod UUID |
| `version` | Mod version |
| `folder` | Mod folder name |
| `md5` | Checksum (optional) |

### Manual Editing

You can manually edit profile.json to:

- Change load order (reorder the array)
- Add/remove mods
- Fix version mismatches

!!! note "GustavDev"
    GustavDev (base game) is never stored in profiles.

## modsettings.lsx

BG3's mod configuration file. Bookeeper reads and modifies this automatically.

### Location

```
{Steam}/steamapps/compatdata/1086940/pfx/drive_c/users/steamuser/AppData/Local/Larian Studios/Baldur's Gate 3/PlayerProfiles/Public/modsettings.lsx
```

### Format

XML-based LSX format:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<save>
  <version major="4" minor="7" revision="1" build="3"/>
  <region id="ModuleSettings">
    <node id="root">
      <children>
        <node id="Mods">
          <children>
            <node id="ModuleShortDesc">
              <attribute id="Folder" type="LSString" value="GustavDev"/>
              <attribute id="MD5" type="LSString" value=""/>
              <attribute id="Name" type="LSString" value="GustavDev"/>
              <attribute id="PublishHandle" type="uint64" value="0"/>
              <attribute id="UUID" type="guid" value="28ac9ce2-2aba-8cda-b3b5-6e922f71b6b8"/>
              <attribute id="Version64" type="int64" value="36028797018963968"/>
            </node>
            <!-- More mods here -->
          </children>
        </node>
      </children>
    </node>
  </region>
</save>
```

### Important Notes

- **Don't edit manually** unless you know what you're doing
- Bookeeper handles this file automatically
- GustavDev must always be present
- Order in this file = load order

## Mod Storage Structure

After importing mods:

```
~/.local/share/bookeeper/downloaded_mods/
├── ModName1/
│   └── version1/
│       ├── ModName1.pak
│       └── info.json
├── ModName2/
│   ├── version1/
│   │   ├── ModName2.pak
│   │   └── info.json
│   └── version2/
│       ├── ModName2.pak
│       └── info.json
```

Multiple versions of the same mod can coexist.
