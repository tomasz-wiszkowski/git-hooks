package hooks

type Hook interface {
	ID() string
	Name() string
	SetSelected(bool)
	IsSelected() bool
	IsAvailable() bool
	SetConfig(Config)
	Run(file []string)
}

type ConfigStore interface {
	GetConfigFor(section, subsection string) Config
}

type Config interface {
	Set(key, value string)
	Has(key string) bool
	GetOrDefault(key, dflt string) string
	Remove(key string)
}
