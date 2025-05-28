package cmd

import (
	"fmt"
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/GiGurra/bookeeper/pkg/config"
	"github.com/GiGurra/bookeeper/pkg/github"
	"github.com/GiGurra/bookeeper/pkg/gui_tree"
	"github.com/GiGurra/bookeeper/pkg/modzip"
	"github.com/spf13/cobra"
	"github.com/xlab/treeprint"
	"os"
	"strconv"
	"strings"
)

type Bg3SeCmdConfig struct {
	Base config.BaseConfig
}

func Bg3SeCmd() *cobra.Command {

	cfg := &config.BaseConfig{}

	return boa.Cmd{
		Use:         "bg3se",
		Short:       "operations related to bg3se (BG3 Script Extender)",
		Params:      cfg,
		ParamEnrich: boa.ParamEnricherDefault,
		SubCmds: []*cobra.Command{
			Bg3SeStatusCmd(),
			Bg3SeInstallCmd(),
		},
	}.ToCobra()
}

func Bg3SeInstallCmd() *cobra.Command {

	cfg := &Bg3SeCmdConfig{}

	return boa.Cmd{
		Use:         "install",
		Short:       "download (from github) and install the latest version of bg3se",
		Params:      cfg,
		ParamEnrich: boa.ParamEnricherDefault,
		RunFunc: func(cmd *cobra.Command, args []string) {

			//Norbyte/bg3se
			fmt.Println("Checking for latest bg3se release at github.com/Norbyte/bg3se/releases/latest")
			release := github.GetLatestRelease("Norbyte", "bg3se")
			if len(release.Assets) == 0 {
				panic("no assets found at github.com/Norbyte/bg3se/releases/latest")
			}
			if len(release.Assets) > 1 {
				panic("multiple assets found at github.com/Norbyte/bg3se/releases/latest. This is not supported (yet?)")
			}
			asset := release.Assets[0]
			// download the asset
			tempDir, err := os.MkdirTemp("", "bg3se")
			if err != nil {
				panic(fmt.Errorf("failed to create temp dir: %w", err))
			}
			defer func() { _ = os.RemoveAll(tempDir) }()

			downloadedZipFilePath := asset.DownloadToDir(tempDir)
			defer func() { _ = os.Remove(downloadedZipFilePath) }()
			if !strings.HasSuffix(strings.ToLower(downloadedZipFilePath), ".zip") {
				panic(fmt.Errorf("downloaded bg3se asset is not a zip file: %s", downloadedZipFilePath))
			}

			// extract the asset on top of our current bg3se installation/bin
			modzip.ExtractSpecificFilesFromZip(downloadedZipFilePath, []string{"DWrite.dll"}, config.Bg3binPath(&cfg.Base))

			fmt.Println("bg3se installed successfully to " + config.Bg3binPath(&cfg.Base))
			fmt.Println("!! Remember to set the command line args in steam for bg3:")
			cmdStr := "%command%"
			//goland:noinspection GoPrintFunctions
			fmt.Println(" WINEDLLOVERRIDES=\"DWrite.dll=n,b\" " + cmdStr)
		},
	}.ToCobra()
}

func Bg3SeStatusCmd() *cobra.Command {

	cfg := &Bg3SeCmdConfig{}

	return boa.Cmd{
		Use:         "status",
		Short:       "show status of bg3se (BG3 Script Extender)",
		Params:      cfg,
		ParamEnrich: boa.ParamEnricherDefault,
		RunFunc: func(cmd *cobra.Command, args []string) {

			bg3SeDllPath := config.Bg3SeDllPath(&cfg.Base)
			bg3SeInstalled := config.ExistsFile(bg3SeDllPath)

			rootNode := treeprint.New() // NewWithRoot("Bookeeper Status")

			///////////// bg3 paths /////////////////////////////////////////
			bg3PathsNode := gui_tree.AddChildStr(rootNode, "bg3 paths")

			installNode := gui_tree.AddKV(bg3PathsNode, "bg3 install dir", config.Bg3Path(&cfg.Base))

			bg3SeNode := gui_tree.AddChildStr(installNode, "bg3se status")
			gui_tree.AddKV(bg3SeNode, "installed", strconv.FormatBool(bg3SeInstalled))
			gui_tree.AddKV(bg3SeNode, "dll path", bg3SeDllPath)
			gui_tree.MakeChildrenSameKeyLen(bg3SeNode)

			/////////////////////////////////////////////////////////////////////
			fmt.Println(rootNode.String())
		},
	}.ToCobra()
}
