package cmd

import (
	"fmt"
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/GiGurra/bookeeper/pkg/config"
	"github.com/GiGurra/bookeeper/pkg/domain"
	"github.com/GiGurra/bookeeper/pkg/gui_tree"
	"github.com/spf13/cobra"
	"github.com/xlab/treeprint"
	"strconv"
)

func StatusCmd() *cobra.Command {

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

			///////////// bookeeper paths /////////////////////////////////////////
			bookeeperPathsNode := gui_tree.AddChildStr(rootNode, "bookeeper paths")
			gui_tree.AddKV(bookeeperPathsNode, "bookeeper", config.BooKeeperDir(cfg))
			gui_tree.AddKV(bookeeperPathsNode, "downloaded mods", config.DownloadedModsDir(cfg))
			gui_tree.AddKV(bookeeperPathsNode, "profiles", config.ProfilesDir(cfg))
			gui_tree.MakeChildrenSameKeyLen(bookeeperPathsNode)

			///////////// bg3 paths /////////////////////////////////////////
			bg3PathsNode := gui_tree.AddChildStr(rootNode, "bg3 paths")

			installNode := gui_tree.AddKV(bg3PathsNode, "bg3 install dir", config.Bg3Path(cfg))
			gui_tree.AddKV(installNode, "bin", config.Bg3binPath(cfg))

			bg3SeNode := gui_tree.AddChildStr(installNode, "bg3se status")
			gui_tree.AddKV(bg3SeNode, "installed", strconv.FormatBool(bg3SeInstalled))
			gui_tree.AddKV(bg3SeNode, "dll path", bg3SeDllPath)
			gui_tree.MakeChildrenSameKeyLen(bg3SeNode)

			///////////// compatdata /////////////////////////////////////////
			compatdataDir := gui_tree.AddChildStr(bg3PathsNode, "compatdata")
			gui_tree.AddKV(compatdataDir, "mod dir", config.Bg3ModInstallDir(cfg))
			gui_tree.AddKV(compatdataDir, "modsettings.lsx", config.Bg3ModsettingsFilePath(cfg))
			gui_tree.MakeChildrenSameKeyLen(compatdataDir)

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
