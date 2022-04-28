package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/tomasz-wiszkowski/git-hooks/hooks"
	"github.com/tomasz-wiszkowski/git-hooks/sort"
)

type hookTreeNodeData struct {
	hook   hooks.Hook
	action hooks.Action
}

func (v *HooksTreeView) add(target *tview.TreeNode, ref *hookTreeNodeData) {
	if ref.hook == nil {
		hks := []hooks.Hook{}
		for _, c := range v.data {
			hks = append(hks, c)
		}
		sort.SortInPlaceByName(hks)

		for _, c := range hks {
			node := tview.NewTreeNode(c.Name()).SetReference(&hookTreeNodeData{c, nil}).SetSelectable(true).SetColor(tcell.ColorGrey)
			target.AddChild(node)
			v.add(node, node.GetReference().(*hookTreeNodeData))
		}
	} else if ref.action == nil {
		actions := ref.hook.Actions()
		sort.SortInPlaceByName(actions)

		for _, h := range actions {
			node := tview.NewTreeNode("").SetReference(&hookTreeNodeData{ref.hook, h}).SetSelectable(true)
			v.updateTreeNode(h, node)
			target.AddChild(node)
		}
	}
}

func (v *HooksTreeView) updateTreeNode(action hooks.Action, node *tview.TreeNode) {
	var marker rune
	if !action.IsSelected() {
		marker = ' '
	} else if !action.IsAvailable() {
		marker = '✘'
	} else /* selected and available */ {
		marker = '✔'
	}

	node.SetText(fmt.Sprintf("[%c] %s", marker, action.Name()))
}

func (v *HooksTreeView) onTreeNodeSelected(node *tview.TreeNode) {
	reference := node.GetReference().(*hookTreeNodeData)

	// Check if node or leaf. Nodes have no hook references.
	if hook := reference.action; hook == nil {
		// This is a node.
		// Collapse if visible, expand if collapsed.
		node.SetExpanded(!node.IsExpanded())
	} else {
		// This is a leaf.
		hook.SetSelected(!hook.IsSelected())
		v.updateTreeNode(hook, node)
	}
}

type HooksTreeView struct {
	*tview.TreeView
	root *tview.TreeNode
	data hooks.Hooks
}

func NewHookTreeView(data hooks.Hooks) *HooksTreeView {
	root := tview.NewTreeNode("Hooks").SetColor(tcell.ColorGrey)

	view := &HooksTreeView{
		tview.NewTreeView().SetRoot(root).SetCurrentNode(root),
		root,
		data,
	}
	view.SetSelectedFunc(view.onTreeNodeSelected)
	view.add(root, &hookTreeNodeData{nil, nil})

	return view
}
