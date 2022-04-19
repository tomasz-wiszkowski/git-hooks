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

type Config interface {
	Set(key, value string)
	Has(key string) bool
	GetOrDefault(key, dflt string) string
	Remove(key string)
}
