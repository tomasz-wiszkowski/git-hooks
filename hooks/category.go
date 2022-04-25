package hooks

type Category interface {
	ID() string
	Name() string
	Hooks() []Hook
	SetConfigStore(ConfigStore)
}

type category struct {
	id    string
	name  string
	hooks []Hook
}

func (c *category) ID() string {
	return c.id
}

func (c *category) Name() string {
	return c.name
}

func (c *category) Hooks() []Hook {
	return c.hooks
}

func (c *category) SetConfigStore(store ConfigStore) {
	for _, h := range c.Hooks() {
		h.SetConfig(store.GetConfigFor(c.ID(), h.ID()))
	}
}
