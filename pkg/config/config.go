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
	UserDataPath boa.Required[string] `default:"${SteamPath}/userdata/[0]" name:"steam-user-id" short-name:"u"`
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

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

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
		dirs := []string{}
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

func Bg3UserDataDir(cfg *BaseConfig) string {
	return filepath.Join(UserDataDir(cfg), SteamBg3AppID)
}

func Bg3SaveRoot(cfg *BaseConfig) string {
	return filepath.Join(Bg3UserDataDir(cfg), "Remote")
}
