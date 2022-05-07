package hooks

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/tomasz-wiszkowski/git-hooks/check"
)

// Resolve executable name into full path.
// Searches for executableName either in PATH or location relative to
// current path. When found, returns true and absolute path of the command.
// Otherwise returns false and a placeholder command.
func getShellCommandAbsolutePath(executableName string) (bool, string) {
	absPath, err := exec.LookPath(executableName)
	if err == nil {
		return true, absPath
	}

	workDir, err := os.Getwd()
	check.Err(err, "Resolve: cannot determine current dir")

	absPath = path.Join(workDir, executableName)
	_, err = os.Stat(absPath)
	if err == nil {
		return true, absPath
	}

	return false, "false"
}

// Execute supplied shell command.
// The command must be supplied in an "exploded" form, where each argument is a
// separate string. Returns a pair of strings: stdout and stderr.
func runShellCommand(args []string) (stdout, stderr string) {
	cmd := exec.Command(args[0], args[1:]...)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	outStr := strings.TrimSpace(outb.String())
	errStr := strings.TrimSpace(errb.String())
	if err != nil {
		log.Printf("Command failed")
		log.Println(outStr)
		log.Println(errStr)
	}

	return outStr, errStr
}

// Substitute arguments and construct a command line.
// The inputCmdLine represents the set of arguments (including command).
// Each argument is matched against items in substituteArgs map. When a
// corresponding entry is found, the argument is the replacement with the
// map value. Replaces single strings and string arrays.
// Returns resulting CommandLine.
func substituteCommandLine(inputCmdLine []string, substituteArgs map[string]interface{}) []string {
	out := []string{}

	for _, arg := range inputCmdLine {
		if sub, ok := substituteArgs[arg]; ok {
			if str, ok := sub.(string); ok {
				out = append(out, str)
			} else if strArr, ok := sub.([]string); ok {
				out = append(out, strArr...)
			} else {
				log.Panicf("Unknown value type for key %s", arg)
			}
		} else {
			out = append(out, arg)
		}
	}
	return out
}
