package hooks

import ()

//
// Golang
//
func GoFmt() Hook {
	res := newHookBase("GoFmt", "Golang Format", `\.go$`, RunPerFile)
	path := res.getExecutablePath("gofmt")
	if path != nil {
		res.setAvailable(true)
		res.setCommand([]string{*path, "-w"})
	}
	return res
}

func GoTidy() Hook {
	res := newHookBase("GoVet", "Golang Tidy", `\.go$`, RunPerCommit)
	path := res.getExecutablePath("go")
	if path != nil {
		res.setAvailable(true)
		res.setCommand([]string{*path, "vet"})
	}
	return res
}

//
// Chrome-specific
//
func ChromeClFmt() Hook {
	res := newHookBase("ChromeClFmt", "Chrome CL Format", `.*`, RunPerCommit)
	path := res.getExecutablePath("git-cl")
	if path == nil {
		path = res.getExecutablePath("third_party/depot_tools/git-cl")
	}
	if path != nil {
		res.setAvailable(true)
		res.setCommand([]string{*path, "format"})
	}
	return res
}

func ChromeClPresubmit() Hook {
	res := newHookBase("ChromeClPresubmit", "Chrome CL Presubmit", `.*`, RunPerCommit)
	path := res.getExecutablePath("git-cl")
	if path == nil {
		path = res.getExecutablePath("third_party/depot_tools/git-cl")
	}
	if path != nil {
		res.setAvailable(true)
		res.setCommand([]string{*path, "presubmit"})
	}
	return res
}

func ChromeGnDeps() Hook {
	res := newHookBase("ChromeGnDeps", "Chrome GN Deps", `^(.*\.gn[i]?|DEPS)$`, RunPerCommit)
	path := res.getExecutablePath("gn")
	if path != nil {
		res.setAvailable(true)
		res.setCommand([]string{*path, "gen", "out/android.debug.arm", "--check"})
	}
	return res
}

func ChromeJsonFmt() Hook {
	res := newHookBase("ChromeJsonFmt", "Chrome JSON Format", `\.json$`, RunPerFile)
	path := res.getExecutablePath("testing/variations/PRESUBMIT.py")
	if path != nil {
		res.setAvailable(true)
		res.setCommand([]string{*path})
	}
	return res
}

func ChromeHistogramFmt() Hook {
	res := newHookBase("ChromeHistogramFmt", "Chrome Histogram Format", `^(histograms|enums)\.xml$`, RunPerFile)
	path := res.getExecutablePath("tools/metrics/histograms/validate_format.py")
	if path != nil {
		res.setAvailable(true)
		res.setCommand([]string{*path})
	}
	return res
}

//
// C++ specific
//
func CppFmt() Hook {
	res := newHookBase("ClangFmt", "Clang Format", `\.(c|cc|h|hh|cpp|hpp)$`, RunPerFile)
	path := res.getExecutablePath("clang-format")
	if path != nil {
		res.setAvailable(true)
		res.setCommand([]string{*path, "-style=file", "-i"})
	}
	return res
}

func CppTidy() Hook {
	res := newHookBase("ClangTidy", "Clang Tidy", `\.(c|cc|h|hh|cpp|hpp)$`, RunPerFile)
	path := res.getExecutablePath("clang-tidy")
	if path != nil {
		res.setAvailable(true)
		res.setCommand([]string{*path, "-format-style=file", "-i"})
	}
	return res
}

//
// Java specific
//
func JavaFmt() Hook {
	res := newHookBase("JavaFmt", "Java Format", `\.java$`, RunPerFile)
	path := res.getExecutablePath("google-java-format")
	if path != nil {
		res.setAvailable(true)
		res.setCommand([]string{*path, "-a", "-r", "--skip-sorting-imports", "--fix-imports-only"})
	}
	return res
}

//
// Rust specific
//
func RustFmt() Hook {
	res := newHookBase("RustFmt", "Rust Format", `\.rs$`, RunPerFile)
	path := res.getExecutablePath("rustfmt")
	if path != nil {
		res.setAvailable(true)
		res.setCommand([]string{*path})
	}
	return res
}

func RustTidy() Hook {
	res := newHookBase("RustTidy", "Rust Tidy (Clippy)", `\.rs$`, RunPerFile)
	path := res.getExecutablePath("cargo")
	if path != nil {
		res.setAvailable(true)
		res.setCommand([]string{*path, "clippy", "--fix"})
	}
	return res
}

//
// Python specific
//
func PythonFmt() Hook {
	res := newHookBase("PythonFmt", "Python Format (black)", `\.py$`, RunPerFile)
	path := res.getExecutablePath("black")
	if path != nil {
		res.setAvailable(true)
		res.setCommand([]string{*path, "-q", "-t", "py310"})
	}
	return res
}
