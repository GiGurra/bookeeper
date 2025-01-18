package main

import (
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/GiGurra/bookeeper/cmd"
	"github.com/GiGurra/bookeeper/pkg/config"
	"github.com/spf13/cobra"
)

func main() {

	cfg := config.BaseConfig{}

	boa.Wrap{
		Use:   "bookeeper",
		Short: "Very basic cli mod manager for Baldur's Gate 3",
		SubCommands: []*cobra.Command{
			cmd.StatusCmd(),
			cmd.GetCmd(),
		},
		Params:      &cfg,
		ParamEnrich: boa.ParamEnricherDefault,
	}.ToApp()
}
