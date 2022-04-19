package hooks

import ()

type Category struct {
	ID    string
	Name  string
	Hooks []Hook
}

func (c *Category) SetConfigStore(store ConfigStore) {
	for _, h := range c.Hooks {
		h.SetConfig(store.GetConfigFor(c.ID, h.ID()))
	}
}
