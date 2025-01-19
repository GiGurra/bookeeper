package main

import (
	"fmt"
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/GiGurra/bookeeper/cmd"
	"github.com/GiGurra/bookeeper/pkg/config"
	"github.com/GiGurra/bookeeper/pkg/gui_tree"
	"github.com/spf13/cobra"
	"github.com/xlab/treeprint"
)

func main() {

	cobra.EnableCommandSorting = false

	cfg := config.BaseConfig{}

	boa.Wrap{
		Use:   "bookeeper",
		Short: "Very basic cli mod manager for Baldur's Gate 3",
		SubCommands: []*cobra.Command{
			cmd.StatusCmd(),
			cmd.GetCmd(),
			cmd.ModsCmd(),
			cmd.Profiles(),
			cmd.Bg3SeCmd(),
			PrintCmdTreeCmd(),
		},
		Params:      &cfg,
		ParamEnrich: boa.ParamEnricherDefault,
	}.ToApp()
}

func makeCmdTree(cmd *cobra.Command, level int) treeprint.Tree {
	tree := &treeprint.Node{
		Meta:  cmd.Name(),
		Value: cmd.Short,
	}
	for _, subCmd := range cmd.Commands() {
		gui_tree.AddChild(tree, makeCmdTree(subCmd, level+1))
	}
	// only if we are 2 levels down
	if level >= 1 {
		gui_tree.MakeChildrenSameKeyLen(tree)
	}
	return tree
}

func PrintCmdTreeCmd() *cobra.Command {
	return boa.Wrap{
		Use:   "print-cmd-tree",
		Short: "print the command tree",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(makeCmdTree(cmd.Root(), 0).String())
		},
	}.ToCmd()
}
