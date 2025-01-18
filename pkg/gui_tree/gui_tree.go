package gui_tree

import (
	"fmt"
	"github.com/GiGurra/bookeeper/pkg/config"
	"github.com/GiGurra/bookeeper/pkg/domain"
	"github.com/xlab/treeprint"
	"strings"
)

func Profiles(cfg *config.BaseConfig) treeprint.Tree {
	return ProfilesN(cfg, "profiles")
}

func ProfilesN(cfg *config.BaseConfig, name string) treeprint.Tree {
	profilesNode := treeprint.NewWithRoot(name)
	profiles := config.ListProfiles(cfg)
	for _, profile := range profiles {
		AddChild(profilesNode, Profile(profile))
	}
	return profilesNode
}

func Profile(
	profile config.Profile,
) treeprint.Tree {
	profileNode := treeprint.NewWithRoot(profile.Name)
	modsNode := profileNode.AddNode("mods")
	for _, mod := range profile.Mods {
		modsNode.AddMetaBranch(fmt.Sprintf("%s@%s", mod.Name, mod.Version), mod.DownloadPath)
	}
	return profileNode
}

func ConfigMod(
	mod config.Mod,
) treeprint.Tree {
	return treeprint.NewWithRoot(fmt.Sprintf("%s@%s, [%s]", mod.Name, mod.Version, mod.DownloadPath))
}

func DomainMod(
	mod domain.Mod,
) treeprint.Tree {
	return treeprint.NewWithRoot(fmt.Sprintf("%s @ %s, [%s]", mod.Name, mod.Version64, mod.UUID))
}

func ModVerbose(
	mod config.Mod,
) treeprint.Tree {
	branch := treeprint.NewWithRoot("mod")
	branch.AddMetaBranch("name", mod.Name)
	branch.AddMetaBranch("version", mod.Version)
	branch.AddMetaBranch("download path", mod.DownloadPath)
	return branch
}

func AddChild(
	atParent treeprint.Tree,
	child treeprint.Tree,
) treeprint.Tree {
	node, ok := child.(*treeprint.Node)
	if !ok {
		panic("expected treeprint.Node")
	}
	branch := atParent.AddBranch(node.Value)
	for _, child := range node.Nodes {
		AddChild(branch, child)
	}
	return branch
}

func AddChildStr(
	atParent treeprint.Tree,
	child string,
) treeprint.Tree {
	branch := atParent.AddBranch(child)
	return branch
}

func AddKV(
	atParent treeprint.Tree,
	key any,
	value any,
) treeprint.Tree {
	return atParent.AddMetaBranch(key, value)
}

func MakeChildrenSameKeyLen(node treeprint.Tree) {

	// pass 1, fetch the max key length
	maxLen := 0
	node.VisitAll(func(n *treeprint.Node) {
		if n.Meta != nil {
			str := fmt.Sprintf("%v", n.Meta)
			if len(str) > maxLen {
				maxLen = len(str)
			}
		}
	})

	// pass 2, pad the keys
	node.VisitAll(func(n *treeprint.Node) {
		if n.Meta != nil {
			str := fmt.Sprintf("%v", n.Meta)
			str = str + strings.Repeat(" ", maxLen-len(str))
			n.Meta = str
		}
	})
}
