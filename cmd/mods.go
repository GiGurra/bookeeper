package cmd

import (
	"fmt"
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/GiGurra/bookeeper/pkg/config"
	"github.com/GiGurra/bookeeper/pkg/domain"
	"github.com/GiGurra/bookeeper/pkg/modsettingslsx"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

func ModsCmd() *cobra.Command {

	cfg := &config.BaseConfig{}

	return boa.Wrap{
		Use:         "mods",
		Short:       "operations on mods",
		Params:      cfg,
		ParamEnrich: boa.ParamEnricherDefault,
		SubCommands: []*cobra.Command{
			ModsActivateCmd(),
			ModsDeactivateCmd(),
			ModsDeactivateAllCmd(),
			ModsMakeAvailableCmd(),
			ModsMakeUnavailableCmd(),
		},
	}.ToCmd()
}

type ModsActivateCmdConfig struct {
	Base       config.BaseConfig
	ModName    boa.Required[string] `positional:"true" description:"mod to activate"`
	ModVersion boa.Required[string] `positional:"true" description:"version of the mod to activate"`
}

type PakFileLink struct {
	PathInStorage   string
	PathInModFolder string
}

func calculatePakFileLinks(
	cfg *config.BaseConfig,
	mod domain.Mod,
) []PakFileLink {

	srcDir := filepath.Join(config.DownloadedModsDir(cfg), mod.Name, mod.Version64)
	trgDir := config.Bg3ModInstallDir(cfg)

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

func setupPakFileLinks(
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

func deletePakFileLinks(
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

func storeModSettings(
	cfg *config.BaseConfig,
	modsettings *modsettingslsx.ModSettingsXml,
) {
	newXml := modsettings.ToXML()
	fmt.Printf("new xml: \n%s\n", newXml)

	xmlSavePath := config.Bg3ModsettingsFilePath(cfg)
	fmt.Printf("saving to %s\n", xmlSavePath)

	err := os.WriteFile(xmlSavePath, []byte(newXml), 0644)
	if err != nil {
		panic(fmt.Errorf("failed to write file: %w", err))
	}
}

func exitWithUserError(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func ModsActivateCmd() *cobra.Command {

	cfg := &ModsActivateCmdConfig{}

	return boa.Wrap{
		Use:           "activate",
		Short:         "activate a specific mod immediately",
		Params:        cfg,
		ParamEnrich:   boa.ParamEnricherDefault,
		ValidArgsFunc: ValidAvailableModNameAndVersionArgsFunc(&cfg.Base),
		Run: func(cmd *cobra.Command, args []string) {

			fmt.Printf("activating mod %s, v %s\n", cfg.ModName.Value(), cfg.ModVersion.Value())

			if cfg.ModName.Value() == "GustavDev" {
				exitWithUserError("Not allowed to activate GustavDev")
			}

			modsettings := modsettingslsx.Load(&cfg.Base)

			activeMods := modsettings.GetMods()
			availableMods := domain.ListAvailableMods(&cfg.Base)
			mod, found := lo.Find(availableMods, func(mod domain.Mod) bool {
				return mod.Name == cfg.ModName.Value() && mod.Version64 == cfg.ModVersion.Value()
			})
			if !found {
				exitWithUserError(fmt.Sprintf("mod %s, v %s not found", cfg.ModName.Value(), cfg.ModVersion.Value()))
			}

			if lo.ContainsBy(activeMods, func(m domain.Mod) bool { return m.Name == cfg.ModName.Value() }) {
				exitWithUserError(fmt.Sprintf("a mod with name %s is already active", cfg.ModName.Value()))
			}

			// First we must copy|symlink the required mod pak files
			setupPakFileLinks(calculatePakFileLinks(&cfg.Base, mod))

			// Then we must update the modsettings file
			modsettings.SetMods(activeMods)
			storeModSettings(&cfg.Base, modsettings)
		},
	}.ToCmd()
}

type ModsDeactivateCmdConfig struct {
	Base       config.BaseConfig
	ModName    boa.Required[string] `positional:"true" description:"mod to deactivate"`
	ModVersion boa.Required[string] `positional:"true" description:"version of the mod to deactivate"`
}

func ModsDeactivateCmd() *cobra.Command {

	cfg := &ModsDeactivateCmdConfig{}

	return boa.Wrap{
		Use:           "deactivate",
		Short:         "deactivate a specific mod immediately",
		Params:        cfg,
		ParamEnrich:   boa.ParamEnricherDefault,
		ValidArgsFunc: ValidActiveModNameAndVersionArgsFunc(&cfg.Base),
		Run: func(cmd *cobra.Command, args []string) {

			if cfg.ModName.Value() == "GustavDev" {
				fmt.Println("Not allowed to deactivate GustavDev")
				os.Exit(1)
			}

			modsettings := modsettingslsx.Load(&cfg.Base)
			modsBefore := modsettings.GetMods()
			modsAfter := lo.Filter(modsBefore, func(mod domain.Mod, _ int) bool {
				return mod.Name != cfg.ModName.Value() || mod.Version64 != cfg.ModVersion.Value()
			})

			if len(modsBefore) == len(modsAfter) {
				fmt.Printf("mod %s, v %s not found\n", cfg.ModName.Value(), cfg.ModVersion.Value())
				os.Exit(1)
			}

			modsettings.SetMods(modsAfter)

			newXml := modsettings.ToXML()
			fmt.Printf("new xml: \n%s\n", newXml)

			xmlSavePath := config.Bg3ModsettingsFilePath(&cfg.Base)
			fmt.Printf("saving to %s\n", xmlSavePath)

			err := os.WriteFile(xmlSavePath, []byte(newXml), 0644)
			if err != nil {
				panic(fmt.Errorf("failed to write file: %w", err))
			}
		},
	}.ToCmd()
}

func ModsDeactivateAllCmd() *cobra.Command {

	cfg := &config.BaseConfig{}

	return boa.Wrap{
		Use:         "deactivate-all",
		Short:       "deactivate all active mods immediately",
		Params:      cfg,
		ParamEnrich: boa.ParamEnricherDefault,
		Run: func(cmd *cobra.Command, args []string) {

			modXml := modsettingslsx.Load(cfg)
			currentMods := modXml.GetMods()
			newModList := lo.Filter(currentMods, func(mod domain.Mod, _ int) bool {
				return mod.Name == "GustavDev"
			})
			modXml.SetMods(newModList)

			newXml := modXml.ToXML()
			fmt.Printf("new xml: \n%s\n", newXml)

			xmlSavePath := config.Bg3ModsettingsFilePath(cfg)
			fmt.Printf("saving to %s\n", xmlSavePath)

			err := os.WriteFile(xmlSavePath, []byte(newXml), 0644)
			if err != nil {
				panic(fmt.Errorf("failed to write file: %w", err))
			}
		},
	}.ToCmd()
}

type ModsMakeAvailableCmdConfig struct {
	Base       config.BaseConfig
	ModZipPath boa.Required[string] `positional:"true" description:"mod zip file path"`
}

func ModsMakeAvailableCmd() *cobra.Command {

	cfg := &ModsMakeAvailableCmdConfig{}

	return boa.Wrap{
		Use:         "make-available",
		Short:       "make a new mod available",
		Params:      cfg,
		ParamEnrich: boa.ParamEnricherDefault,
		Run: func(cmd *cobra.Command, args []string) {
			domain.MakeModAvailable(&cfg.Base, cfg.ModZipPath.Value())
		},
	}.ToCmd()
}

type ModsMakeUnavailableCmdConfig struct {
	Base       config.BaseConfig
	ModName    boa.Required[string] `positional:"true" description:"mod name"`
	ModVersion boa.Required[string] `positional:"true" description:"mod version"`
}

func ModsMakeUnavailableCmd() *cobra.Command {

	cfg := &ModsMakeUnavailableCmdConfig{}

	return boa.Wrap{
		Use:           "make-unavailable",
		Short:         "make a mod unavailable",
		Params:        cfg,
		ParamEnrich:   boa.ParamEnricherDefault,
		ValidArgsFunc: ValidAvailableModNameAndVersionArgsFunc(&cfg.Base),
		Run: func(cmd *cobra.Command, args []string) {
			domain.MakeModUnavailable(&cfg.Base, cfg.ModName.Value(), cfg.ModVersion.Value())
		},
	}.ToCmd()
}

func ValidAvailableModNameAndVersionArgsFunc(cfg *config.BaseConfig) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

		modsByName := lo.GroupBy(domain.ListAvailableMods(cfg), func(item domain.Mod) string {
			return item.Name
		})

		switch len(args) {
		case 0:
			return lo.Keys(modsByName), cobra.ShellCompDirectiveDefault | cobra.ShellCompDirectiveNoFileComp
		default:
			modName := strings.Join(args, " ")
			return lo.Map(modsByName[modName], func(item domain.Mod, _ int) string {
				return item.Version64
			}), cobra.ShellCompDirectiveDefault | cobra.ShellCompDirectiveNoFileComp
		}
	}
}

func ValidActiveModNameAndVersionArgsFunc(cfg *config.BaseConfig) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

		modXml := modsettingslsx.Load(cfg)

		modsByName := lo.GroupBy(modXml.GetMods(), func(item domain.Mod) string {
			return item.Name
		})
		delete(modsByName, "GustavDev")

		switch len(args) {
		case 0:
			return lo.Keys(modsByName), cobra.ShellCompDirectiveDefault | cobra.ShellCompDirectiveNoFileComp
		default:
			modName := strings.Join(args, " ")
			return lo.Map(modsByName[modName], func(item domain.Mod, _ int) string {
				return item.Version64
			}), cobra.ShellCompDirectiveDefault | cobra.ShellCompDirectiveNoFileComp
		}
	}
}
