package command

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/KelvinJRosado/gator/internal/database"
	"github.com/KelvinJRosado/gator/internal/rss"
)

func scrapeFeeds(s *State) error {

	fmt.Printf("Attempting to get the latest feed data at %s\n", time.Now().Format(time.RFC3339))

	// Get the next feed to fetch
	feed, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	slog.Info("Fetching feed", "url", feed.Url, "name", feed.Name)

	// Fetch the feed
	feedData, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	// Mark feed as fetched
	dbArgs := database.MarkFeedFetchedParams{
		UpdatedAt: time.Now(),
		ID:        feed.ID,
	}
	_, err = s.Db.MarkFeedFetched(context.Background(), dbArgs)
	if err != nil {
		return err
	}

	// Print the title of each item in feed
	for _, item := range feedData.Channel.Item {
		fmt.Printf("* %v\n", item.Title)
	}

	return nil
}
