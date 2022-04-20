package hooks

import (
	"bytes"
	"fmt"
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
	outStr := strings.TrimSpace(outb.String())
	errStr := strings.TrimSpace(errb.String())
	if err != nil {
		fmt.Println("*** Command failed ***")
		fmt.Println(outStr)
		fmt.Println(errStr)
	}

	return outStr, errStr
}

/// Substitute arguments and construct a command line
///
/// @param inputCmdLine Input set of arguments (including command)
/// @param substituteArgs Map of substitute arguments. Key is the placeholder value, and Value is the replacement args.
/// @reutrn Resulting CommandLine.
func substituteCommandLine(inputCmdLine []string, substituteArgs map[string]interface{}) []string {
	out := []string{}

	for _, arg := range inputCmdLine {
		if sub, ok := substituteArgs[arg]; ok {
			if str, ok := sub.(string); ok {
				out = append(out, str)
			} else if strArr, ok := sub.([]string); ok {
				out = append(out, strArr...)
			} else {
				panic(fmt.Sprintf("Unknown value type for key %s", arg))
			}
		} else {
			out = append(out, arg)
		}
	}
	return out
}
