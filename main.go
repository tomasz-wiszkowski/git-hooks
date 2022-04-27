package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/rivo/tview"
	"github.com/tomasz-wiszkowski/git-hooks/hooks"
	"github.com/tomasz-wiszkowski/git-hooks/repo"
	"github.com/tomasz-wiszkowski/git-hooks/sort"
	"github.com/tomasz-wiszkowski/git-hooks/try"
)

type reference struct {
	hook   hooks.Hook
	action hooks.Action
}

func add(target *tview.TreeNode, ref *reference) {
	if ref.hook == nil {
		hks := []hooks.Hook{}
		for _, c := range hooks.GetHooks() {
			hks = append(hks, c)
		}
		sort.SortInPlaceByName(hks)

		for _, c := range hks {
			node := tview.NewTreeNode(c.Name()).SetReference(&reference{c, nil}).SetSelectable(true).SetColor(tcell.ColorGrey)
			target.AddChild(node)
			add(node, node.GetReference().(*reference))
		}
	} else if ref.action == nil {
		actions := ref.hook.Actions()
		sort.SortInPlaceByName(actions)

		for _, h := range actions {
			node := tview.NewTreeNode("").SetReference(&reference{ref.hook, h}).SetSelectable(true)
			updateTreeNode(h, node)
			target.AddChild(node)
		}
	}
}

func updateTreeNode(action hooks.Action, node *tview.TreeNode) {
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

func onTreeNodeSelected(node *tview.TreeNode) {
	reference := node.GetReference().(*reference)

	// Check if node or leaf. Nodes have no hook references.
	if hook := reference.action; hook == nil {
		// This is a node.
		// Collapse if visible, expand if collapsed.
		node.SetExpanded(!node.IsExpanded())
	} else {
		// This is a leaf.
		hook.SetSelected(!hook.IsSelected())
		updateTreeNode(hook, node)
	}
}

func openRepo() repo.Repo {
	r := repo.OpenRepo()
	hooks.GetHooks().SetConfigStore(r.GetConfigManager())
	return r
}

func main() {
	log.Default().SetFlags(log.Ltime | log.Lshortfile)

	selfName := path.Base(os.Args[0])
	hks := hooks.GetHooks()

	if h, ok := hks[selfName]; ok {
		runHooks(h, os.Args[1:])
	} else if len(os.Args) == 1 {
		showConfig()
	} else if h, ok := hks[os.Args[1]]; ok {
		runHooks(h, os.Args[2:])
	} else if os.Args[1] == "install" {
		install()
	} else {
		log.Fatalln("Unknown hook type", os.Args[1])
	}
}

func runHooks(hook hooks.Hook, args []string) {
	repo := openRepo()
	files := repo.GetListOfNewAndModifiedFiles()

	// Used by hooks install, file fixing and others
	err := os.Chdir(repo.WorkDir().Root())
	try.CheckErr(err, "Run: cannot open work directory")

	actions := hook.Actions()
	sort.SortInPlaceByPriority(actions)
	for _, h := range actions {
		h.Run(files, args)
	}
}

func showConfig() {
	repo := openRepo()
	root := tview.NewTreeNode("Hooks").SetColor(tcell.ColorGrey)
	tree := tview.NewTreeView().SetRoot(root).SetCurrentNode(root)
	tree.SetSelectedFunc(onTreeNodeSelected)
	add(root, &reference{nil, nil})

	app := tview.NewApplication()
	app.SetRoot(tree, true).EnableMouse(true)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			app.Stop()
		}
		return event
	})

	err := app.Run()
	try.CheckErr(err, "Run: terminated abnormally")
	repo.GetConfigManager().Save()
}

func install() {
	selfAbsolutePath, err := filepath.Abs(os.Args[0])
	try.CheckErr(err, "Install: cannot locate self")

	repo := openRepo()
	configDir := repo.ConfigDir()

	err = configDir.MkdirAll("hooks", 0755)
	try.CheckErr(err, "Install: failed to create hooks directory")

	hookDir, err := configDir.Chroot("hooks")
	try.CheckErr(err, "Install: failed to navigate to hooks directory")

	for _, hook := range hooks.GetHooks() {
		log.Println("Installing", hook.ID(), "in", hookDir.Root(), "pointing to", selfAbsolutePath)
		if _, err = hookDir.Lstat(hook.ID()); err == nil {
			err = hookDir.Remove(hook.ID())
			if err != nil && err != os.ErrNotExist {
				try.CheckErr(err, "Install: failed to remove hook %s", hook.Name())
			}
		}

		err = osfs.Default.Symlink(selfAbsolutePath, hookDir.Join(hookDir.Root(), hook.ID()))
		try.CheckErr(err, "Install: failed to install hook %s", hook.Name())
	}

	showConfig()
}
