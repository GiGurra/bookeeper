package config

import (
	"github.com/GiGurra/boa/pkg/boa"
	"os"
	"strings"
)

type BaseConfig struct {
	SteamPath boa.Required[string] `default:"${HOME}/.local/share/Steam" name:"steam-path" short-name:"s"`
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
	return SteamPath(cfg) + "/steamapps/common/Baldurs Gate 3"
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
