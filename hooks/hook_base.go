package hooks

import (
	"fmt"
	"path"
	"regexp"
)

type RunType int8

const (
	/// Run once per commit.
	runPerCommit RunType = iota
	/// Run once per file.
	runPerFile
)

const (
	//
	// Configuration keys
	//

	/// Configuration key controlling whether hook is enabled.
	keyEnabled = "enabled"
	/// Configuration key controlling the substitute command path.
	keyCommand = "cmd"

	//
	// Values
	//

	/// Value indicating boolean true
	valueTrue = "true"

	//
	// Argument placeholders
	//

	/// Placeholder for a single matching file name.
	placeholderSingleFile = "<file>"

	/// Placeholder for hook arguments, as supplied by Git.
	placeholderGitArgs = "<args>"
)

/// hookBase is a convenient do-it-all class that can be instantiated to execute tools from shell.
type hookBase struct {
	/// Unique ID of the hook. Not enforced.
	id string
	/// Human-readable name of the hook.
	name string
	/// Regexp pattern for file matching. This hook will execute only if appropriate matches are found.
	filePattern *regexp.Regexp
	/// Shell command and arguments.
	shellCommand []string
	/// Execution style, eg. once per file or once per commit.
	runType RunType
	/// Whether the hook is selected to be run.
	selected bool
	/// Whether the hook is available, eg. appropriate tools are installed. This is controlled by the user of the hook.
	available bool
	/// Related configuration section where additional metadata may be stored.
	config Config
}

/// Create a new hookBase object from the supplied pieces.
///
/// @param id The unique ID of the hook.
/// @param name Human readable name.
/// @param filePattern Regexp used for file matching. Hook will only run if a match is detected.
/// @param runType How to execute the hook.
/// @return Newly created hookBase object.
func newHookBase(id, name, filePattern string, shellCmd []string, runType RunType) *hookBase {
	hb := &hookBase{
		id:           id,
		name:         name,
		filePattern:  regexp.MustCompile(filePattern),
		available:    false,
		shellCommand: shellCmd,
		selected:     false,
		runType:      runType,
		config:       nil,
	}

	return hb
}

/// Specify the command to be run for this hook
func (h *hookBase) setShellCmd(cmd string) {
	available, command := getShellCommandAbsolutePath(cmd)
	if available {
		h.shellCommand[0] = command
		h.available = available
	}
}

/// @return An unique ID of this hook.
func (h *hookBase) ID() string {
	return h.id
}

/// @return Human readable name of the hook.
func (h *hookBase) Name() string {
	return h.name
}

/// Execute an action associated with the hook on the supplied list of files.
/// Each file is matched against the previously supplied filePattern.
/// Performs no operation if the hook is not selected, or if the corresponding command does not exist.
func (h *hookBase) Run(files []string, args []string) {
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
			fmt.Println("Running", h.name)
		} else if h.runType == runPerFile {
			fmt.Println("Running", h.name, "on", file)
		}

		sout, serr := runShellCommand(cmd)
		if len(serr) > 0 {
			fmt.Println(serr)
		}
		if len(sout) > 0 {
			fmt.Println(sout)
		}

		if h.runType == runPerCommit {
			return
		}
	}
}

/// @return Whether the hook is requested to be run.
func (h *hookBase) IsSelected() bool {
	return h.selected
}

/// @return Whether the hook can be run.
func (h *hookBase) IsAvailable() bool {
	return h.available
}

/// Modify the selected state of the hook.
///
/// @param wantSelected Whether the hook should be selected.
func (h *hookBase) SetSelected(wantSelected bool) {
	h.selected = wantSelected

	if wantSelected {
		h.config.Set(keyEnabled, valueTrue)
	} else {
		h.config.Remove(keyEnabled)
	}
}

/// Specify the configuration section responsible for managing the hook data.
///
/// @param cfg The configuration storer for the hook data.
func (h *hookBase) SetConfig(cfg Config) {
	h.config = cfg
	if cfg == nil {
		panic("No config section")
	}

	h.SetSelected(cfg.GetOrDefault(keyEnabled, "") == valueTrue)
	h.setShellCmd(cfg.GetOrDefault(keyCommand, h.shellCommand[0]))
}
