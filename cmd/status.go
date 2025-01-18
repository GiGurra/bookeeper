package cmd

import (
	"fmt"
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/GiGurra/bookeeper/pkg/config"
	"github.com/GiGurra/bookeeper/pkg/gui_tree"
	"github.com/GiGurra/bookeeper/pkg/modsettingslsx"
	"github.com/spf13/cobra"
	"github.com/xlab/treeprint"
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

			bookeeperPathsNode := gui_tree.AddChildStr(rootNode, "bookeeper paths")
			gui_tree.AddKV(bookeeperPathsNode, "bookeeper", config.BooKeeperDir(cfg))
			gui_tree.AddKV(bookeeperPathsNode, "downloaded mods", config.DownloadedModsDir(cfg))
			gui_tree.AddKV(bookeeperPathsNode, "profiles", config.ProfilesDir(cfg))
			gui_tree.MakeChildrenSameKeyLen(bookeeperPathsNode)

			bg3PathsNode := gui_tree.AddChildStr(rootNode, "bg3 paths")
			gui_tree.AddKV(bg3PathsNode, "install", config.Bg3Path(cfg))
			gui_tree.AddKV(bg3PathsNode, "bin", config.Bg3binPath(cfg))
			gui_tree.AddKV(bg3PathsNode, "userdata", config.Bg3UserDataDir(cfg))
			gui_tree.AddKV(bg3PathsNode, "mod dir", config.Bg3ModInstallDir(cfg))
			gui_tree.MakeChildrenSameKeyLen(bg3PathsNode)

			bg3SeNode := gui_tree.AddChildStr(rootNode, "bg3se status")
			gui_tree.AddKV(bg3SeNode, "installed", bg3SeInstalled)
			gui_tree.AddKV(bg3SeNode, "dll path", bg3SeDllPath)
			gui_tree.MakeChildrenSameKeyLen(bg3SeNode)

			modXml := modsettingslsx.Load(cfg)
			bg3ActiveModsNode := gui_tree.AddChildStr(rootNode, "bg3 active mods")
			for _, mod := range modXml.GetMods() {
				bg3ActiveModsNode.AddMetaNode(mod.Name, fmt.Sprintf("%s, v%s", mod.UUID, mod.Version64))
			}
			gui_tree.MakeChildrenSameKeyLen(bg3ActiveModsNode)

			bg3CurrentSettings := gui_tree.AddChildStr(rootNode, "current settings")
			gui_tree.AddChildStr(bg3CurrentSettings, "bg3 mod config")
			//bookeeperCurrentProfileNode := gui_tree.AddChildStr(bg3CurrentSettings, "bookeeper current profile")
			//gui_tree.AddChild(bookeeperCurrentProfileNode, gui_tree.Profile(config.GetCurrentProfile(cfg)))

			gui_tree.AddChild(rootNode, gui_tree.ProfilesN(cfg, "available profiles"))

			bg3DownloadedModsNode := rootNode.AddBranch("available mods")
			for _, mod := range config.ListAvailableMods(cfg) {
				gui_tree.AddChild(bg3DownloadedModsNode, gui_tree.ConfigMod(mod))
			}

			fmt.Println(rootNode.String())
		},
	}.ToCmd()
}
