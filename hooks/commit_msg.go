package hooks

var COMMIT_MSG_HOOKS = []Hook{
	// C++
	newHookBase("ChangeId", "Git Add ChangeID", `.`, []string{"git-add-changeid", placeholderGitArgs}, runPerCommit),
	newHookBase("ReflowMsg", "Reflow Git Commit message", `.`, []string{"fmt", "-g", "70", placeholderGitArgs}, runPerCommit),
}
