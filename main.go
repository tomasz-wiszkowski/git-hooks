package main

import (
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
	"github.com/tomasz-wiszkowski/git-hooks/ui"
)

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

	app := tview.NewApplication()
	tree := ui.NewHookTreeView(hooks.GetHooks())
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
