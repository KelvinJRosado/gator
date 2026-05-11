package command

import (
	"errors"
	"sync"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	allCommands map[string]func(*State, Command) error
	mu          sync.RWMutex
}

func (c *Commands) Run(s *State, cmd Command) error {

	// Check for valid command
	f, ok := c.Get(cmd.Name)

	if !ok {
		return errors.New("command not found")
	}

	// Call handler
	return f(s, cmd)
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.allCommands[name] = f
}

func (c *Commands) Get(name string) (func(*State, Command) error, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	res, ok := c.allCommands[name]

	return res, ok
}

func CreateCommandsRegistry() *Commands {
	// Init vars
	cmds := make(map[string]func(*State, Command) error)
	res := Commands{
		allCommands: cmds,
		mu:          sync.RWMutex{},
	}

	return &res
}
