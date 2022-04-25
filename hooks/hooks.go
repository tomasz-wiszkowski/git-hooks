package hooks

type HookConfig map[string]Category

var kKnownHooks = HookConfig{}

func Init() {
	kKnownHooks = loadConfigFile()
}

func GetHookConfig() *HookConfig {
	return &kKnownHooks
}

func SetConfigStore(store ConfigStore) {
	for _, category := range kKnownHooks {
		category.SetConfigStore(store)
	}
}

func GetCategory(name string) (Category, bool) {
	cat, ok := kKnownHooks[name]
	return cat, ok
}
