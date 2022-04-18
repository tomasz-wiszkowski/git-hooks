package hooks

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

type HookBase struct {
	id           string
	name         string
	filePattern  *regexp.Regexp
	available    bool
	shellCommand []string
	selState     SelectedState
	runType      RunType
	config       Config
}

func newHookBase(id, name, filePattern string, runType RunType) *HookBase {
	return &HookBase{
		id:           id,
		name:         name,
		filePattern:  regexp.MustCompile(filePattern),
		available:    false,
		shellCommand: []string{},
		selState:     SelectedStateUnknown,
		runType:      runType,
		config:       nil,
	}
}

func (h *HookBase) setAvailable(available bool) {
	h.available = available
}

func (h *HookBase) setCommand(command []string) {
	h.shellCommand = command
}

func (h *HookBase) RunType() RunType {
	return h.runType
}

func (h *HookBase) ID() string {
	return h.id
}

func (h *HookBase) Name() string {
	return h.name
}

func (h *HookBase) Run(files []string) {
	if !h.IsSelected() {
		return
	}
	if !h.IsAvailable() {
		fmt.Println("Cannot run", h.Name(), "- missing command.")
		return
	}

	for _, file := range files {
		if !h.filePattern.MatchString(file) {
			continue
		}

		if h.RunType() == RunPerCommit {
			fmt.Println("Running", h.name)
			h.runCommand(h.shellCommand)
			return
		} else if h.RunType() == RunPerFile {
			fmt.Println("Running", h.name, "on", file)
			h.runCommand(append(h.shellCommand, file))
		}
	}
}

func (h *HookBase) IsSelected() bool {
	return h.selState == SelectedStateEnabled || h.selState == SelectedStateUnavailable
}

func (h *HookBase) IsAvailable() bool {
	return h.available
}

func (h *HookBase) SetSelected(wantSelected bool) {
	if wantSelected {
		if h.IsAvailable() {
			h.selState = SelectedStateEnabled
		} else {
			h.selState = SelectedStateUnavailable
		}
		h.config.Set("enabled", "true")
	} else {
		h.selState = SelectedStateDisabled
		h.config.Set("enabled", "false")
	}
}

func (h *HookBase) State() SelectedState {
	return h.selState
}

func (h *HookBase) SetConfig(cfg Config) {
	h.config = cfg
	if cfg == nil {
		panic("No config section")
	}
	if !cfg.Has("enabled") {
		h.selState = SelectedStateUnknown
		return
	}
	h.SetSelected(cfg.Get("enabled") == "true")
}

func (h *HookBase) checkAvailable() bool {
	fmt.Println("ffffff")
	return true
}

func (h *HookBase) getExecutablePath(executableName string) *string {
	path, err := exec.LookPath(executableName)
	if err != nil {
		return nil
	}
	return &path
}

func (h *HookBase) runCommand(args []string) (stdout, stderr string) {
	cmd := exec.Command(args[0], args[1:]...)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	outStr := strings.TrimSpace(outb.String())
	if len(outStr) > 0 {
		fmt.Println(h.Name(), "output:", outStr)
	}
	errStr := strings.TrimSpace(errb.String())
	if len(errStr) > 0 {
		fmt.Println(h.Name(), "error:", errStr)
	}

	return outStr, errStr
}
