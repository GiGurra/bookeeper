package config

import (
	"fmt"
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/google/go-cmp/cmp"
	"log/slog"
	"os"
	"path/filepath"
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

	if ExistsDir(Bg3Path(cfg)) {
		slog.Info("Bg3Path() returned a directory")
	} else {
		t.Fatalf("Bg3Path() did not return a directory")
	}
}

func TestBg3binPath(t *testing.T) {
	cfg := validateConfig(&BaseConfig{})
	result := Bg3binPath(cfg)

	slog.Info(fmt.Sprintf("Bg3binPath(): %s", result))

	if !ExistsDir(result) {
		t.Fatalf("Bg3binPath() did not return a directory")
	}
}

func TestUserDataDir(t *testing.T) {
	cfg := validateConfig(&BaseConfig{})
	expectedPrefix := HomeDir() + "/.local/share/Steam/userdata/"
	result := UserDataDir(cfg)
	if len(result) < len(expectedPrefix) {
		t.Fatalf("UserDataDir() returned unexpected value")
	}
	resultPrefix := result[:len(expectedPrefix)]
	if diff := cmp.Diff(resultPrefix, expectedPrefix); diff != "" {
		t.Fatalf("UserDataDir() returned unexpected value, diff: %s", diff)
	}

	if !ExistsDir(result) {
		t.Fatalf("UserDataDir() did not return a directory")
	}

	slog.Info(fmt.Sprintf("UserDataDir(): %s", result))
}

func TestBg3SaveDir(t *testing.T) {
	cfg := validateConfig(&BaseConfig{})
	result := Bg3SaveDir(cfg)

	slog.Info(fmt.Sprintf("Bg3SaveDir(): %s", result))

	if !ExistsDir(result) {
		t.Fatalf("Bg3SaveDir() did not return a directory")
	}
}

func TestBg3ProfileDir(t *testing.T) {
	cfg := validateConfig(&BaseConfig{})
	result := Bg3ProfileDir(cfg)

	slog.Info(fmt.Sprintf("Bg3ProfileDir(): %s", result))

	if !ExistsDir(result) {
		t.Fatalf("Bg3ProfileDir() did not return a directory")
	}

	// should be a modsettings.lsx file in the profile dir
	modSettingsPath := filepath.Join(result, "modsettings.lsx")
	if !ExistsFile(modSettingsPath) {
		t.Fatalf("Bg3ProfileDir() did not return a directory with a modsettings.lsx file")
	}
}

func TestBg3ModsettingsFilePath(t *testing.T) {
	cfg := validateConfig(&BaseConfig{})
	result := Bg3ModsettingsFilePath(cfg)

	slog.Info(fmt.Sprintf("Bg3ModsettingsFilePath(): %s", result))

	if !ExistsFile(result) {
		t.Fatalf("Bg3ModsettingsFilePath() did not return a file")
	}
}

func TestBg3SeDllPath(t *testing.T) {
	cfg := validateConfig(&BaseConfig{})
	result := Bg3SeDllPath(cfg)

	slog.Info(fmt.Sprintf("Bg3SeDllPath(): %s", result))

	// should be BG3Dir/bin
	expected := filepath.Join(Bg3binPath(cfg), "DWrite.dll")
	if result != expected {
		t.Fatalf("Bg3SeDllPath() returned unexpected value")
	}
}

func TestBooKeeperCfgDir(t *testing.T) {
	cfg := validateConfig(&BaseConfig{})
	result := BooKeeperCfgDir(cfg)

	slog.Info(fmt.Sprintf("BooKeeperCfgDir(): %s", result))

	if !ExistsDir(result) {
		t.Fatalf("BooKeeperCfgDir() did not return a directory")
	}
}
