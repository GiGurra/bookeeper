package config

import (
	"fmt"
	"github.com/GiGurra/boa/pkg/boa"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	SteamBg3AppID = "1086940"
)

type BaseConfig struct {
	SteamPath    boa.Required[string] `default:"${HOME}/.local/share/Steam" name:"steam-path" short-name:"s"`
	UserDataPath boa.Required[string] `default:"${SteamPath}/userdata/[0]" name:"user-data-path" short-name:"u"`
}

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

	return strings.ReplaceAll(cfg.SteamPath.Value(), "${HOME}", HomeDir())
}

func Bg3Path(cfg *BaseConfig) string {
	return filepath.Join(SteamPath(cfg), "steamapps", "common", "Baldurs Gate 3")
}

func Bg3binPath(cfg *BaseConfig) string {
	return filepath.Join(Bg3Path(cfg), "bin")
}

func BooKeeperDir(cfg *BaseConfig) string {
	return ensureExistsDir(filepath.Join(HomeDir(), ".local", "share", "bookeeper"))
}

func DownloadedModsDir(cfg *BaseConfig) string {
	return ensureExistsDir(filepath.Join(BooKeeperDir(cfg), "downloaded_mods"))
}

func ProfilesDir(cfg *BaseConfig) string {
	return ensureExistsDir(filepath.Join(BooKeeperDir(cfg), "profiles"))
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

// UserDataDir returns the path to the user data directory for the steam user.
// This is typically ~/.local/share/Steam/userdata/<steam-user-id>.
func UserDataDir(cfg *BaseConfig) string {
	result := cfg.UserDataPath.Value()
	result = strings.ReplaceAll(result, "${Home}", HomeDir())
	result = strings.ReplaceAll(result, "${SteamPath}", SteamPath(cfg))
	// Replace all [x] with the index/file listing
	for strings.Contains(result, "[") {
		start := strings.Index(result, "[")
		end := strings.Index(result, "]")
		indexStr := result[start+1 : end]
		index, err := strconv.Atoi(indexStr)
		if err != nil {
			panic(fmt.Errorf("failed to parse index in path: %s: %w", indexStr, err))
		}
		replace := result[start : end+1]
		result = strings.ReplaceAll(result, replace, "")
		currentPath := result[:start]
		entries, err := os.ReadDir(currentPath)
		if err != nil {
			panic(fmt.Errorf("failed to read directory %s: %w", currentPath, err))
		}
		var dirs []string
		for _, entry := range entries {
			if entry.IsDir() {
				dirs = append(dirs, entry.Name())
			}
		}
		if len(dirs) == 0 {
			panic(fmt.Errorf("no entries found in directory %s", currentPath))
		}
		if index >= len(dirs) {
			panic(fmt.Errorf("index %d out of bounds in directory %s", index, currentPath))
		}
		if index < 0 {
			panic(fmt.Errorf("index %d out of bounds in directory %s", index, currentPath))
		}

		result = filepath.Join(result, dirs[index])
	}
	return result
}

// Bg3UserDataDir returns the path to the user data directory for the steam user for BG3.
// This is typically ~/.local/share/Steam/userdata/<steam-user-id>/1086940.
func Bg3UserDataDir(cfg *BaseConfig) string {
	return filepath.Join(UserDataDir(cfg), SteamBg3AppID)
}

func Bg3SaveDir(cfg *BaseConfig) string {
	return filepath.Join(Bg3UserDataDir(cfg), "remote", "_SAVE_Public", "Savegames", "Story")
}

func Bg3ProfileDir(cfg *BaseConfig) string {
	return filepath.Join(Bg3UserDataDir(cfg), "remote", "_PROFILE_Public")
}

func Bg3ModsettingsFilePath(cfg *BaseConfig) string {
	return filepath.Join(Bg3ProfileDir(cfg), "modsettings.lsx")
}
