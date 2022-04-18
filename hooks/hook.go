package hooks

type SelectedState int8
type RunType int8

const (
	/// Not configured in config file.
	SelectedStateUnknown SelectedState = iota
	/// Not available due to missing dependency.
	SelectedStateUnavailable
	/// Disabled in config file.
	SelectedStateDisabled
	/// Enabled in config file.
	SelectedStateEnabled
)

const (
	/// Run once per commit.
	RunPerCommit = iota
	/// Run once per file.
	RunPerFile
)

type Hook interface {
	ID() string
	Name() string
	State() SelectedState
	RunType() RunType
	SetSelected(bool)
	IsSelected() bool
	IsAvailable() bool
	SetConfig(Config)
	Run(file []string)
}

type Config interface {
	Set(key, value string)
	Has(key string) bool
	Get(key string) string
}
