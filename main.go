package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/tomasz-wiszkowski/go-hookcfg/hooks"
	"github.com/tomasz-wiszkowski/go-hookcfg/log"
)

type Reference struct {
	category *hooks.Category
	hook     hooks.Hook
}

func add(target *tview.TreeNode, ref *Reference) {
	if ref.category == nil {
		for _, c := range *hooks.GetHookConfig() {
			node := tview.NewTreeNode(c.Name).SetReference(&Reference{c, nil}).SetSelectable(true).SetColor(tcell.ColorGrey)
			target.AddChild(node)
			add(node, node.GetReference().(*Reference))
		}
	} else if ref.hook == nil {
		for _, h := range ref.category.Hooks {
			node := tview.NewTreeNode("").SetReference(&Reference{ref.category, h}).SetSelectable(true)
			updateTreeNode(h, node)
			target.AddChild(node)
		}
	}
}

func updateTreeNode(hook hooks.Hook, node *tview.TreeNode) {
	var marker rune
	if !hook.IsSelected() {
		marker = ' '
	} else if !hook.IsAvailable() {
		marker = '✘'
	} else /* selected and available */ {
		marker = '✔'
	}

	node.SetText(fmt.Sprintf("[%c] %s", marker, hook.Name()))
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

	if h, ok := hooks.GetCategory(selfName); ok {
		runHooks(h, os.Args[1:])
	} else if len(os.Args) == 1 {
		showConfig()
	} else if h, ok := hooks.GetCategory(os.Args[1]); ok {
		runHooks(h, os.Args[2:])
	} else if os.Args[1] == "install" {
		install()
	} else {
		fmt.Println("Unknown hook type", os.Args[1])
	}
}

func runHooks(category *hooks.Category, args []string) {
	repo := openRepo()
	files := repo.GetListOfNewAndModifiedFiles()

	// Used by hooks install, file fixing and others
	err := os.Chdir(repo.WorkDir().Root())
	log.Check(err, "Run: cannot open work directory")

	for _, h := range category.Hooks {
		h.Run(files, args)
	}
}

func showConfig() {
	repo := openRepo()
	root := tview.NewTreeNode("Hooks").SetColor(tcell.ColorGrey)
	tree := tview.NewTreeView().SetRoot(root).SetCurrentNode(root)
	tree.SetSelectedFunc(onTreeNodeSelected)
	add(root, &Reference{nil, nil})

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

	for _, category := range *hooks.GetHookConfig() {
		fmt.Println("Installing", category.ID, "in", hookDir.Root(), "pointing to", selfAbsolutePath)
		if _, err = hookDir.Stat(category.ID); err == nil {
			err = hookDir.Remove(category.ID)
			if err != nil && err != os.ErrNotExist {
				log.Check(err, "Install: failed to remove hook %s", category)
			}
		}

		err = os.Symlink(selfAbsolutePath, hookDir.Join(hookDir.Root(), category.ID))
		log.Check(err, "Install: failed to install hook %s", category)
	}

	showConfig()
}
