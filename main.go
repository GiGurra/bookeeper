package main

import (
	"fmt"
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/GiGurra/bookeeper/cmd"
	"github.com/GiGurra/bookeeper/pkg/config"
	"github.com/spf13/cobra"
	"strings"
)

func main() {

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

func printCommandTree(cmd *cobra.Command, level int) {
	indent := strings.Repeat("  ", level)
	fmt.Printf("%s%s - %s\n", indent, cmd.Name(), cmd.Short)
	for _, subCmd := range cmd.Commands() {
		printCommandTree(subCmd, level+1)
	}
}

func PrintCmdTreeCmd() *cobra.Command {
	return boa.Wrap{
		Use:   "print-cmd-tree",
		Short: "print the command tree",
		Run: func(cmd *cobra.Command, args []string) {
			printCommandTree(cmd.Root(), 0)
		},
	}.ToCmd()
}
