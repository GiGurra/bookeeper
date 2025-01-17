package config

import (
	"github.com/GiGurra/boa/pkg/boa"
	"log/slog"
	"os"
	"testing"
)

func validateConfig[CFG any](cfg *CFG) *CFG {
	boa.Wrap{Params: cfg}.ToCmd()
	return cfg
}

func TestSteamPath(t *testing.T) {
	cfg := validateConfig(&BaseConfig{})
	if SteamPath(cfg) != HomeDir()+"/.local/share/Steam" {
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

func TestBg3Path(t *testing.T) {
	cfg := validateConfig(&BaseConfig{})
	if Bg3Path(cfg) != HomeDir()+"/.local/share/Steam/steamapps/common/Baldurs Gate 3" {
		t.Fatalf("Bg3Path() returned unexpected value")
	}

	if isDir(Bg3Path(cfg)) {
		slog.Info("Bg3Path() returned a directory")
	}
}
