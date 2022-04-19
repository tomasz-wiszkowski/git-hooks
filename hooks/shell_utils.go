package hooks

import (
	"bytes"
	"os"
	"os/exec"
	"path"
	"strings"
)

/// Resolve executable name into full path.
///
/// @param executableName Name of the executable command to locate.
/// @return When command is found (either in PATH or location relative to current path), returns true and absolute path of the command.
///         Otherwise returns false and a placeholder command.
func getShellCommandAbsolutePath(executableName string) (bool, string) {
	absPath, err := exec.LookPath(executableName)
	if err == nil {
		return true, absPath
	}

	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	absPath = path.Join(workDir, executableName)
	_, err = os.Stat(absPath)
	if err == nil {
		return true, absPath
	}

	return false, "false"
}

/// Execute supplied shell command.
/// The command must be supplied in an "exploded" form, where each argument is a separate string.
///
/// @return Strings encompassing stdout and stderr.
func runShellCommand(args []string) (stdout, stderr string) {
	cmd := exec.Command(args[0], args[1:]...)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	outStr := strings.TrimSpace(outb.String())
	errStr := strings.TrimSpace(errb.String())
	return outStr, errStr
}