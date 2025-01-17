package main

import (
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/GiGurra/bookeeper/cmd"
	"github.com/spf13/cobra"
)

type Config struct {
	SomeParam boa.Required[string] `default:"default" name:"some-param" short-name:"s"`
}

func main() {

	cfg := Config{}

	boa.Wrap{
		Use:   "bookeeper",
		Short: "Very basic cli mod manager for Baldur's Gate 3",
		SubCommands: []*cobra.Command{
			cmd.StatusCmd(),
		},
		Params:      &cfg,
		ParamEnrich: boa.ParamEnricherDefault,
	}.ToApp()
}
