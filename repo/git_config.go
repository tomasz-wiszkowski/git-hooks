package repo

import (
	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	raw "github.com/go-git/go-git/v5/plumbing/format/config"
	"github.com/tomasz-wiszkowski/git-hooks/check"
	"github.com/tomasz-wiszkowski/git-hooks/config"
)

type gitConfigManager struct {
	repo   *git.Repository
	config *gitconfig.Config
}

// Describes a configuration section (and subsection) within git config.
// Any modifications made to this section will be reflected in .git/config
// file for the current repository as a [<section> "<hook>"] entry.
type gitConfig struct {
	parent     *gitConfigManager
	section    string
	hook       string
	subsection *raw.Subsection
}

func (g *gitConfigManager) Save() {
	err := g.repo.SetConfig(g.config)
	check.Err(err, "Git: failed to save config")
}

func (g *gitConfigManager) GetConfigFor(categoryID, hookID string) config.Config {
	return &gitConfig{
		parent:     g,
		section:    categoryID,
		hook:       hookID,
		subsection: g.config.Raw.Section(categoryID).Subsection(hookID),
	}
}
func (s *gitConfig) Has(key string) bool {
	return s.subsection.HasOption(key)
}

func (s *gitConfig) GetOrDefault(key, dflt string) string {
	if s.Has(key) {
		return s.subsection.Option(key)
	}
	return dflt
}

func (s *gitConfig) Set(key, value string) {
	// Note: AddOption adds multiple keys with same name
	s.subsection.SetOption(key, value)
}

func (s *gitConfig) Remove(key string) {
	s.subsection.RemoveOption(key)
}
