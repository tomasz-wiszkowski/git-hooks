package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Config struct {
	PostCommit []*Hook `json:"postCommit"`
}

const (
	CategoryRoot = iota
	CategoryPostCommit
)

type Reference struct {
	category int
	hook *Hook
}

var config Config

func add(target *tview.TreeNode, ref *Reference) {
	if ref.category == CategoryRoot {
		node := tview.NewTreeNode("Post-commit").SetReference(&Reference{CategoryPostCommit, nil}).SetSelectable(true).SetColor(tcell.ColorGrey)
		target.AddChild(node)
		add(node, node.GetReference().(*Reference))
	}

	if ref.category == CategoryPostCommit && ref.hook == nil {
		for _, h := range config.PostCommit {
			node := tview.NewTreeNode("").SetReference(&Reference{CategoryPostCommit, h}).SetSelectable(true)
			h.updateNode(node)
			target.AddChild(node)
		}
	}
}

func onTreeNodeSelected(node *tview.TreeNode) {
	reference := node.GetReference().(*Reference)

	// Check if node or leaf. Nodes have no hook references.
	if hook := reference.hook; hook == nil {
		// This is a node.
		// Collapse if visible, expand if collapsed.
		node.SetExpanded(!node.IsExpanded())
	} else {
		// This is a leaf.
		hook.toggleSelected()
		hook.updateNode(node)
	}
}

func main() {
	var err error
	var configFile []byte

	if configFile, err = ioutil.ReadFile("hooks.json"); err != nil {
		panic(err)
	}

	if err = json.Unmarshal(configFile, &config); err != nil {
		panic(err)
	}

	root := tview.NewTreeNode("Hooks").SetColor(tcell.ColorGrey)
	tree := tview.NewTreeView().SetRoot(root).SetCurrentNode(root)
	tree.SetSelectedFunc(onTreeNodeSelected)
	add(root, &Reference{CategoryRoot, nil})

	if err := tview.NewApplication().SetRoot(tree, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
