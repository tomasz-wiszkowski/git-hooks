package hooks

var POST_COMMIT_HOOKS = []Hook{
	// C++
	newHookBase("ClangFmt", "Clang Format", `\.(c|cc|h|hh|cpp|hpp)$`, []string{"clang-format", "-style=file", "-i", placeholderSingleFile}, runPerFile),
	newHookBase("ClangTidy", "Clang Tidy", `\.(c|cc|h|hh|cpp|hpp)$`, []string{"clang-tidy", "-format-style=file", "-i", placeholderSingleFile}, runPerFile),

	// Golang
	newHookBase("GoFmt", "Golang Format", `\.go$`, []string{"gofmt", "-w", placeholderSingleFile}, runPerFile),
	newHookBase("GoVet", "Golang Tidy", `\.go$`, []string{"go", "vet"}, runPerCommit),
	newHookBase("GoModTidy", "Golang Module Tidy", `^go.mod$`, []string{"go", "mod", "tidy"}, runPerCommit),

	// Java
	newHookBase("JavaFmt", "Java Format", `\.java$`, []string{"google-java-format", "-a", "-r", "--skip-sorting-imports", "--fix-imports-only", placeholderSingleFile}, runPerFile),

	// Python
	newHookBase("PythonFmt", "Python Format (black)", `\.py$`, []string{"black", "-q", "-t", "py310", placeholderSingleFile}, runPerFile),

	// Rust
	newHookBase("RustFmt", "Rust Format", `\.rs$`, []string{"rustfmt", placeholderSingleFile}, runPerFile),
	newHookBase("RustTidy", "Rust Tidy (Clippy)", `\.rs$`, []string{"cargo", "clippy", "--fix", placeholderSingleFile}, runPerFile),

	// Chrome
	newHookBase("ChromeClFmt", "Chrome CL Format", `.*`, []string{"git-cl", "format"}, runPerCommit),
	newHookBase("ChromeClPresubmit", "Chrome CL Presubmit", `.*`, []string{"git-cl", "presubmit"}, runPerCommit),
	newHookBase("ChromeGnDeps", "Chrome GN Deps", `^(.*\.gn[i]?|DEPS)$`, []string{"gn", "gen", "out/android.debug.arm", "--check"}, runPerCommit),
	newHookBase("ChromeJsonFmt", "Chrome JSON Format", `^fieldtrial_testing_config\.json$`, []string{"testing/variations/PRESUBMIT.py", placeholderSingleFile}, runPerFile),
	newHookBase("ChromeHistogramFmt", "Chrome Histogram Format", `^(histograms|enums)\.xml$`, []string{"tools/metrics/histograms/validate_format.py", placeholderSingleFile}, runPerFile),
}
