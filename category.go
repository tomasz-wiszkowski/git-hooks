package main

import (
	"github.com/tomasz-wiszkowski/go-hookcfg/hooks"
)

type Category struct {
	ID    string
	Name  string
	Hooks []hooks.Hook
}

func (c *Category) readGitSettings(repo *GitRepo) {
	for _, h := range c.Hooks {
		h.SetConfig(repo.getConfigFor(c.ID, h.ID()))
	}
}
