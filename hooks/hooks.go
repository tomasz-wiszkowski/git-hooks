package hooks

type HookConfig map[string]Hook

var kKnownHooks = HookConfig{}

func Init() {
	kKnownHooks = loadConfigFile()
}

func GetHookConfig() HookConfig {
	return kKnownHooks
}

func SetConfigStore(store ConfigStore) {
	for _, hook := range kKnownHooks {
		hook.SetConfigStore(store)
	}
}

func GetHook(name string) (Hook, bool) {
	cat, ok := kKnownHooks[name]
	return cat, ok
}
