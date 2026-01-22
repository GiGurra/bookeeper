# BG3 Script Extender

Some mods require the BG3 Script Extender (bg3se) to function.

## What is bg3se?

BG3 Script Extender adds scripting capabilities to Baldur's Gate 3, enabling more complex mods.

Project: [github.com/Norbyte/bg3se](https://github.com/Norbyte/bg3se)

## Check Status

```bash
bookeeper bg3se status
```

This shows:

- Whether bg3se is installed
- Path to the DLL file

## Install bg3se

```bash
bookeeper bg3se install
```

This will:

1. Fetch the latest release from GitHub
2. Download the zip file
3. Extract `DWrite.dll` to BG3's bin directory

## Steam Launch Configuration

After installing bg3se, you need to configure Steam to load it:

1. Open Steam
2. Right-click Baldur's Gate 3 → Properties
3. In "Launch Options", add:

```
WINEDLLOVERRIDES="DWrite=n,b" %command%
```

This tells Wine/Proton to use the native DWrite.dll (bg3se) instead of the built-in one.

## Verify Installation

After configuring Steam:

1. Launch BG3
2. Check the main menu - bg3se usually shows a version indicator
3. If mods requiring bg3se work, you're good!

## Troubleshooting

### bg3se Not Loading

- Check Steam launch options are set correctly
- Verify `DWrite.dll` exists in BG3's bin directory:
  ```bash
  ls "$(bookeeper get bg3-bin-dir)/DWrite.dll"
  ```

### Wrong Proton Version

Some bg3se versions work better with specific Proton versions. Check:

- [bg3se releases](https://github.com/Norbyte/bg3se/releases) for compatibility notes
- BG3 modding community for recommendations

### Manual Installation

If automatic installation fails:

1. Download from [GitHub releases](https://github.com/Norbyte/bg3se/releases)
2. Extract `DWrite.dll`
3. Copy to: `$(bookeeper get bg3-bin-dir)/`
