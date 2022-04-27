package repo

import (
	billy "github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/go-git/go-git/v5/utils/merkletrie"
	"github.com/tomasz-wiszkowski/git-hooks/config"
	"github.com/tomasz-wiszkowski/git-hooks/try"
)

type gitRepo struct {
	repo   *git.Repository
	config *gitconfig.Config
}

func gitRepoOpen() Repo {
	r, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{DetectDotGit: true})
	try.CheckErr(err, "Git: cannot open repository")

	c, err := r.Config()
	try.CheckErr(err, "Git: cannot query repository config")

	return &gitRepo{
		repo:   r,
		config: c,
	}
}

func (g *gitRepo) WorkDir() billy.Filesystem {
	wt, err := g.repo.Worktree()
	try.CheckErr(err, "Git: no worktree")

	return wt.Filesystem
}

func (g *gitRepo) ConfigDir() billy.Filesystem {
	st := g.repo.Storer.(*filesystem.Storage).Filesystem()
	return st
}

func (g *gitRepo) GetConfigManager() config.ConfigManager {
	c, err := g.repo.Config()
	try.CheckErr(err, "Git: cannot query repository config")

	return &gitConfigManager{
		repo:   g.repo,
		config: c,
	}
}

/// Query the top-most commit and collect the list of modified files.
func (g *gitRepo) GetListOfNewAndModifiedFiles() []string {
	head, err := g.repo.Head()
	try.CheckErr(err, "Git: Can't Query HEAD")

	commit, err := g.repo.CommitObject(head.Hash())
	try.CheckErr(err, "Git: Can't Get top commit")

	parent, err := commit.Parent(0)
	try.CheckErr(err, "Git: Can't Get parent commit")

	tree1, err := commit.Tree()
	try.CheckErr(err, "Git: Can't Get current tree")

	tree2, err := parent.Tree()
	try.CheckErr(err, "Git: Can't Get parent tree")

	// Make sure the order is correct - (from, to)
	changes, err := object.DiffTree(tree2, tree1)
	try.CheckErr(err, "Git: Unable to Diff trees")

	var paths []string
	for _, c := range changes {
		action, err := c.Action()
		try.CheckErr(err, "Git: Unable to query Action on file %s/%s", c.From.Name, c.To.Name)

		switch action {
		case merkletrie.Delete:
			continue
		case merkletrie.Insert:
			fallthrough
		case merkletrie.Modify:
			paths = append(paths, c.To.Name)
		}
	}

	return paths
}