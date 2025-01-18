package domain

import (
	"encoding/json"
	"fmt"
	"github.com/GiGurra/bookeeper/pkg/common"
	"github.com/GiGurra/bookeeper/pkg/config"
	"github.com/GiGurra/bookeeper/pkg/modsettingslsx"
	"github.com/GiGurra/bookeeper/pkg/modzip"
	"github.com/samber/lo"
	"os"
	"path/filepath"
	"strings"
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

func IsModAvailable(cfg *config.BaseConfig, modName string, modVersion string) bool {
	modPath := filepath.Join(config.DownloadedModsDir(cfg), modName, modVersion)
	return config.ExistsDir(modPath)
}

func GetAvailableMod(cfg *config.BaseConfig, modName string, modVersion string) *Mod {
	modPath := filepath.Join(config.DownloadedModsDir(cfg), modName, modVersion)
	if !config.ExistsDir(modPath) {
		return nil
	}

	// read info.json to get mod data
	infoJsonPath := filepath.Join(modPath, "info.json")
	if !config.ExistsFile(infoJsonPath) {
		panic(fmt.Errorf("info.json not found in mod: %s", modPath))
	}
	mod := FromZipData(modzip.ReadInfoJson(infoJsonPath))

	return &mod
}

func IsModActive(cfg *config.BaseConfig, modName string, modVersion string) bool {
	modsettings := modsettingslsx.Load(cfg)
	activeMods := ListActiveModsX(modsettings)
	for _, mod := range activeMods {
		if mod.Name == modName && mod.Version64 == modVersion {
			return true
		}
	}
	return false
}

func IsModActiveByName(cfg *config.BaseConfig, modName string) bool {
	modsettings := modsettingslsx.Load(cfg)
	activeMods := ListActiveModsX(modsettings)
	for _, mod := range activeMods {
		if mod.Name == modName {
			return true
		}
	}
	return false
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
					mod := GetAvailableMod(cfg, modRootEntry.Name(), modVersionEntry.Name())
					if mod != nil {
						mods = append(mods, *mod)
					}
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

type PakFileLink struct {
	PathInStorage   string
	PathInModFolder string
}

func CalculatePakFileLinks(
	cfg *config.BaseConfig,
	mod Mod,
) []PakFileLink {

	srcDir := filepath.Join(config.DownloadedModsDir(cfg), mod.Name, mod.Version64)
	trgDir := config.Bg3ModInstallDir(cfg)

	if !config.ExistsDir(srcDir) {
		panic(fmt.Errorf("mod dir %s does not exist. Cannot calculate pak file links for mod activation/deactivation", srcDir))
	}

	entries, err := os.ReadDir(srcDir)
	if err != nil {
		panic(fmt.Errorf("failed to read dir: %w", err))
	}

	var result []PakFileLink
	for _, entry := range entries {
		if strings.HasSuffix(strings.ToLower(entry.Name()), ".pak") {
			srcPath := filepath.Join(srcDir, entry.Name())
			trgPath := filepath.Join(trgDir, entry.Name())
			result = append(result, PakFileLink{
				PathInStorage:   srcPath,
				PathInModFolder: trgPath,
			})
		}
	}

	return result
}

func SetupPakFileLinks(
	links []PakFileLink,
) {
	for _, link := range links {
		fmt.Printf("symlinking %s -> %s\n", link.PathInModFolder, link.PathInStorage)
		err := os.Symlink(link.PathInStorage, link.PathInModFolder)
		if err != nil {
			panic(fmt.Errorf("failed to symlink file: %w", err))
		}
	}
}

func DeletePakFileLinks(
	links []PakFileLink,
) {
	for _, link := range links {
		fmt.Printf("deleting %s\n", link.PathInModFolder)
		err := os.Remove(link.PathInModFolder)
		if err != nil {
			panic(fmt.Errorf("failed to remove file: %w", err))
		}
	}
}

func StoreModsInBg3Cfg(
	cfg *config.BaseConfig,
	modsettings *modsettingslsx.ModSettingsXml,
) {
	newXml := modsettings.ToXML()
	if cfg.Verbose.Value() {
		fmt.Printf("new xml: \n%s\n", newXml)
	}

	xmlSavePath := config.Bg3ModsettingsFilePath(cfg)
	if cfg.Verbose.Value() {
		fmt.Printf("saving to %s\n", xmlSavePath)
	}

	err := os.WriteFile(xmlSavePath, []byte(newXml), 0644)
	if err != nil {
		panic(fmt.Errorf("failed to write file: %w", err))
	}
}

func ListActiveMods(cfg *config.BaseConfig) []Mod {
	modsettings := modsettingslsx.Load(cfg)
	return ListActiveModsX(modsettings)
}

func GetActiveModByName(cfg *config.BaseConfig, modName string) *Mod {
	activeMods := ListActiveMods(cfg)
	for _, mod := range activeMods {
		if mod.Name == modName {
			return &mod
		}
	}
	return nil
}

func GetActiveMod(cfg *config.BaseConfig, modName string, modVersion string) *Mod {
	activeMods := ListActiveMods(cfg)
	for _, mod := range activeMods {
		if mod.Name == modName && mod.Version64 == modVersion {
			return &mod
		}
	}
	return nil
}

func ActivateMod(cfg *config.BaseConfig, modName string, modValue string) {

	mod := GetAvailableMod(cfg, modName, modValue)

	if mod == nil {
		common.ExitWithUserError(fmt.Sprintf("mod %s, v %s not found", modName, modValue))
	} else {

		if IsModActive(cfg, modName, modValue) {
			common.ExitWithUserError(fmt.Sprintf("mod %s, v %s is already active", modName, modValue))
		}

		SetActiveMods(cfg, append(ListActiveMods(cfg), *mod))
	}
}

func SetActiveMods(cfg *config.BaseConfig, newModList []Mod) {

	// check for duplicates
	if len(newModList) != len(lo.Uniq(newModList)) {
		common.ExitWithUserError("Cannot activate mods with duplicate names")
	}

	oldModList := ListActiveMods(cfg)

	if !lo.ContainsBy(newModList, func(m Mod) bool { return m.Name == "GustavDev" }) {
		common.ExitWithUserError("Not allowed to deactivate GustavDev [newModList]")
	}

	if !lo.ContainsBy(oldModList, func(m Mod) bool { return m.Name == "GustavDev" }) {
		common.ExitWithUserError("Not allowed to deactivate GustavDev [oldModList]")
	}

	newModsLkup := lo.SliceToMap(newModList, func(m Mod) (string, Mod) { return m.Name + m.Version64, m })
	oldModsLkup := lo.SliceToMap(oldModList, func(m Mod) (string, Mod) { return m.Name + m.Version64, m })

	modsToDeactivate := lo.Filter(ListActiveMods(cfg), func(m Mod, _ int) bool {
		_, ok := newModsLkup[m.Name+m.Version64]
		return !ok
	})

	modsToActivate := lo.Filter(newModList, func(m Mod, _ int) bool {
		_, ok := oldModsLkup[m.Name+m.Version64]
		return !ok
	})

	for _, mod := range modsToDeactivate {
		DeletePakFileLinks(CalculatePakFileLinks(cfg, mod))
	}

	for _, mod := range modsToActivate {
		SetupPakFileLinks(CalculatePakFileLinks(cfg, mod))
	}

	modsettings := modsettingslsx.Load(cfg)

	SetActiveModsInBg3Cfg(modsettings, newModList)
	StoreModsInBg3Cfg(cfg, modsettings)
}

///////////////// Bridge to modsettingslsx

func ListActiveModsX(n *modsettingslsx.ModSettingsXml) []Mod {
	return listActiveModsC(&n.Region.Categories)
}

func listActiveModsC(n *modsettingslsx.XmlCategories) []Mod {
	result := make([]Mod, 0)
	xmlMods := n.GetXmlMods()
	order := n.GetXmlModOrder()
	// order the xmlMods according to the order in the ModOrder
	handled := make(map[string]bool)
	for _, mod := range order {
		for _, xmlMod := range xmlMods {
			if mod.GetXmlAttributeValue("UUID") == xmlMod.GetXmlAttributeValue("UUID") {
				result = append(result, Mod{
					Folder:        xmlMod.GetXmlAttributeValue("Folder"),
					MD5:           xmlMod.GetXmlAttributeValue("MD5"),
					Name:          xmlMod.GetXmlAttributeValue("Name"),
					PublishHandle: xmlMod.GetXmlAttributeValue("PublishHandle"),
					UUID:          xmlMod.GetXmlAttributeValue("UUID"),
					Version64:     xmlMod.GetXmlAttributeValue("Version64"),
				})
				handled[mod.GetXmlAttributeValue("UUID")] = true
			}
		}
	}
	for _, xmlMod := range xmlMods {
		if _, ok := handled[xmlMod.GetXmlAttributeValue("UUID")]; !ok {
			result = append(result, Mod{
				Folder:        xmlMod.GetXmlAttributeValue("Folder"),
				MD5:           xmlMod.GetXmlAttributeValue("MD5"),
				Name:          xmlMod.GetXmlAttributeValue("Name"),
				PublishHandle: xmlMod.GetXmlAttributeValue("PublishHandle"),
				UUID:          xmlMod.GetXmlAttributeValue("UUID"),
				Version64:     xmlMod.GetXmlAttributeValue("Version64"),
			})
		}
	}
	return result
}

func SetActiveModsInBg3Cfg(n *modsettingslsx.ModSettingsXml, mods []Mod) {

	xmlMods := lo.Map(mods, func(mod Mod, _ int) modsettingslsx.XmlMod {
		return modsettingslsx.XmlMod{
			ID: "ModuleShortDesc",
			Attributes: []modsettingslsx.XmlAttribute{
				{ID: "Folder", Value: mod.Folder, Type: "LSString"},
				{ID: "MD5", Value: mod.MD5, Type: "LSString"},
				{ID: "Name", Value: mod.Name, Type: "LSString"},
				{ID: "PublishHandle", Value: mod.PublishHandle, Type: "uint64"},
				{ID: "UUID", Value: mod.UUID, Type: "guid"},
				{ID: "Version64", Value: mod.Version64, Type: "int64"},
			},
		}
	})

	xmlModOrder := lo.Map(mods, func(mod Mod, _ int) modsettingslsx.XmlMod {
		return modsettingslsx.XmlMod{
			ID: "Module",
			Attributes: []modsettingslsx.XmlAttribute{
				{ID: "UUID", Value: mod.UUID, Type: "FixedString"},
			},
		}
	})

	n.Region.Categories.SetXmlMods(xmlMods)
	n.Region.Categories.SetXmlModOrder(xmlModOrder)
}

func DeactivateMod(
	c *config.BaseConfig,
	modName string,
	modVersion string,
) {

	mod := GetActiveMod(c, modName, modVersion)
	if mod == nil {
		common.ExitWithUserError(fmt.Sprintf("mod %s, v %s not active, nothing to deactivate", modName, modVersion))
	} else {

		newModList := lo.Filter(ListActiveMods(c), func(m Mod, _ int) bool {
			return m.Name != modName || m.Version64 != modVersion
		})

		SetActiveMods(c, newModList)

	}
}

func DeactivateAllMods(
	c *config.BaseConfig,
) {

	newModList := lo.Filter(ListActiveMods(c), func(m Mod, _ int) bool {
		return m.Name == "GustavDev"
	})

	SetActiveMods(c, newModList)

}
