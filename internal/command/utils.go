package command

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/KelvinJRosado/gator/internal/database"
	"github.com/KelvinJRosado/gator/internal/rss"
	"github.com/google/uuid"
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

	// Save each post to DB
	for _, item := range feedData.Channel.Item {

		parsedTime, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {
			return err
		}

		dbArgs2 := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       sql.NullString{String: item.Title, Valid: len(item.Title) > 0},
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: len(item.Description) > 0},
			PublishedAt: parsedTime,
			FeedID:      feed.ID,
		}
		_, err = s.Db.CreatePost(context.Background(), dbArgs2)
		if err != nil {
			return err
		}
	}

	return nil
}
