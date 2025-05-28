package cmd

import (
	"fmt"
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/GiGurra/bookeeper/pkg/config"
	"github.com/GiGurra/bookeeper/pkg/domain"
	"github.com/GiGurra/bookeeper/pkg/gui_tree"
	"github.com/spf13/cobra"
	"github.com/xlab/treeprint"
)

func Profiles() *cobra.Command {

	cfg := &config.BaseConfig{}

	return boa.Cmd{
		Use:         "profiles",
		Short:       "operations on profiles",
		Params:      cfg,
		ParamEnrich: boa.ParamEnricherDefault,
		SubCmds: []*cobra.Command{
			ProfilesLoadCmd(),
			ProfilesSaveCmd(),
			ProfilesDeactivateCmd(),
			ProfilesDeleteCmd(),
			ProfilesStatusCmd("status"),
			ProfilesStatusCmd("list"),
		},
	}.ToCobra()
}

type ProfilesLoadCmdConfig struct {
	Base        config.BaseConfig
	ProfileName boa.Required[string] `positional:"true" description:"profile to load"`
}

func ProfilesLoadCmd() *cobra.Command {

	cfg := &ProfilesLoadCmdConfig{}

	return boa.Cmd{
		Use:           "load",
		Short:         "load and activate a profile's mods",
		Params:        cfg,
		ParamEnrich:   boa.ParamEnricherDefault,
		ValidArgsFunc: ValidAvailableProfileNameAndVersionArgsFunc(&cfg.Base),
		RunFunc: func(cmd *cobra.Command, args []string) {
			fmt.Printf("loading profile %s\n", cfg.ProfileName.Value())
			domain.LoadProfile(&cfg.Base, cfg.ProfileName.Value())
		},
	}.ToCobra()
}

type ProfilesSaveCmdConfig struct {
	Base        config.BaseConfig
	ProfileName boa.Required[string] `positional:"true" description:"profile to save to"`
}

func ProfilesSaveCmd() *cobra.Command {

	cfg := &ProfilesSaveCmdConfig{}

	return boa.Cmd{
		Use:           "save",
		Short:         "save current active mods to profile",
		Params:        cfg,
		ParamEnrich:   boa.ParamEnricherDefault,
		ValidArgsFunc: ValidAvailableModNameAndVersionArgsFunc(&cfg.Base),
		RunFunc: func(cmd *cobra.Command, args []string) {
			fmt.Printf("saving active mods to profile %s\n", cfg.ProfileName.Value())
			domain.SaveProfile(&cfg.Base, cfg.ProfileName.Value())
		},
	}.ToCobra()
}

type ProfilesDeactivateCmdConfig struct {
	Base config.BaseConfig
}

func ProfilesDeactivateCmd() *cobra.Command {

	cfg := &ProfilesDeactivateCmdConfig{}

	return boa.Cmd{
		Use:         "deactivate-all",
		Short:       "deactivates all active mods, i.e. any profile",
		Params:      cfg,
		ParamEnrich: boa.ParamEnricherDefault,
		RunFunc: func(cmd *cobra.Command, args []string) {
			fmt.Printf("deactivating all active mods\n")
			domain.DeactivateAllMods(&cfg.Base)
		},
	}.ToCobra()
}

type ProfilesDeleteCmdConfig struct {
	Base        config.BaseConfig
	ProfileName boa.Required[string] `positional:"true" description:"profile to delete"`
}

func ProfilesDeleteCmd() *cobra.Command {

	cfg := &ProfilesDeleteCmdConfig{}

	return boa.Cmd{
		Use:           "delete",
		Short:         "delete a profile",
		Params:        cfg,
		ParamEnrich:   boa.ParamEnricherDefault,
		ValidArgsFunc: ValidAvailableProfileNameAndVersionArgsFunc(&cfg.Base),
		RunFunc: func(cmd *cobra.Command, args []string) {
			fmt.Printf("deleting profile %s\n", cfg.ProfileName.Value())
			domain.DeleteProfile(&cfg.Base, cfg.ProfileName.Value())
		},
	}.ToCobra()
}

func ProfilesStatusCmd(name string) *cobra.Command {

	cfg := &config.BaseConfig{}

	return boa.Cmd{
		Use:         name,
		Short:       "status/list of profiles",
		Params:      cfg,
		ParamEnrich: boa.ParamEnricherDefault,
		RunFunc: func(cmd *cobra.Command, args []string) {

			rootNode := treeprint.New() // NewWithRoot("Bookeeper Status")

			///////////// bookeeper paths /////////////////////////////////////////
			bookeeperPathsNode := gui_tree.AddChildStr(rootNode, "bookeeper paths")
			gui_tree.AddKV(bookeeperPathsNode, "profiles", config.ProfilesDir(cfg))
			gui_tree.MakeChildrenSameKeyLen(bookeeperPathsNode)

			///////////// Profiles /////////////////////////////////////////

			gui_tree.AddKV(rootNode, "active profile", domain.ActiveProfileName(cfg))

			availableProfilesTitle := "available profiles"
			if cfg.Verbose.Value() {
				availableProfilesTitle += " (" + config.ProfilesDir(cfg) + ")"
			}
			gui_tree.AddChild(rootNode, gui_tree.DomainProfilesN(availableProfilesTitle, domain.ListProfiles(cfg), cfg.Verbose.Value()))

			/////////////////////////////////////////////////////////////////////
			fmt.Println(rootNode.String())
		},
	}.ToCobra()
}

func ValidAvailableProfileNameAndVersionArgsFunc(cfg *config.BaseConfig) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

		profileNames := domain.ListProfileNames(cfg)

		switch len(args) {
		case 0:
			return profileNames, cobra.ShellCompDirectiveDefault | cobra.ShellCompDirectiveNoFileComp
		default:
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		}
	}
}
