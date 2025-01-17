package cmd

import (
	"fmt"
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/GiGurra/bookeeper/pkg/config"
	"github.com/spf13/cobra"
	"log/slog"
)

func StatusCmd() *cobra.Command {

	cfg := &config.BaseConfig{}

	return boa.Wrap{
		Use:         "status",
		Short:       "print mod status",
		Params:      cfg,
		ParamEnrich: boa.ParamEnricherDefault,
		Run: func(cmd *cobra.Command, args []string) {
			// General:
			//  * bg3se status: installed | not installed
			//  * profile: default | custom1 | custom2...
			//
			// Active mods (profile: default):
			// -----------------------------
			// | Mod Name | Version | Status |
			// etc
			// ------------------------------

			// Inactive mods:
			// -----------------------------
			// | Mod Name | Version | Status |
			// etc
			// ------------------------------

			// Profiles:
			// -----------------------------
			// | Profile Name | Mods |
			// etc
			// ------------------------------

			bg3SeDllPath := config.Bg3SeDllPath(cfg)
			bg3SeInstalled := config.ExistsFile(bg3SeDllPath)

			fmt.Printf("#########################################################################################################################\n")
			fmt.Printf("################################################## General Info #########################################################\n")
			fmt.Printf("            BG3 Path: %s\n", config.Bg3Path(cfg))
			fmt.Printf("         BG3bin Path: %s\n", config.Bg3binPath(cfg))
			fmt.Printf("    Bg3UserData Path: %s\n", config.Bg3UserDataDir(cfg))
			fmt.Printf("     BG3SE installed: %v (%s)\n", bg3SeInstalled, bg3SeDllPath)
			fmt.Printf("#########################################################################################################################\n")

			slog.Error("Not yet implemented")
		},
	}.ToCmd()
}
