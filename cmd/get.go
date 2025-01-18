package cmd

import (
	"fmt"
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/GiGurra/bookeeper/pkg/config"
	"github.com/spf13/cobra"
)

type GetCmdConfig struct {
	Base  config.BaseConfig
	Value boa.Required[string] `positional:"true"`
}

func GetCmd() *cobra.Command {

	cfg := &GetCmdConfig{}

	handlers := map[string]func(*config.BaseConfig) string{
		"bookeeper-dir":                 config.BooKeeperDir,
		"bookeeper-downloaded-mods-dir": config.DownloadedModsDir,
		"bookeeper-profiles-dir":        config.ProfilesDir,
		"bg3-dir":                       config.Bg3Path,
		"bg3-bin-dir":                   config.Bg3binPath,
		"bg3-userdata-dir":              config.Bg3UserDataDir,
		"bg3-mod-dir":                   config.Bg3ModInstallDir,
		"bg3-userdata-profile-modsettings-xml-path": config.Bg3UserdataProfileModsettingsFilePath,
	}

	validArgs := func() []string {
		args := make([]string, 0, len(handlers))
		for k := range handlers {
			args = append(args, k)
		}
		return args
	}()

	return boa.Wrap{
		Use:         "get",
		Short:       "get specific path or info",
		ValidArgs:   validArgs,
		Params:      cfg,
		ParamEnrich: boa.ParamEnricherDefault,
		Run: func(cmd *cobra.Command, args []string) {
			result := handlers[cfg.Value.Value()](&cfg.Base)
			fmt.Println(result)
		},
	}.ToCmd()
}
