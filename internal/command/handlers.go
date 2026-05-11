package command

import (
	"errors"
	"log/slog"
)

func HandlerLogin(s *State, cmd Command) error {
	// Check base case
	if len(cmd.args) == 0 {
		return errors.New("username not provided for login")
	}

	name := cmd.args[0]

	// Attempt to change name
	err := s.Cfg.SetUser(name)
	if err != nil {
		return err
	}

	slog.Info("Successfully logged in", "username", name)

	return nil
}
