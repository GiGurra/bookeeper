package cmd

import (
	"fmt"
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/GiGurra/bookeeper/pkg/config"
	"github.com/GiGurra/bookeeper/pkg/domain"
	"github.com/GiGurra/bookeeper/pkg/gui_tree"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/xlab/treeprint"
)

type GetCmdConfig struct {
	Base config.BaseConfig
}

func GetCmd() *cobra.Command {

	cfg := &GetCmdConfig{}

	type SubCommand struct {
		Name    string
		Handler func(*config.BaseConfig) string
	}
	individualCommands := []SubCommand{
		{"bookeeper-dir", config.BooKeeperDir},
		{"bookeeper-downloaded-mods-dir", config.DownloadedModsDir},
		{"bookeeper-profiles-dir", config.ProfilesDir},
		{"bg3-dir", config.Bg3Path},
		{"bg3-modsettings-path", config.Bg3ModsettingsFilePath},
		{"bg3-bin-dir", config.Bg3binPath},
		{"bg3-mod-dir", config.Bg3ModInstallDir},
		{"active-profile", domain.ActiveProfileName},
	}

	commands := append(individualCommands, SubCommand{"all", func(cfg *config.BaseConfig) string {
		rootNode := treeprint.New()
		for _, command := range individualCommands {
			gui_tree.AddKV(rootNode, command.Name, command.Handler(cfg))
		}
		gui_tree.MakeChildrenSameKeyLen(rootNode)
		return rootNode.String()
	}})

	return boa.Wrap{
		Use:         "get",
		Short:       "get specific path or info",
		Params:      cfg,
		ParamEnrich: boa.ParamEnricherDefault,
		SubCommands: lo.Map(commands, func(command SubCommand, _ int) *cobra.Command {

			cfg := &GetCmdConfig{}

			return boa.Wrap{
				Use:         command.Name,
				Short:       "get value of " + command.Name,
				Params:      cfg,
				ParamEnrich: boa.ParamEnricherDefault,
				Run: func(cmd *cobra.Command, args []string) {
					result := command.Handler(&cfg.Base)
					fmt.Println(result)
				},
			}.ToCmd()
		}),
	}.ToCmd()
}
