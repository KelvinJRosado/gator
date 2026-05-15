package command

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/KelvinJRosado/gator/internal/database"
	"github.com/KelvinJRosado/gator/internal/rss"
	"github.com/google/uuid"
)

func HandlerLogin(s *State, cmd Command) error {
	// Check base case
	if len(cmd.Args) != 1 {
		return errors.New("username not provided for login")
	}

	name := cmd.Args[0]

	// Attempt to grab from DB
	_, err := s.Db.GetUser(context.Background(), name)
	if err != nil {
		return err
	}

	// Attempt to change name
	err = s.Cfg.SetUser(name)
	if err != nil {
		return err
	}

	slog.Info("Successfully logged in", "username", name)

	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	// Check base case
	if len(cmd.Args) != 1 {
		return errors.New("username not provided for login")
	}

	// Get name to create
	name := cmd.Args[0]

	// Create user args
	dbArgs := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}

	// Attempt to insert into DB
	_, err := s.Db.CreateUser(context.Background(), dbArgs)
	if err != nil {
		return err
	}

	// Attempt to change name
	err = s.Cfg.SetUser(name)
	if err != nil {
		return err
	}

	slog.Info("Successfully registered new user", "username", name)

	return nil
}

func HandlerReset(s *State, cmd Command) error {

	// Attempt to delete from DB
	err := s.Db.DeleteAllUsers(context.Background())
	if err != nil {
		return err
	}

	// Attempt to change name
	err = s.Cfg.SetUser("")
	if err != nil {
		return err
	}

	slog.Info("Successfully reset user table")

	return nil
}

func HandlerUsers(s *State, cmd Command) error {
	users, err := s.Db.GetAllUsers(context.Background())
	if err != nil {
		return err
	}

	if len(users) == 0 {
		fmt.Println("No users found")
	}

	for _, usr := range users {
		if s.Cfg.CurrentUserName == usr.Name {
			fmt.Printf("* %v (current)\n", usr.Name)
		} else {
			fmt.Printf("* %v\n", usr.Name)
		}
	}

	return nil
}

func HandlerAgg(s *State, cmd Command) error {

	// Using hard-coded value for testing
	feedUrl := "https://www.wagslane.dev/index.xml"

	// Print rss feed from given URL
	err := rss.PrintFeed(feedUrl)
	if err != nil {
		return err
	}

	return nil
}

func HandlerAddFeed(s *State, cmd Command) error {

	// Check base case
	if len(cmd.Args) != 2 {
		return errors.New("Feed name and URL must be given")
	}

	feedName := cmd.Args[0]
	feedUrl := cmd.Args[1]

	// Get current user
	user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
	if err != nil {
		return err
	}

	dbArgs := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedUrl,
		UserID:    user.ID,
	}

	_, err = s.Db.CreateFeed(context.Background(), dbArgs)
	if err != nil {
		return err
	}

	return nil
}

func HandlerFeeds(s *State, cmd Command) error {

	// Get feeds from db
	feeds, err := s.Db.GetAllFeeds(context.Background())
	if err != nil {
		return err
	}

	// Base case
	if len(feeds) == 0 {
		fmt.Println("No feeds currently saved")
	}

	for _, item := range feeds {
		fmt.Printf("Name: %v, URL: %v, Owner: %v\n", item.Feed.Name, item.Feed.Url, item.User.Name)
	}

	return nil
}
