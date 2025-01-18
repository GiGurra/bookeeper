package domain

import (
	"encoding/json"
	"fmt"
	"github.com/GiGurra/bookeeper/pkg/config"
	"github.com/GiGurra/bookeeper/pkg/modzip"
	"os"
	"path/filepath"
)

type Mod struct {
	Folder        string
	MD5           string
	Name          string
	PublishHandle string
	UUID          string
	Version64     string
}

func MakeModAvailable(cfg *config.BaseConfig, zipFilePath string) {

	modData, pakFiles := modzip.InspectModZip(zipFilePath)
	mod := modData.Entry()

	// calculate the mod path
	modPath := filepath.Join(config.DownloadedModsDir(cfg), mod.Name, mod.Version)

	// create the mod folder
	err := os.MkdirAll(modPath, os.ModePerm)
	if err != nil {
		panic(fmt.Errorf("failed to create mod folder: %w", err))
	}

	// copy the pak files
	modzip.ExtractSpecificFilesFromZip(zipFilePath, append(pakFiles), modPath)

	// write an info.json file
	bsToWrite, err := json.MarshalIndent(modData, "", "  ")
	if err != nil {
		panic(fmt.Errorf("failed to marshal modData: %w", err))
	}
	infoJsonPath := filepath.Join(modPath, "info.json")
	err = os.WriteFile(infoJsonPath, bsToWrite, os.ModePerm)
	if err != nil {
		panic(fmt.Errorf("failed to write info.json: %w", err))
	}

	// DONE!
	fmt.Printf("Mod %s@%s is now available\n", mod.Name, mod.Version)
}

func MakeModUnavailable(cfg *config.BaseConfig, modName string, modVersion string) {

	// calculate the mod path
	modPath := filepath.Join(config.DownloadedModsDir(cfg), modName, modVersion)

	if !config.ExistsDir(modPath) {
		panic(fmt.Errorf("mod %s@%s is not available", modName, modVersion))
	}

	// TODO: Ask if mod is in use in some profile (or actively)

	// remove dir
	err := os.RemoveAll(modPath)
	if err != nil {
		panic(fmt.Errorf("failed to remove mod folder: %w", err))
	}

	// DONE!
	fmt.Printf("Mod %s@%s is now unavailable\n", modName, modVersion)
}

func ListAvailableMods(cfg *config.BaseConfig) []Mod {
	downloadDir := config.DownloadedModsDir(cfg)
	entries, err := os.ReadDir(downloadDir)
	if err != nil {
		panic(fmt.Errorf("failed to read directory %s: %w", downloadDir, err))
	}
	var mods []Mod
	for _, modRootEntry := range entries {
		if modRootEntry.IsDir() {
			// list versions of the mod
			modPathWoV := filepath.Join(downloadDir, modRootEntry.Name())
			modVersions, err := os.ReadDir(modPathWoV)
			if err != nil {
				panic(fmt.Errorf("failed to read directory %s: %w", modPathWoV, err))
			}
			for _, modVersionEntry := range modVersions {
				if modVersionEntry.IsDir() {
					modPath := filepath.Join(modPathWoV, modVersionEntry.Name())
					// read info.json to get mod data
					infoJsonPath := filepath.Join(modPath, "info.json")
					if !config.ExistsFile(infoJsonPath) {
						panic(fmt.Errorf("info.json not found in mod: %s", modPath))
					}
					mod := FromZipData(modzip.ReadInfoJson(infoJsonPath))
					mods = append(mods, mod)
				}
			}
		}
	}
	return mods
}

func FromZipData(zipData modzip.ModData) Mod {
	mod := zipData.Entry()
	return Mod{
		Folder: mod.Folder,
		MD5:    zipData.MD5,
		Name:   mod.Name,
		//PublishHandle: mod.PublishHandle,
		UUID:      mod.UUID,
		Version64: mod.Version,
	}
}
