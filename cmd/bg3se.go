package cmd

import (
	"fmt"
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/GiGurra/bookeeper/pkg/config"
	"github.com/GiGurra/bookeeper/pkg/github"
	"github.com/spf13/cobra"
	"os"
)

type Bg3SeCmdConfig struct {
	Base config.BaseConfig
}

func Bg3SeCmd() *cobra.Command {

	cfg := &config.BaseConfig{}

	return boa.Wrap{
		Use:         "bg3se",
		Short:       "operations related to bg3se (BG3 Script Extender)",
		Params:      cfg,
		ParamEnrich: boa.ParamEnricherDefault,
		SubCommands: []*cobra.Command{
			Bg3SeInstallCmd(),
		},
	}.ToCmd()
}

func Bg3SeInstallCmd() *cobra.Command {

	cfg := &Bg3SeCmdConfig{}

	return boa.Wrap{
		Use:         "install",
		Short:       "download (from github) and install the latest version of bg3se",
		Params:      cfg,
		ParamEnrich: boa.ParamEnricherDefault,
		Run: func(cmd *cobra.Command, args []string) {
			//Norbyte/bg3se
			release := github.GetLatestRelease("Norbyte", "bg3se")
			if len(release.Assets) == 0 {
				panic("no assets found at github.com/Norbyte/bg3se/releases/latest")
			}
			if len(release.Assets) > 1 {
				panic("multiple assets found at github.com/Norbyte/bg3se/releases/latest. This is not supported (yet?)")
			}
			asset := release.Assets[0]
			// download the asset
			tempDir, err := os.CreateTemp("", "bg3se")
			if err != nil {
				panic(fmt.Errorf("failed to create temp dir: %w", err))
			}
			//defer func() { _ = os.RemoveAll(tempDir.Name()) }()

			asset.DownloadToDir(tempDir.Name())
		},
	}.ToCmd()
}
