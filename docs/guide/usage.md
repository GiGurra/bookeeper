# Basic Usage

!!! note "Prototype"
    Remember, this is a weekend hack. Some workflows might be clunky!

## Workflow Overview

1. **Download** a mod zip file manually (from Nexus, etc.)
2. **Import** it with `bookeeper mods make-available`
3. **Activate** it with `bookeeper mods activate`
4. **Save** your mod set as a profile (optional)
5. **Play** Baldur's Gate 3!

## Checking Status

See everything at once:

```bash
bookeeper status
```

This shows:

- Bookeeper paths
- BG3 paths
- BG3SE installation status
- Active mods
- Current profile
- Available profiles
- Available (imported) mods

## Importing Mods

Download a mod zip file, then import it:

```bash
bookeeper mods make-available ~/Downloads/SomeCoolMod.zip
```

The mod is extracted to `~/.local/share/bookeeper/downloaded_mods/ModName/version/`.

!!! warning "info.json Required"
    Mods must have a valid `info.json` file. See [Troubleshooting](troubleshooting.md) if a mod doesn't import properly.

## Activating Mods

```bash
# List available mods
bookeeper mods list-available

# Activate a specific mod and version
bookeeper mods activate "ModName" "1.0"
```

What happens:

1. Symlinks created from mod's `.pak` files to BG3's Mods directory
2. `modsettings.lsx` updated to tell BG3 to load the mod

## Deactivating Mods

```bash
# Deactivate a specific mod
bookeeper mods deactivate "ModName" "1.0"

# Deactivate ALL mods (except GustavDev)
bookeeper mods deactivate-all
```

## Listing Mods

```bash
# Active mods only
bookeeper mods list-active

# Available (imported) mods only
bookeeper mods list-available

# Both
bookeeper mods list
```

## Removing Mods

To remove a mod from bookeeper storage:

```bash
bookeeper mods make-unavailable "ModName" "1.0"
```

This deletes the mod files from `~/.local/share/bookeeper/downloaded_mods/`.

## Load Order

!!! note "Limitation"
    Load order is controlled by **activation order**. The only way to change load order is to deactivate and reactivate mods in the desired order.

Workaround:

1. Activate mods in the order you want
2. Save as a profile: `bookeeper profiles save my-order`
3. Load the profile whenever: `bookeeper profiles load my-order`

Or manually edit the profile JSON file.

## Example Session

```bash
# Check current state
$ bookeeper status

# Import some mods I downloaded
$ bookeeper mods make-available ~/Downloads/5eSpells.zip
$ bookeeper mods make-available ~/Downloads/PartyLimitBegone.zip

# Activate them
$ bookeeper mods activate "5eSpells" "1"
$ bookeeper mods activate "Party Limit Begone" "72902018968059904"

# Save this setup
$ bookeeper profiles save dnd-style

# Later, start fresh
$ bookeeper mods deactivate-all

# And restore my mods
$ bookeeper profiles load dnd-style
```

## Next Steps

- [Profiles](profiles.md) - More on saving and loading mod sets
- [BG3 Script Extender](bg3se.md) - Install bg3se for script mods
- [Troubleshooting](troubleshooting.md) - When things go wrong
