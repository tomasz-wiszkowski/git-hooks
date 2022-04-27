package hooks

import "github.com/tomasz-wiszkowski/git-hooks/config"

type Action interface {
	ID() string
	Name() string
	Priority() int32
	SetSelected(bool)
	IsSelected() bool
	IsAvailable() bool
	SetConfig(config.Config)
	Run(file []string, args []string)
}
