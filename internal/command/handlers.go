package command

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/KelvinJRosado/gator/internal/database"
	"github.com/google/uuid"
)

func HandlerLogin(s *State, cmd Command) error {
	// Check base case
	if len(cmd.Args) == 0 {
		return errors.New("username not provided for login")
	}

	name := cmd.Args[0]

	// Attempt to change name
	err := s.Cfg.SetUser(name)
	if err != nil {
		return err
	}

	slog.Info("Successfully logged in", "username", name)

	return nil
}

func HandlerRegister (s *State, cmd Command) error {
	// Check base case
	if len(cmd.Args) == 0 {
		return errors.New("username not provided for login")
	}

	// Get name to create
	name := cmd.Args[0]

	// Create user args
	dbArgs := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: name,
	}

	// Attempt to insert into DB
	_, err := s.Db.CreateUser(context.Background(), dbArgs)
	if err != nil {
		return UserAlreadyExists
	}

	// Attempt to change name
	err = s.Cfg.SetUser(name)
	if err != nil {
		return err
	}

	slog.Info("Successfully registered new user", "username", name)

	return nil
}
