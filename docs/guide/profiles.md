# Profiles

Profiles let you save and restore different mod combinations.

## Why Profiles?

- **Different playthroughs**: One profile for a modded run, another for vanilla
- **Testing**: Try new mods without losing your working setup
- **Sharing**: Export profiles (they're just JSON files)

## Saving a Profile

With your desired mods active:

```bash
bookeeper profiles save my-profile-name
```

This saves:

- Which mods are active
- Their versions
- The activation order (which determines load order)

Profile saved to: `~/.local/share/bookeeper/profiles/my-profile-name/profile.json`

## Loading a Profile

```bash
bookeeper profiles load my-profile-name
```

This will:

1. Deactivate all current mods
2. Activate the mods from the profile (in order)

!!! note "GustavDev"
    The game's base mod (GustavDev) is never included in profiles and is always preserved.

## Listing Profiles

```bash
bookeeper profiles list
# or
bookeeper profiles status
```

Or see them in `bookeeper status`.

## Deleting a Profile

```bash
bookeeper profiles delete my-profile-name
```

## Deactivating All Mods

To return to "no mods" state:

```bash
bookeeper profiles deactivate-all
# or
bookeeper mods deactivate-all
```

## Profile File Format

Profiles are simple JSON:

```json
{
  "name": "my-profile",
  "mods": [
    {
      "name": "5eSpells",
      "uuid": "fb5f528d-4d48-4bf2-a668-2274d3cfba96",
      "version": "1",
      "folder": "5eSpells"
    },
    {
      "name": "Party Limit Begone",
      "uuid": "1d6c4030-67b9-4b0a-b3ab-caf6dd73d1af",
      "version": "72902018968059904",
      "folder": "Party Limit Begone"
    }
  ]
}
```

You can manually edit this file to:

- Change load order (reorder the mods array)
- Add/remove mods
- Share with others (they need the same mods imported)

## Typical Workflow

```bash
# Start with a working mod set
bookeeper profiles save working-baseline

# Try some new mods
bookeeper mods make-available ~/Downloads/ExperimentalMod.zip
bookeeper mods activate "ExperimentalMod" "1.0"

# Game crashes? Go back to working setup
bookeeper profiles load working-baseline

# Experimental mod works? Update your profile
bookeeper profiles save working-baseline
```
