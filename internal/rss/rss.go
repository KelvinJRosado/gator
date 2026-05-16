package rss

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var feed RSSFeed
	err = xml.Unmarshal(data, &feed)
	if err != nil {
		return nil, err
	}

	// Clean the feed
	cleanFeed(&feed)

	return &feed, nil
}

func cleanFeed(feed *RSSFeed) error {

	// Update channel data
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	// Update each item
	for _, item := range feed.Channel.Item {
		item.Description = html.UnescapeString(item.Description)
		item.Title = html.UnescapeString(item.Title)
	}

	return nil
}

func PrintFeed(feedURL string) error {
	feed, err := FetchFeed(context.Background(), feedURL)
	if err != nil {
		return err
	}

	// Format for readablity before printing
	bytes, _ := json.MarshalIndent(feed, "  ", "  ")
	fmt.Println(string(bytes))

	return nil
}
