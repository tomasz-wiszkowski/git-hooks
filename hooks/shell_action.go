package hooks

import (
	"fmt"
	"log"
	"path"
	"regexp"

	"github.com/tomasz-wiszkowski/git-hooks/check"
	"github.com/tomasz-wiszkowski/git-hooks/config"
)

type RunType int8

const (
	// Run once per commit.
	runPerCommit RunType = iota
	// Run once per file.
	runPerFile
)

const (
	// Configuration key controlling whether hook is enabled.
	keyEnabled = "enabled"
	// Configuration key controlling the substitute command path.
	keyCommand = "cmd"

	// Value indicating boolean true
	valueTrue = "true"

	// Placeholder for a single matching file name.
	placeholderSingleFile = "<file>"

	// Placeholder for hook arguments, as supplied by Git.
	placeholderGitArgs = "<args>"
)

// shellAction is a convenient do-it-all class that can be instantiated to execute tools from shell.
type shellAction struct {
	// Unique ID of the hook. Not enforced.
	id string
	// Human-readable name of the hook.
	name string
	// Execution prioirty.
	priority int32
	// Regexp pattern for file matching. This hook will execute only if appropriate matches are found.
	filePattern *regexp.Regexp
	// Shell command and arguments.
	shellCommand []string
	// Execution style, eg. once per file or once per commit.
	runType RunType
	// Whether the hook is selected to be run.
	selected bool
	// Whether the hook is available, eg. appropriate tools are installed. This is controlled by the user of the hook.
	available bool
	// Related configuration section where additional metadata may be stored.
	config config.Config
}

// Create a new shellAction object from the supplied pieces.
func newShellAction(id, name string, priority int32, filePattern string, shellCmd []string, runType RunType) *shellAction {
	hb := &shellAction{
		id:           id,
		name:         name,
		priority:     priority,
		filePattern:  regexp.MustCompile(filePattern),
		available:    false,
		shellCommand: shellCmd,
		selected:     false,
		runType:      runType,
		config:       nil,
	}

	return hb
}

// Specify the command to be run for this hook
func (h *shellAction) setShellCmd(cmd string) {
	available, command := getShellCommandAbsolutePath(cmd)
	if available {
		h.shellCommand[0] = command
		h.available = available
	}
}

// Return the unique ID of this hook.
func (h *shellAction) ID() string {
	return h.id
}

// Return the human readable name of the hook.
func (h *shellAction) Name() string {
	return h.name
}

// Return priority of the hook. Lower number = higher priority.
func (h *shellAction) Priority() int32 {
	return h.priority
}

// Execute an action associated with the hook on the supplied list of files.
// Each file is matched against the previously supplied filePattern.
// Performs no operation if the hook is not selected, or if the corresponding command does not exist.
func (h *shellAction) Run(files []string, args []string) {
	if !h.IsSelected() {
		return
	}
	if !h.IsAvailable() {
		fmt.Println("Cannot run", h.Name(), "- missing command", h.shellCommand[0])
		return
	}

	substitutions := map[string]interface{}{
		placeholderGitArgs: args,
	}

	for _, file := range files {
		base := path.Base(file)

		if !h.filePattern.MatchString(base) {
			continue
		}

		substitutions[placeholderSingleFile] = file
		cmd := substituteCommandLine(h.shellCommand, substitutions)

		if h.runType == runPerCommit {
			log.Println("Running", h.name)
		} else if h.runType == runPerFile {
			log.Println("Running", h.name, "on", file)
		}

		runShellCommand(cmd)
		if h.runType == runPerCommit {
			return
		}
	}
}

// Return whether the hook is requested to be run.
func (h *shellAction) IsSelected() bool {
	return h.selected
}

// Return whether the hook can be run.
func (h *shellAction) IsAvailable() bool {
	return h.available
}

// Modify the selected state of the hook.
func (h *shellAction) SetSelected(wantSelected bool) {
	h.selected = wantSelected

	if wantSelected {
		h.config.Set(keyEnabled, valueTrue)
	} else {
		h.config.Remove(keyEnabled)
	}
}

// Specify the configuration section responsible for managing the hook data.
func (h *shellAction) SetConfig(cfg config.Config) {
	h.config = cfg
	check.True(cfg != nil, "No config section")

	h.SetSelected(cfg.GetOrDefault(keyEnabled, "") == valueTrue)
	h.setShellCmd(cfg.GetOrDefault(keyCommand, h.shellCommand[0]))
}
