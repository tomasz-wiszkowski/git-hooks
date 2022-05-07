package ui

import (
	"fmt"
	"sort"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/tomasz-wiszkowski/git-hooks/hooks"
)

type hookTreeNodeData struct {
	hook   hooks.Hook
	action hooks.Action
}

// Append individual hook nodes to the root node.
func (v *HooksTreeView) addHookTreeNodes(target *tview.TreeNode) {
	hks := []hooks.Hook{}
	for _, c := range v.data {
		hks = append(hks, c)
	}
	sort.Slice(hks, func(a, b int) bool { return hks[a].Name() < hks[b].Name() })

	for _, c := range hks {
		node := tview.NewTreeNode(c.Name()).SetReference(&hookTreeNodeData{c, nil}).SetSelectable(true).SetColor(tcell.ColorGrey)
		target.AddChild(node)
		v.add(node, node.GetReference().(*hookTreeNodeData))
	}
}

// Append individual action nodes to the hook node.
func (v *HooksTreeView) addActionTreeNodes(target *tview.TreeNode, ref *hookTreeNodeData) {
	actions := ref.hook.Actions()
	sort.Slice(actions, func(a, b int) bool { return actions[a].Name() < actions[b].Name() })

	for _, h := range actions {
		node := tview.NewTreeNode("").SetReference(&hookTreeNodeData{ref.hook, h}).SetSelectable(true)
		v.updateTreeNode(h, node)
		target.AddChild(node)
	}
}

// Determine the type of the node being accessed by the user and populate its
// children.
func (v *HooksTreeView) add(target *tview.TreeNode, ref *hookTreeNodeData) {
	if ref.hook == nil {
		v.addHookTreeNodes(target)
	} else if ref.action == nil {
		v.addActionTreeNodes(target, ref)
	}
}

// Update the tree node's display text.
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

// Respond to user selection. Toggle expanded state of nodes, and
// toggle selected state of leaves.
func (v *HooksTreeView) onTreeNodeSelected(node *tview.TreeNode) {
	reference := node.GetReference().(*hookTreeNodeData)

	// Check if node or leaf. Nodes have no action references.
	if action := reference.action; action == nil {
		// Node (ie. not an action): toggle expanded state.
		node.SetExpanded(!node.IsExpanded())
	} else {
		// Leaf (ie. action): toggle enabled state.
		action.SetSelected(!action.IsSelected())
		v.updateTreeNode(action, node)
	}
}

// TUI TreeView for Hooks and Actions.
type HooksTreeView struct {
	*tview.TreeView
	root *tview.TreeNode
	data hooks.Hooks
}

// Instantiate a new HooksTreeView TUI element. The element is by default popuated
// with all known hooks and actions.
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
