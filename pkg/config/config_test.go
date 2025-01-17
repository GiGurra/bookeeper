package config

import (
	"github.com/GiGurra/boa/pkg/boa"
	"os"
	"testing"
)

func validateConfig[CFG any](cfg *CFG) *CFG {
	boa.Wrap{Params: cfg}.ToCmd()
	return cfg
}

func TestSteamPath(t *testing.T) {
	cfg := validateConfig(&BaseConfig{})
	if SteamPath(cfg) != HomeDir()+"/.local/share/steam" {
		t.Fatalf("SteamPath() returned unexpected value")
	}
}

func TestHomeDir(t *testing.T) {
	if HomeDir() == "" {
		t.Fatalf("HomeDir() returned empty string")
	}

	realHomeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to detect user home dir (used for finding path to game dir)")
	}

	if HomeDir() != realHomeDir {
		t.Fatalf("HomeDir() returned unexpected value")
	}
}
