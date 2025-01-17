package cmd

import (
	"fmt"
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/GiGurra/bookeeper/pkg/config"
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
			// General:
			//  * bg3se status: installed | not installed
			//  * profile: default | custom1 | custom2...
			//
			// Active mods (profile: default):
			// -----------------------------
			// | Mod Name | Version | Status |
			// etc
			// ------------------------------

			// Inactive mods:
			// -----------------------------
			// | Mod Name | Version | Status |
			// etc
			// ------------------------------

			// Profiles:
			// -----------------------------
			// | Profile Name | Mods |
			// etc
			// ------------------------------

			bg3SeDllPath := config.Bg3SeDllPath(cfg)
			bg3SeInstalled := config.ExistsFile(bg3SeDllPath)

			rootNode := treeprint.NewWithRoot("Bookeeper Status")

			bg3PathsNode := rootNode.AddBranch("BG3 Paths")
			bg3PathsNode.AddMetaNode("BG3 Path", config.Bg3Path(cfg))
			bg3PathsNode.AddMetaNode("BG3bin Path", config.Bg3binPath(cfg))
			bg3PathsNode.AddMetaNode("Bg3UserData Path", config.Bg3UserDataDir(cfg))
			makeNodeChildrenSameKeyLen(bg3PathsNode)

			bg3SeNode := rootNode.AddBranch("BG3SE Status")
			bg3SeNode.AddMetaNode("BG3SE installed", bg3SeInstalled)
			bg3SeNode.AddMetaNode("BG3SE DLL Path", bg3SeDllPath)
			makeNodeChildrenSameKeyLen(bg3SeNode)

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
