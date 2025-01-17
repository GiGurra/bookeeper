package config

import (
	"github.com/GiGurra/boa/pkg/boa"
	"os"
	"strings"
)

type BaseConfig struct {
	SteamPath boa.Required[string] `default:"${HOME}/.local/share/steam" name:"steam-path" short-name:"s"`
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
