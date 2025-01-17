package cmd

import (
	"fmt"
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/GiGurra/bookeeper/pkg/config"
	"github.com/GiGurra/bookeeper/pkg/gui_tree"
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
			bookeeperPathsNode.AddMetaNode("bookeeper path", config.BooKeeperDir(cfg))
			bookeeperPathsNode.AddMetaNode("downloaded mods path", config.DownloadedModsDir(cfg))
			bookeeperPathsNode.AddMetaNode("profiles path", config.ProfilesDir(cfg))
			gui_tree.MakeChildrenSameKeyLen(bookeeperPathsNode)

			bg3PathsNode := gui_tree.AddChildStr(rootNode, "bg3")
			gui_tree.AddKV(bg3PathsNode, "install path", config.Bg3Path(cfg))
			gui_tree.AddKV(bg3PathsNode, "bin path", config.Bg3binPath(cfg))
			gui_tree.AddKV(bg3PathsNode, "userdata path", config.Bg3UserDataDir(cfg))
			gui_tree.MakeChildrenSameKeyLen(bg3PathsNode)

			bg3SeNode := gui_tree.AddChildStr(rootNode, "bg3se")
			gui_tree.AddKV(bg3SeNode, "installed", bg3SeInstalled)
			gui_tree.AddKV(bg3SeNode, "dll path", bg3SeDllPath)
			gui_tree.MakeChildrenSameKeyLen(bg3SeNode)

			bg3CurrentSettings := gui_tree.AddChildStr(rootNode, "current settings")
			gui_tree.AddChild(bg3CurrentSettings, gui_tree.Profile(config.GetCurrentProfile(cfg)))

			gui_tree.AddChild(rootNode, gui_tree.Profiles(cfg))

			bg3DownloadedModsNode := rootNode.AddBranch("downloaded/available mods")
			for _, mod := range config.ListInstalledMods(cfg) {
				gui_tree.AddChild(bg3DownloadedModsNode, gui_tree.Mod(mod))
			}

			fmt.Println(rootNode.String())
		},
	}.ToCmd()
}
