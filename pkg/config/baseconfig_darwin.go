package config

import "github.com/GiGurra/boa/pkg/boa"

type BaseConfig struct {
	Verbose   boa.Required[bool]   `name:"verbose" short-name:"v" default:"false"`
	SteamPath boa.Required[string] `default:"${HOME}/.steam/steam" name:"steam-path" short-name:"s"`
	//UserDataPath       boa.Required[string] `default:"${SteamPath}/userdata/[0]" name:"user-data-path" short-name:"u"`
	ModsInstallDir     boa.Required[string] `default:"${SteamPath}/steamapps/compatdata/1086940/pfx/drive_c/users/steamuser/AppData/Local/Larian Studios/Baldur's Gate 3/Mods" name:"mods-install-dir" short-name:"m"`
	ModSettingsLsxPath boa.Required[string] `default:"${SteamPath}/steamapps/compatdata/1086940/pfx/drive_c/users/steamuser/AppData/Local/Larian Studios/Baldur's Gate 3/PlayerProfiles/Public/modsettings.lsx" name:"mod-settings-lsx-path" short-name:"l"`
}
