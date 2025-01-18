package domain

import (
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

	mod := func() modzip.ModDataEntry {
		if len(modData.Mods) == 0 {
			panic("No mods found in zip")
		}
		if len(modData.Mods) > 2 {
			panic("Multiple mods found in zip, not supported")
		}
		if len(modData.Mods) == 2 {
			// only OK if one is GustavDev (which is the game itself)
			first := modData.Mods[0]
			second := modData.Mods[1]
			if first.Name == "GustavDev" {
				return second
			} else if second.Name == "GustavDev" {
				return first
			} else {
				panic("Multiple mods found in zip, not supported")
			}
		}
		return modData.Mods[0]
	}()

	// calculate the mod path
	modPath := filepath.Join(config.DownloadedModsDir(cfg), mod.Name, mod.Version)

	// create the mod folder
	err := os.MkdirAll(modPath, os.ModePerm)
	if err != nil {
		panic(fmt.Errorf("failed to create mod folder: %w", err))
	}

	// copy the pak files
	modzip.ExtractSpecificFilesFromZip(zipFilePath, append(pakFiles, "info.json"), modPath)

	// DONE!
	fmt.Printf("Mod %s@%s is now available\n", mod.Name, mod.Version)
}
