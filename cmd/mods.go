package cmd

import (
	"fmt"
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/GiGurra/bookeeper/pkg/common"
	"github.com/GiGurra/bookeeper/pkg/config"
	"github.com/GiGurra/bookeeper/pkg/domain"
	"github.com/GiGurra/bookeeper/pkg/gui_tree"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/xlab/treeprint"
	"strconv"
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
			ModsStatusCmd(),
			MostList(),
			MostListAvailable(),
			MostListActive(),
		},
	}.ToCmd()
}

type ModsActivateCmdConfig struct {
	Base       config.BaseConfig
	ModName    boa.Required[string] `positional:"true" description:"mod to activate"`
	ModVersion boa.Required[string] `positional:"true" description:"version of the mod to activate"`
}

func ModsActivateCmd() *cobra.Command {

	cfg := &ModsActivateCmdConfig{}

	return boa.Wrap{
		Use:           "activate",
		Short:         "activate a specific mod",
		Params:        cfg,
		ParamEnrich:   boa.ParamEnricherDefault,
		ValidArgsFunc: ValidAvailableModNameAndVersionArgsFunc(&cfg.Base),
		Run: func(cmd *cobra.Command, args []string) {

			fmt.Printf("activating mod %s, v %s\n", cfg.ModName.Value(), cfg.ModVersion.Value())

			if cfg.ModName.Value() == "GustavDev" {
				common.ExitWithUserError("Not allowed to activate GustavDev")
			}

			domain.ActivateMod(&cfg.Base, cfg.ModName.Value(), cfg.ModVersion.Value())

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
		Short:         "deactivate a specific mod",
		Params:        cfg,
		ParamEnrich:   boa.ParamEnricherDefault,
		ValidArgsFunc: ValidActiveModNameAndVersionArgsFunc(&cfg.Base),
		Run: func(cmd *cobra.Command, args []string) {

			fmt.Printf("deactivating mod %s, v %s\n", cfg.ModName.Value(), cfg.ModVersion.Value())

			if cfg.ModName.Value() == "GustavDev" {
				common.ExitWithUserError("Not allowed to deactivate GustavDev")
			}

			domain.DeactivateMod(&cfg.Base, cfg.ModName.Value(), cfg.ModVersion.Value())
		},
	}.ToCmd()
}

func ModsDeactivateAllCmd() *cobra.Command {

	cfg := &config.BaseConfig{}

	return boa.Wrap{
		Use:         "deactivate-all",
		Short:       "deactivate all active mods",
		Params:      cfg,
		ParamEnrich: boa.ParamEnricherDefault,
		Run: func(cmd *cobra.Command, args []string) {
			domain.DeactivateAllMods(cfg)
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

func ModsStatusCmd() *cobra.Command {

	cfg := &config.BaseConfig{}

	return boa.Wrap{
		Use:         "status",
		Short:       "print mod status",
		Params:      cfg,
		ParamEnrich: boa.ParamEnricherDefault,
		Run: func(cmd *cobra.Command, args []string) {

			bg3SeDllPath := config.Bg3SeDllPath(cfg)
			bg3SeInstalled := config.ExistsFile(bg3SeDllPath)

			rootNode := treeprint.New() // NewWithRoot("Bookeeper Status")

			///////////// bg3 paths /////////////////////////////////////////
			bg3PathsNode := gui_tree.AddChildStr(rootNode, "bg3 paths")

			installNode := gui_tree.AddKV(bg3PathsNode, "bg3 install dir", config.Bg3Path(cfg))

			bg3SeNode := gui_tree.AddChildStr(installNode, "bg3se status")
			gui_tree.AddKV(bg3SeNode, "installed", strconv.FormatBool(bg3SeInstalled))
			gui_tree.MakeChildrenSameKeyLen(bg3SeNode)

			///////////// Active mods /////////////////////////////////////////
			activeModsTitle := "active mods"
			if cfg.Verbose.Value() {
				activeModsTitle += " (modsettings.lsx)"
			}
			bg3ActiveModsNode := gui_tree.AddChildStr(rootNode, activeModsTitle)
			for _, mod := range domain.ListActiveMods(cfg) {
				gui_tree.AddChild(bg3ActiveModsNode, gui_tree.DomainMod(mod, cfg.Verbose.Value()))
				//gui_tree.AddKV(bg3ActiveModsNode, mod.Name, fmt.Sprintf("%s, v%s", mod.UUID, mod.Version64))
			}
			gui_tree.MakeChildrenSameKeyLen(bg3ActiveModsNode)

			///////////// Profiles /////////////////////////////////////////

			gui_tree.AddKV(rootNode, "active profile", domain.ActiveProfileName(cfg))

			availableProfilesTitle := "available profiles"
			if cfg.Verbose.Value() {
				availableProfilesTitle += " (" + config.ProfilesDir(cfg) + ")"
			}
			gui_tree.AddChild(rootNode, gui_tree.DomainProfilesN(availableProfilesTitle, domain.ListProfiles(cfg), cfg.Verbose.Value()))

			///////////// Available mods /////////////////////////////////////////
			availableModsTitle := "available mods"
			if cfg.Verbose.Value() {
				availableModsTitle += " (" + config.DownloadedModsDir(cfg) + ")"
			}
			bg3DownloadedModsNode := rootNode.AddBranch(availableModsTitle)
			for _, mod := range domain.ListAvailableMods(cfg) {
				gui_tree.AddChild(bg3DownloadedModsNode, gui_tree.DomainMod(mod, cfg.Verbose.Value()))
			}
			gui_tree.MakeChildrenSameKeyLen(bg3DownloadedModsNode)

			/////////////////////////////////////////////////////////////////////
			fmt.Println(rootNode.String())
		},
	}.ToCmd()
}

func MostListActive() *cobra.Command {

	cfg := &config.BaseConfig{}

	return boa.Wrap{
		Use:         "list-active",
		Short:       "list active mods",
		Params:      cfg,
		ParamEnrich: boa.ParamEnricherDefault,
		Run: func(cmd *cobra.Command, args []string) {

			rootNode := treeprint.New() // NewWithRoot("Bookeeper Status")

			///////////// Active mods /////////////////////////////////////////
			activeModsTitle := "active mods"
			if cfg.Verbose.Value() {
				activeModsTitle += " (modsettings.lsx)"
			}
			bg3ActiveModsNode := gui_tree.AddChildStr(rootNode, activeModsTitle)
			for _, mod := range domain.ListActiveMods(cfg) {
				gui_tree.AddChild(bg3ActiveModsNode, gui_tree.DomainMod(mod, cfg.Verbose.Value()))
				//gui_tree.AddKV(bg3ActiveModsNode, mod.Name, fmt.Sprintf("%s, v%s", mod.UUID, mod.Version64))
			}
			gui_tree.MakeChildrenSameKeyLen(bg3ActiveModsNode)

			/////////////////////////////////////////////////////////////////////
			fmt.Println(rootNode.String())
		},
	}.ToCmd()
}

func MostList() *cobra.Command {

	cfg := &config.BaseConfig{}

	return boa.Wrap{
		Use:         "list",
		Short:       "list active and available mods",
		Params:      cfg,
		ParamEnrich: boa.ParamEnricherDefault,
		Run: func(cmd *cobra.Command, args []string) {

			rootNode := treeprint.New() // NewWithRoot("Bookeeper Status")

			///////////// Active mods /////////////////////////////////////////
			activeModsTitle := "active mods"
			if cfg.Verbose.Value() {
				activeModsTitle += " (modsettings.lsx)"
			}
			bg3ActiveModsNode := gui_tree.AddChildStr(rootNode, activeModsTitle)
			for _, mod := range domain.ListActiveMods(cfg) {
				gui_tree.AddChild(bg3ActiveModsNode, gui_tree.DomainMod(mod, cfg.Verbose.Value()))
				//gui_tree.AddKV(bg3ActiveModsNode, mod.Name, fmt.Sprintf("%s, v%s", mod.UUID, mod.Version64))
			}
			gui_tree.MakeChildrenSameKeyLen(bg3ActiveModsNode)

			///////////// Available mods /////////////////////////////////////////
			availableModsTitle := "available mods"
			if cfg.Verbose.Value() {
				availableModsTitle += " (" + config.DownloadedModsDir(cfg) + ")"
			}
			bg3DownloadedModsNode := rootNode.AddBranch(availableModsTitle)
			for _, mod := range domain.ListAvailableMods(cfg) {
				gui_tree.AddChild(bg3DownloadedModsNode, gui_tree.DomainMod(mod, cfg.Verbose.Value()))
			}
			gui_tree.MakeChildrenSameKeyLen(bg3DownloadedModsNode)

			/////////////////////////////////////////////////////////////////////
			fmt.Println(rootNode.String())
		},
	}.ToCmd()
}

func MostListAvailable() *cobra.Command {

	cfg := &config.BaseConfig{}

	return boa.Wrap{
		Use:         "list-available",
		Short:       "list available mods",
		Params:      cfg,
		ParamEnrich: boa.ParamEnricherDefault,
		Run: func(cmd *cobra.Command, args []string) {

			rootNode := treeprint.New() // NewWithRoot("Bookeeper Status")

			///////////// Available mods /////////////////////////////////////////
			availableModsTitle := "available mods"
			if cfg.Verbose.Value() {
				availableModsTitle += " (" + config.DownloadedModsDir(cfg) + ")"
			}
			bg3DownloadedModsNode := rootNode.AddBranch(availableModsTitle)
			for _, mod := range domain.ListAvailableMods(cfg) {
				gui_tree.AddChild(bg3DownloadedModsNode, gui_tree.DomainMod(mod, cfg.Verbose.Value()))
			}
			gui_tree.MakeChildrenSameKeyLen(bg3DownloadedModsNode)

			/////////////////////////////////////////////////////////////////////
			fmt.Println(rootNode.String())
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

		modsByName := lo.GroupBy(domain.ListActiveMods(cfg), func(item domain.Mod) string {
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
