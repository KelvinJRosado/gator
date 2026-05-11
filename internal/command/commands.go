package command

import (
	"sync"
)

type Command struct {
	name string
	args []string
}

type Commands struct {
	allCommands map[string]func(*State, Command) error
	mu          *sync.RWMutex
}

func (c *Commands) run(s *State, cmd Command) error {
	err := c.run(s, cmd)
	if err != nil {
		return err
	}

	return nil
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.allCommands[name] = f
}

func CreateCommandsRegistry() *Commands {
	// Init vars
	mu := sync.RWMutex{}
	cmds := make(map[string]func(*State, Command) error)
	res := Commands{
		allCommands: cmds,
		mu:          &mu,
	}

	return &res
}
