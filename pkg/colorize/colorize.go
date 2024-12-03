package colorize

type Colorize struct {
	message   string
	modifiers []string
}

func New(message string) *Colorize {
	return &Colorize{
		message:   message,
		modifiers: make([]string, 0),
	}
}

func (c *Colorize) Modify(modifier string) *Colorize {
	c.modifiers = append(c.modifiers, modifier)
	return c
}

func (c *Colorize) Commit() string {
	m := ""
	for i := range c.modifiers {
		m = m + c.modifiers[i]
	}

	return m + c.message + Reset
}
