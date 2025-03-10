package gui_tree

import (
	"fmt"
	"github.com/GiGurra/bookeeper/pkg/domain"
	"github.com/xlab/treeprint"
	"strings"
)

func DomainProfilesN(name string, profiles []domain.Profile, verbose bool) treeprint.Tree {
	profilesNode := treeprint.NewWithRoot(name)
	for _, profile := range profiles {
		AddChild(profilesNode, DomainProfile(profile, verbose))
	}
	return profilesNode
}

func DomainProfile(
	profile domain.Profile,
	verbose bool,
) treeprint.Tree {
	profileNode := treeprint.NewWithRoot(profile.Name)
	for _, mod := range profile.Mods {
		AddChild(profileNode, DomainMod(mod, verbose))
	}
	MakeChildrenSameKeyLen(profileNode)
	return profileNode
}

func DomainMod(
	mod domain.Mod,
	verbose bool,
) treeprint.Tree {
	if verbose {
		branch := treeprint.NewWithRoot(mod.Name)
		branch.AddMetaBranch("Folder", mod.Folder)
		branch.AddMetaBranch("MD5", mod.MD5)
		branch.AddMetaBranch("Name", mod.Name)
		branch.AddMetaBranch("PublishHandle", mod.PublishHandle)
		branch.AddMetaBranch("UUID", mod.UUID)
		branch.AddMetaBranch("Version64", mod.Version64)
		MakeChildrenSameKeyLen(branch)
		return branch
	} else {
		return &treeprint.Node{
			Root:  nil,
			Meta:  mod.Name,
			Value: fmt.Sprintf("%s, v %s", mod.UUID, mod.Version64),
			Nodes: nil,
		}
	}
}

func AddChild(
	atParent treeprint.Tree,
	child treeprint.Tree,
) treeprint.Tree {
	node, ok := child.(*treeprint.Node)
	if !ok {
		panic("expected treeprint.Node")
	}
	branch := func() treeprint.Tree {
		if node.Meta != nil {
			return atParent.AddMetaBranch(node.Meta, node.Value)
		} else {
			return atParent.AddBranch(node.Value)
		}
	}()
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
	value string,
) treeprint.Tree {
	return atParent.AddMetaBranch(key, value)
}

func MakeChildrenSameKeyLen(node treeprint.Tree) {

	// pass 1, fetch the max key length
	maxLen := 0
	node.VisitAll(func(n *treeprint.Node) {
		if n.Root != node {
			return // dont go deeper than the first level
		}
		if n.Meta != nil {
			str := fmt.Sprintf("%v", n.Meta)
			if len(str) > maxLen {
				maxLen = len(str)
			}
		}
	})

	// pass 2, pad the keys
	node.VisitAll(func(n *treeprint.Node) {
		if n.Root != node {
			return // dont go deeper than the first level
		}
		if n.Meta != nil {
			str := fmt.Sprintf("%v", n.Meta)
			str = str + strings.Repeat(" ", maxLen-len(str))
			n.Meta = str
		}
	})
}
