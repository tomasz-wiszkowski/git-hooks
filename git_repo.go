package main

import (
	billy "github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	raw "github.com/go-git/go-git/v5/plumbing/format/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/go-git/go-git/v5/utils/merkletrie"
	"github.com/tomasz-wiszkowski/git-hooks/hooks"
	"github.com/tomasz-wiszkowski/git-hooks/log"
)

type GitRepo struct {
	repo   *git.Repository
	config *config.Config
}

type GitSection struct {
	parent     *GitRepo
	section    string
	hook       string
	subsection *raw.Subsection
}

func GitRepoOpen() *GitRepo {
	r, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{DetectDotGit: true})
	log.Check(err, "Git: cannot open repository")

	c, err := r.Config()
	log.Check(err, "Git: cannot query repository config")

	return &GitRepo{
		repo:   r,
		config: c,
	}
}

func (g *GitRepo) WorkDir() billy.Filesystem {
	wt, err := g.repo.Worktree()
	log.Check(err, "Git: no worktree")

	return wt.Filesystem
}

func (g *GitRepo) GitDir() billy.Filesystem {
	st := g.repo.Storer.(*filesystem.Storage).Filesystem()

	return st
}

func (g *GitRepo) SaveConfig() {
	err := g.repo.SetConfig(g.config)
	log.Check(err, "Git: failed to save config")
}

/// Query the top-most commit and collect the list of modified files.
func (g *GitRepo) GetListOfNewAndModifiedFiles() []string {
	head, err := g.repo.Head()
	log.Check(err, "Git: Can't Query HEAD")

	commit, err := g.repo.CommitObject(head.Hash())
	log.Check(err, "Git: Can't Get top commit")

	parent, err := commit.Parent(0)
	log.Check(err, "Git: Can't Get parent commit")

	tree1, err := commit.Tree()
	log.Check(err, "Git: Can't Get current tree")

	tree2, err := parent.Tree()
	log.Check(err, "Git: Can't Get parent tree")

	// Make sure the order is correct - (from, to)
	changes, err := object.DiffTree(tree2, tree1)
	log.Check(err, "Git: Unable to Diff trees")

	var paths []string
	for _, c := range changes {
		action, err := c.Action()
		log.Check(err, "Git: Unable to query Action on file %s/%s", c.From.Name, c.To.Name)

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

func (g *GitRepo) GetConfigFor(categoryID, hookID string) hooks.Config {
	return &GitSection{
		parent:     g,
		section:    categoryID,
		hook:       hookID,
		subsection: g.config.Raw.Section(categoryID).Subsection(hookID),
	}
}

func (s *GitSection) Has(key string) bool {
	return s.subsection.HasOption(key)
}

func (s *GitSection) GetOrDefault(key, dflt string) string {
	if s.Has(key) {
		return s.subsection.Option(key)
	}
	return dflt
}

func (s *GitSection) Set(key, value string) {
	// Note: AddOption adds multiple keys with same name
	s.subsection.SetOption(key, value)
}

func (s *GitSection) Remove(key string) {
	s.subsection.RemoveOption(key)
}
