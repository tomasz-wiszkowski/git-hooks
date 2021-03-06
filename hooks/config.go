package hooks

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/tomasz-wiszkowski/git-hooks/check"
)

const (
	configRunTypePerFile   = "perFile"
	configRunTypePerCommit = "perCommit"
)

type actionConfig struct {
	Name     string   `json:"name"`
	RunType  string   `json:"runType"`
	Priority int32    `json:"priority"`
	Pattern  string   `json:"filePattern"`
	ShellCmd []string `json:"shellCmd"`
}

type hookConfig struct {
	Name    string                   `json:"name"`
	Actions map[string]*actionConfig `json:"actions"`
}

type topConfig struct {
	Version int32                  `json:"version"`
	Hooks   map[string]*hookConfig `json:"hooks"`
}

// Load user settings from ~/.githooks.config file.
// If the file is installed and valid, returns deserialized content.
// If the file is missing or is empty, returns an empty map.
// All other cases cause assertion failure.
func loadConfigFile() map[string]Hook {
	name, err := os.UserHomeDir()
	check.Err(err, "Unable to query user home directory")

	result := map[string]Hook{}

	content, err := ioutil.ReadFile(path.Join(name, ".githooks.json"))
	if err != nil {
		return result
	}

	var config topConfig
	err = json.Unmarshal(content, &config)
	check.Err(err, "Malformed config file")

	// Assume Version 0 = no config.
	if config.Version == 0 {
		return result
	}

	check.True(config.Version == 1, "Unsupported config file version %d", config.Version)

	for ck, cv := range config.Hooks {
		hooks := []Action{}

		check.True(len(ck) > 0, "Invalid category ID")
		check.True(len(cv.Name) > 0, "Invalid category name for category %s", ck)

		category := &hook{
			id:      ck,
			name:    cv.Name,
			actions: hooks,
		}

		for hk, hv := range cv.Actions {
			runType := runPerFile
			if hv.RunType == configRunTypePerCommit {
				runType = runPerCommit
			} else if hv.RunType != configRunTypePerFile {
				check.True(false, "Invalid runType %s for hook %s", hv.RunType, hk)
			}

			check.True(len(hk) > 0, "Invalid hook ID in category %s", ck)
			check.True(len(hv.Name) > 0, "Invalid hook name for hook %s", hk)
			check.True(len(hv.ShellCmd) > 0, "Invalid shell command for hook %s", hk)

			hook := newShellAction(hk, hv.Name, hv.Priority, hv.Pattern, hv.ShellCmd, runType)
			hooks = append(hooks, hook)
		}

		category.actions = hooks
		result[ck] = category
	}

	return result
}
