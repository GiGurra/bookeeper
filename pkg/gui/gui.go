package gui

import (
	"fmt"
	"github.com/GiGurra/bookeeper/pkg/config"
	"github.com/xlab/treeprint"
)

func Profiles(cfg *config.BaseConfig) treeprint.Tree {
	profilesNode := treeprint.NewWithRoot("profiles")
	profiles := config.ListProfiles(cfg)
	for _, profile := range profiles {
		profilesNode.AddBranch(Profile(profile))
	}
	return profilesNode
}

func Profile(
	profile config.Profile,
) treeprint.Tree {
	profileNode := treeprint.NewWithRoot(profile.Name)
	modsNode := profileNode.AddBranch("mods")
	for _, mod := range profile.Mods {
		modsNode.AddMetaBranch(fmt.Sprintf("%s@%s", mod.Name, mod.Version), mod.DownloadPath)
	}
	return profileNode
}
