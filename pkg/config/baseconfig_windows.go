package config

import "github.com/GiGurra/boa/pkg/boa"

type BaseConfig struct {
	Verbose   boa.Required[bool]   `name:"verbose" short-name:"v" default:"false"`
	SteamPath boa.Required[string] `default:"C:\\Program Files (x86)\\Steam" name:"steam-path" short-name:"s"`
	//UserDataPath       boa.Required[string] `default:"${SteamPath}/userdata/[0]" name:"user-data-path" short-name:"u"`
	ModsInstallDir     boa.Required[string] `default:"${HOME}\\AppData\\Local\\Larian Studios\\Baldur's Gate 3\\Mods" name:"mods-install-dir" short-name:"m"`
	ModSettingsLsxPath boa.Required[string] `default:"${HOME}\\AppData\\Local\\Larian Studios\\Baldur's Gate 3\\PlayerProfiles\\Public\\modsettings.lsx" name:"mod-settings-lsx-path" short-name:"l"`
}
