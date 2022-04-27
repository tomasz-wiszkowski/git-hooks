package repo

import (
	"log"

	billy "github.com/go-git/go-billy/v5"
	"github.com/tomasz-wiszkowski/git-hooks/config"
)

// Simple repository abstraction.
type Repo interface {
	// Return absolute path to repository working directory.
	WorkDir() billy.Filesystem
	// Return absolute path to repository configuration directory.
	ConfigDir() billy.Filesystem
	// Return a list of all modified and added files, relative to
	// the working directory root.
	GetListOfNewAndModifiedFiles() []string
	// Create (if required) and return the configuration manager that
	// can be used to persist configuration for the current repo.
	GetConfigManager() config.ConfigManager
}

// Attempt to identify and open repository under current path.
func OpenRepo() Repo {
	if maybeGit := gitRepoOpen(); maybeGit != nil {
		return maybeGit
	}
	log.Fatalf("Unable to determine repository type")
	return nil
}
