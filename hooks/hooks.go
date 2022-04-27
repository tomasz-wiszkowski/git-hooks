package hooks

import "github.com/tomasz-wiszkowski/git-hooks/config"

// A map of all known and user-defined hooks and their corresponding actions.
// The key is the hook name, and the value is the corresponding Hook definition.
type Hooks map[string]Hook

var kKnownHooks Hooks = nil

// Retrieve the map of user-defined hooks.
// Upon first call the function will attempt to load user-defined hooks from
// the ~/.githooks.json config file.
func GetHooks() Hooks {
	if kKnownHooks == nil {
		kKnownHooks = loadConfigFile()
	}
	return kKnownHooks
}

// Specify the configuration store persisting action configuration relevant to
// the current context (typically the current git repository).
func (h Hooks) SetConfigStore(s config.ConfigManager) {
	for _, hook := range h {
		hook.SetConfigStore(s)
	}
}
