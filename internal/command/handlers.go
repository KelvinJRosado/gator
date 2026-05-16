package command

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/KelvinJRosado/gator/internal/database"
	"github.com/google/uuid"
)

// Login as existing user
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

// Register new user
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

	slog.Info("Successfully registered new user", "username", name, "id", dbArgs.ID)

	return nil
}

// Delete all users
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

// Print all users
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

// Test RSS call
func HandlerAgg(s *State, cmd Command) error {

	// Base case
	if len(cmd.Args) != 1 {
		return errors.New("Fetch frequency must be given")
	}

	// Parse frequency
	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Collecting feeds every %v\n", timeBetweenReqs)

	// Create waitgroup to avoid early exit
	var wg sync.WaitGroup

	// Create ticker and start scraping
	ticker := time.NewTicker(timeBetweenReqs)
	wg.Go(func() {

		// Initial call
		err := scrapeFeeds(s)
		if err != nil {
			slog.Error("Error scraping feeds", "error", err)
		}

		// Continuous call
		for range ticker.C {
			err := scrapeFeeds(s)
			if err != nil {
				slog.Error("Error scraping feeds", "error", err)
			}
		}
	})

	wg.Wait()

	return nil
}

// Save new feed record
func HandlerAddFeed(s *State, cmd Command, user database.User) error {

	// Check base case
	if len(cmd.Args) != 2 {
		return errors.New("Feed name and URL must be given")
	}

	feedName := cmd.Args[0]
	feedUrl := cmd.Args[1]

	dbArgs := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedUrl,
		UserID:    user.ID,
	}

	newFeed, err := s.Db.CreateFeed(context.Background(), dbArgs)
	if err != nil {
		return err
	}

	slog.Info("Successfully added new feed", "feedUrl", feedUrl, "id", dbArgs.ID)

	dbArgs2 := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    newFeed.ID,
	}

	_, err = s.Db.CreateFeedFollow(context.Background(), dbArgs2)
	if err != nil {
		return err
	}

	slog.Info("Successfully followed new feed", "feedUrl", feedUrl, "id", dbArgs2.ID)

	return nil
}

// Get details for all feeds
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
		fmt.Printf("Name: %v, URL: %v, Owner: %v\n", item.Name, item.Url, item.UserName)
	}

	return nil
}

// Follow an existing feed
func HandlerFollow(s *State, cmd Command, user database.User) error {

	// Base case
	if len(cmd.Args) != 1 {
		return errors.New("Feed URL must be provided")
	}

	// Get existing feed
	feed, err := s.Db.GetFeedByUrl(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}

	// Create new follow record
	dbArgs := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	_, err = s.Db.CreateFeedFollow(context.Background(), dbArgs)
	if err != nil {
		return err
	}

	slog.Info("Successfully followed feed", "username", user.Name, "feedUrl", feed.Url, "id", dbArgs.ID)

	return nil
}

// Stop following an existing feed
func HandlerUnfollow(s *State, cmd Command, user database.User) error {

	// Base case
	if len(cmd.Args) != 1 {
		return errors.New("Feed URL must be provided")
	}

	feedFollows, err := s.Db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return err
	}

	unfollowed := false

	for _, feedFollow := range feedFollows {

		// Skip until we get to specified feed
		if feedFollow.FeedUrl != cmd.Args[0] {
			continue
		}

		// If we get here, this is a valid case and we can delete from table
		dbArgs := database.DeleteFeedFollowsForUserParams{
			UserID: user.ID,
			FeedID: feedFollow.FeedID,
		}
		err := s.Db.DeleteFeedFollowsForUser(context.Background(), dbArgs)
		if err != nil {
			return err
		}
		unfollowed = true

	}

	if unfollowed {
		slog.Info("Successfully unfollowed feed", "username", user.Name, "feedUrl", cmd.Args[0])
		return nil
	} else {
		return errors.New("Feed not followed by current user")
	}

}

// Get details for all feeds followed by current user
func HandlerFollowing(s *State, cmd Command, user database.User) error {

	feedFollows, err := s.Db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return err
	}

	// Base case
	if len(feedFollows) == 0 {
		fmt.Println("No feeds currently followed by logged in user")
	}

	for _, item := range feedFollows {
		fmt.Printf("* %v\n", item.FeedName)
	}

	return nil
}
