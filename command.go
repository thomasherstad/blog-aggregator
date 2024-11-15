package main

import "errors"

type command struct {
	name string
	args []string
}

type commands struct {
	actions map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.actions[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.actions[cmd.name]
	if !ok {
		return errors.New("command does not exist")
	}

	err := f(s, cmd)
	if err != nil {
		return err
	}

	return nil
}
