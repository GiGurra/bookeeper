package cmd

import (
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/GiGurra/bookeeper/pkg/config"
	"github.com/GiGurra/bookeeper/pkg/domain"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"strings"
)

func ModsCmd() *cobra.Command {

	cfg := &config.BaseConfig{}

	return boa.Wrap{
		Use:         "mods",
		Short:       "operations on mods",
		Params:      cfg,
		ParamEnrich: boa.ParamEnricherDefault,
		SubCommands: []*cobra.Command{
			ModsMakeAvailableCmd(),
			ModsMakeUnavailableCmd(),
		},
	}.ToCmd()
}

type ModsMakeAvailableCmdConfig struct {
	Base       config.BaseConfig
	ModZipPath boa.Required[string] `positional:"true" description:"mod zip file path"`
}

func ModsMakeAvailableCmd() *cobra.Command {

	cfg := &ModsMakeAvailableCmdConfig{}

	return boa.Wrap{
		Use:         "make-available",
		Short:       "make a new mod available",
		Params:      cfg,
		ParamEnrich: boa.ParamEnricherDefault,
		Run: func(cmd *cobra.Command, args []string) {
			domain.MakeModAvailable(&cfg.Base, cfg.ModZipPath.Value())
		},
	}.ToCmd()
}

type ModsMakeUnavailableCmdConfig struct {
	Base       config.BaseConfig
	ModName    boa.Required[string] `positional:"true" description:"mod name"`
	ModVersion boa.Required[string] `positional:"true" description:"mod version"`
}

func ModsMakeUnavailableCmd() *cobra.Command {

	cfg := &ModsMakeUnavailableCmdConfig{}

	return boa.Wrap{
		Use:         "make-unavailable",
		Short:       "make a mod unavailable",
		Params:      cfg,
		ParamEnrich: boa.ParamEnricherDefault,
		ValidArgsFunc: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

			modsByName := lo.GroupBy(domain.ListAvailableMods(&cfg.Base), func(item domain.Mod) string {
				return item.Name
			})

			switch len(args) {
			case 0:
				return lo.Map(lo.Keys(modsByName), func(item string, _ int) string {
					// replace spaces with \<space>
					return strings.ReplaceAll(item, " ", "\\ ")
				}), cobra.ShellCompDirectiveDefault | cobra.ShellCompDirectiveNoFileComp
			case 1:
				return lo.Map(modsByName[args[0]], func(item domain.Mod, _ int) string {
					return item.Version64
				}), cobra.ShellCompDirectiveDefault | cobra.ShellCompDirectiveNoFileComp
			}

			return []string{}, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			domain.MakeModUnavailable(&cfg.Base, cfg.ModName.Value(), cfg.ModVersion.Value())
		},
	}.ToCmd()
}
