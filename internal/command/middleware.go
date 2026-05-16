package command

import (
	"context"

	"github.com/KelvinJRosado/gator/internal/database"
)

// Checks if the user is logged in before continuing
func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {

	return func(s *State, cmd Command) error {
		// Get current user
		user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}

}
