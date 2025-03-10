package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	SteamBg3AppID = "1086940"
)

func HomeDir() string {
	res, err := os.UserHomeDir()
	if err != nil {
		panic("Failed to detect user home dir (used for finding path to game dir)")
	}

	return res
}

func SteamPath(cfg *BaseConfig) string {
	if !cfg.SteamPath.HasValue() {
		panic("Steam path not set")
	}

	return replaceAllAliases(cfg, cfg.SteamPath.Value())
}

func replaceAllAliases(cfg *BaseConfig, str string) string {

	if !cfg.SteamPath.HasValue() {
		panic("Steam path not set")
	}

	return ResolveStr(str, map[string]string{
		"HOME":      HomeDir(),
		"SteamPath": cfg.SteamPath.Value(),
	})
}

func Bg3Path(cfg *BaseConfig) string {
	return filepath.Join(SteamPath(cfg), "steamapps", "common", "Baldurs Gate 3")
}

func Bg3ModInstallDir(cfg *BaseConfig) string {
	return replaceAllAliases(cfg, cfg.ModsInstallDir.Value())
}

func Bg3binPath(cfg *BaseConfig) string {
	return filepath.Join(Bg3Path(cfg), "bin")
}

func BooKeeperDir(cfg *BaseConfig) string {
	return ensureExistsDir(filepath.Join(HomeDir(), ".local", "share", "bookeeper"))
}

var replaceVarRE = regexp.MustCompile(`\${([^}]+)}`)

func ResolveStr(raw string, lkup map[string]string) string {
	result := raw

	// Keep resolving until no more changes (for nested dependencies)
	for {
		prev := result
		result = replaceVarRE.ReplaceAllStringFunc(result, func(match string) string {
			key := match[2 : len(match)-1] // Remove ${ and }
			if val, ok := lkup[key]; ok {
				return val
			}
			return match
		})
		if prev == result {
			// check if there are any more variables to resolve, if so, panic
			if strings.Contains(result, "$") {
				panic(fmt.Errorf("failed to resolve all variables in string: %s, remaining: %s", raw, result))
			}
			break
		}
	}
	return result
}

// ./.local/share/Steam/steamapps/compatdata/1086940/pfx/drive_c/users/steamuser/AppData/Local/Larian Studios/Baldur's Gate 3/Mods/

func DownloadedModsDir(cfg *BaseConfig) string {
	return ensureExistsDir(filepath.Join(BooKeeperDir(cfg), "downloaded_mods"))
}

func ProfilesDir(cfg *BaseConfig) string {
	return ensureExistsDir(filepath.Join(BooKeeperDir(cfg), "profiles"))
}

func GetProfile(cfg *BaseConfig, profileName string) Profile {
	profilePath := filepath.Join(ProfilesDir(cfg), profileName)
	if !ExistsDir(profilePath) {
		panic(fmt.Errorf("profile %s does not exist", profileName))
	}
	cfgFilePath := filepath.Join(profilePath, "profile.json")
	if !ExistsFile(cfgFilePath) {
		panic(fmt.Errorf("profile %s does not have a config file", profileName))
	}
	bs, err := os.ReadFile(cfgFilePath)
	if err != nil {
		panic(fmt.Errorf("failed to read profile config file %s: %w", cfgFilePath, err))
	}
	var profile Profile
	err = json.Unmarshal(bs, &profile)
	if err != nil {
		panic(fmt.Errorf("failed to unmarshal profile config file %s: %w", cfgFilePath, err))
	}
	return profile
}

func SaveProfile(cfg *BaseConfig, profile Profile) {
	profilePath := filepath.Join(ProfilesDir(cfg), profile.Name)
	if !ExistsDir(profilePath) {
		err := os.MkdirAll(profilePath, 0755)
		if err != nil {
			panic(fmt.Errorf("failed to create profile directory %s: %w", profilePath, err))
		}
	}
	cfgFilePath := filepath.Join(profilePath, "profile.json")
	err := os.WriteFile(cfgFilePath, profile.toJson(), 0644)
	if err != nil {
		panic(fmt.Errorf("failed to write profile config file %s: %w", cfgFilePath, err))
	}
}

func (p *Profile) toJson() []byte {
	bs, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		panic(fmt.Errorf("failed to marshal profile mod list to json: %w", err))
	}
	return bs
}

func ListProfiles(cfg *BaseConfig) []Profile {
	entries, err := os.ReadDir(ProfilesDir(cfg))
	if err != nil {
		panic(fmt.Errorf("failed to read directory %s: %w", ProfilesDir(cfg), err))
	}
	var profiles []Profile
	for _, entry := range entries {
		if entry.IsDir() {
			profile := GetProfile(cfg, entry.Name())
			profiles = append(profiles, profile)
		}
	}
	return profiles
}

func ensureExistsDir(path string) string {
	if !ExistsDir(path) {
		if PathExists(path) {
			panic(fmt.Errorf("path %s exists but is not a directory", path))
		}
		err := os.MkdirAll(path, 0755)
		if err != nil {
			panic(fmt.Errorf("failed to create directory %s: %w", path, err))
		}
	}
	return path
}

func Bg3SeDllPath(cfg *BaseConfig) string {
	return filepath.Join(Bg3binPath(cfg), "/DWrite.dll")
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func ExistsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func ExistsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

//
//// UserDataDir returns the path to the user data directory for the steam user.
//// This is typically ~/.local/share/Steam/userdata/<steam-user-id>.
//func UserDataDir(cfg *BaseConfig) string {
//	result := replaceAllAliases(cfg, cfg.UserDataPath.Value())
//	// Replace all [x] with the index/file listing
//	for strings.Contains(result, "[") {
//		start := strings.Index(result, "[")
//		end := strings.Index(result, "]")
//		indexStr := result[start+1 : end]
//		index, err := strconv.Atoi(indexStr)
//		if err != nil {
//			panic(fmt.Errorf("failed to parse index in path: %s: %w", indexStr, err))
//		}
//		replace := result[start : end+1]
//		result = strings.ReplaceAll(result, replace, "")
//		currentPath := result[:start]
//		entries, err := os.ReadDir(currentPath)
//		if err != nil {
//			panic(fmt.Errorf("failed to read directory %s: %w", currentPath, err))
//		}
//		var dirs []string
//		for _, entry := range entries {
//			if entry.IsDir() {
//				dirs = append(dirs, entry.Name())
//			}
//		}
//		if len(dirs) == 0 {
//			panic(fmt.Errorf("no entries found in directory %s", currentPath))
//		}
//		if index >= len(dirs) {
//			panic(fmt.Errorf("index %d out of bounds in directory %s", index, currentPath))
//		}
//		if index < 0 {
//			panic(fmt.Errorf("index %d out of bounds in directory %s", index, currentPath))
//		}
//
//		result = filepath.Join(result, dirs[index])
//	}
//	return result
//}

// Bg3UserDataDir returns the path to the user data directory for the steam user for BG3.
// This is typically ~/.local/share/Steam/userdata/<steam-user-id>/1086940.
//func Bg3UserDataDir(cfg *BaseConfig) string {
//	return filepath.Join(UserDataDir(cfg), SteamBg3AppID)
//}
//
//func Bg3SaveDir(cfg *BaseConfig) string {
//	return filepath.Join(Bg3UserDataDir(cfg), "remote", "_SAVE_Public", "Savegames", "Story")
//}
//
//func Bg3ProfileDir(cfg *BaseConfig) string {
//	return filepath.Join(Bg3UserDataDir(cfg), "remote", "_PROFILE_Public")
//}
//
//func Bg3UserdataModsettingsFilePath(cfg *BaseConfig) string {
//	return filepath.Join(Bg3ProfileDir(cfg), "modsettings.lsx")
//}

// "/home/johkjo/.local/share/Steam/steamapps/compatdata/1086940/pfx/drive_c/users/steamuser/AppData/Local/Larian Studios/Baldur's Gate 3/PlayerProfiles/Public/modsettings.lsx"
func Bg3ModsettingsFilePath(cfg *BaseConfig) string {
	return replaceAllAliases(cfg, cfg.ModSettingsLsxPath.Value())
}

type Mod struct {
	Name         string `json:"name"`
	DownloadPath string `json:"download_path"`
	Version      string `json:"version"`
}

func ListAvailableMods(cfg *BaseConfig) []Mod {
	entries, err := os.ReadDir(DownloadedModsDir(cfg))
	if err != nil {
		panic(fmt.Errorf("failed to read directory %s: %w", DownloadedModsDir(cfg), err))
	}
	var mods []Mod
	for _, modRootEntry := range entries {
		if modRootEntry.IsDir() {
			// list versions of the mod
			modPath := filepath.Join(DownloadedModsDir(cfg), modRootEntry.Name())
			modVersions, err := os.ReadDir(modPath)
			if err != nil {
				panic(fmt.Errorf("failed to read directory %s: %w", modPath, err))
			}
			for _, modVersionEntry := range modVersions {
				if modVersionEntry.IsDir() {
					mods = append(mods, Mod{
						Name:         modRootEntry.Name(),
						DownloadPath: filepath.Join(modPath, modVersionEntry.Name()),
						Version:      modVersionEntry.Name(),
					})
				}
			}
		}
	}
	return mods
}

type Profile struct {
	Name string
	Path string
	Mods []Mod
}

//
//func BookeperConfigFilePath(cfg *BaseConfig) string {
//	return filepath.Join(BooKeeperDir(cfg), "config.json")
//}
//
//type Config struct {
//	//CurrentProfile string `json:"current_profile"`
//}
//
//func (c *Config) toJson() []byte {
//	bs, err := json.MarshalIndent(c, "", "  ")
//	if err != nil {
//		panic(fmt.Errorf("failed to marshal config to json: %w", err))
//	}
//	return bs
//}

//func GetConfig(cfg *BaseConfig) Config {
//	// <bookeeper_config_file>
//	configFilePath := BookeperConfigFilePath(cfg)
//	if ExistsFile(configFilePath) {
//		config := Config{}
//		bs, err := os.ReadFile(configFilePath)
//		if err != nil {
//			panic(fmt.Errorf("failed to read config file %s: %w", configFilePath, err))
//		}
//		err = json.Unmarshal(bs, &config)
//		if err != nil {
//			panic(fmt.Errorf("failed to unmarshal config file %s: %w", configFilePath, err))
//		}
//		return config
//	} else {
//		// create and save new config
//		fmt.Printf("Creating new config file at %s\n", configFilePath)
//		config := Config{}
//		SaveConfig(cfg, config)
//		return config
//	}
//}
//
//func SaveConfig(cfg *BaseConfig, config Config) {
//	// <bookeeper_config_file>
//	err := os.WriteFile(BookeperConfigFilePath(cfg), config.toJson(), 0644)
//	if err != nil {
//		panic(fmt.Errorf("failed to write config file %s: %w", BookeperConfigFilePath(cfg), err))
//	}
//}

//
//func GetCurrentProfile(cfg *BaseConfig) Profile {
//	name := GetConfig(cfg).CurrentProfile
//	if name == "" {
//		name = "default"
//		fmt.Printf("No current profile set, creating new default profile\n")
//		profile := Profile{
//			Name: name,
//			Path: filepath.Join(ProfilesDir(cfg), name),
//			Mods: []Mod{},
//		}
//		SaveProfile(cfg, profile)
//		currentConfig := GetConfig(cfg)
//		currentConfig.CurrentProfile = name
//		SaveConfig(cfg, currentConfig)
//		return profile
//	} else {
//		return GetProfile(cfg, name)
//	}
//}
