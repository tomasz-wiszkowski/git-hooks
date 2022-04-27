package config

// Abstraction of a configuration manager, allowing access to isolated
// configuration sections.
type ConfigManager interface {
	// Create or access the configuration file for specific section
	// and subsection.
	GetConfigFor(section, subsection string) Config
	// Save the configuration: persist all Config items on disk.
	Save()
}

// Abstraction of a configuration store - simple key/value map where all keys
// and all values are plain strings.
type Config interface {
	// Set value for supplied key.
	Set(key, value string)
	// Check whether key has an associated value.
	Has(key string) bool
	// Retrieve the value for key, substituting dflt if no value is set.
	GetOrDefault(key, dflt string) string
	// Remove value for specified key.
	Remove(key string)
}
