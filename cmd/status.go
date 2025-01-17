package cmd

import (
	"fmt"
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/GiGurra/bookeeper/pkg/config"
	"github.com/GiGurra/bookeeper/pkg/gui"
	"github.com/spf13/cobra"
	"github.com/xlab/treeprint"
	"strings"
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

			bookeeperPathsNode := rootNode.AddBranch("bookeeper paths")
			bookeeperPathsNode.AddMetaNode("bookeeper path", config.BooKeeperDir(cfg))
			bookeeperPathsNode.AddMetaNode("downloaded mods path", config.DownloadedModsDir(cfg))
			bookeeperPathsNode.AddMetaNode("profiles path", config.ProfilesDir(cfg))
			makeNodeChildrenSameKeyLen(bookeeperPathsNode)

			bg3PathsNode := rootNode.AddBranch("bg3")
			bg3PathsNode.AddMetaNode("install path", config.Bg3Path(cfg))
			bg3PathsNode.AddMetaNode("bin DownloadPath", config.Bg3binPath(cfg))
			bg3PathsNode.AddMetaNode("userdata DownloadPath", config.Bg3UserDataDir(cfg))
			makeNodeChildrenSameKeyLen(bg3PathsNode)

			bg3SeNode := rootNode.AddBranch("bg3se")
			bg3SeNode.AddMetaNode("installed", bg3SeInstalled)
			bg3SeNode.AddMetaNode("dll path", bg3SeDllPath)
			makeNodeChildrenSameKeyLen(bg3SeNode)

			bg3CurrentSettings := rootNode.AddBranch("current settings")
			bg3CurrentSettings.AddBranch(gui.Profile(config.GetCurrentProfile(cfg)))

			rootNode.AddBranch(gui.Profiles(cfg))

			bg3DownloadedModsNode := rootNode.AddBranch("downloaded/available mods")
			installedMods := config.ListInstalledMods(cfg)
			for _, mod := range installedMods {
				bg3DownloadedModsNode.AddMetaNode(fmt.Sprintf("%s@%s", mod.Name, mod.Version), mod.DownloadPath)
			}

			fmt.Println(rootNode.String())
		},
	}.ToCmd()
}

func makeNodeChildrenSameKeyLen(node treeprint.Tree) {

	// pass 1, fetch the max key length
	maxLen := 0
	node.VisitAll(func(n *treeprint.Node) {
		if n.Meta != nil {
			str := fmt.Sprintf("%v", n.Meta)
			if len(str) > maxLen {
				maxLen = len(str)
			}
		}
	})

	// pass 2, pad the keys
	node.VisitAll(func(n *treeprint.Node) {
		if n.Meta != nil {
			str := fmt.Sprintf("%v", n.Meta)
			str = str + strings.Repeat(" ", maxLen-len(str))
			n.Meta = str
		}
	})
}
