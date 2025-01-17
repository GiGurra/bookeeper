package cmd

import (
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/spf13/cobra"
	"log/slog"
)

func StatusCmd() *cobra.Command {
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
	return boa.Wrap{
		Use:   "status",
		Short: "print mod status",
		Run: func(cmd *cobra.Command, args []string) {
			slog.Error("Not yet implemented")
		},
	}.ToCmd()
}
