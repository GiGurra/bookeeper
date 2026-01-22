# Troubleshooting

Common issues and how to fix them.

## Mod Import Fails

### Missing info.json

**Problem**: `bookeeper mods make-available` fails with missing metadata.

**Solution**: The mod zip must contain an `info.json` file. If it doesn't, create one manually.

### Bad info.json Keys

**Problem**: Some mods have `info.json` with wrong key names (e.g., `modName` instead of `Name`).

**Expected format**:

```json
{
  "Mods": [
    {
      "Name": "ModName",
      "UUID": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
      "Folder": "ModName",
      "Version": "1"
    }
  ]
}
```

**Common fixes**:

| Wrong Key | Correct Key |
|-----------|-------------|
| `modName` | `Name` |
| `folderName` | `Folder` |
| `version` (lowercase) | `Version` |

**Example - 5eSpells fix**:

Original (broken):
```json
{
  "mods": [
    {
      "modName": "5eSpells",
      "UUID": "fb5f528d-4d48-4bf2-a668-2274d3cfba96",
      "folderName": "5eSpells",
      "version": "1"
    }
  ]
}
```

Fixed:
```json
{
  "Mods": [
    {
      "Name": "5eSpells",
      "UUID": "fb5f528d-4d48-4bf2-a668-2274d3cfba96",
      "Folder": "5eSpells",
      "Version": "1"
    }
  ]
}
```

## BG3 Not Found

**Problem**: Bookeeper can't find your BG3 installation.

**Solution**: Override the Steam path:

```bash
bookeeper --steam-path /path/to/your/steam status
```

Or set environment variable:

```bash
export STEAM_PATH=/path/to/your/steam
bookeeper status
```

## Mods Not Loading in Game

### Check modsettings.lsx

```bash
cat "$(bookeeper get bg3-modsettings-path)"
```

Your mods should be listed in the XML.

### Check Symlinks

```bash
ls -la "$(bookeeper get bg3-mod-dir)"
```

You should see symlinks pointing to your bookeeper mod storage.

### Check Mod Compatibility

- Is the mod compatible with your BG3 version?
- Does the mod require bg3se? (Install it if so)
- Check the mod's page for known issues

## Profile Load Fails

### Missing Mods

**Problem**: Profile references mods you don't have imported.

**Solution**: Import the missing mods first:

```bash
bookeeper mods make-available path/to/missing-mod.zip
```

### Version Mismatch

**Problem**: Profile has different mod version than what's imported.

**Solution**: Either:

1. Import the correct version
2. Edit the profile JSON to use your version

## GustavDev Issues

**Problem**: Something wrong with the base game mod.

**Note**: Bookeeper never touches GustavDev - it's the base game. If GustavDev is corrupted:

1. Verify game files in Steam
2. Check `modsettings.lsx` has GustavDev listed

## Debug: Get All Paths

```bash
bookeeper get all
```

This shows all paths bookeeper uses - helpful for debugging.

## Still Stuck?

Since this is a prototype, error handling is limited. Try:

1. Check `bookeeper status` for clues
2. Look at the actual files in `~/.local/share/bookeeper/`
3. Check BG3's `modsettings.lsx` directly
4. Verify symlinks in BG3's Mods folder

If all else fails, you can always:

```bash
# Reset to no mods
bookeeper mods deactivate-all

# Or manually delete bookeeper data and start fresh
rm -rf ~/.local/share/bookeeper
```
