package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/jwoodsiii/blogator/internal/database"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func timeToNullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: t, Valid: true}
}

func scrapeFeeds(s *state) error {
	// Get the next feed to fetch from the DB.
	// Mark it as fetched.
	// Fetch the feed using the URL (we already wrote this function)
	// Iterate over the items in the feed and print their titles to the console.
	// Update the agg command to now take a single argument: time_between_reqs.

	ctx := context.Background()
	feed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("Error getting next feed to fetch: %v", err)
	}
	log.Println("Found a feed to fetch!")
	scrapeFeed(s.db, feed)

	return nil
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	if err := db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{LastFetchedAt: timeToNullTime(time.Now().UTC()), ID: feed.ID}); err != nil {
		log.Printf("Error marking feed as fetched: %v", err)
		return
	}

	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Error fetching feed at url: %v", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		fmt.Printf("Title: %s", item.Title)
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}

func fetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	if feedUrl == "" {
		return &RSSFeed{}, fmt.Errorf("Need url to fetch")
	}
	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("Error generating request: %v", err)
	}

	req.Header.Set("User-Agent", "blogator")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return &RSSFeed{}, fmt.Errorf("Error executing request: %v", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("Error reading data from request: %v", err)
	}

	var feed RSSFeed
	if err := xml.Unmarshal(data, &feed); err != nil {
		return &RSSFeed{}, fmt.Errorf("Error unmarshaling xml data: %v", err)
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for idx, i := range feed.Channel.Item {
		feed.Channel.Item[idx].Title = html.UnescapeString(i.Title)
		feed.Channel.Item[idx].Description = html.UnescapeString(i.Description)
	}
	return &feed, nil
}
