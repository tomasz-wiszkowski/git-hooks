package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/tomasz-wiszkowski/git-hooks/hooks"
	"github.com/tomasz-wiszkowski/git-hooks/log"
	"github.com/tomasz-wiszkowski/git-hooks/sort"
)

type reference struct {
	hook   hooks.Hook
	action hooks.Action
}

func add(target *tview.TreeNode, ref *reference) {
	if ref.hook == nil {
		hks := []hooks.Hook{}
		for _, c := range hooks.GetHookConfig() {
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

func openRepo() *GitRepo {
	repo := GitRepoOpen()
	hooks.SetConfigStore(repo)
	return repo
}

func main() {
	hooks.Init()

	selfName := path.Base(os.Args[0])

	if h, ok := hooks.GetHook(selfName); ok {
		runHooks(h, os.Args[1:])
	} else if len(os.Args) == 1 {
		showConfig()
	} else if h, ok := hooks.GetHook(os.Args[1]); ok {
		runHooks(h, os.Args[2:])
	} else if os.Args[1] == "install" {
		install()
	} else {
		fmt.Println("Unknown hook type", os.Args[1])
	}
}

func runHooks(hook hooks.Hook, args []string) {
	repo := openRepo()
	files := repo.GetListOfNewAndModifiedFiles()

	// Used by hooks install, file fixing and others
	err := os.Chdir(repo.WorkDir().Root())
	log.Check(err, "Run: cannot open work directory")

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
	log.Check(err, "Run: terminated abnormally")

	repo.SaveConfig()
}

func install() {
	selfAbsolutePath, err := filepath.Abs(os.Args[0])
	log.Check(err, "Install: cannot locate self")

	repo := openRepo()
	gitDir := repo.GitDir()

	err = gitDir.MkdirAll("hooks", 0755)
	log.Check(err, "Install: failed to create hooks directory")

	hookDir, err := gitDir.Chroot("hooks")
	log.Check(err, "Install: failed to navigate to hooks directory")

	for _, hook := range hooks.GetHookConfig() {
		fmt.Println("Installing", hook.ID(), "in", hookDir.Root(), "pointing to", selfAbsolutePath)
		if _, err = hookDir.Stat(hook.ID()); err == nil {
			err = hookDir.Remove(hook.ID())
			if err != nil && err != os.ErrNotExist {
				log.Check(err, "Install: failed to remove hook %s", hook)
			}
		}

		err = os.Symlink(selfAbsolutePath, hookDir.Join(hookDir.Root(), hook.ID()))
		log.Check(err, "Install: failed to install hook %s", hook)
	}

	showConfig()
}
