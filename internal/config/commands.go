package config

import (
	"errors"
	"log/slog"
	"sync"
)

type command struct {
	name string
	args []string
}

type commands struct {
	allCommands map[string]func(*state, command) error
	mu          sync.RWMutex
}

func handlerLogin(s *state, cmd command) error {
	// Check base case
	if len(cmd.args) == 0 {
		return errors.New("username not provided for login")
	}

	name := cmd.args[0]

	// Attempt to change name
	err := s.Config.SetUser(name)
	if err != nil {
		return err
	}

	slog.Info("Successfully logged in", "username", name)

	return nil
}

func (c *commands) run(s *state, cmd command) error {
	err := c.run(s, cmd)
	if err != nil {
		return err
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.allCommands[name] = f
}
