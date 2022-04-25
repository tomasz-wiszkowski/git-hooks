package hooks

type Hook interface {
	ID() string
	Name() string
	Actions() []Action
	SetConfigStore(ConfigStore)
}

type hook struct {
	id      string
	name    string
	actions []Action
}

func (c *hook) ID() string {
	return c.id
}

func (c *hook) Name() string {
	return c.name
}

func (c *hook) Actions() []Action {
	return c.actions
}

func (c *hook) SetConfigStore(store ConfigStore) {
	for _, h := range c.Actions() {
		h.SetConfig(store.GetConfigFor(c.ID(), h.ID()))
	}
}
