package main

import (
	billy "github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	raw "github.com/go-git/go-git/v5/plumbing/format/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"
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
	if err != nil {
		panic(err)
	}

	c, err := r.Config()
	if err != nil {
		panic(err)
	}

	return &GitRepo{
		repo:   r,
		config: c,
	}
}

func (g *GitRepo) WorkDir() billy.Filesystem {
	wt, err := g.repo.Worktree()
	if err != nil {
		panic(err)
	}

	return wt.Filesystem
}

func (g *GitRepo) GitDir() billy.Filesystem {
	st := g.repo.Storer.(*filesystem.Storage).Filesystem()

	return st
}

func (g *GitRepo) SaveConfig() {
	err := g.repo.SetConfig(g.config)
	if err != nil {
		panic(err)
	}
}

func (g *GitRepo) GetListOfNewAndModifiedFiles() []string {
	head, err := g.repo.Head()
	if err != nil {
		panic(err)
	}

	commit, err := g.repo.CommitObject(head.Hash())
	if err != nil {
		panic(err)
	}

	tree, err := commit.Tree()
	if err != nil {
		panic(err)
	}

	iter := tree.Files()
	var paths []string
	iter.ForEach(func(file *object.File) error {
		paths = append(paths, file.Name)
		return nil
	})
	iter.Close()

	return paths
}

func (g *GitRepo) getConfigFor(categoryID, hookID string) *GitSection {
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

func (s *GitSection) Get(key string) string {
	return s.subsection.Option(key)
}

func (s *GitSection) Set(key, value string) {
	// Note: AddOption adds multiple keys with same name
	s.subsection.SetOption(key, value)
}

func (s *GitSection) Remove(key string) {
	s.subsection.RemoveOption(key)
}
